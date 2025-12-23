# Quick Start

Mulai menggunakan Anaphase dalam waktu kurang dari 5 menit. **Tidak perlu API key AI** - langsung bisa dipakai!

## Prasyarat

- Go 1.21 atau lebih tinggi
- PostgreSQL, MySQL, atau SQLite (opsional, untuk fitur database)
- AI API key (opsional, untuk generation berbasis AI)

## Instalasi

### Dari Source

```bash
git clone https://github.com/lisvindanu/anaphase-cli.git
cd anaphase-cli
go install ./cmd/anaphase
```

### Menggunakan Go Install

```bash
go install github.com/lisvindanu/anaphase-cli/cmd/anaphase@latest
```

Verifikasi instalasi:

```bash
anaphase --version
```

## Konfigurasi AI Provider (Opsional)

::: info Dua Mode Tersedia
Anaphase bekerja dalam **dua mode**:
- **Template Mode**: Langsung bisa dipakai tanpa API key (scaffolding CRUD dasar)
- **AI Mode**: Smart generation dari bahasa natural (butuh API key)

**Anda bisa langsung pakai Anaphase dengan Template Mode!**
:::

Jika Anda ingin generation berbasis AI, setup API key:

```bash
export GEMINI_API_KEY="your-api-key-here"
```

Atau buat config file di `~/.anaphase/config.yaml`:

```yaml
ai:
  primary:
    type: gemini
    apiKey: your-api-key-here
    model: gemini-2.5-flash
```

