# Troubleshooting Guide

Mengalami masalah? Panduan ini mencakup masalah umum dan cara memperbaikinya.

::: tip Untuk Pemula
Jika ini pertama kalinya Anda menggunakan Anaphase atau Go, baca bagian [Common Issues](#common-issues) terlebih dahulu. Sebagian besar masalah memiliki solusi sederhana!
:::

## Common Issues

### Interactive Menu Issues

**Masalah:**
```bash
$ anaphase
# Menu muncul tetapi commands tidak bekerja seperti yang diharapkan
```

**Apa artinya:**
Sejak v0.4.0, menjalankan `anaphase` tanpa argumen akan meluncurkan menu TUI interaktif. Jika Anda mengalami masalah, Anda selalu dapat kembali ke perintah langsung.

**Perbaikan:**

**Jika menu tidak responsif:**
```bash
# Gunakan perintah langsung
anaphase init my-project
anaphase gen domain --name user

# Atau periksa kompatibilitas terminal
echo $TERM
# Harus menampilkan sesuatu seperti "xterm-256color"
```

**Jika menu tidak ditampilkan dengan benar:**
```bash
# Update terminal Anda atau gunakan perintah langsung
# Menu memerlukan terminal yang mendukung escape sequences

# Bypass menu sepenuhnya
anaphase --help  # Menampilkan semua perintah yang tersedia
```

**Tips navigasi:**
- Gunakan arrow keys atau j/k untuk navigasi
- Tekan Enter untuk memilih
- Tekan ESC atau q untuk kembali/keluar
- Tekan Ctrl+C untuk keluar langsung

---

### Template Mode vs AI Mode Issues

**Masalah:**
```bash
$ anaphase gen domain --name user --prompt "User with email"
Error: GEMINI_API_KEY environment variable not set
```

**Apa artinya:**
Anda mencoba menggunakan AI mode tanpa API key. Sejak v0.4.0, Anda dapat memilih antara Template dan AI mode.

**Perbaikan:**

**Opsi 1: Gunakan Template Mode (tidak perlu API key)**
```bash
# Template mode cepat dan tidak memerlukan API key
anaphase gen domain --name user --template

# Atau gunakan menu interaktif dan pilih Template Mode
anaphase
```

**Opsi 2: Setup AI Mode**
```bash
# Dapatkan API key dari https://makersuite.google.com/app/apikey
export GEMINI_API_KEY="your-key-here"

# Sekarang Anda dapat menggunakan AI mode
anaphase gen domain --name user --prompt "User with email and profile"
```

**Kapan menggunakan setiap mode:**
- **Template Mode**: Scaffolding cepat, domain sederhana, tidak ada API key tersedia
- **AI Mode**: Logika bisnis kompleks, pola DDD tingkat lanjut, domain events

---

### Auto-Setup Issues

**Masalah:**
```bash
$ anaphase init my-project
# Dependencies tidak ter-install otomatis
```

**Apa artinya:**
Sejak v0.4.0, fitur auto-setup mencoba menjalankan `go mod download` secara otomatis, tetapi mungkin gagal jika Go tidak dikonfigurasi dengan benar.

**Perbaikan:**

**Jika auto-setup gagal:**
```bash
# Selesaikan setup secara manual
cd my-project
go mod download
go mod tidy

# Verifikasi Go dikonfigurasi
go version
go env GOPATH
```

**Nonaktifkan auto-setup jika menyebabkan masalah:**
```bash
# Generate tanpa auto-setup
anaphase init my-project --no-auto-setup

# Kemudian jalankan manual
cd my-project
go mod download
```

---

### "missing go.sum entry" Error

**Masalah:**
```bash
$ make run
internal/config/config.go:6:2: missing go.sum entry for module providing package github.com/spf13/viper
```

**Apa artinya:**
Dependency Go belum diunduh. Ini terjadi setelah `anaphase init`.

**Perbaikan:**
```bash
# Download semua dependencies
go mod download

# Atau jalankan ini (melakukan hal yang sama)
go mod tidy

# Kemudian coba lagi
make run
```

**Mengapa terjadi:**
`anaphase init` membuat struktur proyek tetapi tidak mengunduh dependencies secara otomatis. Anda perlu menjalankan `go mod download` terlebih dahulu.

---

### "command not found: anaphase"

**Masalah:**
```bash
$ anaphase --version
zsh: command not found: anaphase
```

**Apa artinya:**
Binary `anaphase` tidak ada di PATH Anda.

**Perbaikan:**

**Opsi 1: Gunakan install script (rekomendasi)**
```bash
curl -fsSL https://raw.githubusercontent.com/lisvindanu/anaphase-cli/main/install.sh | bash
```

**Opsi 2: Tambahkan ke PATH secara manual**
```bash
# Periksa di mana Go menginstall binaries
go env GOPATH

# Tambahkan ke shell config Anda
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

# Verifikasi
which anaphase
```

---

### "GEMINI_API_KEY not configured"

**Masalah:**
```bash
$ anaphase gen domain --name user --prompt "User with email"
Error: GEMINI_API_KEY environment variable not set
```

**Apa artinya:**
Flag `--prompt` memerlukan AI mode, yang membutuhkan Gemini API key. Namun, AI bersifat opsional di v0.4.0.

**Perbaikan:**

**Opsi 1: Gunakan Template Mode (rekomendasi untuk sebagian besar kasus)**
```bash
# Tidak perlu API key - menggunakan templates
anaphase gen domain --name user --template

# Atau gunakan menu interaktif
anaphase
# Pilih "Template Mode" saat diminta
```

**Opsi 2: Setup AI Mode (untuk fitur lanjutan)**

1. **Dapatkan API key:**
   - Kunjungi https://makersuite.google.com/app/apikey
   - Sign in dengan Google
   - Klik "Create API Key"
   - Copy key tersebut

2. **Set key:**
```bash
# Temporary (sesi ini saja)
export GEMINI_API_KEY="your-key-here"

# Permanent (tambahkan ke shell config)
echo 'export GEMINI_API_KEY="your-key-here"' >> ~/.zshrc
source ~/.zshrc
```

3. **Verifikasi:**
```bash
echo $GEMINI_API_KEY
# Harus mencetak key Anda

# Sekarang gunakan AI mode dengan prompts
anaphase gen domain --name user --prompt "User with email and profile"
```

::: tip AI Bersifat Opsional
Anda tidak perlu API key untuk menggunakan Anaphase! Template mode menyediakan semua fungsionalitas inti tanpa ketergantungan eksternal. Gunakan AI mode hanya ketika Anda membutuhkan domain modeling tingkat lanjut.
:::

---

### Database Connection Failed

**Masalah:**
```bash
$ make run
Error: failed to connect to database: connection refused
```

**Apa artinya:**
PostgreSQL tidak berjalan atau connection string salah.

**Perbaikan:**

**Quick fix dengan Docker:**
```bash
# Start PostgreSQL di Docker
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 \
  postgres:16-alpine

# Tunggu beberapa detik untuk startup
sleep 3

# Set connection string
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable"

# Coba lagi
make run
```

**Periksa apakah PostgreSQL berjalan:**
```bash
# Dengan Docker
docker ps | grep postgres

# Dengan psql
psql -h localhost -U postgres -d mydb -c "SELECT 1"
```

**Format connection string umum:**
```bash
# Local PostgreSQL
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable"

# Docker PostgreSQL
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable"

# Remote PostgreSQL
export DATABASE_URL="postgres://user:password@host:5432/database?sslmode=require"
```

---

### "make: command not found"

**Masalah:**
```bash
$ make run
zsh: command not found: make
```

**Apa artinya:**
`make` tidak terinstall di sistem Anda.

**Perbaikan:**

**macOS:**
```bash
xcode-select --install
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt-get update
sudo apt-get install build-essential
```

**Atau jalankan langsung tanpa make:**
```bash
# Alih-alih: make run
go run cmd/api/main.go

# Alih-alih: make build
go build -o bin/api cmd/api/main.go
```

---

### Import Errors After Generation

**Masalah:**
```bash
$ make run
internal/core/entity/user.go:5:2: no required module provides package github.com/google/uuid
```

**Apa artinya:**
Dependency baru ditambahkan oleh code generation tetapi belum diunduh.

**Perbaikan:**
```bash
# Download missing dependencies
go mod tidy

# Atau
go get ./...

# Kemudian jalankan lagi
make run
```

**Selalu jalankan setelah generation:**
```bash
anaphase gen domain --name user --prompt "..."
go mod tidy  # ← Jalankan ini!
```

---

### Port Already in Use

**Masalah:**
```bash
$ make run
Error: listen tcp :8080: bind: address already in use
```

**Apa artinya:**
Proses lain menggunakan port 8080.

**Perbaikan:**

**Cari apa yang menggunakan port:**
```bash
# macOS/Linux
lsof -i :8080

# Kill proses tersebut
kill -9 <PID>
```

**Atau gunakan port yang berbeda:**
```bash
export PORT=3000
make run
```

---

### Permission Denied

**Masalah:**
```bash
$ make run
zsh: permission denied: ./bin/api
```

**Apa artinya:**
Binary tidak memiliki execute permissions.

**Perbaikan:**
```bash
chmod +x bin/api
./bin/api
```

---

### "go: cannot find main module"

**Masalah:**
```bash
$ go run cmd/api/main.go
go: cannot find main module; see 'go help modules'
```

**Apa artinya:**
Anda tidak berada di direktori Go module (tidak ada file `go.mod`).

**Perbaikan:**
```bash
# Pastikan Anda berada di direktori proyek
cd my-api

# Verifikasi go.mod ada
ls go.mod

# Kemudian jalankan
go run cmd/api/main.go
```

---

## Step-by-Step: First Time Setup

Jika Anda benar-benar baru, ikuti langkah-langkah ini secara berurutan:

### 1. Install Go

**Periksa apakah sudah terinstall:**
```bash
go version
```

**Jika belum terinstall:**
- Download dari https://go.dev/dl/
- Install untuk OS Anda
- Verifikasi: `go version`

### 2. Install Anaphase

```bash
# Quick install
curl -fsSL https://raw.githubusercontent.com/lisvindanu/anaphase-cli/main/install.sh | bash

# Verifikasi
anaphase --version
```

### 3. Pilih Mode Anda

::: tip Baru di v0.4.0
Anda sekarang dapat memilih antara Template Mode (tidak perlu setup) dan AI Mode (memerlukan API key).

**Mulai dengan Template Mode** - lebih sederhana dan tidak memerlukan API key!
:::

**Opsi A: Template Mode (Rekomendasi untuk Pemula)**
```bash
# Tidak perlu setup tambahan!
# Skip ke Langkah 4
```

**Opsi B: AI Mode (Opsional)**

1. Kunjungi https://makersuite.google.com/app/apikey
2. Buat API key
3. Set:
```bash
echo 'export GEMINI_API_KEY="your-key-here"' >> ~/.zshrc
source ~/.zshrc
```

### 4. Start PostgreSQL

**Dengan Docker (termudah):**
```bash
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 \
  postgres:16-alpine
```

**Set connection:**
```bash
echo 'export DATABASE_URL="postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable"' >> ~/.zshrc
source ~/.zshrc
```

### 5. Buat Proyek Pertama Anda

::: tip Gunakan Menu Interaktif
Sejak v0.4.0, Anda dapat menggunakan menu interaktif untuk pembuatan proyek yang lebih mudah:
```bash
anaphase  # Meluncurkan menu interaktif
```
Atau lanjutkan dengan perintah seperti yang ditunjukkan di bawah.
:::

**Menggunakan Template Mode (Tidak Perlu API Key):**
```bash
# Buat proyek
anaphase init my-api
cd my-api

# Auto-setup akan menjalankan go mod download untuk Anda
# Jika gagal, jalankan manual: go mod download

# Generate domain menggunakan templates
anaphase gen domain --name user --template

# Download new dependencies
go mod tidy

# Generate handler
anaphase gen handler --domain user

# Generate repository
anaphase gen repository --domain user --db postgres

# Wire everything
anaphase wire

# Download final dependencies
go mod tidy

# Jalankan!
make run
```

**Menggunakan AI Mode (Memerlukan API Key):**
```bash
# Buat proyek
anaphase init my-api
cd my-api

# Generate domain dengan AI
anaphase gen domain --name user --prompt "User with email, name, and optional profile picture URL"

# Sisanya sama
go mod tidy
anaphase gen handler --domain user
anaphase gen repository --domain user --db postgres
anaphase wire
go mod tidy
make run
```

### 6. Test API Anda

```bash
# Di terminal lain, buat user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "name": "Test User"
  }'
```

**Anda harus melihat:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "test@example.com",
  "name": "Test User",
  "createdAt": "2024-01-01T00:00:00Z"
}
```

🎉 **Berhasil!** API Anda bekerja!

---

## Kesalahan Umum

### ❌ Lupa `go mod tidy`

```bash
# Setelah SETIAP perintah generation, jalankan:
anaphase gen domain --name product --prompt "..."
go mod tidy  # ← Jangan lupa!
```

### ❌ Direktori yang salah

```bash
# Salah - menjalankan dari home directory
~ $ make run
make: *** No rule to make target 'run'.

