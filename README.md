# Velo

Velo is a lightweight, self-hostable deployment and operations platform built on top of Docker Swarm. It's designed for small teams, homelab users, and edge deployments who want PaaS-like simplicity without the complexity of full Kubernetes managed services.

## Features

- **Multi-Interface Deployment**: Deploy services via CLI, Web UI, or chatbot integration
- **Security-First**: Built-in configuration and secret management with automated certificate provisioning
- **Production Ready**: Comprehensive observability with logging, metrics, and alerting capabilities
- **Extensible**: Pluggable architecture supporting custom hooks and CI/CD integration

## Getting Started

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/jasonlovesdoggo/velo.git
   cd velo
   ```

2. Build the project:
   ```bash
   go build -o bin/velo-manager ./cmd/manager
   go build -o bin/velo-client ./cmd/client
   ```

### Running the Server

Start the management server:

```bash
./bin/velo-manager
```

The server will start on port 50051 by default.

Please see the [CLI documentation](./cmd/cli/README.md) for available commands and options.


### Running Tests

Run the tests:

```bash
go test ./...
```

For more detailed documentation, see the [API Documentation](docs/api.md).

## Architecture

Velo follows a modular architecture with the following key components:

- **Control Plane**: Handles API gateway, service management, and core orchestration
- **Management Node**: Manages logging, backups, and configuration
- **Container Hosts**: Runs services with integrated metrics collection

## Development Status

🚧 This project is currently under active development. Features and APIs may change.

## Requirements

- Docker Swarm cluster
- Go 1.20 or later
- Docker Engine 20.10.0 or later
- gRPC tools (for development)

## Contributing

We welcome contributions! Please see our contributing guidelines (coming soon) for more details.

## License
[AGPLv3](LICENSE)

## Support

Please contact me @ velo[at]jasoncameron.dev for info and support
