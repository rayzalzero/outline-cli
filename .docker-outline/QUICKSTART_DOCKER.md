# Quick Start - Deploy Outline dengan Docker

Panduan cepat untuk deploy Outline menggunakan Docker Compose.

## Prerequisites

- Docker dan Docker Compose terinstall
- Port 3000 tersedia
- Minimal 2GB RAM

## Langkah Deploy

### 1. Setup Authentication Provider

Outline **WAJIB** punya minimal 1 authentication provider. Pilih salah satu:

#### Google OAuth (Paling Mudah)

1. Buka [Google Cloud Console](https://console.cloud.google.com/apis/credentials)
2. Create Project baru atau pilih existing
3. Enable "Google+ API"
4. Create OAuth 2.0 Client ID:
   - Application type: Web application
   - Authorized redirect URIs: `http://localhost:3000/auth/google.callback`
5. Copy Client ID dan Client Secret

#### Slack OAuth

1. Buka [Slack API](https://api.slack.com/apps)
2. Create New App → From scratch
3. OAuth & snip Permissions:
   - Redirect URL: `http://localhost:3000/auth/slack.callback`
   - Scopes: `identity.basic`, `identity.email`, `identity.avatar`
4. Install to Workspace
5. Copy Client ID dan Client Secret

### 2. Edit File docker.env

```bash
# Edit file docker.env
nano docker.env

# Update bagian ini:
SECRET_KEY=<hasil_openssl_rand_hex_32>
UTILS_SECRET=<hasil_openssl_rand_hex_32>

# Uncomment dan isi salah satu provider:
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret

# ATAU

SLACK_CLIENT_ID=your_slack_client_id
SLACK_CLIENT_SECRET=your_slack_client_secret
```

### 3. Generate Secret Keys

```bash
# Generate 2 secret keys
openssl rand -hex 32
openssl rand -hex 32

# Copy hasil ke docker.env untuk SECRET_KEY dan UTILS_SECRET
```

### 4. Start Docker Compose

```bash
# Start semua services
make docker-up

# Atau langsung:
docker compose up -d
```

### 5. Check Logs

```bash
# Lihat logs Outline
make docker-logs

# Atau:
docker compose logs -f outline

# Pastikan tidak ada error authentication
```

### 6. Akses Outline

Buka browser: **http://localhost:3000**

Login menggunakan provider yang sudah dikonfigurasi (Google/Slack).

## Commands

```bash
# Start
make docker-up

# Stop
make docker-down

# Restart
make docker-restart

# View logs
make docker-logs

# Clean everything (termasuk data!)
make docker-clean
```

## Troubleshooting

### Error: "No authentication providers configured"

**Solusi**: Pastikan minimal 1 provider sudah dikonfigurasi di `docker.env`:
- `GOOGLE_CLIENT_ID` + `GOOGLE_CLIENT_SECRET`, atau
- `SLACK_CLIENT_ID` + `SLACK_CLIENT_SECRET`, atau
- Provider lain (Azure, OIDC, Discord)

### Error: "Database connection failed"

**Solusi**: 
```bash
# Restart postgres
docker compose restart postgres

# Check postgres logs
docker compose logs postgres
```

### Error: "Redis connection failed"

**Solusi**:
```bash
# Restart redis
docker compose restart redis

# Check redis logs
docker compose logs redis
```

### Port 3000 sudah dipakai

**Solusi**: Edit `docker-compose.yml`, ubah port mapping:
```yaml
ports:
  - "3001:3000"  # Ganti 3001 dengan port yang tersedia
```

Lalu akses di: http://localhost:3001

## Production Deployment

Untuk production, update `docker.env`:

```bash
# Update URL
URL=https://your-domain.com

# Enable HTTPS
FORCE_HTTPS=true

# Ganti password database
POSTGRES_PASSWORD=password_yang_sangat_kuat

# Update DATABASE_URL juga
DATABASE_URL=postgres://outline:password_yang_sangat_kuat@postgres:5432/outline

# Gunakan S3 untuk file storage (recommended)
FILE_STORAGE=s3
AWS_ACCESS_KEY_ID=your_key
AWS_SECRET_ACCESS_KEY=your_secret
AWS_REGION=us-east-1
AWS_S3_UPLOAD_BUCKET_NAME=your-bucket
```

Setup reverse proxy (nginx/caddy) untuk SSL termination.

## Backup

### Database Backup

```bash
# Backup
docker compose exec postgres pg_dump -U outline outline > backup-$(date +%Y%m%d).sql

# Restore
docker compose exec -T postgres psql -U outline outline < backup-20260702.sql
```

### File Storage Backup

```bash
# Backup volume
docker run --rm -v outline-cli_storage-data:/data -v $(pwd):/backup alpine tar czf /backup/storage-backup.tar.gz -C /data .

# Restore volume
docker run --rm -v outline-cli_storage-data:/data -v $(pwd):/backup alpine tar xzf /backup/storage-backup.tar.gz -C /data
```

## Next Steps

Setelah Outline berjalan:

1. **Create Team**: Setup team pertama kali login
2. **Create Collection**: Buat collection untuk organize documents
3. **Invite Users**: Invite team members
4. **Setup CLI**: Install outline-cli untuk sync local

Lihat [DOCKER_SETUP.md](DOCKER_SETUP.md) untuk dokumentasi lengkap.
