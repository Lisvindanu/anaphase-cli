# AI Provider Configuration

::: info AI is OPTIONAL in v0.4.0
**Template Mode is now the default!** Anaphase works perfectly without any AI provider. AI-powered generation is an optional enhancement for custom domain modeling. You can start building immediately using built-in templates.
:::

Configure AI providers for AI-powered domain generation.

## Template Mode (Default)

::: info New in v0.4.0
**No AI provider needed!** Anaphase now includes Template Mode that generates production-ready code using built-in templates. This is perfect for:
- Getting started quickly
- Standard CRUD operations
- Learning DDD patterns
- Building MVPs
:::

Simply run commands without setting up any API keys:

```bash
anaphase gen domain --name customer --template
anaphase gen handler --domain customer
anaphase gen repository --domain customer --db postgres
```

Template Mode generates clean, idiomatic Go code following DDD best practices.

## Supported AI Providers

Anaphase supports **4 major AI providers** for custom domain generation:

1. **Google Gemini** - Recommended for beginners (free tier)
2. **OpenAI** - GPT-4 and GPT-3.5 models
3. **Claude** - Anthropic's Claude models
4. **Groq** - Blazing fast inference (free preview)

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
- `gemini-2.0-flash-exp` (recommended) - Latest, fast, accurate
- `gemini-1.5-pro` - More capable, longer context
- `gemini-1.5-flash` - Fast and efficient

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
- `llama-3.3-70b-versatile` (recommended) - Latest, fast, versatile
- `llama-3.1-70b-versatile` - Stable alternative
- `mixtral-8x7b-32768` - Long context window

### OpenAI

Industry-leading models with excellent code generation.

