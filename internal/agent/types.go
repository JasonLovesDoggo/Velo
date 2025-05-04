package agent

import (
	"context"
	"github.com/docker/docker/client"
	"sync"
	"time"
)

// ContainerAgent is responsible for monitoring the local node and containers
type ContainerAgent struct {
	client        *client.Client
	hostname      string
	nodeID        string
	isManager     bool
	collectTicker *time.Ticker
	healthTicker  *time.Ticker
	ctx           context.Context
	cancel        context.CancelFunc
	containers    []ContainerInfo
	containersMu  sync.RWMutex
}

// ContainerInfo contains information about a container
type ContainerInfo struct {
	ID      string
	Name    string
	Image   string
	Status  string
	Running bool
	Health  string
}
