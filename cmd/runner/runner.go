package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/jasonlovesdoggo/velo/internal/agent"
	"github.com/jasonlovesdoggo/velo/internal/log"
)

func main() {
	flag.Parse()

	log.Info("Starting Velo Runner Service...")

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
	log.Info("Velo Runner Service stopped")
}
