# anaphase gen domain

Generate model domain (entities, value objects, repository ports, service ports) dengan atau tanpa AI.

::: info
**Quick Start**: Jalankan `anaphase` (tanpa argumen) untuk mengakses menu interaktif di mana Anda dapat memilih "Generate Domain" dengan interface visual.
:::

## Synopsis

```bash
anaphase gen domain "<description>" [flags]
anaphase gen domain --interactive
```

## Deskripsi

Generate komponen domain-driven design:

- **Entities**: Entity domain dengan field, constructor, validasi, dan metode bisnis
- **Value Objects**: Objek immutable untuk konsep penting
- **Repository Interface**: Port untuk persistensi data
- **Service Interface**: Port untuk logika bisnis

Semua kode yang dihasilkan mengikuti prinsip Domain-Driven Design (DDD) dan Clean Architecture.

::: info
**AI Bersifat Opsional**: Command ini bekerja dalam dua mode:
- **Template Mode**: Generate kode yang bersih dan berfungsi dari template (tidak perlu AI)
- **AI Mode**: Menggunakan AI untuk menganalisis requirement Anda dan generate kode yang disesuaikan
:::

## Mode Penggunaan

### 1. Menu Interaktif (Disarankan)

Luncurkan menu visual:

```bash
anaphase
```

Kemudian pilih **"Generate Domain"** dari menu. Interface akan memandu Anda melalui:
- Deskripsi domain
- Pemilihan mode Template vs AI
- Pemilihan AI provider (jika menggunakan AI mode)
- Konfigurasi direktori output

### 2. Mode Langsung

Berikan deskripsi sebagai argumen:

```bash
anaphase gen domain "User with email, name, and password"
```

### 3. Mode CLI Interaktif

Gunakan prompt terpandu untuk input:

```bash
anaphase gen domain --interactive
```

**Prompt Interaktif:**
1. **Deskripsi domain** - Requirement bisnis Anda
2. **AI provider** - Pilih dari provider yang tersedia (gemini, groq, openai, claude)
3. **Direktori output** - Di mana generate file (default: internal/core)

**Contoh Sesi:**
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
**Zero Configuration**: Template mode generate kode production-ready tanpa memerlukan setup AI atau API key.
:::

Template mode membuat kode domain yang bersih dan berfungsi berdasarkan pattern yang sudah terbukti:

```bash
# Template mode adalah default ketika tidak ada AI provider yang dikonfigurasi
anaphase gen domain "User with email and password"
```

**Yang Dihasilkan Template Mode:**

```go
// Entity dengan validasi
type User struct {
    ID        uuid.UUID
    Email     string
    Password  string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Constructor dengan validasi
func NewUser(email, password string) (*User, error) {
    if email == "" {
        return nil, ErrInvalidEmail
    }
    // ... logika validasi
}

// Repository interface
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id uuid.UUID) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

**Keuntungan:**
- Tidak perlu API key
- Generate instan
- Struktur kode yang bersih dan dapat diprediksi
- Prinsip DDD dan Clean Architecture
- Siap untuk disesuaikan

## Flag

| Flag | Short | Default | Deskripsi |
|------|-------|---------|-------------|
| `--interactive` | `-i` | `false` | Jalankan dalam mode interaktif dengan prompt terpandu |
| `--provider` | | (config) | AI provider: gemini, groq, openai, claude (opsional) |
| `--output` | | `internal/core` | Direktori output untuk file yang dihasilkan |

## Flag Global

| Flag | Short | Deskripsi |
|------|-------|-------------|
| `--debug` | `-d` | Aktifkan debug mode dengan verbose logging |
| `--verbose` | `-v` | Aktifkan output verbose |

## Contoh

### Quick Start dengan Menu Interaktif

```bash
# Luncurkan menu interaktif
anaphase

# Navigasi ke "Generate Domain" dan ikuti prompt
# Menu menyediakan interface visual untuk semua opsi
```

### Penggunaan Dasar (Template Mode)

```bash
# Tidak perlu AI - generate instan
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

### Dengan AI Provider (Opsional)

::: info
AI provider bersifat opsional. Tanpa konfigurasi, tool menggunakan template mode.
:::

```bash
# Gunakan Groq (opsi AI tercepat)
anaphase gen domain "User with email" --provider groq

# Gunakan OpenAI (opsi AI paling akurat)
anaphase gen domain "Order processing system" --provider openai

# Gunakan Gemini (opsi AI gratis)
anaphase gen domain "Product catalog" --provider gemini
```

