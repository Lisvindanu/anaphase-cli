# Konfigurasi Database

::: info Auto-Setup di v0.4.0
Konfigurasi database sekarang **otomatis**! Saat `anaphase init`, Anda akan diminta memilih database, dan Anaphase akan otomatis:
- Membuat file `.env` dengan connection string
- Membuat contoh nilai DATABASE_URL
- Setup kode repository untuk database yang dipilih

Tidak perlu konfigurasi manual untuk memulai!
:::

Konfigurasi koneksi database untuk microservice Anda.

## Database yang Didukung

- **PostgreSQL** (Direkomendasikan)
- **MySQL** / MariaDB
- **MongoDB**

## Pemilihan Database Interaktif

::: info Baru di v0.4.0
Saat Anda menjalankan `anaphase init`, Anda akan melihat menu pemilihan database interaktif:

```bash
anaphase init my-project

# Prompt interaktif:
? Select database:
  > PostgreSQL (recommended)
    MySQL
    MongoDB
    Skip (configure later)

# Anaphase otomatis:
# 1. Membuat file .env dengan DATABASE_URL
# 2. Generate kode repository untuk database yang dipilih
# 3. Setup connection pooling
# 4. Menyertakan file schema
```
:::

## Connection String

### Environment Variable (Auto-Generated)

Setelah `anaphase init`, cek file `.env` Anda untuk DATABASE_URL yang auto-generated:

```bash
# .env (auto-created oleh Anaphase v0.4.0)
DATABASE_URL="postgres://user:password@localhost:5432/mydb?sslmode=disable"
```

Anda juga bisa set secara manual:

```bash
export DATABASE_URL="postgres://user:password@host:port/database"
```

### File Konfigurasi

`cmd/api/wire.go` yang dihasilkan membaca dari DATABASE_URL:

```go
func InitializeApp(logger *slog.Logger) (*App, error) {
    // Otomatis membaca dari file .env (dibuat oleh anaphase init)
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        // Default fallback
        dbURL = "postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable"
    }

    db, err := pgxpool.New(context.Background(), dbURL)
    // ...
}
```

::: tip Konfigurasi Auto-Generated
File `.env` otomatis dibuat saat Anda menjalankan `anaphase init` dan memilih database. Anda bisa kustomisasi connection string dengan mengedit file ini.
:::

## PostgreSQL

::: info Auto-Setup di v0.4.0
Saat Anda memilih PostgreSQL selama `anaphase init`, file `.env` Anda akan berisi:
```bash
DATABASE_URL="postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable"
```
Edit file ini untuk menyesuaikan kredensial database Anda.
:::

### Format Connection String

```
postgres://username:password@hostname:port/database?options
```

**Contoh:**
```bash
# File .env (auto-generated, lalu dikustomisasi)
DATABASE_URL="postgres://myuser:mypass@localhost:5432/mydb?sslmode=disable"
```

### Opsi Umum

| Opsi | Nilai | Deskripsi |
|------|-------|-----------|
| `sslmode` | `disable`, `require`, `verify-full` | Mode SSL |
| `pool_max_conns` | number | Koneksi maksimal |
| `pool_min_conns` | number | Koneksi minimal |
| `pool_max_conn_lifetime` | duration | Lifetime koneksi maks |

**Contoh production:**
```bash
DATABASE_URL="postgres://user:pass@prod-db.example.com:5432/mydb?\
sslmode=require&\
pool_max_conns=20&\
pool_min_conns=5&\
pool_max_conn_lifetime=1h"
```

### Setup Docker

```bash
docker run -d \
  --name postgres \
  -e POSTGRES_USER=myuser \
  -e POSTGRES_PASSWORD=mypass \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 \
  postgres:16-alpine
```

### Instalasi Lokal

::: code-group

```bash [macOS]
brew install postgresql@16
brew services start postgresql@16
createdb mydb
```

```bash [Ubuntu/Debian]
sudo apt-get install postgresql-16
sudo systemctl start postgresql
sudo -u postgres createdb mydb
```

