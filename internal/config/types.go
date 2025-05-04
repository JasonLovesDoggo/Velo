package config

type ServiceDefinition struct {
	Name         string            `mapstructure:"name"`
	Image        string            `mapstructure:"image"`
	Environment  map[string]string `mapstructure:"environment"`
	Replicas     int               `mapstructure:"replicas"`
	Labels       map[string]string `mapstructure:"labels"`
	Networks     []string          `mapstructure:"networks"`
	Volumes      []VolumeMount     `mapstructure:"volumes"`
	Resources    ResourceConfig    `mapstructure:"resources"`
	HealthCheck  HealthCheckConfig `mapstructure:"healthcheck"`
	Constraints  []string          `mapstructure:"constraints"`
	Dependencies []string          `mapstructure:"dependencies"`
}

type VolumeMount struct {
	Source      string `mapstructure:"source"`
	Destination string `mapstructure:"destination"`
	ReadOnly    bool   `mapstructure:"readonly"`
}

type ResourceConfig struct {
	CPULimit      float64 `mapstructure:"cpu_limit"`
	MemoryLimit   int64   `mapstructure:"memory_limit"`
	CPUReserve    float64 `mapstructure:"cpu_reserve"`
	MemoryReserve int64   `mapstructure:"memory_reserve"`
}

type HealthCheckConfig struct {
	Command     []string `mapstructure:"command"`
	Interval    int      `mapstructure:"interval"`
	Timeout     int      `mapstructure:"timeout"`
	Retries     int      `mapstructure:"retries"`
	StartPeriod int      `mapstructure:"start_period"`
}

type DeploymentStatus struct {
	ID      string
	Service ServiceDefinition
	State   string // pending, running, failed
	Logs    string
}
