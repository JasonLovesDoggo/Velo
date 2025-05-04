package main

import (
	"fmt"
	"github.com/jasonlovesdoggo/velo/internal/config"
	"os"
)

func testConfig() {
	// Create a test TOML file
	testDir := "./test_config"
	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		fmt.Printf("Error creating test directory: %v\n", err)
		return
	}
	defer os.RemoveAll(testDir)

	// Create a test config file
	testConfig := `
name = "test-service"
image = "nginx:latest"
replicas = 2
networks = ["frontend", "backend"]
constraints = ["node.role==worker"]
dependencies = ["database"]

[environment]
ENV_VAR1 = "value1"
ENV_VAR2 = "value2"

[labels]
app = "test"
environment = "development"

[[volumes]]
source = "data-volume"
destination = "/data"
readonly = false

[resources]
cpu_limit = 1.0
memory_limit = 1073741824  # 1GB
cpu_reserve = 0.5
memory_reserve = 536870912  # 512MB

[healthcheck]
command = ["CMD", "curl", "-f", "http://localhost/health"]
interval = 30
timeout = 10
retries = 3
start_period = 5
`

	err = os.WriteFile(testDir+"/velo.toml", []byte(testConfig), 0644)
	if err != nil {
		fmt.Printf("Error writing test config file: %v\n", err)
		return
	}

	// Load the config
	cfg, err := config.LoadConfigFromFile(testDir)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		return
	}

	// Print the config
	fmt.Printf("Loaded config:\n")
	fmt.Printf("  Name: %s\n", cfg.Name)
	fmt.Printf("  Image: %s\n", cfg.Image)
	fmt.Printf("  Replicas: %d\n", cfg.Replicas)

	fmt.Printf("  Environment:\n")
	for k, v := range cfg.Environment {
		fmt.Printf("    %s: %s\n", k, v)
	}

	fmt.Printf("  Labels:\n")
	for k, v := range cfg.Labels {
		fmt.Printf("    %s: %s\n", k, v)
	}

	fmt.Printf("  Networks: %v\n", cfg.Networks)

	fmt.Printf("  Volumes:\n")
	for _, v := range cfg.Volumes {
		fmt.Printf("    %s -> %s (readonly: %v)\n", v.Source, v.Destination, v.ReadOnly)
	}

	fmt.Printf("  Resources:\n")
	fmt.Printf("    CPU Limit: %.2f\n", cfg.Resources.CPULimit)
	fmt.Printf("    Memory Limit: %d\n", cfg.Resources.MemoryLimit)
	fmt.Printf("    CPU Reserve: %.2f\n", cfg.Resources.CPUReserve)
	fmt.Printf("    Memory Reserve: %d\n", cfg.Resources.MemoryReserve)

	fmt.Printf("  Health Check:\n")
	fmt.Printf("    Command: %v\n", cfg.HealthCheck.Command)
	fmt.Printf("    Interval: %d\n", cfg.HealthCheck.Interval)
	fmt.Printf("    Timeout: %d\n", cfg.HealthCheck.Timeout)
	fmt.Printf("    Retries: %d\n", cfg.HealthCheck.Retries)
	fmt.Printf("    Start Period: %d\n", cfg.HealthCheck.StartPeriod)

	fmt.Printf("  Constraints: %v\n", cfg.Constraints)
	fmt.Printf("  Dependencies: %v\n", cfg.Dependencies)
}
