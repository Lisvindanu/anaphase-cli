# anaphase quality

Code quality tools for linting, formatting, and validating generated code.

## Overview

The `quality` command provides integrated tools to maintain code quality, ensure consistency, and catch errors early. It works with both generated and hand-written code.

## Subcommands

### lint

Run linters on your codebase to catch potential issues.

```bash
anaphase quality lint [path] [flags]
```

**Features:**
- Automatic tool selection (golangci-lint or go vet)
- Multiple linter support
- Auto-fix capability
- Detailed error reports

**Examples:**

```bash
# Lint entire project
anaphase quality lint

# Lint specific directory
anaphase quality lint ./internal/core

# Lint single file
anaphase quality lint ./internal/core/entity/user.go

# Auto-fix issues
anaphase quality lint --fix
```

**Output:**
```
⚡ Code Linting

ℹ Linting: .

📋 Running golangci-lint...

internal/core/entity/user.go:15:2: ST1003: should not use underscores in Go names
internal/core/service/user.go:42:1: errcheck: error return value not checked

⚠ Linting found issues

ℹ Run with --fix to automatically fix some issues
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--fix` | Automatically fix issues when possible |

### format

Format code using gofmt and organize imports.

```bash
anaphase quality format [path] [flags]
```

**Features:**
- Code formatting with gofmt
- Import organization with goimports
- Batch processing
- Show diff or write in-place

**Examples:**

```bash
# Format entire project
anaphase quality format

# Format specific directory
anaphase quality format ./internal/core

# Format single file
anaphase quality format ./internal/core/entity/user.go

# Show diff without writing
anaphase quality format --write=false
```

**Output:**
```
⚡ Code Formatting

ℹ Formatting: .

📝 Running gofmt...
✓ Formatted 12 file(s)

📦 Running goimports...
✓ Organized imports in 8 file(s)
```

**Flags:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--write` | `-w` | `true` | Write result to source file |

### validate

Comprehensive validation: syntax check, go vet, and build test.

```bash
anaphase quality validate
```

**Features:**
- 3-step validation process
- Syntax checking
- Static analysis (go vet)
- Build verification

**Examples:**

```bash
# Validate entire project
anaphase quality validate
```

**Output:**
```
⚡ Code Validation

📋 Step 1/3: Checking syntax...
✓ Syntax OK

📋 Step 2/3: Running go vet...
✓ No vet issues found

📋 Step 3/3: Building code...
✓ Build successful

✓ Validation complete! Code is ready to use.
```

## Linters

### golangci-lint (Recommended)

If installed, Anaphase uses golangci-lint which includes 50+ linters:

**Install:**
```bash
# macOS
brew install golangci-lint

# Linux/WSL
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Windows (PowerShell)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

**Configuration (.golangci.yml):**
```yaml
linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    - revive
    - gofmt
    - goimports

linters-settings:
  errcheck:
    check-blank: true

  revive:
    rules:
      - name: var-naming
      - name: exported
```

### go vet (Fallback)

If golangci-lint is not installed, Anaphase falls back to go vet:

**Features:**
- Built-in Go tool
- No installation required
- Catches common mistakes
- Fast and reliable

## Workflow Integration

### Pre-Commit Hook

Create `.git/hooks/pre-commit`:

```bash
#!/bin/bash

echo "Running code quality checks..."

# Format code
anaphase quality format

# Lint code
anaphase quality lint

# Validate
if ! anaphase quality validate; then
    echo "❌ Quality checks failed. Commit aborted."
    exit 1
fi

echo "✅ Quality checks passed!"
```

Make executable:
```bash
chmod +x .git/hooks/pre-commit
```

### CI/CD Integration

**GitHub Actions (.github/workflows/quality.yml):**

```yaml
name: Code Quality

on: [push, pull_request]

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Install Anaphase
        run: |
          go install github.com/lisvindanu/anaphase-cli/cmd/anaphase@latest

      - name: Install golangci-lint
        run: |
          curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

      - name: Format Check
        run: |
          anaphase quality format
          git diff --exit-code

      - name: Lint
        run: anaphase quality lint

      - name: Validate
        run: anaphase quality validate
```

**GitLab CI (.gitlab-ci.yml):**

