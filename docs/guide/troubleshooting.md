# Troubleshooting Guide

Having issues? This guide covers common problems and how to fix them.

::: tip For Beginners
If this is your first time using Anaphase or Go, read the [Common Issues](#common-issues) section first. Most problems have simple fixes!
:::

## Common Issues

### Interactive Menu Issues

**Problem:**
```bash
$ anaphase
# Menu appears but commands don't work as expected
```

**What it means:**
As of v0.4.0, running `anaphase` without arguments launches an interactive TUI menu. If you experience issues, you can always fall back to direct commands.

**Fix:**

**If menu is unresponsive:**
```bash
# Use direct commands instead
anaphase init my-project
anaphase gen domain --name user

# Or check terminal compatibility
echo $TERM
# Should show something like "xterm-256color"
```

**If menu doesn't display correctly:**
```bash
# Update your terminal or use commands directly
# The menu requires a terminal that supports escape sequences

# Bypass menu entirely
anaphase --help  # Shows all available commands
```

**Navigation tips:**
- Use arrow keys or j/k to navigate
- Press Enter to select
- Press ESC or q to go back/quit
- Press Ctrl+C to exit immediately

---

### Template Mode vs AI Mode Issues

**Problem:**
```bash
$ anaphase gen domain --name user --prompt "User with email"
Error: GEMINI_API_KEY environment variable not set
```

**What it means:**
You're trying to use AI mode without an API key. As of v0.4.0, you can choose between Template and AI modes.

**Fix:**

**Option 1: Use Template Mode (no API key needed)**
```bash
# Template mode is fast and requires no API key
anaphase gen domain --name user --template

# Or use the interactive menu and select Template Mode
anaphase
```

**Option 2: Set up AI Mode**
```bash
# Get API key from https://makersuite.google.com/app/apikey
export GEMINI_API_KEY="your-key-here"

# Now you can use AI mode
anaphase gen domain --name user --prompt "User with email and profile"
```

**When to use each mode:**
- **Template Mode**: Quick scaffolding, simple domains, no API key available
- **AI Mode**: Complex business logic, advanced DDD patterns, domain events

---

### Auto-Setup Issues

**Problem:**
```bash
$ anaphase init my-project
# Dependencies not automatically installed
```

**What it means:**
As of v0.4.0, auto-setup features attempt to run `go mod download` automatically, but may fail if Go is not properly configured.

**Fix:**

**If auto-setup fails:**
```bash
# Manually complete the setup
cd my-project
go mod download
go mod tidy

# Verify Go is configured
go version
go env GOPATH
```

**Disable auto-setup if causing issues:**
```bash
# Generate without auto-setup
anaphase init my-project --no-auto-setup

# Then manually run
cd my-project
go mod download
```

---

### "missing go.sum entry" Error

**Problem:**
```bash
$ make run
internal/config/config.go:6:2: missing go.sum entry for module providing package github.com/spf13/viper
```

**What it means:**
Go dependencies aren't downloaded yet. This happens after `anaphase init`.

**Fix:**
```bash
# Download all dependencies
go mod download

# Or run this (does the same thing)
go mod tidy

# Then try again
make run
```

**Why it happens:**
`anaphase init` creates the project structure but doesn't download dependencies automatically. You need to run `go mod download` first.

---

### "command not found: anaphase"

**Problem:**
```bash
$ anaphase --version
zsh: command not found: anaphase
```

**What it means:**
The `anaphase` binary isn't in your PATH.

**Fix:**

**Option 1: Use the install script (recommended)**
```bash
curl -fsSL https://raw.githubusercontent.com/lisvindanu/anaphase-cli/main/install.sh | bash
```

**Option 2: Add to PATH manually**
```bash
# Check where Go installs binaries
go env GOPATH

# Add to your shell config
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# Verify
which anaphase
```

---

### "GEMINI_API_KEY not configured"

**Problem:**
```bash
$ anaphase gen domain --name user --prompt "User with email"
Error: GEMINI_API_KEY environment variable not set
```

**What it means:**
The `--prompt` flag requires AI mode, which needs a Gemini API key. However, AI is optional in v0.4.0.

**Fix:**

**Option 1: Use Template Mode (recommended for most cases)**
```bash
# No API key needed - uses templates
anaphase gen domain --name user --template

# Or use the interactive menu
anaphase
# Select "Template Mode" when prompted
```

**Option 2: Set up AI Mode (for advanced features)**

1. **Get API key:**
   - Go to https://makersuite.google.com/app/apikey
   - Sign in with Google
   - Click "Create API Key"
   - Copy the key

2. **Set the key:**
```bash
# Temporary (this session only)
export GEMINI_API_KEY="your-key-here"

# Permanent (add to shell config)
echo 'export GEMINI_API_KEY="your-key-here"' >> ~/.zshrc
source ~/.zshrc
```

3. **Verify:**
```bash
echo $GEMINI_API_KEY
# Should print your key

# Now use AI mode with prompts
anaphase gen domain --name user --prompt "User with email and profile"
```

::: tip AI is Optional
You don't need an API key to use Anaphase! Template mode provides all core functionality without external dependencies. Use AI mode only when you need advanced domain modeling.
:::

---

### Database Connection Failed

**Problem:**
```bash
$ make run
Error: failed to connect to database: connection refused
```

**What it means:**
PostgreSQL isn't running or wrong connection string.

**Fix:**

**Quick fix with Docker:**
```bash
# Start PostgreSQL in Docker
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 \
  postgres:16-alpine

# Wait a few seconds for startup
sleep 3

# Set connection string
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable"

# Try again
make run
```

**Check if PostgreSQL is running:**
```bash
# With Docker
docker ps | grep postgres

# With psql
psql -h localhost -U postgres -d mydb -c "SELECT 1"
```

**Common connection string formats:**
```bash
# Local PostgreSQL
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable"

# Docker PostgreSQL
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable"

# Remote PostgreSQL
export DATABASE_URL="postgres://user:password@host:5432/database?sslmode=require"
```

---

### "make: command not found"

**Problem:**
```bash
$ make run
zsh: command not found: make
```

**What it means:**
`make` isn't installed on your system.

**Fix:**

**macOS:**
```bash
xcode-select --install
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt-get update
sudo apt-get install build-essential
```

**Or run directly without make:**
```bash
# Instead of: make run
go run cmd/api/main.go

# Instead of: make build
go build -o bin/api cmd/api/main.go
```

---

### Import Errors After Generation

**Problem:**
```bash
$ make run
internal/core/entity/user.go:5:2: no required module provides package github.com/google/uuid
```

**What it means:**
New dependencies were added by code generation but not downloaded.

**Fix:**
```bash
# Download missing dependencies
go mod tidy

# Or
go get ./...

# Then run again
make run
```

**Always run after generation:**
```bash
anaphase gen domain --name user --prompt "..."
go mod tidy  # ← Run this!
```

---

### Port Already in Use

**Problem:**
```bash
$ make run
Error: listen tcp :8080: bind: address already in use
```

**What it means:**
Another process is using port 8080.

**Fix:**

**Find what's using the port:**
```bash
# macOS/Linux
lsof -i :8080

# Kill the process
kill -9 <PID>
```

**Or use a different port:**
```bash
export PORT=3000
make run
```

---

### Permission Denied

**Problem:**
```bash
$ make run
zsh: permission denied: ./bin/api
```

**What it means:**
Binary doesn't have execute permissions.

**Fix:**
```bash
chmod +x bin/api
./bin/api
```

---

### "go: cannot find main module"

**Problem:**
```bash
$ go run cmd/api/main.go
go: cannot find main module; see 'go help modules'
```

**What it means:**
You're not in a Go module directory (no `go.mod` file).

**Fix:**
```bash
# Make sure you're in the project directory
cd my-api

# Verify go.mod exists
ls go.mod

# Then run
go run cmd/api/main.go
```

---

## Step-by-Step: First Time Setup

If you're completely new, follow these steps in order:

### 1. Install Go

**Check if already installed:**
```bash
go version
```

**If not installed:**
- Download from https://go.dev/dl/
- Install for your OS
- Verify: `go version`

### 2. Install Anaphase

```bash
# Quick install
curl -fsSL https://raw.githubusercontent.com/lisvindanu/anaphase-cli/main/install.sh | bash

# Verify
anaphase --version
```

### 3. Choose Your Mode

::: tip New in v0.4.0
You can now choose between Template Mode (no setup needed) and AI Mode (requires API key).

**Start with Template Mode** - it's simpler and doesn't require any API keys!
:::

**Option A: Template Mode (Recommended for Beginners)**
```bash
# No additional setup needed!
# Skip to Step 4
```

**Option B: AI Mode (Optional)**

1. Go to https://makersuite.google.com/app/apikey
2. Create API key
3. Set it:
```bash
echo 'export GEMINI_API_KEY="your-key-here"' >> ~/.zshrc
source ~/.zshrc
```

### 4. Start PostgreSQL

**With Docker (easiest):**
```bash
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 \
  postgres:16-alpine
```

**Set connection:**
```bash
echo 'export DATABASE_URL="postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable"' >> ~/.zshrc
source ~/.zshrc
```

### 5. Create Your First Project

::: tip Use the Interactive Menu
As of v0.4.0, you can use the interactive menu for easier project creation:
```bash
anaphase  # Launches interactive menu
```
Or continue with commands as shown below.
:::

**Using Template Mode (No API Key Needed):**
```bash
# Create project
anaphase init my-api
cd my-api

# Auto-setup will run go mod download for you
# If it fails, run manually: go mod download

# Generate a domain using templates
anaphase gen domain --name user --template

# Download new dependencies
go mod tidy

# Generate handler
anaphase gen handler --domain user

# Generate repository
anaphase gen repository --domain user --db postgres

# Wire everything
anaphase wire

# Download final dependencies
go mod tidy

# Run!
make run
```

**Using AI Mode (Requires API Key):**
```bash
# Create project
anaphase init my-api
cd my-api

# Generate a domain with AI
anaphase gen domain --name user --prompt "User with email, name, and optional profile picture URL"

# Rest is the same
go mod tidy
anaphase gen handler --domain user
anaphase gen repository --domain user --db postgres
anaphase wire
go mod tidy
make run
```

### 6. Test Your API

```bash
# In another terminal, create a user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "name": "Test User"
  }'
```

**You should see:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "test@example.com",
  "name": "Test User",
  "createdAt": "2024-01-01T00:00:00Z"
}
```

🎉 **Success!** Your API is working!

---

## Common Mistakes

### ❌ Forgetting `go mod tidy`

```bash
# After ANY generation command, run:
anaphase gen domain --name product --prompt "..."
go mod tidy  # ← Don't forget!
```

### ❌ Wrong directory

```bash
# Wrong - running from home directory
~ $ make run
make: *** No rule to make target 'run'.

