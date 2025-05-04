package main

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/jasonlovesdoggo/velo/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

// setupBufConn creates a bufconn listener and returns a client connection
func setupBufConn() (*grpc.ClientConn, func()) {
	listener := bufconn.Listen(1024 * 1024)

	// Create a gRPC server and register the test service
	s := grpc.NewServer()

	// Register a mock service
	server.RegisterDeploymentServiceServer(s, &mockDeploymentService{})

	// Start the server
	go func() {
		if err := s.Serve(listener); err != nil {
			panic(err)
		}
	}()

	// Create a client connection
	dialer := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}

	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithInsecure(),
	)
	if err != nil {
		panic(err)
	}

	// Return the connection and a cleanup function
	return conn, func() {
		conn.Close()
		s.Stop()
	}
}

// mockDeploymentService is a mock implementation of the DeploymentServiceServer interface
type mockDeploymentService struct {
	// No embedded type needed, we'll implement all required methods
}

// Deploy implements the Deploy method of the DeploymentServiceServer interface
func (s *mockDeploymentService) Deploy(ctx context.Context, req *server.DeployRequest) (*server.DeployResponse, error) {
	return &server.DeployResponse{
		DeploymentId: "test-deployment-id",
		Status:       "deployed",
	}, nil
}

// Rollback implements the Rollback method of the DeploymentServiceServer interface
func (s *mockDeploymentService) Rollback(ctx context.Context, req *server.RollbackRequest) (*server.GenericResponse, error) {
	return &server.GenericResponse{
		Message: "Deployment rolled back successfully",
		Success: true,
	}, nil
}

// GetStatus implements the GetStatus method of the DeploymentServiceServer interface
func (s *mockDeploymentService) GetStatus(ctx context.Context, req *server.StatusRequest) (*server.StatusResponse, error) {
	return &server.StatusResponse{
		Status: "running",
		Logs:   "Service is running",
	}, nil
}

func TestDeployService(t *testing.T) {
	conn, cleanup := setupBufConn()
	defer cleanup()

	// Create a client
	client := &grpcClient{conn: conn}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a deploy request
	req := &server.DeployRequest{
		ServiceName: "test-service",
		Image:       "nginx:latest",
		Env:         map[string]string{"ENV": "test"},
	}

	// Call the Deploy method
	resp, err := client.Deploy(ctx, req)
	if err != nil {
		t.Fatalf("Failed to deploy service: %v", err)
	}

	// Check the response
	if resp.DeploymentId != "test-deployment-id" {
		t.Errorf("Expected deployment ID %q, got %q", "test-deployment-id", resp.DeploymentId)
	}

	if resp.Status != "deployed" {
		t.Errorf("Expected status %q, got %q", "deployed", resp.Status)
	}
}

func TestGetStatus(t *testing.T) {
	conn, cleanup := setupBufConn()
	defer cleanup()

	// Create a client
	client := &grpcClient{conn: conn}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a status request
	req := &server.StatusRequest{
		DeploymentId: "test-deployment-id",
	}

	// Call the GetStatus method
	resp, err := client.GetStatus(ctx, req)
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}

	// Check the response
	if resp.Status != "running" {
		t.Errorf("Expected status %q, got %q", "running", resp.Status)
	}

	if resp.Logs != "Service is running" {
		t.Errorf("Expected logs %q, got %q", "Service is running", resp.Logs)
	}
}

func TestRollbackDeployment(t *testing.T) {
	conn, cleanup := setupBufConn()
	defer cleanup()

	// Create a client
	client := &grpcClient{conn: conn}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a rollback request
	req := &server.RollbackRequest{
		DeploymentId: "test-deployment-id",
	}

	// Call the Rollback method
	resp, err := client.Rollback(ctx, req)
	if err != nil {
		t.Fatalf("Failed to rollback deployment: %v", err)
	}

	// Check the response
	if !resp.Success {
		t.Errorf("Expected success to be true, got false")
	}

	if resp.Message != "Deployment rolled back successfully" {
		t.Errorf("Expected message %q, got %q", "Deployment rolled back successfully", resp.Message)
	}
}
