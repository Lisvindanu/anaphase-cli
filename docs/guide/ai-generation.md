# AI-Powered Generation (Optional)

::: info AI is Optional in v0.4.0!
Anaphase now has **two modes**:
- **Template Mode**: Works immediately without AI - perfect for standard CRUD
- **AI Mode**: Uses LLMs for intelligent generation - this guide covers AI Mode

**You don't need AI to use Anaphase!** Template Mode works great for most use cases.
:::

When you configure an AI provider, Anaphase can understand natural language and generate intelligent, context-aware code with advanced validation and business logic.

**Supported AI Providers:**
- Google Gemini (recommended, generous free tier)
- OpenAI (GPT-4, GPT-3.5-turbo)
- Anthropic Claude (Claude 3.5 Sonnet)
- Groq (fast inference, free tier)

## How It Works

### 1. Natural Language Input

Describe your domain in plain English:

```bash
anaphase gen domain --name product --prompt \
  "Product with SKU code, name, description, price in USD,
   inventory quantity, and category. Products can be active or discontinued."
```

### 2. AI Processing

The AI analyzes your prompt and identifies:

- **Entities**: Product
- **Value Objects**: Money (for price), SKU
- **Fields**: name, description, quantity, category, status
- **Business Rules**: Active/discontinued status
- **Validation**: SKU format, price > 0, quantity >= 0

### 3. Code Generation

Generates complete, compilable Go code:

```go
// Entity
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

// Value Object
type SKU struct {
    value string
}

func NewSKU(value string) (*SKU, error) {
    if !isValidSKU(value) {
        return nil, ErrInvalidSKU
    }
    return &SKU{value: value}, nil
}

// Repository Interface
type ProductRepository interface {
    Save(ctx context.Context, product *entity.Product) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Product, error)
    FindBySKU(ctx context.Context, sku valueobject.SKU) (*entity.Product, error)
}
```

## AI Prompt Engineering

Anaphase uses carefully crafted system prompts that teach the AI about:

### Domain-Driven Design

The AI understands DDD concepts:

```
You are a Senior Golang Architect specializing in Domain-Driven Design.

Generate code following these patterns:
- Entities: Objects with identity (ID, CreatedAt, UpdatedAt)
- Value Objects: Immutable objects without identity
- Aggregates: Cluster of entities treated as a unit
- Repositories: Interfaces for persistence
```

### Code Structure

The AI knows Go conventions:

```
- Use proper package names
- Follow Go naming conventions (PascalCase, camelCase)
- Add validation in constructors
- Use error wrapping with fmt.Errorf
- Add godoc comments
```

### Best Practices

The AI generates production-ready code:

```
- Add input validation
- Use value objects for important concepts
- Keep entities focused
- Use interfaces for dependencies
- Add proper error handling
```

## Prompt Examples

### Simple Entity

```bash
anaphase gen domain --name user --prompt \
  "User with email and full name"
```

Generates:
- User entity with email (value object), name
- Email value object with validation
- Basic repository methods

### Complex Entity

```bash
anaphase gen domain --name order --prompt \
  "Order with customer reference, multiple line items containing
   products and quantities, shipping address, billing address,
   total amount, and order status (pending, confirmed, shipped, delivered, cancelled)"
```

Generates:
- Order entity (aggregate root)
- LineItem entity
- Address value object
- Money value object
- OrderStatus enum
- Business logic for status transitions

### With Business Rules

```bash
anaphase gen domain --name account --prompt \
  "Bank account with account number, balance, and account type (checking, savings).
   Balance cannot go negative. Savings accounts have interest rate."
```

Generates:
- Account entity with type
- Balance value object with validation
- Business rules enforced in methods
- Repository with FindByAccountNumber

### Multiple Related Entities

```bash
anaphase gen domain --name blog --prompt \
  "Blog post with title, content, author reference, published date,
   and tags. Posts can have multiple comments from users."
```

