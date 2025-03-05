package codepipeline

import (
	"context"
	"fmt"
	"time"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	cpTypes "github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
)

// CloudManualApprovalOperation represents an operation to manage manual approvals in CodePipeline.
type CloudManualApprovalOperation struct {
	profile string
	region  string
}

// NewCloudManualApprovalOperation creates a new manual approval operation.
func NewCloudManualApprovalOperation(profile, region string) *CloudManualApprovalOperation {
	return &CloudManualApprovalOperation{
		profile: profile,
		region:  region,
	}
}

// Name returns the operation's name.
func (o *CloudManualApprovalOperation) Name() string {
	return "Pipeline Approvals"
}

// Description returns the operation's description.
func (o *CloudManualApprovalOperation) Description() string {
	return "Manage Pipeline Approvals"
}

// IsUIVisible returns whether this operation should be visible in the UI.
func (o *CloudManualApprovalOperation) IsUIVisible() bool {
	return true
}

// Execute executes the operation with the given parameters.
func (o *CloudManualApprovalOperation) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	return o.GetPendingApprovals(ctx)
}

// GetPendingApprovals returns all pending manual approval actions.
func (o *CloudManualApprovalOperation) GetPendingApprovals(ctx context.Context) ([]cloud.ApprovalAction, error) {
	// Create a new AWS SDK client
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(o.profile),
		config.WithRegion(o.region),
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

	var approvals []cloud.ApprovalAction

	// Get approvals for each pipeline
	for _, pipeline := range pipelineOutput.Pipelines {
		// Get pipeline details
		pipelineResp, err := client.GetPipeline(ctx, &codepipeline.GetPipelineInput{
			Name: pipeline.Name,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get pipeline details: %w", err)
		}

		// Get pipeline state
		stateResp, err := client.GetPipelineState(ctx, &codepipeline.GetPipelineStateInput{
			Name: pipeline.Name,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get pipeline state: %w", err)
		}

		// Find pending approvals
		pipelineApprovals := findCloudPendingApprovals(*pipeline.Name, pipelineResp.Pipeline.Stages, stateResp.StageStates)
		approvals = append(approvals, pipelineApprovals...)
	}

	return approvals, nil
}

// ApproveAction approves or rejects an approval action.
func (o *CloudManualApprovalOperation) ApproveAction(ctx context.Context, action cloud.ApprovalAction, approved bool, comment string) error {
	// Create a new AWS SDK client
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(o.profile),
		config.WithRegion(o.region),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := codepipeline.NewFromConfig(cfg)

	// Set the approval status
	status := cpTypes.ApprovalStatusRejected
	if approved {
		status = cpTypes.ApprovalStatusApproved
	}

	// Put the approval result
	_, err = client.PutApprovalResult(ctx, &codepipeline.PutApprovalResultInput{
		ActionName:   aws.String(action.ActionName),
		PipelineName: aws.String(action.PipelineName),
		StageName:    aws.String(action.StageName),
		Result: &cpTypes.ApprovalResult{
			Status:  status,
			Summary: aws.String(comment),
		},
		Token: aws.String(action.Token),
	})
	if err != nil {
		return fmt.Errorf("failed to put approval result: %w", err)
	}

	return nil
}

// CloudPipelineStatusOperation represents an operation to view pipeline status.
type CloudPipelineStatusOperation struct {
	profile string
	region  string
}

// NewCloudPipelineStatusOperation creates a new pipeline status operation.
func NewCloudPipelineStatusOperation(profile, region string) *CloudPipelineStatusOperation {
	return &CloudPipelineStatusOperation{
		profile: profile,
		region:  region,
	}
}

// Name returns the operation's name.
func (o *CloudPipelineStatusOperation) Name() string {
	return "Pipeline Status"
}

// Description returns the operation's description.
func (o *CloudPipelineStatusOperation) Description() string {
	return "View Pipeline Status"
}

// IsUIVisible returns whether this operation should be visible in the UI.
func (o *CloudPipelineStatusOperation) IsUIVisible() bool {
	return true
}

// Execute executes the operation with the given parameters.
func (o *CloudPipelineStatusOperation) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	return o.GetPipelineStatus(ctx)
}

// GetPipelineStatus returns the status of all pipelines.
func (o *CloudPipelineStatusOperation) GetPipelineStatus(ctx context.Context) ([]cloud.PipelineStatus, error) {
	// Create a new AWS SDK client
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(o.profile),
		config.WithRegion(o.region),
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

	var pipelineStatuses []cloud.PipelineStatus

	// Get status for each pipeline
	for _, pipeline := range pipelineOutput.Pipelines {
		// Get pipeline state
		stateOutput, err := client.GetPipelineState(ctx, &codepipeline.GetPipelineStateInput{
			Name: pipeline.Name,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get pipeline state: %w", err)
		}

		// Create cloud pipeline status
		status := cloud.PipelineStatus{
			Name:   *pipeline.Name,
			Stages: make([]cloud.StageStatus, len(stateOutput.StageStates)),
		}

		// Fill in stage statuses
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
			status.Stages[i] = cloud.StageStatus{
				Name:        *stage.StageName,
				Status:      stageStatus,
				LastUpdated: lastUpdated,
			}
		}

		pipelineStatuses = append(pipelineStatuses, status)
	}

	return pipelineStatuses, nil
}

// CloudStartPipelineOperation represents an operation to start a pipeline execution.
type CloudStartPipelineOperation struct {
	profile string
	region  string
}

// NewCloudStartPipelineOperation creates a new start pipeline operation.
func NewCloudStartPipelineOperation(profile, region string) *CloudStartPipelineOperation {
	return &CloudStartPipelineOperation{
		profile: profile,
		region:  region,
	}
}

// Name returns the operation's name.
func (o *CloudStartPipelineOperation) Name() string {
	return "Start Pipeline"
}

// Description returns the operation's description.
func (o *CloudStartPipelineOperation) Description() string {
	return "Start Pipeline Execution"
}

// IsUIVisible returns whether this operation should be visible in the UI.
func (o *CloudStartPipelineOperation) IsUIVisible() bool {
	return true
}

// Execute executes the operation with the given parameters.
func (o *CloudStartPipelineOperation) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// Get pipeline name and commit ID from parameters
	pipelineName, ok := params["pipeline_name"].(string)
	if !ok {
		return nil, fmt.Errorf("pipeline_name parameter is required")
	}

	commitID, _ := params["commit_id"].(string)

	// Start the pipeline
	return nil, o.StartPipelineExecution(ctx, pipelineName, commitID)
}

