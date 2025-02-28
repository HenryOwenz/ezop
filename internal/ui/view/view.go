package view

import (
	"fmt"

	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/charmbracelet/lipgloss"
)

// Render renders the UI
func Render(m *model.Model) string {
	if m.Err != nil {
		return renderErrorView(m)
	}

	// Create content array with appropriate spacing
	content := make([]string, constants.AppContentLines)

	// Set the title and context
	content[0] = renderTitle(m)
	content[1] = renderContext(m)
	content[2] = renderLoadingSpinner(m)
	content[3] = renderMainContent(m)
	content[4] = renderHelpText(m)

	// Join all content vertically with consistent spacing
	return m.Styles.App.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			content...,
		),
	)
}

// renderErrorView renders the error view
func renderErrorView(m *model.Model) string {
	return m.Styles.App.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.Styles.Error.Render("Error: "+m.Err.Error()),
			"\n",
			m.Styles.Help.Render(fmt.Sprintf("%s: quit • %s: back", constants.KeyQ, constants.KeyAltBack)),
		),
	)
}

// renderTitle renders the title based on the current view
func renderTitle(m *model.Model) string {
	return m.Styles.Title.Render(getTitleText(m))
}

// renderContext renders the context based on the current view
func renderContext(m *model.Model) string {
	return m.Styles.Context.Render(getContextText(m))
}

// renderLoadingSpinner renders the loading spinner if needed
func renderLoadingSpinner(m *model.Model) string {
	if m.IsLoading {
		return m.Spinner.View()
	}
	return ""
}

// renderMainContent renders the main content area (table or text input)
func renderMainContent(m *model.Model) string {
	// For Summary view with approvals, always show text input
	if m.CurrentView == constants.ViewSummary && m.SelectedApproval != nil {
		return m.TextInput.View()
	}

	// For other views, follow normal logic
	if m.ManualInput {
		return m.TextInput.View()
	}

	return renderTable(m)
}

// renderHelpText renders the help text based on the current view
func renderHelpText(m *model.Model) string {
	return m.Styles.Help.Render(getHelpText(m))
}

// getContextText returns the appropriate context text for the current view
func getContextText(m *model.Model) string {
	switch m.CurrentView {
	case constants.ViewProviders:
		return getProvidersContextText()
	case constants.ViewAWSConfig:
		return getAWSConfigContextText(m)
	case constants.ViewSelectService:
		return getSelectServiceContextText(m)
	case constants.ViewSelectCategory:
		return getSelectCategoryContextText(m)
	case constants.ViewSelectOperation:
		return getSelectOperationContextText(m)
	case constants.ViewApprovals:
		return getApprovalsContextText(m)
	case constants.ViewConfirmation, constants.ViewSummary:
		return getConfirmationSummaryContextText(m)
	case constants.ViewExecutingAction:
		return getExecutingActionContextText(m)
	case constants.ViewPipelineStatus:
		return getPipelineStatusContextText(m)
	case constants.ViewPipelineStages:
		return getPipelineStagesContextText(m)
	default:
		return ""
	}
}

// getProvidersContextText returns the context text for the providers view
func getProvidersContextText() string {
	return constants.MsgAppDescription
}

// getAWSConfigContextText returns the context text for the AWS config view
func getAWSConfigContextText(m *model.Model) string {
	if m.AwsProfile == "" {
		// If in manual entry mode for profile, show the text input in the context
		if m.ManualInput {
			return fmt.Sprintf("Amazon Web Services\n\nEnter AWS Profile: %s", m.TextInput.View())
		}
		return "Amazon Web Services"
	}
	// If in manual entry mode for region, show the text input in the context
	if m.ManualInput {
		return fmt.Sprintf("Profile: %s\n\nEnter AWS Region: %s", m.AwsProfile, m.TextInput.View())
	}
	return fmt.Sprintf("Profile: %s", m.AwsProfile)
}

// getSelectServiceContextText returns the context text for the select service view
func getSelectServiceContextText(m *model.Model) string {
	return fmt.Sprintf("Profile: %s\nRegion: %s",
		m.AwsProfile,
		m.AwsRegion)
}

// getSelectCategoryContextText returns the context text for the select category view
func getSelectCategoryContextText(m *model.Model) string {
	if m.SelectedService == nil {
		return ""
	}
	return fmt.Sprintf("Service: %s",
		m.SelectedService.Name)
}

// getSelectOperationContextText returns the context text for the select operation view
func getSelectOperationContextText(m *model.Model) string {
	if m.SelectedService == nil || m.SelectedCategory == nil {
		return ""
	}
	return fmt.Sprintf("Service: %s\nCategory: %s",
		m.SelectedService.Name,
		m.SelectedCategory.Name)
}

// getApprovalsContextText returns the context text for the approvals view
func getApprovalsContextText(m *model.Model) string {
	return fmt.Sprintf("Profile: %s\nRegion: %s",
		m.AwsProfile,
		m.AwsRegion)
}

// getConfirmationSummaryContextText returns the context text for the confirmation and summary views
func getConfirmationSummaryContextText(m *model.Model) string {
	if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
		if m.SelectedPipeline == nil {
			return ""
		}
		return fmt.Sprintf("Profile: %s\nRegion: %s\nPipeline: %s",
			m.AwsProfile,
			m.AwsRegion,
			m.SelectedPipeline.Name)
	}
	if m.SelectedApproval == nil {
		return ""
	}
	return fmt.Sprintf("Pipeline: %s\nStage: %s\nAction: %s",
		m.SelectedApproval.PipelineName,
		m.SelectedApproval.StageName,
		m.SelectedApproval.ActionName)
}

