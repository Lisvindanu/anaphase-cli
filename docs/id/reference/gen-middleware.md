# anaphase gen middleware

Generate HTTP middleware production-ready untuk microservice Anda.

::: info
**Quick Start**: Jalankan `anaphase` (tanpa argumen) untuk mengakses menu interaktif dan pilih "Generate Middleware" untuk pengalaman terpandu.
:::

## Overview

Command `gen middleware` generate pattern HTTP middleware umum dengan best practice dan opsi konfigurasi. Semua middleware bersifat framework-agnostic dan bekerja dengan package `net/http` standar.

::: info
**Berbasis Template**: Generasi middleware menggunakan template yang sudah terbukti - tidak perlu konfigurasi AI. Cepat, andal, dan production-ready.
:::

## Penggunaan

### Menu Interaktif (Disarankan)

```bash
anaphase
```

Pilih **"Generate Middleware"** dari menu dan pilih tipe middleware Anda dari interface visual.

### Mode CLI Langsung

```bash
anaphase gen middleware --type <type> [flags]
```

## Tipe Middleware yang Tersedia

::: tip
Semua tipe middleware dihasilkan dari template - tidak perlu AI atau API key.
:::

### 1. Authentication (JWT)

Generate middleware autentikasi JWT dengan role-based access control.

```bash
anaphase gen middleware --type auth
```

**Fitur:**
- Validasi JWT token
- Dukungan custom claims
- Otorisasi berbasis role
- Header token dan prefix yang dapat dikonfigurasi
- Skip path untuk endpoint publik

**Contoh Penggunaan:**
```go
import "yourproject/internal/middleware"

config := middleware.AuthConfig{
    SecretKey:   os.Getenv("JWT_SECRET"),
    TokenHeader: "Authorization",
    TokenPrefix: "Bearer ",
    SkipPaths:   []string{"/health", "/login"},
}

router.Use(middleware.AuthMiddleware(config))

// Lindungi route tertentu berdasarkan role
adminRoutes := router.Group("/admin")
adminRoutes.Use(middleware.RequireRole("admin", "superadmin"))
```

### 2. Rate Limiting

Generate middleware rate limiter token bucket.

```bash
anaphase gen middleware --type ratelimit
```

**Fitur:**
- Algoritma token bucket
- Rate limiting per-client
- Rate dan burst size yang dapat dikonfigurasi
- Ekstraksi key kustom (IP, user ID, dll.)
- Cleanup otomatis bucket lama

**Contoh Penggunaan:**
```go
import "yourproject/internal/middleware"

config := middleware.RateLimitConfig{
    Rate:     100,              // 100 request
    Interval: time.Minute,      // per menit
    MaxBurst: 120,              // allow burst of 120
    KeyFunc: func(r *http.Request) string {
        // Rate limit berdasarkan user ID jika terautentikasi
        if userID := r.Header.Get("X-User-ID"); userID != "" {
            return userID
        }
        // Atau berdasarkan IP
        return r.RemoteAddr
    },
}

router.Use(middleware.RateLimitMiddleware(config))
```

### 3. Structured Logging

Generate middleware logging request/response dengan structured logging.

```bash
anaphase gen middleware --type logging
```

**Fitur:**
- Structured logging dengan `log/slog`
- Metrik request/response
- Pelacakan durasi
- Korelasi request ID
- Level log yang dapat dikonfigurasi berdasarkan status code

**Contoh Penggunaan:**
```go
import (
    "log/slog"
    "yourproject/internal/middleware"
)

logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

config := middleware.LoggingConfig{
    Logger:    logger,
    SkipPaths: []string{"/health"},
}

// Tambahkan request ID terlebih dahulu
router.Use(middleware.RequestIDMiddleware())
// Kemudian logging
router.Use(middleware.LoggingMiddleware(config))
```

### 4. CORS

Generate middleware CORS dengan preset development dan production.

```bash
anaphase gen middleware --type cors
```

**Fitur:**
- Penanganan preflight request
- Allowed origins, methods, headers yang dapat dikonfigurasi
- Dukungan credentials
- Exposure response header
- Preset development dan production