// StartPipelineExecution starts a pipeline execution.
func (o *CloudStartPipelineOperation) StartPipelineExecution(ctx context.Context, pipelineName, commitID string) error {
	// Create a new AWS SDK client
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(o.profile),
		config.WithRegion(o.region),
	)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := codepipeline.NewFromConfig(cfg)

	// Create the input
	input := &codepipeline.StartPipelineExecutionInput{
		Name: aws.String(pipelineName),
	}

	// Start the pipeline execution
	_, err = client.StartPipelineExecution(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to start pipeline execution: %w", err)
	}

	return nil
}

// findCloudPendingApprovals finds all pending manual approval actions in a pipeline.
func findCloudPendingApprovals(pipelineName string, stages []cpTypes.StageDeclaration, stageStates []cpTypes.StageState) []cloud.ApprovalAction {
	// Build a map of action types for quick lookup
	actionTypes := buildCloudActionTypeMap(stages)

	// Build a map of stage states for quick lookup
	stageStateMap := buildCloudStageStateMap(stageStates)

	// Find pending approvals in each stage
	var approvals []cloud.ApprovalAction
	for _, stage := range stages {
		if state, ok := stageStateMap[*stage.Name]; ok {
			approvals = append(approvals, findCloudStageApprovals(pipelineName, stage, state, actionTypes)...)
		}
	}

	return approvals
}

// buildCloudActionTypeMap builds a map of action types for quick lookup.
func buildCloudActionTypeMap(stages []cpTypes.StageDeclaration) map[string]cpTypes.ActionCategory {
	actionTypes := make(map[string]cpTypes.ActionCategory)
	for _, stage := range stages {
		for _, action := range stage.Actions {
			actionTypes[*action.Name] = action.ActionTypeId.Category
		}
	}
	return actionTypes
}

// buildCloudStageStateMap builds a map of stage states for quick lookup.
func buildCloudStageStateMap(stageStates []cpTypes.StageState) map[string]cpTypes.StageState {
	stageStateMap := make(map[string]cpTypes.StageState)
	for _, state := range stageStates {
		stageStateMap[*state.StageName] = state
	}
	return stageStateMap
}

// findCloudStageApprovals finds all pending manual approval actions in a stage.
func findCloudStageApprovals(pipelineName string, stage cpTypes.StageDeclaration, state cpTypes.StageState, actionTypes map[string]cpTypes.ActionCategory) []cloud.ApprovalAction {
	var approvals []cloud.ApprovalAction
	for _, actionState := range state.ActionStates {
		if actionState.ActionName != nil && isCloudApprovalAction(*actionState.ActionName, actionTypes) {
			// Check if the action is waiting for approval
			if actionState.LatestExecution != nil && actionState.LatestExecution.Status == cpTypes.ActionExecutionStatusInProgress {
				if actionState.LatestExecution.Token != nil {
					approvals = append(approvals, cloud.ApprovalAction{
						PipelineName: pipelineName,
						StageName:    *state.StageName,
						ActionName:   *actionState.ActionName,
						Token:        *actionState.LatestExecution.Token,
					})
				}
			}
		}
	}
	return approvals
}

// isCloudApprovalAction checks if an action is a manual approval action.
func isCloudApprovalAction(actionName string, actionTypes map[string]cpTypes.ActionCategory) bool {
	category, ok := actionTypes[actionName]
	return ok && category == cpTypes.ActionCategoryApproval
}
