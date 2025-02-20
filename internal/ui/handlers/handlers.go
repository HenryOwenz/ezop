package handlers

import (
	"context"

	"github.com/HenryOwenz/ezop/internal/providers/aws"
	"github.com/HenryOwenz/ezop/internal/ui/model"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model accordingly
func Update(m model.Model, msg tea.Msg) (model.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return handleKeyPress(m, msg)
	default:
		return m, nil
	}
}

// handleKeyPress handles keyboard input events
func handleKeyPress(m model.Model, msg tea.KeyMsg) (model.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "-":
		if m.Step > model.StepSelectProvider {
			m.NavigateBack()
		}
		return m, nil

	case "up", "k":
		if !m.ManualInput && m.Cursor > 0 {
			m.Cursor--
		}

	case "down", "j":
		if !m.ManualInput {
			switch m.Step {
			case model.StepSelectProvider:
				if m.Cursor < len(m.Providers)-1 {
					m.Cursor++
				}
			case model.StepProviderConfig:
				if m.AWSProfile == "" {
					if m.Cursor < len(m.Profiles)-1 {
						m.Cursor++
					}
				} else {
					if m.Cursor < len(m.Regions)-1 {
						m.Cursor++
					}
				}
			case model.StepSelectService:
				if m.Cursor < len(m.Services)-1 {
					m.Cursor++
				}
			case model.StepServiceOperation:
				if m.Cursor < len(m.Operations)-1 {
					m.Cursor++
				}
			case model.StepSelectingApproval:
				if m.Cursor < len(m.Approvals)-1 {
					m.Cursor++
				}
			case model.StepConfirmingAction:
				if m.Cursor < 2 { // Three options: Approve, Reject, Cancel
					m.Cursor++
				}
			case model.StepExecutingAction:
				if m.Cursor < 1 { // Two options: Yes, No
					m.Cursor++
				}
			}
		}

	case "tab":
		if m.Step == model.StepProviderConfig {
			m.ManualInput = !m.ManualInput
			m.InputBuffer = ""
			m.Cursor = 0
		}

	case "enter":
		return handleEnterPress(m)

	case "backspace":
		if m.ManualInput && len(m.InputBuffer) > 0 {
			m.InputBuffer = m.InputBuffer[:len(m.InputBuffer)-1]
		} else if m.Step == model.StepSummaryInput && len(m.Summary) > 0 {
			m.Summary = m.Summary[:len(m.Summary)-1]
		}

	default:
		if m.ManualInput {
			m.InputBuffer += msg.String()
		} else if m.Step == model.StepSummaryInput {
			m.Summary += msg.String()
		}
	}

	return m, nil
}

// handleEnterPress handles the enter key press based on the current step
func handleEnterPress(m model.Model) (model.Model, tea.Cmd) {
	switch m.Step {
	case model.StepSelectProvider:
		provider := m.Providers[m.Cursor]
		if !provider.Available {
			return m, nil
		}
		m.SelectedProvider = &provider
		m.Step = model.StepProviderConfig
		m.Cursor = 0

	case model.StepProviderConfig:
		if m.SelectedProvider.ID == "aws" {
			if m.AWSProfile == "" {
				if m.ManualInput {
					if m.InputBuffer != "" {
						m.AWSProfile = m.InputBuffer
						m.InputBuffer = ""
						m.Cursor = 0
					}
				} else if len(m.Profiles) > 0 {
					m.AWSProfile = m.Profiles[m.Cursor]
					m.Cursor = 0
				}
				return m, nil
			}

			if m.AWSRegion == "" {
				if m.ManualInput {
					if m.InputBuffer != "" {
						m.AWSRegion = m.InputBuffer
						m.InputBuffer = ""
						m.ManualInput = false
						m.Step = model.StepSelectService
						m.Cursor = 0
					}
				} else if len(m.Regions) > 0 {
					m.AWSRegion = m.Regions[m.Cursor]
					m.Step = model.StepSelectService
					m.Cursor = 0
				}

				// Initialize provider after region is set
				if m.AWSRegion != "" {
					provider, err := aws.NewProvider(m.AWSProfile, m.AWSRegion)
					if err != nil {
						m.Error = err
						return m, nil
					}
					m.Services = provider.GetServices()
					m.AWSProvider = provider
				}
				return m, nil
			}
		}

	case model.StepSelectService:
		service := m.Services[m.Cursor]
		if !service.Available {
			return m, nil
		}
		m.SelectedService = &service
		m.Operations = m.AWSProvider.GetOperations(service.ID)
		m.Step = model.StepServiceOperation
		m.Cursor = 0

	case model.StepServiceOperation:
		operation := m.Operations[m.Cursor]
		m.SelectedOperation = &operation
		if operation.ID == "manual-approval" {
			return initApprovals(m)
		}

	case model.StepSelectingApproval:
		if len(m.Approvals) > 0 {
			m.SelectedApproval = &m.Approvals[m.Cursor]
			m.Step = model.StepConfirmingAction
			m.Cursor = 0
		}

	case model.StepConfirmingAction:
		switch m.Cursor {
		case 0: // Approve
			m.Action = "approve"
			m.Step = model.StepSummaryInput
		case 1: // Reject
			m.Action = "reject"
			m.Step = model.StepSummaryInput
		case 2: // Cancel
			return m, tea.Quit
		}

	case model.StepSummaryInput:
		if m.Summary != "" {
			m.Step = model.StepExecutingAction
			m.Cursor = 0
		}

	case model.StepExecutingAction:
		if m.Cursor == 0 { // Yes
			return executeAction(m)
		}
		return m, tea.Quit // No
	}

	return m, nil
}

func initApprovals(m model.Model) (model.Model, tea.Cmd) {
	approvals, err := m.AWSProvider.GetPendingApprovals(context.Background())
	if err != nil {
		m.Error = err
		return m, nil
	}
	m.Approvals = approvals
	m.Step = model.StepSelectingApproval
	return m, nil
}

func executeAction(m model.Model) (model.Model, tea.Cmd) {
	params := map[string]interface{}{
		"pipeline_name": m.SelectedApproval.PipelineName,
		"stage_name":    m.SelectedApproval.StageName,
		"action_name":   m.SelectedApproval.ActionName,
		"token":         m.SelectedApproval.Token,
		"summary":       m.Summary,
		"approve":       m.Action == "approve",
	}

	err := m.AWSProvider.ExecuteOperation(context.Background(),
		m.SelectedService.ID,
		m.SelectedOperation.ID,
		params)

	if err != nil {
		m.Error = err
		return m, nil
	}
	return m, tea.Quit
}
