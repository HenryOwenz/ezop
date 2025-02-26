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
				m.Styles.Help.Render("q: quit • -: back"),
			),
		)
	}

	content := []string{
		m.Styles.Title.Render(getTitleText(m)),
		m.Styles.Context.Render(getContextText(m)),
		"",
		"",
		"", // Empty line for help text
	}

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
			content[3] = m.Table.View()
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
		return "A simple tool to manage your cloud resources"
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
		return "Select Cloud Provider"
	case constants.ViewAWSConfig:
		if m.AwsProfile == "" {
			return "Select AWS Profile"
		}
		return "Select AWS Region"
	case constants.ViewSelectService:
		return "Select AWS Service"
	case constants.ViewSelectCategory:
		return "Select Category"
	case constants.ViewSelectOperation:
		return "Select Operation"
	case constants.ViewApprovals:
		return "Pipeline Approvals"
	case constants.ViewConfirmation:
		return "Execute Action"
	case constants.ViewSummary:
		return "Enter Comment"
	case constants.ViewExecutingAction:
		return "Execute Action"
	case constants.ViewPipelineStatus:
		return "Select Pipeline"
	case constants.ViewPipelineStages:
		return "Pipeline Stages"
	default:
		return ""
	}
}

// getHelpText returns the appropriate help text for the current view
func getHelpText(m *core.Model) string {
	switch {
	case m.CurrentView == constants.ViewProviders:
		return "↑/↓: navigate • enter: select • q: quit"
	case m.CurrentView == constants.ViewAWSConfig && m.ManualInput:
		return "enter: confirm • esc: cancel • ctrl+c: quit"
	case m.CurrentView == constants.ViewAWSConfig:
		return "↑/↓: navigate • enter: select • esc: back • q: quit"
	case m.CurrentView == constants.ViewSummary && m.ManualInput:
		return "enter: confirm • esc: cancel • ctrl+c: quit"
	case m.CurrentView == constants.ViewSummary:
		return "↑/↓: navigate • enter: select • tab: toggle input • esc: back • q: quit"
	default:
		return "↑/↓: navigate • enter: select • esc: back • q: quit"
	}
}
