package config

import (
	"fmt"
	"github.com/jasonlovesdoggo/velo/internal/log"
	"os"
)
import "github.com/BurntSushi/toml"

func LoadConfigFromFile(directoryPath string) (*ServiceDefinition, error) {
	// Load the config file
	config, err := loadConfig(directoryPath)
	if err != nil {
		return nil, err
	}

	// Validate the config
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func validateConfig(config *ServiceDefinition) error {
	// Validate the config fields
	if config.Name == "" {
		return fmt.Errorf("service name is required")
	}
	if config.Image == "" {
		return fmt.Errorf("service image is required")
	}
	if config.Replicas <= 0 {
		return fmt.Errorf("service replicas must be greater than 0")
	}
	return nil
}

func loadConfig(directory string) (*ServiceDefinition, error) {
	for _, dir := range DirNames {
		filePath := fmt.Sprintf("%s/%s/%s", directory, dir, FileName)
		if _, err := os.Stat(filePath); err == nil {
			// File exists, load it
			var config ServiceDefinition
			if _, err := toml.DecodeFile(filePath, &config); err != nil {
				log.Error("Failed to decode config file", "file", filePath, "error", err)
				return nil, ErrInvalidConfig
			}
			return &config, nil
		} else if os.IsNotExist(err) {
			// File does not exist, continue to the next directory
			continue
		} else {
			// Some other error occurred
			return nil, fmt.Errorf("error checking file %s: %w", filePath, err)
		}
	}
	return nil, ErrConfigNotFound
}
