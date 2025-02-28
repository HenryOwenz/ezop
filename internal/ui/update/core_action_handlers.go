package update

import (
	"context"
	"fmt"
	"strings"

	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
)

// UpdateModelForView updates the model based on the current view
func UpdateModelForView(m *model.Model) error {
	switch m.CurrentView {
	case constants.ViewAWSConfig:
		if m.AwsProfile == "" {
			// Get profiles from the registry
			provider, err := m.Registry.Get("AWS")
			if err != nil {
				return err
			}
			profiles, err := provider.GetProfiles()
			if err != nil {
				return err
			}
			m.Profiles = profiles
		} else {
			m.Regions = constants.DefaultAWSRegions
		}
	case constants.ViewApprovals:
		if len(m.Approvals) == 0 {
			// Get the provider from the registry
			provider, err := m.Registry.Get("AWS")
			if err != nil {
				return err
			}

			// Get the CodePipelineManualApprovalOperation from the provider
			approvalOperation, err := provider.GetCodePipelineManualApprovalOperation()
			if err != nil {
				return err
			}

			// Get approvals using the operation
			ctx := context.Background()
			approvals, err := approvalOperation.GetPendingApprovals(ctx)
			if err != nil {
				return err
			}

			m.Provider = provider
			m.Approvals = approvals
		}
	case constants.ViewPipelineStatus:
		if len(m.Pipelines) == 0 {
			// Get the provider from the registry
			provider, err := m.Registry.Get("AWS")
			if err != nil {
				return err
			}

			// Get the PipelineStatusOperation from the provider
			statusOperation, err := provider.GetPipelineStatusOperation()
			if err != nil {
				return err
			}

			// Get pipeline status using the operation
			ctx := context.Background()
			pipelines, err := statusOperation.GetPipelineStatus(ctx)
			if err != nil {
				return err
			}

			m.Provider = provider
			m.Pipelines = pipelines
		}
	case constants.ViewPipelineStages:
		if m.SelectedPipeline != nil && len(m.SelectedPipeline.Stages) == 0 {
			if m.Provider == nil {
				// Get the provider from the registry
				provider, err := m.Registry.Get("AWS")
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
func ExecuteAction(m *model.Model) error {
	if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
		if m.SelectedPipeline == nil {
			return fmt.Errorf(constants.MsgErrorNoPipeline)
		}

		if m.Provider == nil {
			// Get the provider from the registry
			provider, err := m.Registry.Get("AWS")
			if err != nil {
				return err
			}
			m.Provider = provider
		}

		// Get the StartPipelineOperation from the provider
		startOperation, err := m.Provider.GetStartPipelineOperation()
		if err != nil {
			return err
		}

		if m.ManualCommitID && strings.TrimSpace(m.CommitID) == "" {
			return fmt.Errorf(constants.MsgErrorEmptyCommitID)
		}

		// Start the pipeline execution using the operation
		ctx := context.Background()
		err = startOperation.StartPipelineExecution(ctx, m.SelectedPipeline.Name, m.CommitID)

		HandlePipelineExecution(m, err)
		return nil
	}

	if m.SelectedApproval == nil {
		return fmt.Errorf(constants.MsgErrorNoApproval)
	}

	if m.Provider == nil {
		// Get the provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return err
		}
		m.Provider = provider
	}

	// Get the CodePipelineManualApprovalOperation from the provider
	approvalOperation, err := m.Provider.GetCodePipelineManualApprovalOperation()
	if err != nil {
		return err
	}

	if strings.TrimSpace(m.ApprovalComment) == "" {
		return fmt.Errorf(constants.MsgErrorEmptyComment)
	}

	// Execute the approval action using the operation
	ctx := context.Background()
	err = approvalOperation.ApproveAction(ctx, *m.SelectedApproval, m.ApproveAction, m.ApprovalComment)

	HandleApprovalResult(m, err)
	return nil
}
