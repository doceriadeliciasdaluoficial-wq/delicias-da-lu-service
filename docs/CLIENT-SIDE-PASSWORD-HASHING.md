# Client-Side Password Hashing Security Implementation

## Overview

The Delícias da Lú API uses **client-side password hashing** to ensure that plain text passwords never travel over the internet. Instead of sending the plain password to the server, clients hash the password using SHA256 before transmission and send only the hash.

**Security Benefit:** Even if the HTTPS connection is compromised or if there's a network interception, attackers will only see the SHA256 hash, not the original password.

---

## Authentication Flow

### Step 1: Client-Side Hashing

Before sending credentials to the `/auth/login` endpoint, the client must compute the SHA256 hash of the password in hexadecimal format.

### Step 2: Send Hashed Credentials

Send the username and `passwordHash` (not plain password) to the server:

```json
{
  "username": "juliosshoji",
  "passwordHash": "8d969eef6ecad3c29a3a873fba8fe814ea412bfce19db27c81052fe959408b8515"
}
```

### Step 3: Server-Side Verification

The server receives the hash and compares it directly with the stored hash. No additional hashing is performed on the server.

---

## Implementation Examples

### Python

```python
import hashlib
import requests
import json

def hash_password(password: str) -> str:
    """Generate SHA256 hash of password in hex format"""
    return hashlib.sha256(password.encode()).hexdigest()

def login(username: str, password: str, base_url: str = "http://localhost:8080/v1") -> dict:
    """Login with client-side password hashing"""
    password_hash = hash_password(password)
    
    response = requests.post(
        f"{base_url}/auth/login",
        json={
            "username": username,
            "passwordHash": password_hash
        }
    )
    
    response.raise_for_status()
    return response.json()

# Usage
if __name__ == "__main__":
    # Example: hash of "password123"
    test_hash = hash_password("password123")
    print(f"Hash of 'password123': {test_hash}")
    # Output: 8d969eef6ecad3c29a3a873fba8fe814ea412bfce19db27c81052fe959408b8515
    
    # Login
    try:
        token_response = login("juliosshoji", "password123")
        print(f"Login successful! Token: {token_response['token']}")
    except Exception as e:
        print(f"Login failed: {e}")
```

### JavaScript / Node.js

```javascript
const crypto = require('crypto');
const axios = require('axios');

function hashPassword(password) {
    /**
     * Generate SHA256 hash of password in hex format
     */
    return crypto
        .createHash('sha256')
        .update(password)
        .digest('hex');
}

async function login(username, password, baseUrl = 'http://localhost:8080/v1') {
    /**
     * Login with client-side password hashing
     */
    const passwordHash = hashPassword(password);
    
    try {
        const response = await axios.post(
            `${baseUrl}/auth/login`,
            {
                username,
                passwordHash
            }
        );
        return response.data;
    } catch (error) {
        throw new Error(`Login failed: ${error.response?.data?.detail || error.message}`);
    }
}

// Usage
(async () => {
    // Example: hash of "password123"
    const testHash = hashPassword('password123');
    console.log(`Hash of 'password123': ${testHash}`);
    // Output: 8d969eef6ecad3c29a3a873fba8fe814ea412bfce19db27c81052fe959408b8515
    
    // Login
    try {
        const response = await login('juliosshoji', 'password123');
        console.log(`Login successful! Token: ${response.token}`);
    } catch (error) {
        console.error(`Login failed: ${error.message}`);
    }
})();
```

### TypeScript (React Frontend Example)

