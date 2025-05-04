package orchestrator

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"github.com/jasonlovesdoggo/velo/internal/config"
	"github.com/jasonlovesdoggo/velo/pkg/core/node"
	"log"
)

func DeployToSwarm(def config.ServiceDefinition) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

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

	resp, err := cli.ServiceCreate(context.Background(), spec, types.ServiceCreateOptions{})
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

func RollbackDeployment(id string) error {
	// Placeholder for rollback logic
	log.Printf("Rolling back deployment %s", id)
	return nil
}

func GetDeploymentStatus(id string) (config.DeploymentStatus, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return config.DeploymentStatus{}, err
	}
	service, _, err := cli.ServiceInspectWithRaw(context.Background(), id, types.ServiceInspectOptions{})
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
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	swarmNodes, err := cli.NodeList(context.Background(), types.NodeListOptions{})
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

func uint64Ptr(n uint64) *uint64 {
	return &n
}
