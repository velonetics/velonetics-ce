# WebSockets Integration

**Edition:** Enterprise only  
**Namespace:** `websocket` (endpoint `extra_config`)  
**Official docs:** [WebSockets Integration](https://www.krakend.io/docs/enterprise/websockets/)

## What it does

Enables bidirectional real-time communication using the [WebSocket Protocol (RFC-6455)](https://datatracker.ietf.org/doc/html/rfc6455). Clients connect to KrakenD over WebSocket; KrakenD connects to backend hosts via `ws://` or `wss://`.

## Key capabilities

- **Multiplexing** (default): One backend WebSocket shared across many clients
- **Direct communication**: One backend connection per client (transparent proxy)
- Connect/disconnect event notifications to backend
- JWT validation on HTTP upgrade request
- Automatic reconnection with configurable backoff
- Subprotocol support (direct mode only)

## Operating modes

| Mode | Config | Client ↔ Gateway | Gateway ↔ Backend | Best for |
|------|--------|------------------|-------------------|----------|
| **Multiplexing** | `enable_direct_communication: false` | 1 WS per client | **1 shared** WS per endpoint | Chat rooms, fan-out, many clients |
| **Direct** | `enable_direct_communication: true` | 1 WS per client | **1 WS per client** | Binary streams, subprotocols, Socket.IO |

## How multiplexing works

```
1000 clients ──WS──► KrakenD ──1 WS──► Backend
```

KrakenD wraps messages in a JSON envelope for the backend:

```json
{
  "url": "/chat/krakend",
  "session": { "uuid": "0b251b07-...", "Room": "krakend" },
  "body": "SGVsbG8gV29ybGQh"
}
```

- `body` is base64-encoded client payload
- `session.uuid` is assigned by KrakenD per client session
- URL `{placeholders}` appear in `session` with first letter uppercased (e.g. `Room`)

### Backend → client routing

| Filter | Behavior |
|--------|----------|
| No `url` or `session` | **Broadcast** to all clients |
| `"url": "/chat/krakend"` | Multicast to endpoint |
| `"session": { "uuid": "..." }` | Unicast to one client |

### Handshake (multiplexing only)

Backend must reply `OK` to KrakenD's opening message `{"msg":"KrakenD WS proxy starting"}`.

## Configuration

### Requirements

1. `extra_config.websocket` present (empty `{}` is valid)
2. Backend `host` uses `ws://` or `wss://`
3. `disable_host_sanitize: true` on backend
4. Endpoint `method: "GET"` (upgrade uses GET)

### Example (multiplexing)

```json
{
  "endpoint": "/ws/{room}",
  "input_query_strings": ["*"],
  "backend": [{
    "url_pattern": "/ws",
    "disable_host_sanitize": true,
    "host": ["ws://localhost:8081", "ws://localhost:8082"]
  }],
  "extra_config": {
    "websocket": {
      "input_headers": ["Cookie", "Authorization"],
      "connect_event": true,
      "disconnect_event": true,
      "read_buffer_size": 4096,
      "write_buffer_size": 4096,
      "message_buffer_size": 4096,
      "max_message_size": 3200000,
      "write_wait": "10s",
      "pong_wait": "60s",
      "ping_period": "54s",
      "max_retries": 0,
      "backoff_strategy": "exponential"
    }
  }
}
```

### Key `websocket` fields

| Field | Default | Description |
|-------|---------|-------------|
| `enable_direct_communication` | `false` | Disable multiplexing; 1:1 proxy |
| `input_headers` | `[]` | Headers forwarded to backend |
| `max_message_size` | `512` | Max message bytes; oversized = disconnect |
| `max_retries` | `0` | Reconnect attempts (`0` = unlimited) |
| `backoff_strategy` | `fallback` | Reconnect delay strategy |
| `connect_event` / `disconnect_event` | `false` | Notify backend on client connect/disconnect |
| `subprotocols` | `[]` | Allowed subprotocols (direct mode) |
| `return_error_details` | `false` | Send `{"error":"..."}` to client on backend failure |

## Socket.IO note

Socket.IO is **not** plain WebSockets. Use **direct mode** only, with `url_pattern`: `/socket.io/?EIO=4&transport=websocket` and clients restricted to `websocket` transport.

## Log levels

| Level | Meaning |
|-------|---------|
| `WARNING` | Backend connectivity issues |
| `ERROR` | Failed reconnection attempts |
| `CRITICAL` | WebSocket connection lost permanently |
