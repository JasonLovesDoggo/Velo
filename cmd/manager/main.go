package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/jasonlovesdoggo/velo/internal/log"
	"github.com/jasonlovesdoggo/velo/internal/orchestrator/manager"
)

func main() {
	log.Info("Starting Velo Swarm Manager Controller...")

	// Create a new swarm manager
	swarmManager, err := manager.NewSwarmManager()
	if err != nil {
		log.Error("Failed to create swarm manager", "error", err)
		os.Exit(1)
	}

	// Start the manager
	if err := swarmManager.Start(); err != nil {
		log.Error("Failed to start swarm manager", "error", err)
		os.Exit(1)
	}

	// Print node information
	nodes := swarmManager.GetNodes()
	log.Info("Managing nodes in the swarm", "count", len(nodes))
	for _, node := range nodes {
		log.Info("Node details", "hostname", node.Hostname, "id", node.ID, "isManager", node.Manager)
	}

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Stop the manager
	swarmManager.Stop()
	log.Info("Velo Swarm Manager Controller stopped")
}