::: tip Dapatkan API Key Gratis
Dapatkan Gemini API key gratis di [Google AI Studio](https://makersuite.google.com/app/apikey).

Anaphase juga support OpenAI, Claude, dan Groq. [Pelajari lebih lanjut →](/config/ai-providers)
:::

## Buat Project Pertama Anda

### Menu Interaktif (Direkomendasikan)

**Baru di v0.4!** Cukup jalankan `anaphase` untuk akses menu interaktif - tidak perlu hapal command:

```bash
anaphase
```

Menu interaktif muncul:

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

Pilih **"Initialize Project"** dan ikuti prompt:

```bash
Project name: my-app
Database type (postgres/mysql/sqlite) [postgres]: postgres

✅ Project created with auto-generated .env and dependencies!
```

::: tip Pro Tip
Gunakan `/` untuk search/filter command di menu interaktif. Coba ketik "domain" untuk cepat menemukan domain generation!
:::

### Command Line (Alternatif)

Anda juga bisa gunakan command langsung:

```bash
anaphase init my-app --db postgres
cd my-app
```

Kedua metode menghasilkan struktur project lengkap:

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
├── .env          # Auto-generated dengan config database
├── .env.example
├── go.mod
└── README.md
```

### Langkah 2: Generate Domain

#### Menggunakan Menu Interaktif

```bash
cd my-app
anaphase
```

Pilih **"Generate Domain"** dari menu. Anaphase akan otomatis menggunakan:
- **Template Mode** jika AI tidak dikonfigurasi (langsung bisa dipakai!)
- **AI Mode** jika API key sudah di-setup

**Contoh Template Mode:**

```
📝 Template Mode - Domain Generation

Entity name: Customer
Fields: email:string, name:string, phone:string

✅ Generated:
  ✓ internal/core/entity/customer.go
  ✓ internal/core/port/customer_repository.go
  ✓ internal/core/port/customer_service.go
```

**Contoh AI Mode** (dengan API key terkonfigurasi):

```
🧠 AI-Powered Domain Generation

Description: Customer dengan email, name, phone. Bisa melakukan order.

✅ Generated:
  ✓ internal/core/entity/customer.go (dengan validasi)
  ✓ internal/core/valueobject/email.go (validasi email)
  ✓ internal/core/valueobject/phone.go (validasi phone)
  ✓ internal/core/port/customer_repository.go
  ✓ internal/core/port/customer_service.go
```

#### Menggunakan Command Line

**Template Mode:**
```bash
anaphase gen domain "Customer"
# Prompt untuk nama entity dan fields secara interaktif
```

**AI Mode:**
```bash
anaphase gen domain "Customer dengan email, name, dan phone. Bisa melakukan order."
```

Kedua mode membuat domain model DDD-compliant dengan:
- Entity dengan business logic
- Repository interface (port)
- Service interface (port)
- Value objects (AI mode menambahkan validasi smart)

### Langkah 3: Generate Handlers

Menggunakan menu interaktif, pilih **"Generate Handler"**:

```bash
Handler name: customer

✅ Generated:
  ✓ internal/adapter/handler/http/customer_handler.go (CRUD endpoints)
  ✓ internal/adapter/handler/http/customer_dto.go (Request/Response DTOs)
  ✓ internal/adapter/handler/http/customer_handler_test.go
```

Atau via command line:
```bash
anaphase gen handler customer
```

### Langkah 4: Generate Repository

Pilih **"Generate Repository"** dari menu:

```bash
Repository name: customer

✅ Generated:
  ✓ internal/adapter/repository/postgres/customer_repo.go
  ✓ internal/adapter/repository/postgres/schema.sql
  ✓ internal/adapter/repository/postgres/customer_repo_test.go
```

Atau via command line:
```bash
anaphase gen repository customer
```

### Langkah 5: Jalankan Aplikasi Anda

**Auto-setup sudah selesai!** File `.env` Anda sudah dibuat saat `init`. Tinggal jalankan database dan running:

```bash
# Jalankan PostgreSQL dengan Docker
docker run -d \
  --name anaphase-db \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=my-app \
  -p 5432:5432 \
  postgres:16-alpine

# Jalankan API Anda (dependencies sudah terinstall!)
make run
```

::: tip Kredensial Database
File `.env` sudah auto-generated dengan DATABASE_URL yang benar. Tinggal update password jika perlu!
:::

## Test API Anda

API Anda sekarang berjalan di `http://localhost:8080`. Test:

### Buat Customer

```bash
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "name": "John Doe",
    "phone": "+1234567890"
  }'
```

### Ambil Semua Customers

```bash
curl http://localhost:8080/api/v1/customers
```

### Health Check

```bash
curl http://localhost:8080/health
```

## Template Mode vs AI Mode

| Fitur | Template Mode | AI Mode |
|---------|--------------|---------|
| **Setup Required** | ❌ Tidak ada | ✅ API Key |
| **Kecepatan Generation** | ⚡ Instan | 🔄 2-5 detik |
| **Use Case** | Entity CRUD standar | Business logic kompleks |
| **Tipe Field** | Tipe dasar (string, int, dll) | Tipe smart + validasi |
| **Value Objects** | ❌ Tidak termasuk | ✅ Auto-generated |
| **Business Logic** | CRUD dasar | Method spesifik domain |
| **Bahasa Natural** | ❌ Tidak | ✅ Ya |
| **Biaya** | 🆓 Gratis | 🆓 Tier gratis tersedia |

::: tip Kapan Menggunakan Mode Mana
- **Template Mode**: Cocok untuk prototyping cepat, entity standar, dan belajar pola DDD
- **AI Mode**: Terbaik untuk domain kompleks, validasi spesifik bisnis, dan kode production-ready
:::

## Langkah Selanjutnya

- Pelajari tentang [Architecture](/guide/architecture)
- Eksplor [AI-Powered Generation](/guide/ai-generation) (opsional)
- Baca [Command Reference](/reference/commands)
- Lihat [Examples](/examples/basic)
- Coba **fitur Search** (tekan `Ctrl+K` atau `Cmd+K`)

::: tip Pro Tip
Menu interaktif punya fitur search! Tekan `/` untuk filter command dan temukan yang Anda butuhkan dengan cepat.
:::
