package manager

import (
	"github.com/jasonlovesdoggo/velo/internal/deployment"
)

// Manager defines the interface for orchestration managers
type Manager interface {
	// DeployService deploys a service to the orchestration platform
	DeployService(def deployment.ServiceDefinition) (string, error)

	// RemoveService removes a service from the orchestration platform
	RemoveService(serviceID string) error

	// GetServiceStatus returns the status of a service
	GetServiceStatus(serviceID string) (deployment.DeploymentStatus, error)
}
