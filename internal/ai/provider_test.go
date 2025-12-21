package ai

import (
	"testing"
	"time"
)

func TestGenerateRequest(t *testing.T) {
	req := &GenerateRequest{
		SystemPrompt: "You are a test assistant",
		UserPrompt:   "Hello",
		Temperature:  0.7,
		MaxTokens:    100,
		TopP:         0.9,
		Metadata:     map[string]string{"test": "value"},
	}

	if req.SystemPrompt == "" {
		t.Error("SystemPrompt should not be empty")
	}

	if req.Temperature < 0 || req.Temperature > 1 {
		t.Error("Temperature should be between 0 and 1")
	}
}

func TestGenerateResponse(t *testing.T) {
	resp := &GenerateResponse{
		Content:  "Test response",
		Provider: "test",
		Model:    "test-model",
		Duration: 100 * time.Millisecond,
		TokensUsed: TokenUsage{
			PromptTokens:     10,
			CompletionTokens: 20,
			TotalTokens:      30,
		},
		Cost:     0.001,
		CacheHit: false,
	}

	if resp.Content == "" {
		t.Error("Content should not be empty")
	}

	if resp.TokensUsed.TotalTokens != 30 {
		t.Errorf("Expected total tokens 30, got %d", resp.TokensUsed.TotalTokens)
	}

	expectedTotal := resp.TokensUsed.PromptTokens + resp.TokensUsed.CompletionTokens
	if resp.TokensUsed.TotalTokens != expectedTotal {
		t.Errorf("Total tokens mismatch: expected %d, got %d", expectedTotal, resp.TokensUsed.TotalTokens)
	}
}

func TestTokenUsage(t *testing.T) {
	usage := TokenUsage{
		PromptTokens:     100,
		CompletionTokens: 200,
		TotalTokens:      300,
	}

	if usage.PromptTokens+usage.CompletionTokens != usage.TotalTokens {
		t.Error("Token usage calculation is incorrect")
	}
}
