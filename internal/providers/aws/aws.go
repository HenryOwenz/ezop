package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
)

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

// GetPendingApprovals returns all pending manual approval actions.
func (p *Provider) GetPendingApprovals(ctx context.Context) ([]ApprovalAction, error) {
	// Create a new AWS SDK client
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(p.profile),
		config.WithRegion(p.region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := codepipeline.NewFromConfig(cfg)

	// List all pipelines
	pipelineOutput, err := client.ListPipelines(ctx, &codepipeline.ListPipelinesInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pipelines: %w", err)
	}

	var approvals []ApprovalAction

	// Check each pipeline for pending approvals
	for _, pipeline := range pipelineOutput.Pipelines {
		// Get pipeline details
		pipelineOutput, err := client.GetPipeline(ctx, &codepipeline.GetPipelineInput{
			Name: pipeline.Name,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get pipeline details: %w", err)
		}

		// Get pipeline state
		stateOutput, err := client.GetPipelineState(ctx, &codepipeline.GetPipelineStateInput{
			Name: pipeline.Name,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get pipeline state: %w", err)
		}

		// Find pending approvals
		pipelineApprovals := findPendingApprovals(*pipeline.Name, pipelineOutput.Pipeline.Stages, stateOutput.StageStates)
		approvals = append(approvals, pipelineApprovals...)
	}

	return approvals, nil
}

// findPendingApprovals returns a list of pending approval actions from the given stages and their states.
func findPendingApprovals(pipelineName string, stages []types.StageDeclaration, stageStates []types.StageState) []ApprovalAction {
	var approvals []ApprovalAction
	actionTypes := buildActionTypeMap(stages)
	stateMap := buildStageStateMap(stageStates)

	for _, stage := range stages {
		if state, ok := stateMap[*stage.Name]; ok {
			approvals = append(approvals, findStageApprovals(pipelineName, stage, state, actionTypes)...)
		}
	}

	return approvals
}

// buildActionTypeMap creates a map of action names to their categories for quick lookup.
func buildActionTypeMap(stages []types.StageDeclaration) map[string]types.ActionCategory {
	actionTypes := make(map[string]types.ActionCategory)
	for _, stage := range stages {
		for _, action := range stage.Actions {
			actionTypes[*action.Name] = action.ActionTypeId.Category
		}
	}

	return actionTypes
}

// buildStageStateMap creates a map of stage names to their states for quick lookup.
func buildStageStateMap(stageStates []types.StageState) map[string]types.StageState {
	stateMap := make(map[string]types.StageState)
	for _, state := range stageStates {
		stateMap[*state.StageName] = state
	}

	return stateMap
}

// findStageApprovals returns a list of pending approval actions from a single stage.
func findStageApprovals(pipelineName string, stage types.StageDeclaration, state types.StageState, actionTypes map[string]types.ActionCategory) []ApprovalAction {
	var approvals []ApprovalAction
	for _, actionState := range state.ActionStates {
		if isApprovalAction(actionState, actionTypes) {
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
func isApprovalAction(actionState types.ActionState, actionTypes map[string]types.ActionCategory) bool {
	return actionState.LatestExecution != nil &&
		actionState.LatestExecution.Token != nil &&
		actionState.LatestExecution.Status == types.ActionExecutionStatusInProgress &&
		actionTypes[*actionState.ActionName] == types.ActionCategoryApproval
}

// PutApprovalResult handles the approval or rejection of a manual approval action.
func (p *Provider) PutApprovalResult(ctx context.Context, action ApprovalAction, approved bool, comment string) error {
	// Create a new AWS SDK client
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(p.profile),
		config.WithRegion(p.region),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := codepipeline.NewFromConfig(cfg)

	status := types.ApprovalStatusApproved
	if !approved {
		status = types.ApprovalStatusRejected
	}

	_, err = client.PutApprovalResult(ctx, &codepipeline.PutApprovalResultInput{
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
		return fmt.Errorf("failed to put approval result: %w", err)
	}

	return nil
}

// GetPipelineStatus returns the status of all pipelines
func (p *Provider) GetPipelineStatus(ctx context.Context) ([]PipelineStatus, error) {
	// Create a new AWS SDK client
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(p.profile),
		config.WithRegion(p.region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := codepipeline.NewFromConfig(cfg)

	// List all pipelines
	pipelineOutput, err := client.ListPipelines(ctx, &codepipeline.ListPipelinesInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pipelines: %w", err)
	}

	var pipelineStatuses []PipelineStatus

	// Get status for each pipeline
	for _, pipeline := range pipelineOutput.Pipelines {
		status, err := getPipelineStatus(ctx, client, pipeline)
		if err != nil {
			return nil, err
		}
		pipelineStatuses = append(pipelineStatuses, status)
	}

	return pipelineStatuses, nil
}

// getPipelineStatus returns the status of a single pipeline
func getPipelineStatus(ctx context.Context, client *codepipeline.Client, pipeline types.PipelineSummary) (PipelineStatus, error) {
	stateOutput, err := client.GetPipelineState(ctx, &codepipeline.GetPipelineStateInput{
		Name: pipeline.Name,
	})
	if err != nil {
		return PipelineStatus{}, fmt.Errorf("failed to get pipeline state: %w", err)
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
			if len(stage.ActionStates) > 0 {
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

// StartPipelineExecution starts a pipeline execution with optional source revision
func (p *Provider) StartPipelineExecution(ctx context.Context, pipelineName string, commitID string) error {
	// Create a new AWS SDK client
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(p.profile),
		config.WithRegion(p.region),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := codepipeline.NewFromConfig(cfg)

	input := &codepipeline.StartPipelineExecutionInput{
		Name: aws.String(pipelineName),
	}

	// Only add source revision if a specific commitID is provided
	if commitID != "" {
		input.SourceRevisions = []types.SourceRevisionOverride{
			{
				ActionName:    aws.String("Source"), // Assuming standard Source action name
				RevisionType:  types.SourceRevisionTypeCommitId,
				RevisionValue: aws.String(commitID),
			},
		}
	}

	_, err = client.StartPipelineExecution(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to start pipeline execution: %w", err)
	}

	return nil
}
