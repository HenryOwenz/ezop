package aws

import (
	"fmt"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/HenryOwenz/cloudgate/internal/cloud/aws/codepipeline"
	"github.com/HenryOwenz/cloudgate/internal/cloud/aws/lambda"
)

// Common errors
var (
	ErrNotAuthenticated = fmt.Errorf("not authenticated")
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
