# Service E-commerce Multi-Domain

Bangun microservice e-commerce lengkap dengan multiple domain menggunakan Anaphase.

::: info Fitur v0.4.0
Contoh ini menampilkan workflow **Interactive Menu** dan **Template Mode** untuk pengembangan cepat. Bangun API e-commerce lengkap dalam hitungan menit tanpa konfigurasi manual!
:::

## Gambaran Umum

Contoh ini mendemonstrasikan cara generate dan wire multiple domain bersama-sama untuk membuat API e-commerce berfitur lengkap dengan:

- **Customer** domain - Manajemen user dan authentication
- **Product** domain - Katalog product dengan inventory
- **Order** domain - Pemrosesan dan pemenuhan order
- **Payment** domain - Integrasi pemrosesan payment

## Prasyarat

- Go 1.21+
- PostgreSQL (atau Docker)
- **Tidak perlu API key** untuk Template Mode (atau Gemini API key untuk AI Mode)

## Langkah 1: Inisialisasi Project

::: info Interactive Menu
Untuk pengalaman terbaik, gunakan `anaphase` tanpa argumen untuk mengakses interactive menu di seluruh tutorial ini!
:::

```bash
mkdir ecommerce-api
cd ecommerce-api
anaphase init
```

Anaphase secara otomatis:
- Membuat struktur project
- Menginisialisasi `go.mod`
- Menyiapkan `.env` dengan nilai default
- Menginstal dependensi

## Langkah 2: Generate Domain Customer (Template Mode)

::: info Template Mode vs AI Mode
**Template Mode** (direkomendasikan): Gunakan template yang sudah dibuat untuk domain umum - instan, tanpa perlu API key.
**AI Mode**: Deskripsikan kebutuhan khusus dalam bahasa natural - memerlukan Gemini API key.
:::

**Menggunakan Interactive Menu (Direkomendasikan):**
```bash
anaphase
# Select: Generate Domain
# Choose: Template Mode
# Select template: User/Customer
# Enter domain name: customer
```

**Atau menggunakan CLI:**
```bash
anaphase gen domain --name customer --template user
```

**Alternatif: AI Mode**
```bash
anaphase gen domain \
  --name customer \
  --prompt "Customer dengan email (tervalidasi), nama lengkap (depan dan belakang),
            nomor telepon dalam format E.164, alamat billing, alamat shipping,
            dan saldo poin loyalitas. Email harus unik."
```

**File yang dihasilkan:**
- `internal/core/entity/customer.go`
- `internal/core/valueobject/email.go`
- `internal/core/valueobject/phone.go`
- `internal/core/valueobject/address.go`
- `internal/core/port/customer_repo.go`
- `internal/core/port/customer_service.go`

## Langkah 3: Generate Domain Product (Template Mode)

**Menggunakan Interactive Menu:**
```bash
anaphase
# Select: Generate Domain
# Choose: Template Mode
# Select template: Product/Inventory
# Enter domain name: product
```

**Atau menggunakan CLI:**
```bash
anaphase gen domain --name product --template product
```

**File yang dihasilkan:**
- `internal/core/entity/product.go`
- `internal/core/valueobject/sku.go`
- `internal/core/valueobject/money.go`
- `internal/core/port/product_repo.go`
- `internal/core/port/product_service.go`

## Langkah 4: Generate Domain Order (Template Mode)

**Menggunakan Interactive Menu:**
```bash
anaphase
# Select: Generate Domain
# Choose: Template Mode
# Select template: Order/Transaction
# Enter domain name: order
```

**Atau menggunakan CLI:**
```bash
anaphase gen domain --name order --template order
```

**File yang dihasilkan:**
- `internal/core/entity/order.go`
- `internal/core/entity/line_item.go`
- `internal/core/valueobject/order_number.go`
- `internal/core/port/order_repo.go`
- `internal/core/port/order_service.go`

## Langkah 5: Generate Domain Payment (Template Mode)

**Menggunakan Interactive Menu:**
```bash
anaphase
# Select: Generate Domain
# Choose: Template Mode
# Select template: Payment
# Enter domain name: payment
```

**Atau menggunakan CLI:**
```bash
anaphase gen domain --name payment --template payment
```

**File yang dihasilkan:**
- `internal/core/entity/payment.go`
- `internal/core/valueobject/transaction_id.go`
- `internal/core/port/payment_repo.go`
- `internal/core/port/payment_service.go`

## Langkah 6: Generate HTTP Handler

**Menggunakan Interactive Menu (Paling Mudah):**
```bash
anaphase
# Select: Generate Handler
# Select domain: customer
# Ulangi untuk product, order, dan payment
```

**Atau menggunakan CLI:**
```bash
anaphase gen handler --domain customer
anaphase gen handler --domain product
anaphase gen handler --domain order
anaphase gen handler --domain payment
```

Setiap domain menghasilkan:
- Implementasi handler
- DTO untuk request/response
- Registrasi route
- Scaffolding test

## Langkah 7: Generate Repository

