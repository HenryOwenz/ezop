package ui

import (
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/update"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Model is the main UI model that implements the tea.Model interface
type Model struct {
	core *model.Model
}

// New creates a new UI model
func New() Model {
	m := Model{
		core: model.New(),
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
		newModel := m.Clone()
		newModel.core.Width = msg.Width
		newModel.core.Height = msg.Height
		view.UpdateTableForView(newModel.core)
		return newModel, nil
	case model.ErrMsg:
		newModel := m.Clone()
		newModel.core.Err = msg.Err
		newModel.core.IsLoading = false
		return newModel, nil
	case model.ApprovalsMsg:
		newModel := m.Clone()
		newModel.core.Approvals = msg.Approvals
		newModel.core.Provider = msg.Provider
		newModel.core.CurrentView = constants.ViewApprovals
		newModel.core.IsLoading = false
		view.UpdateTableForView(newModel.core)
		return newModel, nil
	case model.ApprovalResultMsg:
		newModel := m.Clone()
		newModel.core.IsLoading = false // Ensure loading is turned off
		update.HandleApprovalResult(newModel.core, msg.Err)
		view.UpdateTableForView(newModel.core)
		return newModel, nil
	case model.PipelineExecutionMsg:
		newModel := m.Clone()
		newModel.core.IsLoading = false // Ensure loading is turned off
		update.HandlePipelineExecution(newModel.core, msg.Err)
		view.UpdateTableForView(newModel.core)
		return newModel, nil
	case spinner.TickMsg:
		newModel := m.Clone()
		var cmd tea.Cmd
		newModel.core.Spinner, cmd = newModel.core.Spinner.Update(msg)
		if newModel.core.IsLoading {
			cmds = append(cmds, cmd)
		}
		return newModel, tea.Batch(cmds...)
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
			modelWrapper, cmd := update.HandleEnter(m.core)
			if wrapper, ok := modelWrapper.(update.ModelWrapper); ok {
				// Since ModelWrapper embeds *model.Model, we can create a new Model with it
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
				newModel := m.Clone()
				var cmd tea.Cmd
				newModel.core.TextInput, cmd = newModel.core.TextInput.Update(msg)
				return newModel, cmd
			}

			// Handle back navigation
			if m.core.ManualInput {
				newModel := m.Clone()
				newModel.core.ManualInput = false
				newModel.core.ResetTextInput()
				view.UpdateTableForView(newModel.core)
				return newModel, nil
			}
			newCore := update.NavigateBack(m.core)
			view.UpdateTableForView(newCore)
			return Model{core: newCore}, nil
		case constants.KeyUp, constants.KeyAltUp:
			newModel := m.Clone()
			newModel.core.Table.MoveUp(1)
			return newModel, nil
		case constants.KeyDown, constants.KeyAltDown:
			newModel := m.Clone()
			newModel.core.Table.MoveDown(1)
			return newModel, nil
		case constants.KeyTab:
			// Tab is now only used for other views, not AWS config
			if m.core.CurrentView == constants.ViewSummary {
				newModel := m.Clone()
				if newModel.core.ManualInput {
					newModel.core.ManualInput = false
					newModel.core.ResetTextInput()
				} else {
					newModel.core.ManualInput = true
					newModel.core.TextInput.Focus()
				}
				return newModel, nil
			}
			return m, nil
		default:
			if m.core.ManualInput {
				newModel := m.Clone()
				var cmd tea.Cmd
				newModel.core.TextInput, cmd = newModel.core.TextInput.Update(msg)

				// If we're in the summary view with manual commit ID
				if newModel.core.CurrentView == constants.ViewSummary && newModel.core.SelectedOperation != nil &&
					newModel.core.SelectedOperation.Name == "Start Pipeline" && newModel.core.ManualInput {
					newModel.core.CommitID = newModel.core.TextInput.Value()
					newModel.core.ManualCommitID = true
				}

				// If we're in the summary view with approval comment
				if newModel.core.CurrentView == constants.ViewSummary && newModel.core.SelectedApproval != nil {
					newModel.core.ApprovalComment = newModel.core.TextInput.Value()
				}

				// For AWS config view, the actual setting happens when Enter is pressed in HandleEnter
				return newModel, cmd
			}
		}
	case model.PipelineStatusMsg:
		newModel := m.Clone()
		newModel.core.Pipelines = msg.Pipelines
		newModel.core.Provider = msg.Provider
		newModel.core.CurrentView = constants.ViewPipelineStatus
		newModel.core.IsLoading = false
		view.UpdateTableForView(newModel.core)
		return newModel, nil
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

// Clone creates a deep copy of the model
func (m Model) Clone() Model {
	return Model{
		core: m.core.Clone(),
	}
}
