package ui

import (
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/core"
	"github.com/HenryOwenz/cloudgate/internal/ui/handlers"
	"github.com/HenryOwenz/cloudgate/internal/ui/navigation"
	"github.com/HenryOwenz/cloudgate/internal/ui/update"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Model is the main UI model that implements the tea.Model interface
type Model struct {
	core *core.Model
}

// New creates a new UI model
func New() Model {
	m := Model{
		core: core.New(),
	}
	// Initialize the table for the current view
	view.UpdateTableForView(m.core)
	return m
}

// Init initializes the UI model
func (m Model) Init() tea.Cmd {
	// Make sure to initialize the table before returning
	view.UpdateTableForView(m.core)
	return m.core.Init()
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.core.Width = msg.Width
		m.core.Height = msg.Height
		view.UpdateTableForView(m.core)
		return m, nil
	case core.ErrMsg:
		m.core.Err = msg.Err
		m.core.IsLoading = false
		return m, nil
	case core.ApprovalsMsg:
		m.core.Approvals = msg.Approvals
		m.core.Provider = msg.Provider
		m.core.CurrentView = constants.ViewApprovals
		m.core.IsLoading = false
		view.UpdateTableForView(m.core)
		return m, nil
	case core.ApprovalResultMsg:
		m.core.IsLoading = false // Ensure loading is turned off
		update.HandleApprovalResult(m.core, msg.Err)
		view.UpdateTableForView(m.core)
		return m, nil
	case core.PipelineExecutionMsg:
		m.core.IsLoading = false // Ensure loading is turned off
		update.HandlePipelineExecution(m.core, msg.Err)
		view.UpdateTableForView(m.core)
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.core.Spinner, cmd = m.core.Spinner.Update(msg)
		if m.core.IsLoading {
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	case tea.KeyMsg:
		// Ignore navigation key presses when loading
		if m.core.IsLoading {
			// Only allow quit commands during loading
			switch msg.String() {
			case constants.KeyCtrlC, constants.KeyQ:
				return m, tea.Quit
			default:
				// Ignore all other key presses during loading
				return m, nil
			}
		}

		// Handle key presses when not loading
		switch msg.String() {
		case constants.KeyCtrlC, constants.KeyQ:
			return m, tea.Quit
		case constants.KeyEnter:
			modelWrapper, cmd := handlers.HandleEnter(m.core)
			if wrapper, ok := modelWrapper.(handlers.ModelWrapper); ok {
				// Since ModelWrapper embeds *core.Model, we can create a new Model with it
				newModel := Model{core: wrapper.Model}
				if newModel.core.IsLoading {
					return newModel, tea.Batch(cmd, newModel.core.Spinner.Tick)
				}
				return newModel, cmd
			}
			return modelWrapper, cmd
		case constants.KeyEsc, constants.KeyAltBack:
			// Only use '-' for back navigation if not in text input mode
			if msg.String() == constants.KeyAltBack && m.core.ManualInput {
				// If in text input mode, '-' should be treated as a character
				var cmd tea.Cmd
				m.core.TextInput, cmd = m.core.TextInput.Update(msg)
				return m, cmd
			}

			// Handle back navigation
			if m.core.ManualInput {
				m.core.ManualInput = false
				m.core.ResetTextInput()
				view.UpdateTableForView(m.core)
				return m, nil
			}
			model := navigation.NavigateBack(m.core)
			view.UpdateTableForView(model)
			return Model{core: model}, nil
		case constants.KeyUp, constants.KeyAltUp:
			m.core.Table.MoveUp(1)
			return m, nil
		case constants.KeyDown, constants.KeyAltDown:
			m.core.Table.MoveDown(1)
			return m, nil
		case constants.KeyTab:
			// Tab is now only used for other views, not AWS config
			if m.core.CurrentView == constants.ViewSummary {
				if m.core.ManualInput {
					m.core.ManualInput = false
					m.core.ResetTextInput()
				} else {
					m.core.ManualInput = true
					m.core.TextInput.Focus()
				}
				return m, nil
			}
			return m, nil
		default:
			if m.core.ManualInput {
				var cmd tea.Cmd
				m.core.TextInput, cmd = m.core.TextInput.Update(msg)

				// If we're in the summary view with manual commit ID
				if m.core.CurrentView == constants.ViewSummary && m.core.SelectedOperation != nil &&
					m.core.SelectedOperation.Name == "Start Pipeline" && m.core.ManualInput {
					m.core.CommitID = m.core.TextInput.Value()
					m.core.ManualCommitID = true
				}

				// If we're in the summary view with approval comment
				if m.core.CurrentView == constants.ViewSummary && m.core.SelectedApproval != nil {
					m.core.ApprovalComment = m.core.TextInput.Value()
				}

				// For AWS config view, the actual setting happens when Enter is pressed in HandleEnter
				return m, cmd
			}
		}
	case core.PipelineStatusMsg:
		m.core.Pipelines = msg.Pipelines
		m.core.Provider = msg.Provider
		m.core.CurrentView = constants.ViewPipelineStatus
		m.core.IsLoading = false
		view.UpdateTableForView(m.core)
		return m, nil
	}

	// If we're loading, make sure to keep the spinner spinning
	if m.core.IsLoading {
		return m, m.core.Spinner.Tick
	}

	return m, nil
}

// View renders the UI
func (m Model) View() string {
	return view.Render(m.core)
}
