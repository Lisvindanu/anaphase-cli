# Referensi Command

Referensi lengkap untuk semua command CLI Anaphase.

::: info
**Baru di v0.4**: Luncurkan menu interaktif dengan menjalankan `anaphase` (tanpa argumen) untuk interface visual dan searchable ke semua command. Tekan `Ctrl+K` untuk mencari.
:::

## Menu Interaktif

### Akses Cepat

Luncurkan menu TUI untuk akses command secara visual:

```bash
anaphase
```

**Fitur:**
- Browser command visual dengan ikon
- Pencarian cepat dengan `Ctrl+K`
- Bantuan kontekstual untuk setiap opsi
- Prompt auto-setup untuk tools yang hilang
- Tidak perlu mengingat flag CLI

**Opsi Menu:**
```
⚡ Anaphase CLI v0.4.0

📋 Generation
  → Generate Domain       Generate domain models (entities, ports)
  → Generate Handler      Generate HTTP handlers with DTOs
  → Generate Repository   Generate database repositories
  → Generate Middleware   Generate HTTP middleware (auth, CORS, etc.)
  → Generate Migration    Generate database migrations

🔧 Tools
  → Wire Dependencies     Auto-wire dependencies and generate main.go
  → Quality Tools         Code quality (lint, format, validate)
  → Configuration         Manage AI providers and settings

📚 Help & Info
  → Documentation         Open documentation website
  → About                 Version and project information

  Press Ctrl+K to search, Ctrl+C to exit
```

### Fitur Pencarian

Tekan `Ctrl+K` di menu untuk mencari command dengan cepat:

```
Search: middleware
→ Generate Middleware
→ Quality Tools
```

## Flag Global

Tersedia untuk semua command:

| Flag | Short | Deskripsi | Default |
|------|-------|-------------|---------|
| `--verbose` | `-v` | Aktifkan output verbose | `false` |
| `--debug` | | Aktifkan debug logging | `false` |
| `--config` | `-c` | Path file konfigurasi | `~/.anaphase/config.yaml` |

## Commands

::: tip
Semua command di bawah ini juga dapat diakses melalui menu interaktif. Jalankan `anaphase` untuk meluncurkannya.
:::

### `anaphase init`

Inisialisasi proyek microservice baru dengan struktur Clean Architecture.

[Dokumentasi Lengkap →](/reference/init)

```bash
anaphase init [project-name] [flags]
```

### `anaphase gen domain`

Generate model domain (entities, value objects, ports) dengan atau tanpa AI.

[Dokumentasi Lengkap →](/reference/gen-domain)

```bash
# Menu interaktif (disarankan)
anaphase

# CLI langsung
anaphase gen domain "<description>" [flags]

# Mode CLI interaktif
anaphase gen domain --interactive
```

::: info
**Template Mode**: Bekerja tanpa konfigurasi AI menggunakan template yang sudah terbukti.
:::

### `anaphase gen handler`

Generate HTTP handler dengan DTO dan test.

[Dokumentasi Lengkap →](/reference/gen-handler)

```bash
# Menu interaktif (disarankan)
anaphase

# CLI langsung
anaphase gen handler --domain <domain> [flags]
```

### `anaphase gen repository`

Generate implementasi repository database dengan schema.

[Dokumentasi Lengkap →](/reference/gen-repository)

```bash
# Menu interaktif (disarankan)
anaphase

# CLI langsung
anaphase gen repository --domain <domain> --db <database> [flags]
```

### `anaphase gen middleware`

Generate HTTP middleware (auth, CORS, rate limiting, logging).

[Dokumentasi Lengkap →](/reference/gen-middleware)

```bash
# Menu interaktif (disarankan)
anaphase

# CLI langsung
anaphase gen middleware --type <type> [flags]
```

### `anaphase gen migration`

Generate file migration database dengan intelligent SQL generation.

