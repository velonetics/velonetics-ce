# gRPC Overview & Service Catalog

**Edition:** Enterprise only  
**Namespace:** `grpc` (service-level `extra_config`)  
**Official docs:** [Introduction to gRPC](https://www.krakend.io/docs/enterprise/grpc/)

## What it does

gRPC integration lets KrakenD act as a **gRPC client** (consume upstream gRPC services) and/or a **gRPC server** (expose gRPC to consumers). Both roles share a common **protocol buffer catalog**.

## Key capabilities

- Unary RPC support (no streaming — server, client, or bidirectional)
- REST ↔ gRPC transparent conversion
- gRPC ↔ gRPC passthrough through the gateway
- Mixed upstream backends (HTTP + gRPC + Lambda + queues) behind gRPC server methods
- Automatic gRPC reflection for debugging (server mode)

## Use cases

| Scenario | Mode |
|----------|------|
| Hide gRPC complexity from REST clients | gRPC client → REST response |
| Offer gRPC when backends are HTTP-only | gRPC server → HTTP backends |
| Full gRPC pipeline | gRPC server + gRPC client |

## How it works

```
                    ┌─────────────────────────────────┐
  gRPC/REST client  │           KrakenD               │
        ──────────► │  catalog (.pb) + grpc config    │ ──► gRPC upstream
                    │                                 │ ──► HTTP upstream
                    └─────────────────────────────────┘
```

1. Define `.proto` services and compile to binary `.pb` descriptor files.
2. Register catalog paths at service level under `extra_config.grpc.catalog`.
3. Configure either `backend/grpc` (client) or `grpc.server` (server).

## Catalog setup

KrakenD uses **binary `.pb` files**, not raw `.proto` files.

### Generate `.pb` from `.proto`

```bash
protoc --descriptor_set_out=file.pb file.proto
```

### Single catalog file

```bash
mkdir -p ./defs
cd contracts && \
  protoc --descriptor_set_out=../fullcatalog.pb $(find . -name '*.proto')
```

### Service-level configuration

```json
{
  "version": 3,
  "extra_config": {
    "grpc": {
      "catalog": [
        "./grpc/flights.pb",
        "./grpc/definitions",
        "/etc/krakend/grpc"
      ]
    }
  }
}
```

| Field | Description |
|-------|-------------|
| `catalog` | Array of `.pb` file paths or directories to scan |

### Dependencies & imports

- If `.proto` files import other definitions, their `.pb` counterparts must also be in the catalog.
- Missing definitions produce **WARNING** logs; data fields will be empty (not a hard failure).
- Well-known types `timestamp.proto` and `duration.proto` are auto-converted to JSON representations.

### Supported proto versions

Both **proto2** and **proto3** are supported.

## Next steps

- [gRPC Client](04-grpc-client.md) — consume upstream gRPC services
- [gRPC Server](05-grpc-server.md) — expose gRPC methods on the gateway port
