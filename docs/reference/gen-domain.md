# anaphase gen domain

Generate domain models (entities, value objects, repository ports, service ports) with or without AI.

::: info
**Quick Start**: Run `anaphase` (no arguments) to access the interactive menu where you can select "Generate Domain" with a visual interface.
:::

## Synopsis

```bash
anaphase gen domain "<description>" [flags]
anaphase gen domain --interactive
```

## Description

Generates domain-driven design components:

- **Entities**: Domain entities with fields, constructors, validation, and business methods
- **Value Objects**: Immutable objects for important concepts
- **Repository Interface**: Port for data persistence
- **Service Interface**: Port for business logic

All generated code follows Domain-Driven Design (DDD) and Clean Architecture principles.

::: info
**AI is Optional**: This command works in two modes:
- **Template Mode**: Generates clean, working code from templates (no AI required)
- **AI Mode**: Uses AI to analyze your requirements and generate customized code
:::

## Usage Modes

### 1. Interactive Menu (Recommended)

Launch the visual menu:

```bash
anaphase
```

Then select **"Generate Domain"** from the menu. The interface guides you through:
- Domain description
- Template vs AI mode selection
- AI provider selection (if using AI mode)
- Output directory configuration

### 2. Direct Mode

Provide description as argument:

```bash
anaphase gen domain "User with email, name, and password"
```

### 3. Interactive CLI Mode

Use guided prompts for input:

```bash
anaphase gen domain --interactive
```

**Interactive Prompts:**
1. **Domain description** - Your business requirement
2. **AI provider** - Select from available providers (gemini, groq, openai, claude)
3. **Output directory** - Where to generate files (default: internal/core)

**Example Session:**
```
⚡ Interactive Domain Generation

Enter domain description: User with email and password. Can login and logout
Select AI provider:
  1) gemini (default)
  2) groq
  3) openai
  4) claude
Enter choice [1]: 2

Output directory [internal/core]:

⚡ AI-Powered Domain Generation
ℹ Description: User with email and password. Can login and logout
ℹ Using provider: groq
...
```

## Template Mode

::: info
**Zero Configuration**: Template mode generates production-ready code without requiring AI setup or API keys.
:::

Template mode creates clean, working domain code based on proven patterns:

```bash
# Template mode is the default when no AI provider is configured
anaphase gen domain "User with email and password"
```

**What Template Mode Generates:**

```go
// Entity with validation
type User struct {
    ID        uuid.UUID
    Email     string
    Password  string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Constructor with validation
func NewUser(email, password string) (*User, error) {
    if email == "" {
        return nil, ErrInvalidEmail
    }
    // ... validation logic
}

// Repository interface
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id uuid.UUID) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

**Benefits:**
- No API keys required
- Instant generation
- Clean, predictable code structure
- DDD and Clean Architecture principles
- Ready to customize

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--interactive` | `-i` | `false` | Run in interactive mode with guided prompts |
| `--provider` | | (config) | AI provider: gemini, groq, openai, claude (optional) |
| `--output` | | `internal/core` | Output directory for generated files |

## Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--debug` | `-d` | Enable debug mode with verbose logging |
| `--verbose` | `-v` | Enable verbose output |

## Examples

### Quick Start with Interactive Menu

```bash
# Launch the interactive menu
anaphase

# Navigate to "Generate Domain" and follow the prompts
# The menu provides a visual interface for all options
```

### Basic Usage (Template Mode)

```bash
# No AI required - instant generation
anaphase gen domain "Cart with Items. User can add, remove, update quantity"
```

Output:
```
⚡ AI-Powered Domain Generation
ℹ Description: Cart with Items. User can add, remove, update quantity

⚙️  Step 1/3: Loading configuration...
ℹ Using provider: gemini

🧠 Step 2/3: Analyzing with AI...
✓ AI Analysis Complete!

Generated Specification:
  📦 Domain: Cart
  📄 Entities: 2
  📄 Value Objects: 1
  ⚙️  Repository: CartRepository
  ⚙️  Service: CartService

📂 Step 3/3: Generating code files...

Generated Files:
✓ internal/core/entity/cart.go
✓ internal/core/entity/item.go
✓ internal/core/valueobject/quantity.go
✓ internal/core/port/cart_repository.go
✓ internal/core/port/cart_service.go

✓ Domain generation complete! 🚀
```

### With AI Provider (Optional)

::: info
AI providers are optional. Without configuration, the tool uses template mode.
:::

```bash
# Use Groq (fastest AI option)
anaphase gen domain "User with email" --provider groq

# Use OpenAI (most accurate AI option)
anaphase gen domain "Order processing system" --provider openai

# Use Gemini (free AI option)
anaphase gen domain "Product catalog" --provider gemini
```

### Interactive Mode

```bash
anaphase gen domain -i
# or
anaphase gen domain --interactive
```

Benefits:
- Guided prompts for all inputs
- Provider selection with descriptions
- Default value suggestions
- Validation of inputs

### Custom Output Directory

```bash
anaphase gen domain "User" --output pkg/domain
```

### Complex Domain Description

```bash
anaphase gen domain "
Order with ID, Total, Status, Items.
Customer can place order, cancel if pending.
Status can be: pending, confirmed, shipped, delivered, cancelled.
Each Item has product reference, quantity, and price.
"
```

## AI Provider Selection

You can override the configured provider:

```bash
# Check available providers
anaphase config show-providers

# Use specific provider
anaphase gen domain "User" --provider groq

# Set default provider
anaphase config set-provider groq
```