# Benar - di dalam direktori proyek
~/my-api $ make run
Starting my-api...
```

### ❌ Tidak set environment variables

```bash
# Tidak akan bekerja
make run

# Akan bekerja
export GEMINI_API_KEY="..."
export DATABASE_URL="..."
make run
```

### ❌ Menggunakan versi Go yang lama

```bash
# Periksa versi
go version
# Harus 1.21 atau lebih tinggi

# Jika terlalu lama, update Go
# Download dari https://go.dev/dl/
```

---

## Mendapatkan Bantuan

### Periksa Logs

```bash
# Jalankan dengan verbose output
go run cmd/api/main.go 2>&1 | tee app.log

# Periksa apa yang salah
cat app.log
```

### Verifikasi Installation

```bash
# Periksa Go
go version
# Harus 1.21+

# Periksa Anaphase
anaphase --version

# Periksa environment
echo $GEMINI_API_KEY
echo $DATABASE_URL

# Periksa database
psql $DATABASE_URL -c "SELECT 1"
```

### Clean Start

Jika semuanya rusak, mulai fresh:

```bash
# Hapus proyek
rm -rf my-api

# Bersihkan module cache
go clean -modcache

# Mulai lagi
anaphase init my-api
cd my-api
go mod download
# ... lanjutkan ...
```

---

## Masih Stuck?

### Baca Dokumentasi
- [Quick Start](/guide/quick-start)
- [Installation](/guide/installation)
- [Architecture](/guide/architecture)

### Periksa Contoh
- [Basic Example](/examples/basic)
- [Multi-Domain](/examples/multi-domain)

### Pertanyaan Umum

**Q: Apakah saya perlu tahu Go?**
A: Pengetahuan dasar Go membantu, tetapi Anaphase menghasilkan sebagian besar kode untuk Anda.

**Q: Apakah saya perlu Gemini API key?**
A: Tidak! Sejak v0.4.0, Template Mode tidak memerlukan API key. AI mode sepenuhnya opsional.

**Q: Apa perbedaan antara Template dan AI mode?**
A: Template mode menggunakan scaffolding yang sudah disiapkan (cepat, deterministik). AI mode menghasilkan model domain yang lebih canggih berdasarkan natural language prompts.

**Q: Apakah Gemini API gratis?**
A: Ya! Free tier mencakup 60 requests/menit. Tapi ingat, Anda dapat menggunakan Template Mode tanpa API key.

**Q: Bisakah saya menggunakan menu interaktif di terminal apapun?**
A: Menu bekerja di sebagian besar terminal modern. Jika Anda mengalami masalah, gunakan perintah langsung.

**Q: Bisakah saya menggunakan MySQL alih-alih PostgreSQL?**
A: Ya! Gunakan `--db mysql` saat menghasilkan repositories.

**Q: Bagaimana cara menambahkan logic kustom?**
A: Edit file service layer yang dihasilkan. Lihat [Custom Handlers](/examples/custom-handlers).

**Q: Bisakah saya menjalankan tanpa Docker?**
A: Ya, install PostgreSQL secara native untuk OS Anda.

**Q: Bisakah saya beralih antara Template dan AI mode?**
A: Ya! Gunakan Template mode untuk scaffolding dasar, lalu tingkatkan secara manual atau gunakan AI mode untuk penambahan kompleks.

---

## Pro Tips

### 1. Gunakan file `.env`

Buat `.env` di proyek Anda:
```bash
GEMINI_API_KEY=your-key-here
DATABASE_URL=postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable
PORT=8080
```

Load:
```bash
# Install dotenv tool
go install github.com/joho/godotenv/cmd/godotenv@latest

