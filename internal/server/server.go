package server

import (
	"context"
	"fmt"
	"net"

	"github.com/jasonlovesdoggo/velo/api/proto"
	"github.com/jasonlovesdoggo/velo/internal/auth"
	"github.com/jasonlovesdoggo/velo/internal/config"
	"github.com/jasonlovesdoggo/velo/internal/log"
	"github.com/jasonlovesdoggo/velo/internal/orchestrator/manager"
	"google.golang.org/grpc"
)

// DeploymentServer implements the proto.DeploymentServiceServer interface
type DeploymentServer struct {
	proto.UnimplementedDeploymentServiceServer
	manager     manager.Manager
	authService *auth.AuthService
	server      *grpc.Server
}

// NewDeploymentServer creates a new DeploymentServer
func NewDeploymentServer(manager manager.Manager, authService *auth.AuthService) *DeploymentServer {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(authService.AuthInterceptor),
	)

	return &DeploymentServer{
		manager:     manager,
		authService: authService,
		server:      server,
	}
}

// Start starts the gRPC server
func (s *DeploymentServer) Start(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	// Register the deployment service
	proto.RegisterDeploymentServiceServer(s.server, s)

	log.Info("Starting gRPC server", "address", address)
	go func() {
		if err := s.server.Serve(lis); err != nil {
			log.Error("Failed to serve gRPC", "error", err)
		}
	}()

	return nil
}

// Stop stops the gRPC server
func (s *DeploymentServer) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
		log.Info("gRPC server stopped")
	}
}

// Deploy handles the Deploy RPC call
func (s *DeploymentServer) Deploy(ctx context.Context, req *proto.DeployRequest) (*proto.DeployResponse, error) {
	log.Info("Received Deploy request", "service", req.ServiceName, "image", req.Image)

	// Convert the request to a ServiceDefinition
	serviceDef := config.ServiceDefinition{
		Name:        req.ServiceName,
		Image:       req.Image,
		Environment: req.Env,
		Replicas:    1, // Default to 1 replica
	}

	// Deploy the service
	deploymentID, err := s.manager.DeployService(serviceDef)
	if err != nil {
		log.Error("Failed to deploy service", "error", err)
		return nil, fmt.Errorf("failed to deploy service: %w", err)
	}

	return &proto.DeployResponse{
		DeploymentId: deploymentID,
		Status:       "deployed",
	}, nil
}

// Rollback handles the Rollback RPC call
func (s *DeploymentServer) Rollback(ctx context.Context, req *proto.RollbackRequest) (*proto.GenericResponse, error) {
	log.Info("Received Rollback request", "deploymentID", req.DeploymentId)

	// For now, just remove the service as a simple rollback
	err := s.manager.RemoveService(req.DeploymentId)
	if err != nil {
		log.Error("Failed to rollback deployment", "error", err)
		return &proto.GenericResponse{
			Message: fmt.Sprintf("Failed to rollback deployment: %v", err),
			Success: false,
		}, nil
	}

	return &proto.GenericResponse{
		Message: "Deployment rolled back successfully",
		Success: true,
	}, nil
}

// GetStatus handles the GetStatus RPC call
func (s *DeploymentServer) GetStatus(ctx context.Context, req *proto.StatusRequest) (*proto.StatusResponse, error) {
	log.Info("Received GetStatus request", "deploymentID", req.DeploymentId)

	// Get the status of the service
	status, err := s.manager.GetServiceStatus(req.DeploymentId)
	if err != nil {
		log.Error("Failed to get service status", "error", err)
		return nil, fmt.Errorf("failed to get service status: %w", err)
	}

	return &proto.StatusResponse{
		Status: status.State,
		Logs:   status.Logs,
	}, nil
}
