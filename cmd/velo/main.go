package main

import (
	"fmt"
	"github.com/jasonlovesdoggo/velo/internal/gateway"
	"github.com/jasonlovesdoggo/velo/internal/log"
	"os"
)

func main() {
	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println(port)

	log.Info("Starting FleetStack API Gateway", "port", port)
	if err := gateway.Start(port); err != nil {
		log.Error("Failed to start gateway", "error", err)
		os.Exit(1)
	}
}
