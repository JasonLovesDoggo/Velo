package manager

import "github.com/jasonlovesdoggo/velo/internal/config"

type DeploymentStatus struct {
	ID      string
	Service config.ServiceDefinition
	State   string // pending, running, failed
	Logs    string
}
