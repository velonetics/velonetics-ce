# Revoke Server

**Edition:** Enterprise (implemented in Pucora CE)  
**Namespace:** `auth/revoker` + CLI `pucora revoker`  
**Modules:** `pucora-revoker`, `bloomfilter/pucora`  
**Official docs:** [Revoke Server](../../../docs/authorization/08-revoke-server.md)

## What it does

Centralized token revocation service that fans out bloomfilter updates to registered gateway instances. Gateways with `revoke_server_ping_url` auto-register and receive push revocations.

## Configuration

Revoke server (`pucora revoker -c config.json`):

```json
{
  "extra_config": {
    "auth/revoker": {
      "n": 1000,
      "p": 0.01,
      "ttl": 3600,
      "hash_name": "optimal",
      "revoke_server_api_key": "secret"
    }
  }
}
```

Gateway client sync:

```json
{
  "extra_config": {
    "auth/revoker": {
      "revoke_server_ping_url": "http://revoker:8080",
      "revoke_server_ping_interval": "30s",
      "revoke_server_api_key": "secret"
    }
  }
}
```

## Parity status

| Area | Status |
|------|--------|
| `POST /revoke` fan-out | Implemented |
| Instance registry (`GET/POST /instances`) | Implemented |
| Health endpoint | Implemented |
| API key auth | Implemented |
| Gateway auto-registration | Implemented |
| DIY bloomfilter client (no revoke server) | Existing CE behavior |

**Fixture:** `tests/fixtures/revoke_server.json`
