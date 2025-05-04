# Velo API Documentation

This document describes the gRPC API for the Velo deployment and operations platform.

## Overview

Velo provides a gRPC API for deploying and managing services on Docker Swarm. The API is defined in the `api/proto/velo.proto` file and implemented by the `DeploymentServer` in the `internal/server` package.

## Service Definition

The gRPC service is defined as follows:

```protobuf
service DeploymentService {
  rpc Deploy (DeployRequest) returns (DeployResponse);
  rpc Rollback (RollbackRequest) returns (GenericResponse);
  rpc GetStatus (StatusRequest) returns (StatusResponse);
}
```

## Methods

### Deploy

Deploys a service to the Docker Swarm cluster.

**Request:**
```protobuf
message DeployRequest {
  string service_name = 1;
  string image = 2;
  map<string, string> env = 3;
}
```

**Response:**
```protobuf
message DeployResponse {
  string deployment_id = 1;
  string status = 2;
}
```

**Example:**
```go
// Create a client
client := &grpcClient{conn: conn}

// Create a deploy request
req := &server.DeployRequest{
    ServiceName: "my-service",
    Image:       "nginx:latest",
    Env:         map[string]string{"ENV": "production"},
}

// Call the Deploy method
resp, err := client.Deploy(ctx, req)
if err != nil {
    log.Fatalf("Failed to deploy service: %v", err)
}

fmt.Printf("Service deployed successfully!\nDeployment ID: %s\nStatus: %s\n", resp.DeploymentId, resp.Status)
```

### Rollback

Rolls back a deployment by removing the service from the Docker Swarm cluster.

**Request:**
```protobuf
message RollbackRequest {
  string deployment_id = 1;
}
```

**Response:**
```protobuf
message GenericResponse {
  string message = 1;
  bool success = 2;
}
```

**Example:**
```go
// Create a client
client := &grpcClient{conn: conn}

// Create a rollback request
req := &server.RollbackRequest{
    DeploymentId: "service-123",
}

// Call the Rollback method
resp, err := client.Rollback(ctx, req)
if err != nil {
    log.Fatalf("Failed to rollback deployment: %v", err)
}

fmt.Printf("Rollback %s: %s\n", map[bool]string{true: "succeeded", false: "failed"}[resp.Success], resp.Message)
```

### GetStatus

Gets the status of a deployed service.

**Request:**
```protobuf
message StatusRequest {
  string deployment_id = 1;
}
```

**Response:**
```protobuf
message StatusResponse {
  string status = 1;
  string logs = 2;
}
```

**Example:**
```go
// Create a client
client := &grpcClient{conn: conn}

// Create a status request
req := &server.StatusRequest{
    DeploymentId: "service-123",
}

// Call the GetStatus method
resp, err := client.GetStatus(ctx, req)
if err != nil {
    log.Fatalf("Failed to get status: %v", err)
}

fmt.Printf("Deployment Status: %s\nLogs: %s\n", resp.Status, resp.Logs)
```

## Status Values

The `status` field in the `StatusResponse` can have the following values:

- `pending`: The service is being created or updated.
- `running`: The service is running.
- `failed`: The service failed to start or has crashed.

## Client Implementation

The client implementation is in the `cmd/client/main.go` file. It provides a command-line interface for interacting with the gRPC API.

### Usage

```
velo-client --server=localhost:37355 --action=deploy --service=my-service --image=nginx:latest
velo-client --server=localhost:37355 --action=status --id=service-123
velo-client --server=localhost:37355 --action=rollback --id=service-123
```

### Options

- `--server`: The address of the gRPC server in the format `host:port`. Default: `localhost:37355`.
- `--action`: The action to perform: `deploy`, `status`, or `rollback`. Default: `deploy`.
- `--service`: The name of the service to deploy. Required for the `deploy` action.
- `--image`: The Docker image to deploy. Required for the `deploy` action.
- `--id`: The deployment ID for the `status` and `rollback` actions.

## Server Implementation

The server implementation is in the `internal/server/server.go` file. It provides a gRPC server that implements the `DeploymentService` interface.

### Starting the Server

```go
// Create a new swarm manager
swarmManager, err := manager.NewSwarmManager()
if err != nil {
    log.Error("Failed to create swarm manager", "error", err)
    os.Exit(1)
}

// Start the manager
if err := swarmManager.Start(); err != nil {
    log.Error("Failed to start swarm manager", "error", err)
    os.Exit(1)
}

// Create and start the gRPC server
deploymentServer := server.NewDeploymentServer(swarmManager)
if err := deploymentServer.Start(":37355"); err != nil {
    log.Error("Failed to start gRPC server", "error", err)
    swarmManager.Stop()
    os.Exit(1)
}
```

### Stopping the Server

```go
// Stop the server and manager
deploymentServer.Stop()
swarmManager.Stop()
```