package update

import (
	"github.com/HenryOwenz/cloudgate/internal/cloud"
	"github.com/HenryOwenz/cloudgate/internal/ui/constants"
	"github.com/HenryOwenz/cloudgate/internal/ui/model"
	"github.com/HenryOwenz/cloudgate/internal/ui/view"
	tea "github.com/charmbracelet/bubbletea"
)

// SelectService handles the selection of a service
func SelectService(m *model.Model) (tea.Model, tea.Cmd) {
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
		var selectedService cloud.Service
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

// SelectCategory handles the selection of a category
func SelectCategory(m *model.Model) (tea.Model, tea.Cmd) {
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
		var selectedService cloud.Service
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
		var selectedCategory cloud.Category
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

// SelectOperation handles the selection of an operation
func SelectOperation(m *model.Model) (tea.Model, tea.Cmd) {
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
		var selectedService cloud.Service
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
		var selectedCategory cloud.Category
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
		var selectedOperation cloud.Operation
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
				return HandlePipelineApprovals(newModel)
			case "Pipeline Status":
				return HandlePipelineStatus(newModel)
			case "Start Pipeline":
				return HandlePipelineStatus(newModel)
			case "Function Status":
				return HandleFunctionStatus(newModel)
			default:
				return WrapModel(newModel), nil
			}
		}
	}
	return WrapModel(m), nil
}
