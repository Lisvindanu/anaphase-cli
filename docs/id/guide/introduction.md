# Apa itu Anaphase?

Anaphase adalah **tool CLI interaktif** yang menghasilkan microservice Golang production-ready mengikuti best practices dalam Domain-Driven Design (DDD) dan Clean Architecture.

**Bekerja dengan atau tanpa AI** - pilih Template Mode untuk scaffolding instan atau AI Mode untuk generation yang intelligent.

## Masalahnya

Membangun microservices melibatkan penulisan boilerplate code yang berulang:

- Definisi entity dengan validasi
- Interface repository dan implementasinya
- Service layer dengan business logic
- HTTP handlers dan DTOs
- Database schemas dan migrations
- Dependency injection
- Tests untuk semua layer

Proses ini:
- **Memakan waktu**: Berjam-jam setup untuk setiap domain
- **Rawan error**: Mudah melewatkan pattern atau membuat kesalahan
- **Tidak konsisten**: Developer berbeda, style berbeda
- **Membosankan**: Pattern yang sama berulang-ulang

## Solusinya

Anaphase menghasilkan semua kode yang diperlukan secara otomatis, mengikuti pattern yang sudah established. Pilih workflow Anda:

- **🎨 Menu Interaktif**: Tidak perlu hapal command (v0.4)
- **📝 Template Mode**: Scaffolding instan tanpa AI
- **🤖 AI Mode**: Smart generation dari bahasa natural

### Fitur Utama

#### 🎨 Menu Interaktif (Baru di v0.4!)

Cukup jalankan `anaphase` - tidak perlu hapal command:

```bash
anaphase
# Menu interaktif dengan search (Ctrl+K), keyboard navigation, dan filtering
```

Pilih yang Anda butuhkan dari interface TUI yang cantik. Cocok untuk:
- **Pemula**: Temukan command yang tersedia
- **Pro**: Akses cepat dengan filter pencarian `/`

#### 🤖 Dual Mode Generation

**Template Mode** (tidak perlu setup):
```bash
anaphase gen domain
# Prompt untuk: Nama Entity, Fields (name:type)
# Generate: Entity, Repository, Service interfaces
```

**AI Mode** (opsional - butuh API key):
```bash
anaphase gen domain "Order dengan items, bisa dibatalkan jika pending"
# AI mengerti: Entities vs Value Objects, Aggregates, Business rules
# Generate: Validasi advanced, domain events, business logic
```

AI Mode tambahan menyediakan:
- **Smart type inference** dari bahasa natural
- **Business rules** dan validation logic
- **Value Objects** dengan immutability
- **Domain events** dan relationships

#### 🏗️ Clean Architecture

Semua kode yang di-generate mengikuti prinsip Clean Architecture:

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
- Independen dari framework
- Testable di setiap layer
- Infrastruktur yang fleksibel
- Business logic terlindungi

#### ⚡ Rapid Development

Dari zero ke running API dalam hitungan menit - **tidak perlu AI**:

| Task | Traditional | Dengan Anaphase (Template) | Dengan Anaphase (AI) |
|------|-------------|--------------------------|---------------------|
| Project setup | 30 menit | 10 detik | 10 detik |
| Domain model | 2 jam | 30 detik | 30 detik (lebih smart) |
| Repository | 1 jam | 10 detik | 10 detik |
| Handler | 1 jam | 10 detik | 10 detik |
| Auto-wiring | 30 menit | 5 detik | 5 detik |
| **Total** | **~5 jam** | **~1 menit** | **~1 menit** |

::: tip
Template dan AI mode sama-sama menghasilkan kode production-ready. Template Mode cocok untuk CRUD standar, AI Mode menambahkan business logic yang intelligent.
:::

#### 🎯 Type-Safe

Strong typing di seluruh aplikasi:

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

### 1. Interactive Selection

Pilih command Anda dari menu TUI atau gunakan CLI langsung. Filter dengan pencarian `/`, navigasi dengan arrow keys.

### 2. AST Analysis

Anaphase menggunakan AST parser Go untuk menganalisis codebase Anda dan menemukan domain yang sudah ada secara otomatis.

### 3. Code Generation (Dual Mode)

**Template Mode** (default - tanpa AI):
- Menggunakan intelligent templates untuk pattern standar
- Prompt untuk nama entity dan tipe field
- Generate kode DDD-compliant secara instan

**AI Mode** (opsional):
- Memanfaatkan LLMs (Gemini, OpenAI, Claude, Groq)
- Memahami deskripsi bahasa natural
- Generate validasi advanced dan business logic

