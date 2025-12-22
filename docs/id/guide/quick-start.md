# Mulai Cepat

Mulai menggunakan Anaphase dalam waktu kurang dari 5 menit.

## Prasyarat

- Go 1.21 atau lebih tinggi
- PostgreSQL (opsional, untuk fitur database)
- Google Gemini API key (tersedia tier gratis)

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

## Konfigurasi AI Provider

Setup Google Gemini API key Anda:

```bash
export GEMINI_API_KEY="your-api-key-here"
```

Atau buat file config di `~/.anaphase/config.yaml`:

```yaml
ai:
  primary:
    type: gemini
    apiKey: your-api-key-here
    model: gemini-2.5-flash
```

::: tip Dapatkan API Key Gratis
Dapatkan Gemini API key gratis di [Google AI Studio](https://makersuite.google.com/app/apikey)
:::

## Buat Project Pertama Anda

### Langkah 1: Inisialisasi

Buat project microservice baru:

```bash
anaphase init my-app
cd my-app
```

Ini akan menghasilkan struktur project lengkap:

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
├── go.mod
└── README.md
```

### Langkah 2: Generate Domain

Gunakan AI untuk generate domain model lengkap:

```bash
anaphase gen domain \
  --name customer \
  --prompt "Customer dengan alamat email, nama lengkap, dan nomor telepon. Customer bisa melakukan order."
```

Ini akan membuat:
- `internal/core/entity/customer.go` - Entity dengan business logic
- `internal/core/valueobject/email.go` - Value objects
- `internal/core/port/customer_repo.go` - Repository interface
- `internal/core/port/customer_service.go` - Service interface

### Langkah 3: Generate Handlers

Buat HTTP handlers untuk domain Anda:

```bash
anaphase gen handler --domain customer
```

File yang dihasilkan:
- `internal/adapter/handler/http/customer_handler.go` - CRUD endpoints
- `internal/adapter/handler/http/customer_dto.go` - Request/Response DTOs
- `internal/adapter/handler/http/customer_handler_test.go` - Tests

### Langkah 4: Generate Repository

Buat implementasi database:

```bash
anaphase gen repository --domain customer --db postgres
```

File yang dihasilkan:
- `internal/adapter/repository/postgres/customer_repo.go` - Repository implementation
- `internal/adapter/repository/postgres/schema.sql` - Database schema
- `internal/adapter/repository/postgres/customer_repo_test.go` - Tests

### Langkah 5: Wire Semua Dependencies

Otomatis wire semua dependencies:

```bash
anaphase wire
```

Ini akan menghasilkan:
- `cmd/api/main.go` - HTTP server dengan graceful shutdown
- `cmd/api/wire.go` - Dependency injection

### Langkah 6: Jalankan

Jalankan database:

```bash
docker run -d \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=anaphase \
  -p 5432:5432 \
  postgres:16-alpine
```

Apply migrations:

```bash
psql -h localhost -U postgres -d anaphase -f internal/adapter/repository/postgres/schema.sql
```

Jalankan API Anda:

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/anaphase?sslmode=disable"
go run cmd/api/main.go
```

## Test API Anda

API Anda sekarang berjalan di `http://localhost:8080`. Coba test:

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

## Selanjutnya Apa?

- Pelajari tentang [Architecture](/guide/architecture)
- Eksplor [AI-Powered Generation](/guide/ai-generation)
- Baca [Command Reference](/reference/commands)
- Lihat [Examples](/examples/basic)

::: tip Pro Tip
Gunakan flag `--verbose` dengan command apapun untuk melihat output detail:
```bash
anaphase gen domain --name product --prompt "..." --verbose
```
:::
