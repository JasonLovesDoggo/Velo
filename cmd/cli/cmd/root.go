package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	serverAddr string
	timeout    time.Duration
)

var rootCmd = &cobra.Command{
	Use:   "veloctl",
	Short: "Velo CLI - A command line interface for Velo",
	Long: `Velo CLI is a command line interface for Velo, a lightweight,
self-hostable deployment and operations platform built on top of Docker Swarm.`,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&serverAddr, "server", "localhost:37355", "The server address in host:port format")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 10*time.Second, "Timeout for API requests")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
