package main

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/jasonlovesdoggo/velo/api/proto"
	"github.com/jasonlovesdoggo/velo/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

// setupBufConn creates a bufconn listener and returns a server and a dialer function
func setupBufConn() (*grpc.Server, func(context.Context, string) (net.Conn, error), func()) {
	listener := bufconn.Listen(1024 * 1024)

	// Create a gRPC server and register the test service
	s := grpc.NewServer()

	// Register a mock service
	proto.RegisterDeploymentServiceServer(s, &mockDeploymentService{})

	// Start the server
	go func() {
		if err := s.Serve(listener); err != nil {
			panic(err)
		}
	}()

	// Create a dialer function
	dialer := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}

	// Return the server, dialer, and a cleanup function
	return s, dialer, func() {
		s.Stop()
	}
}

// mockDeploymentService is a mock implementation of the DeploymentServiceServer interface
type mockDeploymentService struct {
	// Embed the UnimplementedDeploymentServiceServer to ensure forward compatibility
	proto.UnimplementedDeploymentServiceServer
}

// Deploy implements the Deploy method of the DeploymentServiceServer interface
func (s *mockDeploymentService) Deploy(ctx context.Context, req *proto.DeployRequest) (*proto.DeployResponse, error) {
	return &proto.DeployResponse{
		DeploymentId: "test-deployment-id",
		Status:       "deployed",
	}, nil
}

// Rollback implements the Rollback method of the DeploymentServiceServer interface
func (s *mockDeploymentService) Rollback(ctx context.Context, req *proto.RollbackRequest) (*proto.GenericResponse, error) {
	return &proto.GenericResponse{
		Message: "Deployment rolled back successfully",
		Success: true,
	}, nil
}

// GetStatus implements the GetStatus method of the DeploymentServiceServer interface
func (s *mockDeploymentService) GetStatus(ctx context.Context, req *proto.StatusRequest) (*proto.StatusResponse, error) {
	return &proto.StatusResponse{
		Status: "running",
		Logs:   "Service is running",
	}, nil
}

// setupTestClient creates a test client that uses the provided dialer
func setupTestClient(ctx context.Context, dialer func(context.Context, string) (net.Conn, error)) (*client.Client, error) {
	// Create a connection using the dialer
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	// Create a client using the connection
	c := client.NewClientWithConn(conn)

	return c, nil
}

func TestDeployService(t *testing.T) {
	// Set up the test server and client
	_, dialer, cleanup := setupBufConn()
	defer cleanup()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a test client
	c, err := setupTestClient(ctx, dialer)
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}
	defer c.Close()

	// Deploy a service
	resp, err := c.Deploy(ctx, "test-service", "nginx:latest", map[string]string{"ENV": "test"})
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
	// Set up the test server and client
	_, dialer, cleanup := setupBufConn()
	defer cleanup()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a test client
	c, err := setupTestClient(ctx, dialer)
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}
	defer c.Close()

	// Get the status
	resp, err := c.GetStatus(ctx, "test-deployment-id")
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
	// Set up the test server and client
	_, dialer, cleanup := setupBufConn()
	defer cleanup()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create a test client
	c, err := setupTestClient(ctx, dialer)
	if err != nil {
		t.Fatalf("Failed to create test client: %v", err)
	}
	defer c.Close()

	// Rollback a deployment
	resp, err := c.Rollback(ctx, "test-deployment-id")
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
