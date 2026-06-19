# Google Cloud Authentication

**Edition:** Enterprise (implemented in Pucora CE)  
**Namespace:** `auth/gcp`  
**Module:** `pucora-gcp-auth`  
**Official docs:** [Google Cloud Authentication](../../../docs/authorization/12-google-cloud-authentication.md)

## What it does

Injects a GCP ID token into outbound backend requests for service-to-service authentication with Cloud Run, Cloud Functions, and other Google APIs.

## Configuration

```json
{
  "backend": [{
    "host": ["https://my-service.run.app"],
    "extra_config": {
      "auth/gcp": {
        "audience": "https://my-service.run.app",
        "credentials_file": "/etc/gcp/sa.json",
        "s2s_auth_header": "Authorization"
      }
    }
  }]
}
```

Use `X-Serverless-Authorization` for Cloud Functions by setting `s2s_auth_header`.

## Parity status

| Area | Status |
|------|--------|
| ID token via `idtoken` package | Implemented |
| credentials_file / credentials_json | Implemented |
| GOOGLE_APPLICATION_CREDENTIALS fallback | Implemented |
| Custom auth header | Implemented |
| Chained with other backend auth | Implemented |

**Fixture:** `tests/fixtures/backend_gcp.json`
