# Basic E-commerce API

Build a complete e-commerce API with customers, products, and orders.

## Overview

In this example, you'll build a microservice with three domains:

- **Customer**: User accounts with contact information
- **Product**: Items for sale with inventory tracking
- **Order**: Purchase orders linking customers and products

## Prerequisites

- Anaphase CLI installed
- PostgreSQL running (or Docker)
- Google Gemini API key configured

## Step-by-Step

### 1. Initialize Project

```bash
anaphase init ecommerce-api
cd ecommerce-api
```

### 2. Generate Customer Domain

```bash
anaphase gen domain \
  --name customer \
  --prompt "Customer with email address (validated), full name,
            phone number, and default shipping address.
            Customers can be active or inactive."
```

**Generated files:**
- `internal/core/entity/customer.go`
- `internal/core/valueobject/email.go`
- `internal/core/valueobject/phone.go`
- `internal/core/valueobject/address.go`
- `internal/core/port/customer_repo.go`
- `internal/core/port/customer_service.go`

### 3. Generate Product Domain

```bash
anaphase gen domain \
  --name product \
  --prompt "Product with SKU code (unique alphanumeric),
            name, description, price in USD,
            inventory quantity (non-negative),
            and category.
            Products can be active or discontinued."
```

**Generated files:**
- `internal/core/entity/product.go`
- `internal/core/valueobject/sku.go`
- `internal/core/valueobject/money.go`
- `internal/core/port/product_repo.go`
- `internal/core/port/product_service.go`

### 4. Generate Order Domain

```bash
anaphase gen domain \
  --name order \
  --prompt "Order with customer reference,
            multiple line items each containing product reference, quantity, and unit price,
            shipping address, subtotal amount, tax amount, total amount,
            and status (pending, confirmed, processing, shipped, delivered, cancelled).
            Orders can only be cancelled if pending or confirmed."
```

**Generated files:**
- `internal/core/entity/order.go`
- `internal/core/entity/line_item.go`
- `internal/core/port/order_repo.go`
- `internal/core/port/order_service.go`

### 5. Generate Handlers

Create HTTP endpoints for each domain:

```bash
anaphase gen handler --domain customer
anaphase gen handler --domain product
anaphase gen handler --domain order
```

**Generated files per domain:**
- `internal/adapter/handler/http/{domain}_handler.go`
- `internal/adapter/handler/http/{domain}_dto.go`
- `internal/adapter/handler/http/{domain}_handler_test.go`

### 6. Generate Repositories

Create database implementations:

```bash
anaphase gen repository --domain customer --db postgres
anaphase gen repository --domain product --db postgres
anaphase gen repository --domain order --db postgres
```

**Generated files per domain:**
- `internal/adapter/repository/postgres/{domain}_repo.go`
- `internal/adapter/repository/postgres/schema.sql`
- `internal/adapter/repository/postgres/{domain}_repo_test.go`

### 7. Wire Dependencies

Auto-wire everything:

```bash
anaphase wire
```

**Generated files:**
- `cmd/api/main.go`
- `cmd/api/wire.go`

### 8. Set Up Database

Start PostgreSQL:

```bash
docker run -d \
  --name ecommerce-db \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=ecommerce \
  -p 5432:5432 \
  postgres:16-alpine
```

Apply migrations:

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/ecommerce?sslmode=disable"

# Apply all schemas
cat internal/adapter/repository/postgres/schema.sql | psql $DATABASE_URL
```

### 9. Run the API

```bash
go run cmd/api/main.go
```

Output:
```json
{"time":"2024-01-15T10:30:00Z","level":"INFO","msg":"database connected"}
{"time":"2024-01-15T10:30:00Z","level":"INFO","msg":"starting server","port":"8080"}
```

## Testing the API

### Create a Customer

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

### Create a Product

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

### Create an Order

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

### List All Products

```bash
curl http://localhost:8080/api/v1/products
```

### Get Customer by ID

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

### Delete Customer

```bash
curl -X DELETE http://localhost:8080/api/v1/customers/550e8400-e29b-41d4-a716-446655440000
```

## Project Structure

After all generation:

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

## Available Endpoints

Once running, your API has:

### Customers
- `POST /api/v1/customers` - Create customer
- `GET /api/v1/customers/:id` - Get customer
- `PUT /api/v1/customers/:id` - Update customer
- `DELETE /api/v1/customers/:id` - Delete customer

### Products
- `POST /api/v1/products` - Create product
- `GET /api/v1/products/:id` - Get product
- `PUT /api/v1/products/:id` - Update product
- `DELETE /api/v1/products/:id` - Delete product

### Orders
- `POST /api/v1/orders` - Create order
- `GET /api/v1/orders/:id` - Get order
- `PUT /api/v1/orders/:id` - Update order status
- `DELETE /api/v1/orders/:id` - Cancel order

### System
- `GET /health` - Health check

## Next Steps

### Implement Business Logic

Add service layer implementations:

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
    // Validate
    emailVO, err := valueobject.NewEmail(email)
    if err != nil {
        return nil, err
    }

    // Create entity
    customer, err := entity.NewCustomer(emailVO, name)
    if err != nil {
        return nil, err
    }

    // Save
    if err := s.repo.Save(ctx, customer); err != nil {
        return nil, err
    }

    return customer, nil
}
```

### Add Validation

Use a validation library in handlers:

```bash
go get github.com/go-playground/validator/v10
```

### Add Authentication

Implement JWT middleware:

```bash
go get github.com/golang-jwt/jwt/v5
```

### Add Tests

Write integration tests:

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

## Complete Example

Full working example available at:
https://github.com/lisvindanuu/anaphase-examples/tree/main/ecommerce-api

## See Also

- [Multi-Domain Service](/examples/multi-domain)
- [Custom Handlers](/examples/custom-handlers)
- [Architecture Guide](/guide/architecture)
