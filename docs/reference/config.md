# anaphase config

Manage Anaphase configuration including AI providers, cache settings, and more.

## Overview

The `config` command provides tools to view and manage your Anaphase CLI configuration. It supports multiple AI providers, fallback chains, and health monitoring.

## Subcommands

### list

Display current Anaphase configuration.

```bash
anaphase config list
```

**Output:**
```
⚡ Anaphase Configuration

🤖 AI Configuration:
  Primary Provider: gemini
  Fallback Providers: [groq openai]

📡 Configured Providers:
  ✓ Gemini
    Model: gemini-2.5-flash
    Timeout: 30s
    Max Retries: 3

  ✓ Groq
    Model: llama-3.3-70b-versatile
    Timeout: 30s
    Max Retries: 3

💾 Cache Configuration:
  Enabled: true
  Directory: ~/.anaphase/cache
  TTL: 24h

⚙️ Generator Settings:
  Go Version: 1.21
  Code Style: standard
```

### set-provider

Set the default AI provider for code generation.

```bash
anaphase config set-provider <provider>
```

**Available Providers:**
- `gemini` - Google Gemini
- `groq` - Groq (fastest inference)
- `openai` - OpenAI GPT models
- `claude` - Anthropic Claude
- `ollama` - Local AI models

**Example:**
```bash
anaphase config set-provider groq
```

**Output:**
```
✓ Default provider set to: groq
ℹ Note: Config changes are temporary. To persist, edit ~/.anaphase/config.yaml
```

### check

Health check all configured AI providers.

```bash
anaphase config check
```

**Output:**
```
⚡ Provider Health Check

  ✓ gemini - Healthy
  ✓ groq - Healthy
  ✗ openai - API key not configured
  ✗ claude - API key not configured
```

**Use Cases:**
- Verify API keys are working
- Check network connectivity
- Diagnose provider issues
- Validate configuration before generation

### show-providers

List all available AI providers with details.

```bash
anaphase config show-providers
```

**Output:**
```
⚡ Available AI Providers

📡 Gemini
  Google's AI model - Reliable, free tier available
  💰 Free: Yes
  🎯 Models: [gemini-2.0-flash-exp, gemini-pro]

📡 Groq
  Extremely fast inference - Free during preview
  💰 Free: Yes
  🎯 Models: [llama-3.3-70b-versatile, mixtral-8x7b-32768]

📡 OpenAI
  GPT models - Paid service
  💰 Free: No (Paid)
  🎯 Models: [gpt-4o, gpt-4o-mini]

📡 Claude
  Anthropic's AI - Paid service
  💰 Free: No (Paid)
  🎯 Models: [claude-3-5-sonnet-20241022]

📡 Ollama
  Local AI - Run models on your machine
  💰 Free: Yes
  🎯 Models: [qwen2.5-coder:7b, codellama]
```

## Configuration File

Anaphase reads configuration from `~/.anaphase/config.yaml`:

```yaml
ai:
  primary_provider: gemini
  fallback_providers:
    - groq
    - openai

  providers:
    gemini:
      enabled: true
      api_key: ${GEMINI_API_KEY}
      model: gemini-2.5-flash
      timeout: 30s
      max_retries: 3

    groq:
      enabled: true
      api_key: ${GROQ_API_KEY}
      model: llama-3.3-70b-versatile
      timeout: 30s
      max_retries: 3

    openai:
      enabled: false
      api_key: ${OPENAI_API_KEY}
      model: gpt-4o
      timeout: 60s
      max_retries: 3

cache:
  enabled: true
  directory: ~/.anaphase/cache
  ttl: 24h

generator:
  go_version: "1.21"
  code_style: standard
```

## Environment Variables

API keys can be set via environment variables:

```bash
# Gemini (Google AI)
export GEMINI_API_KEY="your-gemini-api-key"

# Groq
export GROQ_API_KEY="your-groq-api-key"

# OpenAI
export OPENAI_API_KEY="your-openai-api-key"

# Anthropic Claude
export ANTHROPIC_API_KEY="your-anthropic-api-key"
```

**Auto-Enable:** When an API key is set via environment variable, the provider is automatically enabled.

## Provider Selection

### Command-Line Flag

Override the default provider for a single command:

```bash
anaphase gen domain "User with email" --provider groq
```

### Interactive Mode

Select provider interactively:

```bash
anaphase gen domain --interactive
```

