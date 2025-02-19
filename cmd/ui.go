package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/HenryOwenz/ciselect/internal/aws"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const (
	stepFetchingApprovals = iota
	stepSelectingApproval
	stepConfirmingAction
	stepExecutingAction
)

type Styles struct {
	Title       lipgloss.Style
	Selected    lipgloss.Style
	Unselected  lipgloss.Style
	Instruction lipgloss.Style
	Error       lipgloss.Style
}

func defaultStyles() Styles {
	return Styles{
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true).
			Padding(1, 0),
		Selected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true),
		Unselected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")),
		Instruction: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFFF")).
			Italic(true),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true),
	}
}

// Model represents the application state
type Model struct {
	profiles         []string
	regions          []string
	approvals        []aws.ApprovalAction
	selectedProfile  string
	selectedRegion   string
	cursor           int
	step             int // 0: profile, 1: region, 2: approvals, 3: action
	err              error
	awsService       *aws.Service
	selectedApproval *aws.ApprovalAction
	summary          string
	action           string // "approve" or "reject"
	styles           Styles
	manualInput      bool
	inputBuffer      string
}

func initialModel(profile, region string) Model {
	m := Model{
		profiles:    getAWSProfiles(),
		regions:     []string{"us-east-1", "us-east-2", "us-west-1", "us-west-2", "eu-west-1", "eu-west-2", "eu-central-1", "ap-southeast-1", "ap-southeast-2", "ap-northeast-1"},
		step:        0,
		cursor:      0,
		styles:      defaultStyles(),
		manualInput: false,
		inputBuffer: "",
	}

	// If profile/region provided via flags, skip to approval fetching
	if profile != "" && region != "" {
		m.selectedProfile = profile
		m.selectedRegion = region
		m.step = stepFetchingApprovals
	}

	return m
}

// Update handles messages and updates the model accordingly
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if !m.manualInput && m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if !m.manualInput {
				switch m.step {
				case 0: // Profile selection
					if m.cursor < len(m.profiles)-1 {
						m.cursor++
					}
				case 1: // Region selection
					if m.cursor < len(m.regions)-1 {
						m.cursor++
					}
				case 2: // Approval selection
					if m.cursor < len(m.approvals)-1 {
						m.cursor++
					}
				case 3, 5: // Action/Confirmation selection
					if m.cursor < 1 {
						m.cursor++
					}
				}
			}

		case "tab":
			if m.step <= 1 { // Only for profile/region selection
				m.manualInput = !m.manualInput
				m.inputBuffer = ""
				m.cursor = 0
			}

		case "enter":
			if m.manualInput {
				input := strings.TrimSpace(m.inputBuffer)
				if input != "" {
					if m.step == 0 {
						m.selectedProfile = input
						m.step++
					} else if m.step == 1 {
						m.selectedRegion = input
						m.step = stepFetchingApprovals
						return m.initAWS()
					}
					m.manualInput = false
					m.inputBuffer = ""
				}
				return m, nil
			}
			switch m.step {
			case 0: // Profile selected
				m.selectedProfile = m.profiles[m.cursor]
				m.step++
				m.cursor = 0

			case 1: // Region selected
				m.selectedRegion = m.regions[m.cursor]
				// Initialize AWS service and fetch approvals
				ctx := context.Background()
				service, err := aws.NewService(ctx, m.selectedProfile, m.selectedRegion)
				if err != nil {
					m.err = err
					return m, nil
				}
				m.awsService = service

				approvals, err := service.ListPendingApprovals(ctx)
				if err != nil {
					m.err = err
					return m, nil
				}
				m.approvals = approvals
				m.step++
				m.cursor = 0

			case 2: // Approval selected
				if len(m.approvals) > 0 {
					m.selectedApproval = &m.approvals[m.cursor]
					m.step++
					m.cursor = 0
				}

			case 3: // Action selection
				switch m.cursor {
				case 0: // Approve
					m.action = "approve"
					m.step++
				case 1: // Reject
					m.action = "reject"
					m.step++
				case 2: // Cancel
					return m, tea.Quit
				}

			case 4: // Summary input
				if m.summary != "" {
					m.step++
				}

			case 5: // Confirmation
				if m.cursor == 0 { // Yes
					ctx := context.Background()
					var err error
					if m.action == "approve" {
						err = m.awsService.ApproveAction(ctx,
							m.selectedApproval.PipelineName,
							m.selectedApproval.StageName,
							m.selectedApproval.ActionName,
							m.selectedApproval.Token,
							m.summary)
					} else {
						err = m.awsService.RejectAction(ctx,
							m.selectedApproval.PipelineName,
							m.selectedApproval.StageName,
							m.selectedApproval.ActionName,
							m.selectedApproval.Token,
							m.summary)
					}
					if err != nil {
						m.err = err
						return m, nil
					}
					return m, tea.Quit
				} else { // No
					return m, tea.Quit
				}
			}

		case "backspace":
			if m.manualInput && len(m.inputBuffer) > 0 {
				m.inputBuffer = m.inputBuffer[:len(m.inputBuffer)-1]
			} else if m.step == 4 && len(m.summary) > 0 {
				m.summary = m.summary[:len(m.summary)-1]
			}

		default:
			if m.manualInput {
				m.inputBuffer += msg.String()
			} else if m.step == 4 {
				m.summary += msg.String()
			}
		}
	}
	return m, nil
}

