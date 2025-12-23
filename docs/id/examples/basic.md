# API E-commerce Dasar

Bangun API e-commerce lengkap dengan customer, product, dan order.

::: info Fitur v0.4.0
Contoh ini menggunakan **Interactive Menu** baru untuk navigasi yang lebih mudah dan **Template Mode** untuk generasi cepat tanpa AI. Anaphase sekarang mengkonfigurasi file .env Anda secara otomatis dan menginstal dependensi secara otomatis - tanpa setup manual!
:::

## Gambaran Umum

Dalam contoh ini, Anda akan membangun microservice dengan tiga domain:

- **Customer**: Akun pengguna dengan informasi kontak
- **Product**: Item untuk dijual dengan pelacakan inventori
- **Order**: Pesanan pembelian yang menghubungkan customer dan product

## Prasyarat

- Anaphase CLI terinstal
- PostgreSQL berjalan (atau Docker)
- **Tidak perlu API key** untuk Template Mode (atau Google Gemini API key untuk AI Mode)

## Langkah demi Langkah

### 1. Inisialisasi Project

::: info Interactive Menu Tersedia
Jalankan `anaphase` tanpa argumen untuk menggunakan interactive menu baru - ini cara paling mudah untuk menavigasi semua command!
:::

```bash
anaphase init ecommerce-api
cd ecommerce-api
```

Anaphase secara otomatis:
- Membuat struktur project
- Menginisialisasi `go.mod`
- Menyiapkan file `.env` dengan nilai default
- Menginstal dependensi yang diperlukan

### 2. Generate Domain Customer (Template Mode)

**Menggunakan Interactive Menu (Paling Mudah):**
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

::: info Mengapa Template Mode?
Template Mode sempurna untuk domain umum seperti user, product, dan order. Template ini menghasilkan kode production-ready secara instan tanpa perlu API key!
:::

**File yang dihasilkan:**
- `internal/core/entity/customer.go`
- `internal/core/valueobject/email.go`
- `internal/core/valueobject/phone.go`
- `internal/core/valueobject/address.go`
- `internal/core/port/customer_repo.go`
- `internal/core/port/customer_service.go`

**Alternatif: AI Mode (untuk kebutuhan khusus)**
```bash
anaphase gen domain \
  --name customer \
  --prompt "Customer dengan alamat email (tervalidasi), nama lengkap,
            nomor telepon, dan alamat pengiriman default.
            Customer dapat aktif atau tidak aktif."
```

### 3. Generate Domain Product (Template Mode)

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

### 4. Generate Domain Order (Template Mode)

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
- `internal/core/port/order_repo.go`
- `internal/core/port/order_service.go`

### 5. Generate Handler

**Menggunakan Interactive Menu:**
```bash
anaphase
# Select: Generate Handler
# Select domain: customer (kemudian ulangi untuk product dan order)
```

**Atau menggunakan CLI:**
```bash
anaphase gen handler --domain customer
anaphase gen handler --domain product
anaphase gen handler --domain order
```

**File yang dihasilkan per domain:**
- `internal/adapter/handler/http/{domain}_handler.go`
- `internal/adapter/handler/http/{domain}_dto.go`
- `internal/adapter/handler/http/{domain}_handler_test.go`

### 6. Generate Repository

**Menggunakan Interactive Menu:**
```bash
anaphase
# Select: Generate Repository
# Select domain: customer
# Select database: PostgreSQL
# (Ulangi untuk product dan order)
```

**Atau menggunakan CLI:**
```bash
anaphase gen repository --domain customer --db postgres
anaphase gen repository --domain product --db postgres
anaphase gen repository --domain order --db postgres
```

**File yang dihasilkan per domain:**
- `internal/adapter/repository/postgres/{domain}_repo.go`
- `internal/adapter/repository/postgres/schema.sql`
- `internal/adapter/repository/postgres/{domain}_repo_test.go`

### 7. Wire Dependencies

**Menggunakan Interactive Menu:**
```bash
anaphase
# Select: Wire Dependencies
```

**Atau menggunakan CLI:**
```bash
anaphase wire
```

::: info Auto-Setup
Command wire secara otomatis menginstal dependensi yang hilang dan memvalidasi struktur project Anda!
:::

**File yang dihasilkan:**
- `cmd/api/main.go`
- `cmd/api/wire.go`

### 8. Setup Database

Jalankan PostgreSQL:

```bash
docker run -d \
  --name ecommerce-db \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=ecommerce \
  -p 5432:5432 \
  postgres:16-alpine
```

Terapkan migrasi:

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/ecommerce?sslmode=disable"

# Terapkan semua schema
cat internal/adapter/repository/postgres/schema.sql | psql $DATABASE_URL
```

### 9. Jalankan API

```bash
go run cmd/api/main.go
```

Output:
```json
{"time":"2024-01-15T10:30:00Z","level":"INFO","msg":"database connected"}
{"time":"2024-01-15T10:30:00Z","level":"INFO","msg":"starting server","port":"8080"}
```

## Testing API

### Buat Customer

```bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "name": "John Doe",
    "phone": "+1234567890",
    "address": {
      "street": "123 Main St",
      "city": "New York",
      "state": "NY",
      "zipCode": "10001",
      "country": "USA"
    }
  }'
```

Response:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john@example.com",
  "name": "John Doe",
  "phone": "+1234567890",
  "address": {
    "street": "123 Main St",
    "city": "New York",
    "state": "NY",
    "zipCode": "10001",
    "country": "USA"
  },
  "createdAt": "2024-01-15T10:35:00Z",
  "updatedAt": "2024-01-15T10:35:00Z"
}
```

