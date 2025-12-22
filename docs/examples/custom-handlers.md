# Custom Handlers and Middleware

Learn how to customize generated handlers and add middleware for authentication, validation, logging, and more.

## Overview

Anaphase generates basic CRUD handlers, but real applications need:

- **Authentication** - JWT, OAuth, API keys
- **Authorization** - Role-based access control
- **Validation** - Request validation and sanitization
- **Rate Limiting** - Protect against abuse
- **Logging** - Structured logging with correlation IDs
- **Error Handling** - Consistent error responses
- **Caching** - Response caching
- **CORS** - Cross-origin resource sharing

This guide shows how to extend generated handlers with these features.

## Authentication Middleware

### JWT Authentication

```go
// internal/middleware/auth.go
package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey contextKey = "user"

type Claims struct {
    UserID string   `json:"user_id"`
    Email  string   `json:"email"`
    Role   string   `json:"role"`
    jwt.RegisteredClaims
}

func JWTAuth(jwtSecret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract token from header
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }

            tokenString := strings.TrimPrefix(authHeader, "Bearer ")

            // Parse and validate token
            claims := &Claims{}
            token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
                return []byte(jwtSecret), nil
            })

            if err != nil || !token.Valid {
                http.Error(w, "invalid token", http.StatusUnauthorized)
                return
            }

            // Add user info to context
            ctx := context.WithValue(r.Context(), UserContextKey, claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// Helper to get user from context
func GetUser(ctx context.Context) (*Claims, bool) {
    user, ok := ctx.Value(UserContextKey).(*Claims)
    return user, ok
}
```

### Apply to Routes

```go
// internal/adapter/handler/http/customer_handler.go
func (h *CustomerHandler) RegisterRoutes(r chi.Router) {
    r.Group(func(r chi.Router) {
        // Apply auth middleware
        r.Use(middleware.JWTAuth(os.Getenv("JWT_SECRET")))

        r.Post("/customers", h.Create)
        r.Get("/customers/{id}", h.GetByID)
        r.Put("/customers/{id}", h.Update)
        r.Delete("/customers/{id}", h.Delete)
    })
}
```

### Use in Handler

```go
func (h *CustomerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
    // Get authenticated user
    user, ok := middleware.GetUser(r.Context())
    if !ok {
        h.respondError(w, http.StatusUnauthorized, "unauthorized")
        return
    }

    // Check authorization
    customerID := chi.URLParam(r, "id")
    if user.UserID != customerID && user.Role != "admin" {
        h.respondError(w, http.StatusForbidden, "forbidden")
        return
    }

    // Continue with normal logic
    customer, err := h.service.GetCustomer(r.Context(), customerID)
    if err != nil {
        h.respondError(w, http.StatusNotFound, "customer not found")
        return
    }

    h.respondJSON(w, http.StatusOK, customer)
}
```

## Request Validation

### Using go-playground/validator

```go
// internal/adapter/handler/http/customer_handler.go
import "github.com/go-playground/validator/v10"

var validate = validator.New()

type CreateCustomerRequest struct {
    Email     string  `json:"email" validate:"required,email"`
    FirstName string  `json:"firstName" validate:"required,min=2,max=50"`
    LastName  string  `json:"lastName" validate:"required,min=2,max=50"`
    Phone     string  `json:"phone" validate:"required,e164"`
    Age       int     `json:"age" validate:"gte=18,lte=120"`
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req CreateCustomerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.respondError(w, http.StatusBadRequest, "invalid request body")
        return
    }

    // Validate request
    if err := validate.Struct(req); err != nil {
        validationErrors := err.(validator.ValidationErrors)
        h.respondValidationError(w, validationErrors)
        return
    }

    // Process request...
}

func (h *CustomerHandler) respondValidationError(w http.ResponseWriter, errors validator.ValidationErrors) {
    errMap := make(map[string]string)
    for _, err := range errors {
        errMap[err.Field()] = formatValidationError(err)
    }

    h.respondJSON(w, http.StatusBadRequest, map[string]interface{}{
        "error": "validation failed",
        "fields": errMap,
    })
}

func formatValidationError(err validator.FieldError) string {
    switch err.Tag() {
    case "required":
        return "this field is required"
    case "email":
        return "invalid email format"
    case "min":
        return fmt.Sprintf("must be at least %s characters", err.Param())
    case "max":
        return fmt.Sprintf("must be at most %s characters", err.Param())
    default:
        return "invalid value"
    }
}
```

## Rate Limiting

