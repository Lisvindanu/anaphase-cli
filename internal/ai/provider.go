package ai

import (
	"context"
	"time"
)

// Provider defines the interface for AI service providers
type Provider interface {
	// Name returns the provider identifier
	Name() string

	// Generate sends a prompt and returns structured response
	Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)

	// Validate checks if provider is properly configured
	Validate() error

	// Health checks if provider is reachable
	Health(ctx context.Context) error

	// EstimateCost calculates approximate cost for a request
	EstimateCost(req *GenerateRequest) (float64, error)
}

// GenerateRequest encapsulates generation parameters
type GenerateRequest struct {
	SystemPrompt string            // System-level instructions
	UserPrompt   string            // User's actual request
	Temperature  float64           // Randomness (0.0-1.0)
	MaxTokens    int               // Maximum output length
	TopP         float64           // Nucleus sampling
	Metadata     map[string]string // Request metadata
}

// GenerateResponse contains the provider's output
type GenerateResponse struct {
	Content      string            // Generated text
	TokensUsed   TokenUsage        // Token consumption
	Provider     string            // Which provider was used
	Model        string            // Specific model used
	Duration     time.Duration     // Request duration
	CacheHit     bool              // Was response cached?
	Cost         float64           // Estimated cost in USD
	FinishReason string            // Why generation stopped
	Metadata     map[string]string // Response metadata
}

// TokenUsage tracks token consumption
type TokenUsage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}