### 4. Code Output

Kedua mode menghasilkan kode Go production-ready:
- Imports dan packages yang proper
- Error handling
- Validasi (Template: basic, AI: advanced)
- Tests
- Dokumentasi

### 5. Auto-Setup & Wiring

Konfigurasi otomatis untuk semuanya:
- File `.env` dengan database URLs
- Dependencies `go.mod` (auto-installed)
- Dependency injection
- HTTP routes dan middleware

## Pattern Arsitektur

Anaphase enforce best practices:

### Domain-Driven Design (DDD)

- **Entities**: Objek dengan identitas
- **Value Objects**: Objek immutable tanpa identitas
- **Aggregates**: Kumpulan entities dan value objects
- **Repositories**: Abstraksi persistence
- **Services**: Koordinasi business logic

### Hexagonal Architecture

- **Ports**: Interface yang mendefinisikan kontrak
- **Adapters**: Implementasi dari ports
- **Core**: Business logic terisolasi dari infrastruktur

### SOLID Principles

- Single Responsibility
- Open/Closed
- Liskov Substitution
- Interface Segregation
- Dependency Inversion

## Use Cases

### Microservices

Generate multiple bounded contexts dengan cepat menggunakan menu interaktif:

```bash
anaphase  # Buka menu interaktif
# Pilih: Generate Domain (ulangi untuk setiap domain)
# Pilih: Auto-Wire Dependencies
```

Atau via CLI:
```bash
anaphase gen domain "user"
anaphase gen domain "product"
anaphase gen domain "order"
# Auto-wiring terjadi otomatis atau jalankan: anaphase wire
```

### API Backends

REST APIs lengkap dengan operasi CRUD - tidak perlu AI:

```bash
anaphase  # Menu interaktif
# 1. Pilih: Initialize Project → Masukkan nama dan database
# 2. Pilih: Generate Domain → Masukkan entity dan fields
# 3. Pilih: Generate Handler → Masukkan nama domain
# 4. Pilih: Generate Repository → Masukkan nama domain
# Selesai! Semua dependencies auto-installed.
```

### Quick Prototyping

Template Mode cocok untuk rapid prototyping:

```bash
anaphase init my-prototype --db sqlite
cd my-prototype
anaphase  # Generate domains secara interaktif
make run  # Langsung jalan!
```

## Perbandingan

| Fitur | Anaphase | Manual | Generator Lain |
|---------|----------|--------|------------------|
| **Menu Interaktif** | ✅ (v0.4) | ❌ | ❌ |
| **Bekerja Tanpa AI** | ✅ Template Mode | N/A | ✅ |
| **AI-Powered (Opsional)** | ✅ | ❌ | ❌ |
| **DDD Support** | ✅ Kedua mode | Tergantung | ⚠️ |
| **Clean Architecture** | ✅ Enforced | Tergantung | ⚠️ |
| **Auto-Setup** | ✅ .env, deps | ❌ | ⚠️ |
| **Auto-Wiring** | ✅ | ❌ | ⚠️ |
| **Type-Safe** | ✅ | Tergantung | ⚠️ |
| **Multi-DB** | ✅ | ❌ | ✅ |
| **Production-Ready** | ✅ | Tergantung | ⚠️ |
| **Learning Curve** | Sangat Rendah | Tinggi | Sedang |
| **Setup Time** | 0 detik | Berjam-jam | Beberapa menit |

## Filosofi

Anaphase v0.4 dibangun dengan prinsip-prinsip ini:

1. **Workflow Fleksibel**: Pilih yang cocok - Menu Interaktif, Template Mode, atau AI Mode
2. **Tidak Ada Hambatan**: Bekerja langsung tanpa API key atau konfigurasi apapun
3. **AI-Assisted, Bukan Required**: AI meningkatkan generation tapi tidak wajib
4. **Patterns Over Configuration**: Enforce DDD best practices secara default
5. **Transparansi**: Lihat persis apa yang di-generate, tanpa magic
6. **Production-First**: Generate kode yang benar-benar bisa Anda deploy
7. **Developer Experience**: TUI cantik, search, auto-setup - buat pengalaman yang menyenangkan

## Langkah Selanjutnya

- [Quick Start](/id/guide/quick-start) - Bangun service pertama Anda
- [Installation](/id/guide/installation) - Setup detail
- [Architecture](/guide/architecture) - Deep dive ke dalam patterns
