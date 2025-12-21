# Command Reference

Complete reference for all Anaphase CLI commands.

## Global Flags

Available for all commands:

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--verbose` | `-v` | Enable verbose output | `false` |
| `--debug` | | Enable debug logging | `false` |
| `--config` | `-c` | Config file path | `~/.anaphase/config.yaml` |

## Commands

### `anaphase init`

Initialize a new microservice project with Clean Architecture structure.

[Full Documentation →](/reference/init)

```bash
anaphase init [project-name] [flags]
```

### `anaphase gen domain`

Generate domain models (entities, value objects, ports) using AI.

[Full Documentation →](/reference/gen-domain)

```bash
anaphase gen domain --name <domain> --prompt <description> [flags]
```

### `anaphase gen handler`

Generate HTTP handlers with DTOs and tests.

[Full Documentation →](/reference/gen-handler)

```bash
anaphase gen handler --domain <domain> [flags]
```

### `anaphase gen repository`

Generate database repository implementations with schema.

[Full Documentation →](/reference/gen-repository)

```bash
anaphase gen repository --domain <domain> --db <database> [flags]
```

### `anaphase wire`

Auto-wire dependencies and generate main.go.

[Full Documentation →](/reference/wire)

```bash
anaphase wire [flags]
```

## Quick Examples

### Initialize and Generate

```bash
# Create project
anaphase init my-api

# Generate domain
cd my-api
anaphase gen domain \
  --name user \
  --prompt "User with email, name, and role (admin, user, guest)"

# Generate infrastructure
anaphase gen handler --domain user
anaphase gen repository --domain user --db postgres

# Wire and run
anaphase wire
go run cmd/api/main.go
```

### Multiple Domains

```bash
# E-commerce example
anaphase gen domain --name product --prompt "Product with SKU, name, price, inventory"
anaphase gen domain --name order --prompt "Order with customer, items, total, status"
anaphase gen domain --name customer --prompt "Customer with email, name, addresses"

# Generate all infrastructure
for domain in product order customer; do
  anaphase gen handler --domain $domain
  anaphase gen repository --domain $domain --db postgres
done

# Wire everything
anaphase wire
```

### With Custom Settings

```bash
# High verbosity
anaphase gen domain \
  --name product \
  --prompt "..." \
  --verbose \
  --debug

# Custom output directory
anaphase wire --output cmd/server

# Different database
anaphase gen repository --domain user --db mysql
```

## Environment Variables

Commands respect these environment variables:

- `GEMINI_API_KEY` - Google Gemini API key
- `ANAPHASE_CONFIG` - Config file path
- `DATABASE_URL` - Default database connection
- `LOG_LEVEL` - Logging level

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
