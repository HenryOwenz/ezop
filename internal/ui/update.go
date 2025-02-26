package ui

import (
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil
	case errMsg:
		return m.handleError(msg)
	case approvalsMsg:
		return m.handleApprovals(msg)
	case approvalResultMsg:
		return m.handleApprovalResult(msg)
	case pipelineExecutionMsg:
		return m.handlePipelineExecution(msg)
	case spinner.TickMsg:
		return m.handleSpinnerTick(msg)
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case pipelineStatusMsg:
		return m.handlePipelineStatus(msg)
	}
	return m, nil
}

func (m *Model) handleError(msg errMsg) (tea.Model, tea.Cmd) {
	newModel := *m
	newModel.err = msg.err
	newModel.isLoading = false
	return &newModel, nil
}

func (m *Model) handleApprovals(msg approvalsMsg) (tea.Model, tea.Cmd) {
	newModel := *m
	newModel.approvals = msg.approvals
	newModel.provider = msg.provider
	newModel.currentView = constants.ViewApprovals
	newModel.isLoading = false
	newModel.updateTableForView()
	return &newModel, nil
}

func (m *Model) handleApprovalResult(msg approvalResultMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		return m, func() tea.Msg {
			return errMsg(msg)
		}
	}
	// First clear loading state
	newModel := *m
	newModel.isLoading = false
	// Then reset approval state and navigate
	newModel.currentView = constants.ViewSelectCategory
	newModel.resetApprovalState()
	// Clear text input
	newModel.resetTextInput()
	newModel.updateTableForView()
	return &newModel, nil
}

func (m *Model) handlePipelineExecution(msg pipelineExecutionMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		return m, func() tea.Msg {
			return errMsg(msg)
		}
	}
	newModel := *m
	newModel.isLoading = false
	newModel.currentView = constants.ViewSelectCategory
	newModel.selectedPipeline = nil
	newModel.selectedOperation = nil
	newModel.resetTextInput()
	newModel.updateTableForView()
	return &newModel, nil
}

func (m *Model) handleSpinnerTick(msg spinner.TickMsg) (tea.Model, tea.Cmd) {
	if m.isLoading {
		var cmd tea.Cmd
		newModel := *m
		newModel.spinner, cmd = m.spinner.Update(msg)
		return &newModel, cmd
	}
	return m, nil
}

func (m *Model) handlePipelineStatus(msg pipelineStatusMsg) (tea.Model, tea.Cmd) {
	newModel := *m
	newModel.pipelines = msg.pipelines
	newModel.provider = msg.provider
	newModel.currentView = constants.ViewPipelineStatus
	newModel.isLoading = false
	newModel.updateTableForView()
	return &newModel, nil
}
