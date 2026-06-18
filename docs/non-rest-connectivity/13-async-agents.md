# Event-Driven Async Agents

**Edition:** Community + Enterprise  
**Namespace:** `async_agent` (root-level config array)  
**Official docs:** [Async Agents](https://www.krakend.io/docs/enterprise/async/)

## What it does

Background workers that listen to message queues or Pub/Sub systems and automatically forward events to backends — **no HTTP client request required**. Ideal for saga patterns, event sourcing, and webhook triggers from queues.

## Key capabilities

- Full endpoint feature set: plugins, transformations, sequential/parallel backends, JSON Schema, OAuth2, rate limiting, circuit breaking, Lambda, etc.
- Multiple workers per agent
- Rate limiting per worker
- NACK/requeue on backend failure
- Health endpoint reporting per agent

## Limitations

- No HTTP request middleware (CORS, JWT validation) — events are automatic, not user-initiated

## How it works

```
KrakenD startup
    │
    ▼
Read async_agent[] → spawn workers per agent
    │
    ▼
Worker connects to queue/topic (via driver in extra_config)
    │
    ▼
Message arrives → send body to backend(s) → ACK or NACK
```

### ACK / NACK behavior

| Outcome | Behavior |
|---------|----------|
| Backend succeeds | Message ACKed |
| Backend fails | Message requeued (NACK) for retry |
| `nack_discard: true` | Failed messages discarded |
| `auto_ack: true` | ACK regardless of backend result |

> Warning: If KrakenD is the only consumer and backend keeps failing, NACK can cause infinite reprocessing loops. Use `nack_discard: true` or fix the backend.

## Configuration

```json
{
  "version": 3,
  "async_agent": [{
    "name": "cool-agent",
    "connection": {
      "max_retries": 10,
      "backoff_strategy": "exponential-jitter",
      "health_interval": "10s"
    },
    "consumer": {
      "topic": "*",
      "workers": 1,
      "timeout": "150ms",
      "max_rate": 0.5
    },
    "backend": [{
      "host": ["http://127.0.0.1:8080"],
      "url_pattern": "/webhook"
    }],
    "extra_config": {
      "async/amqp": {
        "host": "amqp://guest:guest@localhost:5672/",
        "name": "krakend",
        "exchange": "foo",
        "durable": true,
        "prefetch_count": 5,
        "auto_ack": false
      }
    }
  }]
}
```

### Agent fields

| Field | Default | Description |
|-------|---------|-------------|
| `name` * | — | Unique agent name (shown in health endpoint) |
| `connection.max_retries` | `0` | Reconnect attempts (`0` = unlimited) |
| `connection.backoff_strategy` | `fallback` | Reconnect delay |
| `connection.health_interval` | `1s` | Health check ping interval |
| `consumer.topic` * | — | Topic/queue pattern (driver-specific syntax) |
| `consumer.workers` | `1` | Parallel consumer processes |
| `consumer.timeout` | `2s` | Max time to process event before NACK |
| `consumer.max_rate` | `0` | Messages/sec per worker (`0` = unlimited) |
| `encoding` | `json` | Response parsing format |
| `backend` * | — | Full backend definition (same as endpoints) |
| `extra_config` * | — | **Must** include `async/amqp` or `async/kafka` driver |

### Drivers

| Driver | Edition | Doc |
|--------|---------|-----|
| `async/amqp` | CE + EE | [AMQP Async Driver](14-async-amqp-driver.md) |
| `async/kafka` | Enterprise | [Kafka Async Driver](15-async-kafka-driver.md) |

## vs endpoint-backed consumers

| | Endpoint AMQP/Kafka consumer | Async Agent |
|---|------------------------------|-------------|
| Trigger | HTTP client request | Automatic on message |
| Use case | Pull message on demand | Continuous processing |
| HTTP middleware | Yes (JWT, CORS) | No |

## Health monitoring

Each agent reports last-alive timestamp in the KrakenD health endpoint, checked at `health_interval`.

## Backoff strategies

Same as WebSockets/AMQP: `linear`, `linear-jitter`, `exponential`, `exponential-jitter`, `fallback`.
