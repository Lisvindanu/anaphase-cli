# Architecture

Anaphase generates code following Clean Architecture, Hexagonal Architecture, and Domain-Driven Design (DDD) principles.

## Overview

The generated project structure enforces separation of concerns:

```
my-app/
├── cmd/
│   └── api/                    # Application entry points
│       ├── main.go            # HTTP server
│       └── wire.go            # Dependency injection
├── internal/
│   ├── core/                  # Business logic (Domain Layer)
│   │   ├── entity/           # Domain entities
│   │   ├── valueobject/      # Value objects
│   │   └── port/             # Interfaces (ports)
│   └── adapter/              # External concerns (Infrastructure)
│       ├── handler/          # HTTP handlers
│       │   └── http/
│       └── repository/       # Data persistence
│           └── postgres/
├── go.mod
└── go.sum
```

## Layered Architecture

### 1. Domain Layer (`internal/core/`)

The heart of your application. Contains business logic and is independent of any framework or infrastructure.

#### Entities (`entity/`)

Objects with identity that persist over time:

```go
// internal/core/entity/customer.go
package entity

type Customer struct {
    ID        uuid.UUID
    Email     *valueobject.Email
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Business logic methods
func (c *Customer) UpdateEmail(email *valueobject.Email) error {
    if email == nil {
        return ErrInvalidEmail
    }
    c.Email = email
    c.UpdatedAt = time.Now()
    return nil
}
```

**Key characteristics:**
- Has unique identity (ID)
- Contains business rules
- Mutable state
- Lifecycle tracked (CreatedAt, UpdatedAt)

#### Value Objects (`valueobject/`)

Immutable objects without identity, defined by their attributes:

```go
// internal/core/valueobject/email.go
package valueobject

type Email struct {
    value string
}

func NewEmail(value string) (*Email, error) {
    if !isValidEmail(value) {
        return nil, ErrInvalidEmail
    }
    return &Email{value: strings.ToLower(value)}, nil
}

func (e *Email) String() string {
    return e.value
}
```

**Key characteristics:**
- No identity, compared by value
- Immutable
- Self-validating
- Can be shared

#### Ports (`port/`)

Interfaces defining contracts between layers:

```go
// internal/core/port/customer_repo.go
package port

type CustomerRepository interface {
    Save(ctx context.Context, customer *entity.Customer) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Customer, error)
    FindByEmail(ctx context.Context, email valueobject.Email) (*entity.Customer, error)
}

type CustomerService interface {
    CreateCustomer(ctx context.Context, email, name string) (*entity.Customer, error)
    GetCustomer(ctx context.Context, id uuid.UUID) (*entity.Customer, error)
}
```

**Benefits:**
- Dependency inversion
- Easy testing with mocks
- Swap implementations

### 2. Adapter Layer (`internal/adapter/`)

Implements the ports and handles external concerns.

#### Handlers (`handler/http/`)

HTTP request/response handling:

```go
// internal/adapter/handler/http/customer_handler.go
package http

type CustomerHandler struct {
    service port.CustomerService
    logger  *slog.Logger
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req CreateCustomerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusBadRequest, "invalid request")
        return
    }

    customer, err := h.service.CreateCustomer(r.Context(), req.Email, req.Name)
    if err != nil {
        h.respondError(w, http.StatusInternalServerError, err.Error())
        return
    }

    h.respondJSON(w, http.StatusCreated, customer)
}
```

**Responsibilities:**
- Parse HTTP requests
- Validate input
- Call service layer
- Format responses
- Handle HTTP errors

#### Repositories (`repository/postgres/`)

Database implementations:

