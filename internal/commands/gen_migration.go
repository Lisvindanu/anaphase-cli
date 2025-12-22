package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lisvindanu/anaphase-cli/internal/ui"
	"github.com/lisvindanu/anaphase-cli/pkg/fileutil"
	"github.com/spf13/cobra"
)

var (
	migrationOutput string
	migrationDriver string
)

var genMigrationCmd = &cobra.Command{
	Use:   "migration <name>",
	Short: "Generate database migration files",
	Long: `Generate database migration files with up and down scripts.

Supported drivers:
  - postgres (default)
  - mysql
  - sqlite

The migration files will be named with a timestamp prefix for ordering.

Example:
  anaphase gen migration create_users_table
  anaphase gen migration add_email_to_users --driver postgres
  anaphase gen migration create_orders_table --output db/migrations`,
	Args: cobra.ExactArgs(1),
	RunE: runGenMigration,
}

func init() {
	genCmd.AddCommand(genMigrationCmd)

	genMigrationCmd.Flags().StringVar(&migrationOutput, "output", "db/migrations", "Output directory for migration files")
	genMigrationCmd.Flags().StringVar(&migrationDriver, "driver", "postgres", "Database driver (postgres, mysql, sqlite)")
}

func runGenMigration(cmd *cobra.Command, args []string) error {
	name := args[0]

	fmt.Println(ui.RenderTitle("Migration Generator"))
	ui.PrintInfo(fmt.Sprintf("Name: %s", name))
	ui.PrintInfo(fmt.Sprintf("Driver: %s", migrationDriver))
	ui.PrintInfo(fmt.Sprintf("Output: %s", migrationOutput))
	fmt.Println()

	// Validate driver
	validDrivers := []string{"postgres", "mysql", "sqlite"}
	valid := false
	for _, d := range validDrivers {
		if d == migrationDriver {
			valid = true
			break
		}
	}

	if !valid {
		ui.PrintError(fmt.Sprintf("Invalid driver: %s", migrationDriver))
		fmt.Println("Valid drivers:", validDrivers)
		return fmt.Errorf("invalid driver")
	}

	// Create output directory
	if err := fileutil.EnsureDir(migrationOutput); err != nil {
		ui.PrintError(fmt.Sprintf("Failed to create directory: %v", err))
		return err
	}

	// Generate timestamp
	timestamp := time.Now().Format("20060102150405")

	// Generate migration files
	fmt.Println("📝 Generating migration files...")

	upFile := filepath.Join(migrationOutput, fmt.Sprintf("%s_%s.up.sql", timestamp, name))
	downFile := filepath.Join(migrationOutput, fmt.Sprintf("%s_%s.down.sql", timestamp, name))

	// Generate content based on migration name
	upContent := generateUpMigration(name, migrationDriver)
	downContent := generateDownMigration(name, migrationDriver)

	// Write files
	if err := os.WriteFile(upFile, []byte(upContent), 0644); err != nil {
		ui.PrintError(fmt.Sprintf("Failed to write up migration: %v", err))
		return err
	}

	if err := os.WriteFile(downFile, []byte(downContent), 0644); err != nil {
		ui.PrintError(fmt.Sprintf("Failed to write down migration: %v", err))
		return err
	}

	// Show generated files
	fmt.Println(ui.SuccessStyle.Render("\nGenerated Files:"))
	fmt.Println(ui.RenderListItem(upFile, true))
	fmt.Println(ui.RenderListItem(downFile, true))

	fmt.Println()
	ui.PrintSuccess("Migration files generated successfully!")

	// Show usage instructions
	fmt.Println(ui.RenderSubtle("\nUsage Instructions:"))
	fmt.Println("  1. Edit the SQL files to define your migration")
	fmt.Println("  2. Use a migration tool to apply migrations:")
	fmt.Println()
	fmt.Println("     Using golang-migrate:")
	fmt.Println("       migrate -path " + migrationOutput + " -database \"<connection-string>\" up")
	fmt.Println()
	fmt.Println("     Using goose:")
	fmt.Println("       goose -dir " + migrationOutput + " postgres \"<connection-string>\" up")
	fmt.Println()

	fmt.Println(ui.RenderSubtle("Next Steps:"))
	fmt.Println("  1. Edit migration files with your schema changes")
	fmt.Println("  2. Test migrations in development environment")
	fmt.Println("  3. Apply to production using your migration tool")
	fmt.Println()

	return nil
}

