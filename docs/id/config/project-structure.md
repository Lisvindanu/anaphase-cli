# Struktur Project

Memahami layout dan organisasi project yang dihasilkan.

::: info Auto-Setup di v0.4.0
Anaphase sekarang mengotomasi setup project! Saat Anda menjalankan `anaphase init`, tool ini otomatis:
- Membuat struktur direktori lengkap
- Menghasilkan file `.env` dengan konfigurasi database
- Setup Template Mode untuk generasi kode instan
- Inisialisasi git repository (opsional)

Semuanya siap digunakan langsung!
:::

## Overview

Anaphase menghasilkan project mengikuti pola Clean Architecture dan Hexagonal Architecture:

```
my-project/
├── cmd/                      # Entry point aplikasi
│   └── api/
│       ├── main.go          # HTTP server
│       └── wire.go          # Dependency injection
├── internal/                # Kode aplikasi private
│   ├── core/               # Domain Layer (Business Logic)
│   │   ├── entity/         # Entity domain
│   │   ├── port/           # Interface (port)
│   │   └── valueobject/    # Value object
│   └── adapter/            # Infrastructure Layer
│       ├── handler/        # Adapter presentasi
│       │   └── http/
│       └── repository/     # Adapter persistence
│           └── postgres/
├── .env                     # Environment variable (auto-generated di v0.4.0)
├── .gitignore              # File git ignore
├── go.mod                   # File module Go
├── go.sum                   # Dependencies Go
└── README.md               # Dokumentasi project
```

::: tip Baru di v0.4.0
File `.env` otomatis dibuat selama `anaphase init` dengan konfigurasi database Anda. File ini berisi data sensitif dan otomatis ditambahkan ke `.gitignore`.
:::

## Rincian Direktori

### `cmd/`

Entry point aplikasi.

**Tujuan:**
- Program utama
- Perintah yang dapat dieksekusi
- Startup server

**Isi:**
```
cmd/
└── api/
    ├── main.go    # HTTP server dengan graceful shutdown
    └── wire.go    # Wiring dependency injection (auto-generated)
```

::: info Template Mode di v0.4.0
Saat menggunakan Template Mode (flag `--template`), file-file ini dihasilkan secara instan menggunakan template production-ready. Tidak perlu konfigurasi AI!
:::

**Kapan menambahkan:**
- Multiple service (api, worker, cli)
- Mode eksekusi berbeda

**Contoh:**
```
cmd/
├── api/         # REST API server
├── worker/      # Background worker
└── cli/         # CLI tool
```

### `internal/`

Kode aplikasi private (tidak bisa diimport oleh project lain).

**Tujuan:**
- Enkapsulasi implementasi
- Sembunyikan internal
- Cegah dependencies eksternal

### `internal/core/`

Domain Layer - logika dan aturan bisnis.

**Karakteristik:**
- Independen dari framework
- Tidak ada dependencies infrastructure
- Logika bisnis murni
- Kode paling penting

#### `internal/core/entity/`

Entity domain dengan identitas.

**File:**
```
entity/
├── customer.go    # Dihasilkan oleh: anaphase gen domain --name customer
├── product.go     # Dihasilkan oleh: anaphase gen domain --name product
└── order.go       # Dihasilkan oleh: anaphase gen domain --name order
```

::: tip Generasi Template Mode
Di v0.4.0, gunakan flag `--template` untuk menghasilkan file-file ini secara instan:
```bash
anaphase gen domain --name customer --template
# Membuat: internal/core/entity/customer.go
```
:::

**Isi:**
```go
// Entity dengan identitas dan lifecycle (auto-generated)
type Customer struct {
    ID        uuid.UUID
    Email     *valueobject.Email
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Method logika bisnis
func (c *Customer) UpdateEmail(email *valueobject.Email) error {
    // Validasi dan aturan bisnis
}
```

#### `internal/core/valueobject/`

Value object tanpa identitas.

**File:**
```
valueobject/
├── email.go
├── money.go
├── address.go
└── phone.go
```

**Isi:**
```go
// Value object immutable
type Email struct {
    value string
}

func NewEmail(value string) (*Email, error) {
    // Validasi
}
```

#### `internal/core/port/`

Interface yang mendefinisikan kontrak.

**File:**
```
port/
├── customer_repo.go      # Interface repository
├── customer_service.go   # Interface service
├── product_repo.go
└── product_service.go
```

