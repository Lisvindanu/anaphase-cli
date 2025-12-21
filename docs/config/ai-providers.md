# AI Provider Configuration

Configure AI providers for domain generation.

## Supported Providers

### Google Gemini (Recommended)

Free tier with generous limits.

**Get API Key:**
1. Visit [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Sign in with Google account
3. Click "Create API Key"
4. Copy your API key

**Free Tier Limits:**
- 60 requests per minute
- Sufficient for most development

**Models:**
- `gemini-2.5-flash` (recommended) - Fast, accurate
- `gemini-pro` - More capable, slower

## Configuration Methods

### Environment Variable (Simplest)

```bash
export GEMINI_API_KEY="your-api-key-here"
```

Add to shell profile:
```bash
# ~/.bashrc or ~/.zshrc
export GEMINI_API_KEY="your-api-key-here"
```

### Configuration File (Recommended)

Create `~/.anaphase/config.yaml`:

```yaml
ai:
  primary:
    type: gemini
    apiKey: YOUR_API_KEY_HERE
    model: gemini-2.5-flash
    timeout: 30s
    retries: 3
    temperature: 0.3
```

## Advanced Configuration

### Multiple Providers (Fallback)

Set up automatic fallback:

```yaml
ai:
  # Primary provider
  primary:
    type: gemini
    apiKey: PRIMARY_API_KEY
    model: gemini-2.5-flash
    timeout: 30s
    retries: 3

  # Fallback if primary fails
  secondary:
    type: gemini
    apiKey: SECONDARY_API_KEY
    model: gemini-2.5-flash
    timeout: 30s
    retries: 2
```

**When fallback activates:**
- Primary quota exceeded
- Primary timeout
- Primary network error

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
# Override API key
export GEMINI_API_KEY="override-key"

# Override config file path
export ANAPHASE_CONFIG="/custom/path/config.yaml"

# Disable cache
export ANAPHASE_CACHE_ENABLED=false
```

## Verification

Test your configuration:

```bash
# Generate a test domain
anaphase gen domain \
  --name test \
  --prompt "Test entity with name field" \
  --verbose

# Should see:
# - Using provider: gemini
# - Model: gemini-2.5-flash
# - Cache: enabled/disabled
```

## Troubleshooting

### API Key Not Found

```
Error: GEMINI_API_KEY not set
```

**Solution:**
```bash
# Set environment variable
export GEMINI_API_KEY="your-key"

# Or create config file
mkdir -p ~/.anaphase
cat > ~/.anaphase/config.yaml << EOF
ai:
  primary:
    type: gemini
    apiKey: your-key
EOF
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
