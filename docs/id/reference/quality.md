# anaphase quality

Code quality tools untuk linting, formatting, dan validasi kode yang dihasilkan.

::: info
**Akses Cepat**: Jalankan `anaphase` (tanpa argumen) untuk mengakses menu interaktif dan pilih "Quality Tools" untuk interface visual.
:::

## Overview

Command `quality` menyediakan tools terintegrasi untuk menjaga kualitas kode, memastikan konsistensi, dan menangkap error lebih awal. Command ini bekerja dengan kode yang dihasilkan maupun yang ditulis manual.

::: info
**Auto-Install**: Menu interaktif mendeteksi quality tools yang hilang (golangci-lint, goimports) dan menawarkan untuk menginstalnya secara otomatis.
:::

## Penggunaan

### Menu Interaktif (Disarankan)

```bash
anaphase
```

Pilih **"Quality Tools"** dari menu. Interface menampilkan:
- Run Linter (golangci-lint atau go vet)
- Format Code (gofmt + goimports)
- Validate Code (syntax, vet, build)
- Auto-install tools yang hilang

### Mode CLI Langsung

```bash
anaphase quality lint [path]
anaphase quality format [path]
anaphase quality validate
```

## Subcommand

### lint

Jalankan linter pada codebase Anda untuk menangkap potential issue.

```bash
anaphase quality lint [path] [flags]
```

**Fitur:**
- Pemilihan tool otomatis (golangci-lint atau go vet)
- Prompt auto-install jika golangci-lint hilang
- Dukungan multiple linter (50+ dengan golangci-lint)
- Kemampuan auto-fix
- Laporan error detail

**Contoh:**

```bash
# Menu interaktif (disarankan)
anaphase
# Pilih "Quality Tools" → "Run Linter"

# CLI: Lint seluruh proyek
anaphase quality lint

# Lint direktori tertentu
anaphase quality lint ./internal/core

# Lint single file
anaphase quality lint ./internal/core/entity/user.go

# Auto-fix issue
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

**Flag:**

| Flag | Deskripsi |
|------|-------------|
| `--fix` | Otomatis fix issue jika memungkinkan |

### format

Format kode menggunakan gofmt dan organize import.

```bash
anaphase quality format [path] [flags]
```

**Fitur:**
- Code formatting dengan gofmt
- Import organization dengan goimports
- Batch processing
- Tampilkan diff atau tulis in-place

**Contoh:**

```bash
# Menu interaktif (disarankan)
anaphase
# Pilih "Quality Tools" → "Format Code"

# CLI: Format seluruh proyek
anaphase quality format

# Format direktori tertentu
anaphase quality format ./internal/core

# Format single file
anaphase quality format ./internal/core/entity/user.go

# Tampilkan diff tanpa menulis
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

**Flag:**

| Flag | Short | Default | Deskripsi |
|------|-------|---------|-------------|
| `--write` | `-w` | `true` | Tulis hasil ke source file |

### validate

Validasi komprehensif: syntax check, go vet, dan build test.

```bash
anaphase quality validate
```

**Fitur:**
- Proses validasi 3 langkah
- Syntax checking
- Static analysis (go vet)
- Build verification

**Contoh:**

```bash
# Menu interaktif (disarankan)
anaphase
# Pilih "Quality Tools" → "Validate Code"

# CLI: Validasi seluruh proyek
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

## Linter

### golangci-lint (Disarankan)

Jika terinstal, Anaphase menggunakan golangci-lint yang mencakup 50+ linter.

::: tip
**Auto-Install**: Menu interaktif (`anaphase`) mendeteksi jika golangci-lint hilang dan menawarkan untuk menginstalnya secara otomatis.
:::

**Install Manual:**
```bash
# macOS
brew install golangci-lint

# Linux/WSL
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Windows (PowerShell)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

**Konfigurasi (.golangci.yml):**
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

Jika golangci-lint tidak terinstal, Anaphase secara otomatis fallback ke go vet:

**Fitur:**
- Built-in Go tool (selalu tersedia)
- Tidak perlu instalasi
- Menangkap kesalahan umum
- Cepat dan andal
- Bekerja langsung out of the box

## Integrasi Workflow

### Pre-Commit Hook

Buat `.git/hooks/pre-commit`:

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

Buat executable:
```bash
chmod +x .git/hooks/pre-commit
```

### Integrasi CI/CD

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

### Integrasi VS Code

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

## Masalah Umum & Fix

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

Atau manual:
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

Hasil:
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

Hasil:
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

**Fix (dengan golangci-lint):**
```bash
anaphase quality lint
```

Menampilkan:
```
file.go:10:7: Error return value not checked
```

Kode yang benar:
```go
file, err := os.Open("config.yaml")
if err != nil {
    return fmt.Errorf("open config: %w", err)
}
defer file.Close()
```

## Best Practice

### 1. Jalankan Quality Check Sering

```bash
# Sebelum commit
anaphase quality validate

# Selama development
anaphase quality format
anaphase quality lint
```

### 2. Format on Save

Konfigurasi editor Anda untuk format otomatis:
- VS Code: `"editor.formatOnSave": true`
- GoLand: Settings → Tools → File Watchers

### 3. Gunakan di CI/CD

Selalu jalankan quality check di CI/CD pipeline Anda:
```yaml
- name: Quality
  run: |
    anaphase quality format --write=false
    anaphase quality lint
    anaphase quality validate
```

### 4. Fix Issue Segera

Jangan akumulasi technical debt:
```bash
# Setelah generate kode
anaphase gen domain "User"
anaphase quality validate

# Fix issue segera
anaphase quality lint --fix
```

### 5. Konfigurasi Linter

Buat `.golangci.yml` untuk standar tim yang konsisten:
```yaml
linters:
  disable:
    - exhaustruct  # Terlalu strict untuk generated code
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

## Perbandingan Tools

| Tool | Kecepatan | Fitur | Auto-Fix | Install |
|------|-------|----------|----------|---------|
| **gofmt** | ⚡⚡⚡⚡⚡ | Format only | ✅ | Built-in |
| **goimports** | ⚡⚡⚡⚡ | Format + imports | ✅ | go install |
| **go vet** | ⚡⚡⚡⚡ | Static analysis | ❌ | Built-in |
| **golangci-lint** | ⚡⚡⚡ | 50+ linter | ✅ Partial | External |

## Referensi Cepat

```bash
# Workflow harian
anaphase quality format          # Format kode
anaphase quality lint --fix      # Fix issue sederhana
anaphase quality validate        # Full validation

# Sebelum commit
anaphase quality validate

# CI/CD
anaphase quality format --write=false  # Cek format
anaphase quality lint                  # Cek issue
anaphase quality validate              # Full check
```

## Lihat Juga

- [anaphase gen domain](/reference/gen-domain) - Generate kode domain
- [Installation](/guide/installation) - Install quality tools
- [Troubleshooting](/guide/troubleshooting) - Masalah umum
