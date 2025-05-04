package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/jasonlovesdoggo/velo/pkg/client"
	"github.com/spf13/cobra"
)

var statusID string

func init() {
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Get the status of a deployment",
		Long:  `Get the status of a deployment on the Velo platform.`,
		Run:   runStatus,
	}

	statusCmd.Flags().StringVar(&statusID, "id", "", "Deployment ID")
	statusCmd.MarkFlagRequired("id")

	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	c, err := client.NewClient(serverAddr)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	resp, err := c.GetStatus(ctx, statusID)
	if err != nil {
		log.Fatalf("Failed to get status: %v", err)
	}

	fmt.Printf("Deployment Status: %s\nLogs: %s\n", resp.Status, resp.Logs)
}