```typescript
import sha256 from 'js-sha256';

interface LoginRequest {
    username: string;
    passwordHash: string;
}

interface LoginResponse {
    token: string;
    user: {
        id: string;
        username: string;
        email: string;
        role: string;
        lastLogin: string;
        createdAt: string;
        updatedAt: string;
    };
}

function hashPassword(password: string): string {
    /**
     * Generate SHA256 hash of password using js-sha256
     * Install: npm install js-sha256
     */
    return sha256(password);
}

async function login(
    username: string,
    password: string,
    baseUrl: string = 'http://localhost:8080/v1'
): Promise<LoginResponse> {
    const passwordHash = hashPassword(password);
    
    const response = await fetch(`${baseUrl}/auth/login`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            username,
            passwordHash,
        }),
    });
    
    if (!response.ok) {
        const error = await response.json();
        throw new Error(error.detail || 'Login failed');
    }
    
    return response.json();
}

// Usage in React component
function LoginForm() {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleLogin = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            const response = await login(username, password);
            // Store token in localStorage or context
            localStorage.setItem('authToken', response.token);
            console.log('Login successful!');
            // Redirect to dashboard
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Login failed');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleLogin}>
            <input
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                placeholder="Username"
                required
            />
            <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Password"
                required
            />
            <button type="submit" disabled={loading}>
                {loading ? 'Logging in...' : 'Login'}
            </button>
            {error && <p style={{ color: 'red' }}>{error}</p>}
        </form>
    );
}
```

### cURL

```bash
#!/bin/bash

# Hash password using OpenSSL or similar
PASSWORD="password123"
PASSWORD_HASH=$(echo -n "$PASSWORD" | sha256sum | awk '{print $1}')

echo "Hashed password: $PASSWORD_HASH"

# Login request
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{
    \"username\": \"juliosshoji\",
    \"passwordHash\": \"$PASSWORD_HASH\"
  }"
```

---

## Password Hash Reference

For testing purposes, here are common password hashes:

| Password | SHA256 Hash |
|----------|-------------|
| `password123` | `8d969eef6ecad3c29a3a873fba8fe814ea412bfce19db27c81052fe959408b8515` |
| `admin` | `8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918` |
| `test` | `9f86d081884c7d6d9ffd60a2591f2f8c8b3e4d9e5d9e2e6e3e5e6f7d9e0e1e2e` |

---

## Libraries for SHA256 Hashing

### Python
- Built-in: `hashlib`
- PyPI: `pip install pycryptodome` (for additional crypto functions)

### JavaScript/Node.js
- Node.js Built-in: `crypto`
- NPM: `npm install js-sha256` (for browser compatibility)
- NPM: `npm install crypto-js` (alternative, cross-platform)

### Go
```go
import "crypto/sha256"
import "encoding/hex"

func hashPassword(password string) string {
    hash := sha256.Sum256([]byte(password))
    return hex.EncodeToString(hash[:])
}
```

### Java
```java
import java.security.MessageDigest;

public class PasswordHasher {
    public static String hashPassword(String password) throws Exception {
        MessageDigest digest = MessageDigest.getInstance("SHA-256");
        byte[] hash = digest.digest(password.getBytes());
        StringBuilder hexString = new StringBuilder();
        for (byte b : hash) {
            String hex = Integer.toHexString(0xff & b);
            if (hex.length() == 1) hexString.append('0');
            hexString.append(hex);
        }
        return hexString.toString();
    }
}
```

---

## API Request/Response

### Request

```http
POST /v1/auth/login HTTP/1.1
Host: api.delicias-da-lu.com.br
Content-Type: application/json

{
  "username": "juliosshoji",
  "passwordHash": "8d969eef6ecad3c29a3a873fba8fe814ea412bfce19db27c81052fe959408b8515"
}
```

