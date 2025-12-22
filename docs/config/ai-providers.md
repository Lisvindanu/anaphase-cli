# AI Provider Configuration

Configure AI providers for domain generation.

## Supported Providers

### Google Gemini (Recommended for Beginners)

Free tier with generous limits and excellent code generation quality.

**Get API Key:**
1. Visit [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Sign in with Google account
3. Click "Create API Key"
4. Copy your API key

**Free Tier Limits:**
- 60 requests per minute
- Sufficient for most development

**Models:**
- `gemini-2.0-flash-exp` (recommended) - Fast, accurate
- `gemini-pro` - More capable, slower

### Groq (Recommended for Speed)

**Extremely fast inference** with open-source models. Free during preview.

**Get API Key:**
1. Visit [Groq Console](https://console.groq.com/keys)
2. Sign up for free account
3. Create API key
4. Copy your API key

**Free Tier:**
- Very generous rate limits
- Faster than Gemini
- Free during preview period

**Models:**
- `llama-3.3-70b-versatile` (default) - Fast, versatile
- `mixtral-8x7b-32768` - Long context
- `llama-3.1-70b-versatile` - Alternative option

## Configuration Methods

### Environment Variable (Simplest)

**Gemini:**
```bash
export GEMINI_API_KEY="your-gemini-api-key"
```

**Groq:**
```bash
export GROQ_API_KEY="your-groq-api-key"
```

Add to shell profile for persistence:
```bash
# ~/.bashrc or ~/.zshrc
export GEMINI_API_KEY="your-gemini-api-key"
export GROQ_API_KEY="your-groq-api-key"
```

### Configuration File (Recommended)

Anaphase auto-creates `~/.anaphase/config.yaml` on first run. Edit to customize:

```yaml
ai:
  # Which provider to use first
  primary_provider: gemini  # or: groq

  # Fallback providers (tried in order if primary fails)
  fallback_providers:
    - groq
    - openai

  providers:
    # Google Gemini
    gemini:
      enabled: true
      api_key: ${GEMINI_API_KEY}
      model: gemini-2.0-flash-exp
      timeout: 30s
      max_retries: 3

    # Groq (Fast!)
    groq:
      enabled: true
      api_key: ${GROQ_API_KEY}
      model: llama-3.3-70b-versatile
      timeout: 30s
      max_retries: 3
```

**Quick Setup:**
```bash
# 1. Set API keys
export GEMINI_API_KEY="your-key"
export GROQ_API_KEY="your-key"

# 2. Run any command to generate config
anaphase gen domain "test"

# 3. Edit config to choose provider
vim ~/.anaphase/config.yaml
# Change: primary_provider: groq
```

## Advanced Configuration

### Multiple Providers (Fallback)

Anaphase automatically falls back to alternative providers if the primary fails:

```yaml
ai:
  # Try Gemini first
  primary_provider: gemini

  # If Gemini fails, try Groq, then OpenAI
  fallback_providers:
    - groq
    - openai

  providers:
    gemini:
      enabled: true
      api_key: ${GEMINI_API_KEY}
      model: gemini-2.0-flash-exp

    groq:
      enabled: true
      api_key: ${GROQ_API_KEY}
      model: llama-3.3-70b-versatile
```

**When fallback activates:**
- Primary quota exceeded
- Primary timeout
- Primary network error
- Primary API key invalid

**Example fallback scenario:**
1. Try Gemini (primary)
2. Gemini fails with quota exceeded → Try Groq
3. Groq succeeds ✅
4. Generation continues with Groq

**Best practice:** Set up both Gemini and Groq for maximum reliability!

### Caching

Enable response caching to save API calls:

```yaml
cache:
  enabled: true
  ttl: 24h
  dir: ~/.anaphase/cache
```

**Benefits:**
- Faster regeneration
- Save API quota
- Work offline (if cached)

**Cache invalidation:**
```bash
# Clear all cache
rm -rf ~/.anaphase/cache

# Clear specific domain
rm -rf ~/.anaphase/cache/customer*
```

### Request Tuning

Fine-tune AI behavior:

```yaml
ai:
  primary:
    type: gemini
    apiKey: YOUR_API_KEY
    model: gemini-2.5-flash

    # Request timeout (default: 30s)
    timeout: 60s

    # Retry attempts (default: 3)
    retries: 5

    # AI creativity (default: 0.3)
    # Lower = more consistent
    # Higher = more creative
    temperature: 0.1

    # Max tokens in response
    maxTokens: 4096
```

## Configuration Reference

### AI Provider Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `type` | string | - | Provider type (`gemini`) |
| `apiKey` | string | - | API key |
| `model` | string | `gemini-2.5-flash` | Model to use |
| `timeout` | duration | `30s` | Request timeout |
| `retries` | int | `3` | Retry attempts |
| `temperature` | float | `0.3` | AI creativity (0.0-1.0) |
| `maxTokens` | int | `4096` | Max response tokens |

### Cache Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `enabled` | bool | `true` | Enable caching |
| `ttl` | duration | `24h` | Cache lifetime |
| `dir` | string | `~/.anaphase/cache` | Cache directory |

## Environment Variable Overrides

Environment variables take precedence over config file:

```bash
# Set provider API keys (auto-enables the provider)
export GEMINI_API_KEY="your-gemini-key"
export GROQ_API_KEY="your-groq-key"

# Override config file path
export ANAPHASE_CONFIG="/custom/path/config.yaml"

# Disable cache
export ANAPHASE_CACHE_ENABLED=false
```

**Note:** Setting an API key via environment variable automatically enables that provider!

## Verification

Test your configuration:

```bash
# Generate a test domain with verbose logging
anaphase gen domain "Test entity with name field" --verbose

# You should see output like:
# time=... level=INFO msg="attempting generation" provider=gemini
# time=... level=INFO msg="generation successful" provider=gemini tokens=1234 cost=$0.000000 duration=2.5s

# Or if using Groq:
# time=... level=INFO msg="attempting generation" provider=groq
# time=... level=INFO msg="generation successful" provider=groq tokens=1915 cost=$0.000000 duration=1.6s
```

**Check which provider is being used:**
```bash
cat ~/.anaphase/config.yaml | grep primary_provider
```

**Test fallback:**
```bash
# Remove primary provider API key temporarily
unset GEMINI_API_KEY

# Generate domain - should fallback to Groq
anaphase gen domain "Test" --verbose
# You'll see:
# level=WARN msg="provider not available" provider=gemini
# level=INFO msg="attempting generation" provider=groq
```

## Troubleshooting

### API Key Not Found

```
Error: no AI providers configured - please set at least one API key
```

**Solution:**
Set at least one provider's API key:

```bash
# Option 1: Use Gemini
export GEMINI_API_KEY="your-gemini-key"

# Option 2: Use Groq (faster!)
export GROQ_API_KEY="your-groq-key"

# Best: Use both for fallback
export GEMINI_API_KEY="your-gemini-key"
export GROQ_API_KEY="your-groq-key"

# Verify it works
anaphase gen domain "Test" --verbose
```

### Quota Exceeded

```
Error: quota exceeded (429)
```

**Solutions:**
1. Wait (quota resets per minute)
2. Use fallback provider
3. Enable caching to reduce calls
4. Upgrade to paid tier

### Invalid Response

```
Error: failed to parse AI response
```

**Solutions:**
1. Try again (AI can be inconsistent)
2. Lower temperature to 0.1
3. Check model (use gemini-2.5-flash)
4. Enable verbose logging

### Timeout

```
Error: request timeout
```

**Solutions:**
1. Increase timeout in config:
   ```yaml
   timeout: 60s
   ```
2. Check network connection
3. Try different model

## Best Practices

### Security

- **Don't commit API keys** to git
- Use environment variables in CI/CD
- Rotate keys periodically
- Use different keys per environment

```bash
# Development
export GEMINI_API_KEY="dev-key"

# Production
export GEMINI_API_KEY="prod-key"
```

### Performance

- Enable caching
- Use gemini-2.5-flash (fastest)
- Set reasonable timeout
- Use lower temperature

### Reliability

- Configure fallback provider
- Enable retries
- Monitor quota usage
- Keep cache enabled

## See Also

- [AI-Powered Generation](/guide/ai-generation)
- [Installation](/guide/installation)
- [gen domain](/reference/gen-domain)
