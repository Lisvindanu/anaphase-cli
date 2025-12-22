# Multi-Domain E-commerce Service

Build a complete e-commerce microservice with multiple domains using Anaphase.

## Overview

This example demonstrates how to generate and wire multiple domains together to create a full-featured e-commerce API with:

- **Customer** domain - User management and authentication
- **Product** domain - Product catalog with inventory
- **Order** domain - Order processing and fulfillment
- **Payment** domain - Payment processing integration

## Prerequisites

- Go 1.21+
- PostgreSQL (or Docker)
- Gemini API key configured

## Step 1: Initialize Project

```bash
mkdir ecommerce-api
cd ecommerce-api
anaphase init
```

## Step 2: Generate Customer Domain

```bash
anaphase gen domain \
  --name customer \
  --prompt "Customer with email (validated), full name (first and last),
            phone number in E.164 format, billing address, shipping address,
            and loyalty points balance. Email must be unique."
```

**Generated files:**
- `internal/core/entity/customer.go`
- `internal/core/valueobject/email.go`
- `internal/core/valueobject/phone.go`
- `internal/core/valueobject/address.go`
- `internal/core/port/customer_repo.go`
- `internal/core/port/customer_service.go`

## Step 3: Generate Product Domain

```bash
anaphase gen domain \
  --name product \
  --prompt "Product with SKU code (alphanumeric, 8-12 chars),
            name, description, price in USD (must be positive),
            inventory quantity (non-negative), category,
            and status (active, discontinued).
            SKU must be unique."
```

**Generated files:**
- `internal/core/entity/product.go`
- `internal/core/valueobject/sku.go`
- `internal/core/valueobject/money.go`
- `internal/core/port/product_repo.go`
- `internal/core/port/product_service.go`

## Step 4: Generate Order Domain

```bash
anaphase gen domain \
  --name order \
  --prompt "Order containing multiple line items.
            Each line item references a product ID, has quantity, and unit price.
            Order has customer reference, shipping address, total amount,
            and status (pending, confirmed, shipped, delivered, cancelled).
            Order number must be unique."
```

**Generated files:**
- `internal/core/entity/order.go`
- `internal/core/entity/line_item.go`
- `internal/core/valueobject/order_number.go`
- `internal/core/port/order_repo.go`
- `internal/core/port/order_service.go`

## Step 5: Generate Payment Domain

```bash
anaphase gen domain \
  --name payment \
  --prompt "Payment with order reference, amount, currency,
            payment method (credit_card, paypal, bank_transfer),
            status (pending, completed, failed, refunded),
            transaction ID, and timestamp."
```

**Generated files:**
- `internal/core/entity/payment.go`
- `internal/core/valueobject/transaction_id.go`
- `internal/core/port/payment_repo.go`
- `internal/core/port/payment_service.go`

## Step 6: Generate HTTP Handlers

```bash
anaphase gen handler --domain customer
anaphase gen handler --domain product
anaphase gen handler --domain order
anaphase gen handler --domain payment
```

Each generates:
- Handler implementation
- DTOs for requests/responses
- Route registration
- Test scaffolding

## Step 7: Generate Repositories

```bash
anaphase gen repository --domain customer --db postgres
anaphase gen repository --domain product --db postgres
anaphase gen repository --domain order --db postgres
anaphase gen repository --domain payment --db postgres
```

Each generates:
- Repository implementation
- SQL schema with indexes
- CRUD operations

## Step 8: Auto-wire Dependencies

```bash
anaphase wire
```

This generates `cmd/api/main.go` with all dependencies wired:

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

    // Database connection
    db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Initialize repositories
    customerRepository := customerRepo.NewCustomerRepository(db)
    productRepository := productRepo.NewProductRepository(db)
    orderRepository := orderRepo.NewOrderRepository(db)
    paymentRepository := paymentRepo.NewPaymentRepository(db)

    // Initialize handlers
    customerHTTPHandler := customerHandler.NewCustomerHandler(customerRepository, logger)
    productHTTPHandler := productHandler.NewProductHandler(productRepository, logger)
    orderHTTPHandler := orderHandler.NewOrderHandler(orderRepository, logger)
    paymentHTTPHandler := paymentHandler.NewPaymentHandler(paymentRepository, logger)

    // Setup routes
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

## Step 9: Apply Database Schema

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/ecommerce"

# Apply all schemas
psql $DATABASE_URL -f internal/adapter/repository/postgres/schema.sql
```

## Step 10: Run the Service

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/ecommerce"
export GEMINI_API_KEY="your-api-key"

go run cmd/api/main.go
```

## Generated API Endpoints

### Customer API

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/customers` | Create customer |
| GET | `/customers/:id` | Get customer |
| PUT | `/customers/:id` | Update customer |
| DELETE | `/customers/:id` | Delete customer |

### Product API

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/products` | Create product |
| GET | `/products/:id` | Get product |
| GET | `/products/sku/:sku` | Get by SKU |
| PUT | `/products/:id` | Update product |
| DELETE | `/products/:id` | Delete product |

### Order API

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/orders` | Create order |
| GET | `/orders/:id` | Get order |
| PUT | `/orders/:id/status` | Update status |
| GET | `/orders/customer/:customerId` | Customer orders |

### Payment API

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/payments` | Process payment |
| GET | `/payments/:id` | Get payment |
| GET | `/payments/order/:orderId` | Order payments |
| POST | `/payments/:id/refund` | Refund payment |

## Example Requests

### Create Customer

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

### Create Product

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

### Create Order

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

### Process Payment

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

## Project Structure

```
ecommerce-api/
├── cmd/
│   └── api/
│       ├── main.go              # Auto-generated entrypoint
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

## Next Steps

### Add Business Logic

Implement service layer with business rules:

```go
// internal/service/order_service.go
type orderService struct {
    orderRepo   port.OrderRepository
    productRepo port.ProductRepository
    paymentRepo port.PaymentRepository
}

func (s *orderService) CreateOrder(ctx context.Context, req *CreateOrderRequest) (*entity.Order, error) {
    // Validate inventory
    for _, item := range req.Items {
        product, err := s.productRepo.FindByID(ctx, item.ProductID)
        if err != nil {
            return nil, err
        }

        if product.Inventory < item.Quantity {
            return nil, ErrInsufficientInventory
        }
    }

    // Create order
    order := entity.NewOrder(req.CustomerID, req.Items)

    // Reserve inventory
    for _, item := range req.Items {
        if err := s.productRepo.UpdateInventory(ctx, item.ProductID, -item.Quantity); err != nil {
            return nil, err
        }
    }

    // Save order
    if err := s.orderRepo.Save(ctx, order); err != nil {
        return nil, err
    }

    return order, nil
}
```

### Add Tests

```bash
anaphase gen test --domain customer
anaphase gen test --domain product
anaphase gen test --domain order
anaphase gen test --domain payment
```

### Add API Documentation

```bash
anaphase gen swagger
```

### Add Metrics and Monitoring

Integrate Prometheus metrics:

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

## See Also

- [Basic Example](/examples/basic)
- [Custom Handlers](/examples/custom-handlers)
- [Architecture Guide](/guide/architecture)
