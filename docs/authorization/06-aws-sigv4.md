# AWS SigV4 Authentication

**Edition:** Enterprise (implemented in Pucora CE)  
**Namespace:** `auth/aws-sigv4`  
**Module:** `pucora-aws-sigv4`  
**Official docs:** [AWS SigV4 Authentication](../../../docs/authorization/14-aws-sigv4-authentication.md)

## What it does

Signs outbound backend HTTP requests with AWS Signature Version 4 using the default credential chain (env, shared config, IAM role).

## Configuration

```json
{
  "backend": [{
    "host": ["https://abc123.execute-api.us-east-1.amazonaws.com"],
    "extra_config": {
      "auth/aws-sigv4": {
        "service": "execute-api",
        "region": "us-east-1",
        "assume_role_arn": "arn:aws:iam::123456789012:role/GatewayRole"
      }
    }
  }]
}
```

## Parity status

| Area | Status |
|------|--------|
| SigV4 request signing | Implemented |
| Default credential chain | Implemented |
| assume_role_arn (STS) | Implemented |
| Debug logging | Implemented |
| Chained with other backend auth | Implemented |

**Fixture:** `tests/fixtures/backend_sigv4.json`
