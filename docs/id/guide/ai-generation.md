# AI-Powered Generation (Opsional)

::: info AI Bersifat Opsional di v0.4.0!
Anaphase sekarang memiliki **dua mode**:
- **Template Mode**: Bekerja langsung tanpa AI - sempurna untuk CRUD standar
- **AI Mode**: Menggunakan LLM untuk generasi cerdas - panduan ini mencakup AI Mode

**Anda tidak perlu AI untuk menggunakan Anaphase!** Template Mode bekerja dengan baik untuk sebagian besar use case.
:::

Ketika Anda mengkonfigurasi AI provider, Anaphase dapat memahami natural language dan menghasilkan kode yang cerdas dan context-aware dengan validasi tingkat lanjut dan logika bisnis.

**AI Provider yang Didukung:**
- Google Gemini (rekomendasi, free tier yang generous)
- OpenAI (GPT-4, GPT-3.5-turbo)
- Anthropic Claude (Claude 3.5 Sonnet)
- Groq (inference cepat, free tier)

## Cara Kerjanya

### 1. Natural Language Input

Deskripsikan domain Anda dalam bahasa Inggris sederhana:

```bash
anaphase gen domain --name product --prompt \
  "Product with SKU code, name, description, price in USD,
   inventory quantity, and category. Products can be active or discontinued."
```

### 2. AI Processing

AI menganalisis prompt Anda dan mengidentifikasi:

- **Entities**: Product
- **Value Objects**: Money (untuk price), SKU
- **Fields**: name, description, quantity, category, status
- **Business Rules**: Status active/discontinued
- **Validation**: Format SKU, price > 0, quantity >= 0

### 3. Code Generation

Menghasilkan kode Go yang lengkap dan dapat dikompilasi:

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

Anaphase menggunakan system prompts yang dibuat dengan hati-hati yang mengajarkan AI tentang:

### Domain-Driven Design

AI memahami konsep DDD:

```
You are a Senior Golang Architect specializing in Domain-Driven Design.

Generate code following these patterns:
- Entities: Objects with identity (ID, CreatedAt, UpdatedAt)
- Value Objects: Immutable objects without identity
- Aggregates: Cluster of entities treated as a unit
- Repositories: Interfaces for persistence
```

### Code Structure

AI mengetahui konvensi Go:

```
- Use proper package names
- Follow Go naming conventions (PascalCase, camelCase)
- Add validation in constructors
- Use error wrapping with fmt.Errorf
- Add godoc comments
```

### Best Practices

AI menghasilkan kode production-ready:

```
- Add input validation
- Use value objects for important concepts
- Keep entities focused
- Use interfaces for dependencies
- Add proper error handling
```

## Contoh Prompt

### Simple Entity

```bash
anaphase gen domain --name user --prompt \
  "User with email and full name"
```

Menghasilkan:
- User entity dengan email (value object), name
- Email value object dengan validasi
- Basic repository methods

### Complex Entity

```bash
anaphase gen domain --name order --prompt \
  "Order with customer reference, multiple line items containing
   products and quantities, shipping address, billing address,
   total amount, and order status (pending, confirmed, shipped, delivered, cancelled)"
```

Menghasilkan:
- Order entity (aggregate root)
- LineItem entity
- Address value object
- Money value object
- OrderStatus enum
- Logika bisnis untuk status transitions

### With Business Rules

```bash
anaphase gen domain --name account --prompt \
  "Bank account with account number, balance, and account type (checking, savings).
   Balance cannot go negative. Savings accounts have interest rate."
```

Menghasilkan:
- Account entity dengan type
- Balance value object dengan validasi
- Business rules yang diterapkan di methods
- Repository dengan FindByAccountNumber

### Multiple Related Entities

```bash
anaphase gen domain --name blog --prompt \
  "Blog post with title, content, author reference, published date,
   and tags. Posts can have multiple comments from users."
```

Menghasilkan:
- Post entity (aggregate root)
- Comment entity (bagian dari aggregate)
- Tag value object
- Repository untuk posts (tidak untuk comments, mereka bagian dari aggregate)

## Memahami AI Output

### Apa yang Dihasilkan

Untuk `anaphase gen domain --name customer --prompt "Customer with email and name"`:

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

### Apa yang Tidak Dihasilkan

- Service implementations (Anda menulis logika bisnis)
- Handler implementations (gunakan `anaphase gen handler`)
- Repository implementations (gunakan `anaphase gen repository`)

Pemisahan ini memungkinkan Anda:
- Fokus AI pada domain modeling
- Implementasikan logika bisnis kompleks sendiri
- Gunakan templates untuk kode infrastructure

## Kustomisasi Generasi

### Temperature

Kontrol kreativitas vs konsistensi:

```bash
# More creative (may deviate from patterns)
anaphase gen domain --name product --prompt "..." --temperature 0.8

# More consistent (default)
anaphase gen domain --name product --prompt "..." --temperature 0.3

# Very strict (less variation)
anaphase gen domain --name product --prompt "..." --temperature 0.1
```

**Rekomendasi:**
- **0.1-0.3**: Pola yang konsisten (rekomendasi)
- **0.4-0.6**: Seimbang
- **0.7-1.0**: Variasi kreatif

### Caching

Anaphase menyimpan cache dari respons AI untuk menghemat waktu dan quota API:

```bash
# First call - hits AI API
anaphase gen domain --name user --prompt "User with email"

# Second call - uses cache
anaphase gen domain --name user --prompt "User with email"
```