### Mode Interaktif

```bash
anaphase gen domain -i
# atau
anaphase gen domain --interactive
```

Keuntungan:
- Prompt terpandu untuk semua input
- Pemilihan provider dengan deskripsi
- Saran nilai default
- Validasi input

### Direktori Output Kustom

```bash
anaphase gen domain "User" --output pkg/domain
```

### Deskripsi Domain Kompleks

```bash
anaphase gen domain "
Order with ID, Total, Status, Items.
Customer can place order, cancel if pending.
Status can be: pending, confirmed, shipped, delivered, cancelled.
Each Item has product reference, quantity, and price.
"
```

## Pemilihan AI Provider

Anda dapat override provider yang dikonfigurasi:

```bash
# Cek provider yang tersedia
anaphase config show-providers

# Gunakan provider tertentu
anaphase gen domain "User" --provider groq

# Set default provider
anaphase config set-provider groq
```

**Perbandingan Provider:**

| Provider | Kecepatan | Kualitas | Biaya | Terbaik Untuk |
|----------|-------|---------|------|----------|
| **Gemini** | ⚡⚡⚡ | ⭐⭐⭐⭐ | Gratis | Penggunaan umum, pilihan default |
| **Groq** | ⚡⚡⚡⚡⚡ | ⭐⭐⭐ | Gratis | Kecepatan kritis, real-time |
| **OpenAI** | ⚡⚡⚡ | ⭐⭐⭐⭐⭐ | Berbayar | Domain kompleks, akurasi |
| **Claude** | ⚡⚡⚡ | ⭐⭐⭐⭐⭐ | Berbayar | Konteks besar |

## Struktur Kode yang Dihasilkan

### Contoh Entity

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

### Contoh Repository Interface

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

### Contoh Service Interface

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

## Menulis Deskripsi yang Baik

### ✅ Deskripsi Baik

```bash
# Jelas, spesifik, actionable
"User with email, password, and profile picture. Can login and update profile."

# Mencakup aturan bisnis
"Order with items and total. Status: pending, confirmed, shipped. Can be cancelled if pending."

# Menyebutkan relasi
"Cart belongs to User. Cart has many Items. Each Item references a Product."
```

### ❌ Deskripsi Buruk

```bash
# Terlalu samar
"User system"

# Kurang detail
"Order"

# Implementasi teknis (bukan domain bisnis)
"Create a struct with fields id, name, email and CRUD methods"
```

## Tips untuk Hasil Terbaik

1. **Bersikap Spesifik**: Sertakan nama field, tipe, dan aturan bisnis
2. **Jelaskan Perilaku**: Sebutkan apa yang dapat dilakukan user (add, remove, update, dll.)
3. **Sertakan Validasi**: Tentukan constraint dan validasi
4. **Sebutkan Relasi**: Jelaskan bagaimana entity berhubungan satu sama lain
5. **Gunakan Bahasa Bisnis**: Fokus pada konsep domain, bukan implementasi teknis

## Langkah Selanjutnya

Setelah generate kode domain:

1. **Review Kode yang Dihasilkan**
   ```bash
   ls -la internal/core/
   ```

2. **Validasi Kualitas Kode**
   ```bash
   anaphase quality validate
   ```

3. **Generate Implementasi Repository**
   ```bash
   anaphase gen repository Cart
   ```

4. **Generate HTTP Handler**
   ```bash
   anaphase gen handler Cart
   ```

5. **Build dan Test**
   ```bash
   go build ./...
   go test ./...
   ```

## Troubleshooting

### "No AI providers configured"

```bash
# Set API key
export GEMINI_API_KEY="your-key-here"

# Verifikasi
anaphase config check
```

### "AI generation failed"

```bash
# Coba provider berbeda
anaphase gen domain "User" --provider groq

# Cek kesehatan provider
anaphase config check

# Aktifkan debug mode
anaphase gen domain "User" --debug
```

### Kode yang dihasilkan memiliki error

```bash
# Jalankan quality check
anaphase quality lint --fix
anaphase quality format
anaphase quality validate
```

## Lihat Juga

- [anaphase gen handler](/reference/gen-handler) - Generate HTTP handler
- [anaphase gen repository](/reference/gen-repository) - Generate implementasi repository
- [anaphase gen middleware](/reference/gen-middleware) - Generate middleware
- [anaphase config](/reference/config) - Konfigurasi AI provider
- [AI Providers](/config/ai-providers) - Panduan setup provider
- [Domain-Driven Design](/guide/ddd) - Konsep dan pattern DDD
