# anaphase gen migration

Generate file migration database dengan intelligent SQL generation.

::: info
**Quick Start**: Jalankan `anaphase` (tanpa argumen) untuk mengakses menu interaktif dan pilih "Generate Migration" untuk pengalaman terpandu.
:::

## Overview

Command `gen migration` membuat file migration database dengan timestamp beserta script up/down. Command ini secara cerdas generate SQL berdasarkan konvensi penamaan migration dan mendukung multiple database driver.

::: info
**Berbasis Template**: Generasi migration menggunakan smart template berdasarkan konvensi penamaan - tidak perlu konfigurasi AI. Cepat dan andal.
:::

## Penggunaan

### Menu Interaktif (Disarankan)

```bash
anaphase
```

Pilih **"Generate Migration"** dari menu. Interface akan memandu Anda melalui:
- Nama/deskripsi migration
- Pemilihan database driver
- Konfigurasi direktori output

### Mode CLI Langsung

```bash
anaphase gen migration <name> [flags]
```

## Database yang Didukung

- **PostgreSQL** (default) - Dukungan penuh dengan trigger
- **MySQL** - Dukungan SQL standar
- **SQLite** - Dukungan database lightweight

## Konvensi Penamaan

::: tip
Anaphase secara cerdas parse nama migration Anda dan generate SQL yang sesuai - tidak perlu menulis SQL manual untuk pattern umum.
:::

Anaphase secara otomatis generate SQL yang sesuai berdasarkan nama migration Anda:

### Create Table

```bash
anaphase gen migration create_users_table
anaphase gen migration create_orders_table
```

Generate:
- `CREATE TABLE` dengan id, created_at, updated_at
- PostgreSQL: Trigger auto-update untuk updated_at
- Down migration: `DROP TABLE`

### Add Column

```bash
anaphase gen migration add_email_to_users
anaphase gen migration add_total_to_orders
```

Generate:
- `ALTER TABLE ADD COLUMN` dengan tipe yang di-infer
- Down migration: `DROP COLUMN`

**Inferensi Tipe:**
- `*_id` → `BIGINT`
- `*_amount`, `*_price` → `DECIMAL(10,2)`
- `*_count`, `*_quantity` → `INTEGER`
- `is_*`, `has_*` → `BOOLEAN`
- `*_at`, `*_date` → `TIMESTAMP`
- Default → `VARCHAR(255)`

### Drop Table

```bash
anaphase gen migration drop_old_cache_table
```

Generate:
- `DROP TABLE IF EXISTS`
- Down migration: Pengingat untuk backup data

## Flag

| Flag | Default | Deskripsi |
|------|---------|-------------|
| `--output` | `db/migrations` | Direktori output untuk file migration |
| `--driver` | `postgres` | Database driver (postgres, mysql, sqlite) |

## Contoh

### Menggunakan Menu Interaktif

```bash
# Luncurkan menu
anaphase

# Pilih "Generate Migration"
# Masukkan nama migration: create_users_table
# Pilih database: PostgreSQL
# Pilih direktori output: db/migrations
```

### Menggunakan CLI Langsung

**Create Users Table:**

```bash
anaphase gen migration create_users_table
```

Output:
```
✓ db/migrations/20251222101843_create_users_table.up.sql
✓ db/migrations/20251222101843_create_users_table.down.sql
```

SQL yang dihasilkan (up):
```sql
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Trigger untuk auto-update updated_at
CREATE OR REPLACE FUNCTION update_users_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_users_updated_at();
```

SQL yang dihasilkan (down):
```sql
DROP TABLE IF EXISTS users CASCADE;
```

**Add Email Column:**

```bash
anaphase gen migration add_email_to_users
```

SQL yang dihasilkan (up):
```sql
ALTER TABLE users ADD COLUMN email VARCHAR(255);
```

SQL yang dihasilkan (down):
```sql
ALTER TABLE users DROP COLUMN IF EXISTS email;
```

**Direktori Output Kustom:**

```bash
anaphase gen migration create_products_table --output migrations
```

**Database MySQL:**

```bash
anaphase gen migration create_orders_table --driver mysql
```

## Migration Tools

Anaphase generate file migration SQL standar yang kompatibel dengan migration tools populer:

### golang-migrate

```bash
# Install
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Jalankan migration
migrate -path db/migrations -database "postgres://user:pass@localhost/db?sslmode=disable" up

# Rollback
migrate -path db/migrations -database "postgres://user:pass@localhost/db?sslmode=disable" down 1
```

### goose

