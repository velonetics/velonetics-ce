# Basic Authentication

**Edition:** Enterprise (implemented in Pucora CE)  
**Namespace:** `auth/basic`  
**Module:** `pucora-basicauth`  
**Official docs:** [Basic Authentication](../../../docs/authorization/02-basic-authentication.md)

## What it does

Protects endpoints with HTTP Basic Authentication using bcrypt `.htpasswd` files or inline `users` maps. Service-level config is inherited by endpoints that declare `"auth/basic": {}`.

## Configuration

```json
{
  "extra_config": {
    "auth/basic": {
      "htpasswd_path": "/etc/pucora/.htpasswd",
      "users": {
        "alice": "$2a$10$..."
      }
    }
  },
  "endpoints": [{
    "endpoint": "/admin",
    "extra_config": { "auth/basic": {} }
  }]
}
```

## Parity status

| Area | Status |
|------|--------|
| htpasswd (bcrypt) | Implemented |
| Inline users map | Implemented |
| Service → endpoint inheritance | Implemented |
| 401 before backend contact | Implemented |

**Fixture:** `tests/fixtures/basic_auth.json`
