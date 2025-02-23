package aws

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
)

// Provider represents an AWS provider configuration.
type Provider struct {
	client *codepipeline.Client
}

// ApprovalAction represents a pending approval in a pipeline.
type ApprovalAction struct {
	PipelineName string
	StageName    string
	ActionName   string
	Token        string
}

// PipelineStatus represents the status of a pipeline and its stages
type PipelineStatus struct {
	Name   string
	Stages []StageStatus
}

// StageStatus represents the status of a pipeline stage
type StageStatus struct {
	Name        string
	Status      string
	LastUpdated string
}

// Common errors.
var (
	ErrLoadConfig     = errors.New("failed to load AWS config")
	ErrListPipelines  = errors.New("failed to list pipelines")
	ErrGetPipeline    = errors.New("failed to get pipeline details")
	ErrPipelineState  = errors.New("failed to get pipeline state")
	ErrApprovalResult = errors.New("failed to put approval result")
)

// New creates a new AWS provider with the given profile and region.
func New(ctx context.Context, profile string, region string) (*Provider, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrLoadConfig, err)
	}

	return &Provider{
		client: codepipeline.NewFromConfig(cfg),
	}, nil
}

// GetPendingApprovals returns all pending manual approval actions.
func (p *Provider) GetPendingApprovals(ctx context.Context) ([]ApprovalAction, error) {
	pipelines, err := p.listPipelines(ctx)
	if err != nil {
		return nil, err
	}

	var approvals []ApprovalAction

	for _, pipeline := range pipelines {
		pipelineApprovals, err := p.getPipelineApprovals(ctx, pipeline)
		if err != nil {
			return nil, err
		}

		approvals = append(approvals, pipelineApprovals...)
	}

	return approvals, nil
}

// listPipelines returns a list of all pipelines.
func (p *Provider) listPipelines(ctx context.Context) ([]types.PipelineSummary, error) {
	pipelineOutput, err := p.client.ListPipelines(ctx, &codepipeline.ListPipelinesInput{})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrListPipelines, err)
	}

	return pipelineOutput.Pipelines, nil
}

// getPipelineApprovals returns all pending approval actions for a given pipeline.
func (p *Provider) getPipelineApprovals(ctx context.Context, pipeline types.PipelineSummary) ([]ApprovalAction, error) {
	pipelineOutput, err := p.client.GetPipeline(ctx, &codepipeline.GetPipelineInput{
		Name: pipeline.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetPipeline, err)
	}

	stateOutput, err := p.client.GetPipelineState(ctx, &codepipeline.GetPipelineStateInput{
		Name: pipeline.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPipelineState, err)
	}

	return p.findPendingApprovals(*pipeline.Name, pipelineOutput.Pipeline.Stages, stateOutput.StageStates), nil
}

// findPendingApprovals returns a list of pending approval actions from the given stages and their states.
func (p *Provider) findPendingApprovals(pipelineName string, stages []types.StageDeclaration, stageStates []types.StageState) []ApprovalAction {
	var approvals []ApprovalAction
	actionTypes := p.buildActionTypeMap(stages)
	stateMap := p.buildStageStateMap(stageStates)

	for _, stage := range stages {
		if state, ok := stateMap[*stage.Name]; ok {
			approvals = append(approvals, p.findStageApprovals(pipelineName, stage, state, actionTypes)...)
		}
	}

	return approvals
}

// buildActionTypeMap creates a map of action names to their categories for quick lookup.
func (p *Provider) buildActionTypeMap(stages []types.StageDeclaration) map[string]types.ActionCategory {
	actionTypes := make(map[string]types.ActionCategory)
	for _, stage := range stages {
		for _, action := range stage.Actions {
			actionTypes[*action.Name] = action.ActionTypeId.Category
		}
	}

	return actionTypes
}

// buildStageStateMap creates a map of stage names to their states for quick lookup.
func (p *Provider) buildStageStateMap(stageStates []types.StageState) map[string]types.StageState {
	stateMap := make(map[string]types.StageState)
	for _, state := range stageStates {
		stateMap[*state.StageName] = state
	}

	return stateMap
}

