package orchestrator

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/jasonlovesdoggo/velo/internal/config"
	"github.com/jasonlovesdoggo/velo/internal/log"
	"github.com/jasonlovesdoggo/velo/internal/orchestrator/gocker"
	"github.com/jasonlovesdoggo/velo/internal/utils"
	"github.com/jasonlovesdoggo/velo/pkg/core/node"
)

func DeployToSwarm(def config.ServiceDefinition) (string, error) {
	// Create annotations with labels if defined
	annotations := swarm.Annotations{
		Name: def.Name,
	}
	if def.Labels != nil && len(def.Labels) > 0 {
		annotations.Labels = def.Labels
	}

	// Create container spec
	containerSpec := &swarm.ContainerSpec{
		Image: def.Image,
		Env:   def.ToEnv(),
	}

	// Create task template
	taskTemplate := swarm.TaskSpec{
		ContainerSpec: containerSpec,
	}

	// Add networks if defined
	if def.Networks != nil && len(def.Networks) > 0 {
		var networks []swarm.NetworkAttachmentConfig
		for _, network := range def.Networks {
			networks = append(networks, swarm.NetworkAttachmentConfig{
				Target: network,
			})
		}
		taskTemplate.Networks = networks
	}

	// Create service spec
	spec := swarm.ServiceSpec{
		Annotations:  annotations,
		TaskTemplate: taskTemplate,
		Mode: swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{Replicas: utils.Uint64Ptr(uint64(def.Replicas))},
		},
	}

	resp, err := gocker.GetClient().ServiceCreate(context.Background(), spec, types.ServiceCreateOptions{})
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func RollbackDeployment(id string) error {
	// Placeholder for rollback logic
	log.Info("Rolling back deployment", "id", id)
	return nil
}

func GetDeploymentStatus(id string) (config.DeploymentStatus, error) {
	service, _, err := gocker.GetClient().ServiceInspectWithRaw(context.Background(), id, types.ServiceInspectOptions{})
	if err != nil {
		return config.DeploymentStatus{}, err
	}

	data, _ := json.MarshalIndent(service, "", "  ")
	return config.DeploymentStatus{
		ID:    id,
		State: "running",
		Logs:  string(data),
	}, nil
}

func ListNodes() ([]node.Info, error) {
	swarmNodes, err := gocker.GetClient().NodeList(context.Background(), types.NodeListOptions{})
	if err != nil {
		return nil, err
	}

	var nodes []node.Info
	for _, n := range swarmNodes {
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
				Disk:   0, // Disk info not available in Swarm API
				GPU:    0, // GPU info requires deeper integration
			},
		}
		nodes = append(nodes, nodeInfo)
	}
	return nodes, nil
}
