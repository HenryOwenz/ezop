package aws

import (
	"context"
	"fmt"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/HenryOwenz/cloudgate/internal/cloud/aws/codepipeline"
	"github.com/HenryOwenz/cloudgate/internal/cloud/aws/lambda"
)

// Common errors
var (
	ErrNotAuthenticated = fmt.Errorf("not authenticated")
	ErrNotImplemented   = fmt.Errorf("not implemented")
)

// Provider represents the AWS cloud provider.
type Provider struct {
	profile  string
	region   string
	services []cloud.Service
}

// New creates a new AWS provider.
func New() *Provider {
	return &Provider{
		services: make([]cloud.Service, 0),
	}
}

// Name returns the provider's name.
func (p *Provider) Name() string {
	return "AWS"
}

// Description returns the provider's description.
func (p *Provider) Description() string {
	return "Amazon Web Services"
}

// Services returns all available services for this provider.
func (p *Provider) Services() []cloud.Service {
	return p.services
}

// GetProfiles returns all available profiles for this provider.
func (p *Provider) GetProfiles() ([]string, error) {
	return getAWSProfiles()
}

// LoadConfig loads the provider configuration with the given profile and region.
func (p *Provider) LoadConfig(profile, region string) error {
	p.profile = profile
	p.region = region

	// Register services
	p.services = make([]cloud.Service, 0)
	p.services = append(p.services, lambda.NewService(profile, region))
	p.services = append(p.services, codepipeline.NewService(profile, region))

	return nil
}

// GetFunctionStatusOperation returns the function status operation
func (p *Provider) GetFunctionStatusOperation() (cloud.FunctionStatusOperation, error) {
	if p.profile == "" || p.region == "" {
		return nil, ErrNotAuthenticated
	}
	return lambda.NewFunctionStatusOperation(p.profile, p.region), nil
}

// GetCodePipelineManualApprovalOperation returns the CodePipeline manual approval operation
func (p *Provider) GetCodePipelineManualApprovalOperation() (cloud.CodePipelineManualApprovalOperation, error) {
	if p.profile == "" || p.region == "" {
		return nil, ErrNotAuthenticated
	}
	return codepipeline.NewCloudManualApprovalOperation(p.profile, p.region), nil
}

// GetPipelineStatusOperation returns the pipeline status operation
func (p *Provider) GetPipelineStatusOperation() (cloud.PipelineStatusOperation, error) {
	if p.profile == "" || p.region == "" {
		return nil, ErrNotAuthenticated
	}
	return codepipeline.NewCloudPipelineStatusOperation(p.profile, p.region), nil
}

// GetStartPipelineOperation returns the start pipeline operation
func (p *Provider) GetStartPipelineOperation() (cloud.StartPipelineOperation, error) {
	if p.profile == "" || p.region == "" {
		return nil, ErrNotAuthenticated
	}
	return codepipeline.NewCloudStartPipelineOperation(p.profile, p.region), nil
}

// GetAuthenticationMethods returns the available authentication methods
func (p *Provider) GetAuthenticationMethods() []string {
	// AWS only supports profile-based authentication for now
	return []string{"profile"}
}

// GetAuthConfigKeys returns the configuration keys required for an authentication method
func (p *Provider) GetAuthConfigKeys(method string) []string {
	return []string{"profile", "region"}
}

// Authenticate authenticates with the provider using the given method and configuration
func (p *Provider) Authenticate(method string, authConfig map[string]string) error {
	profile, ok := authConfig["profile"]
	if !ok {
		return fmt.Errorf("profile is required")
	}

	region, ok := authConfig["region"]
	if !ok {
		return fmt.Errorf("region is required")
	}

	return p.LoadConfig(profile, region)
}

// IsAuthenticated returns whether the provider is authenticated
func (p *Provider) IsAuthenticated() bool {
	return p.profile != "" && p.region != ""
}

// GetConfigKeys returns the configuration keys required by this provider
func (p *Provider) GetConfigKeys() []string {
	return []string{"profile", "region"}
}

// GetConfigOptions returns the available options for a configuration key
func (p *Provider) GetConfigOptions(key string) ([]string, error) {
	switch key {
	case "profile":
		return p.GetProfiles()
	case "region":
		// Return a list of AWS regions
		return []string{
			"us-east-1", "us-east-2", "us-west-1", "us-west-2",
			"eu-west-1", "eu-west-2", "eu-west-3",
			"eu-central-1",
			"ap-northeast-1", "ap-northeast-2", "ap-northeast-3",
			"ap-southeast-1", "ap-southeast-2",
			"ap-south-1",
			"sa-east-1",
			"ca-central-1",
		}, nil
	default:
		return nil, fmt.Errorf("unknown config key: %s", key)
	}
}

// Configure configures the provider with the given configuration
func (p *Provider) Configure(config map[string]string) error {
	profile, ok := config["profile"]
	if !ok {
		return fmt.Errorf("profile is required")
	}

	region, ok := config["region"]
	if !ok {
		return fmt.Errorf("region is required")
	}

	return p.LoadConfig(profile, region)
}

// GetApprovals returns all pending approvals for the provider
func (p *Provider) GetApprovals(ctx context.Context) ([]cloud.ApprovalAction, error) {
	if p.profile == "" || p.region == "" {
		return nil, ErrNotAuthenticated
	}

	approvalOp, err := p.GetCodePipelineManualApprovalOperation()
	if err != nil {
		return nil, err
	}

	return approvalOp.GetPendingApprovals(ctx)
}

// ApproveAction approves or rejects an approval action
func (p *Provider) ApproveAction(ctx context.Context, action cloud.ApprovalAction, approved bool, comment string) error {
	if p.profile == "" || p.region == "" {
		return ErrNotAuthenticated
	}

	approvalOp, err := p.GetCodePipelineManualApprovalOperation()
	if err != nil {
		return err
	}

	return approvalOp.ApproveAction(ctx, action, approved, comment)
}

// GetStatus returns the status of all pipelines
func (p *Provider) GetStatus(ctx context.Context) ([]cloud.PipelineStatus, error) {
	if p.profile == "" || p.region == "" {
		return nil, ErrNotAuthenticated
	}

	statusOp, err := p.GetPipelineStatusOperation()
	if err != nil {
		return nil, err
	}

	return statusOp.GetPipelineStatus(ctx)
}

// StartPipeline starts a pipeline execution
func (p *Provider) StartPipeline(ctx context.Context, pipelineName string, commitID string) error {
	if p.profile == "" || p.region == "" {
		return ErrNotAuthenticated
	}

	startOp, err := p.GetStartPipelineOperation()
	if err != nil {
		return err
	}

	return startOp.StartPipelineExecution(ctx, pipelineName, commitID)
}
