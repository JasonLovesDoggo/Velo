package cmd

import (
	"errors"
	"fmt"
	"github.com/jasonlovesdoggo/velo/internal/log"
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
				log.Fatal("Failed to get working directory", err)
			}
			if _, err := config.LoadConfigFromFile(pwd); err != nil {
				if errors.Is(err, config.ErrConfigNotFound) {
					fmt.Printf("Config file not found. Please create a %s file.\n", config.FileName)
				} else {
					fmt.Printf("Config file is invalid: %v\n", err)
				}

			} else {
				fmt.Printf("Config file is valid.")
			}
		},
	}
	rootCmd.AddCommand(validateCmd)
}
