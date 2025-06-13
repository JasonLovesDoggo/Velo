package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/jasonlovesdoggo/velo/pkg/client"
	"github.com/spf13/cobra"
)

var rollbackID string

func init() {
	rollbackCmd := &cobra.Command{
		Use:   "rollback",
		Short: "Rollback a deployment",
		Long:  `Rollback a deployment on the Velo platform.`,
		Run:   runRollback,
	}

	rollbackCmd.Flags().StringVar(&rollbackID, "id", "", "Deployment ID")
	rollbackCmd.MarkFlagRequired("id")

	rootCmd.AddCommand(rollbackCmd)
}

func runRollback(cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	c, err := client.NewClient(serverAddr)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer c.Close()

	resp, err := c.Rollback(ctx, rollbackID)
	if err != nil {
		log.Fatalf("Failed to rollback deployment: %v", err)
	}

	fmt.Printf("Rollback %s: %s\n",
		map[bool]string{true: "succeeded", false: "failed"}[resp.Success],
		resp.Message)
}
