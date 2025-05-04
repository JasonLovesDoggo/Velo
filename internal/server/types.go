package server

// DeployRequest represents a request to deploy a service
type DeployRequest struct {
	ServiceName string
	Image       string
	Env         map[string]string
}

// DeployResponse represents a response to a deploy request
type DeployResponse struct {
	DeploymentId string
	Status       string
}

// RollbackRequest represents a request to rollback a deployment
type RollbackRequest struct {
	DeploymentId string
}

// GenericResponse represents a generic response
type GenericResponse struct {
	Message string
	Success bool
}

// StatusRequest represents a request to get the status of a deployment
type StatusRequest struct {
	DeploymentId string
}

// StatusResponse represents a response to a status request
type StatusResponse struct {
	Status string
	Logs   string
}
