package navigation

import (
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/core"
)

// NavigateBack handles navigation to the previous view
func NavigateBack(m *core.Model) *core.Model {
	newModel := *m
	switch m.CurrentView {
	case constants.ViewAWSConfig:
		if m.AwsProfile != "" {
			// If we're in region selection, just clear region and stay in AWS config
			newModel.AwsRegion = ""
			newModel.AwsProfile = ""
			// Don't change the view - we'll stay in AWS config to show profiles
		} else {
			// If we're in profile selection, go back to providers
			newModel.CurrentView = constants.ViewProviders
		}
		newModel.ManualInput = false
		newModel.ResetTextInput()
	case constants.ViewSelectService:
		newModel.CurrentView = constants.ViewAWSConfig
		newModel.SelectedService = nil
	case constants.ViewSelectCategory:
		newModel.CurrentView = constants.ViewSelectService
		newModel.SelectedCategory = nil
	case constants.ViewSelectOperation:
		newModel.CurrentView = constants.ViewSelectCategory
		newModel.SelectedOperation = nil
	case constants.ViewApprovals:
		newModel.CurrentView = constants.ViewSelectOperation
		newModel.ResetApprovalState()
	case constants.ViewConfirmation:
		newModel.CurrentView = constants.ViewApprovals
		newModel.SelectedApproval = nil
	case constants.ViewSummary:
		newModel.CurrentView = constants.ViewConfirmation
		newModel.Summary = ""
		newModel.ResetTextInput()
	case constants.ViewExecutingAction:
		if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
			// For pipeline start flow, go back to pipeline selection
			newModel.CurrentView = constants.ViewPipelineStatus
			newModel.SelectedPipeline = nil
		} else {
			// For approval flow, go back to summary
			newModel.CurrentView = constants.ViewSummary
			// When going back to summary, restore the previous comment and focus
			newModel.TextInput.SetValue(m.Summary)
			newModel.TextInput.Focus()
			if newModel.ApproveAction {
				newModel.TextInput.Placeholder = "Enter approval comment..."
			} else {
				newModel.TextInput.Placeholder = "Enter rejection comment..."
			}
		}
	case constants.ViewPipelineStages:
		newModel.CurrentView = constants.ViewPipelineStatus
		newModel.SelectedPipeline = nil
	case constants.ViewPipelineStatus:
		newModel.CurrentView = constants.ViewSelectOperation
		newModel.Pipelines = nil
		newModel.Provider = nil
	}
	return &newModel
}