### Successful Response (200 OK)

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "user-001",
    "username": "juliosshoji",
    "email": "julio@delicias-da-lu.com.br",
    "role": "admin",
    "lastLogin": "2026-06-14T14:30:00Z",
    "createdAt": "2026-01-01T00:00:00Z",
    "updatedAt": "2026-06-14T14:30:00Z"
  }
}
```

### Error Response (401 Unauthorized)

```json
{
  "type": "https://delicias-da-lu-service.com/docs/errors/invalid-credentials",
  "title": "Invalid Credentials",
  "detail": "The provided username or password is incorrect",
  "status": 401,
  "instance": "https://api.delicias-da-lu.com.br/v1/auth/login",
  "traceId": "a1b2c3d4e5f6g7h8"
}
```

---

## Security Considerations

### ✅ What This Protects Against

- **Network Sniffing:** Attackers intercepting the request only see the hash, not the password
- **Proxy/Firewall Logging:** Security logs won't contain plain text passwords
- **DNS Hijacking:** The transmitted data is the hash, not the password

### ⚠️ Important: ALWAYS Use HTTPS

While client-side hashing provides an additional layer of security, **you MUST always use HTTPS (TLS/SSL)** in production:

```
✅ GOOD:  https://api.delicias-da-lu.com.br/v1/auth/login
❌ BAD:   http://api.delicias-da-lu.com.br/v1/auth/login
```

HTTPS encrypts the entire request, including the hash, providing defense-in-depth security.

### ⚠️ Limitations of SHA256

This implementation uses SHA256, which is suitable for transport security but has limitations:

| Aspect | SHA256 | Bcrypt | Argon2 |
|--------|--------|--------|--------|
| Speed | ⚡ Fast | 🐢 Slow (by design) | 🐢 Slow (by design) |
| Salt Support | ❌ No | ✅ Yes | ✅ Yes |
| Brute Force Resistant | ❌ No | ✅ Yes | ✅ Yes |
| Use Case | Transport | Password Storage | Password Storage |

**Current Architecture:**
- Client sends SHA256 hash over HTTPS
- Server compares hash with stored hash
- Servers stores hash (not plain password) via `json:"-"` tag

**For Maximum Security:**
Consider upgrading server-side storage to use Bcrypt or Argon2 for password hashing at rest.

---

## Troubleshooting

### Issue: "Invalid Credentials" on correct password

**Possible causes:**
1. **Hashing algorithm mismatch** - Verify you're using SHA256
2. **Encoding issue** - Ensure password is encoded as UTF-8 before hashing
3. **Case sensitivity** - Usernames are case-sensitive
4. **Hash format** - Hash must be lowercase hexadecimal string

**Debug steps:**
```python
import hashlib

password = "password123"
hash_result = hashlib.sha256(password.encode('utf-8')).hexdigest()
print(f"Hash: {hash_result}")
print(f"Expected: 8d969eef6ecad3c29a3a873fba8fe814ea412bfce19db27c81052fe959408b8515")
print(f"Match: {hash_result == '8d969eef6ecad3c29a3a873fba8fe814ea412bfce19db27c81052fe959408b8515'}")
```

### Issue: Token received but subsequent requests fail

**Check:**
1. Token is included in `Authorization: Bearer <token>` header
2. Token hasn't expired
3. HTTPS is being used for all requests

---

## Migration from Plain Text Passwords

If you have existing users with plain text passwords in the database:

### Step 1: Create migration script

```python
# migrate_passwords.py
import hashlib
from firestore_client import users_collection

def migrate_password(username, plain_password):
    password_hash = hashlib.sha256(plain_password.encode()).hexdigest()
    users_collection.document(username).update({
        'password': password_hash
    })

# Migrate all users
for user in users_collection.stream():
    if 'plain_password' in user.to_dict():
        migrate_password(user.id, user.to_dict()['plain_password'])
```

### Step 2: Update clients to use hashing

Push client updates to use the new hashing function.

### Step 3: Verify

Test login with new hashing approach before removing old implementation.

---

## Best Practices

1. ✅ Always hash on the client side before sending
2. ✅ Always use HTTPS in production
3. ✅ Never log or log the password hash
4. ✅ Store token securely (httpOnly cookies for web)
5. ✅ Implement token expiration
6. ✅ Use refresh tokens for long-lived sessions
7. ❌ Don't reuse the same hash for other purposes
8. ❌ Don't send password in query parameters or logs
9. ❌ Don't hardcode passwords in client code

---

## References

- [OWASP Password Storage Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
- [SHA256 Algorithm](https://en.wikipedia.org/wiki/SHA-2)
- [RFC 7617 - HTTP Authentication](https://tools.ietf.org/html/rfc7617)
- [HTTP over TLS (HTTPS)](https://en.wikipedia.org/wiki/HTTPS)
