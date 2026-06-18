# Kafka Advanced PubSub

**Edition:** Enterprise only  
**Namespace:** `backend/pubsub/publisher/kafka` and `backend/pubsub/subscriber/kafka`  
**Official docs:** [Kafka PubSub with Extended Connectivity](https://www.krakend.io/docs/enterprise/backends/pubsub/kafka/)

## What it does

Enterprise Kafka integration with **mutual TLS (mTLS)** and **SASL authentication**. Use this instead of the basic `kafka://` Pub/Sub driver when you need secured broker connections.

## Key capabilities

- mTLS with CA and client certificates
- SASL PLAIN and OAUTHBEARER authentication
- Idempotent producer
- Consumer group configuration with isolation levels
- Message key from HTTP headers (`key_meta`)
- Shared `cluster` config with Kafka async agent driver

## How it works

Same request-driven pattern as [basic Pub/Sub](09-pubsub.md), but with explicit broker connection settings instead of environment variables only.

```
Client HTTP request  →  KrakenD Kafka writer/reader  →  Kafka cluster (mTLS/SASL)
```

## Publisher configuration

```json
{
  "backend": [{
    "host": ["ignore"],
    "url_pattern": "/ignore",
    "extra_config": {
      "backend/pubsub/publisher/kafka": {
        "success_status_code": 201,
        "writer": {
          "topic": "orderplacement",
          "key_meta": "X-My-Key",
          "cluster": {
            "brokers": ["localhost:49092"],
            "client_id": "krakend_publisher",
            "client_tls": {
              "allow_insecure_connections": false,
              "ca_certs": ["./config/ca/cacert.pem"],
              "client_certs": [{
                "certificate": "./config/certs/client/client.signed.pem",
                "private_key": "./config/certs/client/client.key"
              }]
            }
          },
          "producer": {
            "idempotent": true
          }
        }
      }
    }
  }]
}
```

### Publisher fields

| Field | Description |
|-------|-------------|
| `success_status_code` | HTTP status on successful publish (default: `200`) |
| `writer.topic` * | Target Kafka topic |
| `writer.key_meta` | HTTP header name for message key |
| `writer.cluster` * | Broker connection settings |
| `writer.producer` | Producer behavior (e.g. `idempotent`) |

## Subscriber configuration

```json
{
  "backend": [{
    "host": ["ignore"],
    "url_pattern": "/ignore",
    "extra_config": {
      "backend/pubsub/subscriber/kafka": {
        "reader": {
          "topics": ["orderplacement"],
          "key_meta": "Message-Id",
          "cluster": {
            "brokers": ["localhost:49092"],
            "client_id": "krakend_subscriber",
            "client_tls": { "...": "..." }
          },
          "group": {
            "group_id": "my_group_id",
            "isolation_level": "read_commited"
          }
        }
      }
    }
  }]
}
```

### Cluster fields (shared)

| Field | Default | Description |
|-------|---------|-------------|
| `brokers` * | — | Kafka broker addresses |
| `client_id` | `KrakenD v[X].[Y].[Z]` | Client identifier |
| `client_tls` | — | mTLS settings |
| `sasl` | — | `{ user, password, mechanism }` |
| `dial_timeout` | `30s` | Connection dial timeout |
| `metadata_retry_max` | `3` | Reconnect attempts |
| `metadata_retry_backoff` | `250ms` | Delay between reconnects |

### SASL example

```json
"cluster": {
  "brokers": ["localhost:9092"],
  "sasl": {
    "user": "johnsmith",
    "password": "myp4ssword"
  }
}
```

### Consumer group fields

| Field | Default | Description |
|-------|---------|-------------|
| `group_id` | — | Consumer group ID |
| `isolation_level` | `read_commited` | `read_commited` or `read_uncommited` |
| `session_timeout` | `10s` | Failure detection timeout |
| `heartbeat_interval` | `3s` | Coordinator heartbeat interval |

## ACK behavior

Kafka has **no per-message ACK**. KrakenD commits partition offsets for every message read, even when the backend pipeline returns a non-success status.

## Reconnection tip

Default `metadata_retry_max=3` with `250ms` backoff only covers ~750ms of outage. For production, increase retries and set backoff to `1s` or higher.

## Related

- [Kafka Async Driver](15-async-kafka-driver.md) — autonomous consumption without HTTP
- [Basic Pub/Sub](09-pubsub.md) — simpler Kafka via `kafka://` + `KAFKA_BROKERS`
