# Error Types Management - Implementation Summary

## Overview

You now have complete functionality to create and delete error types dynamically through API endpoints, along with ready-to-use scripts for automation.

## What Was Implemented

### 1. **API Endpoints** (Go Backend)

#### Create Error Types (Bulk)
```
POST /v1/admin/error-types
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

Request Body:
[
  {
    "identifier": "error-id",
    "title": "Error Title",
    "detail": "Error description",
    "status": 400
  }
]

Response (201 Created):
{
  "message": "error types created successfully",
  "count": 1
}
```

#### Delete Error Type
```
DELETE /v1/admin/error-types/{identifier}
Authorization: Bearer <JWT_TOKEN>

Response (200 OK):
{
  "message": "error type deleted successfully",
  "identifier": "error-id"
}
```

### 2. **Backend Changes**

**Files Modified:**

1. **[internal/repository/errorFirestore/firestore.go](internal/repository/errorFirestore/firestore.go)**
   - Added `DeleteErrorType(ctx, identifier)` method
   - Added `DeleteErrorInstance(ctx, identifier)` method
   - Updated interface to include delete operations

2. **[internal/usecase/errorList/usecase.go](internal/usecase/errorList/usecase.go)**
   - Added `CreateErrorTypeRequest` struct
   - Added `CreateErrorTypesFromList(ctx, requests)` method
   - Added `DeleteErrorType(ctx, identifier)` method
   - Updated `ErrorListUseCase` interface

3. **[internal/controller/system/handler.go](internal/controller/system/handler.go)**
   - Added `CreateErrorTypes(e *echo.Context)` handler
   - Added `DeleteErrorType(e *echo.Context)` handler
   - Updated `Handler` interface

4. **[internal/controller/controller.go](internal/controller/controller.go)**
   - Added routes for error type management
   - Both endpoints are admin-protected (JWT required)
   - Routes:
     - `POST /v1/admin/error-types` → Create
     - `DELETE /v1/admin/error-types/:identifier` → Delete

### 3. **Python Script** (`scripts/seed_error_types.py`)

Full-featured Python 3 script with colored output and multiple modes.

**Features:**
- Authenticate with username/password
- Create all 10 default error types
- Delete specific error types
- Delete all error types
- Create and delete (test mode)
- Colored terminal output
- Error handling with detailed messages
- Configurable wait time between operations

**Usage Examples:**

```bash
# Create error types
python3 scripts/seed_error_types.py \
    --create \
    --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app \
    --username admin \
    --password your_password

# Delete all error types
python3 scripts/seed_error_types.py \
    --delete-all \
    --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app \
    --username admin \
    --password your_password

# Create, wait 5 seconds, then delete (for testing)
python3 scripts/seed_error_types.py \
    --create-and-delete \
    --wait 5 \
    --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app \
    --username admin \
    --password your_password

# Local development
python3 scripts/seed_error_types.py \
    --create \
    --api-url http://localhost:8080 \
    --username admin \
    --password password
```

### 4. **Bash Script** (`scripts/seed_error_types.sh`)

Lightweight bash script using `curl` (no dependencies beyond bash and curl).

**Usage Examples:**

```bash
# Create error types
./scripts/seed_error_types.sh \
    --action create \
    --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app \
    --username admin \
    --password your_password

# Delete all error types
./scripts/seed_error_types.sh \
    --action delete-all \
    --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app \
    --username admin \
    --password your_password

# Create and delete after 10 seconds
./scripts/seed_error_types.sh \
    --action create-and-delete \
    --wait 10 \
    --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app \
    --username admin \
    --password your_password
```

### 5. **Documentation** (`scripts/README.md`)

Comprehensive guide including:
- Installation instructions
- Usage examples
- API endpoint documentation
- List of all 10 error types
- Troubleshooting guide
- Environment variable configuration
- CI/CD integration tips

## Default Error Types

All scripts automatically manage these 10 error types:

| Identifier | Status | Title | Description |
|-----------|--------|-------|-------------|
| `not-found` | 404 | Not Found | Resource not found |
| `invalid-type` | 400 | Invalid Type | Type not supported |
| `invalid-credentials` | 401 | Invalid Credentials | Wrong username/password |
| `invalid-token` | 401 | Invalid Token | Token invalid or expired |
| `unauthorized` | 401 | Unauthorized | Not authorized |
| `type-not-found` | 404 | Error Type Not Found | Error type not found |
| `instance-not-found` | 404 | Error Instance Not Found | Error instance not found |
| `invalidDocumentLenght` | 400 | Invalid Document Length | CPF/CNPJ invalid length |
| `invalidFilter` | 400 | Invalid Filter | Filter parameter invalid |
| `unexpectedUnhandledError` | 500 | Unexpected Error | Unexpected server error |

