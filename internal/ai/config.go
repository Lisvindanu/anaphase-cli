package ai

import (
	"time"
)

// Config holds the AI configuration
type Config struct {
	AI        AIConfig    `yaml:"ai"`
	Cache     CacheConfig `yaml:"cache"`
	Generator GenConfig   `yaml:"generator"`
}

// AIConfig holds AI provider configuration
type AIConfig struct {
	PrimaryProvider   string          `yaml:"primary_provider"`
	FallbackProviders []string        `yaml:"fallback_providers"`
	Providers         ProvidersConfig `yaml:"providers"`
}

// ProvidersConfig holds individual provider configurations
type ProvidersConfig struct {
	Gemini ProviderConfig `yaml:"gemini"`
	Groq   ProviderConfig `yaml:"groq"`
	OpenAI ProviderConfig `yaml:"openai"`
	Claude ProviderConfig `yaml:"claude"`
	Ollama ProviderConfig `yaml:"ollama"`
}

// ProviderConfig holds configuration for a single provider
type ProviderConfig struct {
	Enabled    bool          `yaml:"enabled"`
	APIKey     string        `yaml:"api_key"`
	BaseURL    string        `yaml:"base_url"`
	Model      string        `yaml:"model"`
	Timeout    time.Duration `yaml:"timeout"`
	MaxRetries int           `yaml:"max_retries"`
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	Enabled   bool          `yaml:"enabled"`
	Directory string        `yaml:"directory"`
	TTL       time.Duration `yaml:"ttl"`
	MaxSize   string        `yaml:"max_size"`
}

// GenConfig holds generator configuration
type GenConfig struct {
	OutputLanguage string `yaml:"output_language"`
	GoVersion      string `yaml:"go_version"`
	CodeStyle      string `yaml:"code_style"`
}
