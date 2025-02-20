package model

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/HenryOwenz/ezop/internal/domain"
	"github.com/HenryOwenz/ezop/internal/providers/aws"
	"github.com/HenryOwenz/ezop/internal/ui/styles"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// Step represents the current step in the UI workflow
type Step int

const (
	StepSelectProvider Step = iota
	StepProviderConfig
	StepSelectService
	StepSelectCategory // New step for selecting between workflows and operations
	StepServiceOperation
	StepSelectingApproval
	StepConfirmingAction
	StepSummaryInput
	StepExecutingAction
)

// Model represents the application state
type Model struct {
	Profiles    []string
	Regions     []string
	Approvals   []aws.ApprovalAction
	Cursor      int
	Step        Step
	Error       error
	Styles      styles.Styles
	ManualInput bool
	InputBuffer string
	IsLoading   bool          // Indicates if an API request is in progress
	LoadingMsg  string        // Custom loading message for different operations
	Spinner     spinner.Model // Spinner component for loading states
	Table       table.Model   // Table component for list views

	// Provider selection
	SelectedProvider *domain.Provider
	Providers        []domain.Provider

	// Service selection
	Services          []domain.Service
	SelectedService   *domain.Service
	Categories        []domain.Category // New field for workflow/operations categories
	SelectedCategory  *domain.Category  // New field for selected category
	Operations        []domain.Operation
	SelectedOperation *domain.Operation

	// AWS specific
	AWSProfile       string
	AWSRegion        string
	AWSProvider      *aws.Provider
	SelectedApproval *aws.ApprovalAction
	Summary          string
	Action           string // "approve" or "reject"
}

// NewModel creates a new Model with initial state
func NewModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.DefaultStyles().Loading

	// Initialize table with default styles
	t := table.New(
		table.WithColumns([]table.Column{
			{Title: "Provider", Width: 30},
			{Title: "Description", Width: 50},
		}),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Set table styles
	tableStyles := table.DefaultStyles()
	tableStyles.Header = tableStyles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#00FF00")).
		BorderBottom(true).
		Bold(false)
	tableStyles.Selected = tableStyles.Selected.
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true)
	t.SetStyles(tableStyles)

	// Initialize with provider rows
	rows := []table.Row{}
	for _, provider := range domain.DefaultProviders {
		var name, desc string
		if provider.Available {
			name = provider.Name
			desc = provider.Description
		} else {
			name = provider.Name + " (Coming Soon)"
			desc = provider.Description
		}
		rows = append(rows, table.Row{name, desc})
	}
	t.SetRows(rows)

	return Model{
		Profiles:    getAWSProfiles(),
		Regions:     []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2", "eu-west-1", "eu-west-2", "eu-central-1", "ap-southeast-1", "ap-southeast-2", "ap-northeast-1"},
		Step:        StepSelectProvider,
		Cursor:      0,
		Styles:      styles.DefaultStyles(),
		ManualInput: false,
		InputBuffer: "",
		Spinner:     s,
		Table:       t,
		Providers:   domain.DefaultProviders,
		Categories: []domain.Category{
			{
				ID:          "workflows",
				Name:        "Workflows",
				Description: "Curated workflows for common tasks",
				Available:   true,
			},
			{
				ID:          "operations",
				Name:        "Operations",
				Description: "Direct AWS service operations",
				Available:   false, // Not implemented yet
			},
		},
	}
}

// getAWSProfiles returns a list of AWS profiles from the AWS credentials file
func getAWSProfiles() []string {
	// Read profiles from AWS credentials file
	home, err := os.UserHomeDir()
	if err != nil {
		return []string{"default"}
	}

	// Try both config and credentials files
	configFiles := []string{
		filepath.Join(home, ".aws", "config"),
		filepath.Join(home, ".aws", "credentials"),
	}

	var profiles []string
	for _, file := range configFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		// Parse profiles using regex
		re := regexp.MustCompile(`\[(.*?)\]`)
		matches := re.FindAllStringSubmatch(string(content), -1)
		for _, match := range matches {
			profile := strings.TrimSpace(match[1])
			// Remove "profile " prefix if present (used in config file)
			profile = strings.TrimPrefix(profile, "profile ")
			if profile != "" && !contains(profiles, profile) {
				profiles = append(profiles, profile)
			}
		}
	}

	if len(profiles) == 0 {
		return []string{"default"}
	}

	sort.Strings(profiles)
	return profiles
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// NavigateBack moves to the previous step in the workflow
func (m *Model) NavigateBack() {
	switch m.Step {
	case StepProviderConfig:
		if m.AWSRegion != "" {
			// If we have a region, clear it and stay in provider config to select region
			m.AWSRegion = ""
			m.Cursor = 0
			m.ManualInput = false
			m.InputBuffer = ""
		} else if m.AWSProfile != "" {
			// If we have a profile but no region, clear profile to select profile
			m.AWSProfile = ""
			m.Cursor = 0
			m.ManualInput = false
			m.InputBuffer = ""
		} else {
			// If we have neither, go back to provider selection
			m.Step = StepSelectProvider
			m.SelectedProvider = nil
			m.Cursor = 0
			m.ManualInput = false
			m.InputBuffer = ""
		}
	case StepSelectService:
		m.Step = StepProviderConfig
		m.AWSRegion = "" // Clear region but keep profile
		m.Services = nil
		m.AWSProvider = nil
		m.Cursor = 0
		m.ManualInput = false
		m.InputBuffer = ""
	case StepSelectCategory:
		m.Step = StepSelectService
		m.SelectedCategory = nil
		m.Cursor = 0
	case StepServiceOperation:
		m.Step = StepSelectCategory
		m.SelectedOperation = nil
		m.Operations = nil
		m.Cursor = 0
	case StepSelectingApproval:
		m.Step = StepServiceOperation
		m.SelectedOperation = nil
		m.Approvals = nil
		m.Cursor = 0
	case StepConfirmingAction:
		m.Step = StepSelectingApproval
		m.SelectedApproval = nil
		m.Cursor = 0
	case StepSummaryInput:
		m.Step = StepConfirmingAction
		m.Summary = ""
		m.Action = ""
		m.Cursor = 0
	case StepExecutingAction:
		m.Step = StepSummaryInput
		m.Cursor = 0
	}
}