Prompts:
```
Select AI provider:
  1) gemini (default)
  2) groq
  3) openai
  4) claude
Enter choice [1]:
```

### Configuration File

Set default in `~/.anaphase/config.yaml`:

```yaml
ai:
  primary_provider: groq
```

### Temporary Override

Set for current session:

```bash
anaphase config set-provider groq
```

## Fallback Chain

Anaphase automatically tries fallback providers if the primary fails:

```yaml
ai:
  primary_provider: gemini
  fallback_providers:
    - groq
    - openai
```

**Behavior:**
1. Try `gemini` first
2. If fails (rate limit, error, timeout), try `groq`
3. If still fails, try `openai`
4. If all fail, return error

**Use Cases:**
- **Reliability:** Automatic failover
- **Rate Limits:** Switch when quota exceeded
- **Cost Optimization:** Use free tier first, paid as backup

## Provider Comparison

| Provider | Speed | Cost | Free Tier | Best For |
|----------|-------|------|-----------|----------|
| **Gemini** | ⚡⚡⚡ | Free | ✅ Generous | General use, high quality |
| **Groq** | ⚡⚡⚡⚡⚡ | Free | ✅ Limited | Speed-critical, real-time |
| **OpenAI** | ⚡⚡⚡ | $$$ | ❌ Paid | Complex domains, accuracy |
| **Claude** | ⚡⚡⚡ | $$$ | ❌ Paid | Large contexts, analysis |
| **Ollama** | ⚡⚡ | Free | ✅ Local | Offline, privacy |

## Examples

### Setup Gemini (Free)

1. Get API key from [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Set environment variable:
   ```bash
   export GEMINI_API_KEY="your-key-here"
   ```
3. Verify:
   ```bash
   anaphase config check
   ```

### Setup Groq (Fast & Free)

1. Get API key from [Groq Console](https://console.groq.com/)
2. Set environment variable:
   ```bash
   export GROQ_API_KEY="your-key-here"
   ```
3. Set as default:
   ```bash
   anaphase config set-provider groq
   ```

### Multi-Provider Setup

```bash
# Set all API keys
export GEMINI_API_KEY="gemini-key"
export GROQ_API_KEY="groq-key"
export OPENAI_API_KEY="openai-key"

# Check all providers
anaphase config check

# Use Groq as primary, Gemini as fallback
anaphase config set-provider groq
```

Edit `~/.anaphase/config.yaml`:
```yaml
ai:
  primary_provider: groq
  fallback_providers:
    - gemini
```

### Local Development with Ollama

1. Install Ollama:
   ```bash
   curl -fsSL https://ollama.com/install.sh | sh
   ```

2. Pull a model:
   ```bash
   ollama pull qwen2.5-coder:7b
   ```

3. Use with Anaphase:
   ```bash
   anaphase gen domain "User" --provider ollama
   ```

## Troubleshooting

### Provider Not Working

```bash
# Check configuration
anaphase config list

# Test provider health
anaphase config check

# Verify API key is set
echo $GEMINI_API_KEY
```

### Rate Limit Errors

```bash
# Switch to different provider
anaphase config set-provider groq

# Or use fallback chain
# Edit ~/.anaphase/config.yaml
```

### Invalid API Key

```bash
# Re-export with correct key
export GEMINI_API_KEY="correct-key-here"

# Verify
anaphase config check
```

### Slow Generation

```bash
# Switch to Groq (fastest)
anaphase config set-provider groq

# Or use faster model
# Edit config.yaml:
# model: gemini-2.0-flash-exp
```

## Best Practices

1. **Use Free Tiers First**
   - Start with Gemini or Groq
   - Only use paid providers when needed

2. **Set Up Fallbacks**
   - Configure multiple providers
   - Automatic failover for reliability

3. **Secure API Keys**
   - Use environment variables
   - Never commit keys to git
   - Rotate keys regularly

4. **Monitor Usage**
   - Check provider dashboards
   - Watch for rate limits
   - Track costs (for paid providers)

5. **Test Before Production**
   ```bash
   anaphase config check
   anaphase gen domain "Test" --provider gemini
   ```

## See Also

- [AI Providers Configuration](/config/ai-providers) - Detailed provider setup
- [anaphase gen domain](/reference/gen-domain) - AI-powered domain generation
- [Troubleshooting](/guide/troubleshooting) - Common issues
