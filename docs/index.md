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
      link: https://github.com/lisvindanuu/anaphase-cli

features:
  - icon: 🤖
    title: AI-Powered Generation
    details: Leverage Google Gemini to generate complete domain models from natural language descriptions. Just describe what you need, and get production-ready code.

  - icon: 🏗️
    title: Clean Architecture
    details: Built-in support for Domain-Driven Design (DDD), Hexagonal Architecture, and Clean Architecture patterns. Your code stays maintainable and testable.

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
# Seconds to production
anaphase gen domain --name product --prompt "Product with SKU, price, inventory"
anaphase wire
# Done!
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