**Get API Key:**
1. Visit [OpenAI Platform](https://platform.openai.com/api-keys)
2. Sign up or log in
3. Create API key
4. Copy your API key

**Models:**
- `gpt-4o` (recommended) - Latest, most capable
- `gpt-4-turbo` - Fast and efficient
- `gpt-3.5-turbo` - Budget-friendly option

**Pricing:** Pay-per-use (check [OpenAI Pricing](https://openai.com/pricing))

### Claude (Anthropic)

Advanced reasoning and code generation capabilities.

**Get API Key:**
1. Visit [Anthropic Console](https://console.anthropic.com/)
2. Sign up for account
3. Create API key
4. Copy your API key

**Models:**
- `claude-opus-4-5` (recommended) - Most capable
- `claude-sonnet-4-5` - Balanced performance
- `claude-3-5-sonnet` - Fast and efficient

**Pricing:** Pay-per-use (check [Anthropic Pricing](https://www.anthropic.com/pricing))

## Configuration Methods

::: info New in v0.4.0: CLI Configuration Commands
Use `anaphase config` commands to manage AI providers interactively:

```bash
# Set provider API keys
anaphase config set gemini.api_key "your-key"
anaphase config set openai.api_key "your-key"
anaphase config set claude.api_key "your-key"
anaphase config set groq.api_key "your-key"

# View current configuration
anaphase config show

# List available providers
anaphase config list-providers
```
:::

### Environment Variable (Simplest)

**Gemini:**
```bash
export GEMINI_API_KEY="your-gemini-api-key"
```

**OpenAI:**
```bash
export OPENAI_API_KEY="your-openai-api-key"
```

**Claude:**
```bash
export CLAUDE_API_KEY="your-claude-api-key"
```

**Groq:**
```bash
export GROQ_API_KEY="your-groq-api-key"
```

Add to shell profile for persistence:
```bash
# ~/.bashrc or ~/.zshrc
export GEMINI_API_KEY="your-gemini-api-key"
export OPENAI_API_KEY="your-openai-api-key"
export CLAUDE_API_KEY="your-claude-api-key"
export GROQ_API_KEY="your-groq-api-key"
```

### Configuration File (Recommended)

Anaphase auto-creates `~/.anaphase/config.yaml` on first run. Edit to customize:

```yaml
ai:
  # Which provider to use first
  primary_provider: gemini  # or: openai, claude, groq

  # Fallback providers (tried in order if primary fails)
  fallback_providers:
    - groq
    - openai
    - claude

  providers:
    # Google Gemini
    gemini:
      enabled: true
      api_key: ${GEMINI_API_KEY}
      model: gemini-2.0-flash-exp
      timeout: 30s
      max_retries: 3

    # OpenAI
    openai:
      enabled: true
      api_key: ${OPENAI_API_KEY}
      model: gpt-4o
      timeout: 30s
      max_retries: 3

    # Claude (Anthropic)
    claude:
      enabled: true
      api_key: ${CLAUDE_API_KEY}
      model: claude-opus-4-5
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
# 1. Set API keys (choose one or more)
export GEMINI_API_KEY="your-key"
export OPENAI_API_KEY="your-key"
export CLAUDE_API_KEY="your-key"
export GROQ_API_KEY="your-key"

# 2. Or use CLI commands
anaphase config set gemini.api_key "your-key"
anaphase config set openai.api_key "your-key"

# 3. View configuration
anaphase config show
```

## Advanced Configuration

### Multiple Providers (Fallback)

Anaphase automatically falls back to alternative providers if the primary fails:

```yaml
ai:
  # Try Gemini first
  primary_provider: gemini

  # If Gemini fails, try Groq, then OpenAI, then Claude
  fallback_providers:
    - groq
    - openai
    - claude

  providers:
    gemini:
      enabled: true
      api_key: ${GEMINI_API_KEY}
      model: gemini-2.0-flash-exp

    groq:
      enabled: true
      api_key: ${GROQ_API_KEY}
      model: llama-3.3-70b-versatile

    openai:
      enabled: true
      api_key: ${OPENAI_API_KEY}
      model: gpt-4o

    claude:
      enabled: true
      api_key: ${CLAUDE_API_KEY}
      model: claude-opus-4-5
```

**When fallback activates:**
- Primary quota exceeded
- Primary timeout
- Primary network error
- Primary API key invalid

**Example fallback scenario:**
1. Try Gemini (primary)
2. Gemini fails with quota exceeded → Try Groq
3. Groq succeeds
4. Generation continues with Groq

**Best practice:** Set up multiple providers for maximum reliability! All 4 providers work seamlessly together.

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
| `type` | string | - | Provider type (`gemini`, `openai`, `claude`, `groq`) |
| `apiKey` | string | - | API key |
| `model` | string | varies | Model to use (see provider-specific defaults) |
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
export OPENAI_API_KEY="your-openai-key"
export CLAUDE_API_KEY="your-claude-key"
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

### No AI Provider (Not an Error!)

::: info Template Mode Available
If you see messages about missing AI providers, don't worry! Anaphase v0.4.0 includes Template Mode that works without any AI provider. Simply use the `--template` flag or skip AI generation entirely.

```bash
# Use Template Mode instead
anaphase gen domain --name customer --template

# Or use the interactive menu
anaphase menu
```
:::

### API Key Not Found

```
Error: no AI providers configured - please set at least one API key
```

**Solution:**
Set at least one provider's API key:

```bash
# Option 1: Use Gemini (free tier)
export GEMINI_API_KEY="your-gemini-key"

# Option 2: Use Groq (fastest, free preview)
export GROQ_API_KEY="your-groq-key"

# Option 3: Use OpenAI
export OPENAI_API_KEY="your-openai-key"

# Option 4: Use Claude
export CLAUDE_API_KEY="your-claude-key"

# Best: Use multiple providers for fallback
export GEMINI_API_KEY="your-gemini-key"
export GROQ_API_KEY="your-groq-key"

# Or use Template Mode (no API key needed!)
anaphase gen domain --name customer --template
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
