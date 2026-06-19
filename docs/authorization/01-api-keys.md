# API Keys

**Edition:** Enterprise (implemented in Pucora CE)  
**Namespace:** `auth/api-keys`  
**Module:** `pucora-apikeys`  
**Official docs:** [API-Key authentication](../../../docs/authorization/01-api-keys.md)

## What it does

Validates inbound API keys at the router layer before proxying. Supports header or query-string strategies, optional key hashing, role-based endpoint access, per-key rate limiting, and role propagation to backends.

## Configuration

Service-level keys and defaults:

```json
{
  "extra_config": {
    "auth/api-keys": {
      "strategy": "header",
      "identifier": "Authorization",
      "hash": "plain",
      "keys": [{ "key": "secret", "roles": ["user"] }],
      "propagate_role": "X-Pucora-Role"
    }
  }
}
```

Endpoint-level protection (requires `roles`):

```json
{
  "endpoint": "/protected",
  "extra_config": {
    "auth/api-keys": {
      "roles": ["user"],
      "client_max_rate": 100
    }
  }
}
```

## Parity status

| Area | Status |
|------|--------|
| Header / query strategies | Implemented |
| Bearer / Basic / raw header extraction | Implemented |
| Key hashing (plain, sha256, sha1, fnv128) | Implemented |
| Role checks + propagation | Implemented |
| Per-key rate limit | Implemented |

**Fixture:** `tests/fixtures/api_keys.json`
