# Domain-Driven Design

Understanding DDD concepts and how Anaphase implements them.

## What is DDD?

Domain-Driven Design is an approach to software development that:

1. **Focuses on the core domain** and domain logic
2. **Uses a ubiquitous language** shared by developers and domain experts
3. **Models complex domains** through entities, value objects, and aggregates
4. **Isolates domain logic** from infrastructure concerns

## Core Building Blocks

### Entities

Objects with **identity** that persists over time.

**Characteristics:**
- Has unique identifier (ID)
- Can change state (mutable)
- Tracked lifecycle (CreatedAt, UpdatedAt)
- Contains business logic

**Example:**
```go
type Customer struct {
    ID        uuid.UUID         // Identity
    Email     *Email            // Can change
    Name      string            // Can change
    CreatedAt time.Time         // Lifecycle
    UpdatedAt time.Time         // Lifecycle
}

// Business logic
func (c *Customer) UpdateEmail(email *Email) error {
    if email == nil {
        return ErrInvalidEmail
    }
    c.Email = email
    c.UpdatedAt = time.Now()
    return nil
}
```

**When to use:**
- Object needs unique identity
- Object changes over time
- You care about which specific instance

**Examples:**
- Customer, Order, Product
- User, Invoice, Account
- Booking, Shipment, Payment

### Value Objects

Objects without identity, defined by their **attributes**.

**Characteristics:**
- No ID
- Immutable (cannot change)
- Compared by value, not identity
- Self-validating

**Example:**
```go
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

// Value objects are compared by value
func (e *Email) Equals(other *Email) bool {
    return e.value == other.value
}
```

**When to use:**
- Represents a concept or measurement
- No need for identity
- Can be shared and replaced

**Examples:**
- Email, Phone, Address
- Money, Quantity, Price
- DateRange, Coordinates, URL

### Aggregates

Cluster of entities and value objects treated as a **single unit**.

**Characteristics:**
- Has root entity (aggregate root)
- Enforces invariants
- Transaction boundary
- Consistency boundary

**Example:**
```go
// Order is the aggregate root
type Order struct {
    ID          uuid.UUID         // Root identity
    CustomerID  uuid.UUID         // External reference
    Items       []*LineItem       // Internal entities
    Total       *Money            // Derived value
    Status      OrderStatus
    CreatedAt   time.Time
}

// LineItem is part of the aggregate (no independent existence)
type LineItem struct {
    ProductID uuid.UUID
    Quantity  int
    UnitPrice *Money
}

// Business rule enforced at aggregate boundary
func (o *Order) AddItem(productID uuid.UUID, quantity int, price *Money) error {
    // Validate invariant
    if o.Status != OrderStatusPending {
        return ErrOrderNotEditable
    }

    item := &LineItem{
        ProductID: productID,
        Quantity:  quantity,
        UnitPrice: price,
    }
    o.Items = append(o.Items, item)
    o.recalculateTotal()
    return nil
}
```

**Rules:**
- Only reference by ID, not directly
- Changes go through root
- Repository only for root

**Examples:**
- Order (with LineItems)
- ShoppingCart (with CartItems)
- Invoice (with InvoiceLines)

### Repositories

Abstraction for **persistence and retrieval** of aggregates.

**Characteristics:**
- One repository per aggregate root
- Works with complete aggregates
- Interface (port) in domain layer
- Implementation in infrastructure layer

**Example:**
```go
// Port (interface) in domain layer
package port

type OrderRepository interface {
    Save(ctx context.Context, order *entity.Order) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Order, error)
    FindByCustomer(ctx context.Context, customerID uuid.UUID) ([]*entity.Order, error)
}

// Adapter (implementation) in infrastructure layer
package postgres

type orderRepository struct {
    db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) port.OrderRepository {
    return &orderRepository{db: db}
}

func (r *orderRepository) Save(ctx context.Context, order *entity.Order) error {
    // Save order and all line items in one transaction
    tx, err := r.db.Begin(ctx)
    defer tx.Rollback(ctx)

    // Save order
    _, err = tx.Exec(ctx, "INSERT INTO orders ...", order.ID, ...)

    // Save line items
    for _, item := range order.Items {
        _, err = tx.Exec(ctx, "INSERT INTO line_items ...", item.ProductID, ...)
    }

    return tx.Commit(ctx)
}
```

**Benefits:**
- Swap implementations (Postgres, MySQL, Mock)
- Test without database
- Domain doesn't depend on infrastructure

### Services

Operations that **don't naturally fit** on entities or value objects.

**Domain Services:**
- Coordinate multiple entities
- Implement business processes
- Stateless

**Example:**
```go
package service

type OrderService struct {
    orderRepo   port.OrderRepository
    productRepo port.ProductRepository
}

func (s *OrderService) PlaceOrder(ctx context.Context, customerID uuid.UUID, items []OrderItem) (*entity.Order, error) {
    // Validate products exist and have inventory
    for _, item := range items {
        product, err := s.productRepo.FindByID(ctx, item.ProductID)
        if err != nil {
            return nil, fmt.Errorf("product not found: %w", err)
        }

        if product.Quantity < item.Quantity {
            return nil, ErrInsufficientInventory
        }
    }

    // Create order
    order := entity.NewOrder(customerID)
    for _, item := range items {
        order.AddItem(item.ProductID, item.Quantity, item.Price)
    }

    // Reserve inventory
    for _, item := range items {
        product, _ := s.productRepo.FindByID(ctx, item.ProductID)
        product.ReserveQuantity(item.Quantity)
        s.productRepo.Save(ctx, product)
    }

    // Save order
    if err := s.orderRepo.Save(ctx, order); err != nil {
        return nil, err
    }

    return order, nil
}
```

## DDD Patterns in Anaphase

