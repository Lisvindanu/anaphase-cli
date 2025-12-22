# anaphase gen migration

Generate database migration files with intelligent SQL generation.

## Overview

The `gen migration` command creates timestamped database migration files with up/down scripts. It intelligently generates SQL based on migration naming conventions and supports multiple database drivers.

## Usage

```bash
anaphase gen migration <name> [flags]
```

## Supported Databases

- **PostgreSQL** (default) - Full support with triggers
- **MySQL** - Standard SQL support
- **SQLite** - Lightweight database support

## Naming Conventions

Anaphase automatically generates appropriate SQL based on your migration name:

### Create Table

```bash
anaphase gen migration create_users_table
anaphase gen migration create_orders_table
```

Generates:
- `CREATE TABLE` with id, created_at, updated_at
- PostgreSQL: Auto-update trigger for updated_at
- Down migration: `DROP TABLE`

### Add Column

```bash
anaphase gen migration add_email_to_users
anaphase gen migration add_total_to_orders
```

Generates:
- `ALTER TABLE ADD COLUMN` with inferred type
- Down migration: `DROP COLUMN`

**Type Inference:**
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

Generates:
- `DROP TABLE IF EXISTS`
- Down migration: Reminder to backup data

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--output` | `db/migrations` | Output directory for migration files |
| `--driver` | `postgres` | Database driver (postgres, mysql, sqlite) |

## Examples

### Create Users Table

```bash
anaphase gen migration create_users_table
```

Output:
```
✓ db/migrations/20251222101843_create_users_table.up.sql
✓ db/migrations/20251222101843_create_users_table.down.sql
```

Generated SQL (up):
```sql
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Trigger to auto-update updated_at
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

Generated SQL (down):
```sql
DROP TABLE IF EXISTS users CASCADE;
```

### Add Email Column

```bash
anaphase gen migration add_email_to_users
```

Generated SQL (up):
```sql
ALTER TABLE users ADD COLUMN email VARCHAR(255);
```

Generated SQL (down):
```sql
ALTER TABLE users DROP COLUMN IF EXISTS email;
```

### Custom Output Directory

```bash
anaphase gen migration create_products_table --output migrations
```

### MySQL Database

```bash
anaphase gen migration create_orders_table --driver mysql
```

## Migration Tools

Anaphase generates standard SQL migration files compatible with popular migration tools:

### golang-migrate

```bash
# Install
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path db/migrations -database "postgres://user:pass@localhost/db?sslmode=disable" up

# Rollback
migrate -path db/migrations -database "postgres://user:pass@localhost/db?sslmode=disable" down 1
```

### goose

```bash
# Install
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir db/migrations postgres "user=postgres dbname=mydb" up

# Rollback
goose -dir db/migrations postgres "user=postgres dbname=mydb" down
```

### sql-migrate

```bash
# Install
go install github.com/rubenv/sql-migrate/...@latest

# Run migrations
sql-migrate up -config=dbconfig.yml
```

## Workflow

### 1. Generate Migration

```bash
anaphase gen migration create_users_table
```

### 2. Edit SQL (if needed)

Edit the generated files to add columns, indexes, constraints:

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

### 3. Test Locally

```bash
# Apply migration
migrate -path db/migrations -database "$DATABASE_URL" up

# Verify schema
psql $DATABASE_URL -c "\dt"

# Test rollback
migrate -path db/migrations -database "$DATABASE_URL" down 1

# Re-apply
migrate -path db/migrations -database "$DATABASE_URL" up
```

### 4. Apply to Production

```bash
# Backup first!
pg_dump -U postgres mydb > backup_$(date +%Y%m%d_%H%M%S).sql

# Run migration
migrate -path db/migrations -database "$PROD_DATABASE_URL" up
```

## Best Practices

### 1. Always Reversible

Ensure down migrations can fully reverse up migrations:

```sql
-- ✅ Good - reversible
-- up.sql
ALTER TABLE users ADD COLUMN email VARCHAR(255);

-- down.sql
ALTER TABLE users DROP COLUMN email;

-- ❌ Bad - data loss
-- up.sql
ALTER TABLE users DROP COLUMN old_email;

-- down.sql
-- Can't restore dropped data!
```

### 2. One Change Per Migration

```bash
# ✅ Good
anaphase gen migration add_email_to_users
anaphase gen migration add_phone_to_users

# ❌ Bad - do multiple changes in one migration
anaphase gen migration update_users_table
```

### 3. Test Rollbacks

Always test that down migrations work:

```bash
migrate up
# Verify changes
migrate down 1
# Verify rollback
migrate up
```

### 4. Use Transactions

Wrap complex migrations in transactions (when supported):

```sql
BEGIN;

ALTER TABLE users ADD COLUMN email VARCHAR(255);
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

COMMIT;
```

### 5. Backup Before Production

```bash
# PostgreSQL
pg_dump -U postgres mydb > backup.sql

# MySQL
mysqldump -u root -p mydb > backup.sql

# SQLite
sqlite3 mydb.db ".backup backup.db"
```

## Advanced Examples

### Add Index

```bash
anaphase gen migration create_index_users_email
```

Edit generated file:
```sql
CREATE INDEX idx_users_email ON users(email);
```

### Add Foreign Key

```bash
anaphase gen migration add_user_id_to_orders
```

Edit generated file:
```sql
ALTER TABLE orders ADD COLUMN user_id BIGINT;
ALTER TABLE orders ADD CONSTRAINT fk_orders_user
    FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE;
```

### Modify Column Type

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

### Migration Failed

```bash
# Check current version
migrate -path db/migrations -database "$DATABASE_URL" version

# Force to specific version (use with caution!)
migrate -path db/migrations -database "$DATABASE_URL" force 20251222101843

# Fix issue and retry
migrate -path db/migrations -database "$DATABASE_URL" up
```

### Dirty Database State

```bash
# Check status
migrate -path db/migrations -database "$DATABASE_URL" version

# Output: 20251222101843/d (dirty)

# Manually fix the database state
psql $DATABASE_URL

# Then force clean
migrate -path db/migrations -database "$DATABASE_URL" force 20251222101843
```

## See Also

- [anaphase gen domain](/reference/gen-domain) - Generate domain entities
- [anaphase gen repository](/reference/gen-repository) - Generate repositories
- [Database Configuration](/config/database) - Database setup
