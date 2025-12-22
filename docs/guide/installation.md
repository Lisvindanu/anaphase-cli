# Installation

## Prerequisites

Before installing Anaphase, ensure you have:

- **Go 1.21+**: [Download Go](https://go.dev/dl/)
- **Git**: For cloning repositories
- **PostgreSQL** (Optional): For database features
  - Or Docker to run Postgres in a container

Verify Go installation:

```bash
go version
# Should output: go version go1.21.x or higher
```

## Install Anaphase

### Option 1: Quick Install (Recommended)

Use our install script that automatically configures your PATH:

::: code-group

```bash [macOS/Linux]
curl -fsSL https://raw.githubusercontent.com/lisvindanu/anaphase-cli/main/install.sh | bash
```

```powershell [Windows]
irm https://raw.githubusercontent.com/lisvindanu/anaphase-cli/main/install.ps1 | iex
```

:::

Or download and run manually:

::: code-group

```bash [macOS/Linux]
curl -fsSL https://raw.githubusercontent.com/lisvindanu/anaphase-cli/main/install.sh -o install.sh
chmod +x install.sh
./install.sh
```

```powershell [Windows]
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/lisvindanu/anaphase-cli/main/install.ps1" -OutFile "install.ps1"
powershell -ExecutionPolicy Bypass -File install.ps1
```

:::

The script will:
- Install the latest version via `go install`
- Detect your shell (bash, zsh, fish) or PowerShell on Windows
- Offer to add Go binary directory to your PATH automatically

### Option 2: Manual Install

Install directly using `go install`:

```bash
go install github.com/lisvindanu/anaphase-cli/cmd/anaphase@latest
```

Then add to your PATH:

::: code-group

```bash [Bash]
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

```bash [Zsh]
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

```bash [Fish]
set -gx PATH $HOME/go/bin $PATH
```

```powershell [Windows PowerShell]
# Temporary (current session)
$env:Path += ";$(go env GOPATH)\bin"

# Permanent (all sessions)
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$(go env GOPATH)\bin", "User")
```

:::

### Option 3: From Source

Clone and build from source for development:

```bash
git clone https://github.com/lisvindanu/anaphase-cli.git
cd anaphase-cli
go mod download
go install ./cmd/anaphase
```

## Configure AI Provider

Anaphase requires an AI provider for domain generation. Currently supported:

- **Google Gemini** (Recommended, free tier available)

### Get Gemini API Key

1. Visit [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Sign in with your Google account
3. Click "Create API Key"
4. Copy your API key

::: tip Free Tier
Google Gemini offers a generous free tier:
- 60 requests per minute
- Perfect for development and small projects
:::

### Configure API Key

#### Method 1: Environment Variable

```bash
export GEMINI_API_KEY="your-api-key-here"
```

Add to your shell profile to persist:

```bash
# ~/.bashrc or ~/.zshrc
export GEMINI_API_KEY="your-api-key-here"
```

#### Method 2: Configuration File

Create `~/.anaphase/config.yaml`:

```yaml
ai:
  primary:
    type: gemini
    apiKey: YOUR_API_KEY_HERE
    model: gemini-2.5-flash
    timeout: 30s
    retries: 3

  # Optional: fallback providers
  secondary:
    type: gemini
    apiKey: BACKUP_API_KEY
    model: gemini-2.5-flash

cache:
  enabled: true
  ttl: 24h
  dir: ~/.anaphase/cache
```

::: details Configuration Options

- `type`: AI provider type (`gemini`)
- `apiKey`: Your API key
- `model`: Model to use (`gemini-2.5-flash` recommended)
- `timeout`: Request timeout (default: `30s`)
- `retries`: Number of retries on failure (default: `3`)
- `cache.enabled`: Enable response caching (default: `true`)
- `cache.ttl`: Cache time-to-live (default: `24h`)
:::

## Database Setup (Optional)

For repository generation, you'll need a database.

### PostgreSQL with Docker

Easiest way to get started:

```bash
docker run -d \
  --name anaphase-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=anaphase \
  -p 5432:5432 \
  postgres:16-alpine
```

### Native PostgreSQL

Install PostgreSQL for your system:

::: code-group

```bash [macOS]
brew install postgresql@16
brew services start postgresql@16
createdb anaphase
```

```bash [Ubuntu/Debian]
sudo apt-get install postgresql-16
sudo systemctl start postgresql
sudo -u postgres createdb anaphase
```

```bash [Windows]
# Download from https://www.postgresql.org/download/windows/
# Or use WSL with Ubuntu instructions
```

:::

### MySQL (Alternative)

```bash
docker run -d \
  --name anaphase-mysql \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=anaphase \
  -p 3306:3306 \
  mysql:8
```

### MongoDB (Alternative)

```bash
docker run -d \
  --name anaphase-mongo \
  -p 27017:27017 \
  mongo:7
```

## Verify Installation

Test that everything works:

```bash
# Check version
anaphase --version

# Check help
anaphase --help

# Initialize a test project
mkdir test-project
cd test-project
anaphase init
```

You should see:

```
✅ Project initialized successfully!

Next steps:
  1. Configure your AI provider (see docs)
  2. Generate your first domain:
     anaphase gen domain --name user --prompt "User with email and name"
  3. Run the API:
     go run cmd/api/main.go
```

## Environment Variables

Anaphase respects these environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `GEMINI_API_KEY` | Google Gemini API key | Required |
| `DATABASE_URL` | Database connection string | `postgres://...` |
| `PORT` | HTTP server port | `8080` |
| `LOG_LEVEL` | Logging level (`debug`, `info`, `warn`, `error`) | `info` |
| `ANAPHASE_CONFIG` | Config file path | `~/.anaphase/config.yaml` |

## Troubleshooting

### Command not found

If `anaphase` is not found:

```bash
# Check if installed
ls -la $(go env GOPATH)/bin/anaphase

# Add to PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

### API Key Issues

If you see "API key not configured":

```bash
# Verify environment variable
echo $GEMINI_API_KEY

# Or check config file
cat ~/.anaphase/config.yaml
```

### Database Connection Failed

If database connection fails:

```bash
# Test connection
psql -h localhost -U postgres -d anaphase

# Check if running
docker ps | grep postgres

# View logs
docker logs anaphase-postgres
```

### Import Errors

If you see import errors after generation:

```bash
# Download dependencies
go mod download

# Tidy modules
go mod tidy
```

## Updating Anaphase

Keep Anaphase up-to-date to get the latest features and bug fixes.

### Check Current Version

```bash
anaphase --version
```

### Update to Latest Version

::: code-group

```bash [Quick Update]
# Recommended: Use install script
curl -fsSL https://raw.githubusercontent.com/lisvindanu/anaphase-cli/main/install.sh | bash
```

```bash [Manual Update]
# Using go install
go install github.com/lisvindanu/anaphase-cli/cmd/anaphase@latest
```

```bash [From Source]
# Pull latest changes
cd anaphase-cli
git pull origin main

# Rebuild
go install ./cmd/anaphase
```

```bash [Docker]
# Pull latest code
cd /var/www/anaphase-cli
git pull origin main

# Rebuild image
docker compose build

# Or pull from registry
docker pull ghcr.io/lisvindanu/anaphase-cli:latest
```

:::

### What's New

Check the [changelog](https://github.com/lisvindanu/anaphase-cli/releases) for new features:

**Recent Updates:**
- ✨ Provider Selection CLI - Choose AI provider with `--provider` flag
- ✨ Config Management - Manage providers with `anaphase config`
- ✨ Middleware Generator - Generate auth, rate limit, logging, CORS
- ✨ Interactive Mode - Guided prompts with `-i` flag
- ✨ Code Quality Tools - Lint, format, and validate code
- ✨ Migration Generator - Database migration files with smart SQL

### Verify Update

```bash
# Check new version
anaphase --version

# Test new features
anaphase config show-providers
anaphase gen middleware --help
anaphase quality --help
```

### Update Configuration

After updating, your configuration may need updates:

```bash
# Check current config
anaphase config list

# Update provider if needed
anaphase config set-provider groq

# Health check all providers
anaphase config check
```

### Rollback (If Needed)

If you need to rollback to a specific version:

```bash
# Install specific version
go install github.com/lisvindanu/anaphase-cli/cmd/anaphase@v1.0.0

# Or from source at specific tag
git checkout v1.0.0
go install ./cmd/anaphase
```

## Next Steps

- [Quick Start](/guide/quick-start) - Build your first service
- [Architecture](/guide/architecture) - Understand the patterns
- [AI Generation](/guide/ai-generation) - Learn about AI features
- [Domain-Driven Design](/guide/ddd) - **Our key differentiator**
