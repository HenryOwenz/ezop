package update

import (
	"context"
	"fmt"
	"strings"

	"github.com/HenryOwenz/cloudgate/internal/aws"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/core"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
)

// HandleApprovalResult handles the result of an approval action
func HandleApprovalResult(m *core.Model, err error) {
	if err != nil {
		m.Error = fmt.Sprintf(constants.MsgErrorGeneric, err.Error())
		m.CurrentView = constants.ViewError
		return
	}

	// Use the appropriate message constant based on approval action
	if m.ApproveAction {
		m.Success = fmt.Sprintf(constants.MsgApprovalSuccess,
			m.SelectedApproval.PipelineName,
			m.SelectedApproval.StageName,
			m.SelectedApproval.ActionName)
	} else {
		m.Success = fmt.Sprintf(constants.MsgRejectionSuccess,
			m.SelectedApproval.PipelineName,
			m.SelectedApproval.StageName,
			m.SelectedApproval.ActionName)
	}

	// Reset approval state
	m.SelectedApproval = nil
	m.ApprovalComment = ""

	// Completely reset the text input
	m.ResetTextInput()
	m.TextInput.Placeholder = constants.MsgEnterComment
	m.ManualInput = false

	// Navigate back to the operation selection view
	m.CurrentView = constants.ViewSelectOperation

	// Clear the approvals list to force a refresh next time
	m.Approvals = nil

	// Update the table for the current view
	view.UpdateTableForView(m)
}

// HandlePipelineExecution handles the result of a pipeline execution
func HandlePipelineExecution(m *core.Model, err error) {
	if err != nil {
		m.Error = fmt.Sprintf(constants.MsgErrorGeneric, err.Error())
		m.CurrentView = constants.ViewError
		return
	}

	m.Success = fmt.Sprintf(constants.MsgPipelineStartSuccess, m.SelectedPipeline.Name)

	// Reset pipeline state
	m.SelectedPipeline = nil
	m.CommitID = ""
	m.ManualCommitID = false

	// Completely reset the text input
	m.ResetTextInput()
	m.TextInput.Placeholder = constants.MsgEnterComment
	m.ManualInput = false

	// Navigate back to the operation selection view
	m.CurrentView = constants.ViewSelectOperation

	// Clear the pipelines list to force a refresh next time
	m.Pipelines = nil

	// Update the table for the current view
	view.UpdateTableForView(m)
}

// UpdateModelForView updates the model based on the current view
func UpdateModelForView(m *core.Model) error {
	switch m.CurrentView {
	case constants.ViewAWSConfig:
		if m.AwsProfile == "" {
			profiles := aws.GetProfiles()
			m.Profiles = profiles
		} else {
			m.Regions = []string{
				"us-east-1", "us-east-2", "us-west-1", "us-west-2",
				"eu-west-1", "eu-west-2", "eu-central-1",
				"ap-southeast-1", "ap-southeast-2", "ap-northeast-1",
			}
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