```go
// internal/adapter/repository/postgres/customer_repo.go
package postgres

type customerRepository struct {
    db *pgxpool.Pool
}

func (r *customerRepository) Save(ctx context.Context, c *entity.Customer) error {
    query := `
        INSERT INTO customers (id, email, name, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (id) DO UPDATE
        SET email = $2, name = $3, updated_at = $5
    `
    _, err := r.db.Exec(ctx, query, c.ID, c.Email.String(), c.Name, c.CreatedAt, c.UpdatedAt)
    return err
}
```

**Responsibilities:**
- SQL queries
- Data mapping (DB ↔ Entity)
- Transaction management
- Error handling

### 3. Application Layer (`cmd/api/`)

Wires everything together and starts the application.

#### Main (`main.go`)

Application entry point:

```go
func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    app, err := InitializeApp(logger)
    if err != nil {
        logger.Error("failed to initialize", "error", err)
        os.Exit(1)
    }
    defer app.Cleanup()

    // Setup router
    r := chi.NewRouter()
    r.Use(middleware.Logger)
    app.RegisterRoutes(r)

    // Start server
    srv := &http.Server{Addr: ":8080", Handler: r}
    go srv.ListenAndServe()

    // Graceful shutdown
    <-ctx.Done()
    srv.Shutdown(context.Background())
}
```

#### Wire (`wire.go`)

Dependency injection:

```go
type App struct {
    logger          *slog.Logger
    db              *pgxpool.Pool
    customerHandler *http.CustomerHandler
}

func InitializeApp(logger *slog.Logger) (*App, error) {
    // Database
    db, err := pgxpool.New(context.Background(), dbURL)

    // Repository
    customerRepo := postgres.NewCustomerRepository(db)

    // Service (TODO: implement)
    // customerService := service.NewCustomerService(customerRepo)

    // Handler
    customerHandler := http.NewCustomerHandler(nil, logger)

    return &App{
        logger:          logger,
        db:              db,
        customerHandler: customerHandler,
    }, nil
}
```

## Design Patterns

### Repository Pattern

Abstracts data access:

```
┌──────────────┐
│   Service    │
└──────┬───────┘
       │ uses
       ▼
┌──────────────────┐
│ Repository Port  │ (interface)
└──────┬───────────┘
       │ implements
       ▼
┌──────────────────┐
│ Postgres Adapter │
└──────────────────┘
```

**Benefits:**
- Swap databases easily
- Test without database
- Centralize data access logic

### Dependency Injection

All dependencies injected via constructors:

```go
// Handler depends on Service
func NewCustomerHandler(service port.CustomerService, logger *slog.Logger) *CustomerHandler {
    return &CustomerHandler{
        service: service,
        logger:  logger,
    }
}

// Service depends on Repository
func NewCustomerService(repo port.CustomerRepository) *CustomerService {
    return &CustomerService{
        repo: repo,
    }
}
```

**Benefits:**
- Explicit dependencies
- Easy testing
- Loose coupling

### Factory Pattern

Constructors validate and construct objects:

```go
func NewEmail(value string) (*Email, error) {
    if !isValidEmail(value) {
        return nil, ErrInvalidEmail
    }
    return &Email{value: value}, nil
}

func NewCustomer(email *Email, name string) (*Customer, error) {
    if email == nil || name == "" {
        return nil, ErrInvalidInput
    }
    return &Customer{
        ID:        uuid.New(),
        Email:     email,
        Name:      name,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }, nil
}
```

## SOLID Principles

### Single Responsibility

Each component has one reason to change:

- **Entity**: Business rules for domain object
- **Repository**: Data persistence
- **Handler**: HTTP concerns
- **Service**: Business logic orchestration

### Open/Closed

Open for extension, closed for modification:

```go
// Add new handler without modifying existing code
type OrderHandler struct { /* ... */ }

func (a *App) RegisterRoutes(r chi.Router) {
    a.customerHandler.RegisterRoutes(r)
    a.orderHandler.RegisterRoutes(r)  // Added without changing Customer
}
```

### Liskov Substitution

Interfaces can be swapped:

```go
// Can swap implementations
var repo port.CustomerRepository

repo = postgres.NewCustomerRepository(db)  // Production
repo = mock.NewCustomerRepository()        // Testing
repo = redis.NewCustomerRepository(cache)  // Caching
```

### Interface Segregation

Small, focused interfaces:

```go
// Not this:
type Repository interface {
    SaveCustomer(...)
    FindCustomer(...)
    SaveOrder(...)
    FindOrder(...)
}

// But this:
type CustomerRepository interface {
    Save(...)
    FindByID(...)
}

type OrderRepository interface {
    Save(...)
    FindByID(...)
}
```

### Dependency Inversion

Depend on abstractions, not concretions:

```go
// Service depends on interface, not concrete implementation
type CustomerService struct {
    repo port.CustomerRepository  // Interface, not *postgres.CustomerRepository
}
```

## Data Flow

### Request Flow

```
HTTP Request
    ↓
Handler (Adapter)
    ├─ Parse request
    ├─ Validate input
    └─ Call service
        ↓
Service (Domain)
    ├─ Business logic
    ├─ Call repository
    └─ Return entity
        ↓
Repository (Adapter)
    ├─ Execute SQL
    ├─ Map to entity
    └─ Return result
        ↓
Handler
    ├─ Format response
    └─ Send HTTP response
```

### Error Flow

```
Repository Error
    ↓
Service catches and wraps
    ↓
Handler catches and converts to HTTP error
    ↓
Client receives structured error response
```

## Testing Strategy

### Unit Tests

Test each layer independently:

```go
// Entity tests - pure logic
func TestCustomer_UpdateEmail(t *testing.T) {
    customer := &Customer{Email: oldEmail}
    err := customer.UpdateEmail(newEmail)
    assert.NoError(t, err)
    assert.Equal(t, newEmail, customer.Email)
}

// Service tests - with mock repository
func TestCustomerService_Create(t *testing.T) {
    mockRepo := &MockRepository{}
    service := NewCustomerService(mockRepo)
    // ...
}

// Handler tests - with mock service
func TestCustomerHandler_Create(t *testing.T) {
    mockService := &MockService{}
    handler := NewCustomerHandler(mockService, logger)
    // ...
}
```

### Integration Tests

Test with real database:

```go
func TestCustomerRepository_Save(t *testing.T) {
    db := setupTestDB(t)
    repo := postgres.NewCustomerRepository(db)

    customer := &entity.Customer{/* ... */}
    err := repo.Save(context.Background(), customer)
    assert.NoError(t, err)

    // Verify in DB
    found, err := repo.FindByID(context.Background(), customer.ID)
    assert.Equal(t, customer, found)
}
```

## Next Steps

- [AI-Powered Generation](/guide/ai-generation) - How AI generates this structure
- [DDD Concepts](/guide/ddd) - Deep dive into Domain-Driven Design
- [Command Reference](/reference/commands) - CLI commands
