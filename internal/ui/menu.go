package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color scheme inspired by Claude Code
	primaryColor   = lipgloss.Color("99")  // Purple
	secondaryColor = lipgloss.Color("86")  // Cyan
	accentColor    = lipgloss.Color("212") // Pink
	mutedColor     = lipgloss.Color("241") // Gray
	successColor   = lipgloss.Color("42")  // Green
	warningColor   = lipgloss.Color("214") // Orange

	// Styles
	titleStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true).
			Padding(0, 1)

	categoryStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true).
			Padding(0, 2)

	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			PaddingLeft(4)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				PaddingLeft(2)

	descStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			PaddingLeft(6)

	shortcutStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Padding(1, 2)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(0, 1)

	searchStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(warningColor)

	mutedTextStyle = lipgloss.NewStyle().
			Foreground(mutedColor)
)

type MenuItem struct {
	title        string
	desc         string
	command      string
	subcommands  []string
	needsInput   bool
	inputPrompts []string
	shortcut     string
	requiresAI   bool
	category     string
}

type MenuCategory struct {
	name  string
	items []MenuItem
	icon  string
}

type MenuModel struct {
	categories       []MenuCategory
	selectedCategory int
	selectedItem     int
	filtering        bool
	searchInput      textinput.Model
	filteredItems    []MenuItem
	choice           string
	selectedMenuItem *MenuItem
	quitting         bool
	width            int
	height           int
}

func NewMenuModel() MenuModel {
	// Initialize search input
	ti := textinput.New()
	ti.Placeholder = "Type to search..."
	ti.CharLimit = 50
	ti.Width = 50

	categories := []MenuCategory{
		{
			name: "Generation",
			icon: "📋",
			items: []MenuItem{
				{
					title:        "Generate Domain",
					desc:         "AI-powered domain generation from natural language",
					command:      "gen domain",
					needsInput:   true,
					inputPrompts: []string{"Domain description"},
					shortcut:     "1",
					requiresAI:   true,
					category:     "Generation",
				},
				{
					title:        "Generate Handler",
					desc:         "HTTP handlers with CRUD endpoints",
					command:      "gen handler",
					needsInput:   true,
					inputPrompts: []string{"Handler name"},
					shortcut:     "2",
					category:     "Generation",
				},
				{
					title:        "Generate Repository",
					desc:         "Database repository implementation",
					command:      "gen repository",
					needsInput:   true,
					inputPrompts: []string{"Repository name"},
					shortcut:     "3",
					category:     "Generation",
				},
				{
					title:        "Generate Middleware",
					desc:         "HTTP middleware (auth, ratelimit, logging, cors)",
					command:      "gen middleware",
					needsInput:   true,
					inputPrompts: []string{"Middleware type"},
					shortcut:     "4",
					category:     "Generation",
				},
				{
					title:        "Generate Migration",
					desc:         "Database migration files",
					command:      "gen migration",
					needsInput:   true,
					inputPrompts: []string{"Migration name"},
					shortcut:     "5",
					category:     "Generation",
				},
			},
		},
		{
			name: "Tools",
			icon: "🔧",
			items: []MenuItem{
				{
					title:      "Auto-Wire Dependencies",
					desc:       "Automatic dependency injection with AST discovery",
					command:    "wire",
					needsInput: false,
					shortcut:   "6",
					category:   "Tools",
				},
				{
					title:      "Describe Architecture",
					desc:       "Generate architecture diagrams (Mermaid/ASCII)",
					command:    "describe",
					needsInput: false,
					shortcut:   "7",
					category:   "Tools",
				},
				{
					title:       "Code Quality",
					desc:        "Lint, format, and validate code",
					command:     "quality",
					subcommands: []string{"lint", "format", "validate"},
					needsInput:  false,
					shortcut:    "8",
					category:    "Tools",
				},
			},
		},
		{
			name: "Configuration",
			icon: "⚙️",
			items: []MenuItem{
				{
					title:      "Initialize Project",
					desc:       "Create new microservice project",
					command:    "init",
					needsInput: true,
					inputPrompts: []string{
						"Project name",
						"Database type (postgres/mysql/sqlite)",
					},
					shortcut: "i",
					category: "Configuration",
				},
				{
					title:       "Manage Config",
					desc:        "AI providers and settings",
					command:     "config",
					subcommands: []string{"list", "set-provider", "check"},
					needsInput:  false,
					shortcut:    "c",
					category:    "Configuration",
				},
			},
		},
	}

	return MenuModel{
		categories:       categories,
		selectedCategory: 0,
		selectedItem:     0,
		searchInput:      ti,
		filtering:        false,
		width:            80,
		height:           24,
	}
}

type keyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Tab      key.Binding
	Enter    key.Binding
	Search   key.Binding
	Quit     key.Binding
	Help     key.Binding
	Escape   key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "prev category"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "next category"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next category"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Search: key.NewBinding(
		key.WithKeys("/", "ctrl+k"),
		key.WithHelp("/ or ctrl+k", "search"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel search"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		// Handle search input mode
		if m.filtering {
			switch msg.String() {
			case "esc":
				m.filtering = false
				m.searchInput.SetValue("")
				m.filteredItems = nil
				return m, nil
			case "enter":
				if len(m.filteredItems) > 0 && m.selectedItem < len(m.filteredItems) {
					item := m.filteredItems[m.selectedItem]
					m.choice = item.command
					m.selectedMenuItem = &item
					return m, tea.Quit
				}
				return m, nil
			case "up", "k":
				if m.selectedItem > 0 {
					m.selectedItem--
				}
				return m, nil
			case "down", "j":
				if m.selectedItem < len(m.filteredItems)-1 {
					m.selectedItem++
				}
				return m, nil
			default:
				var cmd tea.Cmd
				m.searchInput, cmd = m.searchInput.Update(msg)
				m.filterItems()
				m.selectedItem = 0
				return m, cmd
			}
		}

		// Normal navigation mode
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "/", "ctrl+k":
			m.filtering = true
			m.searchInput.Focus()
			return m, textinput.Blink

		case "left", "h":
			if m.selectedCategory > 0 {
				m.selectedCategory--
				m.selectedItem = 0
			}
			return m, nil

		case "right", "l", "tab":
			if m.selectedCategory < len(m.categories)-1 {
				m.selectedCategory++
				m.selectedItem = 0
			}
			return m, nil

		case "up", "k":
			if m.selectedItem > 0 {
				m.selectedItem--
			}
			return m, nil

		case "down", "j":
			currentCategory := m.categories[m.selectedCategory]
			if m.selectedItem < len(currentCategory.items)-1 {
				m.selectedItem++
			}
			return m, nil

		case "enter":
			currentCategory := m.categories[m.selectedCategory]
			if m.selectedItem < len(currentCategory.items) {
				item := currentCategory.items[m.selectedItem]
				m.choice = item.command
				m.selectedMenuItem = &item
				return m, tea.Quit
			}
			return m, nil

		default:
			// Handle number shortcuts
			for catIdx, category := range m.categories {
				for itemIdx, item := range category.items {
					if msg.String() == item.shortcut {
						m.selectedCategory = catIdx
						m.selectedItem = itemIdx
						m.choice = item.command
						m.selectedMenuItem = &item
						return m, tea.Quit
					}
				}
			}
		}
	}

	return m, nil
}

func (m *MenuModel) filterItems() {
	query := strings.ToLower(m.searchInput.Value())
	if query == "" {
		m.filteredItems = nil
		return
	}

	m.filteredItems = []MenuItem{}
	for _, category := range m.categories {
		for _, item := range category.items {
			titleMatch := strings.Contains(strings.ToLower(item.title), query)
			descMatch := strings.Contains(strings.ToLower(item.desc), query)
			commandMatch := strings.Contains(strings.ToLower(item.command), query)

			if titleMatch || descMatch || commandMatch {
				m.filteredItems = append(m.filteredItems, item)
			}
		}
	}
}