### Using golang.org/x/time/rate

```go
// internal/middleware/ratelimit.go
package middleware

import (
    "net/http"
    "sync"
    "time"

    "golang.org/x/time/rate"
)

type visitor struct {
    limiter  *rate.Limiter
    lastSeen time.Time
}

type RateLimiter struct {
    visitors map[string]*visitor
    mu       sync.RWMutex
    rate     rate.Limit
    burst    int
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
    rl := &RateLimiter{
        visitors: make(map[string]*visitor),
        rate:     r,
        burst:    b,
    }

    // Cleanup goroutine
    go rl.cleanupVisitors()

    return rl
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
    rl.mu.Lock()
    defer rl.mu.Unlock()

    v, exists := rl.visitors[ip]
    if !exists {
        limiter := rate.NewLimiter(rl.rate, rl.burst)
        rl.visitors[ip] = &visitor{limiter, time.Now()}
        return limiter
    }

    v.lastSeen = time.Now()
    return v.limiter
}

func (rl *RateLimiter) cleanupVisitors() {
    for {
        time.Sleep(time.Minute)

        rl.mu.Lock()
        for ip, v := range rl.visitors {
            if time.Since(v.lastSeen) > 3*time.Minute {
                delete(rl.visitors, ip)
            }
        }
        rl.mu.Unlock()
    }
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr
        limiter := rl.getVisitor(ip)

        if !limiter.Allow() {
            http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

### Apply Rate Limiting

```go
func main() {
    // 10 requests per second with burst of 20
    rateLimiter := middleware.NewRateLimiter(10, 20)

    r := chi.NewRouter()
    r.Use(rateLimiter.Limit)

    // Register routes...
}
```

## Structured Logging

### Request Logging with Correlation IDs

```go
// internal/middleware/logging.go
package middleware

import (
    "log/slog"
    "net/http"
    "time"

    "github.com/google/uuid"
)

type responseWriter struct {
    http.ResponseWriter
    status int
    bytes  int
}

func (rw *responseWriter) WriteHeader(status int) {
    rw.status = status
    rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
    n, err := rw.ResponseWriter.Write(b)
    rw.bytes += n
    return n, err
}

func RequestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()

            // Generate correlation ID
            correlationID := uuid.New().String()
            r.Header.Set("X-Correlation-ID", correlationID)
            w.Header().Set("X-Correlation-ID", correlationID)

            // Wrap response writer
            rw := &responseWriter{
                ResponseWriter: w,
                status:         http.StatusOK,
            }

            // Add correlation ID to context
            ctx := r.Context()
            logger := logger.With(
                "correlation_id", correlationID,
                "method", r.Method,
                "path", r.URL.Path,
                "remote_addr", r.RemoteAddr,
            )

            // Process request
            next.ServeHTTP(rw, r)

            // Log request
            logger.Info("request completed",
                "status", rw.status,
                "bytes", rw.bytes,
                "duration_ms", time.Since(start).Milliseconds(),
            )
        })
    }
}
```

## Error Handling

### Consistent Error Responses

```go
// internal/adapter/handler/http/errors.go
package http

import (
    "encoding/json"
    "log/slog"
    "net/http"
)

type ErrorResponse struct {
    Error          string            `json:"error"`
    Message        string            `json:"message"`
    CorrelationID  string            `json:"correlation_id,omitempty"`
    ValidationErrors map[string]string `json:"validation_errors,omitempty"`
}

func (h *BaseHandler) respondError(w http.ResponseWriter, r *http.Request, status int, message string) {
    correlationID := r.Header.Get("X-Correlation-ID")

    h.logger.Error("request error",
        "correlation_id", correlationID,
        "status", status,
        "error", message,
    )

    resp := ErrorResponse{
        Error:         http.StatusText(status),
        Message:       message,
        CorrelationID: correlationID,
    }

    h.respondJSON(w, status, resp)
}

func (h *BaseHandler) respondValidationError(w http.ResponseWriter, r *http.Request, errors map[string]string) {
    correlationID := r.Header.Get("X-Correlation-ID")

    resp := ErrorResponse{
        Error:            "validation_failed",
        Message:          "request validation failed",
        CorrelationID:    correlationID,
        ValidationErrors: errors,
    }

    h.respondJSON(w, http.StatusBadRequest, resp)
}
```

## Response Caching

### Redis Cache Middleware

```go
// internal/middleware/cache.go
package middleware

import (
    "context"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "net/http"
    "time"

    "github.com/redis/go-redis/v9"
)

