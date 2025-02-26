package ui

import (
	"context"
	"fmt"

	"github.com/HenryOwenz/cloudgate/internal/aws"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	tea "github.com/charmbracelet/bubbletea"
)

// handleEnter processes the enter key press based on the current view
func (m *Model) handleEnter() (tea.Model, tea.Cmd) {
	switch m.currentView {
	case constants.ViewProviders:
		return m.handleProviderSelection()
	case constants.ViewAWSConfig:
		return m.handleAWSConfigSelection()
	case constants.ViewSelectService:
		return m.handleServiceSelection()
	case constants.ViewSelectCategory:
		return m.handleCategorySelection()
	case constants.ViewSelectOperation:
		return m.handleOperationSelection()
	case constants.ViewApprovals:
		return m.handleApprovalSelection()
	case constants.ViewConfirmation:
		return m.handleConfirmationSelection()
	case constants.ViewSummary:
		if !m.manualInput {
			if m.selectedOperation != nil && m.selectedOperation.Name == "Start Pipeline" {
				if selected := m.table.SelectedRow(); len(selected) > 0 {
					newModel := *m
					switch selected[0] {
					case "Latest Commit":
						newModel.currentView = constants.ViewExecutingAction
						newModel.summary = "" // Empty string means use latest commit
						newModel.updateTableForView()
						return &newModel, nil
					case "Manual Input":
						newModel.manualInput = true
						newModel.textInput.Focus()
						newModel.textInput.Placeholder = "Enter commit ID..."
						return &newModel, nil
					}
				}
			}
		}
		return m.handleSummaryConfirmation()
	case constants.ViewExecutingAction:
		return m.handleExecutionSelection()
	case constants.ViewPipelineStatus:
		if selected := m.table.SelectedRow(); len(selected) > 0 {
			newModel := *m
			for _, pipeline := range m.pipelines {
				if pipeline.Name == selected[0] {
					if m.selectedOperation != nil && m.selectedOperation.Name == "Start Pipeline" {
						newModel.currentView = constants.ViewExecutingAction
						newModel.selectedPipeline = &pipeline
						newModel.updateTableForView()
						return &newModel, nil
					}
					newModel.currentView = constants.ViewPipelineStages
					newModel.selectedPipeline = &pipeline
					newModel.updateTableForView()
					return &newModel, nil
				}
			}
		}
	case constants.ViewPipelineStages:
		// Just view only, no action
	}
	return m, nil
}

func (m *Model) handleProviderSelection() (tea.Model, tea.Cmd) {
	if selected := m.table.SelectedRow(); len(selected) > 0 {
		if selected[0] == "Amazon Web Services" {
			newModel := *m
			newModel.currentView = constants.ViewAWSConfig
			newModel.updateTableForView()
			return &newModel, nil
		}
	}
	return m, nil
}

func (m *Model) handleAWSConfigSelection() (tea.Model, tea.Cmd) {
	if selected := m.table.SelectedRow(); len(selected) > 0 {
		newModel := *m
		if m.awsProfile == "" {
			newModel.awsProfile = selected[0]
			newModel.updateTableForView()
		} else {
			newModel.awsRegion = selected[0]
			newModel.currentView = constants.ViewSelectService
			newModel.updateTableForView()
		}
		return &newModel, nil
	}
	return m, nil
}

func (m *Model) handleServiceSelection() (tea.Model, tea.Cmd) {
	if selected := m.table.SelectedRow(); len(selected) > 0 {
		newModel := *m
		newModel.selectedService = &Service{
			Name:        selected[0],
			Description: selected[1],
		}
		newModel.currentView = constants.ViewSelectCategory
		newModel.updateTableForView()
		return &newModel, nil
	}
	return m, nil
}

func (m *Model) handleCategorySelection() (tea.Model, tea.Cmd) {
	if selected := m.table.SelectedRow(); len(selected) > 0 {
		newModel := *m
		newModel.selectedCategory = &Category{
			Name:        selected[0],
			Description: selected[1],
		}
		newModel.currentView = constants.ViewSelectOperation
		newModel.updateTableForView()
		return &newModel, nil
	}
	return m, nil
}

