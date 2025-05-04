# Velo CLI

The Velo CLI is a command-line interface for interacting with the Velo platform. It provides commands for deploying services, checking deployment status, and rolling back deployments.

## Installation

To install the Velo CLI, you need to have Go installed. Then, you can install it using:

```bash
go install github.com/jasonlovesdoggo/velo/cmd/cli@latest
```

## Usage

The Velo CLI provides the following commands:

### Deploy a Service

```bash
veloctl deploy --service <service-name> --image <image-name> --env KEY1=VALUE1 --env KEY2=VALUE2
```

Options:
- `--service`: Name of the service to deploy (default: "test-service")
- `--image`: Docker image to deploy (default: "nginx:latest")
- `--env`: Environment variables in the format KEY=VALUE (can be specified multiple times)

### Check Deployment Status

```bash
veloctl status --id <deployment-id>
```

Options:
- `--id`: Deployment ID (required)

### Rollback a Deployment

```bash
veloctl rollback --id <deployment-id>
```

Options:
- `--id`: Deployment ID (required)

### Validate Configuration

```bash
veloctl validate --config <config-file>
```

This command tests a configuration file for validity.

## Global Options

The following options can be used with any command:

- `--server`: The server address in the format host:port (default: "localhost:37355")
- `--timeout`: Timeout for API requests (default: 10s)

## Examples

Deploy a service:

```bash
veloctl deploy --service my-app --image my-org/my-app:latest --env NODE_ENV=production
```

Check deployment status:

```bash
veloctl status --id deployment-123
```

Rollback a deployment:

```bash
veloctl rollback --id deployment-123
```

## Development

The CLI is built using [Cobra](https://github.com/spf13/cobra), a powerful library for creating modern CLI applications. It communicates with the Velo server using gRPC.

To build the CLI from source:

```bash
go build -o veloctl cmd/cli/main.go
```

To run tests:

```bash
go test ./cmd/cli
```