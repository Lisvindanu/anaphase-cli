# Apa itu Anaphase?

Anaphase adalah **AI-powered CLI tool** yang menghasilkan Golang microservices production-ready dengan mengikuti best practices Domain-Driven Design (DDD) dan Clean Architecture.

## Masalahnya

Membangun microservices melibatkan menulis banyak boilerplate code yang repetitif:

- Definisi entity dengan validasi
- Interface dan implementasi repository
- Service layer dengan business logic
- HTTP handlers dan DTOs
- Database schemas dan migrations
- Dependency injection
- Tests untuk semua layer

Proses ini:
- **Memakan waktu**: Berjam-jam setup untuk setiap domain
- **Rawan error**: Mudah melewatkan pattern atau membuat kesalahan
- **Tidak konsisten**: Developer berbeda, style berbeda
- **Membosankan**: Pattern yang sama berulang kali

## Solusinya

Anaphase menggunakan AI untuk memahami kebutuhan domain Anda dan menghasilkan semua code yang diperlukan secara otomatis, mengikuti pattern yang sudah established.

### Fitur Utama

#### 🤖 AI-Powered Generation

Deskripsikan domain Anda dalam bahasa natural, dan Anaphase akan menghasilkan kode Go yang lengkap dan bisa dikompilasi:

```bash
anaphase gen domain --name order --prompt \
  "Order dengan referensi customer, line items dengan produk dan kuantitas,
   total amount, dan status (pending, confirmed, shipped, delivered)"
```

AI memahami:
- **Entities** vs **Value Objects**
- **Aggregates** dan boundariesnya
- **Business rules** dan validasi
- **Relationships** antar domains

#### 🏗️ Clean Architecture

Semua code yang dihasilkan mengikuti prinsip Clean Architecture:

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

**Keuntungan:**
- Independen dari frameworks
- Testable di setiap layer
- Infrastruktur yang fleksibel
- Business logic terlindungi

#### ⚡ Rapid Development

Dari nol ke running API dalam hitungan menit:

| Task | Tradisional | Dengan Anaphase |
|------|-------------|---------------|
| Setup project | 30 menit | 10 detik |
| Domain model | 2 jam | 30 detik |
| Repository | 1 jam | 10 detik |
| Handler | 1 jam | 10 detik |
| Wiring | 30 menit | 5 detik |
| **Total** | **~5 jam** | **~1 menit** |

#### 🎯 Type-Safe

Strong typing di seluruh codebase:

```go
// Value Objects dengan validasi
type Email struct {
    value string
}

func NewEmail(value string) (*Email, error) {
    if !isValidEmail(value) {
        return nil, ErrInvalidEmail
    }
    return &Email{value: value}, nil
}

// Entities dengan business logic
type Customer struct {
    ID        uuid.UUID
    Email     *Email  // Type-safe value object
    Name      string
    CreatedAt time.Time
}
```

## Cara Kerjanya

### 1. AST Analysis

Anaphase menggunakan Go AST parser untuk menganalisa codebase Anda dan menemukan existing domains secara otomatis.

### 2. AI Generation

Memanfaatkan Google Gemini dengan prompt yang di-engineer dengan hati-hati untuk menghasilkan:
- Domain models mengikuti DDD
- Repository patterns
- Service interfaces
- Handler implementations

### 3. Code Generation

Menghasilkan kode Go production-ready:
- Import dan package yang proper
- Error handling
- Validasi
- Tests
- Dokumentasi

### 4. Auto-Wiring

Scan generated code dan otomatis wire:
- Dependencies
- Database connections
- HTTP routes
- Middleware

## Pattern Arsitektur

Anaphase menerapkan best practices:

### Domain-Driven Design (DDD)

- **Entities**: Object dengan identity
- **Value Objects**: Immutable object tanpa identity
- **Aggregates**: Cluster dari entities dan value objects
- **Repositories**: Abstraksi persistence
- **Services**: Koordinasi business logic

### Hexagonal Architecture

- **Ports**: Interface yang mendefinisikan contract
- **Adapters**: Implementasi dari ports
- **Core**: Business logic terisolasi dari infrastructure

### SOLID Principles

- Single Responsibility
- Open/Closed
- Liskov Substitution
- Interface Segregation
- Dependency Inversion

## Use Cases

### Microservices

Generate multiple bounded contexts dengan cepat:

```bash
anaphase gen domain --name user
anaphase gen domain --name product
anaphase gen domain --name order
anaphase wire
```

### API Backends

REST API lengkap dengan CRUD operations:

```bash
anaphase gen domain --name article --prompt "Blog article dengan title, content, author"
anaphase gen handler --domain article
anaphase gen repository --domain article
anaphase wire
```

### Database Migrations

Ketika Anda perlu mulai fresh atau ubah schemas:

```bash
anaphase gen repository --domain customer --db postgres
# Apply internal/adapter/repository/postgres/schema.sql
```

## Perbandingan

| Fitur | Anaphase | Manual | Generator Lain |
|---------|----------|--------|------------------|
| AI-Powered | ✅ | ❌ | ❌ |
| DDD Support | ✅ | Tergantung | ⚠️ |
| Clean Architecture | ✅ | Tergantung | ⚠️ |
| Auto-Wiring | ✅ | ❌ | ⚠️ |
| Type-Safe | ✅ | Tergantung | ⚠️ |
| Multi-DB | ✅ | ❌ | ✅ |
| Production-Ready | ✅ | Tergantung | ⚠️ |
| Learning Curve | Rendah | Tinggi | Medium |

## Filosofi

Anaphase dibangun berdasarkan prinsip-prinsip ini:

1. **AI-Assisted, Bukan AI-Replaced**: AI membantu task repetitif, Anda fokus ke business logic
2. **Patterns Over Configuration**: Enforce best practices by default
3. **Flexibility**: Generated code adalah milik Anda untuk dimodifikasi
4. **Transparency**: Lihat persis apa yang di-generate, tidak ada magic
5. **Production-First**: Generate code yang benar-benar Anda deploy

## Selanjutnya Apa?

- [Mulai Cepat](/guide/quick-start) - Build service pertama Anda
- [Instalasi](/guide/installation) - Setup detail
- [Architecture](/guide/architecture) - Deep dive ke pattern yang digunakan
