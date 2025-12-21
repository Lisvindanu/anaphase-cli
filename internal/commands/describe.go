package commands

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lisvindanu/anaphase-cli/internal/generator"
	"github.com/lisvindanu/anaphase-cli/internal/ui"
	"github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Generate architecture diagrams and documentation",
	Long: `Generate visual architecture diagrams and documentation for your project.

This command will:
- Scan your project structure
- Generate architecture diagrams (Mermaid, ASCII)
- Document domain relationships
- Create dependency graphs

Example:
  anaphase describe                    # Full architecture diagram
  anaphase describe --format mermaid   # Mermaid diagram
  anaphase describe --format ascii     # ASCII art diagram
  anaphase describe --output arch.md   # Save to file`,
	RunE: runDescribe,
}

var (
	describeFormat string
	describeOutput string
	describeType   string
)

func init() {
	describeCmd.Flags().StringVar(&describeFormat, "format", "mermaid", "Diagram format: mermaid, ascii, or both")
	describeCmd.Flags().StringVar(&describeOutput, "output", "", "Output file (default: stdout)")
	describeCmd.Flags().StringVar(&describeType, "type", "all", "Diagram type: all, domain, layers, dependencies")

	rootCmd.AddCommand(describeCmd)
}

func runDescribe(cmd *cobra.Command, args []string) error {
	// Create steps for progress
	steps := []string{
		"Scanning project structure",
		"Analyzing domains",
		"Building dependency graph",
		"Generating diagrams",
	}

	progress := ui.NewMultiStepProgress(steps)

	// Create Bubble Tea program
	model := &describeModel{
		format:   describeFormat,
		output:   describeOutput,
		diagType: describeType,
		progress: progress,
	}

	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run UI: %w", err)
	}

	// Check for errors
	if m, ok := finalModel.(*describeModel); ok {
		if m.err != nil {
			return m.err
		}

		// Print or save output
		if describeOutput != "" {
			fmt.Println(ui.RenderSuccess("Architecture diagram saved to: " + describeOutput))
		} else {
			fmt.Println()
			fmt.Println(m.diagram)
		}
	}

	return nil
}

// describeModel is the Bubble Tea model for describe
type describeModel struct {
	format   string
	output   string
	diagType string
	progress *ui.MultiStepProgress
	diagram  string
	done     bool
	err      error
}

// Init initializes the model
func (m *describeModel) Init() tea.Cmd {
	return m.runGeneration()
}

// describeDoneMsg indicates generation completion
type describeDoneMsg struct {
	diagram string
	err     error
}

// runGeneration executes the diagram generation
func (m *describeModel) runGeneration() tea.Cmd {
	return func() tea.Msg {
		gen := generator.NewDiagramGenerator(&generator.DiagramConfig{
			Format: m.format,
			Type:   m.diagType,
		})

		diagram, err := gen.Generate()
		if err != nil {
			return describeDoneMsg{err: err}
		}

		// Save to file if output specified
		if m.output != "" {
			if err := os.WriteFile(m.output, []byte(diagram), 0644); err != nil {
				return describeDoneMsg{err: fmt.Errorf("write output file: %w", err)}
			}
		}

		return describeDoneMsg{
			diagram: diagram,
			err:     nil,
		}
	}
}

// Update updates the model
func (m *describeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case describeDoneMsg:
		m.done = true
		m.err = msg.err
		m.diagram = msg.diagram
		if m.err == nil {
			m.progress.SetCurrent(m.progress.Total)
		}
		return m, tea.Quit
	}

	return m, nil
}

// View renders the model
func (m *describeModel) View() string {
	if m.done {
		if m.err != nil {
			return ui.RenderError(fmt.Sprintf("Diagram generation failed: %v", m.err))
		}
		return ""
	}

	header := ui.RenderTitle("Generating Architecture Diagrams")
	return fmt.Sprintf("%s\n\n%s\n", header, m.progress.View())
}
