# anaphase gen repository

Generate database repository implementations with schema and tests.

## Synopsis

```bash
anaphase gen repository --domain <domain-name> --db <database> [flags]
```

## Description

Generates database persistence layer for a domain including:
- Repository implementation
- Database schema (SQL/migrations)
- CRUD operations
- Test scaffolding

## Required Flags

### `--domain` (string)

Domain name to generate repository for.

Must match an existing entity in `internal/core/entity/`.

```bash
--domain customer
--domain product
--domain order
```

### `--db` (string)

Database type to use.

- **Options**: `postgres`, `mysql`, `mongodb`
- **No default** (required)

```bash
--db postgres  # PostgreSQL (recommended)
--db mysql     # MySQL/MariaDB
--db mongodb   # MongoDB
```

## Optional Flags

### `--cache` (boolean)

Enable caching layer.

- **Default**: `false`

```bash
--cache  # Add Redis caching
```

::: tip
Caching support coming soon.
:::

## Examples

### PostgreSQL (Recommended)

```bash
anaphase gen repository --domain customer --db postgres
```

**Generated files:**
```
internal/adapter/repository/postgres/
├── customer_repo.go        # Repository implementation
├── schema.sql              # Database schema
└── customer_repo_test.go   # Test scaffolding
```

### MySQL

```bash
anaphase gen repository --domain product --db mysql
```

### MongoDB

```bash
anaphase gen repository --domain order --db mongodb
```

### Multiple Domains

```bash
for domain in customer product order; do
  anaphase gen repository --domain $domain --db postgres
done
```

## Generated Code

### Repository Implementation

`internal/adapter/repository/postgres/customer_repo.go`:

```go
package postgres

import (
    "context"
    "fmt"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"

    "myapp/internal/core/entity"
    "myapp/internal/core/port"
    "myapp/internal/core/valueobject"
)

type customerRepository struct {
    db *pgxpool.Pool
}

func NewCustomerRepository(db *pgxpool.Pool) port.CustomerRepository {
    return &customerRepository{
        db: db,
    }
}

func (r *customerRepository) Save(ctx context.Context, c *entity.Customer) error {
    query := `
        INSERT INTO customers (id, email, name, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (id) DO UPDATE
        SET email = $2, name = $3, updated_at = $5
    `

    _, err := r.db.Exec(ctx, query,
        c.ID,
        c.Email.String(),
        c.Name,
        c.CreatedAt,
        c.UpdatedAt,
    )

    if err != nil {
        return fmt.Errorf("save customer: %w", err)
    }

    return nil
}

func (r *customerRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Customer, error) {
    var customer entity.Customer
    var emailStr string

    query := `
        SELECT id, email, name, created_at, updated_at
        FROM customers
        WHERE id = $1
    `

    err := r.db.QueryRow(ctx, query, id).Scan(
        &customer.ID,
        &emailStr,
        &customer.Name,
        &customer.CreatedAt,
        &customer.UpdatedAt,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, fmt.Errorf("customer not found")
        }
        return nil, fmt.Errorf("find customer: %w", err)
    }

    // Convert string to value object
    email, err := valueobject.NewEmail(emailStr)
    if err != nil {
        return nil, fmt.Errorf("invalid email: %w", err)
    }
    customer.Email = email

    return &customer, nil
}

func (r *customerRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.Customer, error) {
    // Implementation
}
```

### Database Schema

`internal/adapter/repository/postgres/schema.sql`:

```sql
-- Customer table
CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_customers_email ON customers(email);
CREATE INDEX IF NOT EXISTS idx_customers_created_at ON customers(created_at);
```

### Tests

`internal/adapter/repository/postgres/customer_repo_test.go`:

```go
package postgres

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestCustomerRepository_Save(t *testing.T) {
    // TODO: Setup test database
    // TODO: Test Save method
}

func TestCustomerRepository_FindByID(t *testing.T) {
    // TODO: Test FindByID method
}
```

## Database Setup

### PostgreSQL

Apply schema:

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/mydb"
psql $DATABASE_URL -f internal/adapter/repository/postgres/schema.sql
```

Or with Docker:

```bash
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 \
  postgres:16-alpine

# Wait for startup
sleep 3

# Apply schema
psql -h localhost -U postgres -d mydb -f internal/adapter/repository/postgres/schema.sql
```

### MySQL

```bash
docker run -d \
  --name mysql \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=mydb \
  -p 3306:3306 \
  mysql:8

# Apply schema
mysql -h localhost -u root -proot mydb < internal/adapter/repository/mysql/schema.sql
```

### MongoDB

```bash
docker run -d \
  --name mongo \
  -p 27017:27017 \
  mongo:7

# Collections created automatically
```

## Generated Methods

Each repository implements:

### Core Methods

```go
type CustomerRepository interface {
    Save(ctx context.Context, customer *entity.Customer) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Customer, error)
}
```

### Additional Methods

Based on entity fields:

```go
// If entity has Email field
FindByEmail(ctx context.Context, email valueobject.Email) (*entity.Customer, error)