**Contoh Penggunaan:**
```go
import "yourproject/internal/middleware"

// Development
router.Use(middleware.CORSMiddleware(middleware.DefaultCORSConfig()))

// Production
config := middleware.ProductionCORSConfig([]string{
    "https://example.com",
    "https://app.example.com",
})
router.Use(middleware.CORSMiddleware(config))
```

## Flag

| Flag | Short | Default | Deskripsi |
|------|-------|---------|-------------|
| `--type` | | (diperlukan) | Tipe middleware: auth, ratelimit, logging, cors |
| `--output` | | `internal/middleware` | Direktori output untuk file yang dihasilkan |

## Contoh

### Menggunakan Menu Interaktif

```bash
# Luncurkan menu
anaphase

# Pilih "Generate Middleware"
# Pilih tipe middleware secara visual:
# → Authentication (JWT)
# → Rate Limiting
# → Logging
# → CORS
```

### Menggunakan CLI Langsung

**Generate Auth Middleware:**

```bash
anaphase gen middleware --type auth
```

Output:
```
✓ /internal/middleware/auth.go
```

**Generate Multiple Middleware:**

```bash
anaphase gen middleware --type auth
anaphase gen middleware --type ratelimit
anaphase gen middleware --type logging
anaphase gen middleware --type cors
```

**Direktori Output Kustom:**

```bash
anaphase gen middleware --type auth --output pkg/middleware
```

## Chaining Middleware

Urutan penting saat chaining middleware. Urutan yang disarankan:

```go
// 1. CORS (handle preflight terlebih dahulu)
router.Use(middleware.CORSMiddleware(corsConfig))

// 2. Request ID (untuk korelasi logging)
router.Use(middleware.RequestIDMiddleware())

// 3. Logging (log semua request)
router.Use(middleware.LoggingMiddleware(logConfig))

// 4. Rate Limiting (sebelum auth untuk mencegah brute force)
router.Use(middleware.RateLimitMiddleware(rateLimitConfig))

// 5. Authentication (lindungi route)
router.Use(middleware.AuthMiddleware(authConfig))

// 6. Handler Anda
router.HandleFunc("/api/users", getUsersHandler)
```

## Best Practice

1. **Authentication**
   - Gunakan environment variable untuk secret
   - Rotasi JWT key secara berkala
   - Set token expiration yang sesuai
   - Gunakan HTTPS di production

2. **Rate Limiting**
   - Set limit per endpoint berdasarkan penggunaan
   - Gunakan limit berbeda untuk authenticated vs anonymous user
   - Monitor dan sesuaikan berdasarkan traffic pattern

3. **Logging**
   - Jangan log data sensitif (password, token)
   - Gunakan structured logging untuk query yang lebih baik
   - Sertakan request ID untuk tracing
   - Set level log yang sesuai

4. **CORS**
   - Jangan gunakan `*` di production
   - Whitelist origin tertentu
   - Bersikap restriktif dengan credentials
   - Test secara menyeluruh dengan frontend Anda

## Contoh Integrasi

### Dengan Chi Router

```go
import (
    "github.com/go-chi/chi/v5"
    "yourproject/internal/middleware"
)

r := chi.NewRouter()
r.Use(middleware.CORSMiddleware(corsConfig))
r.Use(middleware.LoggingMiddleware(logConfig))
r.Use(middleware.AuthMiddleware(authConfig))
```

### Dengan Gorilla Mux

```go
import (
    "github.com/gorilla/mux"
    "yourproject/internal/middleware"
)

r := mux.NewRouter()
r.Use(middleware.CORSMiddleware(corsConfig))
r.Use(middleware.LoggingMiddleware(logConfig))
r.Use(middleware.AuthMiddleware(authConfig))
```

### Dengan Standard Library

```go
handler := middleware.CORSMiddleware(corsConfig)(
    middleware.LoggingMiddleware(logConfig)(
        middleware.AuthMiddleware(authConfig)(
            yourHandler,
        ),
    ),
)

http.ListenAndServe(":8080", handler)
```

## Lihat Juga

- [anaphase gen domain](/reference/gen-domain) - Generate entity domain
- [anaphase gen handler](/reference/gen-handler) - Generate HTTP handler
- [anaphase quality](/reference/quality) - Code quality tools