func (m *Model) handleOperationSelection() (tea.Model, tea.Cmd) {
	if selected := m.table.SelectedRow(); len(selected) > 0 {
		newModel := *m
		newModel.selectedOperation = &Operation{
			Name:        selected[0],
			Description: selected[1],
		}

		if selected[0] == "Pipeline Approvals" {
			// Start loading approvals
			newModel.isLoading = true
			newModel.loadingMsg = "Loading approvals..."
			return &newModel, m.fetchApprovals
		} else if selected[0] == "Pipeline Status" || selected[0] == "Start Pipeline" {
			// Start loading pipeline status
			newModel.isLoading = true
			newModel.loadingMsg = "Loading pipelines..."
			return &newModel, m.fetchPipelineStatus
		}
	}
	return m, nil
}

func (m *Model) handleApprovalSelection() (tea.Model, tea.Cmd) {
	if selected := m.table.SelectedRow(); len(selected) > 0 {
		newModel := *m
		for _, approval := range m.approvals {
			if approval.PipelineName == selected[0] &&
				approval.StageName == selected[1] &&
				approval.ActionName == selected[2] {
				newModel.selectedApproval = &approval
				newModel.currentView = constants.ViewConfirmation
				newModel.updateTableForView()
				return &newModel, nil
			}
		}
	}
	return m, nil
}

func (m *Model) handleConfirmationSelection() (tea.Model, tea.Cmd) {
	if selected := m.table.SelectedRow(); len(selected) > 0 {
		newModel := *m
		if selected[0] == "Approve" {
			newModel.approveAction = true
			newModel.currentView = constants.ViewSummary
			newModel.setTextInputForApproval(true)
		} else if selected[0] == "Reject" {
			newModel.approveAction = false
			newModel.currentView = constants.ViewSummary
			newModel.setTextInputForApproval(false)
		}
		newModel.updateTableForView()
		return &newModel, nil
	}
	return m, nil
}

func (m *Model) handleSummaryConfirmation() (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Model) handleExecutionSelection() (tea.Model, tea.Cmd) {
	if selected := m.table.SelectedRow(); len(selected) > 0 {
		if selected[0] == "Execute" {
			newModel := *m
			newModel.isLoading = true
			if m.selectedOperation != nil && m.selectedOperation.Name == "Start Pipeline" {
				newModel.loadingMsg = "Starting pipeline..."
				return &newModel, m.executePipeline
			} else {
				newModel.loadingMsg = "Executing approval action..."
				return &newModel, m.executeApproval
			}
		} else if selected[0] == "Cancel" {
			return m.navigateBack(), nil
		}
	}
	return m, nil
}

// Async operations
func (m *Model) fetchApprovals() tea.Msg {
	// Create a new AWS provider
	provider, err := aws.New(context.Background(), m.awsProfile, m.awsRegion)
	if err != nil {
		return errMsg{err}
	}

	// Get approvals
	approvals, err := provider.GetPendingApprovals(context.Background())
	if err != nil {
		return errMsg{err}
	}

	return approvalsMsg{
		provider:  provider,
		approvals: approvals,
	}
}

func (m *Model) fetchPipelineStatus() tea.Msg {
	// Create a new AWS provider
	provider, err := aws.New(context.Background(), m.awsProfile, m.awsRegion)
	if err != nil {
		return errMsg{err}
	}

	// Get pipeline status
	pipelines, err := provider.GetPipelineStatus(context.Background())
	if err != nil {
		return errMsg{err}
	}

	return pipelineStatusMsg{
		provider:  provider,
		pipelines: pipelines,
	}
}

func (m *Model) executeApproval() tea.Msg {
	if m.provider == nil || m.selectedApproval == nil {
		return errMsg{fmt.Errorf("provider or approval not set")}
	}

	err := m.provider.PutApprovalResult(context.Background(), *m.selectedApproval, m.approveAction, m.summary)

	return approvalResultMsg{err: err}
}

func (m *Model) executePipeline() tea.Msg {
	if m.provider == nil || m.selectedPipeline == nil {
		return errMsg{fmt.Errorf("provider or pipeline not set")}
	}

	err := m.provider.StartPipelineExecution(context.Background(), m.selectedPipeline.Name, m.summary)
	return pipelineExecutionMsg{err: err}
}
