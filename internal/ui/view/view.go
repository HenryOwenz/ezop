package view

import (
	"fmt"

	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/core"
	"github.com/charmbracelet/lipgloss"
)

// Render renders the UI
func Render(m *core.Model) string {
	if m.Err != nil {
		return m.Styles.App.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.Styles.Error.Render("Error: "+m.Err.Error()),
				"\n",
				m.Styles.Help.Render(fmt.Sprintf("%s: quit • %s: back", constants.KeyQ, constants.KeyAltBack)),
			),
		)
	}

	// Create content array with appropriate spacing
	content := make([]string, constants.AppContentLines)

	// Set the title and context
	content[0] = m.Styles.Title.Render(getTitleText(m))
	content[1] = m.Styles.Context.Render(getContextText(m))

	// Add loading spinner if needed
	if m.IsLoading {
		content[2] = m.Spinner.View()
	}

	// For Summary view with approvals, always show text input
	if m.CurrentView == constants.ViewSummary && m.SelectedApproval != nil {
		content[3] = m.TextInput.View()
	} else {
		// For other views, follow normal logic
		if !m.ManualInput {
			content[3] = renderTable(m)
		}
		if m.ManualInput {
			content[3] = m.TextInput.View()
		}
	}

	// Add help text
	content[4] = m.Styles.Help.Render(getHelpText(m))

	// Join all content vertically with consistent spacing
	return m.Styles.App.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			content...,
		),
	)
}

// getContextText returns the appropriate context text for the current view
func getContextText(m *core.Model) string {
	switch m.CurrentView {
	case constants.ViewProviders:
		return constants.MsgAppDescription
	case constants.ViewAWSConfig:
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
	case constants.ViewSelectService:
		return fmt.Sprintf("Profile: %s\nRegion: %s",
			m.AwsProfile,
			m.AwsRegion)
	case constants.ViewSelectCategory:
		if m.SelectedService == nil {
			return ""
		}
		return fmt.Sprintf("Service: %s",
			m.SelectedService.Name)
	case constants.ViewSelectOperation:
		if m.SelectedService == nil || m.SelectedCategory == nil {
			return ""
		}
		return fmt.Sprintf("Service: %s\nCategory: %s",
			m.SelectedService.Name,
			m.SelectedCategory.Name)
	case constants.ViewApprovals:
		return fmt.Sprintf("Profile: %s\nRegion: %s",
			m.AwsProfile,
			m.AwsRegion)
	case constants.ViewConfirmation, constants.ViewSummary:
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
	case constants.ViewExecutingAction:
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
	case constants.ViewPipelineStatus:
		return fmt.Sprintf("Profile: %s\nRegion: %s",
			m.AwsProfile,
			m.AwsRegion)
	case constants.ViewPipelineStages:
		if m.SelectedPipeline == nil {
			return ""
		}
		return fmt.Sprintf("Profile: %s\nRegion: %s\nPipeline: %s",
			m.AwsProfile,
			m.AwsRegion,
			m.SelectedPipeline.Name)
	default:
		return ""
	}
}

// getTitleText returns the appropriate title for the current view
func getTitleText(m *core.Model) string {
	switch m.CurrentView {
	case constants.ViewProviders:
		return constants.TitleProviders
	case constants.ViewAWSConfig:
		if m.AwsProfile == "" {
			return constants.TitleSelectProfile
		}
		return constants.TitleSelectRegion
	case constants.ViewSelectService:
		return constants.TitleSelectService
	case constants.ViewSelectCategory:
		return constants.TitleSelectCategory
	case constants.ViewSelectOperation:
		return constants.TitleSelectOperation
	case constants.ViewApprovals:
		return constants.TitleApprovals
	case constants.ViewConfirmation:
		return constants.TitleConfirmation
	case constants.ViewSummary:
		return constants.TitleSummary
	case constants.ViewExecutingAction:
		return constants.TitleExecutingAction
	case constants.ViewPipelineStatus:
		return constants.TitlePipelineStatus
	case constants.ViewPipelineStages:
		return constants.TitlePipelineStages
	case constants.ViewError:
		return constants.TitleError
	case constants.ViewSuccess:
		return constants.TitleSuccess
	case constants.ViewHelp:
		return constants.TitleHelp
	default:
		return ""
	}
}

// getHelpText returns the appropriate help text for the current view
func getHelpText(m *core.Model) string {
	switch {
	case m.CurrentView == constants.ViewProviders:
		return fmt.Sprintf("↑/↓: navigate • %s: select • %s: quit", constants.KeyEnter, constants.KeyQ)
	case m.CurrentView == constants.ViewAWSConfig && m.ManualInput:
		return fmt.Sprintf("%s: confirm • %s: cancel • %s: quit", constants.KeyEnter, constants.KeyEsc, constants.KeyCtrlC)
	case m.CurrentView == constants.ViewAWSConfig:
		return fmt.Sprintf("↑/↓: navigate • %s: select • %s: back • %s: quit", constants.KeyEnter, constants.KeyEsc, constants.KeyQ)
	case m.CurrentView == constants.ViewSummary && m.ManualInput:
		return fmt.Sprintf("%s: confirm • %s: cancel • %s: quit", constants.KeyEnter, constants.KeyEsc, constants.KeyCtrlC)
	case m.CurrentView == constants.ViewSummary:
		return fmt.Sprintf("↑/↓: navigate • %s: select • %s: toggle input • %s: back • %s: quit", constants.KeyEnter, constants.KeyTab, constants.KeyEsc, constants.KeyQ)
	default:
		return fmt.Sprintf("↑/↓: navigate • %s: select • %s: back • %s: quit", constants.KeyEnter, constants.KeyEsc, constants.KeyQ)
	}
}

// renderTable renders the table for the current view
func renderTable(m *core.Model) string {
	if m.Table.Rows() == nil {
		return ""
	}

	// Create a table style with appropriate height based on the current view
	tableStyle := lipgloss.NewStyle().Padding(1, 2)

	// Use larger height for views that need more space
	if m.CurrentView == constants.ViewPipelineStages {
		tableStyle = tableStyle.Height(constants.TableHeightLarge)
	} else {
		tableStyle = tableStyle.Height(constants.TableHeight)
	}

	// Render the table with the appropriate styles
	return tableStyle.Render(m.Table.View())
}
