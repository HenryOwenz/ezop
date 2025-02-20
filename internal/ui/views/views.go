package views

import (
	"fmt"
	"strings"

	"github.com/HenryOwenz/ezop/internal/ui/model"
	"github.com/charmbracelet/lipgloss"
)

// View renders the current UI state
func View(m model.Model) string {
	if m.Error != nil {
		return m.Styles.Error.Render(fmt.Sprintf("Error: %v", m.Error))
	}

	var content strings.Builder

	// Add loading spinner at the top if needed
	if m.IsLoading {
		content.WriteString(m.Spinner.View())
		content.WriteString("\n\n")
	}

	// Add title based on current step
	switch m.Step {
	case model.StepSelectProvider:
		content.WriteString(m.Styles.Title.Render("Select Cloud Provider"))
	case model.StepProviderConfig:
		if m.AWSProfile == "" {
			content.WriteString(m.Styles.Title.Render("Select AWS Profile"))
		} else {
			content.WriteString(m.Styles.Title.Render("Select AWS Region"))
			content.WriteString("\n")
			content.WriteString(m.Styles.Instruction.Render(fmt.Sprintf("Profile: %s", m.AWSProfile)))
		}
	case model.StepSelectService:
		content.WriteString(m.Styles.Title.Render("Select AWS Service"))
		content.WriteString("\n")
		content.WriteString(m.Styles.Instruction.Render(fmt.Sprintf("Profile: %s | Region: %s",
			m.AWSProfile, m.AWSRegion)))
	case model.StepSelectCategory:
		content.WriteString(m.Styles.Title.Render("Select Category"))
		content.WriteString("\n")
		content.WriteString(m.Styles.Instruction.Render(fmt.Sprintf("Service: %s", m.SelectedService.Name)))
	case model.StepSelectingApproval:
		content.WriteString(m.Styles.Title.Render("Select Pipeline"))
		content.WriteString("\n")
		content.WriteString(m.Styles.Instruction.Render("Service: " + m.SelectedService.Name))
		content.WriteString(m.Styles.Instruction.Render("Operation: " + m.SelectedOperation.Name))
	case model.StepConfirmingAction:
		content.WriteString(m.Styles.Title.Render("Confirm Action"))
		content.WriteString("\n")
		content.WriteString(m.Styles.Instruction.Render("Pipeline: " + m.SelectedApproval.PipelineName))
		content.WriteString(m.Styles.Instruction.Render("Stage: " + m.SelectedApproval.StageName))
		content.WriteString(m.Styles.Instruction.Render("Action: " + m.SelectedApproval.ActionName))
	case model.StepExecutingAction:
		content.WriteString(m.Styles.Title.Render("Execute Action"))
		content.WriteString("\n")
		content.WriteString(m.Styles.Instruction.Render("Pipeline: " + m.SelectedApproval.PipelineName))
		content.WriteString(m.Styles.Instruction.Render("Stage: " + m.SelectedApproval.StageName))
		content.WriteString(m.Styles.Instruction.Render("Action: " + m.SelectedApproval.ActionName))
		content.WriteString(m.Styles.Instruction.Render("Summary: " + m.Summary))
	case model.StepSummaryInput:
		content.WriteString(m.Styles.Title.Render(fmt.Sprintf("%s Pipeline", strings.Title(m.Action))))
		content.WriteString("\n")
		content.WriteString(m.Styles.Instruction.Render("Pipeline: " + m.SelectedApproval.PipelineName))
		content.WriteString(m.Styles.Instruction.Render("Stage: " + m.SelectedApproval.StageName))
		content.WriteString(m.Styles.Instruction.Render("Action: " + m.SelectedApproval.ActionName))
	case model.StepServiceOperation:
		content.WriteString(m.Styles.Title.Render("Select Operation"))
		content.WriteString("\n")
		content.WriteString(m.Styles.Instruction.Render("Service: " + m.SelectedService.Name))
		content.WriteString(m.Styles.Instruction.Render("Category: " + m.SelectedCategory.Name))
	}
	content.WriteString("\n\n")

	// Add table view if we're not in manual input mode
	if !m.ManualInput && m.Step != model.StepSummaryInput {
		content.WriteString(m.Table.View())
	} else if m.ManualInput {
		content.WriteString(m.Styles.Instruction.Render("Enter value: "))
		content.WriteString(m.InputBuffer)
	} else if m.Step == model.StepSummaryInput {
		content.WriteString(m.Styles.Instruction.Render(fmt.Sprintf("Action: %s", m.Action)))
		content.WriteString("\n\nSummary: ")
		content.WriteString(m.Summary)
		content.WriteString("_")
	}

	// Add help text at the bottom
	helpText := renderHelpText(m)

	// Combine content and help text
	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		content.String(),
		"\n"+helpText,
	)

	// Wrap in content area and frame
	return m.Styles.Frame.Render(
		m.Styles.ContentArea.Render(mainContent),
	)
}

// renderHelpText returns the help text based on the current step
func renderHelpText(m model.Model) string {
	if m.Step <= 1 {
		return m.Styles.Instruction.Render("↑/↓: Navigate • Enter: Select • Tab: Toggle Input Mode • q: Quit")
	}
	return m.Styles.Instruction.Render("↑/↓: Navigate • Enter: Select • -: Back • q: Quit")
}
