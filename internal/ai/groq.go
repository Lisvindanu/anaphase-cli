package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GroqProvider implements the Provider interface for Groq API
type GroqProvider struct {
	apiKey  string
	model   string
	timeout time.Duration
	retries int
	baseURL string
}

// NewGroqProvider creates a new Groq provider
func NewGroqProvider(apiKey, model string, timeout time.Duration, retries int) *GroqProvider {
	if model == "" {
		model = "llama-3.3-70b-versatile" // Default to Llama 3.3 70B (fast and capable)
	}

	return &GroqProvider{
		apiKey:  apiKey,
		model:   model,
		timeout: timeout,
		retries: retries,
		baseURL: "https://api.groq.com/openai/v1",
	}
}

func (g *GroqProvider) Name() string {
	return "groq"
}

// GroqRequest represents a request to the Groq API
type groqRequest struct {
	Model       string          `json:"model"`
	Messages    []groqMessage   `json:"messages"`
	Temperature float64         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	TopP        float64         `json:"top_p"`
	Stream      bool            `json:"stream"`
}

type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GroqResponse represents a response from the Groq API
type groqResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
	XGroq             struct {
		ID string `json:"id"`
	} `json:"x_groq"`
}

type groqErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

func (g *GroqProvider) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	startTime := time.Now()

	// Build messages
	messages := []groqMessage{
		{Role: "system", Content: req.SystemPrompt},
		{Role: "user", Content: req.UserPrompt},
	}

	// Prepare request
	groqReq := groqRequest{
		Model:       g.model,
		Messages:    messages,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		TopP:        req.TopP,
		Stream:      false,
	}

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

		resp, err := g.makeRequest(ctx, groqReq)
		if err == nil {
			return g.parseResponse(resp, startTime), nil
		}

		lastErr = err
	}

	return nil, fmt.Errorf("groq request failed after %d retries: %w", g.retries, lastErr)
}

func (g *GroqProvider) makeRequest(ctx context.Context, req groqRequest) (*groqResponse, error) {
	// Marshal request
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", g.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+g.apiKey)

	// Make request with timeout
	client := &http.Client{Timeout: g.timeout}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer httpResp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// Check for errors
	if httpResp.StatusCode != http.StatusOK {
		var errResp groqErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error.Message != "" {
			return nil, fmt.Errorf("groq API error (%d): %s", httpResp.StatusCode, errResp.Error.Message)
		}
		return nil, fmt.Errorf("groq API error (%d): %s", httpResp.StatusCode, string(respBody))
	}

	// Parse response
	var resp groqResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &resp, nil
}

func (g *GroqProvider) parseResponse(resp *groqResponse, startTime time.Time) *GenerateResponse {
	// Extract content
	var content string
	var finishReason string
	if len(resp.Choices) > 0 {
		content = resp.Choices[0].Message.Content
		finishReason = resp.Choices[0].FinishReason
	}

	// Groq is FREE for most models!
	// Pricing varies by model, but generally very cheap or free during preview
	cost := 0.0 // Free for now

	return &GenerateResponse{
		Content:  content,
		Provider: "groq",
		Model:    resp.Model,
		Duration: time.Since(startTime),
		TokensUsed: TokenUsage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
		Cost:         cost,
		FinishReason: finishReason,
		CacheHit:     false,
	}
}

func (g *GroqProvider) Validate() error {
	if g.apiKey == "" {
		return fmt.Errorf("groq API key is required")
	}
	return nil
}

func (g *GroqProvider) Health(ctx context.Context) error {
	req := &GenerateRequest{
		SystemPrompt: "You are a test assistant.",
		UserPrompt:   "Respond with 'OK'",
		Temperature:  0.0,
		MaxTokens:    10,
	}

	_, err := g.Generate(ctx, req)
	return err
}

func (g *GroqProvider) EstimateCost(req *GenerateRequest) (float64, error) {
	// Groq is currently free or very cheap
	return 0.0, nil
}
