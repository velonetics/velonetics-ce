# RabbitMQ Producer (AMQP)

**Edition:** Community + Enterprise  
**Namespace:** `backend/amqp/producer`  
**Official docs:** [AMQP Producer Integration](https://www.krakend.io/docs/enterprise/backends/amqp-producer/)

## What it does

Publishes messages to a RabbitMQ queue when a client calls a KrakenD HTTP endpoint. The request body becomes the message payload.

## Key capabilities

- Dynamic routing keys, message IDs, priority, expiration from URL parameters
- Auto-creates exchange and queue
- Connection retry with backoff

## How it works

```
Client POST /produce/order-123  { "event": "created" }
    │
    ▼
KrakenD publishes body to RabbitMQ exchange/queue
    │
    ▼
HTTP success response to client
```

## Configuration

```json
{
  "endpoint": "/produce/{id}",
  "method": "POST",
  "backend": [{
    "host": ["amqp://guest:guest@localhost:5672"],
    "disable_host_sanitize": true,
    "extra_config": {
      "backend/amqp/producer": {
        "name": "queue-1",
        "exchange": "some-exchange",
        "durable": true,
        "delete": false,
        "no_wait": true,
        "routing_key": "#",
        "mandatory": false,
        "immediate": false,
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
| `exchange` * | — | Exchange name |
| `routing_key` * | `"#"` | Routing key (static or from URL param) |
| `static_routing_key` | `false` | Use fixed `routing_key` vs URL param |
| `mandatory` | `false` | Exchange must have bound queue |
| `immediate` | `false` | Consumer must be connected |
| `msg_id_key` | `""` | URL `{placeholder}` for message ID (capitalized) |
| `priority_key` | `""` | URL param for priority |
| `reply_to_key` | `""` | URL param for reply-to |
| `exp_key` | `""` | URL param for expiration |

### Dynamic parameters example

Endpoint `/produce/{a}/{b}/{id}/{prio}/{route}`:

```json
"backend/amqp/producer": {
  "exp_key": "A",
  "reply_to_key": "B",
  "msg_id_key": "Id",
  "priority_key": "Prio",
  "routing_key": "Route"
}
```

> **Note:** Parameter names in config must have the **first letter uppercased** (`{id}` → `Id`).

## Alternative: Pub/Sub driver

RabbitMQ can also be configured via the generic [Pub/Sub integration](09-pubsub.md) using `rabbit://` host schema.
