package main

import (
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/jasonlovesdoggo/velo/internal/agent"
	"github.com/jasonlovesdoggo/velo/internal/log"
	"github.com/jasonlovesdoggo/velo/internal/orchestrator/manager"
	"github.com/jasonlovesdoggo/velo/internal/server"
	"github.com/jasonlovesdoggo/velo/pkg/core"
)

func main() {
	// Parse command line flags
	isManager := flag.Bool("manager", false, "Run as manager daemon")
	flag.Parse()

	if *isManager {
		runManager()
	} else {
		runWorker()
	}
}

func runManager() {
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

func runWorker() {
	log.Info("Starting Velo Container Agent...")

	// Create a new container agent
	containerAgent, err := agent.NewContainerAgent()
	if err != nil {
		log.Error("Failed to create container agent", "error", err)
		os.Exit(1)
	}

	// Start the agent
	if err := containerAgent.Start(); err != nil {
		log.Error("Failed to start container agent", "error", err)
		os.Exit(1)
	}

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Stop the agent
	containerAgent.Stop()
	log.Info("Velo Container Agent stopped")
}
