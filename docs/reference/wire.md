# anaphase wire

Auto-wire dependencies and generate application entry point.

## Synopsis

```bash
anaphase wire [flags]
```

## Description

Automatically discovers all domains in your project and generates:

- **main.go**: HTTP server with graceful shutdown
- **wire.go**: Dependency injection code

Uses AST (Abstract Syntax Tree) analysis to scan your codebase and detect entities, then wires all components together.

## How It Works

### 1. Domain Discovery

Scans `internal/core/entity/` directory:

```go
// Finds all struct declarations
type Customer struct { ... }  // Discovered: "customer"
type Product struct { ... }   // Discovered: "product"
type Order struct { ... }     // Discovered: "order"
```

### 2. Code Generation

Generates wiring code for each discovered domain:

```go
// App struct
type App struct {
    logger          *slog.Logger
    db              *pgxpool.Pool
    customerHandler *handlerhttp.CustomerHandler
    productHandler  *handlerhttp.ProductHandler
    orderHandler    *handlerhttp.OrderHandler
}

// InitializeApp function
func InitializeApp(logger *slog.Logger) (*App, error) {
    // Database connection
    db, err := pgxpool.New(context.Background(), dbURL)

    // Initialize each domain
    customerRepo := postgres.NewCustomerRepository(db)
    customerHandler := handlerhttp.NewCustomerHandler(nil, logger)

    productRepo := postgres.NewProductRepository(db)
    productHandler := handlerhttp.NewProductHandler(nil, logger)

    // ... etc

    return &App{
        logger:          logger,
        db:              db,
        customerHandler: customerHandler,
        productHandler:  productHandler,
    }, nil
}
```

### 3. Route Registration

Generates route registration:

```go
func (a *App) RegisterRoutes(r chi.Router) {
    a.customerHandler.RegisterRoutes(r)
    a.productHandler.RegisterRoutes(r)
    a.orderHandler.RegisterRoutes(r)
}
```

## Flags

### `--output` (string)

Output directory for generated files.

- **Default**: `cmd/api`

```bash
anaphase wire --output cmd/server
```

## Generated Files

### main.go

Complete HTTP server with:

- Logger setup (JSON structured logging)
- Context with cancellation
- Database connection
- Router setup (Chi)
- Middleware (logger, recoverer, request ID, timeout)
- Health check endpoint
- Graceful shutdown

```go
package main

import (
    "context"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)

func main() {
    // Setup logger
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
    slog.SetDefault(logger)

    // Create context with cancellation
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer cancel()

    // Initialize dependencies
    app, err := InitializeApp(logger)
    if err != nil {
        logger.Error("failed to initialize app", "error", err)
        os.Exit(1)
    }
    defer app.Cleanup()

    // Setup router
    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.RequestID)
    r.Use(middleware.Timeout(60 * time.Second))

    // Health check
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    // API routes
    r.Route("/api/v1", func(r chi.Router) {
        app.RegisterRoutes(r)
    })

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    srv := &http.Server{
        Addr:    ":" + port,
        Handler: r,
    }

    // Start server in goroutine
    go func() {
        logger.Info("starting server", "port", port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Error("server error", "error", err)
        }
    }()

    // Wait for interrupt signal
    <-ctx.Done()
    logger.Info("shutting down gracefully...")

    // Graceful shutdown
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer shutdownCancel()

    if err := srv.Shutdown(shutdownCtx); err != nil {
        logger.Error("shutdown error", "error", err)
    }

    logger.Info("server stopped")
}
```

### wire.go

Dependency injection with:

- App struct holding all dependencies
- InitializeApp function
- RegisterRoutes method
- Cleanup method

