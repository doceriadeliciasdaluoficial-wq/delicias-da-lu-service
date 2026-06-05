# Quick Reference - Error Types Management

## 🚀 Quick Start

### Python Script (Recommended)
```bash
# Create error types
python3 scripts/seed_error_types.py --create \
    --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app \
    --username admin --password password

# Delete all error types
python3 scripts/seed_error_types.py --delete-all \
    --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app \
    --username admin --password password

# Create and delete (test mode)
python3 scripts/seed_error_types.py --create-and-delete --wait 5 \
    --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app \
    --username admin --password password
```

### Bash Script
```bash
# Create error types
./scripts/seed_error_types.sh --action create \
    --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app \
    --username admin --password password

# Delete all error types
./scripts/seed_error_types.sh --action delete-all \
    --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app \
    --username admin --password password
```

## 🔌 API Endpoints

### Create Error Types
```bash
curl -X POST https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/admin/error-types \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '[
    {"identifier":"error-id","title":"Error Title","detail":"Description","status":400}
  ]'
```

### Delete Error Type
```bash
curl -X DELETE https://delicias-da-lu-514609008596.southamerica-east1.run.app/v1/admin/error-types/error-id \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 📝 Error Types

| Identifier | Status | Description |
|-----------|--------|-------------|
| not-found | 404 | Resource not found |
| invalid-type | 400 | Invalid type |
| invalid-credentials | 401 | Wrong credentials |
| invalid-token | 401 | Invalid token |
| unauthorized | 401 | Not authorized |
| type-not-found | 404 | Error type not found |
| instance-not-found | 404 | Error instance not found |
| invalidDocumentLenght | 400 | CPF/CNPJ invalid |
| invalidFilter | 400 | Invalid filter |
| unexpectedUnhandledError | 500 | Server error |

## 🔐 Authentication

1. **Get JWT Token:**
   ```bash
   curl -X POST http://localhost:8080/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username":"admin","password":"password"}'
   ```

2. **Use Token:**
   ```bash
   Authorization: Bearer <token_here>
   ```

## 💻 Local Testing

```bash
# Terminal 1: Start API
go run cmd/main.go

# Terminal 2: Create error types
python3 scripts/seed_error_types.py --create \
    --api-url http://localhost:8080 \
    --username admin --password password

# Terminal 2: Delete error types
python3 scripts/seed_error_types.py --delete-all \
    --api-url http://localhost:8080 \
    --username admin --password password
```

## ✅ Response Examples

### Create Success (201)
```json
{
  "message": "error types created successfully",
  "count": 10
}
```

### Delete Success (200)
```json
{
  "message": "error type deleted successfully",
  "identifier": "not-found"
}
```

## 🐛 Troubleshooting

| Issue | Solution |
|-------|----------|
| Authentication failed | Verify username/password, check user is admin |
| Connection refused | Verify API URL and that API is running |
| 401 Unauthorized | Get a valid JWT token, pass it in Authorization header |
| Error type not found | Verify identifier spelling, create if missing |

## 📚 Full Documentation

See `scripts/README.md` and `IMPLEMENTATION_SUMMARY.md` for complete documentation.