type CacheMiddleware struct {
    redis *redis.Client
    ttl   time.Duration
}

func NewCacheMiddleware(redisURL string, ttl time.Duration) (*CacheMiddleware, error) {
    opt, err := redis.ParseURL(redisURL)
    if err != nil {
        return nil, err
    }

    client := redis.NewClient(opt)
    return &CacheMiddleware{
        redis: client,
        ttl:   ttl,
    }, nil
}

func (c *CacheMiddleware) Cache(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Only cache GET requests
        if r.Method != http.MethodGet {
            next.ServeHTTP(w, r)
            return
        }

        // Generate cache key
        key := c.cacheKey(r)

        // Try to get from cache
        cached, err := c.redis.Get(r.Context(), key).Result()
        if err == nil {
            w.Header().Set("X-Cache", "HIT")
            w.Write([]byte(cached))
            return
        }

        // Cache miss - capture response
        rw := &cachedResponseWriter{
            ResponseWriter: w,
            body:          make([]byte, 0),
        }

        next.ServeHTTP(rw, r)

        // Cache successful responses
        if rw.status >= 200 && rw.status < 300 {
            c.redis.Set(r.Context(), key, rw.body, c.ttl)
            w.Header().Set("X-Cache", "MISS")
        }
    })
}

func (c *CacheMiddleware) cacheKey(r *http.Request) string {
    h := sha256.New()
    h.Write([]byte(r.URL.Path))
    h.Write([]byte(r.URL.RawQuery))
    return "cache:" + hex.EncodeToString(h.Sum(nil))
}
```

## CORS Configuration

```go
// cmd/api/main.go
import "github.com/go-chi/cors"

func main() {
    r := chi.NewRouter()

    // CORS middleware
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"https://example.com", "http://localhost:3000"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
        ExposedHeaders:   []string{"Link", "X-Correlation-ID"},
        AllowCredentials: true,
        MaxAge:           300,
    }))

    // Register routes...
}
```

## Complete Example

```go
// cmd/api/main.go
package main

import (
    "log"
    "log/slog"
    "net/http"
    "os"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/cors"
    "golang.org/x/time/rate"

    "myapp/internal/middleware"
)

func main() {
    // Logger
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))

    // Rate limiter
    rateLimiter := middleware.NewRateLimiter(rate.Limit(10), 20)

    // Router
    r := chi.NewRouter()

    // Global middleware
    r.Use(middleware.RequestLogger(logger))
    r.Use(rateLimiter.Limit)
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
        ExposedHeaders:   []string{"X-Correlation-ID"},
        AllowCredentials: false,
        MaxAge:           300,
    }))

    // Public routes
    r.Post("/auth/login", loginHandler)
    r.Post("/auth/register", registerHandler)

    // Protected routes
    r.Group(func(r chi.Router) {
        r.Use(middleware.JWTAuth(os.Getenv("JWT_SECRET")))

        customerHandler.RegisterRoutes(r)
        productHandler.RegisterRoutes(r)
        orderHandler.RegisterRoutes(r)
    })

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    logger.Info("server starting", "port", port)
    if err := http.ListenAndServe(":"+port, r); err != nil {
        log.Fatal(err)
    }
}
```

## Testing Custom Handlers

```go
// internal/adapter/handler/http/customer_handler_test.go
package http

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestCustomerHandler_Create_WithAuth(t *testing.T) {
    // Setup
    mockService := new(MockCustomerService)
    handler := NewCustomerHandler(mockService, logger)

    // Mock JWT middleware
    authMiddleware := func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx := context.WithValue(r.Context(), middleware.UserContextKey, &middleware.Claims{
                UserID: "test-user-id",
                Role:   "user",
            })
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }

    // Setup router with middleware
    r := chi.NewRouter()
    r.Use(authMiddleware)
    handler.RegisterRoutes(r)

    // Create request
    req := httptest.NewRequest("POST", "/customers", strings.NewReader(`{
        "email": "test@example.com",
        "firstName": "John",
        "lastName": "Doe"
    }`))
    req.Header.Set("Content-Type", "application/json")

    // Execute
    rr := httptest.NewRecorder()
    r.ServeHTTP(rr, req)

    // Assert
    assert.Equal(t, http.StatusCreated, rr.Code)
    mockService.AssertExpectations(t)
}
```

## See Also

- [Basic Example](/examples/basic)
- [Multi-Domain Service](/examples/multi-domain)
- [Architecture Guide](/guide/architecture)
