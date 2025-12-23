# anaphase config

Kelola konfigurasi Anaphase termasuk AI provider, cache settings, dan lainnya.

::: info
**Akses Cepat**: Jalankan `anaphase` (tanpa argumen) untuk mengakses menu interaktif dan pilih "Configuration" untuk interface visual mengelola setting.
:::

## Overview

Command `config` menyediakan tools untuk melihat dan mengelola konfigurasi CLI Anaphase Anda. Command ini mendukung multiple AI provider, fallback chain, dan health monitoring.

::: info
**AI Bersifat Opsional**: Anaphase bekerja tanpa konfigurasi AI menggunakan template mode. AI provider hanya diperlukan jika Anda ingin code generation berbasis AI.
:::

## Subcommand

### list

Tampilkan konfigurasi Anaphase saat ini.

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

Set AI provider default untuk code generation.

```bash
anaphase config set-provider <provider>
```

**Provider yang Tersedia:**
- `gemini` - Google Gemini (tier gratis tersedia)
- `groq` - Groq (inference tercepat, gratis)
- `openai` - OpenAI GPT models (berbayar)
- `claude` - Anthropic Claude (berbayar)

**Contoh:**
```bash
anaphase config set-provider groq
```

**Output:**
```
✓ Default provider set to: groq
ℹ Note: Config changes are temporary. To persist, edit ~/.anaphase/config.yaml
```

### check

Health check semua AI provider yang dikonfigurasi.

```bash
anaphase config check
```

**Output:**
```
⚡ Provider Health Check

  ✓ gemini - Healthy (primary)
  ✓ groq - Healthy
  ✗ openai - API key not configured
  ✗ claude - API key not configured

ℹ 2 of 4 providers configured and healthy
ℹ Template mode is always available (no AI required)
```

**Use Case:**
- Verifikasi API key berfungsi
- Cek konektivitas network
- Diagnosa masalah provider
- Validasi konfigurasi sebelum generation

### show-providers

List semua AI provider yang tersedia dengan detail.

```bash
anaphase config show-providers
```

**Output:**
```
⚡ Available AI Providers

📡 Gemini (Recommended)
  Google's AI model - Reliable, free tier available
  💰 Free: Yes (Generous free tier)
  🎯 Models: [gemini-2.0-flash-exp, gemini-2.5-flash]
  📝 API Key: GEMINI_API_KEY

📡 Groq (Fastest)
  Extremely fast inference - Free during preview
  💰 Free: Yes (Limited preview)
  🎯 Models: [llama-3.3-70b-versatile, mixtral-8x7b-32768]
  📝 API Key: GROQ_API_KEY

📡 OpenAI
  GPT models - Paid service, high quality
  💰 Free: No (Paid)
  🎯 Models: [gpt-4o, gpt-4o-mini]
  📝 API Key: OPENAI_API_KEY

📡 Claude
  Anthropic's AI - Paid service, excellent for complex tasks
  💰 Free: No (Paid)
  🎯 Models: [claude-3-5-sonnet-20241022]
  📝 API Key: ANTHROPIC_API_KEY

ℹ Template mode is always available without any AI provider
```

## File Konfigurasi

Anaphase membaca konfigurasi dari `~/.anaphase/config.yaml`:

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

    claude:
      enabled: false
      api_key: ${ANTHROPIC_API_KEY}
      model: claude-3-5-sonnet-20241022
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

## Environment Variable

API key dapat diset via environment variable:

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

**Auto-Enable:** Ketika API key diset via environment variable, provider secara otomatis diaktifkan.

## Pemilihan Provider

### Flag Command-Line

Override provider default untuk single command:

```bash
anaphase gen domain "User with email" --provider groq
```

### Mode Interaktif

Pilih provider secara interaktif:

```bash
anaphase gen domain --interactive
```

Prompt:
```
Select AI provider:
  1) gemini (default)
  2) groq
  3) openai
  4) claude
Enter choice [1]:
```

### File Konfigurasi

Set default di `~/.anaphase/config.yaml`:

```yaml
ai:
  primary_provider: groq
```

### Override Sementara

Set untuk sesi saat ini:

```bash
anaphase config set-provider groq
```

## Fallback Chain

Anaphase secara otomatis mencoba fallback provider jika primary gagal:

```yaml
ai:
  primary_provider: gemini
  fallback_providers:
    - groq
    - openai
```

