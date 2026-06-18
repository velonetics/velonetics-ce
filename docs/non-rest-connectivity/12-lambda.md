# AWS Lambda Integration

**Edition:** Community + Enterprise  
**Namespace:** `backend/lambda`  
**Official docs:** [AWS Lambda Integration](https://www.krakend.io/docs/enterprise/backends/lambda/)

## What it does

Invokes Amazon Lambda functions directly from KrakenD endpoints — no API Gateway required. Lambda responses are processed like any other backend (manipulation, filtering, etc.).

## Key capabilities

- Direct Lambda invocation (bypasses API Gateway cost)
- Static or dynamic function names from URL parameters
- Configurable retry on invocation failure
- LocalStack support via custom endpoint
- Canary/A/B testing with Lua scripts

## How it works

```
Client HTTP request
    │
    ▼
KrakenD builds Lambda payload from request
    │
    ▼
AWS Lambda Invoke API
    │
    ▼
Lambda response → manipulation → client
```

### Payload rules

| Endpoint method | Lambda payload |
|-----------------|----------------|
| GET | All request parameters |
| POST/PUT/PATCH/etc. | Request body content |

## Configuration

### Fixed function name

```json
{
  "endpoint": "/call-a-lambda",
  "backend": [{
    "host": ["ignore"],
    "url_pattern": "/ignore",
    "extra_config": {
      "backend/lambda": {
        "function_name": "myLambdaFunction",
        "region": "us-east-1",
        "max_retries": 1
      }
    }
  }]
}
```

> `host` and `url_pattern` are required by schema but **never used** for Lambda.

### Dynamic function from URL

```json
{
  "endpoint": "/call-a-lambda/{lambda}",
  "backend": [{
    "host": ["ignore"],
    "url_pattern": "/ignore",
    "extra_config": {
      "backend/lambda": {
        "function_param_name": "Lambda",
        "region": "us-west-1",
        "max_retries": 1
      }
    }
  }]
}
```

`GET /call-a-lambda/my-function` invokes `my-function`.

### Key fields

| Field | Description |
|-------|-------------|
| `function_name` | Static Lambda function name |
| `function_param_name` | URL `{placeholder}` for function name (capitalized) |
| `region` | AWS region (e.g. `us-east-1`) |
| `max_retries` | Retry count (`-1` = defer to service config) |
| `endpoint` | Custom endpoint (LocalStack testing) |

### Function name formats

- Name only: `my-function`
- Name + version: `my-function:v1`
- Full ARN: `arn:aws:lambda:us-west-2:123456789012:function:my-function`
- Partial ARN: `123456789012:function:my-function`

## Authentication

KrakenD needs AWS credentials via one of:

1. `~/.aws/credentials` (mount in Docker)
2. Environment variables: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_REGION`
3. IAM role on the host (EC2/ECS/EKS)

### Docker example

```bash
docker run --rm -it -p 8080:8080 \
  -e AWS_ACCESS_KEY_ID=XXX \
  -e AWS_SECRET_ACCESS_KEY=XXX \
  -e AWS_REGION=eu-west-1 \
  -v "$PWD:/etc/krakend" krakend/krakend-ee:2.13.5
```

## Limitations

- **Header forwarding to Lambda is not supported** by the AWS SDK. Embed headers in a custom payload via Lua if needed.

## Canary testing example

Use Lua to route 20% of traffic to a new Lambda version:

```json
"extra_config": {
  "backend/lambda": {
    "function_param_name": "Function_name",
    "region": "eu-west-1"
  },
  "modifier/lua-backend": {
    "sources": ["canary.lua"],
    "pre": "canaryLambda(request.load())"
  }
}
```

```lua
function canaryLambda(req)
  if math.random(0, 100) < 20 then
    req:params("Function_name", "my-function:3")
  else
    req:params("Function_name", "my-function:2")
  end
end
```
