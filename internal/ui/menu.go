package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("170")).Bold(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")).Bold(true)
	descStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).PaddingLeft(4)
	helpStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type MenuItem struct {
	title       string
	desc        string
	command     string
	subcommands []string
}

func (i MenuItem) Title() string       { return i.title }
func (i MenuItem) Description() string { return i.desc }
func (i MenuItem) FilterValue() string { return i.title }

type MenuModel struct {
	list     list.Model
	choice   string
	quitting bool
}

func NewMenuModel() MenuModel {
	items := []list.Item{
		MenuItem{
			title:   "🚀 Initialize Project",
			desc:    "Create a new microservice project with DDD structure",
			command: "init",
		},
		MenuItem{
			title:   "🤖 Generate Domain",
			desc:    "AI-powered domain generation from natural language",
			command: "gen domain",
		},
		MenuItem{
			title:   "📡 Generate Handler",
			desc:    "Create HTTP handlers with CRUD endpoints",
			command: "gen handler",
		},
		MenuItem{
			title:   "💾 Generate Repository",
			desc:    "Create database repository (Postgres/MySQL/Mongo)",
			command: "gen repository",
		},
		MenuItem{
			title:   "🛡️  Generate Middleware",
			desc:    "Create middleware (auth, ratelimit, logging, cors)",
			command: "gen middleware",
		},
		MenuItem{
			title:   "📊 Generate Migration",
			desc:    "Create database migration files",
			command: "gen migration",
		},
		MenuItem{
			title:   "🔌 Auto-Wire Dependencies",
			desc:    "Automatic dependency injection with AST discovery",
			command: "wire",
		},
		MenuItem{
			title:   "📐 Describe Architecture",
			desc:    "Generate architecture diagrams (Mermaid/ASCII)",
			command: "describe",
		},
		MenuItem{
			title:   "✨ Code Quality",
			desc:    "Lint, format, and validate code",
			command: "quality",
			subcommands: []string{"lint", "format", "validate"},
		},
		MenuItem{
			title:   "⚙️  Configuration",
			desc:    "Manage AI providers and settings",
			command: "config",
			subcommands: []string{"list", "set-provider", "check"},
		},
	}

	const defaultWidth = 80
	const listHeight = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "🎯 Anaphase CLI - AI-Powered Microservice Generator"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	l.Styles.HelpStyle = helpStyle

	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("/"),
				key.WithHelp("/", "filter"),
			),
		}
	}

	return MenuModel{list: l}
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 2 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(MenuItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%s", i.title)
	desc := descStyle.Render(i.desc)

	var output string
	if index == m.Index() {
		output = selectedItemStyle.Render("▶ " + str + "\n" + desc)
	} else {
		output = itemStyle.Render("  " + str + "\n" + desc)
	}

	fmt.Fprint(w, output)
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(MenuItem)
			if ok {
				m.choice = i.command
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m MenuModel) View() string {
	if m.choice != "" {
		return ""
	}
	if m.quitting {
		return "\n  👋 Goodbye!\n\n"
	}

	return "\n" + m.list.View() + "\n\n" + helpStyle.Render("  Press '/' to filter • ↑↓ to navigate • Enter to select • Ctrl+C to quit") + "\n"
}

func (m MenuModel) GetChoice() string {
	return m.choice
}

// FormatCommand formats the selected command for execution
func FormatCommand(choice string) []string {
	parts := strings.Split(choice, " ")
	return parts
}
