# anaphase gen middleware

Generate production-ready HTTP middleware for your microservice.

::: info
**Quick Start**: Run `anaphase` (no arguments) to access the interactive menu and select "Generate Middleware" for a guided experience.
:::

## Overview

The `gen middleware` command generates common HTTP middleware patterns with best practices and configuration options. All middleware is framework-agnostic and works with the standard `net/http` package.

::: info
**Template-Based**: Middleware generation uses proven templates - no AI configuration required. Fast, reliable, and production-ready.
:::

## Usage

### Interactive Menu (Recommended)

```bash
anaphase
```

Select **"Generate Middleware"** from the menu and choose your middleware type from the visual interface.

### CLI Direct Mode

```bash
anaphase gen middleware --type <type> [flags]
```

## Available Middleware Types

::: tip
All middleware types are generated from templates - no AI or API keys needed.
:::

### 1. Authentication (JWT)

Generate JWT authentication middleware with role-based access control.

```bash
anaphase gen middleware --type auth
```

**Features:**
- JWT token validation
- Custom claims support
- Role-based authorization
- Configurable token header and prefix
- Skip paths for public endpoints

**Example Usage:**
```go
import "yourproject/internal/middleware"

config := middleware.AuthConfig{
    SecretKey:   os.Getenv("JWT_SECRET"),
    TokenHeader: "Authorization",
    TokenPrefix: "Bearer ",
    SkipPaths:   []string{"/health", "/login"},
}

router.Use(middleware.AuthMiddleware(config))

// Protect specific routes by role
adminRoutes := router.Group("/admin")
adminRoutes.Use(middleware.RequireRole("admin", "superadmin"))
```

### 2. Rate Limiting

Generate token bucket rate limiter middleware.

```bash
anaphase gen middleware --type ratelimit
```

**Features:**
- Token bucket algorithm
- Per-client rate limiting
- Configurable rate and burst size
- Custom key extraction (IP, user ID, etc.)
- Automatic cleanup of old buckets

**Example Usage:**
```go
import "yourproject/internal/middleware"

config := middleware.RateLimitConfig{
    Rate:     100,              // 100 requests
    Interval: time.Minute,      // per minute
    MaxBurst: 120,              // allow burst of 120
    KeyFunc: func(r *http.Request) string {
        // Rate limit by user ID if authenticated
        if userID := r.Header.Get("X-User-ID"); userID != "" {
            return userID
        }
        // Otherwise by IP
        return r.RemoteAddr
    },
}

router.Use(middleware.RateLimitMiddleware(config))
```

### 3. Structured Logging

Generate request/response logging middleware with structured logging.

```bash
anaphase gen middleware --type logging
```

**Features:**
- Structured logging with `log/slog`
- Request/response metrics
- Duration tracking
- Request ID correlation
- Configurable log levels based on status code

**Example Usage:**
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

// Add request ID first
router.Use(middleware.RequestIDMiddleware())
// Then logging
router.Use(middleware.LoggingMiddleware(config))
```

### 4. CORS

Generate CORS middleware with development and production presets.

```bash
anaphase gen middleware --type cors
```

**Features:**
- Preflight request handling
- Configurable allowed origins, methods, headers
- Credentials support
- Response header exposure
- Development and production presets

**Example Usage:**
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

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--type` | | (required) | Middleware type: auth, ratelimit, logging, cors |
| `--output` | | `internal/middleware` | Output directory for generated files |

## Examples

### Using Interactive Menu

```bash
# Launch menu
anaphase

# Select "Generate Middleware"
# Choose middleware type visually:
# → Authentication (JWT)
# → Rate Limiting
# → Logging
# → CORS
```

### Using CLI Directly

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

**Custom Output Directory:**

```bash
anaphase gen middleware --type auth --output pkg/middleware
```

## Middleware Chaining

Order matters when chaining middleware. Recommended order:

```go
// 1. CORS (handle preflight first)
router.Use(middleware.CORSMiddleware(corsConfig))

// 2. Request ID (for logging correlation)
router.Use(middleware.RequestIDMiddleware())

// 3. Logging (log all requests)
router.Use(middleware.LoggingMiddleware(logConfig))

// 4. Rate Limiting (before auth to prevent brute force)
router.Use(middleware.RateLimitMiddleware(rateLimitConfig))

// 5. Authentication (protect routes)
router.Use(middleware.AuthMiddleware(authConfig))

// 6. Your handlers
router.HandleFunc("/api/users", getUsersHandler)
```

## Best Practices

1. **Authentication**
   - Use environment variables for secrets
   - Rotate JWT keys regularly
   - Set appropriate token expiration
   - Use HTTPS in production

2. **Rate Limiting**
   - Set limits per endpoint based on usage
   - Use different limits for authenticated vs anonymous users
   - Monitor and adjust based on traffic patterns

3. **Logging**
   - Don't log sensitive data (passwords, tokens)
   - Use structured logging for better querying
   - Include request IDs for tracing
   - Set appropriate log levels

4. **CORS**
   - Never use `*` in production
   - Whitelist specific origins
   - Be restrictive with credentials
   - Test thoroughly with your frontend

## Integration Examples

### With Chi Router

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

### With Gorilla Mux

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

### With Standard Library

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

## See Also

- [anaphase gen domain](/reference/gen-domain) - Generate domain entities
- [anaphase gen handler](/reference/gen-handler) - Generate HTTP handlers
- [anaphase quality](/reference/quality) - Code quality tools
