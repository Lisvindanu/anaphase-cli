package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lisvindanu/anaphase-cli/pkg/fileutil"
)

// MiddlewareType represents the type of middleware to generate
type MiddlewareType string

const (
	MiddlewareAuth      MiddlewareType = "auth"
	MiddlewareRateLimit MiddlewareType = "ratelimit"
	MiddlewareLogging   MiddlewareType = "logging"
	MiddlewareCORS      MiddlewareType = "cors"
)

// MiddlewareGenerator generates middleware code
type MiddlewareGenerator struct {
	middlewareType MiddlewareType
	outputDir      string
	packageName    string
}

// NewMiddlewareGenerator creates a new middleware generator
func NewMiddlewareGenerator(middlewareType MiddlewareType, outputDir string) *MiddlewareGenerator {
	return &MiddlewareGenerator{
		middlewareType: middlewareType,
		outputDir:      outputDir,
		packageName:    "middleware",
	}
}

// Generate creates middleware files
func (g *MiddlewareGenerator) Generate() ([]string, error) {
	var generatedFiles []string

	// Ensure output directory exists
	if err := fileutil.EnsureDir(g.outputDir); err != nil {
		return nil, fmt.Errorf("ensure directory: %w", err)
	}

	// Generate middleware based on type
	switch g.middlewareType {
	case MiddlewareAuth:
		file, err := g.generateAuthMiddleware()
		if err != nil {
			return nil, err
		}
		generatedFiles = append(generatedFiles, file)

	case MiddlewareRateLimit:
		file, err := g.generateRateLimitMiddleware()
		if err != nil {
			return nil, err
		}
		generatedFiles = append(generatedFiles, file)

	case MiddlewareLogging:
		file, err := g.generateLoggingMiddleware()
		if err != nil {
			return nil, err
		}
		generatedFiles = append(generatedFiles, file)

	case MiddlewareCORS:
		file, err := g.generateCORSMiddleware()
		if err != nil {
			return nil, err
		}
		generatedFiles = append(generatedFiles, file)

	default:
		return nil, fmt.Errorf("unknown middleware type: %s", g.middlewareType)
	}

	return generatedFiles, nil
}

func (g *MiddlewareGenerator) generateAuthMiddleware() (string, error) {
	filename := filepath.Join(g.outputDir, "auth.go")

	content := `package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents JWT claims
type Claims struct {
	UserID string ` + "`json:\"user_id\"`" + `
	Email  string ` + "`json:\"email\"`" + `
	Role   string ` + "`json:\"role\"`" + `
	jwt.RegisteredClaims
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	SecretKey     string
	TokenHeader   string // Default: "Authorization"
	TokenPrefix   string // Default: "Bearer "
	SkipPaths     []string
	ContextKey    string // Default: "user"
}

// AuthMiddleware creates a JWT authentication middleware
func AuthMiddleware(config AuthConfig) func(http.Handler) http.Handler {
	// Set defaults
	if config.TokenHeader == "" {
		config.TokenHeader = "Authorization"
	}
	if config.TokenPrefix == "" {
		config.TokenPrefix = "Bearer "
	}
	if config.ContextKey == "" {
		config.ContextKey = "user"
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if path should skip authentication
			for _, path := range config.SkipPaths {
				if r.URL.Path == path {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Extract token from header
			authHeader := r.Header.Get(config.TokenHeader)
			if authHeader == "" {
				http.Error(w, "Missing authentication token", http.StatusUnauthorized)
				return
			}

			// Remove "Bearer " prefix
			tokenString := strings.TrimPrefix(authHeader, config.TokenPrefix)
			if tokenString == authHeader {
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}

			// Parse and validate token
			token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				// Validate signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(config.SecretKey), nil
			})

			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Extract claims
			claims, ok := token.Claims.(*Claims)
			if !ok || !token.Valid {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			// Add claims to context
			ctx := context.WithValue(r.Context(), config.ContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetClaimsFromContext extracts claims from request context
func GetClaimsFromContext(ctx context.Context, key string) (*Claims, bool) {
	if key == "" {
		key = "user"
	}
	claims, ok := ctx.Value(key).(*Claims)
	return claims, ok
}

// RequireRole creates a middleware that requires specific role
func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetClaimsFromContext(r.Context(), "user")
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if user has required role
			roleAllowed := false
			for _, role := range allowedRoles {
				if claims.Role == role {
					roleAllowed = true
					break
				}
			}

			if !roleAllowed {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
`

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return "", err
	}

	return filename, nil
}

