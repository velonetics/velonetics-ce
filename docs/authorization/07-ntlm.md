# NTLM Authentication

**Edition:** Enterprise (implemented in Pucora CE)  
**Namespace:** `auth/ntlm`  
**Module:** `pucora-ntlm`  
**Official docs:** [NTLM authentication](../../../docs/authorization/17-ntlm-authentication.md)

## What it does

Wraps the backend HTTP client transport with NTLMv2 challenge-response authentication for legacy Windows/IIS backends.

## Configuration

```json
{
  "backend": [{
    "host": ["http://iis.internal"],
    "extra_config": {
      "auth/ntlm": {
        "user": "DOMAIN\\serviceaccount",
        "password": "secret"
      }
    }
  }]
}
```

## Parity status

| Area | Status |
|------|--------|
| NTLMv2 via `go-ntlmssp` | Implemented |
| HTTP client factory wrapping | Implemented |
| Chained with other backend auth | Implemented |

**Fixture:** `tests/fixtures/backend_ntlm.json`

**Note:** Full integration testing requires a Windows/IIS backend; unit coverage focuses on config parsing and transport wiring.
