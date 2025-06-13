# Velo

A lightweight, self-hostable deployment platform built on Docker Swarm. PaaS simplicity without the complexity.

## What is Velo?

Velo transforms Docker Swarm into a user-friendly deployment platform for small teams and homelab enthusiasts. Deploy services through a clean web interface, powerful CLI, or API - all while maintaining the simplicity and reliability of Docker Swarm.

## Quick Start

### Automated Installation
```bash
curl -sSL https://raw.githubusercontent.com/jasonlovesdoggo/velo/main/install.sh | bash
```

### Manual Installation
```bash
# Clone and build
git clone https://github.com/jasonlovesdoggo/velo.git
cd velo
make build

# Start manager node
./bin/velo --manager

# Access web interface at http://localhost:8080
# Default login: admin / admin
```

## Features

- **Web Interface**: Deploy and manage services through a modern web UI
- **CLI Tools**: Full-featured command line interface for automation
- **Authentication**: Built-in user management and session handling
- **Docker Swarm**: Leverages proven container orchestration
- **State Management**: Persistent service configuration and deployment history

## Usage

### Deploy via Web UI
1. Navigate to `http://localhost:8080`
2. Login with default credentials (admin/admin)
3. Use the Deploy Service form to create new deployments

### Deploy via CLI
```bash
# Deploy a service
./bin/veloctl deploy --name nginx --image nginx:latest --replicas 2

# Check status
./bin/veloctl status nginx

# Scale service
./bin/veloctl scale nginx --replicas 5
```

### Deploy via API
```bash
curl -X POST http://localhost:37355/api/deploy \
  -H "Content-Type: application/json" \
  -d '{"serviceName":"nginx","image":"nginx:latest","replicas":1}'
```

## Architecture

- **Server** (`velo`): Main application with manager and worker modes
- **CLI** (`veloctl`): Client interface for deployments and management
- **Web Interface**: Built-in HTTP server for browser-based management
- **gRPC API**: High-performance API for programmatic access

## Requirements

- Docker 20.10+ with Swarm mode enabled
- Go 1.24+ (for building from source)
- Linux/macOS (Windows support coming soon)

## Development

```bash
# Install dependencies
make deps

# Run tests
make test

# Start development server
make dev

# View all available commands
make help
```

## Status

Velo is production-ready for small to medium deployments. Active development continues with new features being added regularly. See [roadmap](docs/roadmap.md) for planned features.

## License

AGPL-3.[AGPLv3](LICENSE)

