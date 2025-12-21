package ai

import (
	"context"
	"fmt"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GeminiProvider implements the Provider interface for Google Gemini API
type GeminiProvider struct {
	apiKey  string
	model   string
	timeout time.Duration
	retries int
}

// NewGeminiProvider creates a new Gemini provider
func NewGeminiProvider(apiKey, model string, timeout time.Duration, retries int) *GeminiProvider {
	if model == "" {
		model = "gemini-2.0-flash-exp" // Default to Gemini 2.0 Flash (free tier)
	}

	return &GeminiProvider{
		apiKey:  apiKey,
		model:   model,
		timeout: timeout,
		retries: retries,
	}
}

func (g *GeminiProvider) Name() string {
	return "gemini"
}

func (g *GeminiProvider) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	startTime := time.Now()

	// Create client
	client, err := genai.NewClient(ctx, option.WithAPIKey(g.apiKey))
	if err != nil {
		return nil, fmt.Errorf("create gemini client: %w", err)
	}
	defer client.Close()

	// Get model
	model := client.GenerativeModel(g.model)

	// Configure generation
	model.SetTemperature(float32(req.Temperature))
	model.SetTopP(float32(req.TopP))
	model.SetMaxOutputTokens(int32(req.MaxTokens))

	// Combine system and user prompts
	fullPrompt := req.SystemPrompt + "\n\n" + req.UserPrompt

	// Retry logic
	var lastErr error
	for attempt := 0; attempt <= g.retries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := time.Duration(attempt) * time.Second
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		// Generate content
		resp, err := model.GenerateContent(ctx, genai.Text(fullPrompt))
		if err == nil && resp != nil {
			return g.parseResponse(resp, startTime), nil
		}

		lastErr = err
	}

	return nil, fmt.Errorf("gemini request failed after %d retries: %w", g.retries, lastErr)
}

func (g *GeminiProvider) parseResponse(resp *genai.GenerateContentResponse, startTime time.Time) *GenerateResponse {
	// Extract text from response
	var content string
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				content += string(txt)
			}
		}
	}

	// Extract token usage
	var promptTokens, completionTokens int
	if resp.UsageMetadata != nil {
		promptTokens = int(resp.UsageMetadata.PromptTokenCount)
		completionTokens = int(resp.UsageMetadata.CandidatesTokenCount)
	}

	// Gemini 1.5 Flash is FREE
	// Gemini 1.5 Pro pricing: $0.35/1M input, $1.05/1M output
	var cost float64
	if g.model == "gemini-1.5-pro" || g.model == "gemini-1.5-pro-latest" {
		inputCost := float64(promptTokens) * 0.35 / 1_000_000
		outputCost := float64(completionTokens) * 1.05 / 1_000_000
		cost = inputCost + outputCost
	} else {
		// Flash and other models are free
		cost = 0.0
	}

	// Determine finish reason
	finishReason := "stop"
	if len(resp.Candidates) > 0 {
		switch resp.Candidates[0].FinishReason {
		case genai.FinishReasonStop:
			finishReason = "stop"
		case genai.FinishReasonMaxTokens:
			finishReason = "length"
		case genai.FinishReasonSafety:
			finishReason = "content_filter"
		case genai.FinishReasonRecitation:
			finishReason = "recitation"
		default:
			finishReason = "unknown"
		}
	}

	return &GenerateResponse{
		Content:  content,
		Provider: "gemini",
		Model:    g.model,
		Duration: time.Since(startTime),
		TokensUsed: TokenUsage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      promptTokens + completionTokens,
		},
		Cost:         cost,
		FinishReason: finishReason,
		CacheHit:     false,
	}
}

func (g *GeminiProvider) Validate() error {
	if g.apiKey == "" {
		return fmt.Errorf("gemini API key is required")
	}
	return nil
}

func (g *GeminiProvider) Health(ctx context.Context) error {
	req := &GenerateRequest{
		SystemPrompt: "You are a test assistant.",
		UserPrompt:   "Respond with 'OK'",
		Temperature:  0.0,
		MaxTokens:    10,
	}

	_, err := g.Generate(ctx, req)
	return err
}

func (g *GeminiProvider) EstimateCost(req *GenerateRequest) (float64, error) {
	// Flash models are free
	if g.model == "gemini-1.5-flash" || g.model == "gemini-1.5-flash-latest" {
		return 0.0, nil // Free tier
	}

	// Rough estimation for Pro model
	estimatedInputTokens := len(req.SystemPrompt+req.UserPrompt) / 4
	estimatedOutputTokens := req.MaxTokens

	inputCost := float64(estimatedInputTokens) * 0.35 / 1_000_000
	outputCost := float64(estimatedOutputTokens) * 1.05 / 1_000_000

	return inputCost + outputCost, nil
}
