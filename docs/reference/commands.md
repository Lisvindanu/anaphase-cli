# Command Reference

Complete reference for all Anaphase CLI commands.

::: info
**New in v0.4**: Launch the interactive menu by running `anaphase` (no arguments) for a visual, searchable interface to all commands. Press `Ctrl+K` to search.
:::

## Interactive Menu

### Quick Access

Launch the TUI menu for visual command access:

```bash
anaphase
```

**Features:**
- Visual command browser with icons
- Quick search with `Ctrl+K`
- Contextual help for each option
- Auto-setup prompts for missing tools
- No need to remember CLI flags

**Menu Options:**
```
⚡ Anaphase CLI v0.4.0

📋 Generation
  → Generate Domain       Generate domain models (entities, ports)
  → Generate Handler      Generate HTTP handlers with DTOs
  → Generate Repository   Generate database repositories
  → Generate Middleware   Generate HTTP middleware (auth, CORS, etc.)
  → Generate Migration    Generate database migrations

🔧 Tools
  → Wire Dependencies     Auto-wire dependencies and generate main.go
  → Quality Tools         Code quality (lint, format, validate)
  → Configuration         Manage AI providers and settings

📚 Help & Info
  → Documentation         Open documentation website
  → About                 Version and project information

  Press Ctrl+K to search, Ctrl+C to exit
```

### Search Feature

Press `Ctrl+K` in the menu to quickly find commands:

```
Search: middleware
→ Generate Middleware
→ Quality Tools
```

## Global Flags

Available for all commands:

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--verbose` | `-v` | Enable verbose output | `false` |
| `--debug` | | Enable debug logging | `false` |
| `--config` | `-c` | Config file path | `~/.anaphase/config.yaml` |

## Commands

::: tip
All commands below can also be accessed through the interactive menu. Run `anaphase` to launch it.
:::

### `anaphase init`

Initialize a new microservice project with Clean Architecture structure.

[Full Documentation →](/reference/init)

```bash
anaphase init [project-name] [flags]
```

### `anaphase gen domain`

Generate domain models (entities, value objects, ports) with or without AI.

[Full Documentation →](/reference/gen-domain)

```bash
# Interactive menu (recommended)
anaphase

# Direct CLI
anaphase gen domain "<description>" [flags]

# Interactive CLI mode
anaphase gen domain --interactive
```

::: info
**Template Mode**: Works without AI configuration using proven templates.
:::

### `anaphase gen handler`

Generate HTTP handlers with DTOs and tests.

[Full Documentation →](/reference/gen-handler)

```bash
# Interactive menu (recommended)
anaphase

# Direct CLI
anaphase gen handler --domain <domain> [flags]
```

### `anaphase gen repository`

Generate database repository implementations with schema.

[Full Documentation →](/reference/gen-repository)

```bash
# Interactive menu (recommended)
anaphase

# Direct CLI
anaphase gen repository --domain <domain> --db <database> [flags]
```

### `anaphase gen middleware`

Generate HTTP middleware (auth, CORS, rate limiting, logging).

[Full Documentation →](/reference/gen-middleware)

```bash
# Interactive menu (recommended)
anaphase

# Direct CLI
anaphase gen middleware --type <type> [flags]
```

### `anaphase gen migration`

Generate database migration files with intelligent SQL generation.

[Full Documentation →](/reference/gen-migration)

```bash
# Interactive menu (recommended)
anaphase

# Direct CLI
anaphase gen migration <name> [flags]
```

### `anaphase wire`

Auto-wire dependencies and generate main.go.

[Full Documentation →](/reference/wire)

```bash
# Interactive menu (recommended)
anaphase

# Direct CLI
anaphase wire [flags]
```

::: info
Auto-wiring happens automatically after generating handlers and repositories.
:::

### `anaphase quality`

Code quality tools (lint, format, validate).

[Full Documentation →](/reference/quality)

```bash
# Interactive menu (recommended)
anaphase

# Direct CLI
anaphase quality lint [path]
anaphase quality format [path]
anaphase quality validate
```

### `anaphase config`

Manage AI providers and configuration.

[Full Documentation →](/reference/config)

```bash
anaphase config list
anaphase config set-provider <provider>
anaphase config check
anaphase config show-providers
```

## Quick Examples

### Using Interactive Menu (Recommended)

```bash
# Launch interactive menu
anaphase

# Select options visually:
# 1. Select "Generate Domain"
# 2. Enter description: "User with email, name, and role"
# 3. Choose Template or AI mode
# 4. Follow prompts for handlers and repositories
# 5. Auto-wiring happens automatically

# Run the app
go run cmd/api/main.go
```

### Using CLI Directly

```bash
# Create project
anaphase init my-api

# Generate domain (template mode - no AI needed)
cd my-api
anaphase gen domain "User with email, name, and role (admin, user, guest)"

# Generate infrastructure
anaphase gen handler --domain user
anaphase gen repository --domain user --db postgres

# Wire and run (auto-wiring also happens automatically)
anaphase wire
go run cmd/api/main.go
```

### Multiple Domains Workflow

```bash
# Option 1: Use interactive menu for each domain
anaphase
# Select "Generate Domain" → Enter details → Repeat

# Option 2: Use CLI directly
anaphase gen domain "Product with SKU, name, price, inventory"
anaphase gen domain "Order with customer, items, total, status"
anaphase gen domain "Customer with email, name, addresses"

# Generate all infrastructure
for domain in product order customer; do
  anaphase gen handler --domain $domain
  anaphase gen repository --domain $domain --db postgres
done

# Wire everything (or it auto-wires after each generation)
anaphase wire
```

### Quick Search in Menu

```bash
# Launch menu
anaphase

# Press Ctrl+K, type "quality"
# → Quickly jump to Quality Tools

# Press Ctrl+K, type "config"
# → Quickly access Configuration
```

## Environment Variables

Commands respect these environment variables:

- `GEMINI_API_KEY` - Google Gemini API key
- `GROQ_API_KEY` - Groq API key
- `OPENAI_API_KEY` - OpenAI API key
- `ANTHROPIC_API_KEY` - Claude API key
- `ANAPHASE_CONFIG` - Config file path
- `DATABASE_URL` - Default database connection
- `LOG_LEVEL` - Logging level

::: info
**AI is Optional**: Most commands work in template mode without any API keys configured.
:::

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Configuration error |
| 3 | AI provider error |
| 4 | File system error |

## See Also

- [Quick Start](/guide/quick-start)
- [Configuration](/config/ai-providers)
- [Examples](/examples/basic)