// UpdateTableForStep updates the table based on the current step
func (m *Model) UpdateTableForStep() {
	var columns []table.Column
	var rows []table.Row

	switch m.Step {
	case StepSelectProvider:
		columns = []table.Column{
			{Title: "Provider", Width: 30},
			{Title: "Description", Width: 50},
		}
		for _, provider := range m.Providers {
			if provider.Available {
				rows = append(rows, table.Row{provider.Name, provider.Description})
			} else {
				rows = append(rows, table.Row{
					provider.Name + " (Coming Soon)",
					provider.Description,
				})
			}
		}

	case StepProviderConfig:
		if m.AWSProfile == "" {
			columns = []table.Column{
				{Title: "Profile", Width: 30},
			}
			for _, profile := range m.Profiles {
				rows = append(rows, table.Row{profile})
			}
		} else {
			columns = []table.Column{
				{Title: "Region", Width: 30},
			}
			for _, region := range m.Regions {
				rows = append(rows, table.Row{region})
			}
		}

	case StepSelectService:
		columns = []table.Column{
			{Title: "Service", Width: 30},
			{Title: "Description", Width: 50},
		}
		for _, service := range m.Services {
			if service.Available {
				rows = append(rows, table.Row{service.Name, service.Description})
			} else {
				rows = append(rows, table.Row{
					service.Name + " (Coming Soon)",
					service.Description,
				})
			}
		}

	case StepSelectCategory:
		columns = []table.Column{
			{Title: "Category", Width: 30},
			{Title: "Description", Width: 50},
		}
		for _, category := range m.Categories {
			if category.Available {
				rows = append(rows, table.Row{category.Name, category.Description})
			} else {
				rows = append(rows, table.Row{
					category.Name + " (Coming Soon)",
					category.Description,
				})
			}
		}

	case StepServiceOperation:
		columns = []table.Column{
			{Title: "Operation", Width: 30},
			{Title: "Description", Width: 50},
		}
		for _, operation := range m.Operations {
			rows = append(rows, table.Row{operation.Name, operation.Description})
		}

	case StepSelectingApproval:
		columns = []table.Column{
			{Title: "Pipeline", Width: 30},
			{Title: "Stage", Width: 20},
			{Title: "Action", Width: 20},
		}
		for _, approval := range m.Approvals {
			rows = append(rows, table.Row{
				approval.PipelineName,
				approval.StageName,
				approval.ActionName,
			})
		}

	case StepConfirmingAction:
		columns = []table.Column{
			{Title: "Action", Width: 30},
			{Title: "Description", Width: 50},
		}
		rows = []table.Row{
			{"Approve", "Approve the pipeline stage"},
			{"Reject", "Reject the pipeline stage"},
			{"Cancel", "Cancel and return to main menu"},
		}
		m.Table = table.New(
			table.WithColumns(columns),
			table.WithRows(rows),
			table.WithFocused(true),
			table.WithHeight(10),
		)
		m.Table.SetCursor(m.Cursor)

	case StepSummaryInput:
		columns = []table.Column{
			{Title: "Action", Width: 30},
			{Title: "Description", Width: 50},
		}
		rows = []table.Row{
			{"Confirm", fmt.Sprintf("Proceed with %s action", m.Action)},
			{"Cancel", "Cancel and return to main menu"},
		}

	case StepExecutingAction:
		columns = []table.Column{
			{Title: "Action", Width: 30},
			{Title: "Description", Width: 50},
		}
		rows = []table.Row{
			{"Execute", fmt.Sprintf("Execute %s action with summary", m.Action)},
			{"Cancel", "Cancel and return to main menu"},
		}
	}

	// Create new table with current state
	m.Table = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Set table styles
	tableStyles := table.DefaultStyles()
	tableStyles.Header = tableStyles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#00FF00")).
		BorderBottom(true).
		Bold(false)
	tableStyles.Selected = tableStyles.Selected.
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true)
	m.Table.SetStyles(tableStyles)

	// Ensure cursor state is valid
	if m.Cursor >= len(rows) {
		m.Cursor = 0
	}
	m.Table.SetCursor(m.Cursor)
}

// GetSelectedRow returns the currently selected row based on the step
func (m *Model) GetSelectedRow() interface{} {
	if len(m.Table.Rows()) == 0 {
		return nil
	}

	switch m.Step {
	case StepSelectProvider:
		if m.Cursor < len(m.Providers) {
			return m.Providers[m.Cursor]
		}
	case StepProviderConfig:
		if m.AWSProfile == "" {
			if m.Cursor < len(m.Profiles) {
				return m.Profiles[m.Cursor]
			}
		} else {
			if m.Cursor < len(m.Regions) {
				return m.Regions[m.Cursor]
			}
		}
	case StepSelectService:
		if m.Cursor < len(m.Services) {
			return m.Services[m.Cursor]
		}
	case StepSelectCategory:
		if m.Cursor < len(m.Categories) {
			return m.Categories[m.Cursor]
		}
	case StepServiceOperation:
		if m.Cursor < len(m.Operations) {
			return m.Operations[m.Cursor]
		}
	case StepSelectingApproval:
		if m.Cursor < len(m.Approvals) {
			return m.Approvals[m.Cursor]
		}
	}
	return nil
}
