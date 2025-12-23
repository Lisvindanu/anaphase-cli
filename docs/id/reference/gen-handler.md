# anaphase gen handler

Generate HTTP handler dengan DTO dan test untuk domain.

::: info
**Quick Start**: Jalankan `anaphase` (tanpa argumen) untuk mengakses menu interaktif di mana Anda dapat memilih "Generate Handler" dengan interface visual.
:::

## Synopsis

```bash
anaphase gen handler --domain <domain-name> [flags]
```

## Deskripsi

Generate HTTP request handler untuk domain termasuk:
- Struct handler dengan service dependency
- Endpoint CRUD (Create, Read, Update, Delete)
- Request/Response DTO
- Registrasi route
- Scaffolding test

::: info
**Tidak Perlu AI**: Generasi handler menggunakan template untuk membuat kode yang bersih dan berfungsi secara instan.
:::

## Flag yang Diperlukan

### `--domain` (string)

Nama domain untuk generate handler.

Harus sesuai dengan entity yang ada di `internal/core/entity/`.

```bash
--domain customer
--domain product
--domain order
```

## Flag Opsional

### `--protocol` (string)

Protocol yang digunakan untuk handler.

- **Opsi**: `http`, `grpc`, `graphql`
- **Default**: `http`

```bash
--protocol http    # REST API (default)
--protocol grpc    # gRPC service
--protocol graphql # GraphQL resolvers
```

::: tip
Saat ini hanya HTTP yang fully supported.
:::

## Contoh

### Menu Interaktif (Disarankan)

```bash
# Luncurkan menu interaktif
anaphase

# Navigasi ke "Generate Handler" dan ikuti prompt:
# - Pilih domain dari entity yang tersedia
# - Pilih protocol (HTTP, gRPC, GraphQL)
# - Review dan konfirmasi
```

### Penggunaan Dasar (CLI)

```bash
# Generate HTTP handler untuk domain customer
anaphase gen handler --domain customer
```

**File yang dihasilkan:**
```
internal/adapter/handler/http/
├── customer_handler.go       # Implementasi handler
├── customer_dto.go           # Request/Response DTO
└── customer_handler_test.go  # Scaffolding test
```

::: info
**Generate Instan**: Tidak perlu setup AI. Handler dihasilkan dari template secara langsung.
:::

### Multiple Domain

```bash
anaphase gen handler --domain customer
anaphase gen handler --domain product
anaphase gen handler --domain order
```

## Kode yang Dihasilkan

### Handler

`internal/adapter/handler/http/customer_handler.go`:

```go
package http

import (
    "encoding/json"
    "log/slog"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
)

type CustomerHandler struct {
    service port.CustomerService
    logger  *slog.Logger
}

func NewCustomerHandler(service port.CustomerService, logger *slog.Logger) *CustomerHandler {
    return &CustomerHandler{
        service: service,
        logger:  logger,
    }
}

func (h *CustomerHandler) RegisterRoutes(r chi.Router) {
    r.Post("/customers", h.Create)
    r.Get("/customers/{id}", h.GetByID)
    r.Put("/customers/{id}", h.Update)
    r.Delete("/customers/{id}", h.Delete)
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req CreateCustomerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusBadRequest, "invalid request")
        return
    }

    // Call service
    customer, err := h.service.CreateCustomer(r.Context(), req.Email, req.Name)
    if err != nil {
        h.respondError(w, http.StatusInternalServerError, err.Error())
        return
    }

    h.respondJSON(w, http.StatusCreated, customer)
}

// Metode tambahan: GetByID, Update, Delete
```

### DTO

`internal/adapter/handler/http/customer_dto.go`:

```go
package http

type CreateCustomerRequest struct {
    Email string `json:"email"`
    Name  string `json:"name"`
}

type UpdateCustomerRequest struct {
    Name string `json:"name"`
}

type CustomerResponse struct {
    ID        string `json:"id"`
    Email     string `json:"email"`
    Name      string `json:"name"`
    CreatedAt string `json:"createdAt"`
    UpdatedAt string `json:"updatedAt"`
}
```

### Test

`internal/adapter/handler/http/customer_handler_test.go`:

```go
package http

import (
    "testing"
)

func TestCustomerHandler_Create(t *testing.T) {
    // TODO: Implement test
}

func TestCustomerHandler_GetByID(t *testing.T) {
    // TODO: Implement test
}
```

## Endpoint yang Dihasilkan

Untuk domain `customer`:

| Method | Path | Deskripsi |
|--------|------|-------------|
| POST | `/customers` | Buat customer |
| GET | `/customers/:id` | Dapatkan customer berdasarkan ID |
| PUT | `/customers/:id` | Update customer |
| DELETE | `/customers/:id` | Hapus customer |

Semua route didaftarkan di bawah `/api/v1` oleh command wire.

## Integrasi dengan Wire

Setelah generate handler, jalankan wire untuk mendaftarkan route:

```bash
anaphase gen handler --domain customer
anaphase wire
```

Wire secara otomatis mendaftarkan:
```go
func (a *App) RegisterRoutes(r chi.Router) {
    a.customerHandler.RegisterRoutes(r)
}
```

## Kustomisasi

Handler yang dihasilkan adalah starting point. Kustomisasi umum:

### Tambahkan Validasi

```go
import "github.com/go-playground/validator/v10"

var validate = validator.New()

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req CreateCustomerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusBadRequest, "invalid request")
        return
    }

    // Tambahkan validasi
    if err := validate.Struct(req); err != nil {
        h.respondError(w, http.StatusBadRequest, err.Error())
        return
    }

    // ...
}
```

### Tambahkan Authentication

```go
func (h *CustomerHandler) RegisterRoutes(r chi.Router) {
    r.Group(func(r chi.Router) {
        r.Use(AuthMiddleware) // Tambahkan auth middleware

        r.Post("/customers", h.Create)
        r.Get("/customers/{id}", h.GetByID)
        r.Put("/customers/{id}", h.Update)
        r.Delete("/customers/{id}", h.Delete)
    })
}
```

### Tambahkan Pagination

```go
func (h *CustomerHandler) List(w http.ResponseWriter, r *http.Request) {
    page := r.URL.Query().Get("page")
    limit := r.URL.Query().Get("limit")

    customers, err := h.service.ListCustomers(r.Context(), page, limit)
    if err != nil {
        h.respondError(w, http.StatusInternalServerError, err.Error())
        return
    }

    h.respondJSON(w, http.StatusOK, customers)
}
```

## Lihat Juga

- [gen domain](/reference/gen-domain)
- [gen repository](/reference/gen-repository)
- [wire](/reference/wire)
- [Examples](/examples/basic)
