# Error Types Management Scripts

This directory contains scripts for managing error types in the Delícias da Lu API.

## Available Scripts

### 1. Python Script (Recommended)

**File:** `seed_error_types.py`

#### Prerequisites
```bash
pip install requests
```

#### Usage

**Create error types:**
```bash
python3 seed_error_types.py \
    --create \
    --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app \
    --username admin \
    --password your_password
```

**Delete a specific error type:**
```bash
python3 seed_error_types.py \
    --delete not-found \
    --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app \
    --username admin \
    --password your_password
```

**Delete all error types:**
```bash
python3 seed_error_types.py \
    --delete-all \
    --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app \
    --username admin \
    --password your_password
```

**Create, wait, then delete all error types:**
```bash
python3 seed_error_types.py \
    --create-and-delete \
    --wait 5 \
    --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app \
    --username admin \
    --password your_password
```

### 2. Bash Script

**File:** `seed_error_types.sh`

#### Prerequisites
- `curl` command-line tool (usually pre-installed)
- `bash` shell

#### Usage

**Create error types:**
```bash
./seed_error_types.sh \
    --action create \
    --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app \
    --username admin \
    --password your_password
```

**Delete all error types:**
```bash
./seed_error_types.sh \
    --action delete-all \
    --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app \
    --username admin \
    --password your_password
```

**Create, wait, then delete all error types:**
```bash
./seed_error_types.sh \
    --action create-and-delete \
    --wait 5 \
    --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app \
    --username admin \
    --password your_password
```

## API Endpoints

### Create Error Types
```
POST /v1/admin/error-types

Headers:
  Authorization: Bearer <JWT_TOKEN>
  Content-Type: application/json

Body:
[
  {
    "identifier": "error-type-1",
    "title": "Error Type Title",
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

### Delete Error Type
```
DELETE /v1/admin/error-types/{identifier}

Headers:
  Authorization: Bearer <JWT_TOKEN>

Response (200 OK):
{
  "message": "error type deleted successfully",
  "identifier": "error-type-1"
}
```

## Error Types Included

The scripts automatically seed all 10 default error types:

1. **not-found** (404) - Resource not found
2. **invalid-type** (400) - Invalid type provided
3. **invalid-credentials** (401) - Wrong username/password
4. **invalid-token** (401) - Invalid or expired token
5. **unauthorized** (401) - Not authorized to access
6. **type-not-found** (404) - Error type not found
7. **instance-not-found** (404) - Error instance not found
8. **invalidDocumentLenght** (400) - CPF/CNPJ invalid length
9. **invalidFilter** (400) - Invalid filter parameter
10. **unexpectedUnhandledError** (500) - Unexpected server error

## Examples

### Local Testing (localhost)

```bash
# Create error types locally
python3 seed_error_types.py \
    --create \
    --api-url http://localhost:8080 \
    --username admin \
    --password password

# Then delete them
python3 seed_error_types.py \
    --delete-all \
    --api-url http://localhost:8080 \
    --username admin \
    --password password
```

### Production Testing

```bash
# Create and immediately delete (test run)
python3 seed_error_types.py \
    --create-and-delete \
    --wait 10 \
    --api-url https://delicias-da-lu-514609008596.southamerica-east1.run.app \
    --username your_username \
    --password your_password
```

## Troubleshooting

### Authentication Failed
- Verify username and password are correct
- Ensure user has admin privileges

### Connection Refused
- Check API is running and accessible at the provided URL
- Verify network connectivity

### Permission Denied
- The endpoints are admin-protected (require JWT token)
- Ensure authenticated user has admin role

### Error Type Already Exists
- The script will silently skip if error type already exists
- Use `--delete-all` to clear before creating

## Environment Variables (Optional)

You can set default values via environment variables:

```bash
export API_URL="https://delicias-da-lu-514609008596.southamerica-east1.run.app"
export API_USERNAME="admin"
export API_PASSWORD="your_password"

# Then script can use simplified syntax
python3 seed_error_types.py --create
```

## Integration

These scripts can be integrated into:
- CI/CD pipelines for testing
- Deployment automation
- Database seeding procedures
- Administrative tools