// View renders the UI
func (m Model) View() string {
	var s strings.Builder

	if m.err != nil {
		return m.styles.Error.Render(fmt.Sprintf("Error: %v", m.err))
	}

	switch m.step {
	case 0, 1: // Profile or Region selection
		title := "Select AWS Profile"
		if m.step == 1 {
			title = "Select AWS Region"
			s.WriteString(m.styles.Instruction.Render(fmt.Sprintf("Profile: %s\n", m.selectedProfile)))
		}
		s.WriteString(m.styles.Title.Render(title))

		if m.manualInput {
			s.WriteString("\nEnter name: ")
			s.WriteString(m.inputBuffer)
			s.WriteString("_")
		} else {
			items := m.profiles
			if m.step == 1 {
				items = m.regions
			}
			for i, item := range items {
				if m.cursor == i {
					s.WriteString("\n> " + m.styles.Selected.Render(item))
				} else {
					s.WriteString("\n  " + item)
				}
			}
		}

	case 2: // Approval selection
		s.WriteString(m.styles.Title.Render("Select Approval"))
		s.WriteString("\n")
		s.WriteString(m.styles.Instruction.Render(fmt.Sprintf("Profile: %s | Region: %s",
			m.selectedProfile, m.selectedRegion)))

		if len(m.approvals) == 0 {
			s.WriteString("\n\n")
			s.WriteString(m.styles.Instruction.Render("No pending approvals found"))
			s.WriteString("\n\nPress q to quit")
			return s.String()
		}

		for i, approval := range m.approvals {
			text := fmt.Sprintf("%s → %s → %s",
				approval.PipelineName,
				approval.StageName,
				approval.ActionName)

			if m.cursor == i {
				s.WriteString("\n> " + m.styles.Selected.Render(text))
			} else {
				s.WriteString("\n  " + text)
			}
		}

	case 3: // Action selection
		s.WriteString(m.styles.Title.Render("Choose Action"))
		s.WriteString("\n")
		s.WriteString(m.styles.Instruction.Render(fmt.Sprintf("Pipeline: %s\nStage: %s\nAction: %s",
			m.selectedApproval.PipelineName,
			m.selectedApproval.StageName,
			m.selectedApproval.ActionName)))

		options := []string{"Approve", "Reject", "Cancel"}
		for i, option := range options {
			if m.cursor == i {
				s.WriteString("\n> " + m.styles.Selected.Render(option))
			} else {
				s.WriteString("\n  " + option)
			}
		}

	case 4: // Summary input
		s.WriteString(m.styles.Title.Render("Enter Summary"))
		s.WriteString("\n")
		s.WriteString(m.styles.Instruction.Render(fmt.Sprintf("Action: %s", m.action)))
		s.WriteString("\n\nSummary: ")
		s.WriteString(m.summary)
		s.WriteString("_")

	case 5: // Confirmation
		s.WriteString(m.styles.Title.Render("Confirm Action"))
		s.WriteString("\n")
		s.WriteString(m.styles.Instruction.Render(fmt.Sprintf(`Pipeline: %s
Stage: %s
Action: %s
Operation: %s
Summary: %s

Are you sure?`,
			m.selectedApproval.PipelineName,
			m.selectedApproval.StageName,
			m.selectedApproval.ActionName,
			m.action,
			m.summary)))

		options := []string{"Yes", "No"}
		for i, option := range options {
			if m.cursor == i {
				s.WriteString("\n> " + m.styles.Selected.Render(option))
			} else {
				s.WriteString("\n  " + option)
			}
		}
	}

	// Help text
	s.WriteString("\n\n")
	if m.step <= 1 {
		s.WriteString(m.styles.Instruction.Render("↑/↓: Navigate • Enter: Select • Tab: Toggle Input Mode • q: Quit"))
	} else {
		s.WriteString(m.styles.Instruction.Render("↑/↓: Navigate • Enter: Select • q: Quit"))
	}

	return s.String()
}

func getAWSProfiles() []string {
	// Read profiles from AWS credentials file
	home, err := os.UserHomeDir()
	if err != nil {
		return []string{}
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

var interactiveCmd = &cobra.Command{
	Use:     "interactive",
	Aliases: []string{"i"},
	Short:   "Interactive mode for managing AWS CodePipeline approvals",
	Long: `Interactive mode provides a user-friendly interface for managing AWS CodePipeline approvals.
It guides you through selecting your AWS profile, region, and approval actions with a beautiful interface.

You can either:
1. Select from a list of available profiles and regions
2. Type your own profile and region names`,
	Run: func(cmd *cobra.Command, args []string) {
		profile, _ := cmd.Flags().GetString("profile")
		region, _ := cmd.Flags().GetString("region")

		p := tea.NewProgram(initialModel(profile, region))
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running program: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	interactiveCmd.Flags().StringP("profile", "p", "", "AWS profile to use (optional)")
	interactiveCmd.Flags().StringP("region", "r", "", "AWS region to use (optional)")
	rootCmd.AddCommand(interactiveCmd)
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) initAWS() (Model, tea.Cmd) {
	ctx := context.Background()
	service, err := aws.NewService(ctx, m.selectedProfile, m.selectedRegion)
	if err != nil {
		m.err = err
		return m, nil
	}
	m.awsService = service

	approvals, err := service.ListPendingApprovals(ctx)
	if err != nil {
		m.err = err
		return m, nil
	}
	m.approvals = approvals
	return m, nil
}
