package server

import (
	"context"
	"errors"
	"testing"

	"github.com/jasonlovesdoggo/velo/internal/config"
	"github.com/jasonlovesdoggo/velo/internal/orchestrator/manager"
)

// MockManager is a mock implementation of the Manager interface for testing
type MockManager struct {
	// Mock return values
	DeployServiceID  string
	DeployServiceErr error
	RemoveServiceErr error
	ServiceStatus    config.DeploymentStatus
	ServiceStatusErr error
}

// Ensure MockManager implements manager.Manager
var _ manager.Manager = (*MockManager)(nil)

// DeployService mocks the Manager's DeployService method
func (m *MockManager) DeployService(def config.ServiceDefinition) (string, error) {
	return m.DeployServiceID, m.DeployServiceErr
}

// RemoveService mocks the Manager's RemoveService method
func (m *MockManager) RemoveService(serviceID string) error {
	return m.RemoveServiceErr
}

// GetServiceStatus mocks the Manager's GetServiceStatus method
func (m *MockManager) GetServiceStatus(serviceID string) (config.DeploymentStatus, error) {
	return m.ServiceStatus, m.ServiceStatusErr
}

func TestDeploy(t *testing.T) {
	tests := []struct {
		name           string
		req            *DeployRequest
		mockID         string
		mockErr        error
		expectedID     string
		expectedStatus string
		expectError    bool
	}{
		{
			name: "Successful deployment",
			req: &DeployRequest{
				ServiceName: "test-service",
				Image:       "nginx:latest",
				Env:         map[string]string{"ENV": "test"},
			},
			mockID:         "service-123",
			mockErr:        nil,
			expectedID:     "service-123",
			expectedStatus: "deployed",
			expectError:    false,
		},
		{
			name: "Failed deployment",
			req: &DeployRequest{
				ServiceName: "test-service",
				Image:       "nginx:latest",
				Env:         map[string]string{"ENV": "test"},
			},
			mockID:      "",
			mockErr:     errors.New("deployment failed"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock Manager
			mockManager := &MockManager{
				DeployServiceID:  tt.mockID,
				DeployServiceErr: tt.mockErr,
			}

			// Create a server with the mock manager
			server := NewDeploymentServer(mockManager)

			// Call the Deploy method
			resp, err := server.Deploy(context.Background(), tt.req)

			// Check the error
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check the response
			if resp.DeploymentId != tt.expectedID {
				t.Errorf("Expected deployment ID %q, got %q", tt.expectedID, resp.DeploymentId)
			}

			if resp.Status != tt.expectedStatus {
				t.Errorf("Expected status %q, got %q", tt.expectedStatus, resp.Status)
			}
		})
	}
}

func TestRollback(t *testing.T) {
	tests := []struct {
		name          string
		req           *RollbackRequest
		mockErr       error
		expectSuccess bool
	}{
		{
			name:          "Successful rollback",
			req:           &RollbackRequest{DeploymentId: "service-123"},
			mockErr:       nil,
			expectSuccess: true,
		},
		{
			name:          "Failed rollback",
			req:           &RollbackRequest{DeploymentId: "service-123"},
			mockErr:       errors.New("rollback failed"),
			expectSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock Manager
			mockManager := &MockManager{
				RemoveServiceErr: tt.mockErr,
			}

			// Create a server with the mock manager
			server := NewDeploymentServer(mockManager)

			// Call the Rollback method
			resp, err := server.Rollback(context.Background(), tt.req)

			// We don't expect errors from the Rollback method itself
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check the response
			if resp.Success != tt.expectSuccess {
				t.Errorf("Expected success %v, got %v", tt.expectSuccess, resp.Success)
			}
		})
	}
}

func TestGetStatus(t *testing.T) {
	tests := []struct {
		name           string
		req            *StatusRequest
		mockStatus     config.DeploymentStatus
		mockErr        error
		expectedStatus string
		expectedLogs   string
		expectError    bool
	}{
		{
			name: "Successful status retrieval",
			req:  &StatusRequest{DeploymentId: "service-123"},
			mockStatus: config.DeploymentStatus{
				ID:    "service-123",
				State: "running",
				Logs:  "Service is running",
			},
			mockErr:        nil,
			expectedStatus: "running",
			expectedLogs:   "Service is running",
			expectError:    false,
		},
		{
			name:        "Failed status retrieval",
			req:         &StatusRequest{DeploymentId: "service-123"},
			mockStatus:  config.DeploymentStatus{},
			mockErr:     errors.New("status retrieval failed"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock Manager
			mockManager := &MockManager{
				ServiceStatus:    tt.mockStatus,
				ServiceStatusErr: tt.mockErr,
			}

			// Create a server with the mock manager
			server := NewDeploymentServer(mockManager)

			// Call the GetStatus method
			resp, err := server.GetStatus(context.Background(), tt.req)

			// Check the error
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected an error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Check the response
			if resp.Status != tt.expectedStatus {
				t.Errorf("Expected status %q, got %q", tt.expectedStatus, resp.Status)
			}

			if resp.Logs != tt.expectedLogs {
				t.Errorf("Expected logs %q, got %q", tt.expectedLogs, resp.Logs)
			}
		})
	}
}