func (g *MiddlewareGenerator) generateRateLimitMiddleware() (string, error) {
	filename := filepath.Join(g.outputDir, "ratelimit.go")

	content := `package middleware

import (
	"net/http"
	"sync"
	"time"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*bucket
	rate     int           // requests per interval
	interval time.Duration
	maxBurst int
}

type bucket struct {
	tokens   int
	lastSeen time.Time
}

// NewRateLimiter creates a new rate limiter
// rate: number of requests allowed per interval
// interval: time window for rate limiting (e.g., 1 minute)
// maxBurst: maximum number of requests in a burst
func NewRateLimiter(rate int, interval time.Duration, maxBurst int) *RateLimiter {
	rl := &RateLimiter{
		buckets:  make(map[string]*bucket),
		rate:     rate,
		interval: interval,
		maxBurst: maxBurst,
	}

	// Cleanup goroutine to remove old buckets
	go rl.cleanup()

	return rl
}

// Allow checks if a request from the given key is allowed
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Get or create bucket for this key
	b, exists := rl.buckets[key]
	if !exists {
		b = &bucket{
			tokens:   rl.maxBurst,
			lastSeen: now,
		}
		rl.buckets[key] = b
	}

	// Calculate tokens to add based on elapsed time
	elapsed := now.Sub(b.lastSeen)
	tokensToAdd := int(elapsed / rl.interval * time.Duration(rl.rate))

	// Update bucket
	b.tokens += tokensToAdd
	if b.tokens > rl.maxBurst {
		b.tokens = rl.maxBurst
	}
	b.lastSeen = now

	// Check if we have tokens available
	if b.tokens > 0 {
		b.tokens--
		return true
	}

	return false
}

// cleanup removes old buckets periodically
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, b := range rl.buckets {
			if now.Sub(b.lastSeen) > 10*time.Minute {
				delete(rl.buckets, key)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Rate       int           // Requests per interval
	Interval   time.Duration // Time window
	MaxBurst   int           // Maximum burst size
	KeyFunc    func(*http.Request) string // Function to extract key from request
	OnLimitExceeded func(http.ResponseWriter, *http.Request) // Custom handler for rate limit exceeded
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(config RateLimitConfig) func(http.Handler) http.Handler {
	limiter := NewRateLimiter(config.Rate, config.Interval, config.MaxBurst)

	// Default key function: use IP address
	if config.KeyFunc == nil {
		config.KeyFunc = func(r *http.Request) string {
			// Try X-Forwarded-For first (for proxies)
			if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
				return xff
			}
			// Fall back to RemoteAddr
			return r.RemoteAddr
		}
	}

	// Default handler for rate limit exceeded
	if config.OnLimitExceeded == nil {
		config.OnLimitExceeded = func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Rate limit exceeded. Please try again later.", http.StatusTooManyRequests)
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := config.KeyFunc(r)

			if !limiter.Allow(key) {
				config.OnLimitExceeded(w, r)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Example usage:
// middleware := RateLimitMiddleware(RateLimitConfig{
//     Rate:     100,                    // 100 requests
//     Interval: time.Minute,            // per minute
//     MaxBurst: 120,                    // allow burst of 120
//     KeyFunc: func(r *http.Request) string {
//         // Rate limit by user ID if authenticated
//         if userID := r.Header.Get("X-User-ID"); userID != "" {
//             return userID
//         }
//         // Otherwise by IP
//         return r.RemoteAddr
//     },
// })
`

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return "", err
	}

	return filename, nil
}