**Menggunakan Interactive Menu:**
```bash
anaphase
# Select: Generate Repository
# Select domain and database type (PostgreSQL)
# Ulangi untuk semua domain
```

**Atau menggunakan CLI:**
```bash
anaphase gen repository --domain customer --db postgres
anaphase gen repository --domain product --db postgres
anaphase gen repository --domain order --db postgres
anaphase gen repository --domain payment --db postgres
```

Setiap domain menghasilkan:
- Implementasi repository
- SQL schema dengan index
- Operasi CRUD

## Langkah 8: Auto-wire Dependencies

**Menggunakan Interactive Menu:**
```bash
anaphase
# Select: Wire Dependencies
```

**Atau menggunakan CLI:**
```bash
anaphase wire
```

::: info Auto-Install
Command wire secara otomatis menginstal dependensi yang hilang dan memvalidasi struktur project Anda!
:::

Ini menghasilkan `cmd/api/main.go` dengan semua dependencies yang sudah di-wire:

```go
package main

import (
    "context"
    "log"
    "log/slog"
    "os"

    "github.com/go-chi/chi/v5"
    "github.com/jackc/pgx/v5/pgxpool"

    customerHandler "myapp/internal/adapter/handler/http"
    productHandler "myapp/internal/adapter/handler/http"
    orderHandler "myapp/internal/adapter/handler/http"
    paymentHandler "myapp/internal/adapter/handler/http"

    customerRepo "myapp/internal/adapter/repository/postgres"
    productRepo "myapp/internal/adapter/repository/postgres"
    orderRepo "myapp/internal/adapter/repository/postgres"
    paymentRepo "myapp/internal/adapter/repository/postgres"
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    // Koneksi database
    db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Inisialisasi repository
    customerRepository := customerRepo.NewCustomerRepository(db)
    productRepository := productRepo.NewProductRepository(db)
    orderRepository := orderRepo.NewOrderRepository(db)
    paymentRepository := paymentRepo.NewPaymentRepository(db)

    // Inisialisasi handler
    customerHTTPHandler := customerHandler.NewCustomerHandler(customerRepository, logger)
    productHTTPHandler := productHandler.NewProductHandler(productRepository, logger)
    orderHTTPHandler := orderHandler.NewOrderHandler(orderRepository, logger)
    paymentHTTPHandler := paymentHandler.NewPaymentHandler(paymentRepository, logger)

    // Setup route
    r := chi.NewRouter()
    customerHTTPHandler.RegisterRoutes(r)
    productHTTPHandler.RegisterRoutes(r)
    orderHTTPHandler.RegisterRoutes(r)
    paymentHTTPHandler.RegisterRoutes(r)

    // Start server
    log.Println("Server starting on :8080")
    http.ListenAndServe(":8080", r)
}
```

## Langkah 9: Terapkan Database Schema

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/ecommerce"

# Terapkan semua schema
psql $DATABASE_URL -f internal/adapter/repository/postgres/schema.sql
```

## Langkah 10: Jalankan Service

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/ecommerce"
export GEMINI_API_KEY="your-api-key"

go run cmd/api/main.go
```

## Endpoint API yang Dihasilkan

### Customer API

| Method | Endpoint | Deskripsi |
|--------|----------|-------------|
| POST | `/customers` | Buat customer |
| GET | `/customers/:id` | Dapatkan customer |
| PUT | `/customers/:id` | Update customer |
| DELETE | `/customers/:id` | Hapus customer |

### Product API

| Method | Endpoint | Deskripsi |
|--------|----------|-------------|
| POST | `/products` | Buat product |
| GET | `/products/:id` | Dapatkan product |
| GET | `/products/sku/:sku` | Dapatkan berdasarkan SKU |
| PUT | `/products/:id` | Update product |
| DELETE | `/products/:id` | Hapus product |

### Order API

| Method | Endpoint | Deskripsi |
|--------|----------|-------------|
| POST | `/orders` | Buat order |
| GET | `/orders/:id` | Dapatkan order |
| PUT | `/orders/:id/status` | Update status |
| GET | `/orders/customer/:customerId` | Order customer |

### Payment API

| Method | Endpoint | Deskripsi |
|--------|----------|-------------|
| POST | `/payments` | Proses payment |
| GET | `/payments/:id` | Dapatkan payment |
| GET | `/payments/order/:orderId` | Payment order |
| POST | `/payments/:id/refund` | Refund payment |

## Contoh Request

### Buat Customer

```bash
curl -X POST http://localhost:8080/customers \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "phone": "+12125551234",
    "billingAddress": {
      "street": "123 Main St",
      "city": "New York",
      "state": "NY",
      "zipCode": "10001",
      "country": "USA"
    }
  }'
```

### Buat Product

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "sku": "PROD123456",
    "name": "Wireless Headphones",
    "description": "High-quality wireless headphones",
    "price": 99.99,
    "currency": "USD",
    "inventory": 100,
    "category": "Electronics",
    "status": "active"
  }'