# Right - inside project directory
~/my-api $ make run
Starting my-api...
```

### ❌ Not setting environment variables

```bash
# Won't work
make run

# Will work
export GEMINI_API_KEY="..."
export DATABASE_URL="..."
make run
```

### ❌ Using old Go version

```bash
# Check version
go version
# Should be 1.21 or higher

# If too old, update Go
# Download from https://go.dev/dl/
```

---

## Getting Help

### Check Logs

```bash
# Run with verbose output
go run cmd/api/main.go 2>&1 | tee app.log

# Check what's wrong
cat app.log
```

### Verify Installation

```bash
# Check Go
go version
# Should be 1.21+

# Check Anaphase
anaphase --version

# Check environment
echo $GEMINI_API_KEY
echo $DATABASE_URL

# Check database
psql $DATABASE_URL -c "SELECT 1"
```

### Clean Start

If everything is broken, start fresh:

```bash
# Remove the project
rm -rf my-api

# Clear module cache
go clean -modcache

# Start over
anaphase init my-api
cd my-api
go mod download
# ... continue ...
```

---

## Still Stuck?

### Read the Docs
- [Quick Start](/guide/quick-start)
- [Installation](/guide/installation)
- [Architecture](/guide/architecture)

### Check Examples
- [Basic Example](/examples/basic)
- [Multi-Domain](/examples/multi-domain)

### Common Questions

**Q: Do I need to know Go?**
A: Basic Go knowledge helps, but Anaphase generates most code for you.

**Q: Do I need a Gemini API key?**
A: No! As of v0.4.0, Template Mode requires no API key. AI mode is completely optional.

**Q: What's the difference between Template and AI mode?**
A: Template mode uses predefined scaffolding (fast, deterministic). AI mode generates more sophisticated domain models based on natural language prompts.

**Q: Is Gemini API free?**
A: Yes! Free tier includes 60 requests/minute. But remember, you can use Template Mode without any API key.

**Q: Can I use the interactive menu on any terminal?**
A: The menu works on most modern terminals. If you have issues, use direct commands instead.

**Q: Can I use MySQL instead of PostgreSQL?**
A: Yes! Use `--db mysql` when generating repositories.

**Q: How do I add custom logic?**
A: Edit the generated service layer files. See [Custom Handlers](/examples/custom-handlers).

**Q: Can I run without Docker?**
A: Yes, install PostgreSQL natively for your OS.

**Q: Can I switch between Template and AI mode?**
A: Yes! Use Template mode for basic scaffolding, then manually enhance or use AI mode for complex additions.

---

## Pro Tips

### 1. Use `.env` file

Create `.env` in your project:
```bash
GEMINI_API_KEY=your-key-here
DATABASE_URL=postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable
PORT=8080
```

Load it:
```bash
# Install dotenv tool
go install github.com/joho/godotenv/cmd/godotenv@latest

