package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ProgressBar represents a progress bar
type ProgressBar struct {
	width    int
	progress float64
	style    lipgloss.Style
	fillChar string
	emptyChar string
}

// NewProgressBar creates a new progress bar
func NewProgressBar(width int) *ProgressBar {
	return &ProgressBar{
		width:     width,
		progress:  0,
		style:     ProgressBarStyle,
		fillChar:  "█",
		emptyChar: "░",
	}
}

// SetProgress sets the progress (0.0 to 1.0)
func (p *ProgressBar) SetProgress(progress float64) {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	p.progress = progress
}

// Increment increments the progress
func (p *ProgressBar) Increment(amount float64) {
	p.SetProgress(p.progress + amount)
}

// View renders the progress bar
func (p *ProgressBar) View() string {
	filled := int(float64(p.width) * p.progress)
	empty := p.width - filled

	bar := strings.Repeat(p.fillChar, filled) + strings.Repeat(p.emptyChar, empty)
	percentage := fmt.Sprintf(" %.0f%%", p.progress*100)

	style := p.style
	if p.progress >= 1.0 {
		style = ProgressCompleteStyle
	}

	return style.Render(bar) + percentage
}

// ViewWithLabel renders the progress bar with a label
func (p *ProgressBar) ViewWithLabel(label string) string {
	return fmt.Sprintf("%s\n%s", label, p.View())
}

// MultiStepProgress represents a multi-step progress tracker
type MultiStepProgress struct {
	steps   []string
	current int
	Total   int // Exported for external access
}

// NewMultiStepProgress creates a new multi-step progress tracker
func NewMultiStepProgress(steps []string) *MultiStepProgress {
	return &MultiStepProgress{
		steps:   steps,
		current: 0,
		Total:   len(steps),
	}
}

// Next moves to the next step
func (m *MultiStepProgress) Next() {
	if m.current < m.Total {
		m.current++
	}
}

// SetCurrent sets the current step
func (m *MultiStepProgress) SetCurrent(step int) {
	if step >= 0 && step <= m.Total {
		m.current = step
	}
}

// View renders the multi-step progress
func (m *MultiStepProgress) View() string {
	var lines []string

	lines = append(lines, SubtitleStyle.Render(
		fmt.Sprintf("Progress: %d/%d steps", m.current, m.Total),
	))
	lines = append(lines, "")

	for i, step := range m.steps {
		var line string
		if i < m.current {
			// Completed
			line = CheckmarkStyle.Render() + " " + SubtleStyle.Render(step)
		} else if i == m.current {
			// Current
			line = "▶ " + lipgloss.NewStyle().Bold(true).Render(step)
		} else {
			// Pending
			line = "  " + SubtleStyle.Render(step)
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// IsComplete returns true if all steps are complete
func (m *MultiStepProgress) IsComplete() bool {
	return m.current >= m.Total
}

// Progress returns the current progress (0.0 to 1.0)
func (m *MultiStepProgress) Progress() float64 {
	if m.Total == 0 {
		return 1.0
	}
	return float64(m.current) / float64(m.Total)
}
