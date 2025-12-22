package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Spinner represents a loading spinner
type Spinner struct {
	frames   []string
	index    int
	message  string
	active   bool
	style    lipgloss.Style
	duration time.Duration
}

// NewSpinner creates a new spinner
func NewSpinner(message string) *Spinner {
	return &Spinner{
		frames:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		index:   0,
		message: message,
		active:  true,
		style: lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true),
		duration: 80 * time.Millisecond,
	}
}

// tickMsg represents a spinner tick
type tickMsg time.Time

// Init initializes the spinner
func (s *Spinner) Init() tea.Cmd {
	return s.tick()
}

// Update updates the spinner state
func (s *Spinner) Update(msg tea.Msg) (*Spinner, tea.Cmd) {
	switch msg.(type) {
	case tickMsg:
		if s.active {
			s.index = (s.index + 1) % len(s.frames)
			return s, s.tick()
		}
	}
	return s, nil
}

// View renders the spinner
func (s *Spinner) View() string {
	if !s.active {
		return ""
	}
	frame := s.frames[s.index]
	return s.style.Render(frame) + " " + s.message
}

// tick generates the next tick
func (s *Spinner) tick() tea.Cmd {
	return tea.Tick(s.duration, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// SetMessage updates the spinner message
func (s *Spinner) SetMessage(msg string) {
	s.message = msg
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	s.active = false
}

// SpinnerModel is a bubbletea model for spinner
type SpinnerModel struct {
	spinner *Spinner
	done    bool
	err     error
	task    func() error
}

// NewSpinnerModel creates a new spinner model
func NewSpinnerModel(message string, task func() error) SpinnerModel {
	return SpinnerModel{
		spinner: NewSpinner(message),
		task:    task,
	}
}

// Init initializes the model
func (m SpinnerModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Init(),
		m.runTask(),
	)
}

// taskDoneMsg indicates task completion
type taskDoneMsg struct {
	err error
}

// runTask executes the task
func (m SpinnerModel) runTask() tea.Cmd {
	return func() tea.Msg {
		err := m.task()
		return taskDoneMsg{err: err}
	}
}

// Update updates the model
func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case taskDoneMsg:
		m.done = true
		m.err = msg.err
		m.spinner.Stop()
		return m, tea.Quit

	case tickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the model
func (m SpinnerModel) View() string {
	if m.done {
		if m.err != nil {
			return RenderError(fmt.Sprintf("Failed: %v", m.err))
		}
		return RenderSuccess("Done!")
	}
	return m.spinner.View()
}
