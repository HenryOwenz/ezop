package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
)

// CodePipelineService handles AWS CodePipeline operations
type CodePipelineService struct {
	client *codepipeline.Client
}

// NewCodePipelineService creates a new AWS CodePipeline service
func NewCodePipelineService(profile, region string) (*CodePipelineService, error) {
	if profile == "" {
		return nil, fmt.Errorf("AWS profile must be specified")
	}
	if region == "" {
		return nil, fmt.Errorf("AWS region must be specified")
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config with profile %s in region %s: %w", profile, region, err)
	}

	client := codepipeline.NewFromConfig(cfg)
	return &CodePipelineService{client: client}, nil
}

// ListPendingApprovals returns all pending manual approval actions
func (s *CodePipelineService) ListPendingApprovals(ctx context.Context) ([]ApprovalAction, error) {
	var approvals []ApprovalAction

	// List all pipelines first
	pipelineOutput, err := s.client.ListPipelines(ctx, &codepipeline.ListPipelinesInput{})
	if err != nil {
		return nil, err
	}

	// Check each pipeline for pending approvals
	for _, pipeline := range pipelineOutput.Pipelines {
		pipelineName := *pipeline.Name
		stateOutput, err := s.client.GetPipelineState(ctx, &codepipeline.GetPipelineStateInput{
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

				if actionState.ActionName != nil && actionState.LatestExecution.Status == types.ActionExecutionStatusInProgress {
					// Get the approval token
					token := ""
					if actionState.LatestExecution.Token != nil {
						token = *actionState.LatestExecution.Token
					}

					approvals = append(approvals, ApprovalAction{
						PipelineName: pipelineName,
						StageName:    *stageState.StageName,
						ActionName:   *actionState.ActionName,
						Token:        token,
					})
				}
			}
		}
	}

	return approvals, nil
}

// HandleApproval handles a manual approval action
func (s *CodePipelineService) HandleApproval(ctx context.Context, params map[string]interface{}) error {
	// Extract parameters
	pipelineName, ok := params["pipeline_name"].(string)
	if !ok {
		return fmt.Errorf("pipeline_name parameter is required")
	}
	stageName, ok := params["stage_name"].(string)
	if !ok {
		return fmt.Errorf("stage_name parameter is required")
	}
	actionName, ok := params["action_name"].(string)
	if !ok {
		return fmt.Errorf("action_name parameter is required")
	}
	token, ok := params["token"].(string)
	if !ok {
		return fmt.Errorf("token parameter is required")
	}
	summary, ok := params["summary"].(string)
	if !ok {
		return fmt.Errorf("summary parameter is required")
	}
	approve, ok := params["approve"].(bool)
	if !ok {
		return fmt.Errorf("approve parameter is required")
	}

	status := types.ApprovalStatusRejected
	if approve {
		status = types.ApprovalStatusApproved
	}

	input := &codepipeline.PutApprovalResultInput{
		PipelineName: &pipelineName,
		StageName:    &stageName,
		ActionName:   &actionName,
		Token:        &token,
		Result: &types.ApprovalResult{
			Summary: &summary,
			Status:  status,
		},
	}

	_, err := s.client.PutApprovalResult(ctx, input)
	return err
}
