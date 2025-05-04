# Config Package

This package provides functionality for loading and validating service configuration files.

## ServiceDefinition

The `ServiceDefinition` struct represents a service to be deployed. It includes the following fields:

- `Name` (string): The name of the service
- `Image` (string): The Docker image to use
- `Environment` (map[string]string): Environment variables for the service
- `Replicas` (int): Number of replicas to deploy
- `Labels` (map[string]string): Docker labels for the service
- `Networks` ([]string): Networks to attach to the service
- `Volumes` ([]VolumeMount): Volumes to mount in the service
- `Resources` (ResourceConfig): CPU and memory limits and reservations
- `HealthCheck` (HealthCheckConfig): Health check configuration
- `Constraints` ([]string): Placement constraints for the service
- `Dependencies` ([]string): Services that this service depends on

## Configuration File Format

The configuration file is in TOML format. Here's an example:

```toml
name = "my-service"
image = "nginx:latest"
replicas = 2
networks = ["frontend", "backend"]
constraints = ["node.role==worker"]
dependencies = ["database"]

[environment]
ENV_VAR1 = "value1"
ENV_VAR2 = "value2"

[labels]
app = "my-app"
environment = "production"

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
```

## Usage

To load a configuration file:

```go
import "github.com/jasonlovesdoggo/velo/internal/config"

// Load the config file from the specified directory
def, err := config.LoadConfigFromFile("/path/to/directory")
if err != nil {
    // Handle error
}

// Use the config
fmt.Println(def.Name)
```

The `LoadConfigFromFile` function will search for a file named `velo.toml` in the specified directory and its subdirectories.