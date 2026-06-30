# Authentication Guide

Outline CLI mendukung 2 jenis token authentication:

## 1. API Key (Recommended) ✅

API Key adalah token khusus untuk API access yang dibuat di Outline settings.

### Format
```
ol_api_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### Cara Mendapatkan
1. Login ke Outline: https://outline-rbi.jatismobile.com
2. Buka Settings → API: https://outline-rbi.jatismobile.com/settings/api
3. Klik "Create new token" atau "New API token"
4. Copy token yang muncul (format: `ol_api_...`)
5. Export ke environment:
   ```bash
   export OUTLINE_API_KEY='ol_api_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
   # atau
   export OUTLINE_TOKEN='ol_api_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
   ```

### Keuntungan
- ✅ Khusus untuk API access
- ✅ Tidak expire (sampai di-revoke manual)
- ✅ Bisa di-revoke kapan saja
- ✅ Lebih aman untuk automation

### Usage
```bash
export OUTLINE_API_KEY='ol_api_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
outline clone <collection-id> <directory>
```

---

## 2. JWT Session Token (Limited) ⚠️

JWT token adalah session token dari login web browser.

### Format
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6Ii...
```

### Cara Mendapatkan
1. Login ke Outline via browser
2. Buka Developer Tools (F12)
3. Cek localStorage atau cookies untuk token
4. Copy JWT token

### Keterbatasan
- ⚠️ **TIDAK BISA DIPAKAI UNTUK API** - Outline API menolak JWT session token dengan error 403
- ⚠️ Expire dalam waktu tertentu (biasanya beberapa hari/minggu)
- ⚠️ Tied to browser session
- ⚠️ Tidak cocok untuk automation

### Status
JWT token detection sudah implemented di CLI, tapi **Outline API server menolak JWT token**.
Kamu tetap harus menggunakan API Key (`ol_api_...`).

---

## Environment Variables

CLI mendukung 2 environment variable (priority order):

```bash
# Priority 1: OUTLINE_API_KEY
export OUTLINE_API_KEY='ol_api_xxxxxxxx'

# Priority 2: OUTLINE_TOKEN (fallback)
export OUTLINE_TOKEN='ol_api_xxxxxxxx'

# Base URL (optional, default: https://outline-rbi.jatismobile.com)
export OUTLINE_BASE_URL='https://outline.example.com'
```

---

## Token Detection

CLI automatically detects token type:

| Token Format | Detected As | Works? |
|--------------|-------------|--------|
| `ol_api_...` | API Key | ✅ Yes |
| `eyJ...` (3 dots) | JWT | ❌ No (403 from server) |

---

## Troubleshooting

### Error: "OUTLINE_API_KEY or OUTLINE_TOKEN not set"
**Solution**: Export environment variable
```bash
export OUTLINE_API_KEY='ol_api_xxxxxxxx'
```

### Error: "HTTP 403: authorization_error"
**Possible causes**:
1. ❌ Using JWT session token instead of API key
2. ❌ API key expired or revoked
3. ❌ Wrong API key

**Solution**: Generate new API key dari Outline settings

### Error: "dial tcp: lookup outline-rbi.jatismobile.com: no such host"
**Possible causes**:
1. DNS resolution issue
2. Network connectivity problem

**Solution**: 
```bash
# Test connectivity
ping outline-rbi.jatismobile.com
curl -I https://outline-rbi.jatismobile.com

# Check DNS
cat /etc/resolv.conf
```

---

## Security Best Practices

1. **Never commit tokens** to git
   ```bash
   # Add to .gitignore
   echo ".env" >> .gitignore
   ```

2. **Use environment variables**
   ```bash
   # In ~/.bashrc or ~/.zshrc
   export OUTLINE_API_KEY='ol_api_xxxxxxxx'
   ```

3. **Rotate tokens regularly**
   - Create new token
   - Update environment variable
   - Revoke old token

4. **Use different tokens** for different environments
   ```bash
   # Development
   export OUTLINE_API_KEY='ol_api_dev_xxxxxxxx'
   
   # Production
   export OUTLINE_API_KEY='ol_api_prod_xxxxxxxx'
   ```

---

## Quick Start

```bash
# 1. Get API key from Outline settings
#    https://outline-rbi.jatismobile.com/settings/api

# 2. Export to environment
export OUTLINE_API_KEY='ol_api_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'

# 3. Test connection
outline clone test-collection-id test-dir

# 4. If it works, add to shell config for persistence
echo "export OUTLINE_API_KEY='ol_api_xxxxxxxx'" >> ~/.bashrc
source ~/.bashrc
```
