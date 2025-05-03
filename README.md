# Velo

Velo is a lightweight, self-hostable deployment and operations platform built on top of Docker Swarm. It's designed for small teams, homelab users, and edge deployments who want PaaS-like simplicity without the complexity of full Kubernetes managed services.

## Features

- **Multi-Interface Deployment**: Deploy services via CLI, Web UI, or chatbot integration
- **Security-First**: Built-in configuration and secret management with automated certificate provisioning
- **Production Ready**: Comprehensive observability with logging, metrics, and alerting capabilities
- **Extensible**: Pluggable architecture supporting custom hooks and CI/CD integration

## Getting Started

[Documentation and installation instructions coming soon]

## Architecture

Velo follows a modular architecture with the following key components:

- **Control Plane**: Handles API gateway, service management, and core orchestration
- **Management Node**: Manages logging, backups, and configuration
- **Container Hosts**: Runs services with integrated metrics collection

## Development Status

🚧 This project is currently under active development. Features and APIs may change.

## Requirements

- Docker Swarm cluster
- Go 1.24.2 or later
- [Additional requirements to be specified]

## Contributing

We welcome contributions! Please see our contributing guidelines (coming soon) for more details.

## License
[AGPLv3](LICENSE)

## Support

Please contact me @ velo[at]jasoncameron.dev for info and support