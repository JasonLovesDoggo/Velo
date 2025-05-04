package manager

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/jasonlovesdoggo/velo/internal/config"
	"github.com/jasonlovesdoggo/velo/internal/log"
	"github.com/jasonlovesdoggo/velo/pkg/core/node"
)

// SwarmManager handles Docker Swarm cluster management operations
type SwarmManager struct {
	client        *client.Client
	nodeCache     map[string]node.Info
	nodeCacheMu   sync.RWMutex
	refreshTicker *time.Ticker
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewSwarmManager creates a new SwarmManager instance
func NewSwarmManager() (*SwarmManager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	manager := &SwarmManager{
		client:    cli,
		nodeCache: make(map[string]node.Info),
		ctx:       ctx,
		cancel:    cancel,
	}

	return manager, nil
}

// Start begins the manager's background operations
func (m *SwarmManager) Start() error {
	// Initial node refresh
	if err := m.RefreshNodes(); err != nil {
		return fmt.Errorf("failed initial node refresh: %w", err)
	}

	// Start periodic node refresh
	m.refreshTicker = time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-m.refreshTicker.C:
				if err := m.RefreshNodes(); err != nil {
					log.Error("Error refreshing nodes", "error", err)
				}
			case <-m.ctx.Done():
				return
			}
		}
	}()

	return nil
}

// Stop stops the manager's background operations
func (m *SwarmManager) Stop() {
	if m.refreshTicker != nil {
		m.refreshTicker.Stop()
	}
	m.cancel()
}

// InitSwarm initializes a new Swarm cluster
func (m *SwarmManager) InitSwarm(advertiseAddr string) (string, error) {
	req := swarm.InitRequest{
		ListenAddr:    "0.0.0.0:2377",
		AdvertiseAddr: advertiseAddr,
	}

	swarmID, err := m.client.SwarmInit(context.Background(), req)
	if err != nil {
		return "", fmt.Errorf("failed to initialize swarm: %w", err)
	}

	// Refresh nodes after init
	if err := m.RefreshNodes(); err != nil {
		log.Warn("Failed to refresh nodes after swarm init", "error", err)
	}

	return swarmID, nil
}

// GetJoinToken returns the token needed for a node to join the swarm
func (m *SwarmManager) GetJoinToken(isManager bool) (string, error) {
	swarmInfo, err := m.client.SwarmInspect(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to inspect swarm: %w", err)
	}

	if isManager {
		return swarmInfo.JoinTokens.Manager, nil
	}
	return swarmInfo.JoinTokens.Worker, nil
}

// RefreshNodes updates the node cache with the current state of the swarm
func (m *SwarmManager) RefreshNodes() error {
	nodes, err := m.client.NodeList(context.Background(), types.NodeListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list nodes: %w", err)
	}

	m.nodeCacheMu.Lock()
	defer m.nodeCacheMu.Unlock()

	// Clear the cache
	m.nodeCache = make(map[string]node.Info)

	// Update with fresh data
	for _, n := range nodes {
		nodeInfo := node.Info{
			ID:           n.ID,
			Hostname:     n.Description.Hostname,
			Address:      n.Status.Addr,
			Labels:       n.Spec.Labels,
			Role:         string(n.Spec.Role),
			Manager:      n.ManagerStatus != nil,
			Availability: string(n.Spec.Availability),
			Conditions:   []string{string(n.Status.State)},
			Capacity: node.Resources{
				CPU:    int(n.Description.Resources.NanoCPUs / 1e9),
				Memory: n.Description.Resources.MemoryBytes,
				Disk:   0, // Not available from Swarm API
				GPU:    0, // Not available from Swarm API
			},
		}
		m.nodeCache[n.ID] = nodeInfo
	}

	return nil
}

// GetNodes returns all nodes in the swarm
func (m *SwarmManager) GetNodes() []node.Info {
	m.nodeCacheMu.RLock()
	defer m.nodeCacheMu.RUnlock()

	nodes := make([]node.Info, 0, len(m.nodeCache))
	for _, nodeEntry := range m.nodeCache {
		nodes = append(nodes, nodeEntry)
	}

	return nodes
}

// GetNode returns information about a specific node
func (m *SwarmManager) GetNode(nodeID string) (node.Info, error) {
	m.nodeCacheMu.RLock()
	defer m.nodeCacheMu.RUnlock()

	nodeInfo, exists := m.nodeCache[nodeID]
	if !exists {
		return node.Info{}, errors.New("node not found")
	}

	return nodeInfo, nil
}

// UpdateNodeLabels updates the labels for a node
func (m *SwarmManager) UpdateNodeLabels(nodeID string, labels map[string]string) error {
	// Get current node spec
	swarmNode, _, err := m.client.NodeInspectWithRaw(context.Background(), nodeID)
	if err != nil {
		return fmt.Errorf("failed to inspect node: %w", err)
	}

	// Update labels
	spec := swarmNode.Spec
	spec.Labels = labels

	// Update node
	err = m.client.NodeUpdate(context.Background(), nodeID, swarmNode.Version, spec)
	if err != nil {
		return fmt.Errorf("failed to update node: %w", err)
	}

	// Refresh nodes after update
	if err := m.RefreshNodes(); err != nil {
		log.Warn("Failed to refresh nodes after label update", "error", err)
	}

	return nil
}

// DrainNode puts a node in drain state
func (m *SwarmManager) DrainNode(nodeID string) error {
	return m.updateNodeAvailability(nodeID, swarm.NodeAvailabilityDrain)
}

// ActivateNode puts a node in active state
func (m *SwarmManager) ActivateNode(nodeID string) error {
	return m.updateNodeAvailability(nodeID, swarm.NodeAvailabilityActive)
}