// If entity has unique fields
FindBySKU(ctx context.Context, sku valueobject.SKU) (*entity.Product, error)
```

## Integration with Wire

After generating repositories:

```bash
anaphase gen repository --domain customer --db postgres
anaphase wire
```

Wire automatically injects:
```go
func InitializeApp(logger *slog.Logger) (*App, error) {
    db, err := pgxpool.New(context.Background(), dbURL)

    customerRepo := postgres.NewCustomerRepository(db)
    // ...
}
```

## Connection String

Configure via environment:

```bash
# PostgreSQL
export DATABASE_URL="postgres://user:pass@host:5432/dbname?sslmode=disable"

# MySQL
export DATABASE_URL="user:pass@tcp(host:3306)/dbname?parseTime=true"

# MongoDB
export DATABASE_URL="mongodb://host:27017/dbname"
```

## Testing

### Unit Tests

Mock the repository:

```go
type MockCustomerRepository struct {
    SaveFunc     func(context.Context, *entity.Customer) error
    FindByIDFunc func(context.Context, uuid.UUID) (*entity.Customer, error)
}

func (m *MockCustomerRepository) Save(ctx context.Context, c *entity.Customer) error {
    return m.SaveFunc(ctx, c)
}
```

### Integration Tests

Use test database:

```go
func setupTestDB(t *testing.T) *pgxpool.Pool {
    db, err := pgxpool.New(context.Background(), "postgres://localhost/test")
    require.NoError(t, err)

    // Clean tables
    db.Exec(context.Background(), "TRUNCATE customers CASCADE")

    return db
}

func TestCustomerRepository_Save_Integration(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()

    repo := postgres.NewCustomerRepository(db)

    customer := &entity.Customer{
        ID:   uuid.New(),
        Email: email,
        Name: "John Doe",
    }

    err := repo.Save(context.Background(), customer)
    assert.NoError(t, err)

    // Verify
    found, err := repo.FindByID(context.Background(), customer.ID)
    assert.NoError(t, err)
    assert.Equal(t, customer.ID, found.ID)
}
```

## Customization

### Add Transactions

```go
func (r *customerRepository) SaveWithOrders(ctx context.Context, customer *entity.Customer, orders []*entity.Order) error {
    tx, err := r.db.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    // Save customer
    _, err = tx.Exec(ctx, "INSERT INTO customers ...", customer.ID)
    if err != nil {
        return err
    }

    // Save orders
    for _, order := range orders {
        _, err = tx.Exec(ctx, "INSERT INTO orders ...", order.ID)
        if err != nil {
            return err
        }
    }

    return tx.Commit(ctx)
}
```

### Add Pagination

```go
func (r *customerRepository) List(ctx context.Context, page, limit int) ([]*entity.Customer, error) {
    offset := (page - 1) * limit

    query := `
        SELECT id, email, name, created_at, updated_at
        FROM customers
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

    rows, err := r.db.Query(ctx, query, limit, offset)
    defer rows.Close()

    var customers []*entity.Customer
    for rows.Next() {
        var customer entity.Customer
        // Scan...
        customers = append(customers, &customer)
    }

    return customers, nil
}
```

### Add Filtering

```go
func (r *customerRepository) FindByFilters(ctx context.Context, filters map[string]interface{}) ([]*entity.Customer, error) {
    query := "SELECT * FROM customers WHERE 1=1"
    args := []interface{}{}
    argCount := 1

    if email, ok := filters["email"]; ok {
        query += fmt.Sprintf(" AND email = $%d", argCount)
        args = append(args, email)
        argCount++
    }

    // ...
}
```

## See Also

- [gen domain](/reference/gen-domain)
- [gen handler](/reference/gen-handler)
- [wire](/reference/wire)
- [Architecture](/guide/architecture)
