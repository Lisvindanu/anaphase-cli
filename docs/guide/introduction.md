# What is Anaphase?

Anaphase is an **interactive CLI tool** that generates production-ready Golang microservices following best practices in Domain-Driven Design (DDD) and Clean Architecture.

**Works with or without AI** - choose Template Mode for instant scaffolding or AI Mode for intelligent generation.

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

Anaphase generates all the necessary code automatically, following established patterns. Choose your workflow:

- **🎨 Interactive Menu**: No commands to memorize (v0.4)
- **📝 Template Mode**: Instant scaffolding without AI
- **🤖 AI Mode**: Smart generation from natural language

### Key Features

#### 🎨 Interactive Menu (New in v0.4!)

Just run `anaphase` - no commands to memorize:

```bash
anaphase
# Interactive menu with search (Ctrl+K), keyboard navigation, and filtering
```

Select what you need from a beautiful TUI interface. Perfect for:
- **Beginners**: Discover available commands
- **Pros**: Quick access with `/` search filter

#### 🤖 Dual Mode Generation

**Template Mode** (no setup required):
```bash
anaphase gen domain
# Prompts for: Entity name, Fields (name:type)
# Generates: Entity, Repository, Service interfaces
```

**AI Mode** (optional - requires API key):
```bash
anaphase gen domain "Order with items, can be cancelled if pending"
# AI understands: Entities vs Value Objects, Aggregates, Business rules
# Generates: Advanced validation, domain events, business logic
```

The AI Mode additionally provides:
- **Smart type inference** from natural language
- **Business rules** and validation logic
- **Value Objects** with immutability
- **Domain events** and relationships

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

From zero to running API in minutes - **no AI required**:

| Task | Traditional | With Anaphase (Template) | With Anaphase (AI) |
|------|-------------|--------------------------|---------------------|
| Project setup | 30 min | 10 seconds | 10 seconds |
| Domain model | 2 hours | 30 seconds | 30 seconds (smarter) |
| Repository | 1 hour | 10 seconds | 10 seconds |
| Handler | 1 hour | 10 seconds | 10 seconds |
| Auto-wiring | 30 min | 5 seconds | 5 seconds |
| **Total** | **~5 hours** | **~1 minute** | **~1 minute** |

::: tip
Both Template and AI modes generate production-ready code. Template Mode is perfect for standard CRUD, AI Mode adds intelligent business logic.
:::

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

### 1. Interactive Selection

Choose your command from the TUI menu or use CLI directly. Filter with `/` search, navigate with arrows.

### 2. AST Analysis

Anaphase uses Go's AST parser to analyze your codebase and discover existing domains automatically.

### 3. Code Generation (Dual Mode)

**Template Mode** (default - no AI):
- Uses intelligent templates for standard patterns
- Prompts for entity names and field types
- Generates DDD-compliant code instantly

**AI Mode** (optional):
- Leverages LLMs (Gemini, OpenAI, Claude, Groq)
- Understands natural language descriptions
- Generates advanced validation and business logic

### 4. Code Output

Both modes generate production-ready Go code:
- Proper imports and packages
- Error handling
- Validation (Template: basic, AI: advanced)
- Tests
- Documentation

### 5. Auto-Setup & Wiring

Automatically configures everything:
- `.env` files with database URLs
- `go.mod` dependencies (auto-installed)
- Dependency injection
- HTTP routes and middleware

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

Generate multiple bounded contexts quickly with the interactive menu:

```bash
anaphase  # Opens interactive menu
# Select: Generate Domain (repeat for each domain)
# Select: Auto-Wire Dependencies
```

Or via CLI:
```bash
anaphase gen domain "user"
anaphase gen domain "product"
anaphase gen domain "order"
# Auto-wiring happens automatically or run: anaphase wire
```

### API Backends

Complete REST APIs with CRUD operations - no AI required:

```bash
anaphase  # Interactive menu
# 1. Select: Initialize Project → Enter name and database
# 2. Select: Generate Domain → Enter entity and fields
# 3. Select: Generate Handler → Enter domain name
# 4. Select: Generate Repository → Enter domain name
# Done! All dependencies auto-installed.
```

### Quick Prototyping

Template Mode is perfect for rapid prototyping:

```bash
anaphase init my-prototype --db sqlite
cd my-prototype
anaphase  # Generate domains interactively
make run  # It just works!
```

## Comparison

| Feature | Anaphase | Manual | Other Generators |
|---------|----------|--------|------------------|
| **Interactive Menu** | ✅ (v0.4) | ❌ | ❌ |
| **Works Without AI** | ✅ Template Mode | N/A | ✅ |
| **AI-Powered (Optional)** | ✅ | ❌ | ❌ |
| **DDD Support** | ✅ Both modes | Depends | ⚠️ |
| **Clean Architecture** | ✅ Enforced | Depends | ⚠️ |
| **Auto-Setup** | ✅ .env, deps | ❌ | ⚠️ |
| **Auto-Wiring** | ✅ | ❌ | ⚠️ |
| **Type-Safe** | ✅ | Depends | ⚠️ |
| **Multi-DB** | ✅ | ❌ | ✅ |
| **Production-Ready** | ✅ | Depends | ⚠️ |
| **Learning Curve** | Very Low | High | Medium |
| **Setup Time** | 0 seconds | Hours | Minutes |

## Philosophy

Anaphase v0.4 is built on these principles:

1. **Flexible Workflows**: Choose what fits - Interactive Menu, Template Mode, or AI Mode
2. **No Barriers to Entry**: Works immediately without any API keys or configuration
3. **AI-Assisted, Not Required**: AI enhances generation but isn't mandatory
4. **Patterns Over Configuration**: Enforce DDD best practices by default
5. **Transparency**: See exactly what's generated, no magic
6. **Production-First**: Generate code you'd actually deploy
7. **Developer Experience**: Beautiful TUI, search, auto-setup - make it delightful

## What's Next?

- [Quick Start](/guide/quick-start) - Build your first service
- [Installation](/guide/installation) - Detailed setup
- [Architecture](/guide/architecture) - Deep dive into patterns
