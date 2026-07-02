# Outline Docker Setup

Setup untuk deploy Outline menggunakan Docker Compose.

## Prerequisites

- Docker dan Docker Compose terinstall
- Port 3000 tersedia
- Minimal 2GB RAM

## Quick Start

### 1. Konfigurasi Environment

File `docker.env` sudah dibuat dengan konfigurasi default. **PENTING**: Sebelum production, ubah:

```bash
# Generate secret keys baru
openssl rand -hex 32  # untuk SECRET_KEY
openssl rand -hex 32  # untuk UTILS_SECRET

# Update di docker.env:
SECRET_KEY=<key_baru>
UTILS_SECRET=<key_baru>

# Ubah password database
POSTGRES_PASSWORD=<password_kuat>
# Dan update di DATABASE_URL juga
```

### 2. Konfigurasi Authentication

Outline memerlukan **minimal 1 provider authentication**. Pilih salah satu:

#### Option A: Google OAuth
1. Buat OAuth credentials di [Google Cloud Console](https://console.cloud.google.com/apis/credentials)
2. Authorized redirect URI: `http://localhost:3000/auth/google.callback`
3. Tambahkan ke `docker.env`:
```bash
GOOGLE_CLIENT_ID=your_client_id
GOOGLE_CLIENT_SECRET=your_client_secret
```

#### Option B: Slack OAuth
1. Buat Slack App di [api.slack.com/apps](https://api.slack.com/apps)
2. Redirect URL: `http://localhost:3000/auth/slack.callback`
3. Tambahkan ke `docker.env`:
```bash
SLACK_CLIENT_ID=your_client_id
SLACK_CLIENT_SECRET=your_client_secret
```

#### Option C: Azure AD / Microsoft Entra
1. Register app di Azure Portal
2. Redirect URI: `http://localhost:3000/auth/azure.callback`
3. Tambahkan ke `docker.env`:
```bash
AZURE_CLIENT_ID=your_client_id
AZURE_CLIENT_SECRET=your_client_secret
AZURE_RESOURCE_APP_ID=your_resource_app_id
```

### 3. Jalankan Docker Compose

```bash
# Start semua services
docker compose up -d

# Check logs
docker compose logs -f outline

# Check status
docker compose ps
```

### 4. Akses Outline

Buka browser: http://localhost:3000

## Production Deployment

Untuk production, update `docker.env`:

```bash
# Update URL ke domain production
URL=https://your-domain.com

# Enable HTTPS
FORCE_HTTPS=true

# Gunakan reverse proxy (nginx/caddy) untuk SSL termination
```

### Dengan HTTPS Portal (Optional)

Uncomment service `https-portal` di `docker-compose.yml` dan update:

```yaml
environment:
  DOMAINS: 'your-domain.com -> http://outline:3000'
  STAGE: 'production'
```

## File Storage

### Local Storage (Default)
Data disimpan di Docker volume `storage-data`. Untuk backup:

```bash
docker run --rm -v outline-cli_storage-data:/data -v $(pwd):/backup alpine tar czf /backup/outline-storage-backup.tar.gz -C /data .
```

### S3 Storage (Recommended untuk Production)

Update `docker.env`:

```bash
FILE_STORAGE=s3
AWS_ACCESS_KEY_ID=your_key
AWS_SECRET_ACCESS_KEY=your_secret
AWS_REGION=us-east-1
AWS_S3_UPLOAD_BUCKET_NAME=your-bucket
AWS_S3_ACL=private
```

## Database Backup

```bash
# Backup database
docker compose exec postgres pg_dump -U outline outline > outline-backup.sql

# Restore database
docker compose exec -T postgres psql -U outline outline < outline-backup.sql
```

## Troubleshooting

### Check logs
```bash
docker compose logs -f outline
docker compose logs -f postgres
docker compose logs -f redis
```

### Reset database
```bash
docker compose down -v
docker compose up -d
```

### Update Outline
```bash
docker compose pull outline
docker compose up -d outline
```

## Commands

```bash
# Start
docker compose up -d

# Stop
docker compose down

# Restart
docker compose restart outline

# View logs
docker compose logs -f

# Remove everything (including data)
docker compose down -v
```

## Struktur File

```
outline-cli/
├── docker-compose.yml    # Docker Compose configuration
├── docker.env           # Environment variables (JANGAN commit ke git!)
├── .env.example         # Template environment variables
└── DOCKER_SETUP.md      # Dokumentasi ini
```

## Security Notes

1. **JANGAN commit `docker.env` ke git** - file ini berisi secrets
2. Gunakan `.env.example` sebagai template
3. Untuk production:
   - Gunakan password database yang kuat
   - Generate SECRET_KEY dan UTILS_SECRET yang unik
   - Enable HTTPS
   - Gunakan S3 untuk file storage
   - Setup regular database backups

## Resources

- [Outline Documentation](https://docs.getoutline.com)
- [Docker Documentation](https://docs.getoutline.com/s/hosting/doc/docker-7pfeLP5a8t)
- [Configuration Guide](https://docs.getoutline.com/s/hosting/doc/configuration-509J4lAzjo)
