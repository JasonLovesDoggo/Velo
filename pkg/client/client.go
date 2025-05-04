package client

import (
	"context"
	"net"
	"time"

	"github.com/jasonlovesdoggo/velo/api/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client represents a client for the Velo API
type Client struct {
	conn   *grpc.ClientConn
	client proto.DeploymentServiceClient
}

// NewClientWithConn creates a new client with an existing connection (for testing)
func NewClientWithConn(conn *grpc.ClientConn) *Client {
	return &Client{
		conn:   conn,
		client: proto.NewDeploymentServiceClient(conn),
	}
}

// NewClient creates a new client for the Velo API
func NewClient(serverAddr string) (*Client, error) {
	// Set up a connection to the server
	conn, err := grpc.NewClient(
		serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			d := &net.Dialer{}
			return d.DialContext(ctx, "tcp", addr)
		}),
	)
	if err != nil {
		return nil, err
	}

	// Create a client using the generated client interface
	client := proto.NewDeploymentServiceClient(conn)

	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

// Close closes the client connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Deploy deploys a service
func (c *Client) Deploy(ctx context.Context, serviceName, image string, env map[string]string) (*proto.DeployResponse, error) {
	// Create a deploy request
	req := &proto.DeployRequest{
		ServiceName: serviceName,
		Image:       image,
		Env:         env,
	}

	// Call the Deploy method
	return c.client.Deploy(ctx, req)
}

// GetStatus gets the status of a deployment
func (c *Client) GetStatus(ctx context.Context, deploymentID string) (*proto.StatusResponse, error) {
	// Create a status request
	req := &proto.StatusRequest{
		DeploymentId: deploymentID,
	}

	// Call the GetStatus method
	return c.client.GetStatus(ctx, req)
}

// Rollback rolls back a deployment
func (c *Client) Rollback(ctx context.Context, deploymentID string) (*proto.GenericResponse, error) {
	// Create a rollback request
	req := &proto.RollbackRequest{
		DeploymentId: deploymentID,
	}

	// Call the Rollback method
	return c.client.Rollback(ctx, req)
}

// WithTimeout creates a new context with a timeout
func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, timeout)
}