```bash [Windows]
# Download dari postgresql.org
# Atau gunakan Docker
```

:::

### Apply Schema

```bash
psql $DATABASE_URL -f internal/adapter/repository/postgres/schema.sql
```

## MySQL

::: info Auto-Setup di v0.4.0
Saat Anda memilih MySQL selama `anaphase init`, file `.env` Anda akan berisi:
```bash
DATABASE_URL="root:password@tcp(localhost:3306)/mydb?parseTime=true"
```
Edit file ini untuk menyesuaikan kredensial MySQL Anda.
:::

### Format Connection String

```
username:password@tcp(hostname:port)/database?options
```

**Contoh:**
```bash
# File .env (auto-generated, lalu dikustomisasi)
DATABASE_URL="myuser:mypass@tcp(localhost:3306)/mydb?parseTime=true"
```

### Opsi Umum

| Opsi | Nilai | Deskripsi |
|------|-------|-----------|
| `parseTime` | `true`, `false` | Parse nilai TIME |
| `charset` | `utf8mb4` | Character set |
| `collation` | `utf8mb4_unicode_ci` | Collation |

**Contoh production:**
```bash
DATABASE_URL="user:pass@tcp(prod-db:3306)/mydb?\
parseTime=true&\
charset=utf8mb4&\
collation=utf8mb4_unicode_ci&\
maxAllowedPacket=67108864"
```

### Setup Docker

```bash
docker run -d \
  --name mysql \
  -e MYSQL_ROOT_PASSWORD=rootpass \
  -e MYSQL_USER=myuser \
  -e MYSQL_PASSWORD=mypass \
  -e MYSQL_DATABASE=mydb \
  -p 3306:3306 \
  mysql:8
```

### Apply Schema

```bash
mysql -h localhost -u myuser -p mydb < internal/adapter/repository/mysql/schema.sql
```

## MongoDB

::: info Auto-Setup di v0.4.0
Saat Anda memilih MongoDB selama `anaphase init`, file `.env` Anda akan berisi:
```bash
DATABASE_URL="mongodb://localhost:27017/mydb"
```
Edit file ini untuk menambahkan autentikasi dan opsi lainnya.
:::

### Format Connection String

```
mongodb://username:password@hostname:port/database?options
```

**Contoh:**
```bash
# File .env (auto-generated, lalu dikustomisasi)
DATABASE_URL="mongodb://myuser:mypass@localhost:27017/mydb"
```

### Opsi Umum

| Opsi | Nilai | Deskripsi |
|------|-------|-----------|
| `authSource` | database | Database auth |
| `replicaSet` | name | Replica set |
| `ssl` | `true`, `false` | Gunakan SSL |

**Contoh production:**
```bash
DATABASE_URL="mongodb://user:pass@mongo1:27017,mongo2:27017,mongo3:27017/mydb?\
replicaSet=rs0&\
ssl=true&\
authSource=admin"
```

### Setup Docker

```bash
docker run -d \
  --name mongo \
  -e MONGO_INITDB_ROOT_USERNAME=admin \
  -e MONGO_INITDB_ROOT_PASSWORD=pass \
  -e MONGO_INITDB_DATABASE=mydb \
  -p 27017:27017 \
  mongo:7
```

### Collection

Collection MongoDB dibuat otomatis saat insert pertama.

## Connection Pooling

### PostgreSQL (pgxpool)

Konfigurasi di kode:

```go
import "github.com/jackc/pgx/v5/pgxpool"

config, err := pgxpool.ParseConfig(dbURL)
if err != nil {
    return err
}

// Pengaturan connection pool
config.MaxConns = 20
config.MinConns = 5
config.MaxConnLifetime = time.Hour
config.MaxConnIdleTime = 30 * time.Minute

db, err := pgxpool.NewWithConfig(context.Background(), config)
```

Atau via connection string:
```bash
DATABASE_URL="postgres://...?pool_max_conns=20&pool_min_conns=5"
```