func (g *MiddlewareGenerator) generateLoggingMiddleware() (string, error) {
	filename := filepath.Join(g.outputDir, "logging.go")

	content := `package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.written += n
	return n, err
}

// LoggingConfig holds logging middleware configuration
type LoggingConfig struct {
	Logger       *slog.Logger
	SkipPaths    []string
	LogRequestBody  bool
	LogResponseBody bool
}

// LoggingMiddleware creates a structured logging middleware
func LoggingMiddleware(config LoggingConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip logging for certain paths
			for _, path := range config.SkipPaths {
				if r.URL.Path == path {
					next.ServeHTTP(w, r)
					return
				}
			}

			start := time.Now()

			// Wrap response writer to capture status code
			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     0,
				written:        0,
			}

			// Process request
			next.ServeHTTP(wrapped, r)

			// Calculate duration
			duration := time.Since(start)

			// Build log attributes
			attrs := []slog.Attr{
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.Int("status", wrapped.statusCode),
				slog.Duration("duration", duration),
				slog.Int("response_size", wrapped.written),
			}

			// Add query parameters if present
			if r.URL.RawQuery != "" {
				attrs = append(attrs, slog.String("query", r.URL.RawQuery))
			}

			// Add user agent
			if ua := r.Header.Get("User-Agent"); ua != "" {
				attrs = append(attrs, slog.String("user_agent", ua))
			}

			// Add request ID if present
			if reqID := r.Header.Get("X-Request-ID"); reqID != "" {
				attrs = append(attrs, slog.String("request_id", reqID))
			}

			// Determine log level based on status code
			var logLevel slog.Level
			switch {
			case wrapped.statusCode >= 500:
				logLevel = slog.LevelError
			case wrapped.statusCode >= 400:
				logLevel = slog.LevelWarn
			default:
				logLevel = slog.LevelInfo
			}

			// Log the request
			config.Logger.LogAttrs(
				r.Context(),
				logLevel,
				"HTTP request",
				attrs...,
			)
		})
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				// Generate simple request ID (in production, use UUID)
				requestID = time.Now().Format("20060102150405") + "-" + r.RemoteAddr
			}

			// Add to response header
			w.Header().Set("X-Request-ID", requestID)

			// Add to request header for downstream handlers
			r.Header.Set("X-Request-ID", requestID)

			next.ServeHTTP(w, r)
		})
	}
}
`

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return "", err
	}

	return filename, nil
}

func (g *MiddlewareGenerator) generateCORSMiddleware() (string, error) {
	filename := filepath.Join(g.outputDir, "cors.go")

	content := `package middleware

import (
	"net/http"
	"strings"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string // e.g., ["https://example.com", "http://localhost:3000"]
	AllowedMethods   []string // e.g., ["GET", "POST", "PUT", "DELETE"]
	AllowedHeaders   []string // e.g., ["Authorization", "Content-Type"]
	ExposedHeaders   []string // Headers that browsers are allowed to access
	AllowCredentials bool     // Allow cookies and auth headers
	MaxAge           int      // Preflight cache duration in seconds (default: 86400 = 24h)
}

// DefaultCORSConfig returns a permissive CORS configuration for development
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: false,
		MaxAge:           86400, // 24 hours
	}
}

// ProductionCORSConfig returns a stricter CORS configuration
// You should customize allowedOrigins for your production domains
func ProductionCORSConfig(allowedOrigins []string) CORSConfig {
	return CORSConfig{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           3600, // 1 hour
	}
}

// CORSMiddleware creates a CORS middleware
func CORSMiddleware(config CORSConfig) func(http.Handler) http.Handler {
	// Set defaults
	if len(config.AllowedMethods) == 0 {
		config.AllowedMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	}
	if len(config.AllowedHeaders) == 0 {
		config.AllowedHeaders = []string{"Accept", "Content-Type"}
	}
	if config.MaxAge == 0 {
		config.MaxAge = 86400
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			allowed := false
			if len(config.AllowedOrigins) > 0 {
				for _, allowedOrigin := range config.AllowedOrigins {
					if allowedOrigin == "*" || allowedOrigin == origin {
						allowed = true
						break
					}
				}
			}

			if allowed {
				// Set CORS headers
				if len(config.AllowedOrigins) == 1 && config.AllowedOrigins[0] == "*" {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				}

				w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))

				if len(config.ExposedHeaders) > 0 {
					w.Header().Set("Access-Control-Expose-Headers", strings.Join(config.ExposedHeaders, ", "))
				}

				if config.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}

				if config.MaxAge > 0 {
					w.Header().Set("Access-Control-Max-Age", string(rune(config.MaxAge)))
				}
			}

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Example usage:
//
// Development:
// middleware := CORSMiddleware(DefaultCORSConfig())
//
// Production:
// middleware := CORSMiddleware(ProductionCORSConfig([]string{
//     "https://example.com",
//     "https://app.example.com",
// }))
`

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return "", err
	}

	return filename, nil
}
