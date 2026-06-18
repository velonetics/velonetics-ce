# gRPC Client

**Edition:** Enterprise only  
**Namespace:** `backend/grpc` (backend `extra_config`) + `grpc.catalog` (service level)  
**Official docs:** [gRPC client and gRPC to REST conversion](https://www.krakend.io/docs/enterprise/backends/grpc/)

## What it does

Connects KrakenD to upstream gRPC services. When no gRPC server is enabled, responses are automatically transformed to REST/JSON for end users. With gRPC server enabled, supports gRPC-to-gRPC passthrough.

## Key capabilities

- Unary RPC calls to external gRPC services
- Query strings, URL placeholders, headers, and body → gRPC message fields
- Header → gRPC metadata mapping
- TLS client configuration
- Full response manipulation (same as REST backends)
- Alternate host retry on transient failures

## How it works

```
Client GET /flights?lat=123&lon=456
    │
    ▼
KrakenD maps params → gRPC message (via input_mapping)
    │
    ▼
Unary RPC: flight_finder.Flights/FindFlight @ localhost:4242
    │
    ▼
Protobuf response → JSON (snake_case or camelCase) → client
```

## Configuration

### Prerequisites

1. Service-level `grpc.catalog` with `.pb` files (see [gRPC Overview](03-grpc-overview.md))
2. Backend `extra_config.backend/grpc`

### Minimal example

```json
{
  "version": 3,
  "extra_config": {
    "grpc": {
      "catalog": ["grpcatalog/flights/flights.pb"]
    }
  },
  "endpoints": [{
    "endpoint": "/flights",
    "input_query_strings": ["*"],
    "backend": [{
      "host": ["localhost:4242"],
      "url_pattern": "/flight_finder.Flights/FindFlight",
      "extra_config": { "backend/grpc": {} }
    }]
  }]
}
```

### Important rules

| Setting | Rule |
|---------|------|
| `host` | `host:port` only — **no** `http://` or `https://` prefix |
| `url_pattern` | Full gRPC method: `/package.Service/Method` — no variables |
| `output_encoding` | Must **not** be `no-op` |
| `method`, `is_collection` | Ignored for gRPC backends |

### Key `backend/grpc` fields

| Field | Default | Description |
|-------|---------|-------------|
| `use_request_body` | `false` | Fill gRPC request from HTTP body (consumed by first backend) |
| `disable_query_params` | `false` | Ignore URL params and query strings |
| `input_mapping` | — | Rename query/placeholder params to nested gRPC fields |
| `header_mapping` | — | Rename HTTP headers to gRPC metadata |
| `client_tls` | — | TLS settings for gRPC connection |
| `request_naming_convention` | `snake_case` | `snake_case` or `camelCase` for request |
| `response_naming_convention` | `snake_case` | `snake_case` or `camelCase` for response |
| `output_timestamp_as_string` | `false` | RFC3339 string for Timestamp types |
| `output_enum_as_string` | `false` | String representation of enum values |
| `max_call_recv_msg_size` | 4MB | Max receive message size |

### Passing parameters

| Source | How |
|--------|-----|
| Query strings | Declare in `input_query_strings`; use `input_mapping` for nested fields |
| URL `{placeholders}` | First letter capitalized: `{date}` → `Date` |
| Headers | Declare in `input_headers`; map via `header_mapping` |
| Body | Set `use_request_body: true` |

**Precedence:** params/query first, then body (body overwrites collisions).

### input_mapping example

Endpoint `/flights/{date}` with `?lat=123&lon=456`:

```json
"backend/grpc": {
  "input_mapping": {
    "lat": "where.latitude",
    "lon": "where.longitude",
    "Date": "when.departure"
  }
}
```

gRPC receives:
```json
{
  "where": { "latitude": 123, "longitude": 456 },
  "when": { "departure": "2023-07-09" }
}
```

### Limitations

- Dot notation for nested fields: `some_field.child.grand_child=10`
- Repeated basic types: `a=1&a=2&a=3`
- **Cannot** fill arrays of objects or nested arrays via query strings
