package codepipeline

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
)

// Common errors.
var (
	ErrLoadConfig     = errors.New("failed to load AWS config")
	ErrListPipelines  = errors.New("failed to list pipelines")
	ErrGetPipeline    = errors.New("failed to get pipeline details")
	ErrPipelineState  = errors.New("failed to get pipeline state")
	ErrApprovalResult = errors.New("failed to put approval result")
)

// Service represents the CodePipeline service.
type Service struct {
	profile    string
	region     string
	categories []cloud.Category
}

// NewService creates a new CodePipeline service.
func NewService(profile, region string) *Service {
	service := &Service{
		profile:    profile,
		region:     region,
		categories: make([]cloud.Category, 0),
	}

	// Register categories
	service.categories = append(service.categories, NewWorkflowsCategory(profile, region))
	service.categories = append(service.categories, NewInternalOperationsCategory(profile, region))

	return service
}

// Name returns the service's name.
func (s *Service) Name() string {
	return "CodePipeline"
}

// Description returns the service's description.
func (s *Service) Description() string {
	return "Continuous Delivery Service"
}

// Categories returns all available categories for this service.
func (s *Service) Categories() []cloud.Category {
	return s.categories
}

// getClient creates a new CodePipeline client.
func getClient(ctx context.Context, profile, region string) (*codepipeline.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrLoadConfig, err)
	}

	return codepipeline.NewFromConfig(cfg), nil
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

// GetPendingApprovals returns all pending manual approval actions.
func GetPendingApprovals(ctx context.Context, profile, region string) ([]ApprovalAction, error) {
	client, err := getClient(ctx, profile, region)
	if err != nil {
		return nil, err
	}

	pipelines, err := listPipelines(ctx, client)
	if err != nil {
		return nil, err
	}

	var approvals []ApprovalAction

	for _, pipeline := range pipelines {
		pipelineApprovals, err := getPipelineApprovals(ctx, client, pipeline)
		if err != nil {
			return nil, err
		}

		approvals = append(approvals, pipelineApprovals...)
	}

	return approvals, nil
}

// listPipelines returns a list of all pipelines.
func listPipelines(ctx context.Context, client *codepipeline.Client) ([]types.PipelineSummary, error) {
	pipelineOutput, err := client.ListPipelines(ctx, &codepipeline.ListPipelinesInput{})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrListPipelines, err)
	}

	return pipelineOutput.Pipelines, nil
}

// getPipelineApprovals returns all pending approval actions for a given pipeline.
func getPipelineApprovals(ctx context.Context, client *codepipeline.Client, pipeline types.PipelineSummary) ([]ApprovalAction, error) {
	pipelineOutput, err := client.GetPipeline(ctx, &codepipeline.GetPipelineInput{
		Name: pipeline.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetPipeline, err)
	}

	stateOutput, err := client.GetPipelineState(ctx, &codepipeline.GetPipelineStateInput{
		Name: pipeline.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrPipelineState, err)
	}

	return findPendingApprovals(*pipeline.Name, pipelineOutput.Pipeline.Stages, stateOutput.StageStates), nil
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
	if actionState.ActionName == nil || actionState.LatestExecution == nil || actionState.LatestExecution.Token == nil {
		return false
	}

	category, ok := actionTypes[*actionState.ActionName]
	if !ok {
		return false
	}

	return category == types.ActionCategoryApproval &&
		actionState.LatestExecution.Status == types.ActionExecutionStatusInProgress
}

// PutApprovalResult handles the approval or rejection of a manual approval action.
func PutApprovalResult(ctx context.Context, profile, region string, action ApprovalAction, approved bool, comment string) error {
	client, err := getClient(ctx, profile, region)
	if err != nil {
		return err
	}

	status := types.ApprovalStatusRejected
	if approved {
		status = types.ApprovalStatusApproved
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
		return fmt.Errorf("%w: %w", ErrApprovalResult, err)
	}

	return nil
}

// GetPipelineStatus returns the status of all pipelines
func GetPipelineStatus(ctx context.Context, profile, region string) ([]PipelineStatus, error) {
	client, err := getClient(ctx, profile, region)
	if err != nil {
		return nil, err
	}

	pipelines, err := listPipelines(ctx, client)
	if err != nil {
		return nil, err
	}

	var pipelineStatuses []PipelineStatus

	for _, pipeline := range pipelines {
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
func StartPipelineExecution(ctx context.Context, profile, region, pipelineName, commitID string) error {
	client, err := getClient(ctx, profile, region)
	if err != nil {
		return err
	}

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
