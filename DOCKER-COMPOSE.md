# Docker Compose Documentation - Tugas 2

## ✅ Requirements Checklist

### 1. Service Definition
✅ **Implemented** dengan 2 services:
- **db**: PostgreSQL 18.1 Alpine
- **app**: Wallet API application

### 2. Healthchecks
✅ **Implemented** untuk kedua service:

#### PostgreSQL Healthcheck
```yaml
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U wallet_user -d wallet_db"]
  interval: 10s
  timeout: 5s
  retries: 5
  start_period: 10s
```

#### App Healthcheck
```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8000/healthz"]
  interval: 30s
  timeout: 3s
  retries: 3
  start_period: 10s
```

### 3. Depends_on with Healthcheck Condition
✅ **Implemented**:
```yaml
depends_on:
  db:
    condition: service_healthy
```

**Bukti berhasil**:
```
Container wallet_db Waiting
Container wallet_db Healthy
Container wallet_app Starting
```

App hanya start setelah database healthy ✅

### 4. Persistence
✅ **Implemented** dengan named volume:
```yaml
volumes:
  db_data:
    driver: local
```

Mount: `wallet_api_db_data` → `/var/lib/postgresql/data`

## Usage

### Start Services
```bash
# Start di background
docker compose up -d

# Start di foreground (lihat logs real-time)
docker compose up

# Setelah start, app menunggu DB healthy sebelum running
```

### Stop Services
```bash
# Stop tapi preserve volumes
docker compose down

# Stop dan remove volumes (hapus semua data)
docker compose down -v
```

### View Logs
```bash
# Lihat semua logs
docker compose logs -f

# Lihat app logs saja
docker compose logs -f app

# Lihat db logs saja
docker compose logs -f db
```

### Check Status
```bash
# Lihat status semua services
docker compose ps

# Lihat detail health status
docker inspect wallet_db --format='{{.State.Health.Status}}'
docker inspect wallet_app --format='{{.State.Health.Status}}'
```

## Environment Variables

### Database (PostgreSQL)
- `POSTGRES_DB`: wallet_db
- `POSTGRES_USER`: wallet_user
- `POSTGRES_PASSWORD`: wallet_password

### Application
- `APP_NAME`: wallet_api
- `APP_VERSION`: 1.0.0
- `HTTP_PORT`: 8000
- `LOG_LEVEL`: debug
- `PG_URL`: postgres://wallet_user:wallet_password@db:5432/wallet_db
- `PG_POOL_MAX`: 2

## Networks & Volumes

### Networks
- **app_network**: Bridge network untuk komunikasi antar services

### Volumes
- **db_data**: Named volume untuk PostgreSQL data persistence
  - Location: `/var/lib/docker/volumes/wallet_api_db_data/_data`
  - Data persist meski container di-remove

## Testing

### Test Health Endpoint
```bash
# Test health check
curl http://localhost:8000/healthz

# Expected response
{"status":"ok"}
```

### Test API Endpoints
```bash
# Register user
curl -X POST http://localhost:8000/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'

# Login
curl -X POST http://localhost:8000/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

### Test Persistence
```bash
# 1. Insert data
docker compose exec db psql -U wallet_user -d wallet_db -c "INSERT INTO users (username, password_hash) VALUES ('persist_test', 'hash');"

# 2. Stop containers (volume tetap ada)
docker compose down

# 3. Start lagi
docker compose up -d

# 4. Cek data masih ada
docker compose exec db psql -U wallet_user -d wallet_db -c "SELECT * FROM users;"
```

## Resource Limits

App service memiliki resource limits:
- **CPU Limit**: 1.0 core
- **Memory Limit**: 512MB
- **CPU Reservation**: 0.5 core
- **Memory Reservation**: 256MB

## Port Mappings

- **App**: `8000:8000` (host:container)
- **DB**: `5432:5432` (host:container)

## Restart Policy

App service menggunakan: `restart: unless-stopped`

Container akan otomatis restart jika crash, kecuali di-stop manually.

## Next Steps

Untuk melanjutkan ke **Tugas 3 (CI/CD Pipeline)**:
1. Setup GitHub Actions workflow
2. Configure Docker Hub / GHCR credentials
3. Implement build & push setelah test sukses
