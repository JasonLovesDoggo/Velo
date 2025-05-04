package main

import (
	"github.com/jasonlovesdoggo/velo/pkg/core"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/jasonlovesdoggo/velo/internal/log"
	"github.com/jasonlovesdoggo/velo/internal/orchestrator/manager"
	"github.com/jasonlovesdoggo/velo/internal/server"
)

func main() {
	log.Info("Starting Velo Management Server...")

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

	// Create and start the gRPC server
	deploymentServer := server.NewDeploymentServer(swarmManager)
	portstring := strconv.Itoa(core.Port)
	if err := deploymentServer.Start(":" + portstring); err != nil {
		log.Error("Failed to start gRPC server", "error", err)
		swarmManager.Stop()
		os.Exit(1)
	}
	log.Info("gRPC server started", "address", ":"+portstring)

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Stop the server and manager
	deploymentServer.Stop()
	swarmManager.Stop()
	log.Info("Velo Management Server stopped")
}