### Entity Generation

```bash
anaphase gen domain --name customer --prompt "Customer with email and name"
```

Generates entity with:
- Unique ID (uuid.UUID)
- Lifecycle tracking (CreatedAt, UpdatedAt)
- Constructor with validation
- Validate() method

### Value Object Detection

AI recognizes value objects:

```bash
--prompt "Customer with email (validated), billing address"
```

Generates:
- `Email` value object with validation
- `Address` value object (composite)

### Aggregate Modeling

```bash
--prompt "Order with line items. Each line item has product and quantity."
```

AI understands:
- Order is aggregate root
- LineItem is part of aggregate
- Only Order gets a repository

### Repository Interfaces

Generated repositories follow aggregate rules:

```go
// ✅ Repository for aggregate root
type OrderRepository interface {
    Save(ctx context.Context, order *entity.Order) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Order, error)
}

// ❌ No repository for LineItem (part of aggregate)
```

## Common Patterns

### Money Pattern

Always use value object for money:

```go
type Money struct {
    amount   int64  // Store in smallest unit (cents)
    currency string
}

func NewMoney(amount float64, currency string) *Money {
    return &Money{
        amount:   int64(amount * 100),
        currency: currency,
    }
}

func (m *Money) Add(other *Money) (*Money, error) {
    if m.currency != other.currency {
        return nil, ErrCurrencyMismatch
    }
    return &Money{
        amount:   m.amount + other.amount,
        currency: m.currency,
    }, nil
}
```

**Why:**
- Avoid floating-point errors
- Enforce currency matching
- Encapsulate money operations

### Enum Pattern

Use typed constants for status:

```go
type OrderStatus string

const (
    OrderStatusPending    OrderStatus = "pending"
    OrderStatusConfirmed  OrderStatus = "confirmed"
    OrderStatusShipped    OrderStatus = "shipped"
    OrderStatusDelivered  OrderStatus = "delivered"
)

func (o *Order) Confirm() error {
    if o.Status != OrderStatusPending {
        return ErrInvalidStatusTransition
    }
    o.Status = OrderStatusConfirmed
    return nil
}
```

### Factory Pattern

Use constructors for validation:

```go
func NewOrder(customerID uuid.UUID) (*Order, error) {
    if customerID == uuid.Nil {
        return nil, ErrInvalidCustomerID
    }

    return &Order{
        ID:         uuid.New(),
        CustomerID: customerID,
        Items:      []*LineItem{},
        Status:     OrderStatusPending,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }, nil
}
```

## Anti-Patterns to Avoid

### Anemic Domain Model

❌ **Don't:**
```go
type Customer struct {
    ID    uuid.UUID
    Email string
    Name  string
}

// Logic in service, not entity
func (s *CustomerService) UpdateEmail(customer *Customer, email string) {
    customer.Email = email
}
```

✅ **Do:**
```go
type Customer struct {
    ID    uuid.UUID
    Email *Email
    Name  string
}

// Logic in entity
func (c *Customer) UpdateEmail(email *Email) error {
    if email == nil {
        return ErrInvalidEmail
    }
    c.Email = email
    c.UpdatedAt = time.Now()
    return nil
}
```

### Exposing Internals

❌ **Don't:**
```go
type Order struct {
    Items []*LineItem  // Direct access
}

// External code modifies directly
order.Items = append(order.Items, newItem)
```

✅ **Do:**
```go
type Order struct {
    items []*LineItem  // Private
}

// Controlled access
func (o *Order) AddItem(item *LineItem) error {
    if o.Status != OrderStatusPending {
        return ErrOrderNotEditable
    }
    o.items = append(o.items, item)
    return nil
}

func (o *Order) Items() []*LineItem {
    return append([]*LineItem{}, o.items...)  // Return copy
}
```

### Large Aggregates

❌ **Don't:**
```go
type Customer struct {
    ID      uuid.UUID
    Orders  []*Order   // ❌ Too large
    Invoices []*Invoice // ❌ Too large
}
```

✅ **Do:**
```go
type Customer struct {
    ID    uuid.UUID
    Email *Email
    Name  string
}

type Order struct {
    CustomerID uuid.UUID  // ✅ Reference by ID
}
```

## Best Practices

### 1. Use Ubiquitous Language

Use domain terminology:

```bash
# Finance domain
"account" not "thing"
"transaction" not "record"
"balance" not "amount"

# E-commerce domain
"order" not "purchase"
"inventory" not "stock count"
"SKU" not "product code"
```

### 2. Make Invariants Explicit

Encode business rules:

```go
func (a *Account) Withdraw(amount *Money) error {
    // Invariant: balance cannot go negative
    if a.Balance.Amount < amount.Amount {
        return ErrInsufficientFunds
    }

    a.Balance = a.Balance.Subtract(amount)
    return nil
}
```

### 3. Keep Aggregates Small

Only include what must be consistent:

```go
// ✅ Good: Small aggregate
type Order struct {
    ID     uuid.UUID
    Items  []*LineItem
    Total  *Money
}

// ❌ Bad: Too large
type Order struct {
    Customer  *Customer  // Should be ID
    Products  []*Product // Should be IDs
    Warehouse *Warehouse // Should be ID
}
```

### 4. Validate at Boundaries

```go
// Constructor validates
func NewEmail(value string) (*Email, error) {
    if !isValidEmail(value) {
        return nil, ErrInvalidEmail
    }
    return &Email{value: value}, nil
}

// Can't create invalid email
email, err := NewEmail("invalid")  // Returns error
```

## Next Steps

- [Architecture](/guide/architecture) - See how DDD fits in Clean Architecture
- [AI Generation](/guide/ai-generation) - How AI generates DDD code
- [Examples](/examples/basic) - DDD in practice