**Isi:**
```go
// Interface repository (port)
type CustomerRepository interface {
    Save(ctx context.Context, c *entity.Customer) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Customer, error)
}
```

### `internal/adapter/`

Infrastructure Layer - urusan eksternal.

**Karakteristik:**
- Mengimplementasi port
- Kode spesifik framework
- Database, HTTP, dll.

#### `internal/adapter/handler/`

Adapter presentation layer.

**Struktur:**
```
handler/
└── http/
    ├── customer_handler.go
    ├── customer_dto.go
    ├── product_handler.go
    └── product_dto.go
```

**Isi:**
```go
// HTTP handler (adapter)
type CustomerHandler struct {
    service port.CustomerService
    logger  *slog.Logger
}

// Urusan spesifik HTTP
func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
    // Parse HTTP request
    // Panggil service
    // Return HTTP response
}
```

#### `internal/adapter/repository/`

Adapter persistence layer.

**Struktur:**
```
repository/
└── postgres/                # Database dipilih saat anaphase init
    ├── customer_repo.go    # Dihasilkan oleh: anaphase gen repository
    ├── product_repo.go     # Dihasilkan oleh: anaphase gen repository
    └── schema.sql          # Schema auto-generated
```

::: info Auto-Selection Database
Di v0.4.0, saat Anda menjalankan `anaphase init`, Anda memilih database (PostgreSQL, MySQL, atau MongoDB), dan semua kode repository otomatis dikonfigurasi untuk database tersebut.
:::

**Isi:**
```go
// PostgreSQL repository (adapter - auto-configured)
type customerRepository struct {
    db *pgxpool.Pool  // Koneksi dari .env DATABASE_URL
}

// Mengimplementasi port.CustomerRepository
func (r *customerRepository) Save(ctx context.Context, c *entity.Customer) error {
    // SQL query (auto-generated untuk database yang dipilih)
}
```

## Aturan Dependency

### Layer Dependencies

```
┌─────────────────┐
│   cmd/api       │  Application Layer
│   (main.go)     │  Wiring semuanya
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  adapter/       │  Infrastructure Layer
│  (handler, repo)│  Bergantung pada Core
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  core/          │  Domain Layer
│  (entity, port) │  Tidak ada dependencies
└─────────────────┘
```

**Aturan:**
- Core TIDAK memiliki dependencies
- Adapter bergantung pada Core
- Application mewiring Adapter ke Core

**Contoh:**
```go
// ✅ Baik: Adapter mengimport Core
package postgres
import "myapp/internal/core/port"

// ❌ Buruk: Core mengimport Adapter
package port
import "myapp/internal/adapter/repository/postgres"
```

## Organisasi Package

### Berdasarkan Domain (Direkomendasikan)

Organisasi berdasarkan domain bisnis:

```
internal/
├── core/
│   ├── entity/
│   │   ├── customer.go
│   │   ├── order.go
│   │   └── product.go
│   └── port/
│       ├── customer_repo.go
│       ├── order_repo.go
│       └── product_repo.go
└── adapter/
    ├── handler/
    │   └── http/
    │       ├── customer_handler.go
    │       ├── order_handler.go
    │       └── product_handler.go
    └── repository/
        └── postgres/
            ├── customer_repo.go
            ├── order_repo.go
            └── product_repo.go
```

### Berdasarkan Fitur (Alternatif)

Untuk domain kompleks:

```
internal/
├── customer/
│   ├── entity/
│   │   └── customer.go
│   ├── port/
│   │   ├── repository.go
│   │   └── service.go
│   ├── handler/
│   │   └── http_handler.go
│   └── repository/
│       └── postgres.go
├── order/
│   └── ...
└── product/
    └── ...
```

## Penamaan File

### Konvensi

**Entity:**
```
customer.go      # Singular, lowercase
product.go
order.go
```

**Repository:**
```
customer_repo.go           # Interface
postgres/customer_repo.go  # Implementasi
```

**Handler:**
```
customer_handler.go   # Handler
customer_dto.go       # DTO
customer_test.go      # Test
```

**Test:**
```
customer_test.go           # Unit test (package sama)
customer_integration_test.go  # Integration test
```

## Menambahkan Komponen Baru

### Tambah Domain Baru