## Testing Instructions

### 1. **Local Testing**

```bash
# Terminal 1: Start the API
go run cmd/main.go

# Terminal 2: Create error types
python3 scripts/seed_error_types.py \
    --create \
    --api-url http://localhost:8080 \
    --username juliosshoji \
    --password password123

# Verify in Firestore or with curl:
curl -s http://localhost:8080/v1/error?filter=type&identifier=not-found | head -c 200

# Delete error types
python3 scripts/seed_error_types.py \
    --delete-all \
    --api-url http://localhost:8080 \
    --username juliosshoji \
    --password password123
```

### 2. **Production Testing**

```bash
# Create and then immediately delete (verify both endpoints work)
python3 scripts/seed_error_types.py \
    --create-and-delete \
    --wait 2 \
    --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app \
    --username your_admin_user \
    --password your_password
```

### 3. **Using cURL Directly**

```bash
# Get JWT token
TOKEN=$(curl -s -X POST http://localhost:8080/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"username":"juliosshoji","password":"password123"}' \
    | grep -o '"token":"[^"]*' | cut -d'"' -f4)

# Create error types
curl -X POST http://localhost:8080/v1/admin/error-types \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '[{"identifier":"test-error","title":"Test Error","detail":"Test","status":400}]'

# Delete error type
curl -X DELETE http://localhost:8080/v1/admin/error-types/test-error \
    -H "Authorization: Bearer $TOKEN"
```

## Integration Examples

### CI/CD Pipeline (GitHub Actions)

```yaml
- name: Seed Error Types to Production
  run: |
    python3 scripts/seed_error_types.py \
      --create \
      --api-url ${{ secrets.API_URL }} \
      --username ${{ secrets.API_USER }} \
      --password ${{ secrets.API_PASSWORD }}
```

### Docker Initialization

```dockerfile
# In your Dockerfile or docker-compose.yml
RUN python3 -m pip install requests && \
    python3 scripts/seed_error_types.py \
      --create \
      --api-url http://localhost:8080 \
      --username admin \
      --password password
```

### Deployment Script

```bash
#!/bin/bash
# deploy.sh

# 1. Deploy new API version
kubectl set image deployment/delicias-da-lu api=$NEW_IMAGE

# 2. Wait for rollout
kubectl rollout status deployment/delicias-da-lu

# 3. Seed error types
python3 scripts/seed_error_types.py \
    --create \
    --api-url https://api.delicias-da-lu.com \
    --username $ADMIN_USER \
    --password $ADMIN_PASSWORD
```

## Security Notes

1. **JWT Authentication**: Both endpoints require valid JWT tokens
2. **Admin Role**: Only authenticated admin users can access these endpoints
3. **Credentials**: 
   - Never commit credentials to version control
   - Use environment variables or secret management tools
   - Scripts only accept credentials as command-line arguments
4. **Network**: Use HTTPS in production (endpoint is already HTTPS)

## Troubleshooting

### "Authentication failed"
- Verify username/password are correct
- Check user has admin privileges
- Ensure API is accessible

### "Connection refused"
- Verify API URL is correct and accessible
- Check API is running (for local: `go run cmd/main.go`)
- Verify network connectivity

### "Permission denied"
- Endpoints require JWT authentication
- User must have admin role
- Check token is being passed correctly

### Error types not appearing in Firestore
- Verify the create endpoint returned 201 status
- Check logs for errors: `logging.ErrorEventFromEcho(e, err)`
- Verify Firestore connection works

## Files Created/Modified

**New Files:**
- `scripts/seed_error_types.py` - Python script
- `scripts/seed_error_types.sh` - Bash script
- `scripts/README.md` - Documentation
- `IMPLEMENTATION_SUMMARY.md` - This file

**Modified Files:**
- `internal/repository/errorFirestore/firestore.go` - Delete methods
- `internal/usecase/errorList/usecase.go` - Create/Delete logic
- `internal/controller/system/handler.go` - API handlers
- `internal/controller/controller.go` - Routes registration

## Next Steps

1. **Test locally**: 
   ```bash
   go run cmd/main.go &
   python3 scripts/seed_error_types.py --create --api-url http://localhost:8080 --username admin --password password
   ```

2. **Deploy to production**: Use the scripts in your CI/CD pipeline

3. **Monitor**: Check logs and Firestore collections for operations

4. **Automate**: Integrate into deployment/initialization scripts

## Support

For issues or questions:
1. Check `scripts/README.md` troubleshooting section
2. Review API logs: `go run cmd/main.go`
3. Inspect Firestore collections: Console → Firestore → Collections

