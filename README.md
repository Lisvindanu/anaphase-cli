# Anaphase CLI

> AI-Powered Golang Microservice Generator with DDD Architecture Enforcement

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Overview

Anaphase is an intelligent code scaffolding tool that generates production-ready Golang microservices following Domain-Driven Design (DDD) and Clean Architecture principles.

### Key Features

- 🤖 **AI-Powered Generation** - Natural language to code using LLM
- 🏗️ **Architecture Enforcement** - Strict DDD/Clean Architecture patterns
- 🔌 **Auto-Wiring** - AST-based dependency injection
- 🧪 **Complete Testing** - Unit and integration tests generated
- 📚 **Auto Documentation** - Swagger/OpenAPI specs
- 🆓 **Zero Cost** - Free tier using Groq API
- 🚫 **Zero Lock-in** - Generated code works standalone

## Quick Start

### Installation

**Quick install (recommended):**

```bash
curl -fsSL https://raw.githubusercontent.com/lisvindanu/anaphase-cli/main/install.sh | bash
```

**Manual install:**

```bash
go install github.com/lisvindanu/anaphase-cli/cmd/anaphase@latest

# Add to PATH if needed
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

### Configuration

```bash
# Initialize configuration
anaphase config init

# Set Groq API key (FREE)
export GROQ_API_KEY=gsk_your_key_here
```

### Usage

```bash
# Create a new project
anaphase init my-shop --module github.com/mycompany/my-shop

# Generate a domain
cd my-shop
anaphase gen domain "Create a Cart domain with Items"

# Generate HTTP handlers
anaphase gen handler cart

# Generate repository
anaphase gen repository cart --db postgres

# Generate tests
anaphase gen test cart
```

## Architecture

```
┌─────────────┐
│   Handler   │ ← HTTP/gRPC Adapter
└──────┬──────┘
       │
┌──────▼──────┐
│   Service   │ ← Business Logic
└──────┬──────┘
       │
┌──────▼──────┐
│ Repository  │ ← Data Access
└──────┬──────┘
       │
┌──────▼──────┐
│  Database   │
└─────────────┘
```

## Project Status

This is a thesis project (Tugas Akhir) developed as part of academic research in AI-assisted software engineering.

**Current Phase:** Phase 1 - Skeleton Implementation

## Documentation

- [Technical Specification](TA.md)
- [AI Integration Strategy](ADDENDUM.MD)

## License

MIT License - see LICENSE file for details

## Author

Built with ❤️ for thesis research
