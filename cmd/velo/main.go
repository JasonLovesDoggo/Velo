package main

import (
	"fmt"
	"github.com/jasonlovesdoggo/velo/internal/gateway"
	"log"
	"os"
)

func main() {
	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println(port)

	log.Println("Starting FleetStack API Gateway on port", port)
	if err := gateway.Start(port); err != nil {
		log.Fatal("Failed to start gateway:", err)
	}
}