Lokasi cache: `~/.anaphase/cache/`

Bersihkan cache:
```bash
rm -rf ~/.anaphase/cache
```

Nonaktifkan cache:
```yaml
# ~/.anaphase/config.yaml
cache:
  enabled: false
```

## Best Practices

### Tulis Prompt yang Jelas

Bagus:
```bash
"Customer with email address, full name, and phone number.
 Customers can have a billing address and shipping address."
```

Kurang optimal:
```bash
"customer stuff"
```

### Sebutkan Detail Penting

Sertakan:
- Nama field dan tipe
- Aturan validasi
- Relasi ke entities lain
- Business rules
- Status/state jika ada

Contoh:
```bash
"Invoice with invoice number (unique), customer reference,
 line items with products and quantities, subtotal, tax amount,
 total amount, and status (draft, sent, paid, overdue).
 Invoices can't be edited once sent."
```

### Gunakan Domain Language

Gunakan istilah dari domain bisnis Anda:

```bash
# E-commerce
"Product with SKU, price, inventory"

# Healthcare
"Patient with medical record number, diagnosis history"

# Finance
"Transaction with amount, currency, timestamp, type (debit/credit)"
```

### Iterasi dan Perbaiki

Mulai sederhana, lalu regenerate dengan lebih banyak detail:

```bash
# First iteration
anaphase gen domain --name order --prompt "Order with products"

# Refined
anaphase gen domain --name order --prompt \
  "Order with customer, line items (product, quantity, price),
   shipping address, payment status, fulfillment status"
```

## Konfigurasi AI Provider

Konfigurasi menggunakan CLI atau config file:

### Menggunakan CLI (Termudah)

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

Dapatkan API key: [Google AI Studio](https://makersuite.google.com/app/apikey)

### OpenAI

```yaml
ai:
  primary:
    type: openai
    apiKey: YOUR_API_KEY
    model: gpt-4o-mini
    timeout: 30s
```

Dapatkan API key: [OpenAI Platform](https://platform.openai.com/api-keys)

### Anthropic Claude

```yaml
ai:
  primary:
    type: claude
    apiKey: YOUR_API_KEY
    model: claude-3-5-sonnet-20241022
    timeout: 30s
```

Dapatkan API key: [Anthropic Console](https://console.anthropic.com/)

### Groq

```yaml
ai:
  primary:
    type: groq
    apiKey: YOUR_API_KEY
    model: llama-3.3-70b-versatile
    timeout: 30s
```

Dapatkan API key: [Groq Console](https://console.groq.com/)

### Konfigurasi Fallback

Setup backup providers:

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

Jika primary gagal (quota exceeded, network error), otomatis fallback ke secondary.

## Troubleshooting

### API Quota Exceeded

```
Error: quota exceeded
```

**Solusi:**
1. Tunggu (quota reset per menit)
2. Gunakan API key yang berbeda (secondary provider)
3. Enable caching untuk mengurangi API calls
4. Upgrade ke paid tier

### Invalid Response Format

```
Error: failed to parse AI response
```

**Solusi:**
1. Coba lagi (AI kadang output JSON yang invalid)
2. Kurangi temperature untuk output yang lebih konsisten
3. Periksa prompt untuk karakter yang tidak biasa

### Poor Code Quality

Jika kode yang dihasilkan tidak sesuai kebutuhan Anda:

1. **Lebih spesifik dalam prompt**
   ```bash
   # Vague
   "user with profile"

   # Specific
   "user with email (validated), display name,
    profile picture URL, and bio (max 500 chars)"
   ```

2. **Gunakan terminologi domain**
   ```bash
   # Generic
   "item with price"

   # Domain-specific
   "product with SKU, retail price in USD, and wholesale price"
   ```

3. **Sertakan business rules**
   ```bash
   # No rules
   "order with status"

   # With rules
   "order with status (pending -> confirmed -> shipped -> delivered).
    Orders can be cancelled only when pending or confirmed."
   ```

## AI Mode vs Template Mode

| Fitur | Template Mode | AI Mode |
|---------|--------------|---------|
| **Setup** | Tidak diperlukan | API key diperlukan |
| **Kecepatan** | Instan | 2-5 detik |
| **Input** | Nama entity + fields | Natural language description |
| **Value Objects** | ❌ Tidak dihasilkan | ✅ Auto-generated |
| **Validation** | Basic (type checking) | Advanced (business rules) |
| **Business Logic** | Standard CRUD | Domain-specific methods |
| **Relationships** | Manual | Terdeteksi dari description |
| **Cost** | Gratis | Free tier tersedia |
| **Use Case** | Entities standar | Domain kompleks |

### Kapan Menggunakan AI Mode

✅ **Gunakan AI Mode ketika:**
- Logika bisnis dan aturan validasi yang kompleks
- Membutuhkan value objects dengan validasi cerdas
- Ingin nama method yang spesifik untuk bisnis
- Berurusan dengan konsep yang spesifik untuk domain
- Membutuhkan deteksi relationship

✅ **Gunakan Template Mode ketika:**
- Entities CRUD sederhana
- Prototyping dengan cepat
- Model data standar
- Belajar pola DDD
- Tidak ada API key yang tersedia

## Next Steps

- [Quick Start](/guide/quick-start) - Coba kedua mode
- [DDD Concepts](/guide/ddd) - Pelajari DDD secara mendalam
- [Command Reference](/reference/gen-domain) - Opsi command lengkap
- [Examples](/examples/basic) - Lihat contoh real-world
