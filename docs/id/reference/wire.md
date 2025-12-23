# anaphase wire

Auto-wire dependencies dan generate application entry point.

::: info
**Akses Cepat**: Jalankan `anaphase` (tanpa argumen) untuk mengakses menu interaktif dan pilih "Wire Dependencies" untuk interface visual.
:::

## Synopsis

```bash
anaphase wire [flags]
```

## Deskripsi

Secara otomatis menemukan semua domain di proyek Anda dan generate:

- **main.go**: HTTP server dengan graceful shutdown
- **wire.go**: Kode dependency injection

Menggunakan analisis AST (Abstract Syntax Tree) untuk scan codebase Anda dan mendeteksi entity, kemudian menyambungkan semua komponen.

::: info
**Auto-Wiring**: Di v0.4.0, wiring terjadi secara otomatis setelah generate handler dan repository. Anda jarang perlu menjalankan command ini secara manual.
:::

## Penggunaan

### Menu Interaktif (Disarankan)

```bash
anaphase
```

Pilih **"Wire Dependencies"** dari menu. Interface:
- Secara otomatis mendeteksi semua domain
- Menampilkan apa yang akan di-wire
- Memungkinkan direktori output kustom
- Menampilkan progress

### Mode CLI Langsung

```bash
anaphase wire [flags]
```

::: tip
**Biasanya Otomatis**: Setelah menjalankan `anaphase gen handler` atau `anaphase gen repository`, wiring terjadi secara otomatis. Anda hanya perlu menjalankan ini secara manual ketika menambahkan komponen kustom.
:::

## Cara Kerja

### 1. Domain Discovery

Scan direktori `internal/core/entity/`:

```go
// Menemukan semua struct declaration
type Customer struct { ... }  // Ditemukan: "customer"
type Product struct { ... }   // Ditemukan: "product"
type Order struct { ... }     // Ditemukan: "order"
```

### 2. Code Generation

Generate kode wiring untuk setiap domain yang ditemukan:

```go
// Struct App
type App struct {
    logger          *slog.Logger
    db              *pgxpool.Pool
    customerHandler *handlerhttp.CustomerHandler
    productHandler  *handlerhttp.ProductHandler
    orderHandler    *handlerhttp.OrderHandler
}

// Fungsi InitializeApp
func InitializeApp(logger *slog.Logger) (*App, error) {
    // Koneksi database
    db, err := pgxpool.New(context.Background(), dbURL)

    // Inisialisasi setiap domain
    customerRepo := postgres.NewCustomerRepository(db)
    customerHandler := handlerhttp.NewCustomerHandler(nil, logger)

    productRepo := postgres.NewProductRepository(db)
    productHandler := handlerhttp.NewProductHandler(nil, logger)

    // ... dst

    return &App{
        logger:          logger,
        db:              db,
        customerHandler: customerHandler,
        productHandler:  productHandler,
    }, nil
}
```

### 3. Registrasi Route

Generate registrasi route:

```go
func (a *App) RegisterRoutes(r chi.Router) {
    a.customerHandler.RegisterRoutes(r)
    a.productHandler.RegisterRoutes(r)
    a.orderHandler.RegisterRoutes(r)
}
```

## Flag

### `--output` (string)

Direktori output untuk file yang dihasilkan.

- **Default**: `cmd/api`

```bash
anaphase wire --output cmd/server
```

## File yang Dihasilkan

### main.go

HTTP server lengkap dengan:

- Setup logger (JSON structured logging)
- Context dengan cancellation
- Koneksi database
- Setup router (Chi)
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

    // Buat context dengan cancellation
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer cancel()

    // Inisialisasi dependencies
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

    // Start server di goroutine
    go func() {
        logger.Info("starting server", "port", port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Error("server error", "error", err)
        }
    }()

    // Tunggu interrupt signal
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

Dependency injection dengan:

- Struct App yang menyimpan semua dependency
- Fungsi InitializeApp
- Metode RegisterRoutes
- Metode Cleanup

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

// App menyimpan semua application dependency
type App struct {
    logger *slog.Logger
    db     *pgxpool.Pool

    customerHandler *handlerhttp.CustomerHandler
}

// InitializeApp menginisialisasi semua application dependency
func InitializeApp(logger *slog.Logger) (*App, error) {
    // Koneksi database
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

    // Inisialisasi customer dependency
    customerRepo := postgres.NewCustomerRepository(db)
    _ = customerRepo // TODO: Pass to service when implemented
    // TODO: Buat customer service implementation
    // customerService := service.NewCustomerService(customerRepo)
    customerHandler := handlerhttp.NewCustomerHandler(nil, logger)

    return &App{
        logger:          logger,
        db:              db,
        customerHandler: customerHandler,
    }, nil
}

// RegisterRoutes mendaftarkan semua HTTP route
func (a *App) RegisterRoutes(r chi.Router) {
    a.customerHandler.RegisterRoutes(r)
}

