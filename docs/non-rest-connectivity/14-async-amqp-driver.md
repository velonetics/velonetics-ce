# AMQP Driver for Async Agents

**Edition:** Community + Enterprise  
**Namespace:** `async/amqp` (inside `async_agent[].extra_config`)  
**Official docs:** [AMQP driver for Async Agent](https://www.krakend.io/docs/async/amqp/)

## What it does

Connects an [Async Agent](13-async-agents.md) to a RabbitMQ (AMQP) queue for autonomous message consumption. Unlike [endpoint AMQP consumers](07-rabbitmq-consumer.md), no HTTP request triggers processing.

## Key capabilities

- Autonomous queue consumption on KrakenD startup
- Prefetch control
- ACK on success / NACK on failure
- `nack_discard` to prevent infinite reprocessing loops
- Durable queue support

## How it works

```
RabbitMQ queue ──message──► async/amqp driver ──body──► backend webhook(s)
                                    │
                              ACK or NACK
```

The `consumer.topic` field in the agent config defines which messages to consume (AMQP routing patterns like `*`, `lazy.#`, `foo.*`).

## Configuration

Place inside `async_agent[].extra_config`:

```json
{
  "async/amqp": {
    "host": "amqp://guest:guest@localhost:5672/",
    "name": "krakend",
    "exchange": "foo",
    "durable": true,
    "delete": false,
    "exclusive": false,
    "no_wait": true,
    "prefetch_count": 5,
    "auto_ack": false,
    "no_local": true,
    "nack_discard": false
  }
}
```

### Fields

| Field | Default | Description |
|-------|---------|-------------|
| `host` * | — | AMQP connection string ending in `/` |
| `name` * | — | Queue name |
| `exchange` * | — | Exchange name (created or must be topic type) |
| `durable` | `false` | Queue survives broker restart |
| `delete` | `false` | Delete queue when no connections |
| `exclusive` | `false` | Single KrakenD instance only |
| `no_wait` | — | Don't wait for server confirm |
| `prefetch_count` | `10` | Messages to prefetch |
| `prefetch_size` | `0` | Bytes to prefetch |
| `auto_ack` | `false` | ACK regardless of backend success |
| `nack_discard` | `false` | Discard failed messages instead of requeue |
| `no_local` | — | Not supported by RabbitMQ |

## Full agent example

```json
{
  "version": 3,
  "async_agent": [{
    "name": "order-processor",
    "connection": {
      "max_retries": 10,
      "backoff_strategy": "exponential-jitter"
    },
    "consumer": {
      "topic": "orders.*",
      "workers": 3,
      "timeout": "5s",
      "max_rate": 0
    },
    "backend": [{
      "host": ["http://internal-api:8080"],
      "url_pattern": "/process-order"
    }],
    "extra_config": {
      "async/amqp": {
        "host": "amqp://user:pass@rabbitmq:5672/",
        "name": "order-queue",
        "exchange": "orders",
        "durable": true,
        "prefetch_count": 10,
        "auto_ack": false,
        "nack_discard": true
      }
    }
  }]
}
```

## vs endpoint AMQP consumer

| Feature | `backend/amqp/consumer` | `async/amqp` |
|---------|-------------------------|--------------|
| Trigger | HTTP GET/POST | Automatic |
| Multiple workers | Via multiple backends | Via `consumer.workers` |
| Rate limiting | Endpoint rate limit | `consumer.max_rate` |
| Best for | On-demand pull | Continuous processing |
