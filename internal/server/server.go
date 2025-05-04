package server

import (
	"context"
	"fmt"
	"net"

	"github.com/jasonlovesdoggo/velo/internal/deployment"
	"github.com/jasonlovesdoggo/velo/internal/log"
	"github.com/jasonlovesdoggo/velo/internal/orchestrator/manager"
	"google.golang.org/grpc"
)

// DeploymentServiceServer is the interface that must be implemented by the server
type DeploymentServiceServer interface {
	Deploy(context.Context, *DeployRequest) (*DeployResponse, error)
	Rollback(context.Context, *RollbackRequest) (*GenericResponse, error)
	GetStatus(context.Context, *StatusRequest) (*StatusResponse, error)
}

// DeploymentServer implements the DeploymentServiceServer interface
type DeploymentServer struct {
	manager manager.Manager
	server  *grpc.Server
}

// NewDeploymentServer creates a new DeploymentServer
func NewDeploymentServer(manager manager.Manager) *DeploymentServer {
	return &DeploymentServer{
		manager: manager,
		server:  grpc.NewServer(),
	}
}

// Start starts the gRPC server
func (s *DeploymentServer) Start(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	// Register the deployment service
	RegisterDeploymentServiceServer(s.server, s)

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
func (s *DeploymentServer) Deploy(ctx context.Context, req *DeployRequest) (*DeployResponse, error) {
	log.Info("Received Deploy request", "service", req.ServiceName, "image", req.Image)

	// Convert the request to a ServiceDefinition
	serviceDef := deployment.ServiceDefinition{
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

	return &DeployResponse{
		DeploymentId: deploymentID,
		Status:       "deployed",
	}, nil
}

// Rollback handles the Rollback RPC call
func (s *DeploymentServer) Rollback(ctx context.Context, req *RollbackRequest) (*GenericResponse, error) {
	log.Info("Received Rollback request", "deploymentID", req.DeploymentId)

	// For now, just remove the service as a simple rollback
	err := s.manager.RemoveService(req.DeploymentId)
	if err != nil {
		log.Error("Failed to rollback deployment", "error", err)
		return &GenericResponse{
			Message: fmt.Sprintf("Failed to rollback deployment: %v", err),
			Success: false,
		}, nil
	}

	return &GenericResponse{
		Message: "Deployment rolled back successfully",
		Success: true,
	}, nil
}

// GetStatus handles the GetStatus RPC call
func (s *DeploymentServer) GetStatus(ctx context.Context, req *StatusRequest) (*StatusResponse, error) {
	log.Info("Received GetStatus request", "deploymentID", req.DeploymentId)

	// Get the status of the service
	status, err := s.manager.GetServiceStatus(req.DeploymentId)
	if err != nil {
		log.Error("Failed to get service status", "error", err)
		return nil, fmt.Errorf("failed to get service status: %w", err)
	}

	return &StatusResponse{
		Status: status.State,
		Logs:   status.Logs,
	}, nil
}

// RegisterDeploymentServiceServer registers the DeploymentServiceServer with the gRPC server
func RegisterDeploymentServiceServer(s *grpc.Server, srv DeploymentServiceServer) {
	s.RegisterService(&_DeploymentService_serviceDesc, srv)
}

// _DeploymentService_serviceDesc is the gRPC service descriptor for DeploymentService
var _DeploymentService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "velo.DeploymentService",
	HandlerType: (*DeploymentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Deploy",
			Handler:    _DeploymentService_Deploy_Handler,
		},
		{
			MethodName: "Rollback",
			Handler:    _DeploymentService_Rollback_Handler,
		},
		{
			MethodName: "GetStatus",
			Handler:    _DeploymentService_GetStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/velo.proto",
}

// Handler functions for the DeploymentService methods
func _DeploymentService_Deploy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeployRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeploymentServiceServer).Deploy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/velo.DeploymentService/Deploy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeploymentServiceServer).Deploy(ctx, req.(*DeployRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeploymentService_Rollback_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RollbackRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeploymentServiceServer).Rollback(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/velo.DeploymentService/Rollback",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeploymentServiceServer).Rollback(ctx, req.(*RollbackRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeploymentService_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeploymentServiceServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/velo.DeploymentService/GetStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeploymentServiceServer).GetStatus(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}
