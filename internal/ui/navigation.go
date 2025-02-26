package ui

import (
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
)

func (m *Model) navigateBack() *Model {
	newModel := *m
	switch m.currentView {
	case constants.ViewAWSConfig:
		if m.awsProfile != "" {
			// If we're in region selection, just clear region and stay in AWS config
			newModel.awsRegion = ""
			newModel.awsProfile = ""
			// Don't change the view - we'll stay in AWS config to show profiles
		} else {
			// If we're in profile selection, go back to providers
			newModel.currentView = constants.ViewProviders
		}
		newModel.manualInput = false
		newModel.resetTextInput()
	case constants.ViewSelectService:
		newModel.currentView = constants.ViewAWSConfig
		newModel.selectedService = nil
	case constants.ViewSelectCategory:
		newModel.currentView = constants.ViewSelectService
		newModel.selectedCategory = nil
	case constants.ViewSelectOperation:
		newModel.currentView = constants.ViewSelectCategory
		newModel.selectedOperation = nil
	case constants.ViewApprovals:
		newModel.currentView = constants.ViewSelectOperation
		newModel.resetApprovalState()
	case constants.ViewConfirmation:
		newModel.currentView = constants.ViewApprovals
		newModel.selectedApproval = nil
	case constants.ViewSummary:
		newModel.currentView = constants.ViewConfirmation
		newModel.summary = ""
		newModel.resetTextInput()
	case constants.ViewExecutingAction:
		if m.selectedOperation != nil && m.selectedOperation.Name == "Start Pipeline" {
			// For pipeline start flow, go back to pipeline selection
			newModel.currentView = constants.ViewPipelineStatus
			newModel.selectedPipeline = nil
		} else {
			// For approval flow, go back to summary
			newModel.currentView = constants.ViewSummary
			// When going back to summary, restore the previous comment and focus
			newModel.textInput.SetValue(m.summary)
			newModel.textInput.Focus()
			if newModel.approveAction {
				newModel.textInput.Placeholder = "Enter approval comment..."
			} else {
				newModel.textInput.Placeholder = "Enter rejection comment..."
			}
		}
	case constants.ViewPipelineStages:
		newModel.currentView = constants.ViewPipelineStatus
		newModel.selectedPipeline = nil
	case constants.ViewPipelineStatus:
		newModel.currentView = constants.ViewSelectOperation
		newModel.pipelines = nil
		newModel.provider = nil
	}
	newModel.updateTableForView()
	return &newModel
}

// Helper functions for common operations
func (m *Model) resetApprovalState() {
	m.approvals = nil
	m.provider = nil
	m.selectedApproval = nil
	m.summary = ""
}

func (m *Model) resetTextInput() {
	m.textInput.SetValue("")
	m.textInput.Blur()
}

func (m *Model) setTextInputForApproval(isApproval bool) {
	m.textInput.Focus()
	if isApproval {
		m.textInput.Placeholder = "Enter approval comment..."
	} else {
		m.textInput.Placeholder = "Enter rejection comment..."
	}
} 