# Database Configuration

Configure database connections for your microservice.

## Supported Databases

- **PostgreSQL** (Recommended)
- **MySQL** / MariaDB
- **MongoDB**

## Connection String

### Environment Variable

Set `DATABASE_URL` environment variable:

```bash
export DATABASE_URL="postgres://user:password@host:port/database"
```

### Configuration File

Or configure in generated `cmd/api/wire.go`:

```go
func InitializeApp(logger *slog.Logger) (*App, error) {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        dbURL = "postgres://postgres:postgres@localhost:5432/mydb?sslmode=disable"
    }

    db, err := pgxpool.New(context.Background(), dbURL)
    // ...
}
```

## PostgreSQL

### Connection String Format

```
postgres://username:password@hostname:port/database?options
```

**Example:**
```bash
export DATABASE_URL="postgres://myuser:mypass@localhost:5432/mydb?sslmode=disable"
```

### Common Options

| Option | Values | Description |
|--------|--------|-------------|
| `sslmode` | `disable`, `require`, `verify-full` | SSL mode |
| `pool_max_conns` | number | Max connections |
| `pool_min_conns` | number | Min connections |
| `pool_max_conn_lifetime` | duration | Max conn lifetime |

**Production example:**
```bash
DATABASE_URL="postgres://user:pass@prod-db.example.com:5432/mydb?\
sslmode=require&\
pool_max_conns=20&\
pool_min_conns=5&\
pool_max_conn_lifetime=1h"
```

### Docker Setup

```bash
docker run -d \
  --name postgres \
  -e POSTGRES_USER=myuser \
  -e POSTGRES_PASSWORD=mypass \
  -e POSTGRES_DB=mydb \
  -p 5432:5432 \
  postgres:16-alpine
```

### Local Installation

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
# Download from postgresql.org
# Or use Docker
```

:::

### Apply Schema

```bash
psql $DATABASE_URL -f internal/adapter/repository/postgres/schema.sql
```

## MySQL

### Connection String Format

```
username:password@tcp(hostname:port)/database?options
```

**Example:**
```bash
export DATABASE_URL="myuser:mypass@tcp(localhost:3306)/mydb?parseTime=true"
```

### Common Options

| Option | Values | Description |
|--------|--------|-------------|
| `parseTime` | `true`, `false` | Parse TIME values |
| `charset` | `utf8mb4` | Character set |
| `collation` | `utf8mb4_unicode_ci` | Collation |

**Production example:**
```bash
DATABASE_URL="user:pass@tcp(prod-db:3306)/mydb?\
parseTime=true&\
charset=utf8mb4&\
collation=utf8mb4_unicode_ci&\
maxAllowedPacket=67108864"
```

### Docker Setup

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

### Connection String Format

```
mongodb://username:password@hostname:port/database?options
```

**Example:**
```bash
export DATABASE_URL="mongodb://myuser:mypass@localhost:27017/mydb"
```

### Common Options

| Option | Values | Description |
|--------|--------|-------------|
| `authSource` | database | Auth database |
| `replicaSet` | name | Replica set |
| `ssl` | `true`, `false` | Use SSL |

**Production example:**
```bash
DATABASE_URL="mongodb://user:pass@mongo1:27017,mongo2:27017,mongo3:27017/mydb?\
replicaSet=rs0&\
ssl=true&\
authSource=admin"
```

### Docker Setup

```bash
docker run -d \
  --name mongo \
  -e MONGO_INITDB_ROOT_USERNAME=admin \
  -e MONGO_INITDB_ROOT_PASSWORD=pass \
  -e MONGO_INITDB_DATABASE=mydb \
  -p 27017:27017 \
  mongo:7
```

### Collections

MongoDB collections are created automatically on first insert.

## Connection Pooling

### PostgreSQL (pgxpool)

Configure in code:

```go
import "github.com/jackc/pgx/v5/pgxpool"

config, err := pgxpool.ParseConfig(dbURL)
if err != nil {
    return err
}

// Connection pool settings
config.MaxConns = 20
config.MinConns = 5
config.MaxConnLifetime = time.Hour
config.MaxConnIdleTime = 30 * time.Minute

db, err := pgxpool.NewWithConfig(context.Background(), config)
```

Or via connection string:
```bash
DATABASE_URL="postgres://...?pool_max_conns=20&pool_min_conns=5"
```

### MySQL (sql.DB)

```go
import "database/sql"

db, err := sql.Open("mysql", dbURL)

// Connection pool settings
db.SetMaxOpenConns(20)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(time.Hour)
db.SetConnMaxIdleTime(30 * time.Minute)
```

## Environment-Specific Configuration

### Development

```bash
# .env.development
DATABASE_URL="postgres://postgres:postgres@localhost:5432/myapp_dev?sslmode=disable"
```

### Testing

```bash
# .env.test
DATABASE_URL="postgres://postgres:postgres@localhost:5432/myapp_test?sslmode=disable"
```

### Production

```bash
# .env.production
DATABASE_URL="postgres://user:pass@prod-db:5432/myapp?sslmode=require&pool_max_conns=50"
```

## Migrations

### Initial Schema

Generated with repository:

```bash
anaphase gen repository --domain customer --db postgres
# Creates: internal/adapter/repository/postgres/schema.sql

psql $DATABASE_URL -f internal/adapter/repository/postgres/schema.sql
```

### Migration Tools

Use migration tool for production:

**golang-migrate:**
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Create migration
migrate create -ext sql -dir db/migrations -seq add_customers

# Run migrations
migrate -database $DATABASE_URL -path db/migrations up
```

**goose:**
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir db/migrations postgres $DATABASE_URL up
```

## Health Checks

Verify database connectivity:

### In Code

```go
// Ping database
if err := db.Ping(context.Background()); err != nil {
    logger.Error("database ping failed", "error", err)
    return err
}
```

### Endpoint

Generated main.go includes health check:

```bash
curl http://localhost:8080/health
# Returns: OK (if database connected)
```

## Troubleshooting

### Connection Refused

```
Error: connection refused
```

**Solutions:**
1. Check database is running:
   ```bash
   # PostgreSQL
   docker ps | grep postgres
   pg_isready -h localhost

   # MySQL
   docker ps | grep mysql
   mysqladmin ping -h localhost
   ```

2. Check port is correct
3. Check firewall rules

### Authentication Failed

```
Error: authentication failed
```

**Solutions:**
1. Verify username/password
2. Check database exists:
   ```bash
   psql -h localhost -U postgres -l
   ```
3. Verify user permissions

### Too Many Connections

```
Error: sorry, too many clients already
```

**Solutions:**
1. Reduce pool size:
   ```bash
   DATABASE_URL="...?pool_max_conns=10"
   ```
2. Increase database max connections:
   ```sql
   ALTER SYSTEM SET max_connections = 200;
   ```

### SSL Required

```
Error: SSL required
```

**Solution:**
```bash
DATABASE_URL="...?sslmode=require"
```

## Best Practices

### Security

- **Never commit credentials** to git
- Use environment variables
- Rotate passwords regularly
- Use SSL in production
- Limit user permissions

```sql
-- Create app user with limited permissions
CREATE USER myapp WITH PASSWORD 'secure_password';
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO myapp;
```

### Performance

- Use connection pooling
- Set appropriate pool size
- Monitor slow queries
- Add indexes for frequent queries
- Use prepared statements

### Reliability

- Enable automatic reconnection
- Set connection timeouts
- Monitor connection pool
- Use read replicas for scaling

## See Also

- [gen repository](/reference/gen-repository)
- [Architecture](/guide/architecture)
- [Examples](/examples/basic)
