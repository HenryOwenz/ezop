package aws

import (
	"context"
	"fmt"

	"github.com/HenryOwenz/ezop/internal/domain"
)

// Provider implements the CloudProvider interface for AWS
type Provider struct {
	profile      string
	region       string
	codePipeline *CodePipelineService
}

// ApprovalAction represents a pending approval in a pipeline
type ApprovalAction struct {
	PipelineName string
	StageName    string
	ActionName   string
	Token        string
}

// Services available in AWS
var awsServices = []domain.Service{
	{
		ID:          "codepipeline",
		Name:        "CodePipeline",
		Description: "Continuous Delivery Service",
		Available:   true,
	},
	// Placeholder for future AWS services
}

// Operations available for CodePipeline
var codePipelineOperations = []domain.Operation{
	{
		ID:          "manual-approval",
		Name:        "Manual Approval",
		Description: "Manage manual approval actions",
	},
	// Placeholder for future CodePipeline operations
}

// NewProvider creates a new AWS provider instance
func NewProvider(profile, region string) (*Provider, error) {
	codePipeline, err := NewCodePipelineService(profile, region)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CodePipeline service: %w", err)
	}

	return &Provider{
		profile:      profile,
		region:       region,
		codePipeline: codePipeline,
	}, nil
}

// GetServices returns the list of available AWS services
func (p *Provider) GetServices() []domain.Service {
	return awsServices
}

// GetOperations returns the list of available operations for a service
func (p *Provider) GetOperations(serviceID string) []domain.Operation {
	switch serviceID {
	case "codepipeline":
		return codePipelineOperations
	default:
		return nil
	}
}

// ExecuteOperation executes an operation on a service
func (p *Provider) ExecuteOperation(ctx context.Context, serviceID, operationID string, params map[string]interface{}) error {
	switch serviceID {
	case "codepipeline":
		switch operationID {
		case "manual-approval":
			return p.codePipeline.HandleApproval(ctx, params)
		default:
			return fmt.Errorf("operation %s not supported for CodePipeline", operationID)
		}
	default:
		return fmt.Errorf("service %s not supported", serviceID)
	}
}

// GetPendingApprovals returns all pending manual approval actions
func (p *Provider) GetPendingApprovals(ctx context.Context) ([]ApprovalAction, error) {
	return p.codePipeline.ListPendingApprovals(ctx)
}