**Perilaku:**
1. Coba `gemini` terlebih dahulu
2. Jika gagal (rate limit, error, timeout), coba `groq`
3. Jika masih gagal, coba `openai`
4. Jika semua gagal, return error

**Use Case:**
- **Reliability:** Automatic failover
- **Rate Limit:** Switch ketika quota terlampaui
- **Optimisasi Biaya:** Gunakan tier gratis dahulu, berbayar sebagai backup

## Perbandingan Provider

| Provider | Kecepatan | Biaya | Tier Gratis | Terbaik Untuk |
|----------|-------|------|-----------|----------|
| **Template Mode** | ⚡⚡⚡⚡⚡ | Gratis | ✅ Selalu | Quick generation, no setup |
| **Gemini** | ⚡⚡⚡⚡ | Gratis | ✅ Generous | Penggunaan umum, kualitas tinggi |
| **Groq** | ⚡⚡⚡⚡⚡ | Gratis | ✅ Limited | Kecepatan kritis, real-time |
| **OpenAI** | ⚡⚡⚡ | $$$ | ❌ Berbayar | Domain kompleks, akurasi |
| **Claude** | ⚡⚡⚡ | $$$ | ❌ Berbayar | Konteks besar, analisis |

## Contoh

### Setup Gemini (Gratis)

1. Dapatkan API key dari [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Set environment variable:
   ```bash
   export GEMINI_API_KEY="your-key-here"
   ```
3. Verifikasi:
   ```bash
   anaphase config check
   ```

### Setup Groq (Cepat & Gratis)

1. Dapatkan API key dari [Groq Console](https://console.groq.com/)
2. Set environment variable:
   ```bash
   export GROQ_API_KEY="your-key-here"
   ```
3. Set sebagai default:
   ```bash
   anaphase config set-provider groq
   ```

### Setup Multi-Provider

```bash
# Set semua API key
export GEMINI_API_KEY="gemini-key"
export GROQ_API_KEY="groq-key"
export OPENAI_API_KEY="openai-key"

# Cek semua provider
anaphase config check

# Gunakan Groq sebagai primary, Gemini sebagai fallback
anaphase config set-provider groq
```

Edit `~/.anaphase/config.yaml`:
```yaml
ai:
  primary_provider: groq
  fallback_providers:
    - gemini
```

### Tidak Perlu Setup AI

Gunakan template mode tanpa konfigurasi apapun:

```bash
# Bekerja langsung tanpa API key
anaphase gen domain "User with email and password"

# Template mode cepat dan andal
anaphase gen handler --domain user
anaphase gen repository --domain user --db postgres
```

::: tip
Template mode generate kode production-ready berdasarkan pattern yang sudah terbukti. Anda tidak perlu AI untuk kebanyakan use case.
:::

## Troubleshooting

### Provider Tidak Bekerja

```bash
# Cek konfigurasi
anaphase config list

# Test kesehatan provider
anaphase config check

# Verifikasi API key sudah diset
echo $GEMINI_API_KEY
```

### Error Rate Limit

```bash
# Ganti ke provider berbeda
anaphase config set-provider groq

# Atau gunakan fallback chain
# Edit ~/.anaphase/config.yaml
```

### API Key Tidak Valid

```bash
# Re-export dengan key yang benar
export GEMINI_API_KEY="correct-key-here"

# Verifikasi
anaphase config check
```

### Generasi Lambat

```bash
# Ganti ke Groq (tercepat)
anaphase config set-provider groq

# Atau gunakan model lebih cepat
# Edit config.yaml:
# model: gemini-2.0-flash-exp
```

## Best Practice

1. **Gunakan Tier Gratis Terlebih Dahulu**
   - Mulai dengan Gemini atau Groq
   - Hanya gunakan provider berbayar ketika diperlukan

2. **Setup Fallback**
   - Konfigurasi multiple provider
   - Automatic failover untuk reliability

3. **Amankan API Key**
   - Gunakan environment variable
   - Jangan commit key ke git
   - Rotasi key secara berkala

4. **Monitor Penggunaan**
   - Cek dashboard provider
   - Perhatikan rate limit
   - Track biaya (untuk provider berbayar)

5. **Test Sebelum Production**
   ```bash
   anaphase config check
   anaphase gen domain "Test" --provider gemini
   ```

## Lihat Juga

- [AI Providers Configuration](/config/ai-providers) - Setup provider detail
- [anaphase gen domain](/reference/gen-domain) - AI-powered domain generation
- [Troubleshooting](/guide/troubleshooting) - Masalah umum
