# Konfigurasi AI Provider

::: info AI OPSIONAL di v0.4.0
**Template Mode sekarang menjadi default!** Anaphase bekerja sempurna tanpa AI provider apapun. Generasi bertenaga AI adalah peningkatan opsional untuk pemodelan domain kustom. Anda bisa langsung mulai membangun menggunakan template bawaan.
:::

Konfigurasi AI provider untuk generasi domain bertenaga AI.

## Template Mode (Default)

::: info Baru di v0.4.0
**Tidak perlu AI provider!** Anaphase sekarang menyertakan Template Mode yang menghasilkan kode production-ready menggunakan template bawaan. Ini sempurna untuk:
- Memulai dengan cepat
- Operasi CRUD standar
- Mempelajari pola DDD
- Membangun MVP
:::

Cukup jalankan perintah tanpa perlu setup API key:

```bash
anaphase gen domain --name customer --template
anaphase gen handler --domain customer
anaphase gen repository --domain customer --db postgres
```

Template Mode menghasilkan kode Go yang bersih dan idiomatis mengikuti best practice DDD.

## AI Provider yang Didukung

Anaphase mendukung **4 AI provider utama** untuk generasi domain kustom:

1. **Google Gemini** - Direkomendasikan untuk pemula (tier gratis)
2. **OpenAI** - Model GPT-4 dan GPT-3.5
3. **Claude** - Model Claude dari Anthropic
4. **Groq** - Inferensi super cepat (preview gratis)

### Google Gemini (Direkomendasikan untuk Pemula)

Tier gratis dengan limit yang besar dan kualitas generasi kode yang sangat baik.

