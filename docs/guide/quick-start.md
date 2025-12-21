# Quick Start

Get started with Anaphase in under 5 minutes.

## Prerequisites

- Go 1.21 or higher
- PostgreSQL (optional, for database features)
- Google Gemini API key (free tier available)

## Installation

### From Source

```bash
git clone https://github.com/lisvindanu/anaphase-cli.git
cd anaphase-cli
go install ./cmd/anaphase
```

### Using Go Install

```bash
go install github.com/lisvindanu/anaphase-cli/cmd/anaphase@latest
```

Verify installation:

```bash
anaphase --version
```

## Configure AI Provider

Set up your Google Gemini API key:

```bash
export GEMINI_API_KEY="your-api-key-here"
```

Or create a config file at `~/.anaphase/config.yaml`:

```yaml
ai:
  primary:
    type: gemini
    apiKey: your-api-key-here
    model: gemini-2.5-flash
```

::: tip Get a Free API Key
Get a free Gemini API key at [Google AI Studio](https://makersuite.google.com/app/apikey)
:::

## Create Your First Project

### Step 1: Initialize

Create a new microservice project:

```bash
anaphase init my-app
cd my-app
```

This generates a complete project structure:

```
my-app/
├── cmd/
│   └── api/
├── internal/
│   ├── core/
│   │   ├── entity/
│   │   ├── port/
│   │   └── valueobject/
│   └── adapter/
│       ├── handler/
│       └── repository/
├── go.mod
└── README.md
```

### Step 2: Generate a Domain

Use AI to generate a complete domain model:

```bash
anaphase gen domain \
  --name customer \
  --prompt "Customer with email address, full name, and phone number. Customers can place orders."
```

This creates:
- `internal/core/entity/customer.go` - Entity with business logic
- `internal/core/valueobject/email.go` - Value objects
- `internal/core/port/customer_repo.go` - Repository interface
- `internal/core/port/customer_service.go` - Service interface

### Step 3: Generate Handlers

Create HTTP handlers for your domain:

```bash
anaphase gen handler --domain customer
```

Generated files:
- `internal/adapter/handler/http/customer_handler.go` - CRUD endpoints
- `internal/adapter/handler/http/customer_dto.go` - Request/Response DTOs
- `internal/adapter/handler/http/customer_handler_test.go` - Tests

### Step 4: Generate Repository

Create database implementation:

```bash
anaphase gen repository --domain customer --db postgres
```

Generated files:
- `internal/adapter/repository/postgres/customer_repo.go` - Repository implementation
- `internal/adapter/repository/postgres/schema.sql` - Database schema
- `internal/adapter/repository/postgres/customer_repo_test.go` - Tests

### Step 5: Wire Everything

Automatically wire all dependencies:

```bash
anaphase wire
```

This generates:
- `cmd/api/main.go` - HTTP server with graceful shutdown
- `cmd/api/wire.go` - Dependency injection

### Step 6: Run

Start the database:

```bash
docker run -d \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=anaphase \
  -p 5432:5432 \
  postgres:16-alpine
```

Apply migrations:

```bash
psql -h localhost -U postgres -d anaphase -f internal/adapter/repository/postgres/schema.sql
```

Run your API:

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/anaphase?sslmode=disable"
go run cmd/api/main.go
```

## Test Your API

Your API is now running on `http://localhost:8080`. Test it:

### Create a Customer

```bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "name": "John Doe",
    "phone": "+1234567890"
  }'
```

### Get All Customers

```bash
curl http://localhost:8080/api/v1/customers
```

### Health Check

```bash
curl http://localhost:8080/health
```

## What's Next?

- Learn about [Architecture](/guide/architecture)
- Explore [AI-Powered Generation](/guide/ai-generation)
- Read the [Command Reference](/reference/commands)
- Check out [Examples](/examples/basic)

::: tip Pro Tip
Use the `--verbose` flag with any command to see detailed output:
```bash
anaphase gen domain --name product --prompt "..." --verbose
```
:::
