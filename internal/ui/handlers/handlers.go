package handlers

import (
	"context"
	"fmt"

	"github.com/HenryOwenz/ezop/internal/domain"
	"github.com/HenryOwenz/ezop/internal/providers/aws"
	"github.com/HenryOwenz/ezop/internal/ui/model"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model accordingly
func Update(m model.Model, msg tea.Msg) (model.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Always handle quit, even during loading
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			return m, tea.Quit
		}
		if m.IsLoading {
			return m, m.Spinner.Tick
		}

		// Handle special keys first
		switch msg.String() {
		case "-":
			if m.Step > model.StepSelectProvider {
				m.NavigateBack()
				m.UpdateTableForStep()
			}
			return m, nil
		case "tab":
			if m.Step == model.StepProviderConfig {
				m.ManualInput = !m.ManualInput
				m.InputBuffer = ""
				return m, nil
			}
		}

		// If we're in manual input or summary input mode, handle those specially
		if m.ManualInput {
			switch msg.String() {
			case "enter":
				return handleEnterPress(m)
			case "backspace":
				if len(m.InputBuffer) > 0 {
					m.InputBuffer = m.InputBuffer[:len(m.InputBuffer)-1]
				}
			default:
				m.InputBuffer += msg.String()
			}
			return m, nil
		} else if m.Step == model.StepSummaryInput {
			switch msg.String() {
			case "enter":
				return handleEnterPress(m)
			case "backspace":
				if len(m.Summary) > 0 {
					m.Summary = m.Summary[:len(m.Summary)-1]
				}
			default:
				m.Summary += msg.String()
			}
			return m, nil
		}

		// For table views, let the table handle all navigation
		var cmd tea.Cmd
		m.Table, cmd = m.Table.Update(msg)

		// After table update, sync our cursor
		m.Cursor = m.Table.Cursor()

		// Handle enter separately
		if msg.String() == "enter" {
			if m.Step == model.StepConfirmingAction {
				switch m.Table.SelectedRow()[0] {
				case "Approve":
					m.Action = "approve"
					m.Step = model.StepSummaryInput
					m.UpdateTableForStep()
					return m, nil
				case "Reject":
					m.Action = "reject"
					m.Step = model.StepSummaryInput
					m.UpdateTableForStep()
					return m, nil
				case "Cancel":
					return m, tea.Quit
				}
			} else if m.Step == model.StepExecutingAction {
				switch m.Table.SelectedRow()[0] {
				case "Execute":
					m.IsLoading = true
					m.LoadingMsg = fmt.Sprintf("%sing pipeline...", m.Action)
					return m, tea.Batch(
						m.Spinner.Tick,
						func() tea.Msg {
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
								return err
							}
							return actionCompleteMsg{}
						},
					)
				case "Cancel":
					return m, tea.Quit
				}
			} else if selected := m.GetSelectedRow(); selected != nil {
				return handleEnterPress(m)
			}
		}

		return m, cmd

	case spinner.TickMsg:
		if m.IsLoading {
			var cmd tea.Cmd
			m.Spinner, cmd = m.Spinner.Update(msg)
			return m, cmd
		}
		return m, nil
	case error:
		m.Error = msg
		m.IsLoading = false
		m.LoadingMsg = ""
		return m, nil
	case awsProviderMsg:
		m.Services = msg.services
		m.AWSProvider = msg.provider
		m.IsLoading = false
		m.LoadingMsg = ""
		m.Step = model.StepSelectService
		m.UpdateTableForStep()
		return m, nil
	case approvalsMsg:
		m.Approvals = msg.approvals
		m.IsLoading = false
		m.LoadingMsg = ""
		m.Step = model.StepSelectingApproval
		m.UpdateTableForStep()
		return m, nil
	case actionCompleteMsg:
		m.IsLoading = false
		m.LoadingMsg = ""
		return m, tea.Quit
	default:
		return m, nil
	}
}

// Custom message types for handling async operations
type awsProviderMsg struct {
	provider *aws.Provider
	services []domain.Service
}

type approvalsMsg struct {
	approvals []aws.ApprovalAction
}

type actionCompleteMsg struct{}

