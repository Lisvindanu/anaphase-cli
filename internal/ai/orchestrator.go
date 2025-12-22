package ai

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// Orchestrator manages multiple AI providers with fallback logic
type Orchestrator struct {
	providers       map[string]Provider
	primaryProvider string
	fallbackChain   []string
	cache           *Cache
	logger          *slog.Logger
}

// NewOrchestrator creates a new orchestrator
func NewOrchestrator(cfg *Config, logger *slog.Logger) (*Orchestrator, error) {
	// Initialize cache
	cache := NewCache(cfg.Cache.Directory, cfg.Cache.TTL, cfg.Cache.Enabled)

	// Initialize providers
	providerMap := make(map[string]Provider)

	// Gemini
	if cfg.AI.Providers.Gemini.Enabled && cfg.AI.Providers.Gemini.APIKey != "" {
		providerMap["gemini"] = NewGeminiProvider(
			cfg.AI.Providers.Gemini.APIKey,
			cfg.AI.Providers.Gemini.Model,
			cfg.AI.Providers.Gemini.Timeout,
			cfg.AI.Providers.Gemini.MaxRetries,
		)
	}

	// Groq
	if cfg.AI.Providers.Groq.Enabled && cfg.AI.Providers.Groq.APIKey != "" {
		providerMap["groq"] = NewGroqProvider(
			cfg.AI.Providers.Groq.APIKey,
			cfg.AI.Providers.Groq.Model,
			cfg.AI.Providers.Groq.Timeout,
			cfg.AI.Providers.Groq.MaxRetries,
		)
	}

	// TODO: Add OpenAI provider

	// Validate at least one provider is available
	if len(providerMap) == 0 {
		return nil, fmt.Errorf("no AI providers configured - please set at least one API key")
	}

	return &Orchestrator{
		providers:       providerMap,
		primaryProvider: cfg.AI.PrimaryProvider,
		fallbackChain:   cfg.AI.FallbackProviders,
		cache:           cache,
		logger:          logger,
	}, nil
}

// Generate attempts generation with fallback logic
func (o *Orchestrator) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	// Check cache first
	if cached, hit := o.cache.Get(req); hit {
		o.logger.Info("cache hit",
			"provider", cached.Provider,
			"tokens", cached.TokensUsed.TotalTokens,
		)
		return cached, nil
	}

	// Build provider chain (primary + fallbacks)
	providerChain := []string{o.primaryProvider}
	providerChain = append(providerChain, o.fallbackChain...)

	var lastErr error
	for _, providerName := range providerChain {
		provider, exists := o.providers[providerName]
		if !exists {
			o.logger.Warn("provider not available",
				"provider", providerName,
			)
			continue
		}

		o.logger.Info("attempting generation",
			"provider", providerName,
		)

		startTime := time.Now()
		resp, err := provider.Generate(ctx, req)

		if err != nil {
			o.logger.Warn("provider failed",
				"provider", providerName,
				"error", err,
				"duration", time.Since(startTime),
			)
			lastErr = err
			continue
		}

		// Success!
		o.logger.Info("generation successful",
			"provider", providerName,
			"tokens", resp.TokensUsed.TotalTokens,
			"cost", fmt.Sprintf("$%.6f", resp.Cost),
			"duration", resp.Duration,
		)

		// Cache the response
		if err := o.cache.Set(req, resp); err != nil {
			o.logger.Warn("failed to cache response", "error", err)
		}

		return resp, nil
	}

	return nil, fmt.Errorf("all providers failed, last error: %w", lastErr)
}

// ValidateProviders checks health of all configured providers
func (o *Orchestrator) ValidateProviders(ctx context.Context) map[string]error {
	results := make(map[string]error)

	for name, provider := range o.providers {
		if err := provider.Validate(); err != nil {
			results[name] = fmt.Errorf("validation failed: %w", err)
			continue
		}

		if err := provider.Health(ctx); err != nil {
			results[name] = fmt.Errorf("health check failed: %w", err)
			continue
		}

		results[name] = nil // Success
	}

	return results
}

// EstimateCost estimates cost for a generation request
func (o *Orchestrator) EstimateCost(req *GenerateRequest) (float64, error) {
	provider, exists := o.providers[o.primaryProvider]
	if !exists {
		return 0, fmt.Errorf("primary provider not available")
	}

	return provider.EstimateCost(req)
}
