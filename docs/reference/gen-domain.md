# anaphase gen domain

Generate domain models (entities, value objects, repository ports, service ports) using AI.

## Synopsis

```bash
anaphase gen domain --name <domain-name> --prompt <description> [flags]
```

## Description

Uses AI (Google Gemini) to analyze your domain description and generate:

- **Entity**: Domain entity with fields, constructors, and validation
- **Value Objects**: Immutable objects for important concepts
- **Repository Port**: Interface for data persistence
- **Service Port**: Interface for business logic

All generated code follows Domain-Driven Design (DDD) and Clean Architecture principles.

## Required Flags

### `--name` (string)

Domain name in singular form.

```bash
# Good
--name customer
--name product
--name order

# Avoid
--name customers  # Plural
--name Customer   # Capitalized (will be normalized)
```

### `--prompt` (string)

Natural language description of your domain.

```bash
--prompt "Customer with email, name, and billing address"
```

## Optional Flags

### `--temperature` (float)

AI creativity level. Lower = more consistent, higher = more creative.

- **Range**: 0.0 to 1.0
- **Default**: 0.3
- **Recommended**: 0.1-0.3 for production

```bash
# Very consistent
--temperature 0.1

# Balanced (default)
--temperature 0.3

# More creative
--temperature 0.7
```

### `--output` (string)

Output directory for generated files.

- **Default**: Current directory
- Generated files go to `internal/core/`

```bash
--output /path/to/project
```

## Examples

### Basic Entity

```bash
anaphase gen domain \
  --name user \
  --prompt "User with email address and full name"
```

**Generates:**
```
internal/core/
├── entity/
│   └── user.go
├── valueobject/
│   └── email.go
└── port/
    ├── user_repo.go
    └── user_service.go
```

### Complex Entity

```bash
anaphase gen domain \
  --name product \
  --prompt "Product with SKU code, name, description, price in USD,
            inventory quantity, category, and status (active, discontinued).
            Products must have unique SKU."
```

**Generates:**
```go
// entity/product.go
type Product struct {
    ID          uuid.UUID
    SKU         *valueobject.SKU
    Name        string
    Description string
    Price       *valueobject.Money
    Quantity    int
    Category    string
    Status      ProductStatus
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// valueobject/sku.go
type SKU struct {
    value string
}

// valueobject/money.go
type Money struct {
    amount   int64  // cents
    currency string
}

// port/product_repo.go
type ProductRepository interface {
    Save(ctx context.Context, product *entity.Product) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
    FindBySKU(ctx context.Context, sku valueobject.SKU) (*entity.Product, error)
}
```

### With Relationships

```bash
anaphase gen domain \
  --name order \
  --prompt "Order with customer reference, multiple line items containing
            products and quantities, shipping address, total amount,
            and status (pending, confirmed, shipped, delivered)"
```

**Generates:**
- Order entity (aggregate root)
- LineItem entity (part of aggregate)
- Address value object
- Money value object
- OrderStatus enum
- OrderRepository interface

### E-commerce Example

```bash
# Customer domain
anaphase gen domain \
  --name customer \
  --prompt "Customer with email, name, phone, billing address, and shipping address"

# Product domain
anaphase gen domain \
  --name product \
  --prompt "Product with SKU, name, price, inventory, and category"

# Order domain
anaphase gen domain \
  --name order \
  --prompt "Order with customer, products with quantities, total, and status"
```

## Writing Good Prompts

### Be Specific

Good:
```bash
"Customer with validated email address, full name (first and last),
 phone number in E.164 format, and loyalty points balance"
```

Vague:
```bash
"Customer with info"
```

### Include Validation Rules

```bash
"Product with SKU (alphanumeric, 8-12 chars), price (must be positive),
 inventory (non-negative integer)"
```

### Specify Relationships

```bash
"Order containing multiple line items. Each line item references a product,
 has quantity, and unit price. Order has single customer reference."
```

### Mention Business Rules

```bash
"Account with balance. Balance cannot go negative.
 Account can be active, suspended, or closed.
 Closed accounts cannot be reactivated."
```

## AI Understanding

The AI recognizes:

### Entity vs Value Object

- **Entity**: Has ID, mutable
  - Customer, Order, Product
- **Value Object**: No ID, immutable
  - Email, Money, Address

### Field Types

- Simple: string, int, float, bool
- Time: dates, timestamps
- Complex: nested objects, arrays
- References: IDs to other entities

### Enums

```bash
"status (pending, approved, rejected)"
→ Generates ProductStatus enum
```

### Validation

```bash
"email (validated)"
→ Generates Email value object with validation

"price (positive)"
→ Adds validation in constructor
```

## Output Files

### Entity

`internal/core/entity/{name}.go`

```go
package entity

type Customer struct {
    ID        uuid.UUID
    Email     *valueobject.Email
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

func NewCustomer(email *valueobject.Email, name string) (*Customer, error) {
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

func (c *Customer) Validate() error {
    if c.Email == nil {
        return ErrInvalidEmail
    }
    if c.Name == "" {
        return ErrInvalidName
    }
    return nil
}
```

### Value Object

`internal/core/valueobject/{type}.go`

```go
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

func (e *Email) Validate() error {
    if !isValidEmail(e.value) {
        return ErrInvalidEmail
    }
    return nil
}
```

### Repository Port

`internal/core/port/{name}_repo.go`

```go
package port

type CustomerRepository interface {
    Save(ctx context.Context, customer *entity.Customer) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Customer, error)
    FindByEmail(ctx context.Context, email valueobject.Email) (*entity.Customer, error)
}
```

### Service Port

`internal/core/port/{name}_service.go`

```go
package port

type CustomerService interface {
    CreateCustomer(ctx context.Context, email, name string) (*entity.Customer, error)
    GetCustomer(ctx context.Context, id uuid.UUID) (*entity.Customer, error)
    UpdateCustomer(ctx context.Context, id uuid.UUID, name string) error
}
```

## Caching

Responses are cached by default:

```bash
# First call - hits AI API (~1-2s)
anaphase gen domain --name user --prompt "User with email"

# Second call with same prompt - uses cache (~0.1s)
anaphase gen domain --name user --prompt "User with email"
```

Cache key includes:
- Domain name
- Prompt text
- Temperature

Clear cache:
```bash
rm -rf ~/.anaphase/cache
```

## Troubleshooting

### API Quota Exceeded

```
Error: quota exceeded
```

**Solution**: Wait or use secondary API key in config.

### Invalid JSON Response

```
Error: failed to parse AI response
```

**Solution**: Try again or lower temperature to 0.1.

### Missing Value Objects

If AI doesn't generate value objects for important concepts:

```bash
# Add explicit hint
--prompt "Customer with email (value object), name, phone (value object)"
```

## See Also

- [AI-Powered Generation](/guide/ai-generation)
- [gen handler](/reference/gen-handler)
- [gen repository](/reference/gen-repository)
- [Examples](/examples/basic)