```yaml
quality:
  image: golang:1.21
  stage: test
  before_script:
    - go install github.com/lisvindanu/anaphase-cli/cmd/anaphase@latest
    - curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
  script:
    - anaphase quality format
    - anaphase quality lint
    - anaphase quality validate
```

### VS Code Integration

**settings.json:**

```json
{
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "package",
  "editor.formatOnSave": true,
  "go.formatTool": "goimports",

  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  }
}
```

**tasks.json:**

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Quality Check",
      "type": "shell",
      "command": "anaphase quality validate",
      "problemMatcher": "$go",
      "group": {
        "kind": "test",
        "isDefault": true
      }
    }
  ]
}
```

## Common Issues & Fixes

### 1. Unused Variables

**Issue:**
```go
func GetUser(id int) (*User, error) {
    user, err := repo.FindByID(id)  // err unused
    return user, nil
}
```

**Fix:**
```bash
anaphase quality lint --fix
```

Or manually:
```go
func GetUser(id int) (*User, error) {
    user, err := repo.FindByID(id)
    if err != nil {
        return nil, err
    }
    return user, nil
}
```

### 2. Formatting Issues

**Issue:**
```go
func GetUser(id int)(*User,error){
return &User{ID:id},nil
}
```

**Fix:**
```bash
anaphase quality format
```

Result:
```go
func GetUser(id int) (*User, error) {
    return &User{ID: id}, nil
}
```

### 3. Import Organization

**Issue:**
```go
import (
    "github.com/lisvindanu/myapp/internal/core"
    "fmt"
    "time"
    "github.com/google/uuid"
)
```

**Fix:**
```bash
anaphase quality format
```

Result:
```go
import (
    "fmt"
    "time"

    "github.com/google/uuid"

    "github.com/lisvindanu/myapp/internal/core"
)
```

### 4. Error Checking

**Issue:**
```go
file, _ := os.Open("config.yaml")
```

**Fix (with golangci-lint):**
```bash
anaphase quality lint
```

Shows:
```
file.go:10:7: Error return value not checked
```

Correct code:
```go
file, err := os.Open("config.yaml")
if err != nil {
    return fmt.Errorf("open config: %w", err)
}
defer file.Close()
```

## Best Practices

### 1. Run Quality Checks Often

```bash
# Before committing
anaphase quality validate

# During development
anaphase quality format
anaphase quality lint
```

### 2. Format on Save

Configure your editor to format automatically:
- VS Code: `"editor.formatOnSave": true`
- GoLand: Settings → Tools → File Watchers

### 3. Use in CI/CD

Always run quality checks in your CI/CD pipeline:
```yaml
- name: Quality
  run: |
    anaphase quality format --write=false
    anaphase quality lint
    anaphase quality validate
```

### 4. Fix Issues Immediately

Don't accumulate technical debt:
```bash
# After generating code
anaphase gen domain "User"
anaphase quality validate

# Fix any issues immediately
anaphase quality lint --fix
```

### 5. Configure Linters

Create `.golangci.yml` for consistent team standards:
```yaml
linters:
  disable:
    - exhaustruct  # Too strict for generated code
  enable:
    - errcheck
    - gosimple
    - govet

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - errcheck
```

## Tools Comparison

| Tool | Speed | Features | Auto-Fix | Install |
|------|-------|----------|----------|---------|
| **gofmt** | ⚡⚡⚡⚡⚡ | Format only | ✅ | Built-in |
| **goimports** | ⚡⚡⚡⚡ | Format + imports | ✅ | go install |
| **go vet** | ⚡⚡⚡⚡ | Static analysis | ❌ | Built-in |
| **golangci-lint** | ⚡⚡⚡ | 50+ linters | ✅ Partial | External |

## Quick Reference

```bash
# Daily workflow
anaphase quality format          # Format code
anaphase quality lint --fix      # Fix simple issues
anaphase quality validate        # Full validation

# Before commit
anaphase quality validate

# CI/CD
anaphase quality format --write=false  # Check format
anaphase quality lint                  # Check issues
anaphase quality validate              # Full check
```

## See Also

- [anaphase gen domain](/reference/gen-domain) - Generate domain code
- [Installation](/guide/installation) - Install quality tools
- [Troubleshooting](/guide/troubleshooting) - Common issues
