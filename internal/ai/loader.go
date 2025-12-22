package ai

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// LoadConfig loads AI configuration from file and environment
func LoadConfig() (*Config, error) {
	// Set config file path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".anaphase")
	configFile := filepath.Join(configDir, "config.yaml")

	// Create default config if not exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := createDefaultConfig(configDir, configFile); err != nil {
			return nil, fmt.Errorf("create default config: %w", err)
		}
	}

	// Load config
	v := viper.New()
	v.SetConfigFile(configFile)
	v.SetConfigType("yaml")

	// Read environment variables
	v.AutomaticEnv()
	v.SetEnvPrefix("ANAPHASE")

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	// Unmarshal config
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// Override with environment variables
	overrideWithEnv(&config)

	return &config, nil
}

func createDefaultConfig(configDir, configFile string) error {
	// Create config directory
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	// Default configuration
	defaultConfig := `# Anaphase CLI Configuration
version: "1.0"

# AI Provider Configuration
ai:
  # Primary provider (will be tried first)
  primary_provider: gemini

  # Fallback chain (tried in order if primary fails)
  fallback_providers:
    - groq
    - openai

  # Provider-specific settings
  providers:
    gemini:
      enabled: true
      api_key: ${GEMINI_API_KEY}
      base_url: https://generativelanguage.googleapis.com
      model: gemini-2.0-flash-exp
      timeout: 30s
      max_retries: 3

    groq:
      enabled: false
      api_key: ${GROQ_API_KEY}
      base_url: https://api.groq.com/openai/v1
      model: llama-3.3-70b-versatile
      timeout: 30s
      max_retries: 3

    openai:
      enabled: false
      api_key: ${OPENAI_API_KEY}
      base_url: https://api.openai.com/v1
      model: gpt-4o-mini
      timeout: 30s
      max_retries: 3

    claude:
      enabled: false
      api_key: ${CLAUDE_API_KEY}
      base_url: https://api.anthropic.com/v1
      model: claude-3-5-sonnet-20241022
      timeout: 45s
      max_retries: 2

    ollama:
      enabled: false
      base_url: http://localhost:11434
      model: qwen2.5-coder:7b
      timeout: 60s
      max_retries: 1

# Cache Configuration
cache:
  enabled: true
  directory: ~/.anaphase/cache
  ttl: 24h
  max_size: 100MB

# Generator Settings
generator:
  output_language: go
  go_version: "1.22"
  code_style: gofmt
`

	// Write config file
	if err := os.WriteFile(configFile, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	return nil
}

func overrideWithEnv(config *Config) {
	// Gemini
	if key := os.Getenv("GEMINI_API_KEY"); key != "" {
		config.AI.Providers.Gemini.APIKey = key
		config.AI.Providers.Gemini.Enabled = true
	}

	// Groq
	if key := os.Getenv("GROQ_API_KEY"); key != "" {
		config.AI.Providers.Groq.APIKey = key
		config.AI.Providers.Groq.Enabled = true
	}

	// OpenAI
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		config.AI.Providers.OpenAI.APIKey = key
		config.AI.Providers.OpenAI.Enabled = true
	}

	// Claude
	if key := os.Getenv("CLAUDE_API_KEY"); key != "" {
		config.AI.Providers.Claude.APIKey = key
		config.AI.Providers.Claude.Enabled = true
	}

	// Parse durations (viper doesn't auto-parse to time.Duration from env)
	if config.AI.Providers.Gemini.Timeout == 0 {
		config.AI.Providers.Gemini.Timeout = 30 * time.Second
	}
	if config.AI.Providers.Groq.Timeout == 0 {
		config.AI.Providers.Groq.Timeout = 30 * time.Second
	}
	if config.AI.Providers.OpenAI.Timeout == 0 {
		config.AI.Providers.OpenAI.Timeout = 30 * time.Second
	}
	if config.AI.Providers.Claude.Timeout == 0 {
		config.AI.Providers.Claude.Timeout = 45 * time.Second
	}
	if config.AI.Providers.Ollama.Timeout == 0 {
		config.AI.Providers.Ollama.Timeout = 60 * time.Second
	}

	// Set default primary provider if not set
	if config.AI.PrimaryProvider == "" {
		config.AI.PrimaryProvider = "gemini"
	}

	// Set default fallback chain if empty
	if len(config.AI.FallbackProviders) == 0 {
		config.AI.FallbackProviders = []string{"groq", "openai"}
	}

	// Parse cache TTL
	if config.Cache.TTL == 0 {
		config.Cache.TTL = 24 * time.Hour
	}

	// Expand home directory in cache path
	if config.Cache.Directory != "" {
		if config.Cache.Directory[:2] == "~/" {
			homeDir, _ := os.UserHomeDir()
			config.Cache.Directory = filepath.Join(homeDir, config.Cache.Directory[2:])
		}
	}
}
