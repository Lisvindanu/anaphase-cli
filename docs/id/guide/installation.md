# Instalasi

## Prasyarat

Sebelum install Anaphase, pastikan Anda punya:

- **Go 1.21+**: [Download Go](https://go.dev/dl/)
- **Git**: Untuk cloning repositories
- **PostgreSQL** (Opsional): Untuk fitur database
  - Atau Docker untuk menjalankan Postgres dalam container

Verifikasi instalasi Go:

```bash
go version
# Harus output: go version go1.21.x atau lebih tinggi
```

## Install Anaphase

### Opsi 1: Quick Install (Direkomendasikan)

Gunakan install script kami yang otomatis konfigurasi PATH:

::: code-group

```bash [macOS/Linux]
curl -fsSL https://raw.githubusercontent.com/lisvindanu/anaphase-cli/main/install.sh | bash
```

```powershell [Windows]
irm https://raw.githubusercontent.com/lisvindanu/anaphase-cli/main/install.ps1 | iex
```

:::

Atau download dan jalankan manual:

::: code-group

```bash [macOS/Linux]
curl -fsSL https://raw.githubusercontent.com/lisvindanu/anaphase-cli/main/install.sh -o install.sh
chmod +x install.sh
./install.sh
```

```powershell [Windows]
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/lisvindanu/anaphase-cli/main/install.ps1" -OutFile "install.ps1"
powershell -ExecutionPolicy Bypass -File install.ps1
```

:::

Script ini akan:
- Install versi terbaru via `go install`
- Deteksi shell Anda (bash, zsh, fish) atau PowerShell di Windows
- Menawarkan untuk menambahkan Go binary directory ke PATH secara otomatis

### Opsi 2: Manual Install

Install langsung menggunakan `go install`:

```bash
go install github.com/lisvindanu/anaphase-cli/cmd/anaphase@latest
```

Kemudian tambahkan ke PATH:

::: code-group

```bash [Bash]
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

```bash [Zsh]
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

```bash [Fish]
set -gx PATH $HOME/go/bin $PATH
```

```powershell [Windows PowerShell]
# Temporary (session sekarang)
$env:Path += ";$(go env GOPATH)\bin"

# Permanent (semua session)
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$(go env GOPATH)\bin", "User")
```

:::

### Opsi 3: Dari Source

Clone dan build dari source untuk development:

```bash
git clone https://github.com/lisvindanu/anaphase-cli.git
cd anaphase-cli
go mod download
go install ./cmd/anaphase
```

## Konfigurasi AI Provider

Anaphase membutuhkan AI provider untuk domain generation. Saat ini support:

- **Google Gemini** (Direkomendasikan, tersedia tier gratis)

### Dapatkan Gemini API Key

1. Kunjungi [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Sign in dengan akun Google Anda
3. Klik "Create API Key"
4. Copy API key Anda

::: tip Free Tier
Google Gemini menawarkan free tier yang cukup generous:
- 60 requests per menit
- Sempurna untuk development dan project kecil
:::

### Konfigurasi API Key

#### Metode 1: Environment Variable

```bash
export GEMINI_API_KEY="your-api-key-here"
```

Tambahkan ke shell profile agar persisten:

```bash
# ~/.bashrc atau ~/.zshrc
export GEMINI_API_KEY="your-api-key-here"
```

#### Metode 2: Configuration File

Buat `~/.anaphase/config.yaml`:

```yaml
ai:
  primary:
    type: gemini
    apiKey: YOUR_API_KEY_HERE
    model: gemini-2.5-flash
    timeout: 30s
    retries: 3

  # Opsional: fallback providers
  secondary:
    type: gemini
    apiKey: BACKUP_API_KEY
    model: gemini-2.5-flash

cache:
  enabled: true
  ttl: 24h
  dir: ~/.anaphase/cache
```

::: details Opsi Konfigurasi

- `type`: Tipe AI provider (`gemini`)
- `apiKey`: API key Anda
- `model`: Model yang digunakan (`gemini-2.5-flash` direkomendasikan)
- `timeout`: Request timeout (default: `30s`)
- `retries`: Jumlah retry saat gagal (default: `3`)
- `cache.enabled`: Aktifkan response caching (default: `true`)
- `cache.ttl`: Cache time-to-live (default: `24h`)
:::

## Setup Database (Opsional)

Untuk repository generation, Anda perlu database.

### PostgreSQL dengan Docker

Cara paling mudah untuk memulai:

```bash
docker run -d \
  --name anaphase-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=anaphase \
  -p 5432:5432 \
  postgres:16-alpine
```

### PostgreSQL Native

Install PostgreSQL untuk sistem Anda:

::: code-group

```bash [macOS]
brew install postgresql@16
brew services start postgresql@16
createdb anaphase
```

```bash [Ubuntu/Debian]
sudo apt-get install postgresql-16
sudo systemctl start postgresql
sudo -u postgres createdb anaphase
```

```bash [Windows]
# Download dari https://www.postgresql.org/download/windows/
# Atau gunakan WSL dengan instruksi Ubuntu
```

:::

### MySQL (Alternatif)

```bash
docker run -d \
  --name anaphase-mysql \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=anaphase \
  -p 3306:3306 \
  mysql:8
```

### MongoDB (Alternatif)

```bash
docker run -d \
  --name anaphase-mongo \
  -p 27017:27017 \
  mongo:7
```

## Verifikasi Instalasi

Test bahwa semuanya berjalan:

```bash
# Check version
anaphase --version

# Check help
anaphase --help

# Initialize project test
mkdir test-project
cd test-project
anaphase init
```

Anda akan melihat:

```
✅ Project initialized successfully!

Next steps:
  1. Configure your AI provider (see docs)
  2. Generate your first domain:
     anaphase gen domain --name user --prompt "User dengan email dan nama"
  3. Run the API:
     go run cmd/api/main.go
```

## Environment Variables

Anaphase menggunakan environment variables berikut:

| Variable | Deskripsi | Default |
|----------|-------------|---------|
| `GEMINI_API_KEY` | Google Gemini API key | Required |
| `DATABASE_URL` | Database connection string | `postgres://...` |
| `PORT` | HTTP server port | `8080` |
| `LOG_LEVEL` | Logging level (`debug`, `info`, `warn`, `error`) | `info` |
| `ANAPHASE_CONFIG` | Config file path | `~/.anaphase/config.yaml` |

## Troubleshooting

### Command not found

Jika `anaphase` tidak ditemukan:

```bash
# Check jika sudah terinstall
ls -la $(go env GOPATH)/bin/anaphase

# Tambahkan ke PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

### Masalah API Key

Jika Anda melihat "API key not configured":

```bash
# Verifikasi environment variable
echo $GEMINI_API_KEY

# Atau cek config file
cat ~/.anaphase/config.yaml
```

### Database Connection Failed

Jika koneksi database gagal:

```bash
# Test connection
psql -h localhost -U postgres -d anaphase

# Check jika berjalan
docker ps | grep postgres

# Lihat logs
docker logs anaphase-postgres
```

### Import Errors

Jika Anda melihat import errors setelah generation:

```bash
# Download dependencies
go mod download

# Tidy modules
go mod tidy
```

## Update Anaphase

Jaga Anaphase tetap up-to-date untuk mendapatkan fitur terbaru dan bug fixes.

### Cek Versi Saat Ini

```bash
anaphase --version
```

### Update ke Versi Terbaru

::: code-group

```bash [Quick Update]
# Direkomendasikan: Gunakan install script
curl -fsSL https://raw.githubusercontent.com/lisvindanu/anaphase-cli/main/install.sh | bash
```

```bash [Manual Update]
# Menggunakan go install
go install github.com/lisvindanu/anaphase-cli/cmd/anaphase@latest
```

```bash [Dari Source]
# Pull perubahan terbaru
cd anaphase-cli
git pull origin main

# Rebuild
go install ./cmd/anaphase
```

```bash [Docker]
# Pull code terbaru
cd /var/www/anaphase-cli
git pull origin main

# Rebuild image
docker compose build

# Atau pull dari registry
docker pull ghcr.io/lisvindanu/anaphase-cli:latest
```

:::

### Apa yang Baru

Cek [changelog](https://github.com/lisvindanu/anaphase-cli/releases) untuk fitur baru:

**Update Terbaru:**
- ✨ Provider Selection CLI - Pilih AI provider dengan flag `--provider`
- ✨ Config Management - Kelola providers dengan `anaphase config`
- ✨ Middleware Generator - Generate auth, rate limit, logging, CORS
- ✨ Interactive Mode - Guided prompts dengan flag `-i`
- ✨ Code Quality Tools - Lint, format, dan validasi code
- ✨ Migration Generator - File database migration dengan smart SQL

### Verifikasi Update

```bash
# Check versi baru
anaphase --version

# Test fitur baru
anaphase config show-providers
anaphase gen middleware --help
anaphase quality --help
```

### Update Konfigurasi

Setelah update, konfigurasi Anda mungkin perlu diperbarui:

```bash
# Check config saat ini
anaphase config list

# Update provider jika diperlukan
anaphase config set-provider groq

# Health check semua providers
anaphase config check
```

### Rollback (Jika Diperlukan)

Jika Anda perlu rollback ke versi spesifik:

```bash
# Install versi spesifik
go install github.com/lisvindanu/anaphase-cli/cmd/anaphase@v1.0.0

# Atau dari source di tag spesifik
git checkout v1.0.0
go install ./cmd/anaphase
```

## Langkah Selanjutnya

- [Mulai Cepat](/guide/quick-start) - Build service pertama Anda
- [Architecture](/guide/architecture) - Pahami pattern yang digunakan
- [AI Generation](/guide/ai-generation) - Pelajari fitur AI
- [Domain-Driven Design](/guide/ddd) - **Key differentiator kami**
