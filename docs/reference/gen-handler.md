# anaphase gen handler

Generate HTTP handlers with DTOs and tests for a domain.

## Synopsis

```bash
anaphase gen handler --domain <domain-name> [flags]
```

## Description

Generates HTTP request handlers for a domain including:
- Handler struct with service dependency
- CRUD endpoints (Create, Read, Update, Delete)
- Request/Response DTOs
- Route registration
- Test scaffolding

## Required Flags

### `--domain` (string)

Domain name to generate handlers for.

Must match an existing entity in `internal/core/entity/`.

```bash
--domain customer
--domain product
--domain order
```

## Optional Flags

### `--protocol` (string)

Protocol to use for handlers.

- **Options**: `http`, `grpc`, `graphql`
- **Default**: `http`

```bash
--protocol http    # REST API (default)
--protocol grpc    # gRPC service
--protocol graphql # GraphQL resolvers
```

::: tip
Currently only HTTP is fully supported.
:::

## Examples

### Basic Usage

```bash
# Generate HTTP handlers for customer domain
anaphase gen handler --domain customer
```

**Generated files:**
```
internal/adapter/handler/http/
├── customer_handler.go       # Handler implementation
├── customer_dto.go           # Request/Response DTOs
└── customer_handler_test.go  # Test scaffolding
```

### Multiple Domains

```bash
anaphase gen handler --domain customer
anaphase gen handler --domain product
anaphase gen handler --domain order
```

## Generated Code

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

// Additional methods: GetByID, Update, Delete
```

### DTOs

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

### Tests

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

## Generated Endpoints

For domain `customer`:

| Method | Path | Description |
|--------|------|-------------|
| POST | `/customers` | Create customer |
| GET | `/customers/:id` | Get customer by ID |
| PUT | `/customers/:id` | Update customer |
| DELETE | `/customers/:id` | Delete customer |

All routes are registered under `/api/v1` by the wire command.

## Integration with Wire

After generating handlers, run wire to register routes:

```bash
anaphase gen handler --domain customer
anaphase wire
```

Wire automatically registers:
```go
func (a *App) RegisterRoutes(r chi.Router) {
    a.customerHandler.RegisterRoutes(r)
}
```

## Customization

Generated handlers are starting points. Common customizations:

### Add Validation

```go
import "github.com/go-playground/validator/v10"

var validate = validator.New()

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req CreateCustomerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusBadRequest, "invalid request")
        return
    }

    // Add validation
    if err := validate.Struct(req); err != nil {
        h.respondError(w, http.StatusBadRequest, err.Error())
        return
    }

    // ...
}
```

### Add Authentication

```go
func (h *CustomerHandler) RegisterRoutes(r chi.Router) {
    r.Group(func(r chi.Router) {
        r.Use(AuthMiddleware) // Add auth middleware

        r.Post("/customers", h.Create)
        r.Get("/customers/{id}", h.GetByID)
        r.Put("/customers/{id}", h.Update)
        r.Delete("/customers/{id}", h.Delete)
    })
}
```

### Add Pagination

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

## See Also

- [gen domain](/reference/gen-domain)
- [gen repository](/reference/gen-repository)
- [wire](/reference/wire)
- [Examples](/examples/basic)
