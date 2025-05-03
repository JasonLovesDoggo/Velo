package deployment

type ServiceDefinition struct {
	Name        string            `json:"name"`
	Image       string            `json:"image"`
	Environment map[string]string `json:"environment"`
	Replicas    int               `json:"replicas"`
}

type DeploymentStatus struct {
	ID      string
	Service ServiceDefinition
	State   string // pending, running, failed
	Logs    string
}
