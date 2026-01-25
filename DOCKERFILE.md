# Dockerfile Documentation - Tugas 1

## ✅ Requirements Checklist

### 1. Multi-stage Build
✅ **Implemented** dengan 3 stages:
- **Stage 1 (modules)**: Download dependencies menggunakan `golang:1.25-alpine3.21`
- **Stage 2 (builder)**: Compile binary menggunakan cached dependencies
- **Stage 3 (runner)**: Runtime environment menggunakan `alpine:latest`

### 2. Security
✅ **Non-root user**:
- Membuat user `appuser` dengan UID 1000
- Menggunakan `USER appuser` instruction
- File permissions di-set dengan `chown -R appuser:appuser /app`

### 3. Optimization
✅ **Layer Caching**:
- Dependencies di-download terpisah di stage 1 (modules)
- `go.mod` dan `go.sum` di-copy sebelum source code
- Cached dependencies disalin ke builder stage

✅ **Minimal Image Size**:
- Base image: `alpine:latest` (±5MB)
- Final image size: **66.5MB**
- Hanya menyertakan binary dan migrations

## Build Commands

```bash
# Build image
docker build -t wallet_api:test .

# Build dengan tag version
docker build -t wallet_api:v1.0.0 .

# Build untuk specific platform
docker buildx build --platform linux/amd64 -t wallet_api:latest .
```

## Run Commands

```bash
# Run container (butuh database)
docker run -p 8000:8000 \
  -e DB_HOST=host.docker.internal \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASS=postgres \
  -e DB_NAME=wallet_db \
  wallet_api:test
```

## Verification

```bash
# Cek ukuran image
docker images wallet_api:test

# Cek user (harus appuser, bukan root)
docker inspect wallet_api:test --format='User: {{.Config.User}}'

# Cek layers
docker history wallet_api:test

# Test healthcheck
docker run --rm wallet_api:test curl -f http://localhost:8000/healthz
```

## Health Check

✅ **Docker Healthcheck**:
- Interval: 30s
- Timeout: 3s
- Start period: 5s
- Retries: 3
- Command: `curl -f http://localhost:8000/healthz`

## Files Created/Modified

1. **Dockerfile** - Multi-stage build dengan 3 stages
2. **.dockerignore** - Optimize build context dan exclude files

## Next Steps

Untuk melanjutkan ke **Tugas 2 (Docker Compose)**:
1. Update `docker-compose.yml` untuk service app
2. Tambah healthchecks untuk database dan app
3. Implement depends_on dengan condition
4. Setup named volumes untuk persistence