### MySQL (sql.DB)

```go
import "database/sql"

db, err := sql.Open("mysql", dbURL)

// Pengaturan connection pool
db.SetMaxOpenConns(20)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(time.Hour)
db.SetConnMaxIdleTime(30 * time.Minute)
```

## Konfigurasi Per-Environment

::: tip Fitur v0.4.0
File `.env` yang auto-generated oleh `anaphase init` sempurna untuk development. Buat file `.env` tambahan untuk environment berbeda.
:::

### Development

```bash
# .env (auto-created oleh anaphase init)
DATABASE_URL="postgres://postgres:postgres@localhost:5432/myapp_dev?sslmode=disable"
```

### Testing

```bash
# .env.test (buat manual)
DATABASE_URL="postgres://postgres:postgres@localhost:5432/myapp_test?sslmode=disable"
```

### Production

```bash
# .env.production (buat manual)
DATABASE_URL="postgres://user:pass@prod-db:5432/myapp?sslmode=require&pool_max_conns=50"
```

## Migrasi

### Schema Awal

Dihasilkan bersama repository:

```bash
anaphase gen repository --domain customer --db postgres
# Membuat: internal/adapter/repository/postgres/schema.sql

psql $DATABASE_URL -f internal/adapter/repository/postgres/schema.sql
```

### Tool Migrasi

Gunakan tool migrasi untuk production:

**golang-migrate:**
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Buat migrasi
migrate create -ext sql -dir db/migrations -seq add_customers

# Jalankan migrasi
migrate -database $DATABASE_URL -path db/migrations up
```

**goose:**
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest

# Jalankan migrasi
goose -dir db/migrations postgres $DATABASE_URL up
```

## Health Check

Verifikasi konektivitas database:

### Di Kode

```go
// Ping database
if err := db.Ping(context.Background()); err != nil {
    logger.Error("database ping failed", "error", err)
    return err
}
```

### Endpoint

main.go yang dihasilkan menyertakan health check:

```bash
curl http://localhost:8080/health
# Return: OK (jika database terhubung)
```

## Troubleshooting

### Connection Refused

```
Error: connection refused
```

**Solusi:**
1. Cek apakah database berjalan:
   ```bash
   # PostgreSQL
   docker ps | grep postgres
   pg_isready -h localhost

   # MySQL
   docker ps | grep mysql
   mysqladmin ping -h localhost
   ```

2. Cek apakah port benar
3. Cek aturan firewall

### Authentication Failed

```
Error: authentication failed
```

**Solusi:**
1. Verifikasi username/password
2. Cek apakah database ada:
   ```bash
   psql -h localhost -U postgres -l
   ```
3. Verifikasi permission user

### Too Many Connection

```
Error: sorry, too many clients already
```

**Solusi:**
1. Kurangi ukuran pool:
   ```bash
   DATABASE_URL="...?pool_max_conns=10"
   ```
2. Tingkatkan max connection database:
   ```sql
   ALTER SYSTEM SET max_connections = 200;
   ```

### SSL Required

```
Error: SSL required
```

**Solusi:**
```bash
DATABASE_URL="...?sslmode=require"
```

## Best Practice

### Keamanan

- **Jangan commit kredensial** ke git
- Gunakan environment variable
- Rotasi password secara berkala
- Gunakan SSL di production
- Batasi permission user

```sql
-- Buat app user dengan permission terbatas
CREATE USER myapp WITH PASSWORD 'secure_password';
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO myapp;
```

### Performa

- Gunakan connection pooling
- Set ukuran pool yang sesuai
- Monitor query lambat
- Tambahkan index untuk query sering
- Gunakan prepared statement

### Reliabilitas

- Aktifkan automatic reconnection
- Set connection timeout
- Monitor connection pool
- Gunakan read replica untuk scaling

## Lihat Juga

- [gen repository](/reference/gen-repository)
- [Arsitektur](/guide/architecture)
- [Contoh](/examples/basic)
