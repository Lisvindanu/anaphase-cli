# Quick Start

Get started with Anaphase in under 5 minutes. **No AI API key required** - works out of the box!

## Prerequisites

- Go 1.21 or higher
- PostgreSQL, MySQL, or SQLite (optional, for database features)
- AI API key (optional, for AI-powered generation)

## Installation

### From Source

```bash
git clone https://github.com/lisvindanu/anaphase-cli.git
cd anaphase-cli
go install ./cmd/anaphase
```

### Using Go Install

```bash
go install github.com/lisvindanu/anaphase-cli/cmd/anaphase@latest
```

Verify installation:

```bash
anaphase --version
```

## Configure AI Provider (Optional)

::: info Two Modes Available
Anaphase works in **two modes**:
- **Template Mode**: Works immediately without any API key (basic CRUD scaffolding)
- **AI Mode**: Smart generation from natural language (requires API key)

**You can start using Anaphase right away with Template Mode!**
:::

If you want AI-powered generation, set up an API key:

```bash
export GEMINI_API_KEY="your-api-key-here"
```

Or create a config file at `~/.anaphase/config.yaml`:

```yaml
ai:
  primary:
    type: gemini
    apiKey: your-api-key-here
    model: gemini-2.5-flash
```

::: tip Get a Free API Key
Get a free Gemini API key at [Google AI Studio](https://makersuite.google.com/app/apikey).

Anaphase also supports OpenAI, Claude, and Groq. [Learn more →](/config/ai-providers)
:::

## Create Your First Project

### Interactive Menu (Recommended)

**New in v0.4!** Just run `anaphase` to access the interactive menu - no commands to memorize:

```bash
anaphase
```

The interactive menu appears:

```
⚡ Anaphase CLI - DDD Microservice Generator
   💡 Commands marked [AI] require API key setup

▶ 🚀 Initialize Project
  🤖 Generate Domain [AI]
  📡 Generate Handler
  💾 Generate Repository
  🛡️  Generate Middleware
  📊 Generate Migration
  🔌 Auto-Wire Dependencies
  📐 Describe Architecture
  ✨ Code Quality
  ⚙️  Configuration

⌨️  Keys: ↑↓ navigate • / filter • Enter select • q quit
```

Select **"Initialize Project"** and follow the prompts:

```bash
Project name: my-app
Database type (postgres/mysql/sqlite) [postgres]: postgres

✅ Project created with auto-generated .env and dependencies!
```

::: tip Pro Tip
Use `/` to search/filter commands in the interactive menu. Try typing "domain" to quickly find domain generation!
:::

### Command Line (Alternative)

You can also use direct commands:

```bash
anaphase init my-app --db postgres
cd my-app
```

Both methods generate a complete project structure:

```
my-app/
├── cmd/
│   └── api/
├── internal/
│   ├── core/
│   │   ├── entity/
│   │   ├── port/
│   │   └── valueobject/
│   └── adapter/
│       ├── handler/
│       └── repository/
├── .env          # Auto-generated with database config
├── .env.example
├── go.mod
└── README.md
```

### Step 2: Generate a Domain

#### Using Interactive Menu

```bash
cd my-app
anaphase
```

Select **"Generate Domain"** from the menu. Anaphase will automatically use:
- **Template Mode** if no AI is configured (works immediately!)
- **AI Mode** if API key is set up

**Template Mode Example:**

```
📝 Template Mode - Domain Generation

Entity name: Customer
Fields: email:string, name:string, phone:string

✅ Generated:
  ✓ internal/core/entity/customer.go
  ✓ internal/core/port/customer_repository.go
  ✓ internal/core/port/customer_service.go
```

**AI Mode Example** (with API key configured):

```
🧠 AI-Powered Domain Generation

Description: Customer with email, name, phone. Can place orders.

✅ Generated:
  ✓ internal/core/entity/customer.go (with validation)
  ✓ internal/core/valueobject/email.go (email validation)
  ✓ internal/core/valueobject/phone.go (phone validation)
  ✓ internal/core/port/customer_repository.go
  ✓ internal/core/port/customer_service.go
```

#### Using Command Line

**Template Mode:**
```bash
anaphase gen domain "Customer"
# Prompts for entity name and fields interactively
```

**AI Mode:**
```bash
anaphase gen domain "Customer with email, name, and phone. Can place orders."
```

Both modes create DDD-compliant domain models with:
- Entity with business logic
- Repository interface (port)
- Service interface (port)
- Value objects (AI mode adds smart validation)

### Step 3: Generate Handlers

Using the interactive menu, select **"Generate Handler"**:

```bash
Handler name: customer

✅ Generated:
  ✓ internal/adapter/handler/http/customer_handler.go (CRUD endpoints)
  ✓ internal/adapter/handler/http/customer_dto.go (Request/Response DTOs)
  ✓ internal/adapter/handler/http/customer_handler_test.go
```

Or via command line:
```bash
anaphase gen handler customer
```

### Step 4: Generate Repository

Select **"Generate Repository"** from the menu:

```bash
Repository name: customer

✅ Generated:
  ✓ internal/adapter/repository/postgres/customer_repo.go
  ✓ internal/adapter/repository/postgres/schema.sql
  ✓ internal/adapter/repository/postgres/customer_repo_test.go
```

Or via command line:
```bash
anaphase gen repository customer
```

### Step 5: Run Your Application

**Auto-setup is already done!** Your `.env` file was created during `init`. Just start the database and run:

```bash
# Start PostgreSQL with Docker
docker run -d \
  --name anaphase-db \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=my-app \
  -p 5432:5432 \
  postgres:16-alpine

# Run your API (dependencies already installed!)
make run
```

::: tip Database Credentials
The `.env` file is auto-generated with the correct DATABASE_URL. Just update the password if needed!
:::

## Test Your API

Your API is now running at `http://localhost:8080`. Test it:

### Create a Customer

```bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "name": "John Doe",
    "phone": "+1234567890"
  }'
```

### Get All Customers

```bash
curl http://localhost:8080/api/v1/customers
```

### Health Check

```bash
curl http://localhost:8080/health
```

## Template Mode vs AI Mode

| Feature | Template Mode | AI Mode |
|---------|--------------|---------|
| **Setup Required** | ❌ None | ✅ API Key |
| **Generation Speed** | ⚡ Instant | 🔄 2-5 seconds |
| **Use Case** | Standard CRUD entities | Complex business logic |
| **Field Types** | Basic types (string, int, etc.) | Smart types + validation |
| **Value Objects** | ❌ Not included | ✅ Auto-generated |
| **Business Logic** | Basic CRUD | Domain-specific methods |
| **Natural Language** | ❌ No | ✅ Yes |
| **Cost** | 🆓 Free | 🆓 Free tier available |

::: tip When to Use Each Mode
- **Template Mode**: Perfect for quick prototyping, standard entities, and learning DDD patterns
- **AI Mode**: Best for complex domains, business-specific validation, and production-ready code
:::

## What's Next?

- Learn about [Architecture](/guide/architecture)
- Explore [AI-Powered Generation](/guide/ai-generation) (optional)
- Read the [Command Reference](/reference/commands)
- Check out [Examples](/examples/basic)
- Try the **Search feature** (press `Ctrl+K` or `Cmd+K`)

::: tip Pro Tip
The interactive menu has a search feature! Press `/` to filter commands and find what you need quickly.
:::