// Cleanup membersihkan application resource
func (a *App) Cleanup() {
    if a.db != nil {
        a.db.Close()
        a.logger.Info("database connection closed")
    }
}
```

## Contoh

### Menggunakan Menu Interaktif

```bash
# Luncurkan menu
anaphase

# Ikuti workflow:
# 1. Generate Domain (auto-wire)
# 2. Generate Handler (auto-wire)
# 3. Generate Repository (auto-wire)
# → Wiring terjadi secara otomatis!

# Atau pilih manual "Wire Dependencies" untuk rewire
```

### Menggunakan CLI (Auto-Wiring)

```bash
# Generate domain dan infrastructure
anaphase gen domain "Customer with email and name"
anaphase gen handler --domain customer
anaphase gen repository --domain customer
# → Auto-wiring terjadi setelah setiap command

anaphase gen domain "Product with SKU and price"
anaphase gen handler --domain product
anaphase gen repository --domain product
# → Auto-wiring terjadi lagi

# Sudah ter-wire! Tidak perlu menjalankan `anaphase wire`
```

### Manual Wiring (Ketika Diperlukan)

```bash
# Hanya diperlukan ketika menambahkan komponen kustom atau troubleshooting
anaphase wire
```

### Output Kustom

```bash
# Generate ke direktori kustom
anaphase wire --output cmd/server

# File yang dibuat:
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

## Yang Di-wire

### Database

- PostgreSQL connection pool (pgxpool)
- Connection dari env var `DATABASE_URL`
- Default: `postgres://postgres:postgres@localhost:5432/anaphase?sslmode=disable`
- Ping check saat startup

### Repository

Untuk setiap domain:
```go
customerRepo := postgres.NewCustomerRepository(db)
```

::: tip Service Layer
Repository dibuat tetapi service layer ditinggalkan untuk Anda implementasikan:
```go
// TODO: Buat customer service implementation
// customerService := service.NewCustomerService(customerRepo)
```
:::

### Handler

Untuk setiap domain:
```go
customerHandler := handlerhttp.NewCustomerHandler(nil, logger)
```

Nil diteruskan untuk service (implementasikan service layer untuk digunakan).

### Route

Semua handler didaftarkan di bawah `/api/v1`:

```
GET  /health                    → Health check
POST /api/v1/customers          → Buat customer
GET  /api/v1/customers/:id      → Dapatkan customer
PUT  /api/v1/customers/:id      → Update customer
DELETE /api/v1/customers/:id    → Hapus customer
```

## Environment Variable

Aplikasi yang dihasilkan menggunakan:

| Variable | Deskripsi | Default |
|----------|-------------|---------|
| `DATABASE_URL` | Koneksi PostgreSQL | `postgres://...` |
| `PORT` | Port HTTP server | `8080` |

## Menjalankan Aplikasi

Setelah wiring:

```bash
# Setup database
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/anaphase"

# Terapkan migration
psql $DATABASE_URL -f internal/adapter/repository/postgres/schema.sql

# Jalankan
go run cmd/api/main.go
```

Output:
```json
{"time":"...","level":"INFO","msg":"database connected"}
{"time":"...","level":"INFO","msg":"starting server","port":"8080"}
```

## Regenerasi

Aman untuk dijalankan berkali-kali:

```bash
# Tambahkan domain baru
anaphase gen domain "Order with items and total"
anaphase gen handler --domain order
# → Auto-wiring terjadi secara otomatis

# Atau rewire manual jika diperlukan
anaphase wire
```

`main.go` dan `wire.go` yang sudah ada akan ditimpa.

::: warning Perubahan Kustom
Jangan edit manual `main.go` atau `wire.go` yang dihasilkan - perubahan akan hilang saat rewire.

Untuk kustomisasi:
- Extend handler
- Buat file middleware (lihat `anaphase gen middleware`)
- Modifikasi di service layer
:::

::: tip Auto-Wiring di v0.4.0
Auto-wiring terjadi secara otomatis setelah generate handler dan repository. Sistem mendeteksi perubahan dan rewire dependency tanpa intervensi manual.
:::

## Troubleshooting

### Tidak Ada Domain yang Ditemukan

```
discovered domains: count=0
```

**Penyebab**: Tidak ada entity di `internal/core/entity/`

**Solusi**: Generate domain terlebih dahulu:
```bash
# Menu interaktif (disarankan)
anaphase
# Pilih "Generate Domain"

# Atau CLI
anaphase gen domain "Customer with email and name"
```

### Error Import

```
could not import github.com/lisvindanu/anaphase-cli/internal/adapter/handler/http
```

**Penyebab**: Handler belum dihasilkan

**Solusi**: Generate handler:
```bash
anaphase gen handler --domain customer
```

### Koneksi Database Gagal

```
Error: connect to database: connection refused
```

**Penyebab**: Database tidak berjalan

**Solusi**: Start database:
```bash
docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres:16
```

## Lihat Juga

- [gen domain](/reference/gen-domain)
- [gen handler](/reference/gen-handler)
- [gen repository](/reference/gen-repository)
- [Architecture](/guide/architecture)