func (m MenuModel) View() string {
	if m.choice != "" {
		return ""
	}

	if m.quitting {
		return helpStyle.Render("\n  👋 Goodbye!\n\n")
	}

	// Header
	header := titleStyle.Render("⚡ Anaphase CLI") +
		headerStyle.Render(" - DDD Microservice Generator") +
		"\n" +
		helpStyle.Render("  💡 Commands marked with [AI] require API key setup\n")

	var content string

	// Search mode
	if m.filtering {
		content = m.renderSearchView()
	} else {
		content = m.renderCategoryView()
	}

	// Footer with help
	footer := m.renderFooter()

	return "\n" + header + "\n" + content + "\n" + footer + "\n"
}

func (m MenuModel) renderCategoryView() string {
	var b strings.Builder

	// Category tabs
	var tabs []string
	for i, category := range m.categories {
		tab := category.icon + " " + category.name
		if i == m.selectedCategory {
			tabs = append(tabs, selectedItemStyle.Render("▶ "+tab))
		} else {
			tabs = append(tabs, itemStyle.Render("  "+tab))
		}
	}
	b.WriteString("  " + strings.Join(tabs, "  ") + "\n\n")

	// Current category items
	currentCategory := m.categories[m.selectedCategory]
	for i, item := range currentCategory.items {
		var line strings.Builder

		// Shortcut
		shortcut := shortcutStyle.Render(fmt.Sprintf("[%s]", item.shortcut))

		// Title with AI indicator
		title := item.title
		if item.requiresAI {
			title += " " + warningStyle.Render("[AI]")
		}

		// Selected indicator
		if i == m.selectedItem {
			line.WriteString(selectedItemStyle.Render(fmt.Sprintf("  ▶ %s %s", shortcut, title)))
		} else {
			line.WriteString(itemStyle.Render(fmt.Sprintf("    %s %s", shortcut, title)))
		}

		line.WriteString("\n")
		line.WriteString(descStyle.Render(item.desc))
		line.WriteString("\n")

		b.WriteString(line.String())
	}

	return boxStyle.Render(b.String())
}

func (m MenuModel) renderSearchView() string {
	var b strings.Builder

	// Search input
	b.WriteString(searchStyle.Render("  🔍 Search: "))
	b.WriteString(m.searchInput.View())
	b.WriteString("\n\n")

	// Search results
	if len(m.filteredItems) == 0 {
		b.WriteString(mutedTextStyle.Render("  No results found"))
	} else {
		for i, item := range m.filteredItems {
			var line strings.Builder

			// Category badge
			badge := categoryStyle.Render(fmt.Sprintf("[%s]", item.category))

			// Title
			title := item.title
			if item.requiresAI {
				title += " " + warningStyle.Render("[AI]")
			}

			// Selected indicator
			if i == m.selectedItem {
				line.WriteString(selectedItemStyle.Render(fmt.Sprintf("  ▶ %s %s", badge, title)))
			} else {
				line.WriteString(itemStyle.Render(fmt.Sprintf("    %s %s", badge, title)))
			}

			line.WriteString("\n")
			line.WriteString(descStyle.Render(item.desc))
			line.WriteString("\n")

			b.WriteString(line.String())
		}
	}

	return boxStyle.Render(b.String())
}

func (m MenuModel) renderFooter() string {
	if m.filtering {
		return helpStyle.Render(
			"  ↑↓/jk navigate • enter select • esc cancel search",
		)
	}

	return helpStyle.Render(
		"  ↑↓/jk navigate • ←→/hl/tab switch category • / or ctrl+k search • enter select • q quit",
	)
}

func (m MenuModel) GetChoice() string {
	return m.choice
}

func (m MenuModel) GetSelectedItem() *MenuItem {
	return m.selectedMenuItem
}

func (m MenuItem) Title() string         { return m.title }
func (m MenuItem) Description() string   { return m.desc }
func (m MenuItem) FilterValue() string   { return m.title }
func (m MenuItem) NeedsInput() bool      { return m.needsInput }
func (m MenuItem) InputPrompts() []string { return m.inputPrompts }
func (m MenuItem) Command() string       { return m.command }

// FormatCommand formats the selected command for execution
func FormatCommand(choice string) []string {
	parts := strings.Split(choice, " ")
	return parts
}