func generateUpMigration(name, driver string) string {
	var content strings.Builder

	content.WriteString(fmt.Sprintf("-- Migration: %s\n", name))
	content.WriteString(fmt.Sprintf("-- Created: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("-- Driver: %s\n\n", driver))

	// Try to infer migration type from name
	nameLower := strings.ToLower(name)

	if strings.HasPrefix(nameLower, "create_") && strings.HasSuffix(nameLower, "_table") {
		// Create table migration
		tableName := extractTableName(name, "create_", "_table")
		content.WriteString(generateCreateTable(tableName, driver))
	} else if strings.HasPrefix(nameLower, "add_") && strings.Contains(nameLower, "_to_") {
		// Add column migration
		parts := strings.Split(name, "_to_")
		if len(parts) == 2 {
			columnName := strings.TrimPrefix(parts[0], "add_")
			tableName := strings.TrimSuffix(parts[1], "_table")
			content.WriteString(generateAddColumn(tableName, columnName, driver))
		} else {
			content.WriteString("-- TODO: Add your migration SQL here\n\n")
		}
	} else if strings.HasPrefix(nameLower, "drop_") && strings.HasSuffix(nameLower, "_table") {
		// Drop table migration
		tableName := extractTableName(name, "drop_", "_table")
		content.WriteString(generateDropTable(tableName, driver))
	} else if strings.HasPrefix(nameLower, "create_index_") {
		// Create index migration
		content.WriteString("-- TODO: Add CREATE INDEX statement\n")
		content.WriteString("-- Example:\n")
		content.WriteString("-- CREATE INDEX idx_table_column ON table_name(column_name);\n\n")
	} else {
		// Generic migration
		content.WriteString("-- TODO: Add your migration SQL here\n\n")
	}

	return content.String()
}

func generateDownMigration(name, driver string) string {
	var content strings.Builder

	content.WriteString(fmt.Sprintf("-- Rollback: %s\n", name))
	content.WriteString(fmt.Sprintf("-- Created: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	content.WriteString(fmt.Sprintf("-- Driver: %s\n\n", driver))

	// Try to infer migration type from name
	nameLower := strings.ToLower(name)

	if strings.HasPrefix(nameLower, "create_") && strings.HasSuffix(nameLower, "_table") {
		// Drop table for rollback
		tableName := extractTableName(name, "create_", "_table")
		content.WriteString(generateDropTable(tableName, driver))
	} else if strings.HasPrefix(nameLower, "add_") && strings.Contains(nameLower, "_to_") {
		// Drop column for rollback
		parts := strings.Split(name, "_to_")
		if len(parts) == 2 {
			columnName := strings.TrimPrefix(parts[0], "add_")
			tableName := strings.TrimSuffix(parts[1], "_table")
			content.WriteString(generateDropColumn(tableName, columnName, driver))
		} else {
			content.WriteString("-- TODO: Add your rollback SQL here\n\n")
		}
	} else if strings.HasPrefix(nameLower, "drop_") && strings.HasSuffix(nameLower, "_table") {
		// Create table for rollback (need to restore)
		content.WriteString("-- TODO: Restore table structure\n")
		content.WriteString("-- You may need to backup data before dropping the table\n\n")
	} else {
		// Generic rollback
		content.WriteString("-- TODO: Add your rollback SQL here\n\n")
	}

	return content.String()
}

func extractTableName(name, prefix, suffix string) string {
	tableName := strings.TrimPrefix(strings.ToLower(name), prefix)
	tableName = strings.TrimSuffix(tableName, suffix)
	return tableName
}

func generateCreateTable(tableName, driver string) string {
	var content strings.Builder

	content.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n", tableName))
	content.WriteString("    id BIGSERIAL PRIMARY KEY,\n")
	content.WriteString("    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,\n")
	content.WriteString("    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP\n")
	content.WriteString(");\n\n")

	// Add trigger for updated_at (PostgreSQL)
	if driver == "postgres" {
		content.WriteString("-- Trigger to auto-update updated_at\n")
		content.WriteString(fmt.Sprintf("CREATE OR REPLACE FUNCTION update_%s_updated_at()\n", tableName))
		content.WriteString("RETURNS TRIGGER AS $$\n")
		content.WriteString("BEGIN\n")
		content.WriteString("    NEW.updated_at = CURRENT_TIMESTAMP;\n")
		content.WriteString("    RETURN NEW;\n")
		content.WriteString("END;\n")
		content.WriteString("$$ LANGUAGE plpgsql;\n\n")

		content.WriteString(fmt.Sprintf("CREATE TRIGGER trigger_%s_updated_at\n", tableName))
		content.WriteString(fmt.Sprintf("    BEFORE UPDATE ON %s\n", tableName))
		content.WriteString("    FOR EACH ROW\n")
		content.WriteString(fmt.Sprintf("    EXECUTE FUNCTION update_%s_updated_at();\n\n", tableName))
	}

	return content.String()
}

func generateDropTable(tableName, driver string) string {
	if driver == "postgres" {
		return fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;\n\n", tableName)
	}
	return fmt.Sprintf("DROP TABLE IF EXISTS %s;\n\n", tableName)
}

func generateAddColumn(tableName, columnName, driver string) string {
	// Determine column type based on name
	columnType := "VARCHAR(255)"
	if strings.Contains(columnName, "id") {
		columnType = "BIGINT"
	} else if strings.Contains(columnName, "amount") || strings.Contains(columnName, "price") {
		columnType = "DECIMAL(10,2)"
	} else if strings.Contains(columnName, "count") || strings.Contains(columnName, "quantity") {
		columnType = "INTEGER"
	} else if strings.Contains(columnName, "is_") || strings.Contains(columnName, "has_") {
		columnType = "BOOLEAN DEFAULT FALSE"
	} else if strings.Contains(columnName, "date") || strings.Contains(columnName, "at") {
		columnType = "TIMESTAMP"
	}

	return fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;\n\n", tableName, columnName, columnType)
}

func generateDropColumn(tableName, columnName, driver string) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN IF EXISTS %s;\n\n", tableName, columnName)
}
