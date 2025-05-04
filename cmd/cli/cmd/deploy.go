package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jasonlovesdoggo/velo/pkg/client"
	"github.com/spf13/cobra"
)

var (
	deployService string
	deployImage   string
	deployEnv     []string
)

func init() {
	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a service",
		Long:  `Deploy a service to the Velo platform.`,
		Run:   runDeploy,
	}

	deployCmd.Flags().StringVar(&deployService, "service", "test-service", "Name of the service to deploy")
	deployCmd.Flags().StringVar(&deployImage, "image", "nginx:latest", "Docker image to deploy")
	deployCmd.Flags().StringArrayVar(&deployEnv, "env", []string{}, "Environment variables (KEY=VALUE)")

	rootCmd.AddCommand(deployCmd)
}

func runDeploy(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	c, err := client.NewClient(serverAddr)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	envMap := make(map[string]string)
	for _, kv := range deployEnv {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			log.Fatalf("Invalid env var %q; expected KEY=VALUE", kv)
		}
		envMap[parts[0]] = parts[1]
	}

	resp, err := c.Deploy(ctx, deployService, deployImage, envMap)
	if err != nil {
		log.Fatalf("Failed to deploy service: %v", err)
	}

	fmt.Printf("Service deployed successfully!\nDeployment ID: %s\nStatus: %s\n",
		resp.DeploymentId, resp.Status)
}