// updateNodeAvailability updates the availability of a node
func (m *SwarmManager) updateNodeAvailability(nodeID string, availability swarm.NodeAvailability) error {
	// Get current node spec
	swarmNode, _, err := m.client.NodeInspectWithRaw(context.Background(), nodeID)
	if err != nil {
		return fmt.Errorf("failed to inspect node: %w", err)
	}

	// Update availability
	spec := swarmNode.Spec
	spec.Availability = availability

	// Update node
	err = m.client.NodeUpdate(context.Background(), nodeID, swarmNode.Version, spec)
	if err != nil {
		return fmt.Errorf("failed to update node: %w", err)
	}

	// Refresh nodes after update
	if err := m.RefreshNodes(); err != nil {
		log.Warn("Failed to refresh nodes after availability update", "error", err)
	}

	return nil
}

// RemoveNode removes a node from the swarm
func (m *SwarmManager) RemoveNode(nodeID string, force bool) error {
	err := m.client.NodeRemove(context.Background(), nodeID, types.NodeRemoveOptions{Force: force})
	if err != nil {
		return fmt.Errorf("failed to remove node: %w", err)
	}

	// Refresh nodes after removal
	if err := m.RefreshNodes(); err != nil {
		log.Warn("Failed to refresh nodes after node removal", "error", err)
	}

	return nil
}

// DeployService deploys a service to the swarm
func (m *SwarmManager) DeployService(def config.ServiceDefinition) (string, error) {
	spec := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: def.Name,
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image: def.Image,
				Env:   def.ToEnv(),
			},
		},
		Mode: swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{Replicas: uint64Ptr(uint64(def.Replicas))},
		},
	}

	resp, err := m.client.ServiceCreate(context.Background(), spec, types.ServiceCreateOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create service: %w", err)
	}

	return resp.ID, nil
}

// UpdateService updates an existing service
func (m *SwarmManager) UpdateService(serviceID string, def config.ServiceDefinition) error {
	// Get current service spec
	service, _, err := m.client.ServiceInspectWithRaw(context.Background(), serviceID, types.ServiceInspectOptions{})
	if err != nil {
		return fmt.Errorf("failed to inspect service: %w", err)
	}

	// Update spec
	spec := service.Spec
	spec.TaskTemplate.ContainerSpec.Image = def.Image
	spec.TaskTemplate.ContainerSpec.Env = def.ToEnv()

	if def.Replicas > 0 {
		replicas := uint64(def.Replicas)
		if spec.Mode.Replicated != nil {
			spec.Mode.Replicated.Replicas = &replicas
		} else {
			spec.Mode = swarm.ServiceMode{
				Replicated: &swarm.ReplicatedService{Replicas: &replicas},
			}
		}
	}

	// Update service
	response, err := m.client.ServiceUpdate(context.Background(), serviceID, service.Version, spec, types.ServiceUpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}

	// Log warnings if any
	for _, warning := range response.Warnings {
		log.Warn("Warning during service update", "warning", warning)
	}

	return nil
}

// RemoveService removes a service from the swarm
func (m *SwarmManager) RemoveService(serviceID string) error {
	err := m.client.ServiceRemove(context.Background(), serviceID)
	if err != nil {
		return fmt.Errorf("failed to remove service: %w", err)
	}

	return nil
}

// GetServiceStatus returns the status of a service
func (m *SwarmManager) GetServiceStatus(serviceID string) (config.DeploymentStatus, error) {
	service, _, err := m.client.ServiceInspectWithRaw(context.Background(), serviceID, types.ServiceInspectOptions{})
	if err != nil {
		return config.DeploymentStatus{}, fmt.Errorf("failed to inspect service: %w", err)
	}

	// Get all tasks and filter for this service
	allTasks, err := m.client.TaskList(context.Background(), types.TaskListOptions{})
	if err != nil {
		return config.DeploymentStatus{}, fmt.Errorf("failed to list tasks: %w", err)
	}

	// Filter tasks for this service
	var serviceTasks []swarm.Task
	for _, task := range allTasks {
		if task.ServiceID == serviceID {
			serviceTasks = append(serviceTasks, task)
		}
	}

	// Determine overall state
	state := "running"
	if len(serviceTasks) == 0 {
		state = "pending"
	} else {
		failedTasks := 0
		for _, task := range serviceTasks {
			if task.Status.State == swarm.TaskStateFailed {
				failedTasks++
			}
		}

		// If all tasks failed, the service is failed
		if failedTasks == len(serviceTasks) {
			state = "failed"
		}
	}

	return config.DeploymentStatus{
		ID:    serviceID,
		State: state,
		Service: config.ServiceDefinition{
			Name:  service.Spec.Annotations.Name,
			Image: service.Spec.TaskTemplate.ContainerSpec.Image,
			// Environment would need parsing from env strings back to map
			Replicas: getReplicaCount(service.Spec),
		},
	}, nil
}

// ListServices returns all services in the swarm
func (m *SwarmManager) ListServices() ([]config.DeploymentStatus, error) {
	services, err := m.client.ServiceList(context.Background(), types.ServiceListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	result := make([]config.DeploymentStatus, 0, len(services))
	for _, service := range services {
		status, err := m.GetServiceStatus(service.ID)
		if err != nil {
			log.Warn("Failed to get status for service", "serviceID", service.ID, "error", err)
			continue
		}
		result = append(result, status)
	}

	return result, nil
}

// Helper functions

func uint64Ptr(n uint64) *uint64 {
	return &n
}

func getReplicaCount(spec swarm.ServiceSpec) int {
	if spec.Mode.Replicated != nil && spec.Mode.Replicated.Replicas != nil {
		return int(*spec.Mode.Replicated.Replicas)
	}
	return 0
}
