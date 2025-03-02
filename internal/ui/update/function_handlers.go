package update

import (
	"context"
	"fmt"

	"github.com/HenryOwenz/cloudgate/internal/providers"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
	tea "github.com/charmbracelet/bubbletea"
)

// HandleFunctionStatus handles the function status operation
func HandleFunctionStatus(m *model.Model) (tea.Model, tea.Cmd) {
	newModel := m.Clone()
	newModel.IsLoading = true
	newModel.LoadingMsg = constants.MsgLoadingFunctions

	return WrapModel(newModel), func() tea.Msg {
		// Get the provider
		provider, err := m.Registry.Get(m.ProviderState.ProviderName)
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Get the FunctionStatusOperation from the provider
		functionOperation, err := provider.GetFunctionStatusOperation()
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		// Get function status using the operation
		ctx := context.Background()
		functions, err := functionOperation.GetFunctionStatus(ctx)
		if err != nil {
			return model.ErrMsg{Err: err}
		}

		return model.FunctionStatusMsg{
			Functions: functions,
			Provider:  provider,
		}
	}
}

// HandleFunctionSelection handles the selection of a function
func HandleFunctionSelection(m *model.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		functionName := selected[0]

		// Find the selected function
		var selectedFunction *providers.FunctionStatus
		for _, function := range m.Functions {
			if function.Name == functionName {
				selectedFunction = &function
				break
			}
		}

		if selectedFunction == nil {
			return WrapModel(m), func() tea.Msg {
				return model.ErrMsg{Err: fmt.Errorf(constants.MsgErrorNoFunction)}
			}
		}

		// Update the model
		newModel := m.Clone()
		newModel.SetSelectedFunction(selectedFunction)
		newModel.CurrentView = constants.ViewFunctionDetails
		view.UpdateTableForView(newModel)
		return WrapModel(newModel), nil
	}
	return WrapModel(m), nil
}
