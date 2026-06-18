# RabbitMQ Consumer (AMQP)

**Edition:** Community + Enterprise  
**Namespace:** `backend/amqp/consumer`  
**Official docs:** [AMQP Consumer Integration](https://www.krakend.io/docs/enterprise/backends/amqp-consumer/)

## What it does

Pulls messages from a RabbitMQ queue when a client calls a KrakenD HTTP endpoint. The message body becomes the backend response (or input for further processing).

## Key capabilities

- Consume from N queues via N backends on one endpoint
- Auto-creates exchange and queue
- Prefetch control, auto-ack, NACK/requeue on failure
- Connection retry with backoff strategies

## How it works

```
Client GET /events
    │
    ▼
KrakenD pulls message from RabbitMQ queue
    │
    ▼
Message body returned to client (or forwarded to additional backends)
```

> For **autonomous** consumption without HTTP requests, use [Async Agents](13-async-agents.md) with the [AMQP driver](14-async-amqp-driver.md).

## Configuration

```json
{
  "endpoint": "/events",
  "method": "GET",
  "backend": [{
    "host": ["amqp://guest:guest@localhost:5672"],
    "disable_host_sanitize": true,
    "extra_config": {
      "backend/amqp/consumer": {
        "name": "queue-1",
        "exchange": "some-exchange",
        "durable": true,
        "delete": false,
        "no_wait": true,
        "no_local": false,
        "routing_key": ["#"],
        "prefetch_count": 10,
        "auto_ack": false,
        "backoff_strategy": "exponential-jitter"
      }
    }
  }]
}
```

### Key fields

| Field | Default | Description |
|-------|---------|-------------|
| `name` * | — | Queue name |
| `exchange` * | — | Exchange name (topic type if exists) |
| `routing_key` * | — | Array of routing keys (e.g. `["#"]`) |
| `durable` | `false` | Queue survives broker restart |
| `delete` | `false` | Delete queue when no connections (not recommended) |
| `exclusive` | `false` | Single KrakenD instance only |
| `prefetch_count` | `0` | Messages to prefetch before consuming |
| `auto_ack` | `false` | ACK regardless of processing success |
| `nack_discard` | `false` | Discard failed messages instead of requeue |
| `max_retries` | `0` | Reconnect attempts (`0` = unlimited) |
| `backoff_strategy` | `fallback` | Reconnect delay strategy |

### Recommendations

- Connect consumers to **GET** endpoints
- Set `durable: true` for production queues
- Avoid `delete: true` — queue may be deleted during connectivity blips
- Avoid `exclusive: true` in HA multi-instance deployments

## Connection retries

KrakenD auto-reconnects on AMQP connection loss. Configure `backoff_strategy` and `max_retries` to match your SLA.
