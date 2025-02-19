package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
)

// Service represents an AWS service client
type Service struct {
	codePipelineClient *codepipeline.Client
}

// NewService creates a new AWS service client
func NewService(profile, region string) (*Service, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return nil, err
	}

	return &Service{
		codePipelineClient: codepipeline.NewFromConfig(cfg),
	}, nil
}

// ListPipelines returns a list of CodePipeline pipelines
func (s *Service) ListPipelines(ctx context.Context) ([]string, error) {
	input := &codepipeline.ListPipelinesInput{}
	output, err := s.codePipelineClient.ListPipelines(ctx, input)
	if err != nil {
		return nil, err
	}

	var pipelines []string
	for _, pipeline := range output.Pipelines {
		pipelines = append(pipelines, *pipeline.Name)
	}

	return pipelines, nil
}

// GetPipelineState returns the state of a pipeline
func (s *Service) GetPipelineState(ctx context.Context, pipelineName string) (*codepipeline.GetPipelineStateOutput, error) {
	input := &codepipeline.GetPipelineStateInput{
		Name: &pipelineName,
	}
	return s.codePipelineClient.GetPipelineState(ctx, input)
}

// PutApprovalResult submits an approval result for a manual approval action
func (s *Service) PutApprovalResult(ctx context.Context, pipelineName, stageName, actionName, token, summary string, status bool) error {
	result := types.ApprovalStatusRejected
	if status {
		result = types.ApprovalStatusApproved
	}

	input := &codepipeline.PutApprovalResultInput{
		PipelineName: &pipelineName,
		StageName:    &stageName,
		ActionName:   &actionName,
		Result: &types.ApprovalResult{
			Summary: &summary,
			Status:  result,
		},
		Token: &token,
	}

	_, err := s.codePipelineClient.PutApprovalResult(ctx, input)
	return err
}

// GetPendingApprovals returns a list of pending manual approval actions
func (s *Service) GetPendingApprovals(ctx context.Context) ([]ApprovalAction, error) {
	pipelines, err := s.ListPipelines(ctx)
	if err != nil {
		return nil, err
	}

	var approvals []ApprovalAction
	for _, pipeline := range pipelines {
		state, err := s.GetPipelineState(ctx, pipeline)
		if err != nil {
			continue
		}

		for _, stage := range state.StageStates {
			for _, action := range stage.ActionStates {
				if action.ActionName != nil && action.LatestExecution != nil &&
					action.LatestExecution.Status == types.ActionExecutionStatusInProgress &&
					action.LatestExecution.Token != nil {
					approvals = append(approvals, ApprovalAction{
						PipelineName: pipeline,
						StageName:    *stage.StageName,
						ActionName:   *action.ActionName,
						Token:        *action.LatestExecution.Token,
					})
				}
			}
		}
	}

	return approvals, nil
}

func (s *Service) ListPendingApprovals(ctx context.Context) ([]ApprovalAction, error) {
	// Implementation of ListPendingApprovals method
	return nil, nil
}

// HandleApproval handles the approval or rejection of a manual approval action
func (s *Service) HandleApproval(ctx context.Context, action *ApprovalAction, approve bool, summary string) error {
	return s.PutApprovalResult(ctx, action.PipelineName, action.StageName, action.ActionName, action.Token, summary, approve)
}
