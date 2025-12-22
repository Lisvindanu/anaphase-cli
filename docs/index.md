---
layout: home

hero:
  name: Anaphase
  text: AI-Powered Microservice Generator
  tagline: Generate production-ready Golang microservices with AI. From idea to deployment in minutes.
  image:
    src: /hero-image.svg
    alt: Anaphase
  actions:
    - theme: brand
      text: Get Started
      link: /guide/quick-start
    - theme: alt
      text: View on GitHub
      link: https://github.com/lisvindanu/anaphase-cli

features:
  - icon: 🎯
    title: Domain-Driven Design First
    details: "**Our key differentiator.** True DDD with Aggregates, Entities, Value Objects, and Bounded Contexts. Not just MVC with extra layers - actual tactical DDD patterns that scale."

  - icon: 🤖
    title: AI-Powered Generation
    details: Leverage multiple AI providers (Gemini, Groq, OpenAI, Claude) to generate complete domain models from natural language. Just describe your business logic, get DDD-compliant code.

  - icon: ⚡
    title: Lightning Fast
    details: Generate complete CRUD APIs with handlers, repositories, and tests in seconds. Auto-wire dependencies and get running immediately.

  - icon: 🎯
    title: Type-Safe
    details: Strong typing throughout. Value objects, entities, and aggregates are generated with proper validation and business logic.

  - icon: 🔌
    title: Database Agnostic
    details: Support for PostgreSQL, MySQL, and MongoDB out of the box. Switch databases with a single flag.

  - icon: 📦
    title: Production Ready
    details: Generated code includes error handling, logging, graceful shutdown, health checks, and comprehensive tests.

  - icon: 🔄
    title: Auto-Wiring
    details: Automatic dependency injection with AST-based domain discovery. No manual wiring needed.

  - icon: 🛠️
    title: Extensible
    details: Customize generators, add your own templates, and integrate with your existing tools and workflows.
---

## Quick Example

Generate a complete e-commerce microservice in 3 commands:

```bash
# Initialize project
anaphase init my-ecommerce

# Generate domain with AI
anaphase gen domain --name customer --prompt "Customer with email, name, and billing address"

# Auto-wire and run
anaphase wire
go run cmd/api/main.go
```

Your API is now running at `http://localhost:8080` with:
- ✅ CRUD endpoints for customers
- ✅ PostgreSQL repository with migrations
- ✅ Input validation and error handling
- ✅ Structured logging
- ✅ Graceful shutdown
- ✅ Health checks

## Why Anaphase Over Goravel?

### 🎯 True Domain-Driven Design

**Anaphase** enforces tactical DDD patterns that protect your business logic:

```go
// ✅ Anaphase: Rich Domain Model (DDD)
type Order struct {
    ID          uuid.UUID
    Customer    Customer          // Aggregate Root
    Items       []OrderItem       // Entities
    ShippingAddress Address       // Value Object
    Status      OrderStatus       // Value Object
}

// Business logic IN the domain
func (o *Order) Cancel() error {
    if o.Status != Pending {
        return ErrCannotCancelOrder
    }
    o.Status = Cancelled
    o.RecordEvent(OrderCancelledEvent{...})
    return nil
}
```

**Goravel**: MVC with Active Record pattern:

```go
// ❌ Goravel: Anemic Domain Model (MVC)
type Order struct {
    orm.Model
    CustomerID  uint
    TotalAmount float64
    Status      string
}

// Business logic scattered in services/controllers
func CancelOrder(orderID uint) error {
    order := facades.Orm().Find(&Order{}, orderID)
    order.Status = "cancelled"
    order.Save()
}
```

### Key Architectural Differences

| Feature | Anaphase (DDD) | Goravel (MVC) |
|---------|----------------|---------------|
| **Architecture** | Hexagonal + DDD | MVC + Active Record |
| **Domain Model** | Rich (business logic in domain) | Anemic (logic in services) |
| **Aggregates** | ✅ First-class concept | ❌ No concept |
| **Value Objects** | ✅ Immutable, validated | ❌ Primitive types |
| **Bounded Contexts** | ✅ Explicit boundaries | ❌ No boundaries |
| **Domain Events** | ✅ Built-in support | ⚠️ Manual implementation |
| **Dependency Direction** | ✅ Inward (to domain) | ❌ Outward (from domain) |
| **Testability** | ✅ Pure domain, no DB | ⚠️ Coupled to framework |
| **Scalability** | ✅ Micro services ready | ⚠️ Monolith oriented |

### When to Choose Anaphase

✅ **Use Anaphase when you need:**
- Complex business logic that changes frequently
- Multiple microservices with clear boundaries
- Long-term maintainability (enterprise projects)
- Team scalability (multiple teams, different domains)
- True separation of concerns
- Framework independence

### When Goravel Works

✅ **Use Goravel when you need:**
- Simple CRUD applications
- Rapid prototyping
- Laravel-like development experience in Go
- Monolithic applications
- Small team with full-stack developers

## Why Anaphase?

### Traditional Approach
```bash
# Hours of boilerplate
mkdir -p internal/{domain,handler,repository}
# Write entity
# Write repository interface
# Write repository implementation
# Write handler
# Write DTOs
# Write tests
# Wire dependencies manually
# ... repeat for each domain
```

### With Anaphase
```bash
# Seconds to production-ready DDD code
anaphase gen domain "Order with items, can be cancelled if pending"
anaphase gen middleware --type auth
anaphase wire
# Done! Full DDD architecture ready
```

## Trusted by Developers

> "Anaphase transformed our development workflow. What used to take days now takes minutes."

> "The AI generation is incredibly accurate. It understands DDD patterns and generates clean code."

> "Best tool for bootstrapping microservices. The auto-wiring feature alone saves hours."

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