### Buat Product

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "sku": "LAPTOP-001",
    "name": "MacBook Pro 16\"",
    "description": "Apple MacBook Pro with M3 chip",
    "price": 2499.99,
    "quantity": 50,
    "category": "Electronics"
  }'
```

Response:
```json
{
  "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "sku": "LAPTOP-001",
  "name": "MacBook Pro 16\"",
  "description": "Apple MacBook Pro with M3 chip",
  "price": 2499.99,
  "quantity": 50,
  "category": "Electronics",
  "status": "active",
  "createdAt": "2024-01-15T10:36:00Z"
}
```

### Buat Order

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "customerId": "550e8400-e29b-41d4-a716-446655440000",
    "items": [
      {
        "productId": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
        "quantity": 1,
        "unitPrice": 2499.99
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

Response:
```json
{
  "id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "customerId": "550e8400-e29b-41d4-a716-446655440000",
  "items": [
    {
      "productId": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
      "quantity": 1,
      "unitPrice": 2499.99,
      "subtotal": 2499.99
    }
  ],
  "subtotal": 2499.99,
  "tax": 199.99,
  "total": 2699.98,
  "status": "pending",
  "createdAt": "2024-01-15T10:37:00Z"
}
```

### Lihat Semua Product

```bash
curl http://localhost:8080/api/v1/products
```

### Dapatkan Customer berdasarkan ID

```bash
curl http://localhost:8080/api/v1/customers/550e8400-e29b-41d4-a716-446655440000
```

### Update Product

```bash
curl -X PUT http://localhost:8080/api/v1/products/7c9e6679-7425-40de-944b-e07fc1f90ae7 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro 16\" (Updated)",
    "price": 2399.99,
    "quantity": 45
  }'
```

### Hapus Customer

```bash
curl -X DELETE http://localhost:8080/api/v1/customers/550e8400-e29b-41d4-a716-446655440000
```

## Struktur Project

Setelah semua generasi:

```
ecommerce-api/
├── cmd/
│   └── api/
│       ├── main.go              # HTTP server
│       └── wire.go              # Dependency injection
├── internal/
│   ├── core/
│   │   ├── entity/
│   │   │   ├── customer.go
│   │   │   ├── product.go
│   │   │   ├── order.go
│   │   │   └── line_item.go
│   │   ├── valueobject/
│   │   │   ├── email.go
│   │   │   ├── phone.go
│   │   │   ├── address.go
│   │   │   ├── sku.go
│   │   │   └── money.go
│   │   └── port/
│   │       ├── customer_repo.go
│   │       ├── customer_service.go
│   │       ├── product_repo.go
│   │       ├── product_service.go
│   │       ├── order_repo.go
│   │       └── order_service.go
│   └── adapter/
│       ├── handler/
│       │   └── http/
│       │       ├── customer_handler.go
│       │       ├── customer_dto.go
│       │       ├── product_handler.go
│       │       ├── product_dto.go
│       │       ├── order_handler.go
│       │       └── order_dto.go
│       └── repository/
│           └── postgres/
│               ├── customer_repo.go
│               ├── product_repo.go
│               ├── order_repo.go
│               └── schema.sql
├── go.mod
└── go.sum
```

## Endpoint yang Tersedia

Setelah berjalan, API Anda memiliki:

### Customer
- `POST /api/v1/customers` - Buat customer
- `GET /api/v1/customers/:id` - Dapatkan customer
- `PUT /api/v1/customers/:id` - Update customer
- `DELETE /api/v1/customers/:id` - Hapus customer

### Product
- `POST /api/v1/products` - Buat product
- `GET /api/v1/products/:id` - Dapatkan product
- `PUT /api/v1/products/:id` - Update product
- `DELETE /api/v1/products/:id` - Hapus product

### Order
- `POST /api/v1/orders` - Buat order
- `GET /api/v1/orders/:id` - Dapatkan order
- `PUT /api/v1/orders/:id` - Update status order
- `DELETE /api/v1/orders/:id` - Batalkan order

### System
- `GET /health` - Health check

## Langkah Selanjutnya

### Implementasi Business Logic

Tambahkan implementasi service layer:

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
    // Validasi
    emailVO, err := valueobject.NewEmail(email)
    if err != nil {
        return nil, err
    }

    // Buat entity
    customer, err := entity.NewCustomer(emailVO, name)
    if err != nil {
        return nil, err
    }

    // Simpan
    if err := s.repo.Save(ctx, customer); err != nil {
        return nil, err
    }

    return customer, nil
}
```

### Tambahkan Validasi

Gunakan validation library di handler:

```bash
go get github.com/go-playground/validator/v10
```

### Tambahkan Authentication

Implementasi JWT middleware:

```bash
go get github.com/golang-jwt/jwt/v5
```

### Tambahkan Test

Tulis integration test:

```go
func TestCreateCustomer(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()

    // Test repository
    repo := postgres.NewCustomerRepository(db)
    service := service.NewCustomerService(repo)

    // Test creation
    customer, err := service.CreateCustomer(context.Background(), "test@example.com", "Test User")
    assert.NoError(t, err)
    assert.NotNil(t, customer)
}
```

## Contoh Lengkap

Contoh lengkap yang berfungsi tersedia di:
https://github.com/lisvindanu/anaphase-examples/tree/main/ecommerce-api

## Lihat Juga

- [Panduan Arsitektur](/guide/architecture)
- [Panduan Cepat](/guide/quick-start)
- [Referensi Command](/reference/commands)
