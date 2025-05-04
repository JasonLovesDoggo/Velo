package main

import (
	"context"
	"fmt"
	"github.com/jasonlovesdoggo/velo/internal/config"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jasonlovesdoggo/velo/pkg/client"
	"github.com/spf13/cobra"
)

var (
	serverAddr   string
	serviceName  string
	image        string
	deploymentID string
	timeout      time.Duration
	envVars      []string
)

func main() {
	// Create the root command
	rootCmd := &cobra.Command{
		Use:   "veloctl",
		Short: "Velo CLI - A command line interface for Velo",
		Long: `Velo CLI is a command line interface for Velo, a lightweight, 
self-hostable deployment and operations platform built on top of Docker Swarm.`,
	}

	// Add global flags
	rootCmd.PersistentFlags().StringVar(&serverAddr, "server", "localhost:50051", "The server address in the format host:port")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 10*time.Second, "Timeout for API requests")

	// Create the deploy command
	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a service",
		Long:  `Deploy a service to the Velo platform.`,
		Run:   runDeploy,
	}
	deployCmd.Flags().StringVar(&serviceName, "service", "test-service", "Name of the service to deploy")
	deployCmd.Flags().StringVar(&image, "image", "nginx:latest", "Docker image to deploy")
	deployCmd.Flags().StringArrayVar(&envVars, "env", []string{}, "Environment variables in the format KEY=VALUE")

	// Create the status command
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Get the status of a deployment",
		Long:  `Get the status of a deployment on the Velo platform.`,
		Run:   runStatus,
	}
	statusCmd.Flags().StringVar(&deploymentID, "id", "", "Deployment ID")
	statusCmd.MarkFlagRequired("id")

	// Create the rollback command
	rollbackCmd := &cobra.Command{
		Use:   "rollback",
		Short: "Rollback a deployment",
		Long:  `Rollback a deployment on the Velo platform.`,
		Run:   runRollback,
	}
	rollbackCmd.Flags().StringVar(&deploymentID, "id", "", "Deployment ID")
	rollbackCmd.MarkFlagRequired("id")

	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a configuration file",
		Long:  `Test a configuration file for validity.`,
		Run: func(cmd *cobra.Command, args []string) {
			pwd, err := os.Getwd()
			if err != nil {
				log.Fatalf("Failed to get current working directory: %v", err)
			}
			_, err = config.LoadConfigFromFile(pwd)
			if err != nil {
				log.Fatalf("Failed to load config file: %v", err)
			} else {
				fmt.Println("Config file is valid.")
			}

		},
	}

	// Add commands to the root command
	rootCmd.AddCommand(deployCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(rollbackCmd)
	rootCmd.AddCommand(validateCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// runDeploy handles the deploy command
func runDeploy(cmd *cobra.Command, args []string) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create a client
	c, err := client.NewClient(serverAddr)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	// Parse environment variables
	env := make(map[string]string)
	for _, e := range envVars {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) != 2 {
			log.Fatalf("Invalid environment variable format: %s. Expected KEY=VALUE", e)
		}
		env[parts[0]] = parts[1]
	}

	// Deploy the service
	resp, err := c.Deploy(ctx, serviceName, image, env)
	if err != nil {
		log.Fatalf("Failed to deploy service: %v", err)
	}

	fmt.Printf("Service deployed successfully!\nDeployment ID: %s\nStatus: %s\n", resp.DeploymentId, resp.Status)
}

// runStatus handles the status command
func runStatus(cmd *cobra.Command, args []string) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create a client
	c, err := client.NewClient(serverAddr)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	// Get the status
	resp, err := c.GetStatus(ctx, deploymentID)
	if err != nil {
		log.Fatalf("Failed to get status: %v", err)
	}

	fmt.Printf("Deployment Status: %s\nLogs: %s\n", resp.Status, resp.Logs)
}

// runRollback handles the rollback command
func runRollback(cmd *cobra.Command, args []string) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create a client
	c, err := client.NewClient(serverAddr)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	// Rollback the deployment
	resp, err := c.Rollback(ctx, deploymentID)
	if err != nil {
		log.Fatalf("Failed to rollback deployment: %v", err)
	}

	fmt.Printf("Rollback %s: %s\n", map[bool]string{true: "succeeded", false: "failed"}[resp.Success], resp.Message)
}
