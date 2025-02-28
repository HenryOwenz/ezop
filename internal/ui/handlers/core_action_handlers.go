package handlers

import (
	"context"
	"fmt"
	"strings"

	"github.com/HenryOwenz/cloudgate/internal/aws"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/core"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
)

// UpdateModelForView updates the model based on the current view
func UpdateModelForView(m *core.Model) error {
	switch m.CurrentView {
	case constants.ViewAWSConfig:
		if m.AwsProfile == "" {
			profiles := aws.GetProfiles()
			m.Profiles = profiles
		} else {
			m.Regions = constants.DefaultAWSRegions
		}
	case constants.ViewApprovals:
		if len(m.Approvals) == 0 {
			provider, err := aws.New(context.Background(), m.AwsProfile, m.AwsRegion)
			if err != nil {
				return err
			}

			approvals, err := provider.GetPendingApprovals(context.Background())
			if err != nil {
				return err
			}

			m.Provider = provider
			m.Approvals = approvals
		}
	case constants.ViewPipelineStatus:
		if len(m.Pipelines) == 0 {
			provider, err := aws.New(context.Background(), m.AwsProfile, m.AwsRegion)
			if err != nil {
				return err
			}

			pipelines, err := provider.GetPipelineStatus(context.Background())
			if err != nil {
				return err
			}

			m.Provider = provider
			m.Pipelines = pipelines
		}
	case constants.ViewPipelineStages:
		if m.SelectedPipeline != nil && len(m.SelectedPipeline.Stages) == 0 {
			if m.Provider == nil {
				provider, err := aws.New(context.Background(), m.AwsProfile, m.AwsRegion)
				if err != nil {
					return err
				}
				m.Provider = provider
			}

			for _, pipeline := range m.Pipelines {
				if pipeline.Name == m.SelectedPipeline.Name {
					m.SelectedPipeline = &pipeline
					break
				}
			}
		}
	}

	// Update the table for the current view
	view.UpdateTableForView(m)
	return nil
}

// ExecuteAction executes the selected action
func ExecuteAction(m *core.Model) error {
	if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
		if m.SelectedPipeline == nil {
			return fmt.Errorf(constants.MsgErrorNoPipeline)
		}

		if m.Provider == nil {
			provider, err := aws.New(context.Background(), m.AwsProfile, m.AwsRegion)
			if err != nil {
				return err
			}
			m.Provider = provider
		}

		var err error
		if m.ManualCommitID && strings.TrimSpace(m.CommitID) == "" {
			return fmt.Errorf(constants.MsgErrorEmptyCommitID)
		}

		if m.ManualCommitID {
			err = m.Provider.StartPipelineExecution(context.Background(), m.SelectedPipeline.Name, m.CommitID)
		} else {
			err = m.Provider.StartPipelineExecution(context.Background(), m.SelectedPipeline.Name, "")
		}

		HandlePipelineExecution(m, err)
		return nil
	}

	if m.SelectedApproval == nil {
		return fmt.Errorf(constants.MsgErrorNoApproval)
	}

	if m.Provider == nil {
		provider, err := aws.New(context.Background(), m.AwsProfile, m.AwsRegion)
		if err != nil {
			return err
		}
		m.Provider = provider
	}

	var err error
	if m.ApproveAction {
		if strings.TrimSpace(m.ApprovalComment) == "" {
			return fmt.Errorf(constants.MsgErrorEmptyComment)
		}
		err = m.Provider.PutApprovalResult(context.Background(), *m.SelectedApproval, true, m.ApprovalComment)
	} else {
		if strings.TrimSpace(m.ApprovalComment) == "" {
			return fmt.Errorf(constants.MsgErrorEmptyComment)
		}
		err = m.Provider.PutApprovalResult(context.Background(), *m.SelectedApproval, false, m.ApprovalComment)
	}

	HandleApprovalResult(m, err)
	return nil
}