**Provider Comparison:**

| Provider | Speed | Quality | Cost | Best For |
|----------|-------|---------|------|----------|
| **Gemini** | ⚡⚡⚡ | ⭐⭐⭐⭐ | Free | General use, default choice |
| **Groq** | ⚡⚡⚡⚡⚡ | ⭐⭐⭐ | Free | Speed-critical, real-time |
| **OpenAI** | ⚡⚡⚡ | ⭐⭐⭐⭐⭐ | Paid | Complex domains, accuracy |
| **Claude** | ⚡⚡⚡ | ⭐⭐⭐⭐⭐ | Paid | Large contexts |

## Generated Code Structure

### Entity Example

```go
// internal/core/entity/cart.go
package entity

import (
    "errors"
    "time"
    "github.com/google/uuid"
    "yourproject/internal/core/valueobject"
)

var (
    ErrCartNotFound = errors.New("cart not found")
    ErrInvalidCart = errors.New("invalid cart")
)

// Cart is an aggregate root
type Cart struct {
    ID        uuid.UUID
    UserID    uuid.UUID
    Items     []Item
    Total     float64
    CreatedAt time.Time
    UpdatedAt time.Time
}

// NewCart creates a new cart
func NewCart() *Cart {
    return &Cart{
        ID:        uuid.New(),
        Items:     []Item{},
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

// AddItem adds an item to the cart
func (c *Cart) AddItem(item Item) error {
    // Business logic here
    return nil
}

// RemoveItem removes an item from the cart
func (c *Cart) RemoveItem(itemID uuid.UUID) error {
    // Business logic here
    return nil
}

// Validate validates the cart
func (c *Cart) Validate() error {
    if c.ID == uuid.Nil {
        return ErrInvalidCart
    }
    return nil
}
```

### Repository Interface Example

```go
// internal/core/port/cart_repository.go
package port

import (
    "context"
    "github.com/google/uuid"
    "yourproject/internal/core/entity"
)

// CartRepository defines the contract for cart persistence
type CartRepository interface {
    // Create creates a new cart
    Create(ctx context.Context, cart *entity.Cart) error

    // FindByID finds a cart by ID
    FindByID(ctx context.Context, id uuid.UUID) (*entity.Cart, error)

    // Update updates an existing cart
    Update(ctx context.Context, cart *entity.Cart) error

    // Delete deletes a cart
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### Service Interface Example

```go
// internal/core/port/cart_service.go
package port

import (
    "context"
    "github.com/google/uuid"
    "yourproject/internal/core/entity"
)

// CartService defines the contract for cart business logic
type CartService interface {
    // AddItemToCart adds an item to the cart
    AddItemToCart(ctx context.Context, cartID uuid.UUID, item entity.Item) error

    // RemoveItemFromCart removes an item from the cart
    RemoveItemFromCart(ctx context.Context, cartID uuid.UUID, itemID uuid.UUID) error

    // GetCart retrieves a cart by ID
    GetCart(ctx context.Context, id uuid.UUID) (*entity.Cart, error)
}
```

## Writing Good Descriptions

### ✅ Good Descriptions

```bash
# Clear, specific, actionable
"User with email, password, and profile picture. Can login and update profile."

# Includes business rules
"Order with items and total. Status: pending, confirmed, shipped. Can be cancelled if pending."

# Mentions relationships
"Cart belongs to User. Cart has many Items. Each Item references a Product."
```

### ❌ Poor Descriptions

```bash
# Too vague
"User system"

# Missing details
"Order"

# Technical implementation (not business domain)
"Create a struct with fields id, name, email and CRUD methods"
```

## Tips for Best Results

1. **Be Specific**: Include field names, types, and business rules
2. **Describe Behavior**: Mention what users can do (add, remove, update, etc.)
3. **Include Validations**: Specify constraints and validations
4. **Mention Relationships**: Describe how entities relate to each other
5. **Use Business Language**: Focus on domain concepts, not technical implementation

## Next Steps

After generating domain code:

1. **Review Generated Code**
   ```bash
   ls -la internal/core/
   ```

2. **Validate Code Quality**
   ```bash
   anaphase quality validate
   ```

3. **Generate Repository Implementation**
   ```bash
   anaphase gen repository Cart
   ```

4. **Generate HTTP Handlers**
   ```bash
   anaphase gen handler Cart
   ```

5. **Build and Test**
   ```bash
   go build ./...
   go test ./...
   ```

## Troubleshooting

### "No AI providers configured"

```bash
# Set API key
export GEMINI_API_KEY="your-key-here"

# Verify
anaphase config check
```

### "AI generation failed"

```bash
# Try different provider
anaphase gen domain "User" --provider groq

# Check provider health
anaphase config check

# Enable debug mode
anaphase gen domain "User" --debug
```

### Generated code has errors

```bash
# Run quality checks
anaphase quality lint --fix
anaphase quality format
anaphase quality validate
```

## See Also

- [anaphase gen handler](/reference/gen-handler) - Generate HTTP handlers
- [anaphase gen repository](/reference/gen-repository) - Generate repository implementation
- [anaphase gen middleware](/reference/gen-middleware) - Generate middleware
- [anaphase config](/reference/config) - Configure AI providers
- [AI Providers](/config/ai-providers) - Provider setup guide
- [Domain-Driven Design](/guide/ddd) - DDD concepts and patterns