// findStageApprovals returns a list of pending approval actions from a single stage.
func (p *Provider) findStageApprovals(pipelineName string, stage types.StageDeclaration, state types.StageState, actionTypes map[string]types.ActionCategory) []ApprovalAction {
	var approvals []ApprovalAction
	for _, actionState := range state.ActionStates {
		if p.isApprovalAction(actionState, actionTypes) {
			approval := ApprovalAction{
				PipelineName: pipelineName,
				StageName:    *stage.Name,
				ActionName:   *actionState.ActionName,
				Token:        *actionState.LatestExecution.Token,
			}
			approvals = append(approvals, approval)
		}
	}

	return approvals
}

// isApprovalAction checks if the given action state represents a pending manual approval.
func (p *Provider) isApprovalAction(actionState types.ActionState, actionTypes map[string]types.ActionCategory) bool {
	return actionState.LatestExecution != nil &&
		actionState.LatestExecution.Token != nil &&
		actionState.LatestExecution.Status == types.ActionExecutionStatusInProgress &&
		actionTypes[*actionState.ActionName] == types.ActionCategoryApproval
}

// PutApprovalResult handles the approval or rejection of a manual approval action.
func (p *Provider) PutApprovalResult(ctx context.Context, action ApprovalAction, approved bool, comment string) error {
	status := types.ApprovalStatusApproved
	if !approved {
		status = types.ApprovalStatusRejected
	}

	_, err := p.client.PutApprovalResult(ctx, &codepipeline.PutApprovalResultInput{
		ActionName:   aws.String(action.ActionName),
		PipelineName: aws.String(action.PipelineName),
		Result: &types.ApprovalResult{
			Status:  status,
			Summary: aws.String(comment),
		},
		StageName: aws.String(action.StageName),
		Token:     aws.String(action.Token),
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrApprovalResult, err)
	}

	return nil
}

// GetProfiles returns a list of available AWS profiles.
func GetProfiles() []string {
	// Get user's home directory
	home, err := os.UserHomeDir()
	if err != nil {
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

	// If no profiles found, return default
	if len(profiles) == 0 {
		return []string{"default"}
	}

	sort.Strings(profiles)
	return profiles
}

// GetPipelineStatus returns the status of all pipelines
func (p *Provider) GetPipelineStatus(ctx context.Context) ([]PipelineStatus, error) {
	pipelines, err := p.listPipelines(ctx)
	if err != nil {
		return nil, err
	}

	var pipelineStatuses []PipelineStatus

	for _, pipeline := range pipelines {
		status, err := p.getPipelineStatus(ctx, pipeline)
		if err != nil {
			return nil, err
		}
		pipelineStatuses = append(pipelineStatuses, status)
	}

	return pipelineStatuses, nil
}

// getPipelineStatus returns the status of a single pipeline
func (p *Provider) getPipelineStatus(ctx context.Context, pipeline types.PipelineSummary) (PipelineStatus, error) {
	stateOutput, err := p.client.GetPipelineState(ctx, &codepipeline.GetPipelineStateInput{
		Name: pipeline.Name,
	})
	if err != nil {
		return PipelineStatus{}, fmt.Errorf("%w: %w", ErrPipelineState, err)
	}

	status := PipelineStatus{
		Name:   *pipeline.Name,
		Stages: make([]StageStatus, len(stateOutput.StageStates)),
	}

	for i, stage := range stateOutput.StageStates {
		stageStatus := "Unknown"
		lastUpdated := "N/A"
		if stage.LatestExecution != nil {
			stageStatus = string(stage.LatestExecution.Status)
			if stage.ActionStates != nil && len(stage.ActionStates) > 0 {
				// Find the most recent action update time
				var latestTime *time.Time
				for _, action := range stage.ActionStates {
					if action.LatestExecution != nil && action.LatestExecution.LastStatusChange != nil {
						if latestTime == nil || action.LatestExecution.LastStatusChange.After(*latestTime) {
							latestTime = action.LatestExecution.LastStatusChange
						}
					}
				}
				if latestTime != nil {
					lastUpdated = latestTime.UTC().Format("Jan 02 15:04:05") + " UTC"
				}
			}
		}
		status.Stages[i] = StageStatus{
			Name:        *stage.StageName,
			Status:      stageStatus,
			LastUpdated: lastUpdated,
		}
	}

	return status, nil
}
