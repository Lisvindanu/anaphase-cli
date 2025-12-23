# Domain-Driven Design

**Inilah yang membuat Anaphase berbeda.** Sementara framework lain menggunakan pola MVC atau Active Record, Anaphase memaksakan prinsip Domain-Driven Design yang sesungguhnya untuk menjaga logika bisnis Anda tetap bersih, mudah ditest, dan maintainable.

::: tip DDD di Kedua Mode
Sejak v0.4.0, Anaphase menawarkan dua mode generasi, keduanya menghasilkan kode yang sesuai dengan DDD:

- **Template Mode**: Menghasilkan struktur DDD yang bersih dengan entities, value objects, repositories, dan ports
- **AI Mode**: Menghasilkan pola DDD tingkat lanjut termasuk aggregates, domain events, dan aturan bisnis yang kompleks

Terlepas dari mode yang digunakan, semua kode yang dihasilkan mengikuti prinsip arsitektur yang dijelaskan dalam panduan ini.
:::

## Apa itu DDD?

Domain-Driven Design (DDD) adalah pendekatan pengembangan software untuk **domain yang kompleks** yang:

1. **Fokus pada core domain** dan logika domain
2. **Menggunakan ubiquitous language** yang dibagikan oleh developer dan domain experts
3. **Memodelkan domain yang kompleks** melalui entities, value objects, dan aggregates
4. **Mengisolasi logika domain** dari concern infrastruktur

## Mengapa DDD Lebih Baik dari MVC?

### Masalah dengan MVC Tradisional

Kebanyakan framework Go (seperti Goravel) menggunakan **pola MVC dengan Active Record**:

```go
// ❌ MVC/Active Record: Business logic tersebar di mana-mana
type Order struct {
    orm.Model
    UserID      uint
    TotalAmount float64
    Status      string  // Just a string, no validation
}

// Logic in Controller
func (c *OrderController) Cancel(ctx *gin.Context) {
    order := models.Order.Find(id)
    if order.Status != "pending" {  // Business rule in controller
        return errors.New("cannot cancel")
    }
    order.Status = "cancelled"
    order.Save()  // Coupled to database
}

// Logic in Service
func (s *OrderService) CalculateTotal(order *Order) {
    // More business logic scattered in services
}
```

**Masalah:**
- ❌ Logika bisnis tersebar di controllers, services, models
- ❌ Pengetahuan domain tercampur dengan concern teknis (DB, HTTP)
- ❌ Sulit untuk test (semuanya coupled ke framework)
- ❌ Sulit dipahami (di mana logika bisnisnya?)
- ❌ Tidak bisa diubah (memodifikasi satu rule mempengaruhi banyak file)

### Solusi DDD

Anaphase memaksakan **Rich Domain Models** di mana logika bisnis berada di domain:

```go
// ✅ DDD: Business logic encapsulated in domain
type Order struct {
    ID              uuid.UUID
    Customer        *Customer      // Aggregate Root
    Items           []OrderItem    // Entities
    ShippingAddress Address        // Value Object
    Status          OrderStatus    // Value Object (type-safe)
    Total           Money          // Value Object (immutable)
}

// Business logic IN the domain
func (o *Order) Cancel() error {
    // Business rule enforced here
    if !o.CanBeCancelled() {
        return ErrCannotCancelOrder
    }

    o.Status = OrderCancelled
    o.RecordEvent(OrderCancelledEvent{OrderID: o.ID})
    return nil
}

func (o *Order) CanBeCancelled() bool {
    // Complex business rules in one place
    return o.Status == OrderPending || o.Status == OrderConfirmed
}

// Invariants protected
func (o *Order) AddItem(product *Product, quantity int) error {
    if quantity <= 0 {
        return ErrInvalidQuantity
    }

    item := NewOrderItem(product, quantity)
    o.Items = append(o.Items, item)
    o.RecalculateTotal() // Aggregate maintains consistency
    return nil
}
```

**Keuntungan:**
- ✅ Semua logika bisnis di domain (mudah ditemukan dan dipahami)
- ✅ Pure Go (tidak ada ketergantungan framework, mudah untuk test)
- ✅ Type-safe (compiler menangkap error)
- ✅ Self-documenting (kode terbaca seperti bahasa bisnis)
- ✅ Change-friendly (modifikasi rules di satu tempat)

## Kapan Menggunakan DDD vs MVC

### Gunakan DDD (Anaphase) Ketika:

✅ **Logika Bisnis yang Kompleks**
- Banyak aturan bisnis per operasi
- Rules yang sering berubah
- Domain experts terlibat dalam requirements

✅ **Proyek Jangka Panjang**
- Aplikasi enterprise
- Lifecycle 5+ tahun
- Banyak tim bekerja pada domain yang berbeda

✅ **Arsitektur Microservices**
- Bounded contexts yang jelas
- Independent deployments
- Tim yang berbeda memiliki service yang berbeda

✅ **Kompleksitas Domain**
- E-commerce dengan inventory, pricing rules, promotions
- Sistem finansial dengan kalkulasi kompleks
- Healthcare dengan persyaratan regulasi
- Logistik dengan optimisasi rute

### Gunakan MVC (Goravel) Ketika:

✅ **CRUD Sederhana**
- Basic create, read, update, delete
- Sedikit aturan bisnis
- Aplikasi yang data-centric

✅ **Rapid Prototyping**
- MVP atau proof of concept
- Proyek jangka pendek
- Demo cepat

