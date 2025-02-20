package aws

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
)

// Provider represents an AWS provider configuration
type Provider struct {
	Profile string
	Region  string
	client  *codepipeline.Client
}

// ApprovalAction represents a pending approval in a pipeline
type ApprovalAction struct {
	PipelineName string
	StageName    string
	ActionName   string
	Token        string
}

// New creates a new AWS provider with the given profile and region
func New(profile, region string) (*Provider, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %v", err)
	}

	return &Provider{
		Profile: profile,
		Region:  region,
		client:  codepipeline.NewFromConfig(cfg),
	}, nil
}

// GetPendingApprovals returns all pending manual approval actions
func (p *Provider) GetPendingApprovals(ctx context.Context) ([]ApprovalAction, error) {
	// List all pipelines first
	pipelineOutput, err := p.client.ListPipelines(ctx, &codepipeline.ListPipelinesInput{})
	if err != nil {
		return nil, fmt.Errorf("unable to list pipelines: %v", err)
	}

	var approvals []ApprovalAction
	for _, pipeline := range pipelineOutput.Pipelines {
		pipelineName := *pipeline.Name

		// Get pipeline definition to check action types
		pipelineOutput, err := p.client.GetPipeline(ctx, &codepipeline.GetPipelineInput{
			Name: pipeline.Name,
		})
		if err != nil {
			continue
		}

		// Create a map of action names to their types
		actionTypes := make(map[string]types.ActionCategory)
		for _, stage := range pipelineOutput.Pipeline.Stages {
			for _, action := range stage.Actions {
				if action.Name != nil && action.ActionTypeId != nil {
					actionTypes[*action.Name] = action.ActionTypeId.Category
				}
			}
		}

		// Get pipeline state
		stateOutput, err := p.client.GetPipelineState(ctx, &codepipeline.GetPipelineStateInput{
			Name: pipeline.Name,
		})
		if err != nil {
			continue
		}

		for _, stageState := range stateOutput.StageStates {
			if stageState.StageName == nil {
				continue
			}

			for _, actionState := range stageState.ActionStates {
				if actionState.ActionName == nil || actionState.LatestExecution == nil {
					continue
				}

				// Only include actions that are:
				// 1. In progress
				// 2. Have an approval token
				// 3. Are of type Approval
				if actionState.LatestExecution.Status == types.ActionExecutionStatusInProgress &&
					actionState.LatestExecution.Token != nil &&
					actionTypes[*actionState.ActionName] == types.ActionCategoryApproval {

					approvals = append(approvals, ApprovalAction{
						PipelineName: pipelineName,
						StageName:    *stageState.StageName,
						ActionName:   *actionState.ActionName,
						Token:        *actionState.LatestExecution.Token,
					})
				}
			}
		}
	}

	return approvals, nil
}

// HandleApproval handles the approval or rejection of a manual approval action
func (p *Provider) HandleApproval(ctx context.Context, action *ApprovalAction, approve bool, summary string) error {
	status := types.ApprovalStatusRejected
	if approve {
		status = types.ApprovalStatusApproved
	}

	input := &codepipeline.PutApprovalResultInput{
		PipelineName: &action.PipelineName,
		StageName:    &action.StageName,
		ActionName:   &action.ActionName,
		Result: &types.ApprovalResult{
			Summary: &summary,
			Status:  status,
		},
		Token: &action.Token,
	}

	_, err := p.client.PutApprovalResult(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to put approval result: %v", err)
	}
	return nil
}

// GetProfiles returns a list of available AWS profiles
func GetProfiles() []string {
	// Get user's home directory
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println("Unable to get user home directory, using default profile")
		return []string{"default"}
	}

	// Try both config and credentials files
	configFiles := []string{
		filepath.Join(home, ".aws", "config"),
		filepath.Join(home, ".aws", "credentials"),
	}

	var profiles []string
	profileMap := make(map[string]bool)

	for _, file := range configFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		// Parse profiles using regex
		re := regexp.MustCompile(`\[(.*?)\]`)
		matches := re.FindAllStringSubmatch(string(content), -1)
		for _, match := range matches {
			profile := strings.TrimSpace(match[1])
			// Remove "profile " prefix if present (used in config file)
			profile = strings.TrimPrefix(profile, "profile ")
			if profile != "" && !profileMap[profile] {
				profileMap[profile] = true
				profiles = append(profiles, profile)
			}
		}
	}

	if len(profiles) == 0 {
		log.Println("No AWS profiles found, using default")
		return []string{"default"}
	}

	// Sort profiles for consistent display
	sort.Strings(profiles)
	return profiles
}
