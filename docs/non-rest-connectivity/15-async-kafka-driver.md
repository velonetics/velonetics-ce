# Kafka Driver for Async Agents

**Edition:** Enterprise only  
**Namespace:** `async/kafka` (inside `async_agent[].extra_config`)  
**Official docs:** [Kafka Async Driver](https://www.krakend.io/docs/enterprise/async/kafka/)

## What it does

Connects an [Async Agent](13-async-agents.md) to Kafka topics for autonomous event consumption. Unlike [endpoint Kafka subscribers](10-kafka-advanced-pubsub.md), no HTTP request triggers processing.

## Key capabilities

- mTLS and SASL authentication (shared config with Kafka Advanced PubSub)
- Consumer group management
- Message key exposed as HTTP header (`key_meta`)
- Configurable isolation level
- Latest-offset policy (only messages after startup)

## How it works

```
Kafka topic ──message──► async/kafka driver ──body──► backend webhook(s)
                                │
                          offset committed (always)
```

Topics are defined in the agent's `consumer.topic` field (not in the driver config).

KrakenD follows the **latest offset policy** — only messages produced after agent startup are processed.

## Configuration

```json
{
  "version": 3,
  "async_agent": [{
    "name": "stock-consumer",
    "consumer": {
      "topic": "stock-prices",
      "workers": 2,
      "timeout": "2s"
    },
    "backend": [{
      "host": ["http://analytics:8080"],
      "url_pattern": "/ingest"
    }],
    "extra_config": {
      "async/kafka": {
        "cluster": {
          "brokers": ["localhost:9092"],
          "client_id": "cid_stocksconsumer",
          "client_tls": {
            "ca_certs": ["./config/ca/cacert.pem"],
            "client_certs": [{
              "certificate": "./config/certs/client.pem",
              "private_key": "./config/certs/client.key"
            }]
          },
          "sasl": {
            "user": "johnsmith",
            "password": "myp4ssword"
          }
        },
        "group": {
          "group_id": "my_group_id",
          "isolation_level": "read_commited"
        },
        "key_meta": "Message-Id"
      }
    }
  }]
}
```

### Driver fields

| Field | Description |
|-------|-------------|
| `cluster` * | Broker connection (brokers, TLS, SASL, timeouts) |
| `group` | Consumer group settings |
| `key_meta` | HTTP header name for Kafka message key |

> Driver config matches `reader` in `backend/pubsub/subscriber/kafka`, except `topics` comes from `consumer.topic` in the agent.

### Cluster fields

| Field | Default | Description |
|-------|---------|-------------|
| `brokers` * | — | Broker addresses |
| `client_id` | `KrakenD v[X].[Y].[Z]` | Client identifier |
| `client_tls` | — | mTLS configuration |
| `sasl` | — | `{ user, password, mechanism }` |
| `metadata_retry_max` | `3` | Reconnect attempts |
| `metadata_retry_backoff` | `250ms` | Reconnect delay |

### SASL mechanisms

- `PLAIN` (default)
- `OAUTHBEARER`
- Azure Event Hub: set `azure_event_hub: true` for SASL V0

### Consumer group fields

| Field | Default | Description |
|-------|---------|-------------|
| `group_id` | — | Consumer group ID |
| `isolation_level` | `read_commited` | `read_commited` or `read_uncommited` |
| `session_timeout` | `10s` | Failure detection |
| `heartbeat_interval` | `3s` | Coordinator heartbeat |
| `rebalance_strategies` | `["range"]` | `range`, `roundrobin`, `sticky` |

## ACK behavior

Kafka has **no per-message ACK**. KrakenD commits partition offsets for every message read, **even when the backend pipeline fails**. Design backends for idempotency.

## Reconnection

Default `metadata_retry_max=3` × `250ms` = ~750ms recovery window. For production:

```json
"cluster": {
  "metadata_retry_max": 30,
  "metadata_retry_backoff": "1s"
}
```

## Related

- [Kafka Advanced PubSub](10-kafka-advanced-pubsub.md) — HTTP-triggered Kafka publish/subscribe with same `cluster` config
- [Async Agents overview](13-async-agents.md) — agent lifecycle and worker configuration
