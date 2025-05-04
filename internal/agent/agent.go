package agent

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/jasonlovesdoggo/velo/internal/logs"
)

// NewContainerAgent creates a new ContainerAgent
func NewContainerAgent() (*ContainerAgent, error) {
	// Create Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	agent := &ContainerAgent{
		client:     cli,
		hostname:   hostname,
		ctx:        ctx,
		cancel:     cancel,
		containers: []ContainerInfo{},
	}

	// Determine if this node is a manager
	info, err := cli.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Docker info: %w", err)
	}

	if info.Swarm.NodeID != "" {
		agent.nodeID = info.Swarm.NodeID
		agent.isManager = info.Swarm.ControlAvailable
	} else {
		return nil, fmt.Errorf("node is not part of a swarm")
	}

	return agent, nil
}

// Start begins the agent's background operations
func (a *ContainerAgent) Start() error {
	// Initial container collection
	if err := a.collectContainers(); err != nil {
		return fmt.Errorf("failed initial container collection: %w", err)
	}

	// Start periodic container collection
	a.collectTicker = time.NewTicker(30 * time.Second)
	go func() {
		for {
			select {
			case <-a.collectTicker.C:
				if err := a.collectContainers(); err != nil {
					logs.Error("Error collecting containers", "error", err)
				}
			case <-a.ctx.Done():
				return
			}
		}
	}()

	// Start periodic health checks
	a.healthTicker = time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case <-a.healthTicker.C:
				if err := a.checkContainerHealth(); err != nil {
					logs.Error("Error checking container health", "error", err)
				}
			case <-a.ctx.Done():
				return
			}
		}
	}()

	logs.Info("Container agent started",
		"node", a.hostname, "nodeID", a.nodeID, "isManager", a.isManager)
	return nil
}

// Stop stops the agent's background operations
func (a *ContainerAgent) Stop() {
	if a.collectTicker != nil {
		a.collectTicker.Stop()
	}
	if a.healthTicker != nil {
		a.healthTicker.Stop()
	}
	a.cancel()
	logs.Info("Container agent stopped", "node", a.hostname)
}

// collectContainers collects information about running containers
func (a *ContainerAgent) collectContainers() error {
	// For now, just log that we're collecting containers
	// In a real implementation, this would use the Docker API to list containers
	logs.Info("Collecting container information", "node", a.hostname)

	// Simulate container collection with a placeholder
	a.containersMu.Lock()
	a.containers = []ContainerInfo{
		{
			ID:      "placeholder",
			Name:    "placeholder",
			Image:   "placeholder",
			Status:  "running",
			Running: true,
			Health:  "healthy",
		},
	}
	a.containersMu.Unlock()

	return nil
}

// checkContainerHealth checks the health of all containers
func (a *ContainerAgent) checkContainerHealth() error {
	// For now, just log that we're checking container health
	// In a real implementation, this would use the Docker API to check container health
	logs.Info("Checking container health", "node", a.hostname)
	return nil
}

// GetContainers returns information about all containers
func (a *ContainerAgent) GetContainers() []ContainerInfo {
	a.containersMu.RLock()
	defer a.containersMu.RUnlock()

	// Return a copy to avoid race conditions
	result := make([]ContainerInfo, len(a.containers))
	copy(result, a.containers)
	return result
}

// GetNodeInfo returns information about this node
func (a *ContainerAgent) GetNodeInfo() map[string]interface{} {
	a.containersMu.RLock()
	defer a.containersMu.RUnlock()

	return map[string]interface{}{
		"id":         a.nodeID,
		"hostname":   a.hostname,
		"is_manager": a.isManager,
		"containers": len(a.containers),
		"timestamp":  time.Now(),
	}
}

// RestartContainer restarts a container
func (a *ContainerAgent) RestartContainer(containerID string) error {
	timeout := container.StopOptions{
		Timeout: &[]int{10}[0],
	}
	return a.client.ContainerRestart(a.ctx, containerID, timeout)
}

// StopContainer stops a container
func (a *ContainerAgent) StopContainer(containerID string) error {
	timeout := container.StopOptions{
		Timeout: &[]int{10}[0],
	}
	return a.client.ContainerStop(a.ctx, containerID, timeout)
}

// StartContainer starts a container
func (a *ContainerAgent) StartContainer(containerID string) error {
	return a.client.ContainerStart(a.ctx, containerID, container.StartOptions{})
}
