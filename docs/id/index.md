---
layout: home

hero:
  name: Anaphase
  text: Generator Microservice dengan AI
  tagline: Generate microservice Golang production-ready dengan AI. Dari ide ke deployment dalam hitungan menit.
  image:
    src: /hero-image.svg
    alt: Anaphase
  actions:
    - theme: brand
      text: Mulai Sekarang
      link: /id/guide/quick-start
    - theme: alt
      text: Lihat di GitHub
      link: https://github.com/lisvindanu/anaphase-cli

features:
  - icon: 🎯
    title: Domain-Driven Design First
    details: **Pembeda utama kami.** DDD sejati dengan Aggregates, Entities, Value Objects, dan Bounded Contexts. Bukan sekedar MVC dengan layer tambahan - tapi pola DDD taktis yang bisa scale.

  - icon: 🤖
    title: AI-Powered Generation
    details: Gunakan berbagai AI provider (Gemini, Groq, OpenAI, Claude) untuk generate domain model lengkap dari bahasa natural. Deskripsikan logika bisnis Anda, dapat code yang sesuai DDD.

  - icon: ⚡
    title: Super Cepat
    details: Generate CRUD API lengkap dengan handlers, repositories, dan tests dalam hitungan detik. Auto-wire dependencies dan langsung running.

  - icon: 🎯
    title: Type-Safe
    details: Strong typing di semua layer. Value objects, entities, dan aggregates di-generate dengan validasi dan business logic yang proper.

  - icon: 🔌
    title: Database Agnostic
    details: Support PostgreSQL, MySQL, dan MongoDB out of the box. Ganti database cukup dengan satu flag.

  - icon: 📦
    title: Production Ready
    details: Code yang di-generate sudah include error handling, logging, graceful shutdown, health checks, dan comprehensive tests.

  - icon: 🔄
    title: Auto-Wiring
    details: Dependency injection otomatis dengan AST-based domain discovery. Tidak perlu manual wiring.

  - icon: 🛠️
    title: Extensible
    details: Customize generators, tambah template sendiri, dan integrasikan dengan tools dan workflow yang sudah ada.
---

## Contoh Cepat

Generate microservice e-commerce lengkap dalam 3 perintah:

```bash
# Initialize project
anaphase init my-ecommerce

# Generate domain dengan AI
anaphase gen domain "Customer dengan email, nama, dan alamat billing"

# Auto-wire dan jalankan
anaphase wire
go run cmd/api/main.go
```

API Anda sekarang running di `http://localhost:8080` dengan:
- ✅ CRUD endpoints untuk customers
- ✅ PostgreSQL repository dengan migrations
- ✅ Input validation dan error handling
- ✅ Structured logging
- ✅ Graceful shutdown
- ✅ Health checks

## Kenapa Anaphase vs Goravel?

### 🎯 True Domain-Driven Design

**Anaphase** enforce pola DDD taktis yang melindungi business logic Anda:

```go
// ✅ Anaphase: Rich Domain Model (DDD)
type Order struct {
    ID          uuid.UUID
    Customer    Customer          // Aggregate Root
    Items       []OrderItem       // Entities
    ShippingAddress Address       // Value Object
    Status      OrderStatus       // Value Object
}

// Business logic ADA DI domain
func (o *Order) Cancel() error {
    if o.Status != Pending {
        return ErrCannotCancelOrder
    }
    o.Status = Cancelled
    o.RecordEvent(OrderCancelledEvent{...})
    return nil
}
```

**Goravel**: MVC dengan Active Record pattern:

```go
// ❌ Goravel: Anemic Domain Model (MVC)
type Order struct {
    orm.Model
    CustomerID  uint
    TotalAmount float64
    Status      string
}

// Business logic tersebar di services/controllers
func CancelOrder(orderID uint) error {
    order := facades.Orm().Find(&Order{}, orderID)
    order.Status = "cancelled"
    order.Save()
}
```

### Perbedaan Arsitektur Utama

| Fitur | Anaphase (DDD) | Goravel (MVC) |
|---------|----------------|---------------|
| **Arsitektur** | Hexagonal + DDD | MVC + Active Record |
| **Domain Model** | Rich (logic di domain) | Anemic (logic di services) |
| **Aggregates** | ✅ Konsep utama | ❌ Tidak ada |
| **Value Objects** | ✅ Immutable, validated | ❌ Primitive types |
| **Bounded Contexts** | ✅ Boundary eksplisit | ❌ Tidak ada boundary |
| **Domain Events** | ✅ Built-in support | ⚠️ Manual implementation |
| **Dependency Direction** | ✅ Ke dalam (ke domain) | ❌ Ke luar (dari domain) |
| **Testability** | ✅ Pure domain, tanpa DB | ⚠️ Coupled ke framework |
| **Scalability** | ✅ Siap microservices | ⚠️ Oriented ke monolith |

### Kapan Pilih Anaphase

✅ **Gunakan Anaphase ketika:**
- Business logic yang kompleks dan sering berubah
- Multiple microservices dengan boundary yang jelas
- Maintainability jangka panjang (proyek enterprise)
- Scalability team (banyak team, domain berbeda)
- Separation of concerns yang benar
- Framework independence

### Kapan Goravel Cocok

✅ **Gunakan Goravel ketika:**
- Aplikasi CRUD sederhana
- Rapid prototyping
- Pengalaman development ala Laravel di Go
- Aplikasi monolithic
- Team kecil dengan full-stack developers

## Kenapa Anaphase?

### Cara Traditional
```bash
# Berjam-jam boilerplate
mkdir -p internal/{domain,handler,repository}
# Tulis entity
# Tulis repository interface
# Tulis repository implementation
# Tulis handler
# Tulis DTOs
# Tulis tests
# Wire dependencies manual
# ... ulangi untuk setiap domain
```

### Dengan Anaphase
```bash
# Detik ke production-ready DDD code
anaphase gen domain "Order dengan items, bisa dibatalkan jika pending"
anaphase gen middleware --type auth
anaphase wire
# Selesai! Arsitektur DDD lengkap siap pakai
```

## Dipercaya Developer

> "Anaphase mengubah workflow development kami. Yang dulu butuh berhari-hari sekarang hanya beberapa menit."

> "AI generation-nya sangat akurat. Mengerti pola DDD dan generate code yang clean."

> "Tool terbaik untuk bootstrapping microservices. Fitur auto-wiring saja menghemat berjam-jam."

<style>
:root {
  --vp-home-hero-name-color: transparent;
  --vp-home-hero-name-background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);

  --vp-c-brand: #667eea;
  --vp-c-brand-light: #764ba2;
  --vp-c-brand-lighter: #8b7fc5;
  --vp-c-brand-dark: #5568d3;
  --vp-c-brand-darker: #4451b8;
}
</style>
