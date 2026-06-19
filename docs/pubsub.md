# Publisher/Subscribe (Kafka, NATS, Cloud)

Pucora connects HTTP endpoints to publish/subscribe messaging systems using the `backend/pubsub/publisher` and `backend/pubsub/subscriber` namespaces.

## What it does

- **Publisher:** `POST /events` publishes the HTTP body to a topic or queue.
- **Subscriber:** `GET /events` pulls one message and returns it as the HTTP response.

All standard endpoint middleware (validation, rate limiting, manipulation) still applies.

## Supported drivers

| Driver | `host` scheme | Environment variable |
|--------|---------------|----------------------|
| Apache Kafka | `kafka://` | `KAFKA_BROKERS` |
| NATS.io | `nats://` | `NATS_SERVER_URL` |
| RabbitMQ | `rabbit://` | `RABBIT_SERVER_URL` |
| GCP Pub/Sub | `gcppubsub://` | `GOOGLE_APPLICATION_CREDENTIALS` / `PUBSUB_EMULATOR_HOST` |
| AWS SNS | `awssns:///` + ARN in host | AWS credentials |
| AWS SQS | `awssqs://` + queue URL | AWS credentials |
| Azure Service Bus | `azuresb://` | `SERVICEBUS_CONNECTION_STRING` |

## Configuration

Set the broker scheme in `host[0]`, set `disable_host_sanitize: true`, and put the topic or subscription path in `extra_config`. `url_pattern` is required by schema but unused.

### Publisher

```json
{
  "host": ["kafka://"],
  "url_pattern": "/ignored",
  "disable_host_sanitize": true,
  "extra_config": {
    "backend/pubsub/publisher": {
      "topic_url": "mytopic"
    }
  }
}
```

### Subscriber

```json
{
  "host": ["kafka://"],
  "url_pattern": "/ignored",
  "disable_host_sanitize": true,
  "encoding": "json",
  "extra_config": {
    "backend/pubsub/subscriber": {
      "subscription_url": "group?topic=mytopic"
    }
  }
}
```

## Examples

Docker Compose stacks for each driver live under [`examples/pubsub/`](../examples/pubsub/README.md).

```bash
make pubsub-nats-compose-test
make pubsub-kafka-compose-test
```

## Kafka with mTLS/SASL

For secured Kafka with explicit cluster settings, use [Kafka Advanced Pub/Sub](pubsub-kafka-advanced.md) (`backend/pubsub/publisher/kafka` and `backend/pubsub/subscriber/kafka`).

## Related

- [Kafka Advanced Pub/Sub](pubsub-kafka-advanced.md)
- [Async Kafka agent](async-kafka.md)
- For full RabbitMQ control (prefetch, ack, routing keys), use `backend/amqp/consumer` and `backend/amqp/producer`.
