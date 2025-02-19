package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
)

// ApprovalAction represents a manual approval action in CodePipeline
type ApprovalAction struct {
	PipelineName string
	StageName    string
	ActionName   string
	Token        string
}

// Service handles AWS CodePipeline operations
type Service struct {
	client *codepipeline.Client
}

// NewService creates a new AWS CodePipeline service
func NewService(ctx context.Context, profile, region string) (*Service, error) {
	if profile == "" {
		return nil, fmt.Errorf("AWS profile must be specified")
	}
	if region == "" {
		return nil, fmt.Errorf("AWS region must be specified")
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(profile),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config with profile %s in region %s: %w", profile, region, err)
	}

	client := codepipeline.NewFromConfig(cfg)
	return &Service{client: client}, nil
}

// ListPendingApprovals returns all pending manual approval actions
func (s *Service) ListPendingApprovals(ctx context.Context) ([]ApprovalAction, error) {
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

// ApproveAction approves a manual approval action
func (s *Service) ApproveAction(ctx context.Context, pipelineName, stageName, actionName, token, summary string) error {
	input := &codepipeline.PutApprovalResultInput{
		PipelineName: &pipelineName,
		StageName:    &stageName,
		ActionName:   &actionName,
		Token:        &token,
		Result: &types.ApprovalResult{
			Summary: &summary,
			Status:  types.ApprovalStatusApproved,
		},
	}

	_, err := s.client.PutApprovalResult(ctx, input)
	return err
}

// RejectAction rejects a manual approval action
func (s *Service) RejectAction(ctx context.Context, pipelineName, stageName, actionName, token, summary string) error {
	input := &codepipeline.PutApprovalResultInput{
		PipelineName: &pipelineName,
		StageName:    &stageName,
		ActionName:   &actionName,
		Token:        &token,
		Result: &types.ApprovalResult{
			Summary: &summary,
			Status:  types.ApprovalStatusRejected,
		},
	}

	_, err := s.client.PutApprovalResult(ctx, input)
	return err
}
