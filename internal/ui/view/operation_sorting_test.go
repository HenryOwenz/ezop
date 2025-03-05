package view

import (
	"context"
	"testing"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
)

// MockOperation implements the cloud.Operation interface for testing
type MockOperation struct {
	name        string
	description string
}

func (o *MockOperation) Name() string {
	return o.name
}

func (o *MockOperation) Description() string {
	return o.description
}

func (o *MockOperation) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	return nil, nil
}

func (o *MockOperation) IsUIVisible() bool {
	return true
}

// MockCategory implements the cloud.Category interface for testing
type MockCategory struct {
	name        string
	description string
	operations  []cloud.Operation
}

func (c *MockCategory) Name() string {
	return c.name
}

func (c *MockCategory) Description() string {
	return c.description
}

func (c *MockCategory) Operations() []cloud.Operation {
	return c.operations
}

func (c *MockCategory) IsUIVisible() bool {
	return true
}

// MockService implements the cloud.Service interface for testing
type MockService struct {
	name        string
	description string
	categories  []cloud.Category
}

func (s *MockService) Name() string {
	return s.name
}

func (s *MockService) Description() string {
	return s.description
}

func (s *MockService) Categories() []cloud.Category {
	return s.categories
}

// MockProvider implements the cloud.Provider interface for testing
type MockProvider struct {
	name        string
	description string
	services    []cloud.Service
}

func (p *MockProvider) Name() string {
	return p.name
}

func (p *MockProvider) Description() string {
	return p.description
}

func (p *MockProvider) Services() []cloud.Service {
	return p.services
}

// Implement other required methods with minimal functionality
func (p *MockProvider) GetProfiles() ([]string, error) {
	return []string{}, nil
}

func (p *MockProvider) LoadConfig(profile, region string) error {
	return nil
}

func (p *MockProvider) GetFunctionStatusOperation() (cloud.FunctionStatusOperation, error) {
	return nil, nil
}

func (p *MockProvider) GetCodePipelineManualApprovalOperation() (cloud.CodePipelineManualApprovalOperation, error) {
	return nil, nil
}

func (p *MockProvider) GetPipelineStatusOperation() (cloud.PipelineStatusOperation, error) {
	return nil, nil
}

func (p *MockProvider) GetStartPipelineOperation() (cloud.StartPipelineOperation, error) {
	return nil, nil
}

func (p *MockProvider) GetAuthenticationMethods() []string {
	return []string{}
}

func (p *MockProvider) GetAuthConfigKeys(method string) []string {
	return []string{}
}

func (p *MockProvider) Authenticate(method string, authConfig map[string]string) error {
	return nil
}

func (p *MockProvider) IsAuthenticated() bool {
	return true
}

func (p *MockProvider) GetConfigKeys() []string {
	return []string{}
}

func (p *MockProvider) GetConfigOptions(key string) ([]string, error) {
	return []string{}, nil
}

func (p *MockProvider) Configure(config map[string]string) error {
	return nil
}

func (p *MockProvider) GetApprovals(ctx context.Context) ([]cloud.ApprovalAction, error) {
	return []cloud.ApprovalAction{}, nil
}

func (p *MockProvider) ApproveAction(ctx context.Context, action cloud.ApprovalAction, approved bool, comment string) error {
	return nil
}

func (p *MockProvider) GetStatus(ctx context.Context) ([]cloud.PipelineStatus, error) {
	return []cloud.PipelineStatus{}, nil
}

func (p *MockProvider) StartPipeline(ctx context.Context, pipelineName string, commitID string) error {
	return nil
}

// TestOperationSorting verifies that operations are sorted alphabetically
func TestOperationSorting(t *testing.T) {
	// Create mock operations in unsorted order
	operations := []cloud.Operation{
		&MockOperation{name: "Start Pipeline", description: "Start Pipeline Execution"},
		&MockOperation{name: "Pipeline Approvals", description: "Manage Pipeline Approvals"},
		&MockOperation{name: "Pipeline Status", description: "View Pipeline Status"},
	}

	// Create a mock category with the operations
	category := &MockCategory{
		name:        "Workflows",
		description: "CodePipeline Workflows",
		operations:  operations,
	}

	// Create a mock service with the category
	service := &MockService{
		name:        "CodePipeline",
		description: "AWS CodePipeline",
		categories:  []cloud.Category{category},
	}

	// Create a mock provider with the service
	provider := &MockProvider{
		name:        "AWS",
		description: "Amazon Web Services",
		services:    []cloud.Service{service},
	}

	// Create a mock registry
	registry := cloud.NewProviderRegistry()
	registry.Register(provider)

	// Create a model with the registry
	m := model.New()
	m.Registry = registry
	m.CurrentView = constants.ViewSelectOperation
	m.SelectedService = &model.Service{
		Name:        "CodePipeline",
		Description: "AWS CodePipeline",
	}
	m.SelectedCategory = &model.Category{
		Name:        "Workflows",
		Description: "CodePipeline Workflows",
	}

	// Get the rows for the view
	rows := getRowsForView(m)

	// Verify that the operations are sorted alphabetically
	if len(rows) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(rows))
	}

	// Expected order: Pipeline Approvals, Pipeline Status, Start Pipeline
	if rows[0][0] != "Pipeline Approvals" {
		t.Errorf("Expected first operation to be 'Pipeline Approvals', got '%s'", rows[0][0])
	}
	if rows[1][0] != "Pipeline Status" {
		t.Errorf("Expected second operation to be 'Pipeline Status', got '%s'", rows[1][0])
	}
	if rows[2][0] != "Start Pipeline" {
		t.Errorf("Expected third operation to be 'Start Pipeline', got '%s'", rows[2][0])
	}
}