```

### Buat Order

```bash
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customerId": "550e8400-e29b-41d4-a716-446655440000",
    "items": [
      {
        "productId": "660e8400-e29b-41d4-a716-446655440000",
        "quantity": 2,
        "unitPrice": 99.99
      }
    ],
    "shippingAddress": {
      "street": "123 Main St",
      "city": "New York",
      "state": "NY",
      "zipCode": "10001",
      "country": "USA"
    }
  }'
```

### Proses Payment

```bash
curl -X POST http://localhost:8080/payments \
  -H "Content-Type: application/json" \
  -d '{
    "orderId": "770e8400-e29b-41d4-a716-446655440000",
    "amount": 199.98,
    "currency": "USD",
    "paymentMethod": "credit_card",
    "cardToken": "tok_visa"
  }'
```

## Struktur Project

```
ecommerce-api/
├── cmd/
│   └── api/
│       ├── main.go              # Entrypoint yang di-generate otomatis
│       └── wire.go              # Dependency wiring
├── internal/
│   ├── core/
│   │   ├── entity/
│   │   │   ├── customer.go
│   │   │   ├── product.go
│   │   │   ├── order.go
│   │   │   ├── line_item.go
│   │   │   └── payment.go
│   │   ├── valueobject/
│   │   │   ├── email.go
│   │   │   ├── phone.go
│   │   │   ├── address.go
│   │   │   ├── sku.go
│   │   │   ├── money.go
│   │   │   ├── order_number.go
│   │   │   └── transaction_id.go
│   │   └── port/
│   │       ├── customer_repo.go
│   │       ├── customer_service.go
│   │       ├── product_repo.go
│   │       ├── product_service.go
│   │       ├── order_repo.go
│   │       ├── order_service.go
│   │       ├── payment_repo.go
│   │       └── payment_service.go
│   └── adapter/
│       ├── handler/
│       │   └── http/
│       │       ├── customer_handler.go
│       │       ├── customer_dto.go
│       │       ├── product_handler.go
│       │       ├── product_dto.go
│       │       ├── order_handler.go
│       │       ├── order_dto.go
│       │       ├── payment_handler.go
│       │       └── payment_dto.go
│       └── repository/
│           └── postgres/
│               ├── customer_repo.go
│               ├── product_repo.go
│               ├── order_repo.go
│               ├── payment_repo.go
│               └── schema.sql
├── go.mod
└── go.sum
```

## Ringkasan Quick Start

::: info Workflow Lengkap
Seluruh service multi-domain dapat dibangun hanya menggunakan interactive menu:
```bash
anaphase  # Jalankan untuk setiap langkah:
# 1. Initialize Project
# 2. Generate Domain (customer, product, order, payment)
# 3. Generate Handlers (untuk setiap domain)
# 4. Generate Repositories (untuk setiap domain)
# 5. Wire Dependencies
```
Total waktu: ~5 menit menggunakan Template Mode!
:::

## Langkah Selanjutnya

### Tambahkan Business Logic

Implementasi service layer dengan business rule:

```go
// internal/service/order_service.go
type orderService struct {
    orderRepo   port.OrderRepository
    productRepo port.ProductRepository
    paymentRepo port.PaymentRepository
}

func (s *orderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*entity.Order, error) {
    // Validasi inventory
    for _, item := range req.Items {
        product, err := s.productRepo.FindByID(ctx, item.ProductID)
        if err != nil {
            return nil, err
        }

        if product.Inventory < item.Quantity {
            return nil, ErrInsufficientInventory
        }
    }

    // Buat order
    order := entity.NewOrder(req.CustomerID, req.Items)

    // Reserve inventory
    for _, item := range req.Items {
        if err := s.productRepo.UpdateInventory(ctx, item.ProductID, -item.Quantity); err != nil {
            return nil, err
        }
    }

    // Simpan order
    if err := s.orderRepo.Save(ctx, order); err != nil {
        return nil, err
    }

    return order, nil
}
```

### Tambahkan Test

**Menggunakan Interactive Menu:**
```bash
anaphase
# Select: Generate Tests
# Select domains to test
```

**Atau menggunakan CLI:**
```bash
anaphase gen test --domain customer
anaphase gen test --domain product
anaphase gen test --domain order
anaphase gen test --domain payment
```

### Tambahkan Dokumentasi API

**Menggunakan Interactive Menu:**
```bash
anaphase
# Select: Generate Documentation
# Choose: Swagger/OpenAPI
```

**Atau menggunakan CLI:**
```bash
anaphase gen swagger
```

### Tambahkan Metrics dan Monitoring

Integrasikan Prometheus metrics:

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    ordersCreated = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "orders_created_total",
            Help: "Total number of orders created",
        },
    )
)
```

::: info Pro Tip
Gunakan Template Mode untuk domain standar dan AI Mode hanya ketika Anda memerlukan business logic yang sangat khusus. Ini menghemat waktu dan biaya API!
:::

## Lihat Juga

- [Contoh Dasar](/examples/basic)
- [Custom Handler](/examples/custom-handlers)
- [Panduan Arsitektur](/guide/architecture)
