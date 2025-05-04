package config

type ServiceDefinition struct {
	Name        string            `toml:"name"`
	Image       string            `toml:"image"`
	Environment map[string]string `toml:"environment"`
	Replicas    int               `toml:"replicas"`
}

type DeploymentStatus struct {
	ID      string
	Service ServiceDefinition
	State   string // pending, running, failed
	Logs    string
}