```bash
# Install
go install github.com/pressly/goose/v3/cmd/goose@latest

# Jalankan migration
goose -dir db/migrations postgres "user=postgres dbname=mydb" up

# Rollback
goose -dir db/migrations postgres "user=postgres dbname=mydb" down
```

### sql-migrate

```bash
# Install
go install github.com/rubenv/sql-migrate/...@latest

# Jalankan migration
sql-migrate up -config=dbconfig.yml
```

## Workflow

### 1. Generate Migration

```bash
anaphase gen migration create_users_table
```

### 2. Edit SQL (jika diperlukan)

Edit file yang dihasilkan untuk menambahkan kolom, index, constraint:

```sql
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
```

### 3. Test Secara Lokal

```bash
# Terapkan migration
migrate -path db/migrations -database "$DATABASE_URL" up

# Verifikasi schema
psql $DATABASE_URL -c "\dt"

# Test rollback
migrate -path db/migrations -database "$DATABASE_URL" down 1

# Re-apply
migrate -path db/migrations -database "$DATABASE_URL" up
```

### 4. Terapkan ke Production

```bash
# Backup terlebih dahulu!
pg_dump -U postgres mydb > backup_$(date +%Y%m%d_%H%M%S).sql

# Jalankan migration
migrate -path db/migrations -database "$PROD_DATABASE_URL" up
```

## Best Practice

### 1. Selalu Reversible

Pastikan down migration dapat sepenuhnya reverse up migration:

```sql
-- ✅ Baik - reversible
-- up.sql
ALTER TABLE users ADD COLUMN email VARCHAR(255);

-- down.sql
ALTER TABLE users DROP COLUMN email;

-- ❌ Buruk - data loss
-- up.sql
ALTER TABLE users DROP COLUMN old_email;

-- down.sql
-- Tidak bisa restore data yang di-drop!
```

### 2. Satu Perubahan Per Migration

```bash
# ✅ Baik
anaphase gen migration add_email_to_users
anaphase gen migration add_phone_to_users

# ❌ Buruk - multiple perubahan dalam satu migration
anaphase gen migration update_users_table
```

### 3. Test Rollback

Selalu test bahwa down migration bekerja:

```bash
migrate up
# Verifikasi perubahan
migrate down 1
# Verifikasi rollback
migrate up
```

### 4. Gunakan Transaction

Wrap migration kompleks dalam transaction (jika didukung):

```sql
BEGIN;

ALTER TABLE users ADD COLUMN email VARCHAR(255);
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

COMMIT;
```

### 5. Backup Sebelum Production

```bash
# PostgreSQL
pg_dump -U postgres mydb > backup.sql

# MySQL
mysqldump -u root -p mydb > backup.sql

# SQLite
sqlite3 mydb.db ".backup backup.db"
```

## Contoh Advanced

### Tambahkan Index

```bash
anaphase gen migration create_index_users_email
```

Edit file yang dihasilkan:
```sql
CREATE INDEX idx_users_email ON users(email);
```

### Tambahkan Foreign Key

```bash
anaphase gen migration add_user_id_to_orders
```

Edit file yang dihasilkan:
```sql
ALTER TABLE orders ADD COLUMN user_id BIGINT;
ALTER TABLE orders ADD CONSTRAINT fk_orders_user
    FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE;
```

### Ubah Tipe Column

```bash
anaphase gen migration change_email_length
```

```sql
-- up.sql
ALTER TABLE users ALTER COLUMN email TYPE VARCHAR(320);

-- down.sql
ALTER TABLE users ALTER COLUMN email TYPE VARCHAR(255);
```

## Troubleshooting

### Migration Gagal

```bash
# Cek versi saat ini
migrate -path db/migrations -database "$DATABASE_URL" version

# Force ke versi tertentu (gunakan dengan hati-hati!)
migrate -path db/migrations -database "$DATABASE_URL" force 20251222101843

# Fix issue dan retry
migrate -path db/migrations -database "$DATABASE_URL" up
```

### Dirty Database State

```bash
# Cek status
migrate -path db/migrations -database "$DATABASE_URL" version

# Output: 20251222101843/d (dirty)

# Fix state database secara manual
psql $DATABASE_URL

# Kemudian force clean
migrate -path db/migrations -database "$DATABASE_URL" force 20251222101843
```

## Lihat Juga

- [anaphase gen domain](/reference/gen-domain) - Generate entity domain
- [anaphase gen repository](/reference/gen-repository) - Generate repository
- [Database Configuration](/config/database) - Setup database
