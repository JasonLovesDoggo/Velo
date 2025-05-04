package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/jasonlovesdoggo/velo/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a configuration file",
		Long:  `Test a configuration file for validity.`,
		Run: func(cmd *cobra.Command, args []string) {
			pwd, err := os.Getwd()
			if err != nil {
				log.Fatalf("Failed to get working directory: %v", err)
			}
			if _, err := config.LoadConfigFromFile(pwd); err != nil {
				log.Fatalf("Invalid config: %v", err)
			}
			fmt.Println("Config file is valid.")
		},
	}
	rootCmd.AddCommand(validateCmd)
}
