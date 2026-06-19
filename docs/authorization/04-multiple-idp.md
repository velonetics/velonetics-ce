# Multiple Identity Providers

**Edition:** Enterprise (implemented in Pucora CE)  
**Namespace:** `plugin/http-server` → `jwk-aggregator`  
**Module:** `pucora-jwk-aggregator`  
**Official docs:** [Multiple Identity Providers](../../../docs/authorization/09-multiple-identity-providers.md)

## What it does

HTTP server plugin that starts a localhost JWK aggregation sidecar. Fetches JWK sets from multiple IdP origins in parallel, merges them, and serves the combined set for JWT validation.

## Configuration

```json
{
  "extra_config": {
    "plugin/http-server": {
      "name": ["jwk-aggregator"],
      "jwk-aggregator": {
        "port": 9876,
        "cache": true,
        "origins": [
          "https://idp1.example.com/.well-known/jwks.json",
          "https://idp2.example.com/.well-known/jwks.json"
        ]
      }
    }
  }
}
```

Point JWT validator at the sidecar:

```json
{
  "auth/validator": {
    "jwk_url": "http://localhost:9876",
    "disable_jwk_security": true
  }
}
```

## Build

```bash
make -C pucora-ce jwk-aggregator-plugin
```

## Parity status

| Area | Status |
|------|--------|
| Plugin load (`jwk-aggregator.so`) | Implemented |
| Localhost-only JWK server | Implemented |
| Multi-origin merge | Implemented |
| Cache-Control honoring | Implemented |

**Fixture:** `tests/fixtures/multi_idp.json`