**Dapatkan API Key:**
1. Kunjungi [Google AI Studio](https://makersuite.google.com/app/apikey)
2. Login dengan akun Google
3. Klik "Create API Key"
4. Salin API key Anda

**Limit Tier Gratis:**
- 60 request per menit
- Cukup untuk sebagian besar development

**Model:**
- `gemini-2.0-flash-exp` (direkomendasikan) - Terbaru, cepat, akurat
- `gemini-1.5-pro` - Lebih capable, konteks lebih panjang
- `gemini-1.5-flash` - Cepat dan efisien

### Groq (Direkomendasikan untuk Kecepatan)

**Inferensi sangat cepat** dengan model open-source. Gratis selama periode preview.

**Dapatkan API Key:**
1. Kunjungi [Groq Console](https://console.groq.com/keys)
2. Daftar akun gratis
3. Buat API key
4. Salin API key Anda

**Tier Gratis:**
- Rate limit yang sangat besar
- Lebih cepat dari Gemini
- Gratis selama periode preview

**Model:**
- `llama-3.3-70b-versatile` (direkomendasikan) - Terbaru, cepat, serbaguna
- `llama-3.1-70b-versatile` - Alternatif yang stabil
- `mixtral-8x7b-32768` - Context window panjang

### OpenAI

Model terdepan di industri dengan generasi kode yang sangat baik.

**Dapatkan API Key:**
1. Kunjungi [OpenAI Platform](https://platform.openai.com/api-keys)
2. Daftar atau login
3. Buat API key
4. Salin API key Anda

**Model:**
- `gpt-4o` (direkomendasikan) - Terbaru, paling capable
- `gpt-4-turbo` - Cepat dan efisien
- `gpt-3.5-turbo` - Opsi hemat budget

**Harga:** Pay-per-use (cek [OpenAI Pricing](https://openai.com/pricing))

### Claude (Anthropic)

Kemampuan reasoning dan generasi kode yang canggih.

**Dapatkan API Key:**
1. Kunjungi [Anthropic Console](https://console.anthropic.com/)
2. Daftar akun
3. Buat API key
4. Salin API key Anda

**Model:**
- `claude-opus-4-5` (direkomendasikan) - Paling capable
- `claude-sonnet-4-5` - Performa seimbang
- `claude-3-5-sonnet` - Cepat dan efisien

**Harga:** Pay-per-use (cek [Anthropic Pricing](https://www.anthropic.com/pricing))

## Metode Konfigurasi

::: info Baru di v0.4.0: Perintah Konfigurasi CLI
Gunakan perintah `anaphase config` untuk mengelola AI provider secara interaktif:

```bash
# Set API key provider
anaphase config set gemini.api_key "your-key"
anaphase config set openai.api_key "your-key"
anaphase config set claude.api_key "your-key"
anaphase config set groq.api_key "your-key"

# Lihat konfigurasi saat ini
anaphase config show

# Daftar provider yang tersedia
anaphase config list-providers
```
:::

### Environment Variable (Paling Sederhana)

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

Tambahkan ke shell profile agar persisten:
```bash
# ~/.bashrc atau ~/.zshrc
export GEMINI_API_KEY="your-gemini-api-key"
export OPENAI_API_KEY="your-openai-api-key"
export CLAUDE_API_KEY="your-claude-api-key"
export GROQ_API_KEY="your-groq-api-key"
```

### File Konfigurasi (Direkomendasikan)

Anaphase secara otomatis membuat `~/.anaphase/config.yaml` pada first run. Edit untuk kustomisasi:

```yaml
ai:
  # Provider mana yang digunakan pertama
  primary_provider: gemini  # atau: openai, claude, groq

  # Provider fallback (dicoba berurutan jika primary gagal)
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

    # Groq (Cepat!)
    groq:
      enabled: true
      api_key: ${GROQ_API_KEY}
      model: llama-3.3-70b-versatile
      timeout: 30s
      max_retries: 3
```

**Setup Cepat:**
```bash
# 1. Set API key (pilih satu atau lebih)
export GEMINI_API_KEY="your-key"
export OPENAI_API_KEY="your-key"
export CLAUDE_API_KEY="your-key"
export GROQ_API_KEY="your-key"

# 2. Atau gunakan perintah CLI
anaphase config set gemini.api_key "your-key"
anaphase config set openai.api_key "your-key"

# 3. Lihat konfigurasi
anaphase config show
```

## Konfigurasi Advanced

### Multiple Provider (Fallback)

Anaphase otomatis fallback ke provider alternatif jika primary gagal:

```yaml
ai:
  # Coba Gemini dulu
  primary_provider: gemini

  # Jika Gemini gagal, coba Groq, lalu OpenAI, lalu Claude
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

**Kapan fallback aktif:**
- Quota primary terlampaui
- Primary timeout
- Primary network error
- API key primary tidak valid

**Contoh skenario fallback:**
1. Coba Gemini (primary)
2. Gemini gagal dengan quota exceeded → Coba Groq
3. Groq berhasil
4. Generasi berlanjut dengan Groq

**Best practice:** Setup multiple provider untuk reliabilitas maksimal! Semua 4 provider bekerja mulus bersama.

### Caching

Aktifkan response caching untuk menghemat API call:

```yaml
cache:
  enabled: true
  ttl: 24h
  dir: ~/.anaphase/cache
```

**Keuntungan:**
- Regenerasi lebih cepat
- Hemat quota API
- Bisa kerja offline (jika di-cache)

**Invalidasi cache:**
```bash
# Hapus semua cache
rm -rf ~/.anaphase/cache

# Hapus cache domain tertentu
rm -rf ~/.anaphase/cache/customer*
```

### Request Tuning

Fine-tune perilaku AI:

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

    # Kreativitas AI (default: 0.3)
    # Lebih rendah = lebih konsisten
    # Lebih tinggi = lebih kreatif
    temperature: 0.1

    # Max token dalam response
    maxTokens: 4096
```

## Referensi Konfigurasi

### Opsi AI Provider

| Opsi | Tipe | Default | Deskripsi |
|------|------|---------|-----------|
| `type` | string | - | Tipe provider (`gemini`, `openai`, `claude`, `groq`) |
| `apiKey` | string | - | API key |
| `model` | string | varies | Model yang digunakan (lihat default per provider) |
| `timeout` | duration | `30s` | Request timeout |
| `retries` | int | `3` | Retry attempts |
| `temperature` | float | `0.3` | Kreativitas AI (0.0-1.0) |
| `maxTokens` | int | `4096` | Max response token |

### Opsi Cache

| Opsi | Tipe | Default | Deskripsi |
|------|------|---------|-----------|
| `enabled` | bool | `true` | Aktifkan caching |
| `ttl` | duration | `24h` | Lifetime cache |
| `dir` | string | `~/.anaphase/cache` | Direktori cache |

## Override Environment Variable

Environment variable lebih diutamakan daripada config file:

```bash
# Set API key provider (otomatis mengaktifkan provider)
export GEMINI_API_KEY="your-gemini-key"
export OPENAI_API_KEY="your-openai-key"
export CLAUDE_API_KEY="your-claude-key"
export GROQ_API_KEY="your-groq-key"

# Override path config file
export ANAPHASE_CONFIG="/custom/path/config.yaml"

# Matikan cache
export ANAPHASE_CACHE_ENABLED=false
```

**Catatan:** Setting API key via environment variable otomatis mengaktifkan provider tersebut!

## Verifikasi

Test konfigurasi Anda:

```bash
# Generate test domain dengan verbose logging
anaphase gen domain "Test entity with name field" --verbose

# Anda akan melihat output seperti:
# time=... level=INFO msg="attempting generation" provider=gemini
# time=... level=INFO msg="generation successful" provider=gemini tokens=1234 cost=$0.000000 duration=2.5s

# Atau jika menggunakan Groq:
# time=... level=INFO msg="attempting generation" provider=groq
# time=... level=INFO msg="generation successful" provider=groq tokens=1915 cost=$0.000000 duration=1.6s
```

**Cek provider mana yang digunakan:**
```bash
cat ~/.anaphase/config.yaml | grep primary_provider
```

**Test fallback:**
```bash
# Hapus API key provider primary sementara
unset GEMINI_API_KEY

# Generate domain - seharusnya fallback ke Groq
anaphase gen domain "Test" --verbose
# Anda akan melihat:
# level=WARN msg="provider not available" provider=gemini
# level=INFO msg="attempting generation" provider=groq
```

## Troubleshooting

### Tidak Ada AI Provider (Bukan Error!)

::: info Template Mode Tersedia
Jika Anda melihat pesan tentang AI provider yang hilang, jangan khawatir! Anaphase v0.4.0 menyertakan Template Mode yang bekerja tanpa AI provider apapun. Cukup gunakan flag `--template` atau skip generasi AI sepenuhnya.

```bash
# Gunakan Template Mode sebagai gantinya
anaphase gen domain --name customer --template

# Atau gunakan menu interaktif
anaphase menu
```
:::

### API Key Tidak Ditemukan

```
Error: no AI providers configured - please set at least one API key
```

**Solusi:**
Set API key setidaknya satu provider:

```bash
# Opsi 1: Gunakan Gemini (tier gratis)
export GEMINI_API_KEY="your-gemini-key"

# Opsi 2: Gunakan Groq (tercepat, preview gratis)
export GROQ_API_KEY="your-groq-key"

# Opsi 3: Gunakan OpenAI
export OPENAI_API_KEY="your-openai-key"

# Opsi 4: Gunakan Claude
export CLAUDE_API_KEY="your-claude-key"

# Terbaik: Gunakan multiple provider untuk fallback
export GEMINI_API_KEY="your-gemini-key"
export GROQ_API_KEY="your-groq-key"

# Atau gunakan Template Mode (tidak perlu API key!)
anaphase gen domain --name customer --template
```

### Quota Terlampaui

```
Error: quota exceeded (429)
```

**Solusi:**
1. Tunggu (quota reset per menit)
2. Gunakan fallback provider
3. Aktifkan caching untuk mengurangi call
4. Upgrade ke paid tier

### Response Tidak Valid

```
Error: failed to parse AI response
```

**Solusi:**
1. Coba lagi (AI bisa tidak konsisten)
2. Turunkan temperature ke 0.1
3. Cek model (gunakan gemini-2.5-flash)
4. Aktifkan verbose logging

### Timeout

```
Error: request timeout
```

**Solusi:**
1. Tingkatkan timeout di config:
   ```yaml
   timeout: 60s
   ```
2. Cek koneksi internet
3. Coba model yang berbeda

## Best Practice

### Keamanan

- **Jangan commit API key** ke git
- Gunakan environment variable di CI/CD
- Rotasi key secara berkala
- Gunakan key berbeda per environment

```bash
# Development
export GEMINI_API_KEY="dev-key"

# Production
export GEMINI_API_KEY="prod-key"
```

### Performa

- Aktifkan caching
- Gunakan gemini-2.5-flash (tercepat)
- Set timeout yang wajar
- Gunakan temperature lebih rendah

### Reliabilitas

- Konfigurasi fallback provider
- Aktifkan retries
- Monitor penggunaan quota
- Biarkan cache aktif

## Lihat Juga

- [Generasi Bertenaga AI](/guide/ai-generation)
- [Instalasi](/guide/installation)
- [gen domain](/reference/gen-domain)
