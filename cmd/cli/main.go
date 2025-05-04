package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/jasonlovesdoggo/velo/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Define command-line flags
	serverAddr := flag.String("server", "localhost:50051", "The server address in the format host:port")
	action := flag.String("action", "deploy", "Action to perform: deploy, status, rollback, test-config")
	serviceName := flag.String("service", "test-service", "Name of the service to deploy")
	image := flag.String("image", "nginx:latest", "Docker image to deploy")
	deploymentID := flag.String("id", "", "Deployment ID for status or rollback actions")
	flag.Parse()

	// Set up a connection to the server
	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Perform the requested action
	switch *action {
	case "deploy":
		deployService(ctx, conn, *serviceName, *image)
	case "status":
		if *deploymentID == "" {
			log.Fatal("Deployment ID is required for status action")
		}
		getStatus(ctx, conn, *deploymentID)
	case "rollback":
		if *deploymentID == "" {
			log.Fatal("Deployment ID is required for rollback action")
		}
		rollbackDeployment(ctx, conn, *deploymentID)
	case "test-config":
		// This action doesn't require a connection to the server
		cancel() // Cancel the context since we don't need it
		testConfig()
	default:
		log.Fatalf("Unknown action: %s", *action)
	}
}

func deployService(ctx context.Context, conn *grpc.ClientConn, serviceName, image string) {
	// Create a client using the generated client interface
	client := proto.NewDeploymentServiceClient(conn)

	// Create a deploy request
	req := &proto.DeployRequest{
		ServiceName: serviceName,
		Image:       image,
		Env:         map[string]string{"ENV": "test"},
	}

	// Call the Deploy method
	resp, err := client.Deploy(ctx, req)
	if err != nil {
		log.Fatalf("Failed to deploy service: %v", err)
	}

	fmt.Printf("Service deployed successfully!\nDeployment ID: %s\nStatus: %s\n", resp.DeploymentId, resp.Status)
}

func getStatus(ctx context.Context, conn *grpc.ClientConn, deploymentID string) {
	// Create a client using the generated client interface
	client := proto.NewDeploymentServiceClient(conn)

	// Create a status request
	req := &proto.StatusRequest{
		DeploymentId: deploymentID,
	}

	// Call the GetStatus method
	resp, err := client.GetStatus(ctx, req)
	if err != nil {
		log.Fatalf("Failed to get status: %v", err)
	}

	fmt.Printf("Deployment Status: %s\nLogs: %s\n", resp.Status, resp.Logs)
}

func rollbackDeployment(ctx context.Context, conn *grpc.ClientConn, deploymentID string) {
	// Create a client using the generated client interface
	client := proto.NewDeploymentServiceClient(conn)

	// Create a rollback request
	req := &proto.RollbackRequest{
		DeploymentId: deploymentID,
	}

	// Call the Rollback method
	resp, err := client.Rollback(ctx, req)
	if err != nil {
		log.Fatalf("Failed to rollback deployment: %v", err)
	}

	fmt.Printf("Rollback %s: %s\n", map[bool]string{true: "succeeded", false: "failed"}[resp.Success], resp.Message)
}
