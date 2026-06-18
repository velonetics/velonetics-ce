# gRPC Server

**Edition:** Enterprise only  
**Namespace:** `grpc.server` (service-level `extra_config`)  
**Official docs:** [gRPC Server](https://www.krakend.io/docs/enterprise/grpc/server/)

## What it does

Exposes gRPC services on the **same port** as the HTTP gateway. Each gRPC method maps to one or more backends (gRPC, HTTP, Lambda, queues, etc.). HTTP-specific middleware (CORS, HTTP security) is not applied to gRPC requests.

## Key capabilities

- Publish selected methods from the protobuf catalog
- Aggregate multiple backends per gRPC method (same as HTTP endpoints)
- JWT validation per gRPC method
- gRPC metadata from client headers
- Map gRPC payload fields to URL `{placeholders}` in HTTP backends
- Built-in gRPC reflection (no config needed)

## How it works

```
gRPC client ──► KrakenD gRPC server (same port)
                    │
                    ├──► gRPC backend (flight_finder.Flights/FindFlight)
                    └──► HTTP backend (GET /articles.json?q={cursor})
```

## Configuration

```json
{
  "version": 3,
  "extra_config": {
    "grpc": {
      "catalog": ["./grpc/definitions"],
      "server": {
        "services": [{
          "name": "flight_finder.Flights",
          "methods": [{
            "name": "FindFlight",
            "input_headers": ["*"],
            "payload_params": {
              "page.cursor": "cursor"
            },
            "backend": [
              {
                "host": ["example.com:4242"],
                "url_pattern": "/flight_finder.Flights/FindFlight",
                "extra_config": {
                  "backend/grpc": { "use_request_body": true }
                }
              },
              {
                "method": "GET",
                "host": ["http://example.com:8000"],
                "url_pattern": "/articles.json?q={cursor}"
              }
            ]
          }]
        }]
      }
    }
  }
}
```

### Structure

| Level | Field | Description |
|-------|-------|-------------|
| `grpc` | `catalog` | `.pb` descriptor paths |
| `grpc.server` | `services` | gRPC services to expose |
| `services[]` | `name` | Full service name (e.g. `flight_finder.Flights`) |
| `services[].methods[]` | `name` | Method to publish |
| `methods[]` | `backend` | Same backend config as HTTP endpoints |
| `methods[]` | `input_headers` | Client headers → gRPC metadata |
| `methods[]` | `payload_params` | Map gRPC field (dot notation) → `{placeholder}` |
| `methods[]` | `extra_config` | Per-method middleware (e.g. JWT) |

### JWT authorization

Place `auth/validator` under the method's `extra_config`:

```json
{
  "name": "FindFlight",
  "extra_config": {
    "auth/validator": {
      "alg": "RS256",
      "audience": ["http://api.example.com"],
      "roles_key": "http://api.example.com/custom/roles",
      "roles": ["user", "admin"],
      "jwk_url": "https://identity.example.com/.well-known/jwks.json",
      "cache": true
    }
  },
  "backend": [...]
}
```

### gRPC reflection

Enabled automatically. Tools like `grpcurl` can discover services without local `.proto` files.

### OpenTelemetry overrides

Under `grpc.server`:

| Field | Description |
|-------|-------------|
| `opentelemetry.disable_metrics` | Disable gRPC server metrics |
| `opentelemetry.disable_traces` | Disable gRPC server traces |
