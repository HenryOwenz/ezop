package views

import (
	"fmt"
	"strings"

	"github.com/HenryOwenz/ezop/internal/ui/model"
)

// View renders the current UI state
func View(m model.Model) string {
	var s strings.Builder

	if m.Error != nil {
		return m.Styles.Error.Render(fmt.Sprintf("Error: %v", m.Error))
	}

	switch m.Step {
	case model.StepSelectProvider:
		s.WriteString(m.Styles.Title.Render("Select Cloud Provider"))
		s.WriteString("\n")

		for i, provider := range m.Providers {
			text := fmt.Sprintf("%s - %s", provider.Name, provider.Description)

			if m.Cursor == i {
				if provider.Available {
					s.WriteString("\n> " + m.Styles.Selected.Render(text))
				} else {
					s.WriteString("\n> " + m.Styles.Disabled.Render(text))
				}
			} else {
				if provider.Available {
					s.WriteString("\n  " + text)
				} else {
					s.WriteString("\n  " + m.Styles.Disabled.Render(text))
				}
			}
		}

	case model.StepProviderConfig:
		if m.SelectedProvider.ID == "aws" {
			if m.AWSProfile == "" {
				s.WriteString(m.Styles.Title.Render("Select AWS Profile"))
				s.WriteString("\n")

				if m.ManualInput {
					s.WriteString(m.Styles.Instruction.Render("Enter AWS Profile: "))
					s.WriteString(m.InputBuffer)
				} else {
					for i, profile := range m.Profiles {
						if m.Cursor == i {
							s.WriteString("\n> " + m.Styles.Selected.Render(profile))
						} else {
							s.WriteString("\n  " + profile)
						}
					}
				}
			} else {
				s.WriteString(m.Styles.Title.Render("Select AWS Region"))
				s.WriteString("\n")
				s.WriteString(m.Styles.Instruction.Render(fmt.Sprintf("Profile: %s", m.AWSProfile)))
				s.WriteString("\n")

				if m.ManualInput {
					s.WriteString(m.Styles.Instruction.Render("Enter AWS Region: "))
					s.WriteString(m.InputBuffer)
				} else {
					for i, region := range m.Regions {
						if m.Cursor == i {
							s.WriteString("\n> " + m.Styles.Selected.Render(region))
						} else {
							s.WriteString("\n  " + region)
						}
					}
				}
			}
		}

	case model.StepSelectService:
		if m.SelectedProvider.ID == "aws" {
			s.WriteString(m.Styles.Title.Render("Select AWS Service"))
			s.WriteString("\n")
			s.WriteString(m.Styles.Instruction.Render(fmt.Sprintf("Profile: %s | Region: %s",
				m.AWSProfile, m.AWSRegion)))

			for i, service := range m.Services {
				text := fmt.Sprintf("%s - %s", service.Name, service.Description)

				if m.Cursor == i {
					if service.Available {
						s.WriteString("\n> " + m.Styles.Selected.Render(text))
					} else {
						s.WriteString("\n> " + m.Styles.Disabled.Render(text))
					}
				} else {
					if service.Available {
						s.WriteString("\n  " + text)
					} else {
						s.WriteString("\n  " + m.Styles.Disabled.Render(text))
					}
				}
			}
		}

	case model.StepServiceOperation:
		if m.SelectedService.ID == "codepipeline" {
			s.WriteString(m.Styles.Title.Render("Select Operation"))
			s.WriteString("\n")
			s.WriteString(m.Styles.Instruction.Render(fmt.Sprintf("Service: %s", m.SelectedService.Name)))

			for i, operation := range m.Operations {
				text := fmt.Sprintf("%s - %s", operation.Name, operation.Description)

				if m.Cursor == i {
					s.WriteString("\n> " + m.Styles.Selected.Render(text))
				} else {
					s.WriteString("\n  " + text)
				}
			}
		}

	case model.StepSelectingApproval:
		s.WriteString(m.Styles.Title.Render("Select Approval"))
		s.WriteString("\n")
		s.WriteString(m.Styles.Instruction.Render(fmt.Sprintf("Profile: %s | Region: %s",
			m.AWSProfile, m.AWSRegion)))

		if len(m.Approvals) == 0 {
			s.WriteString("\n\n")
			s.WriteString(m.Styles.Instruction.Render("No pending approvals found"))
			s.WriteString("\n\nPress q to quit")
			return s.String()
		}

		for i, approval := range m.Approvals {
			text := fmt.Sprintf("%s → %s → %s",
				approval.PipelineName,
				approval.StageName,
				approval.ActionName)

			if m.Cursor == i {
				s.WriteString("\n> " + m.Styles.Selected.Render(text))
			} else {
				s.WriteString("\n  " + text)
			}
		}

	case model.StepConfirmingAction:
		s.WriteString(m.Styles.Title.Render("Choose Action"))
		s.WriteString("\n")
		s.WriteString(m.Styles.Instruction.Render(fmt.Sprintf("Pipeline: %s\nStage: %s\nAction: %s",
			m.SelectedApproval.PipelineName,
			m.SelectedApproval.StageName,
			m.SelectedApproval.ActionName)))

		options := []string{"Approve", "Reject", "Cancel"}
		for i, option := range options {
			if m.Cursor == i {
				s.WriteString("\n> " + m.Styles.Selected.Render(option))
			} else {
				s.WriteString("\n  " + option)
			}
		}

	case model.StepSummaryInput:
		s.WriteString(m.Styles.Title.Render("Enter Summary"))
		s.WriteString("\n")
		s.WriteString(m.Styles.Instruction.Render(fmt.Sprintf("Action: %s", m.Action)))
		s.WriteString("\n\nSummary: ")
		s.WriteString(m.Summary)
		s.WriteString("_")

	case model.StepExecutingAction:
		s.WriteString(m.Styles.Title.Render("Confirm Action"))
		s.WriteString("\n")
		s.WriteString(m.Styles.Instruction.Render(fmt.Sprintf(`Pipeline: %s
Stage: %s
Action: %s
Operation: %s
Summary: %s

Are you sure?`,
			m.SelectedApproval.PipelineName,
			m.SelectedApproval.StageName,
			m.SelectedApproval.ActionName,
			m.Action,
			m.Summary)))

		options := []string{"Yes", "No"}
		for i, option := range options {
			if m.Cursor == i {
				s.WriteString("\n> " + m.Styles.Selected.Render(option))
			} else {
				s.WriteString("\n  " + option)
			}
		}
	}

	// Help text
	s.WriteString("\n\n")
	if m.Step <= 1 {
		s.WriteString(m.Styles.Instruction.Render("↑/↓: Navigate • Enter: Select • Tab: Toggle Input Mode • q: Quit"))
	} else {
		s.WriteString(m.Styles.Instruction.Render("↑/↓: Navigate • Enter: Select • q: Quit"))
	}

	return s.String()
}
