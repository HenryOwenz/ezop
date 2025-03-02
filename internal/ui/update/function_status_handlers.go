package update

import (
	"context"

	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	tea "github.com/charmbracelet/bubbletea"
)

// HandleFunctionStatusOperation handles the function status operation
func HandleFunctionStatusOperation(m *model.Model) (tea.Model, tea.Cmd) {
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
