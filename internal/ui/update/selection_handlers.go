package update

import (
	"github.com/HenryOwenz/cloudgate/internal/providers"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
	tea "github.com/charmbracelet/bubbletea"
)

// HandleServiceSelection handles the selection of a service
func HandleServiceSelection(m *model.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		serviceName := selected[0]

		// Get the AWS provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return WrapModel(m), func() tea.Msg {
				return model.ErrMsg{Err: err}
			}
		}

		// Find the selected service
		var selectedService providers.Service
		for _, service := range provider.Services() {
			if service.Name() == serviceName {
				selectedService = service
				break
			}
		}

		if selectedService != nil {
			newModel := m.Clone()
			newModel.SelectedService = &model.Service{
				Name:        selectedService.Name(),
				Description: selectedService.Description(),
			}
			newModel.CurrentView = constants.ViewSelectCategory
			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil
		}
	}
	return WrapModel(m), nil
}

// HandleCategorySelection handles the selection of a category
func HandleCategorySelection(m *model.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		categoryName := selected[0]

		// Get the AWS provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return WrapModel(m), func() tea.Msg {
				return model.ErrMsg{Err: err}
			}
		}

		// Find the selected service
		var selectedService providers.Service
		for _, service := range provider.Services() {
			if service.Name() == m.SelectedService.Name {
				selectedService = service
				break
			}
		}

		if selectedService == nil {
			return WrapModel(m), nil
		}

		// Find the selected category
		var selectedCategory providers.Category
		for _, category := range selectedService.Categories() {
			if category.Name() == categoryName {
				selectedCategory = category
				break
			}
		}

		if selectedCategory != nil {
			newModel := m.Clone()
			newModel.SelectedCategory = &model.Category{
				Name:        selectedCategory.Name(),
				Description: selectedCategory.Description(),
			}
			newModel.CurrentView = constants.ViewSelectOperation
			view.UpdateTableForView(newModel)
			return WrapModel(newModel), nil
		}
	}
	return WrapModel(m), nil
}

// HandleOperationSelection handles the selection of an operation
func HandleOperationSelection(m *model.Model) (tea.Model, tea.Cmd) {
	if selected := m.Table.SelectedRow(); len(selected) > 0 {
		operationName := selected[0]

		// Get the AWS provider from the registry
		provider, err := m.Registry.Get("AWS")
		if err != nil {
			return WrapModel(m), func() tea.Msg {
				return model.ErrMsg{Err: err}
			}
		}

		// Find the selected service
		var selectedService providers.Service
		for _, service := range provider.Services() {
			if service.Name() == m.SelectedService.Name {
				selectedService = service
				break
			}
		}

		if selectedService == nil {
			return WrapModel(m), nil
		}

		// Find the selected category
		var selectedCategory providers.Category
		for _, category := range selectedService.Categories() {
			if category.Name() == m.SelectedCategory.Name {
				selectedCategory = category
				break
			}
		}

		if selectedCategory == nil {
			return WrapModel(m), nil
		}

		// Find the selected operation
		var selectedOperation providers.Operation
		for _, operation := range selectedCategory.Operations() {
			if operation.Name() == operationName {
				selectedOperation = operation
				break
			}
		}

		if selectedOperation != nil {
			newModel := m.Clone()
			newModel.SelectedOperation = &model.Operation{
				Name:        selectedOperation.Name(),
				Description: selectedOperation.Description(),
			}

			// Handle different operations
			switch operationName {
			case "Pipeline Approvals":
				newModel.IsLoading = true
				newModel.LoadingMsg = "Fetching approvals..."
				return WrapModel(newModel), FetchApprovals(newModel)
			case "Pipeline Status":
				newModel.IsLoading = true
				newModel.LoadingMsg = "Fetching pipeline status..."
				return WrapModel(newModel), FetchPipelineStatus(newModel)
			case "Start Pipeline":
				newModel.CurrentView = constants.ViewPipelineStatus
				newModel.IsLoading = true
				newModel.LoadingMsg = "Fetching pipelines..."
				return WrapModel(newModel), FetchPipelineStatus(newModel)
			default:
				newModel.CurrentView = constants.ViewConfirmation
				view.UpdateTableForView(newModel)
				return WrapModel(newModel), nil
			}
		}
	}
	return WrapModel(m), nil
}
