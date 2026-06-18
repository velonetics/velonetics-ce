# Publisher/Subscribe (Kafka, NATS, Cloud)

**Edition:** Community + Enterprise  
**Namespace:** `backend/pubsub/subscriber` and `backend/pubsub/publisher`  
**Official docs:** [Publisher/Subscribe](https://www.krakend.io/docs/enterprise/backends/pubsub/)

## What it does

Connects HTTP endpoints to publish/subscribe messaging systems. Clients can push events via REST or consume the latest events from a topic/queue through a REST call.

## Supported backends

| Driver | `host` schema | Env variable |
|--------|---------------|--------------|
| Apache Kafka | `kafka://` | `KAFKA_BROKERS` |
| NATS.io | `nats://` | `NATS_SERVER_URL` |
| RabbitMQ | `rabbit://` | `RABBIT_SERVER_URL` |
| GCP PubSub | `gcppubsub://` | `GOOGLE_APPLICATION_CREDENTIALS` |
| AWS SNS | `awssns:///` + ARN in host | AWS credentials |
| AWS SQS | `awssqs://` + queue URL | AWS credentials |
| Azure Service Bus | `azuresb://` | `SERVICEBUS_CONNECTION_STRING` |

## How it works

### Publisher (HTTP → message broker)

```
Client POST /events  →  KrakenD  →  publish to topic/queue
```

### Subscriber (message broker → HTTP response)

```
Client GET /events  →  KrakenD pulls from subscription  →  return message to client
```

All standard KrakenD middleware (validation, rate limiting, manipulation) applies.

## Configuration

### Subscriber

```json
{
  "host": ["kafka://"],
  "url_pattern": "/ignored",
  "disable_host_sanitize": true,
  "extra_config": {
    "backend/pubsub/subscriber": {
      "subscription_url": "group?topic=mytopic"
    }
  }
}
```

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

> `url_pattern` is required by schema but unused — set any value. `disable_host_sanitize: true` is required.

## Driver-specific examples

### Kafka

Env: `KAFKA_BROKERS=192.168.1.100:9092`

```json
// Subscribe
"subscription_url": "group?topic=mytopic"
// Publish
"topic_url": "mytopic"
```

### NATS

Env: `NATS_SERVER_URL=nats://localhost:4222`

```json
"subscription_url": "mysubject"
```

### GCP PubSub

```json
"host": ["gcppubsub://"],
"subscription_url": "projects/myproject/subscriptions/mysub"
// or shortened: "myproject/mysub"
```

### AWS SNS

```json
"host": ["awssns:///arn:aws:sns:us-east-2:123456789012:mytopic"],
"subscription_url": "?region=us-east-2"
```

### AWS SQS

```json
"host": ["awssqs://sqs.us-east-2.amazonaws.com/123456789012"],
"subscription_url": "/myqueue?region=us-east-2"
```

### Azure Service Bus

Env: `SERVICEBUS_CONNECTION_STRING`

```json
"host": ["azuresb://"],
"subscription_url": "mytopic"
// Subscriptions: "mytopic?subscription=mysubscription"
```

### RabbitMQ (Pub/Sub style)

Env: `RABBIT_SERVER_URL='guest:guest@localhost:5672'`

```json
"host": ["rabbit://"],
"subscription_url": "myexchange"
```

> For full AMQP control (prefetch, ack, routing keys), use [RabbitMQ Consumer/Producer](07-rabbitmq-consumer.md) instead.

## Kafka with mTLS/SASL

For secured Kafka connections, use [Kafka Advanced PubSub](10-kafka-advanced-pubsub.md) (`backend/pubsub/*/kafka` namespaces).