Generates:
- Post entity (aggregate root)
- Comment entity (part of aggregate)
- Tag value object
- Repository for posts (not comments, they're part of aggregate)

## Understanding AI Output

### What Gets Generated

For `anaphase gen domain --name customer --prompt "Customer with email and name"`:

1. **Entity** (`entity/customer.go`)
```go
type Customer struct {
    ID        uuid.UUID
    Email     *valueobject.Email
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

func NewCustomer(email *valueobject.Email, name string) (*Customer, error)
func (c *Customer) Validate() error
```

2. **Value Objects** (`valueobject/email.go`)
```go
type Email struct {
    value string
}

func NewEmail(value string) (*Email, error)
func (e *Email) String() string
func (e *Email) Validate() error
```

3. **Repository Port** (`port/customer_repo.go`)
```go
type CustomerRepository interface {
    Save(ctx context.Context, customer *entity.Customer) error
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Customer, error)
    FindByEmail(ctx context.Context, email valueobject.Email) (*entity.Customer, error)
}
```

4. **Service Port** (`port/customer_service.go`)
```go
type CustomerService interface {
    CreateCustomer(ctx context.Context, email, name string) (*entity.Customer, error)
    GetCustomer(ctx context.Context, id uuid.UUID) (*entity.Customer, error)
}
```

### What Doesn't Get Generated

- Service implementations (you write business logic)
- Handler implementations (use `anaphase gen handler`)
- Repository implementations (use `anaphase gen repository`)

This separation lets you:
- Focus AI on domain modeling
- Implement complex business logic yourself
- Use templates for infrastructure code

## Customizing Generation

### Temperature

Control creativity vs consistency:

```bash
# More creative (may deviate from patterns)
anaphase gen domain --name product --prompt "..." --temperature 0.8

# More consistent (default)
anaphase gen domain --name product --prompt "..." --temperature 0.3

# Very strict (less variation)
anaphase gen domain --name product --prompt "..." --temperature 0.1
```

**Recommendations:**
- **0.1-0.3**: Consistent patterns (recommended)
- **0.4-0.6**: Balanced
- **0.7-1.0**: Creative variations

### Caching

Anaphase caches AI responses to save time and API quota:

```bash
# First call - hits AI API
anaphase gen domain --name user --prompt "User with email"

# Second call - uses cache
anaphase gen domain --name user --prompt "User with email"
```

Cache location: `~/.anaphase/cache/`

Clear cache:
```bash
rm -rf ~/.anaphase/cache
```

Disable cache:
```yaml
# ~/.anaphase/config.yaml
cache:
  enabled: false
```

## Best Practices

### Write Clear Prompts

Good:
```bash
"Customer with email address, full name, and phone number.
 Customers can have a billing address and shipping address."
```

Less optimal:
```bash
"customer stuff"
```

### Specify Important Details

Include:
- Field names and types
- Validation rules
- Relationships to other entities
- Business rules
- Status/state if applicable

Example:
```bash
"Invoice with invoice number (unique), customer reference,
 line items with products and quantities, subtotal, tax amount,
 total amount, and status (draft, sent, paid, overdue).
 Invoices can't be edited once sent."
```

### Use Domain Language

Use terms from your business domain:

```bash
# E-commerce
"Product with SKU, price, inventory"

# Healthcare
"Patient with medical record number, diagnosis history"

# Finance
"Transaction with amount, currency, timestamp, type (debit/credit)"
```

### Iterate and Refine

Start simple, then regenerate with more details:

```bash
# First iteration
anaphase gen domain --name order --prompt "Order with products"

# Refined
anaphase gen domain --name order --prompt \
  "Order with customer, line items (product, quantity, price),
   shipping address, payment status, fulfillment status"
```

## AI Provider Configuration

Configure using the CLI or config file:

### Using CLI (Easiest)

```bash
# Set provider interactively
anaphase config set-provider

# Or directly
anaphase config set-provider gemini
anaphase config set-provider openai
anaphase config set-provider claude
anaphase config set-provider groq
```

### Google Gemini

```yaml
# ~/.anaphase/config.yaml
ai:
  primary:
    type: gemini
    apiKey: YOUR_API_KEY
    model: gemini-2.0-flash-exp
    timeout: 30s
```

Get API key: [Google AI Studio](https://makersuite.google.com/app/apikey)

### OpenAI

```yaml
ai:
  primary:
    type: openai
    apiKey: YOUR_API_KEY
    model: gpt-4o-mini
    timeout: 30s
```

Get API key: [OpenAI Platform](https://platform.openai.com/api-keys)

### Anthropic Claude

```yaml
ai:
  primary:
    type: claude
    apiKey: YOUR_API_KEY
    model: claude-3-5-sonnet-20241022
    timeout: 30s
```

Get API key: [Anthropic Console](https://console.anthropic.com/)

### Groq

```yaml
ai:
  primary:
    type: groq
    apiKey: YOUR_API_KEY
    model: llama-3.3-70b-versatile
    timeout: 30s
```

Get API key: [Groq Console](https://console.groq.com/)

### Fallback Configuration

Set up backup providers:

```yaml
ai:
  primary:
    type: gemini
    apiKey: PRIMARY_KEY
    model: gemini-2.0-flash-exp

  secondary:
    type: openai
    apiKey: BACKUP_KEY
    model: gpt-4o-mini
```

If primary fails (quota exceeded, network error), automatically falls back to secondary.

## Troubleshooting

### API Quota Exceeded

```
Error: quota exceeded
```

**Solutions:**
1. Wait (quota resets per minute)
2. Use a different API key (secondary provider)
3. Enable caching to reduce API calls
4. Upgrade to paid tier

### Invalid Response Format

```
Error: failed to parse AI response
```

**Solutions:**
1. Try again (AI sometimes outputs invalid JSON)
2. Lower temperature for more consistent output
3. Check prompt for unusual characters

### Poor Code Quality

If generated code doesn't match your needs:

1. **Be more specific in prompt**
   ```bash
   # Vague
   "user with profile"

   # Specific
   "user with email (validated), display name,
    profile picture URL, and bio (max 500 chars)"
   ```

2. **Use domain terminology**
   ```bash
   # Generic
   "item with price"

   # Domain-specific
   "product with SKU, retail price in USD, and wholesale price"
   ```

3. **Include business rules**
   ```bash
   # No rules
   "order with status"

   # With rules
   "order with status (pending -> confirmed -> shipped -> delivered).
    Orders can be cancelled only when pending or confirmed."
   ```

## AI Mode vs Template Mode

| Feature | Template Mode | AI Mode |
|---------|--------------|---------|
| **Setup** | None required | API key needed |
| **Speed** | Instant | 2-5 seconds |
| **Input** | Entity name + fields | Natural language description |
| **Value Objects** | ❌ Not generated | ✅ Auto-generated |
| **Validation** | Basic (type checking) | Advanced (business rules) |
| **Business Logic** | Standard CRUD | Domain-specific methods |
| **Relationships** | Manual | Detected from description |
| **Cost** | Free | Free tier available |
| **Use Case** | Standard entities | Complex domains |

### When to Use AI Mode

✅ **Use AI Mode when:**
- Complex business logic and validation rules
- Need value objects with smart validation
- Want business-specific method names
- Dealing with domain-specific concepts
- Need relationship detection

✅ **Use Template Mode when:**
- Simple CRUD entities
- Prototyping quickly
- Standard data models
- Learning DDD patterns
- No API key available

## Next Steps

- [Quick Start](/guide/quick-start) - Try both modes
- [DDD Concepts](/guide/ddd) - Learn DDD in depth
- [Command Reference](/reference/gen-domain) - Full command options
- [Examples](/examples/basic) - See real-world examples