::: info Template Mode di v0.4.0
Gunakan flag `--template` untuk generasi instan tanpa AI:
```bash
# Cepat, tidak perlu AI!
anaphase gen domain --name inventory --template
```

Atau gunakan AI untuk domain kustom:
```bash
# Domain kustom dengan AI
anaphase gen domain --name inventory --prompt "Inventory with SKU, quantity, location"
```
:::

```bash
# 1. Generate domain (Template Mode - instan!)
anaphase gen domain --name inventory --template

# Hasil:
internal/core/
├── entity/
│   └── inventory.go      # Dibuat instan
└── port/
    └── inventory_repo.go # Dibuat instan

# 2. Generate infrastructure
anaphase gen handler --domain inventory
anaphase gen repository --domain inventory --db postgres

# 3. Wire (opsional, bisa diotomasi)
anaphase wire
```

### Tambah Service Layer

Buat implementasi service secara manual:

```bash
mkdir -p internal/core/service
```

```go
// internal/core/service/customer_service.go
package service

type customerService struct {
    repo port.CustomerRepository
}

func NewCustomerService(repo port.CustomerRepository) port.CustomerService {
    return &customerService{repo: repo}
}

func (s *customerService) CreateCustomer(ctx context.Context, email, name string) (*entity.Customer, error) {
    // Logika bisnis
}
```

Update wire.go:
```go
customerService := service.NewCustomerService(customerRepo)
customerHandler := http.NewCustomerHandler(customerService, logger)
```

### Tambah Middleware

```bash
mkdir -p internal/adapter/middleware
```

```go
// internal/adapter/middleware/auth.go
package middleware

func Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Logika auth
        next.ServeHTTP(w, r)
    })
}
```

Gunakan di main.go:
```go
r.Use(middleware.Auth)
```

## Best Practice

### Jaga Core Tetap Murni

```go
// ✅ Baik: Core tidak mengimport infrastructure
package entity
import "github.com/google/uuid"

// ❌ Buruk: Core mengimport database
package entity
import "github.com/jackc/pgx/v5"
```

### File Kecil dan Fokus

```go
// ✅ Baik: Satu entity per file
// customer.go
type Customer struct { ... }

// ❌ Buruk: Multiple entity per file
// entities.go
type Customer struct { ... }
type Product struct { ... }
type Order struct { ... }
```

### Penamaan Konsisten

```go
// Repository
type CustomerRepository interface    // Interface
type customerRepository struct       // Implementasi

// Service
type CustomerService interface       // Interface
type customerService struct          // Implementasi
```

### Komentar Package

```go
// Package entity berisi entity domain.
//
// Entity adalah objek dengan identitas yang persisten seiring waktu.
// Mereka berisi logika bisnis dan menegakkan aturan bisnis.
package entity
```

## File Auto-Generated

::: info Baru di v0.4.0
Anaphase otomatis membuat beberapa file untuk membantu Anda memulai dengan cepat.
:::

### File `.env`

Auto-created selama `anaphase init` dengan konfigurasi database Anda:

```bash
# .env (auto-generated berdasarkan pilihan database Anda)
DATABASE_URL="postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable"

# Tambahkan environment variable Anda sendiri sesuai kebutuhan
# API_KEY=your-api-key
# JWT_SECRET=your-secret
```

**Fitur:**
- Pre-configured untuk database yang dipilih
- Otomatis ditambahkan ke `.gitignore`
- Siap dikustomisasi untuk environment Anda
- Dimuat otomatis oleh kode yang dihasilkan

### File `.gitignore`

Auto-created untuk mengecualikan file sensitif dan generated:

```
# .gitignore (auto-generated)
.env
.env.*
*.log
tmp/
dist/
```

### File Generated Template Mode

Saat menggunakan Template Mode (flag `--template`), semua file kode dihasilkan dari template production-ready:

- Entity domain (`internal/core/entity/*.go`)
- Repository port (`internal/core/port/*.go`)
- HTTP handler (`internal/adapter/handler/http/*.go`)
- Database repository (`internal/adapter/repository/*/*.go`)
- Schema (`internal/adapter/repository/*/schema.sql`)

Semua file mengikuti best practice DDD dan siap digunakan!

## Lihat Juga

- [Arsitektur](/guide/architecture)
- [Konsep DDD](/guide/ddd)
- [Quick Start](/guide/quick-start)
