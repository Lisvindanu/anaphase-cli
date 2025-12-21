package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	Primary   = lipgloss.Color("#667eea")
	Secondary = lipgloss.Color("#764ba2")
	Success   = lipgloss.Color("#10b981")
	Error     = lipgloss.Color("#ef4444")
	Warning   = lipgloss.Color("#f59e0b")
	Info      = lipgloss.Color("#3b82f6")
	Subtle    = lipgloss.Color("#6b7280")

	// Text Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Primary).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Italic(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning)

	InfoStyle = lipgloss.NewStyle().
			Foreground(Info)

	SubtleStyle = lipgloss.NewStyle().
			Foreground(Subtle)

	// Box Styles
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)

	// List Styles
	ListItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	CheckmarkStyle = lipgloss.NewStyle().
			Foreground(Success).
			SetString("✓")

	CrossStyle = lipgloss.NewStyle().
			Foreground(Error).
			SetString("✗")

	// Progress Styles
	ProgressBarStyle = lipgloss.NewStyle().
				Foreground(Primary)

	ProgressCompleteStyle = lipgloss.NewStyle().
				Foreground(Success)
)

// RenderTitle renders a styled title
func RenderTitle(title string) string {
	return TitleStyle.Render("⚡ " + title)
}

// RenderSuccess renders a success message
func RenderSuccess(msg string) string {
	return SuccessStyle.Render("✓ " + msg)
}

// RenderError renders an error message
func RenderError(msg string) string {
	return ErrorStyle.Render("✗ " + msg)
}

// RenderWarning renders a warning message
func RenderWarning(msg string) string {
	return WarningStyle.Render("⚠ " + msg)
}

// RenderInfo renders an info message
func RenderInfo(msg string) string {
	return InfoStyle.Render("ℹ " + msg)
}

// RenderSubtle renders subtle text
func RenderSubtle(msg string) string {
	return SubtleStyle.Render(msg)
}

// RenderBox renders content in a box
func RenderBox(content string) string {
	return BoxStyle.Render(content)
}

// RenderListItem renders a list item with checkmark
func RenderListItem(text string, checked bool) string {
	if checked {
		return CheckmarkStyle.Render() + " " + text
	}
	return "  " + text
}