[Dokumentasi Lengkap →](/reference/gen-migration)

```bash
# Menu interaktif (disarankan)
anaphase

# CLI langsung
anaphase gen migration <name> [flags]
```

### `anaphase wire`

Auto-wire dependencies dan generate main.go.

[Dokumentasi Lengkap →](/reference/wire)

```bash
# Menu interaktif (disarankan)
anaphase

# CLI langsung
anaphase wire [flags]
```

::: info
Auto-wiring terjadi secara otomatis setelah generate handler dan repository.
:::

### `anaphase quality`

Code quality tools (lint, format, validate).

[Dokumentasi Lengkap →](/reference/quality)

```bash
# Menu interaktif (disarankan)
anaphase

# CLI langsung
anaphase quality lint [path]
anaphase quality format [path]
anaphase quality validate
```

### `anaphase config`

Kelola AI provider dan konfigurasi.

[Dokumentasi Lengkap →](/reference/config)

```bash
anaphase config list
anaphase config set-provider <provider>
anaphase config check
anaphase config show-providers
```

## Contoh Cepat

### Menggunakan Menu Interaktif (Disarankan)

```bash
# Luncurkan menu interaktif
anaphase

# Pilih opsi secara visual:
# 1. Pilih "Generate Domain"
# 2. Masukkan deskripsi: "User with email, name, and role"
# 3. Pilih mode Template atau AI
# 4. Ikuti prompt untuk handler dan repository
# 5. Auto-wiring terjadi secara otomatis

# Jalankan aplikasi
go run cmd/api/main.go
```

### Menggunakan CLI Langsung

```bash
# Buat proyek
anaphase init my-api

# Generate domain (template mode - tidak perlu AI)
cd my-api
anaphase gen domain "User with email, name, and role (admin, user, guest)"

# Generate infrastructure
anaphase gen handler --domain user
anaphase gen repository --domain user --db postgres

# Wire dan jalankan (auto-wiring juga terjadi otomatis)
anaphase wire
go run cmd/api/main.go
```

### Workflow Multiple Domain

```bash
# Opsi 1: Gunakan menu interaktif untuk setiap domain
anaphase
# Pilih "Generate Domain" → Masukkan detail → Ulangi

# Opsi 2: Gunakan CLI langsung
anaphase gen domain "Product with SKU, name, price, inventory"
anaphase gen domain "Order with customer, items, total, status"
anaphase gen domain "Customer with email, name, addresses"

# Generate semua infrastructure
for domain in product order customer; do
  anaphase gen handler --domain $domain
  anaphase gen repository --domain $domain --db postgres
done

# Wire semuanya (atau auto-wire setelah setiap generasi)
anaphase wire
```

### Pencarian Cepat di Menu

```bash
# Luncurkan menu
anaphase

# Tekan Ctrl+K, ketik "quality"
# → Langsung ke Quality Tools

# Tekan Ctrl+K, ketik "config"
# → Langsung ke Configuration
```

## Environment Variables

Command menghormati environment variable berikut:

- `GEMINI_API_KEY` - Google Gemini API key
- `GROQ_API_KEY` - Groq API key
- `OPENAI_API_KEY` - OpenAI API key
- `ANTHROPIC_API_KEY` - Claude API key
- `ANAPHASE_CONFIG` - Path file konfigurasi
- `DATABASE_URL` - Koneksi database default
- `LOG_LEVEL` - Level logging

::: info
**AI Bersifat Opsional**: Sebagian besar command bekerja dalam template mode tanpa API key yang dikonfigurasi.
:::

## Exit Code

| Code | Arti |
|------|---------|
| 0 | Sukses |
| 1 | Error umum |
| 2 | Error konfigurasi |
| 3 | Error AI provider |
| 4 | Error file system |

## Lihat Juga

- [Quick Start](/guide/quick-start)
- [Configuration](/config/ai-providers)
- [Examples](/examples/basic)