// getExecutingActionContextText returns the context text for the executing action view
func getExecutingActionContextText(m *model.Model) string {
	if m.SelectedOperation != nil && m.SelectedOperation.Name == "Start Pipeline" {
		if m.SelectedPipeline == nil {
			return ""
		}
		return fmt.Sprintf("Profile: %s\nRegion: %s\nPipeline: %s\nRevisionID: Latest commit",
			m.AwsProfile,
			m.AwsRegion,
			m.SelectedPipeline.Name)
	}
	if m.SelectedApproval == nil {
		return ""
	}
	return fmt.Sprintf("Pipeline: %s\nStage: %s\nAction: %s\nComment: %s",
		m.SelectedApproval.PipelineName,
		m.SelectedApproval.StageName,
		m.SelectedApproval.ActionName,
		m.ApprovalComment)
}

// getPipelineStatusContextText returns the context text for the pipeline status view
func getPipelineStatusContextText(m *model.Model) string {
	return fmt.Sprintf("Profile: %s\nRegion: %s",
		m.AwsProfile,
		m.AwsRegion)
}

// getPipelineStagesContextText returns the context text for the pipeline stages view
func getPipelineStagesContextText(m *model.Model) string {
	if m.SelectedPipeline == nil {
		return ""
	}
	return fmt.Sprintf("Profile: %s\nRegion: %s\nPipeline: %s",
		m.AwsProfile,
		m.AwsRegion,
		m.SelectedPipeline.Name)
}

// getTitleText returns the appropriate title for the current view
func getTitleText(m *model.Model) string {
	// Map of view types to their corresponding titles
	titleMap := map[constants.View]string{
		constants.ViewProviders:       constants.TitleProviders,
		constants.ViewSelectService:   constants.TitleSelectService,
		constants.ViewSelectCategory:  constants.TitleSelectCategory,
		constants.ViewSelectOperation: constants.TitleSelectOperation,
		constants.ViewApprovals:       constants.TitleApprovals,
		constants.ViewConfirmation:    constants.TitleConfirmation,
		constants.ViewSummary:         constants.TitleSummary,
		constants.ViewExecutingAction: constants.TitleExecutingAction,
		constants.ViewPipelineStatus:  constants.TitlePipelineStatus,
		constants.ViewPipelineStages:  constants.TitlePipelineStages,
		constants.ViewError:           constants.TitleError,
		constants.ViewSuccess:         constants.TitleSuccess,
		constants.ViewHelp:            constants.TitleHelp,
	}

	// Special case for AWS config view
	if m.CurrentView == constants.ViewAWSConfig {
		if m.AwsProfile == "" {
			return constants.TitleSelectProfile
		}
		return constants.TitleSelectRegion
	}

	// Return the title from the map, or empty string if not found
	if title, ok := titleMap[m.CurrentView]; ok {
		return title
	}
	return ""
}

// getHelpText returns the appropriate help text for the current view
func getHelpText(m *model.Model) string {
	// Define common help text patterns
	const (
		defaultHelpText     = "↑/↓: navigate • %s: select • %s: back • %s: quit"
		manualInputHelpText = "%s: confirm • %s: cancel • %s: quit"
		summaryHelpText     = "↑/↓: navigate • %s: select • %s: toggle input • %s: back • %s: quit"
		providersHelpText   = "↑/↓: navigate • %s: select • %s: quit"
	)

	// Special cases based on view and state
	switch {
	case m.CurrentView == constants.ViewProviders:
		return fmt.Sprintf(providersHelpText, constants.KeyEnter, constants.KeyQ)
	case m.CurrentView == constants.ViewAWSConfig && m.ManualInput:
		return fmt.Sprintf(manualInputHelpText, constants.KeyEnter, constants.KeyEsc, constants.KeyCtrlC)
	case m.CurrentView == constants.ViewSummary && m.ManualInput:
		return fmt.Sprintf(manualInputHelpText, constants.KeyEnter, constants.KeyEsc, constants.KeyCtrlC)
	case m.CurrentView == constants.ViewSummary:
		return fmt.Sprintf(summaryHelpText, constants.KeyEnter, constants.KeyTab, constants.KeyEsc, constants.KeyQ)
	default:
		return fmt.Sprintf(defaultHelpText, constants.KeyEnter, constants.KeyEsc, constants.KeyQ)
	}
}

// renderTable renders the table for the current view
func renderTable(m *model.Model) string {
	if m.Table.Rows() == nil {
		return ""
	}

	// Create a table style with appropriate height based on the current view
	// Use padding(top, right, bottom, left) to control spacing
	tableStyle := lipgloss.NewStyle().PaddingTop(1).PaddingRight(2).PaddingBottom(0).PaddingLeft(0)

	// Use larger height for views that need more space
	if m.CurrentView == constants.ViewPipelineStages {
		tableStyle = tableStyle.Height(constants.TableHeightLarge)
	} else {
		tableStyle = tableStyle.Height(constants.TableHeight)
	}

	// Render the table with the appropriate styles
	return tableStyle.Render(m.Table.View())
}