✅ **Tim Kecil**
- 1-3 developers
- Full-stack developers
- Semua orang tahu segalanya

## Konsep Kunci DDD

Memahami tactical patterns ini sangat penting untuk menggunakan Anaphase secara efektif:

## Core Building Blocks

### Entities

Objek dengan **identitas** yang persisten dari waktu ke waktu.

**Karakteristik:**
- Memiliki identifier unik (ID)
- Dapat berubah state (mutable)
- Lifecycle yang dilacak (CreatedAt, UpdatedAt)
- Berisi logika bisnis

**Contoh:**
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

**Kapan menggunakan:**
- Objek membutuhkan identitas unik
- Objek berubah dari waktu ke waktu
- Anda peduli tentang instance spesifik mana

**Contoh:**
- Customer, Order, Product
- User, Invoice, Account
- Booking, Shipment, Payment

### Value Objects

Objek tanpa identitas, didefinisikan berdasarkan **atribut** mereka.

**Karakteristik:**
- Tidak ada ID
- Immutable (tidak dapat berubah)
- Dibandingkan berdasarkan nilai, bukan identitas
- Validasi mandiri (self-validating)

**Contoh:**
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

**Kapan menggunakan:**
- Merepresentasikan konsep atau pengukuran
- Tidak memerlukan identitas
- Dapat dibagikan dan diganti

**Contoh:**
- Email, Phone, Address
- Money, Quantity, Price
- DateRange, Coordinates, URL

### Aggregates

Cluster dari entities dan value objects yang diperlakukan sebagai **satu unit**.

**Karakteristik:**
- Memiliki root entity (aggregate root)
- Memaksakan invariants
- Transaction boundary
- Consistency boundary

**Contoh:**
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

**Aturan:**
- Hanya referensi by ID, tidak secara langsung
- Perubahan melalui root
- Repository hanya untuk root

**Contoh:**
- Order (dengan LineItems)
- ShoppingCart (dengan CartItems)
- Invoice (dengan InvoiceLines)

### Repositories

Abstraksi untuk **persistence dan retrieval** dari aggregates.

**Karakteristik:**
- Satu repository per aggregate root
- Bekerja dengan complete aggregates
- Interface (port) di domain layer
- Implementation di infrastructure layer

**Contoh:**
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

**Keuntungan:**
- Swap implementations (Postgres, MySQL, Mock)
- Test tanpa database
- Domain tidak bergantung pada infrastructure

### Services

Operasi yang **tidak cocok secara natural** di entities atau value objects.

**Domain Services:**
- Mengkoordinasikan banyak entities
- Mengimplementasikan proses bisnis
- Stateless

**Contoh:**
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

## Pola DDD di Anaphase

### Entity Generation

```bash
anaphase gen domain --name customer --prompt "Customer with email and name"
```

Menghasilkan entity dengan:
- Unique ID (uuid.UUID)
- Lifecycle tracking (CreatedAt, UpdatedAt)
- Constructor dengan validasi
- Validate() method

### Value Object Detection

AI mengenali value objects:

```bash
--prompt "Customer with email (validated), billing address"
```

Menghasilkan:
- `Email` value object dengan validasi
- `Address` value object (composite)

### Aggregate Modeling

```bash
--prompt "Order with line items. Each line item has product and quantity."
```

AI memahami:
- Order adalah aggregate root
- LineItem adalah bagian dari aggregate
- Hanya Order yang mendapat repository

### Repository Interfaces

Repository yang dihasilkan mengikuti aturan aggregate:

```go
// ✅ Repository for aggregate root
type OrderRepository interface {
    Save(ctx context.Context, order *entity.Order) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Order, error)
}

// ❌ No repository for LineItem (part of aggregate)
```

## Pola Umum

### Money Pattern

Selalu gunakan value object untuk uang:

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

**Mengapa:**
- Hindari error floating-point
- Memaksakan pencocokan mata uang
- Encapsulate operasi uang

### Enum Pattern

Gunakan typed constants untuk status:

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

Gunakan constructors untuk validasi:

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

## Anti-Patterns yang Harus Dihindari

### Anemic Domain Model

❌ **Jangan:**
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

✅ **Lakukan:**
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

❌ **Jangan:**
```go
type Order struct {
    Items []*LineItem  // Direct access
}

// External code modifies directly
order.Items = append(order.Items, newItem)
```

✅ **Lakukan:**
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

❌ **Jangan:**
```go
type Customer struct {
    ID      uuid.UUID
    Orders  []*Order   // ❌ Too large
    Invoices []*Invoice // ❌ Too large
}
```

✅ **Lakukan:**
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

### 1. Gunakan Ubiquitous Language

Gunakan terminologi domain:

```bash
# Finance domain
"account" bukan "thing"
"transaction" bukan "record"
"balance" bukan "amount"

# E-commerce domain
"order" bukan "purchase"
"inventory" bukan "stock count"
"SKU" bukan "product code"
```

### 2. Buat Invariants Eksplisit

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

### 3. Jaga Aggregates Tetap Kecil

Hanya sertakan apa yang harus konsisten:

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

### 4. Validasi di Boundaries

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

- [Architecture](/guide/architecture) - Lihat bagaimana DDD cocok di Clean Architecture
- [AI Generation](/guide/ai-generation) - Bagaimana AI menghasilkan kode DDD
- [Examples](/examples/basic) - DDD dalam praktik