```go
package main

import (
    "context"
    "fmt"
    "log/slog"
    "os"

    "github.com/go-chi/chi/v5"
    "github.com/jackc/pgx/v5/pgxpool"

    handlerhttp "github.com/lisvindanu/anaphase-cli/internal/adapter/handler/http"
    "github.com/lisvindanu/anaphase-cli/internal/adapter/repository/postgres"
)

// App holds all application dependencies
type App struct {
    logger *slog.Logger
    db     *pgxpool.Pool

    customerHandler *handlerhttp.CustomerHandler
}

// InitializeApp initializes all application dependencies
func InitializeApp(logger *slog.Logger) (*App, error) {
    // Database connection
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        dbURL = "postgres://postgres:postgres@localhost:5432/anaphase?sslmode=disable"
    }

    db, err := pgxpool.New(context.Background(), dbURL)
    if err != nil {
        return nil, fmt.Errorf("connect to database: %w", err)
    }

    // Ping database
    if err := db.Ping(context.Background()); err != nil {
        return nil, fmt.Errorf("ping database: %w", err)
    }

    logger.Info("database connected")

    // Initialize customer dependencies
    customerRepo := postgres.NewCustomerRepository(db)
    _ = customerRepo // TODO: Pass to service when implemented
    // TODO: Create customer service implementation
    // customerService := service.NewCustomerService(customerRepo)
    customerHandler := handlerhttp.NewCustomerHandler(nil, logger)

    return &App{
        logger:          logger,
        db:              db,
        customerHandler: customerHandler,
    }, nil
}

// RegisterRoutes registers all HTTP routes
func (a *App) RegisterRoutes(r chi.Router) {
    a.customerHandler.RegisterRoutes(r)
}

// Cleanup cleans up application resources
func (a *App) Cleanup() {
    if a.db != nil {
        a.db.Close()
        a.logger.Info("database connection closed")
    }
}
```

## Examples

### Basic Usage

```bash
# Generate domains
anaphase gen domain --name customer --prompt "..."
anaphase gen domain --name product --prompt "..."

# Generate infrastructure
anaphase gen handler --domain customer
anaphase gen handler --domain product
anaphase gen repository --domain customer
anaphase gen repository --domain product

# Wire everything
anaphase wire
```

### Custom Output

```bash
# Generate to custom directory
anaphase wire --output cmd/server

# Files created:
# - cmd/server/main.go
# - cmd/server/wire.go
```

### Multi-Service

```bash
# API service
anaphase wire --output cmd/api

# Worker service
anaphase wire --output cmd/worker

# Admin service
anaphase wire --output cmd/admin
```

## What Gets Wired

### Database

- PostgreSQL connection pool (pgxpool)
- Connection from `DATABASE_URL` env var
- Default: `postgres://postgres:postgres@localhost:5432/anaphase?sslmode=disable`
- Ping check on startup

### Repositories

For each domain:
```go
customerRepo := postgres.NewCustomerRepository(db)
```

::: tip Service Layer
Repository is created but service layer is left for you to implement:
```go
// TODO: Create customer service implementation
// customerService := service.NewCustomerService(customerRepo)
```
:::

### Handlers

For each domain:
```go
customerHandler := handlerhttp.NewCustomerHandler(nil, logger)
```

Nil passed for service (implement service layer to use).

### Routes

All handlers registered under `/api/v1`:

```
GET  /health                    → Health check
POST /api/v1/customers          → Create customer
GET  /api/v1/customers/:id      → Get customer
PUT  /api/v1/customers/:id      → Update customer
DELETE /api/v1/customers/:id    → Delete customer
```

## Environment Variables

The generated app uses:

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection | `postgres://...` |
| `PORT` | HTTP server port | `8080` |

## Running the App

After wiring:

```bash
# Set up database
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/anaphase"

# Apply migrations
psql $DATABASE_URL -f internal/adapter/repository/postgres/schema.sql

# Run
go run cmd/api/main.go
```

Output:
```json
{"time":"...","level":"INFO","msg":"database connected"}
{"time":"...","level":"INFO","msg":"starting server","port":"8080"}
```

## Regenerating

Safe to run multiple times:

```bash
# Add new domain
anaphase gen domain --name order --prompt "..."
anaphase gen handler --domain order

# Rewire (will include new domain)
anaphase wire
```

Existing `main.go` and `wire.go` will be overwritten.

::: warning Custom Changes
Don't manually edit generated `main.go` or `wire.go` - changes will be lost on rewire.

For customizations:
- Extend handlers
- Create middleware files
- Modify in service layer
:::

## Troubleshooting

### No Domains Found

```
discovered domains: count=0
```

**Cause**: No entities in `internal/core/entity/`

**Solution**: Generate a domain first:
```bash
anaphase gen domain --name customer --prompt "..."
```

### Import Errors

```
could not import github.com/lisvindanu/anaphase-cli/internal/adapter/handler/http
```

**Cause**: Handlers not generated yet

**Solution**: Generate handlers:
```bash
anaphase gen handler --domain customer
```

### Database Connection Failed

```
Error: connect to database: connection refused
```

**Cause**: Database not running

**Solution**: Start database:
```bash
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres:16
```

## See Also

- [gen domain](/reference/gen-domain)
- [gen handler](/reference/gen-handler)
- [gen repository](/reference/gen-repository)
- [Architecture](/guide/architecture)