# Jalankan dengan .env
godotenv go run cmd/api/main.go
```

### 2. Buat setup script

Buat `setup.sh`:
```bash
#!/bin/bash
set -e

echo "Setting up project..."

# Download dependencies
go mod download

# Start database
docker run -d --name postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 \
  postgres:16-alpine

sleep 3

# Apply migrations
psql "postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable" \
  -f internal/adapter/repository/postgres/schema.sql

echo "✅ Setup complete!"
```

Jalankan:
```bash
chmod +x setup.sh
./setup.sh
```

### 3. Tambahkan ke Makefile

Edit `Makefile`:
```makefile
.PHONY: setup deps clean

setup: deps
    @echo "Starting PostgreSQL..."
    docker run -d --name postgres \
      -e POSTGRES_PASSWORD=postgres \
      -p 5432:5432 postgres:16-alpine

deps:
    @echo "Downloading dependencies..."
    go mod download
    go mod tidy

clean:
    @echo "Cleaning up..."
    docker stop postgres || true
    docker rm postgres || true
    go clean
```

Sekarang Anda dapat:
```bash
make setup   # First time setup
make deps    # Download dependencies
make clean   # Clean everything
```

---

::: tip Ingat (v0.4.0+)
1. Coba menu interaktif: `anaphase` (tanpa argumen)
2. Gunakan Template Mode dulu - tidak perlu API key!
3. Selalu jalankan `go mod tidy` setelah generating kode
4. Auto-setup menjalankan `go mod download` secara otomatis, tapi verifikasi berhasil
5. AI Mode bersifat opsional - Template Mode mencakup sebagian besar use case
6. Set environment variables sebelum menjalankan
7. Periksa database berjalan
8. Baca pesan error - biasanya memberi tahu apa yang salah!
:::