// handleEnterPress handles the enter key press based on the current step
func handleEnterPress(m model.Model) (model.Model, tea.Cmd) {
	switch m.Step {
	case model.StepSelectProvider:
		if selected := m.GetSelectedRow(); selected != nil {
			provider := selected.(domain.Provider)
			if !provider.Available {
				return m, nil
			}
			m.SelectedProvider = &provider
			m.Step = model.StepProviderConfig
			m.UpdateTableForStep()
		}

	case model.StepProviderConfig:
		if m.SelectedProvider.ID == "aws" {
			if m.AWSProfile == "" {
				if m.ManualInput {
					if m.InputBuffer != "" {
						m.AWSProfile = m.InputBuffer
						m.InputBuffer = ""
					}
				} else if selected := m.GetSelectedRow(); selected != nil {
					m.AWSProfile = selected.(string)
				}
				m.UpdateTableForStep()
				return m, nil
			}

			if m.AWSRegion == "" {
				if m.ManualInput {
					if m.InputBuffer != "" {
						m.AWSRegion = m.InputBuffer
						m.InputBuffer = ""
						m.ManualInput = false
						m.Step = model.StepSelectService
					}
				} else if selected := m.GetSelectedRow(); selected != nil {
					m.AWSRegion = selected.(string)
					m.Step = model.StepSelectService
				}

				// Initialize provider after region is set
				if m.AWSRegion != "" {
					m.IsLoading = true
					m.LoadingMsg = "Initializing AWS provider..."
					return m, tea.Batch(
						m.Spinner.Tick,
						func() tea.Msg {
							provider, err := aws.NewProvider(m.AWSProfile, m.AWSRegion)
							if err != nil {
								return err
							}
							services := provider.GetServices()
							return awsProviderMsg{
								provider: provider,
								services: services,
							}
						},
					)
				}
				return m, nil
			}
		}

	case model.StepSelectService:
		if selected := m.GetSelectedRow(); selected != nil {
			service := selected.(domain.Service)
			if !service.Available {
				return m, nil
			}
			m.SelectedService = &service
			m.Step = model.StepSelectCategory
			m.UpdateTableForStep()
		}

	case model.StepSelectCategory:
		if selected := m.GetSelectedRow(); selected != nil {
			category := selected.(domain.Category)
			if !category.Available {
				return m, nil
			}
			m.SelectedCategory = &category

			if category.ID == "workflows" {
				m.Operations = m.AWSProvider.GetOperations(m.SelectedService.ID)
				m.Step = model.StepServiceOperation
				m.UpdateTableForStep()
			} else {
				// Operations mode not implemented yet
				m.Error = fmt.Errorf("direct operations mode not yet implemented")
			}
		}

	case model.StepServiceOperation:
		if selected := m.GetSelectedRow(); selected != nil {
			operation := selected.(domain.Operation)
			m.SelectedOperation = &operation
			if operation.ID == "manual-approval" {
				m.IsLoading = true
				m.LoadingMsg = "Fetching pending approvals..."
				return m, tea.Batch(
					m.Spinner.Tick,
					func() tea.Msg {
						approvals, err := m.AWSProvider.GetPendingApprovals(context.Background())
						if err != nil {
							return err
						}
						return approvalsMsg{approvals: approvals}
					},
				)
			}
		}

	case model.StepSelectingApproval:
		if selected := m.GetSelectedRow(); selected != nil {
			approval := selected.(aws.ApprovalAction)
			m.SelectedApproval = &approval
			m.Step = model.StepConfirmingAction
			m.Cursor = 0
			m.UpdateTableForStep()
		}

	case model.StepConfirmingAction:
		switch m.Cursor {
		case 0: // Approve
			m.Action = "approve"
			m.Step = model.StepSummaryInput
			m.Cursor = 0
			m.UpdateTableForStep()
		case 1: // Reject
			m.Action = "reject"
			m.Step = model.StepSummaryInput
			m.Cursor = 0
			m.UpdateTableForStep()
		case 2: // Cancel
			return m, tea.Quit
		}

	case model.StepSummaryInput:
		if !m.ManualInput {
			switch m.Cursor {
			case 0: // Confirm
				if m.Summary != "" {
					m.Step = model.StepExecutingAction
					m.Cursor = 0
					m.UpdateTableForStep()
				}
			case 1: // Cancel
				return m, tea.Quit
			}
		}

	case model.StepExecutingAction:
		// Handled in Update function
		return m, nil
	}

	return m, nil
}
