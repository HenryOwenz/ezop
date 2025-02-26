package ui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
)

func (m *Model) View() string {
	if m.err != nil {
		return m.styles.App.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.styles.Error.Render("Error: "+m.err.Error()),
				"\n",
				m.styles.Help.Render("q: quit • -: back"),
			),
		)
	}

	content := []string{
		m.styles.Title.Render(m.getTitleText()),
		m.styles.Context.Render(m.getContextText()),
		"",
		"",
		"", // Empty line for help text
	}

	// Add loading spinner if needed
	if m.isLoading {
		content[2] = m.spinner.View()
	}

	// For Summary view with approvals, always show text input
	if m.currentView == constants.ViewSummary && m.selectedApproval != nil {
		content[3] = m.textInput.View()
	} else {
		// For other views, follow normal logic
		if !m.manualInput {
			content[3] = m.table.View()
		}
		if m.manualInput {
			content[3] = m.textInput.View()
		}
	}

	// Add help text
	content[4] = m.styles.Help.Render(m.getHelpText())

	// Join all content vertically with consistent spacing
	return m.styles.App.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			content...,
		),
	)
}

// getTitleText returns the appropriate title for the current view
func (m *Model) getTitleText() string {
	switch m.currentView {
	case constants.ViewProviders:
		return "Select Cloud Provider"
	case constants.ViewAWSConfig:
		if m.awsProfile == "" {
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

// getContextText returns the appropriate context text for the current view
func (m *Model) getContextText() string {
	switch m.currentView {
	case constants.ViewProviders:
		return "A simple tool to manage your cloud resources"
	case constants.ViewAWSConfig:
		if m.awsProfile == "" {
			return "Amazon Web Services"
		}
		return fmt.Sprintf("Profile: %s", m.awsProfile)
	case constants.ViewSelectService:
		return fmt.Sprintf("Profile: %s\nRegion: %s",
			m.awsProfile,
			m.awsRegion)
	case constants.ViewSelectCategory:
		if m.selectedService == nil {
			return ""
		}
		return fmt.Sprintf("Service: %s",
			m.selectedService.Name)
	case constants.ViewSelectOperation:
		if m.selectedService == nil || m.selectedCategory == nil {
			return ""
		}
		return fmt.Sprintf("Service: %s\nCategory: %s",
			m.selectedService.Name,
			m.selectedCategory.Name)
	case constants.ViewApprovals:
		return fmt.Sprintf("Profile: %s\nRegion: %s",
			m.awsProfile,
			m.awsRegion)
	case constants.ViewConfirmation, constants.ViewSummary:
		if m.selectedOperation != nil && m.selectedOperation.Name == "Start Pipeline" {
			if m.selectedPipeline == nil {
				return ""
			}
			return fmt.Sprintf("Profile: %s\nRegion: %s\nPipeline: %s",
				m.awsProfile,
				m.awsRegion,
				m.selectedPipeline.Name)
		}
		if m.selectedApproval == nil {
			return ""
		}
		return fmt.Sprintf("Pipeline: %s\nStage: %s\nAction: %s",
			m.selectedApproval.PipelineName,
			m.selectedApproval.StageName,
			m.selectedApproval.ActionName)
	case constants.ViewExecutingAction:
		if m.selectedOperation != nil && m.selectedOperation.Name == "Start Pipeline" {
			if m.selectedPipeline == nil {
				return ""
			}
			return fmt.Sprintf("Profile: %s\nRegion: %s\nPipeline: %s\nRevisionID: Latest commit",
				m.awsProfile,
				m.awsRegion,
				m.selectedPipeline.Name)
		}
		if m.selectedApproval == nil {
			return ""
		}
		return fmt.Sprintf("Pipeline: %s\nStage: %s\nAction: %s\nComment: %s",
			m.selectedApproval.PipelineName,
			m.selectedApproval.StageName,
			m.selectedApproval.ActionName,
			m.summary)
	case constants.ViewPipelineStatus:
		return fmt.Sprintf("Profile: %s\nRegion: %s",
			m.awsProfile,
			m.awsRegion)
	case constants.ViewPipelineStages:
		if m.selectedPipeline == nil {
			return ""
		}
		return fmt.Sprintf("Profile: %s\nRegion: %s\nPipeline: %s",
			m.awsProfile,
			m.awsRegion,
			m.selectedPipeline.Name)
	default:
		return ""
	}
}

// getHelpText returns the appropriate help text for the current view
func (m *Model) getHelpText() string {
	switch {
	case m.currentView == constants.ViewProviders:
		return "↑/↓: navigate • enter: select • q: quit"
	case m.currentView == constants.ViewAWSConfig && m.manualInput:
		return "enter: confirm • esc: cancel • ctrl+c: quit"
	case m.currentView == constants.ViewAWSConfig:
		return "↑/↓: navigate • enter: select • tab: toggle input • esc: back • q: quit"
	case m.currentView == constants.ViewSummary && m.manualInput:
		return "enter: confirm • esc: cancel • ctrl+c: quit"
	case m.currentView == constants.ViewSummary:
		return "↑/↓: navigate • enter: select • tab: toggle input • esc: back • q: quit"
	default:
		return "↑/↓: navigate • enter: select • esc: back • q: quit"
	}
} 