# What is Anaphase?

Anaphase is an **AI-powered CLI tool** that generates production-ready Golang microservices following best practices in Domain-Driven Design (DDD) and Clean Architecture.

## The Problem

Building microservices involves writing repetitive boilerplate code:

- Entity definitions with validation
- Repository interfaces and implementations
- Service layer with business logic
- HTTP handlers and DTOs
- Database schemas and migrations
- Dependency injection
- Tests for all layers

This process is:
- **Time-consuming**: Hours of setup for each domain
- **Error-prone**: Easy to miss patterns or make mistakes
- **Inconsistent**: Different developers, different styles
- **Tedious**: Same patterns over and over

## The Solution

Anaphase uses AI to understand your domain requirements and generates all the necessary code automatically, following established patterns.

### Key Features

#### 🤖 AI-Powered Generation

Describe your domain in natural language, and Anaphase generates complete, compilable Go code:

```bash
anaphase gen domain --name order --prompt \
  "Order with customer reference, line items with products and quantities,
   total amount, and status (pending, confirmed, shipped, delivered)"
```

The AI understands:
- **Entities** vs **Value Objects**
- **Aggregates** and their boundaries
- **Business rules** and validation
- **Relationships** between domains

#### 🏗️ Clean Architecture

All generated code follows Clean Architecture principles:

```
┌─────────────────────────────────────┐
│         Presentation Layer          │
│  (HTTP Handlers, gRPC, GraphQL)     │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│        Application Layer            │
│      (Use Cases, Services)          │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│          Domain Layer               │
│   (Entities, Value Objects, Ports)  │
└──────────────┬──────────────────────┘
               │
┌──────────────▼──────────────────────┐
│      Infrastructure Layer           │
│  (Database, External APIs, Cache)   │
└─────────────────────────────────────┘
```

**Benefits:**
- Independent of frameworks
- Testable at every layer
- Flexible infrastructure
- Business logic protected

#### ⚡ Rapid Development

From zero to running API in minutes:

| Task | Traditional | With Anaphase |
|------|-------------|---------------|
| Project setup | 30 min | 10 seconds |
| Domain model | 2 hours | 30 seconds |
| Repository | 1 hour | 10 seconds |
| Handler | 1 hour | 10 seconds |
| Wiring | 30 min | 5 seconds |
| **Total** | **~5 hours** | **~1 minute** |

#### 🎯 Type-Safe

Strong typing throughout:

```go
// Value Objects with validation
type Email struct {
    value string
}

func NewEmail(value string) (*Email, error) {
    if !isValidEmail(value) {
        return nil, ErrInvalidEmail
    }
    return &Email{value: value}, nil
}

// Entities with business logic
type Customer struct {
    ID        uuid.UUID
    Email     *Email  // Type-safe value object
    Name      string
    CreatedAt time.Time
}
```

## How It Works

### 1. AST Analysis

Anaphase uses Go's AST parser to analyze your codebase and discover existing domains automatically.

### 2. AI Generation

Leverages Google Gemini with carefully engineered prompts to generate:
- Domain models following DDD
- Repository patterns
- Service interfaces
- Handler implementations

### 3. Code Generation

Generates production-ready Go code:
- Proper imports and packages
- Error handling
- Validation
- Tests
- Documentation

### 4. Auto-Wiring

Scans generated code and automatically wires:
- Dependencies
- Database connections
- HTTP routes
- Middleware

## Architecture Patterns

Anaphase enforces best practices:

### Domain-Driven Design (DDD)

- **Entities**: Objects with identity
- **Value Objects**: Immutable objects without identity
- **Aggregates**: Clusters of entities and value objects
- **Repositories**: Persistence abstraction
- **Services**: Business logic coordination

### Hexagonal Architecture

- **Ports**: Interfaces defining contracts
- **Adapters**: Implementations of ports
- **Core**: Business logic isolated from infrastructure

### SOLID Principles

- Single Responsibility
- Open/Closed
- Liskov Substitution
- Interface Segregation
- Dependency Inversion

## Use Cases

### Microservices

Generate multiple bounded contexts quickly:

```bash
anaphase gen domain --name user
anaphase gen domain --name product
anaphase gen domain --name order
anaphase wire
```

### API Backends

Complete REST APIs with CRUD operations:

```bash
anaphase gen domain --name article --prompt "Blog article with title, content, author"
anaphase gen handler --domain article
anaphase gen repository --domain article
anaphase wire
```

### Database Migrations

When you need to start fresh or change schemas:

```bash
anaphase gen repository --domain customer --db postgres
# Apply internal/adapter/repository/postgres/schema.sql
```

## Comparison

| Feature | Anaphase | Manual | Other Generators |
|---------|----------|--------|------------------|
| AI-Powered | ✅ | ❌ | ❌ |
| DDD Support | ✅ | Depends | ⚠️ |
| Clean Architecture | ✅ | Depends | ⚠️ |
| Auto-Wiring | ✅ | ❌ | ⚠️ |
| Type-Safe | ✅ | Depends | ⚠️ |
| Multi-DB | ✅ | ❌ | ✅ |
| Production-Ready | ✅ | Depends | ⚠️ |
| Learning Curve | Low | High | Medium |

## Philosophy

Anaphase is built on these principles:

1. **AI-Assisted, Not AI-Replaced**: AI helps with repetitive tasks, you focus on business logic
2. **Patterns Over Configuration**: Enforce best practices by default
3. **Flexibility**: Generated code is yours to modify
4. **Transparency**: See exactly what's generated, no magic
5. **Production-First**: Generate code you'd actually deploy

## What's Next?

- [Quick Start](/guide/quick-start) - Build your first service
- [Installation](/guide/installation) - Detailed setup
- [Architecture](/guide/architecture) - Deep dive into patterns