# Run with .env
godotenv go run cmd/api/main.go
```

### 2. Create a setup script

Create `setup.sh`:
```bash
#!/bin/bash
set -e

echo "Setting up project..."

# Download dependencies
go mod download

# Start database
docker run -d --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 \
  postgres:16-alpine

sleep 3

# Apply migrations
psql "postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable" \
  -f internal/adapter/repository/postgres/schema.sql

echo "✅ Setup complete!"
```

Run it:
```bash
chmod +x setup.sh
./setup.sh
```

### 3. Add to Makefile

Edit `Makefile`:
```makefile
.PHONY: setup deps clean

setup: deps
    @echo "Starting PostgreSQL..."
    docker run -d --name postgres \
      -e POSTGRES_PASSWORD=postgres \
      -p 5432:5432 postgres:16-alpine

deps:
    @echo "Downloading dependencies..."
    go mod download
    go mod tidy

clean:
    @echo "Cleaning up..."
    docker stop postgres || true
    docker rm postgres || true
    go clean
```

Now you can:
```bash
make setup   # First time setup
make deps    # Download dependencies
make clean   # Clean everything
```

---

::: tip Remember (v0.4.0+)
1. Try the interactive menu: `anaphase` (no arguments)
2. Use Template Mode first - no API key needed!
3. Always run `go mod tidy` after generating code
4. Auto-setup runs `go mod download` automatically, but verify it worked
5. AI Mode is optional - Template Mode covers most use cases
6. Set environment variables before running
7. Check database is running
8. Read error messages - they usually tell you what's wrong!
:::
