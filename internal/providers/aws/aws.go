package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/HenryOwenz/cloudgate/internal/providers"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	cpTypes "github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
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
func findPendingApprovals(pipelineName string, stages []cpTypes.StageDeclaration, stageStates []cpTypes.StageState) []ApprovalAction {
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
func buildActionTypeMap(stages []cpTypes.StageDeclaration) map[string]cpTypes.ActionCategory {
	actionTypes := make(map[string]cpTypes.ActionCategory)
	for _, stage := range stages {
		for _, action := range stage.Actions {
			actionTypes[*action.Name] = action.ActionTypeId.Category
		}
	}

	return actionTypes
}

// buildStageStateMap creates a map of stage names to their states for quick lookup.
func buildStageStateMap(stageStates []cpTypes.StageState) map[string]cpTypes.StageState {
	stateMap := make(map[string]cpTypes.StageState)
	for _, state := range stageStates {
		stateMap[*state.StageName] = state
	}

	return stateMap
}

// findStageApprovals returns a list of pending approval actions from a single stage.
func findStageApprovals(pipelineName string, stage cpTypes.StageDeclaration, state cpTypes.StageState, actionTypes map[string]cpTypes.ActionCategory) []ApprovalAction {
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
func isApprovalAction(actionState cpTypes.ActionState, actionTypes map[string]cpTypes.ActionCategory) bool {
	return actionState.LatestExecution != nil &&
		actionState.LatestExecution.Token != nil &&
		actionState.LatestExecution.Status == cpTypes.ActionExecutionStatusInProgress &&
		actionTypes[*actionState.ActionName] == cpTypes.ActionCategoryApproval
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

	status := cpTypes.ApprovalStatusApproved
	if !approved {
		status = cpTypes.ApprovalStatusRejected
	}

	_, err = client.PutApprovalResult(ctx, &codepipeline.PutApprovalResultInput{
		ActionName:   aws.String(action.ActionName),
		PipelineName: aws.String(action.PipelineName),
		Result: &cpTypes.ApprovalResult{
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
func getPipelineStatus(ctx context.Context, client *codepipeline.Client, pipeline cpTypes.PipelineSummary) (PipelineStatus, error) {
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
		input.SourceRevisions = []cpTypes.SourceRevisionOverride{
			{
				ActionName:    aws.String("Source"), // Assuming standard Source action name
				RevisionType:  cpTypes.SourceRevisionTypeCommitId,
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

// GetFunctionStatus returns the status of all Lambda functions
func (p *Provider) GetFunctionStatus(ctx context.Context) ([]providers.FunctionStatus, error) {
	// Create a new AWS SDK client
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(p.profile),
		config.WithRegion(p.region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := lambda.NewFromConfig(cfg)

	// List all functions
	var functions []providers.FunctionStatus
	var marker *string

	for {
		output, err := client.ListFunctions(ctx, &lambda.ListFunctionsInput{
			Marker: marker,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list functions: %w", err)
		}

		// Convert Lambda functions to FunctionStatus
		for _, function := range output.Functions {
			memory := int32(0)
			if function.MemorySize != nil {
				memory = *function.MemorySize
			}

			timeout := int32(0)
			if function.Timeout != nil {
				timeout = *function.Timeout
			}

			// Get architecture (default to x86_64 if not specified)
			architecture := "x86_64"
			if len(function.Architectures) > 0 {
				architecture = string(function.Architectures[0])
			}

			// Get log group if available
			logGroup := ""
			if function.LoggingConfig != nil && function.LoggingConfig.LogGroup != nil {
				logGroup = *function.LoggingConfig.LogGroup
			}

			functions = append(functions, providers.FunctionStatus{
				Name:         aws.ToString(function.FunctionName),
				Runtime:      string(function.Runtime),
				Memory:       memory,
				Timeout:      timeout,
				LastUpdate:   aws.ToString(function.LastModified),
				Role:         aws.ToString(function.Role),
				Handler:      aws.ToString(function.Handler),
				Description:  aws.ToString(function.Description),
				FunctionArn:  aws.ToString(function.FunctionArn),
				CodeSize:     function.CodeSize,
				Version:      aws.ToString(function.Version),
				PackageType:  string(function.PackageType),
				Architecture: architecture,
				LogGroup:     logGroup,
			})
		}

		if output.NextMarker == nil {
			break
		}
		marker = output.NextMarker
	}

	return functions, nil
}
