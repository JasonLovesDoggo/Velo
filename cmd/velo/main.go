package main

import (
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/jasonlovesdoggo/velo/internal/agent"
	"github.com/jasonlovesdoggo/velo/internal/auth"
	"github.com/jasonlovesdoggo/velo/internal/log"
	"github.com/jasonlovesdoggo/velo/internal/orchestrator/manager"
	"github.com/jasonlovesdoggo/velo/internal/server"
	"github.com/jasonlovesdoggo/velo/internal/state"
	"github.com/jasonlovesdoggo/velo/internal/web"
	"github.com/jasonlovesdoggo/velo/pkg/core"
)

func main() {
	// Parse command line flags
	isManager := flag.Bool("manager", false, "Run as manager daemon")
	webPort := flag.String("web-port", "8080", "Web interface port")
	flag.Parse()

	if *isManager {
		runManager(*webPort)
	} else {
		runWorker()
	}
}

func runManager(webPort string) {
	log.Info("Starting Velo Management Server...")

	// Initialize state store
	stateStore, err := state.NewDefaultStateStore()
	if err != nil {
		log.Error("Failed to create state store", "error", err)
		os.Exit(1)
	}

	// Initialize auth service
	authService := auth.NewAuthService(stateStore)
	if err := authService.Initialize(); err != nil {
		log.Error("Failed to initialize auth service", "error", err)
		os.Exit(1)
	}

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
	deploymentServer := server.NewDeploymentServer(swarmManager, authService)
	portstring := strconv.Itoa(core.Port)
	if err := deploymentServer.Start(":" + portstring); err != nil {
		log.Error("Failed to start gRPC server", "error", err)
		swarmManager.Stop()
		os.Exit(1)
	}
	log.Info("gRPC server started", "address", ":"+portstring)

	// Create and start the web server
	webServer := web.NewWebServer(swarmManager, authService, webPort)
	go func() {
		if err := webServer.Start(); err != nil {
			log.Error("Failed to start web server", "error", err)
		}
	}()
	log.Info("Web server started", "address", ":"+webPort)

	// Wait for termination signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// Stop the servers and manager
	deploymentServer.Stop()
	if err := webServer.Stop(); err != nil {
		log.Error("Error stopping web server", "error", err)
	}
	swarmManager.Stop()
	stateStore.Close()
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
