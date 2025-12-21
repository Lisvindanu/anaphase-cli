package commands

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lisvindanu/anaphase-cli/internal/generator"
	"github.com/lisvindanu/anaphase-cli/internal/ui"
	"github.com/spf13/cobra"
)

var genTestCmd = &cobra.Command{
	Use:   "test",
	Short: "Generate comprehensive tests for a domain",
	Long: `Generate unit and integration tests for domain entities, repositories, and handlers.

This command will:
- Generate entity tests (constructors, validation, business logic)
- Generate repository tests (unit and integration)
- Generate handler tests (HTTP endpoints)
- Generate test helpers and mocks

Example:
  anaphase gen test --domain customer
  anaphase gen test --domain product --type unit
  anaphase gen test --domain order --type integration`,
	RunE: runGenTest,
}

var (
	testDomain   string
	testType     string
	testCoverage bool
)

func init() {
	genTestCmd.Flags().StringVar(&testDomain, "domain", "", "Domain to generate tests for (required)")
	genTestCmd.Flags().StringVar(&testType, "type", "all", "Test type: unit, integration, or all")
	genTestCmd.Flags().BoolVar(&testCoverage, "coverage", false, "Generate coverage report")

	genTestCmd.MarkFlagRequired("domain")
	genCmd.AddCommand(genTestCmd)
}

func runGenTest(cmd *cobra.Command, args []string) error {
	if testDomain == "" {
		return fmt.Errorf("domain is required")
	}

	// Create steps for progress
	steps := []string{
		"Scanning domain structure",
		"Generating entity tests",
		"Generating repository tests",
		"Generating handler tests",
		"Generating test helpers",
		"Formatting test files",
	}

	// Create progress model
	progress := ui.NewMultiStepProgress(steps)

	// Create Bubble Tea program
	model := &testGeneratorModel{
		domain:   testDomain,
		testType: testType,
		progress: progress,
	}

	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run UI: %w", err)
	}

	// Check for errors
	if m, ok := finalModel.(*testGeneratorModel); ok {
		if m.err != nil {
			return m.err
		}

		// Print success summary
		fmt.Println()
		fmt.Println(ui.RenderSuccess("Test generation complete!"))
		fmt.Println()
		fmt.Println(ui.RenderSubtle("Generated files:"))
		for _, file := range m.generatedFiles {
			fmt.Println(ui.RenderListItem(file, true))
		}
		fmt.Println()
		fmt.Println(ui.RenderInfo("Run tests with: go test ./..."))
		if testCoverage {
			fmt.Println(ui.RenderInfo("Run coverage: go test -cover ./..."))
		}
	}

	return nil
}

// testGeneratorModel is the Bubble Tea model for test generation
type testGeneratorModel struct {
	domain         string
	testType       string
	progress       *ui.MultiStepProgress
	done           bool
	err            error
	generatedFiles []string
}

// Init initializes the model
func (m *testGeneratorModel) Init() tea.Cmd {
	return m.runGeneration()
}

// generationDoneMsg indicates generation completion
type generationDoneMsg struct {
	files []string
	err   error
}

// runGeneration executes the test generation
func (m *testGeneratorModel) runGeneration() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		// Create test generator
		gen := generator.NewTestGenerator(&generator.TestConfig{
			Domain:   m.domain,
			TestType: m.testType,
		})

		// Step 1: Scan domain
		if err := gen.ScanDomain(); err != nil {
			return generationDoneMsg{err: err}
		}

		// Generate tests based on type
		var files []string
		var err error

		switch m.testType {
		case "unit":
			files, err = gen.GenerateUnitTests(ctx)
		case "integration":
			files, err = gen.GenerateIntegrationTests(ctx)
		default: // "all"
			files, err = gen.GenerateAllTests(ctx)
		}

		return generationDoneMsg{
			files: files,
			err:   err,
		}
	}
}

// Update updates the model
func (m *testGeneratorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case generationDoneMsg:
		m.done = true
		m.err = msg.err
		m.generatedFiles = msg.files
		if m.err == nil {
			m.progress.SetCurrent(m.progress.total)
		}
		return m, tea.Quit
	}

	return m, nil
}

// View renders the model
func (m *testGeneratorModel) View() string {
	if m.done {
		if m.err != nil {
			return ui.RenderError(fmt.Sprintf("Test generation failed: %v", m.err))
		}
		return ""
	}

	header := ui.RenderTitle(fmt.Sprintf("Generating tests for %s domain", m.domain))
	return fmt.Sprintf("%s\n\n%s\n", header, m.progress.View())
}
