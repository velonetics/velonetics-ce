# KrakenD Non-REST Connectivity

Reference documentation for all **15 non-REST connectivity features** in [KrakenD Enterprise](https://www.krakend.io/docs/enterprise/non-rest-connectivity/). These features let the API gateway integrate with message brokers, event-driven systems, and alternative protocols beyond traditional REST/HTTP.

> **Source:** [KrakenD Enterprise — Non-REST Connectivity](https://www.krakend.io/docs/enterprise/non-rest-connectivity/)

## How it fits together

KrakenD exposes non-REST systems through two main patterns:

| Pattern | Trigger | Example |
|---------|---------|---------|
| **Endpoint-backed** | Client HTTP request | `POST /subscribe` → publish to RabbitMQ |
| **Async agent** | Message arrives on queue/topic | Kafka event → call internal webhook |

You can mix multiple backends and protocols in a single endpoint (e.g. REST + gRPC + Lambda + queue).

```
┌─────────────┐     HTTP/WS/gRPC      ┌──────────────┐     REST/gRPC/SOAP/Queue     ┌─────────────┐
│   Client    │ ◄──────────────────► │   KrakenD    │ ◄──────────────────────────► │  Backends   │
└─────────────┘                       └──────────────┘                               └─────────────┘
                                            │
                                            │ async_agent (no HTTP client)
                                            ▼
                                      ┌──────────────┐
                                      │ Kafka/Rabbit │
                                      └──────────────┘
```

## Feature index

| # | Feature | Edition | Config namespace | Doc |
|---|---------|---------|------------------|-----|
| 1 | [SOAP Integration](01-soap.md) | CE | `backend/soap` | Legacy XML → modern REST/JSON |
| 2 | [WebSockets](02-websockets.md) | Enterprise | `websocket` (endpoint) | Bidirectional real-time (RFC-6455) |
| 3 | [gRPC Overview & Catalog](03-grpc-overview.md) | Enterprise | `grpc` (service) | Protocol buffers catalog setup |
| 4 | [gRPC Client](04-grpc-client.md) | Enterprise | `backend/grpc` | Consume gRPC upstreams; REST conversion |
| 5 | [gRPC Server](05-grpc-server.md) | Enterprise | `grpc.server` (service) | Expose gRPC on gateway port |
| 6 | [HTTP Streaming & SSE](06-streaming-sse.md) | Enterprise | `output_encoding: no-op` | Long-lived HTTP streams |
| 7 | [RabbitMQ Consumer](07-rabbitmq-consumer.md) | CE + EE | `backend/amqp/consumer` | Pull messages on HTTP request |
| 8 | [RabbitMQ Producer](08-rabbitmq-producer.md) | CE + EE | `backend/amqp/producer` | Push messages on HTTP request |
| 9 | [Pub/Sub (Kafka, NATS, cloud)](09-pubsub.md) | CE + EE | `backend/pubsub/*` | Multi-broker publish/subscribe |
| 10 | [Kafka Advanced PubSub](10-kafka-advanced-pubsub.md) | Enterprise | `backend/pubsub/*/kafka` | Kafka with mTLS/SASL |
| 11 | [GraphQL](11-graphql.md) | CE + EE | `backend/graphql` | REST↔GraphQL adapter or proxy |
| 12 | [AWS Lambda](12-lambda.md) | CE + EE | `backend/lambda` | Invoke serverless functions |
| 13 | [Async Agents](13-async-agents.md) | CE + EE | `async_agent` (root) | Event-driven background workers |
| 14 | [AMQP Async Driver](14-async-amqp-driver.md) | CE + EE | `async/amqp` | RabbitMQ consumer for agents |
| 15 | [Kafka Async Driver](15-async-kafka-driver.md) | Enterprise | `async/kafka` | Kafka consumer for agents |

**CE + EE** = available in Community Edition and Enterprise. **Enterprise** = Enterprise Edition only.

## Edition summary

| Available in CE | Enterprise-only |
|---------------|-----------------|
| RabbitMQ Consumer/Producer | WebSockets |
| Pub/Sub (basic drivers) | gRPC (client + server) |
| GraphQL | HTTP Streaming / SSE |
| Lambda | Kafka Advanced PubSub |
| Async Agents + AMQP driver | Kafka Async Driver |
| SOAP |

## Velonetics bugfix changelog

See [BUGFIX-CHANGELOG.md](BUGFIX-CHANGELOG.md) for module versions and fixes applied during non-REST connectivity hardening (websocket, gRPC, pubsub/kafka, SOAP, AMQP, GraphQL, streaming).

## Common configuration concepts

- **`extra_config`**: Feature-specific settings live under namespaced keys (e.g. `backend/soap`, `websocket`).
- **`disable_host_sanitize: true`**: Required for non-HTTP schemes (`amqp://`, `ws://`, `kafka://`, etc.).
- **Backoff strategies**: `linear`, `linear-jitter`, `exponential`, `exponential-jitter`, `fallback` — used by WebSockets, AMQP, and async agents for reconnection.
- **Placeholder capitalization**: URL params like `{country}` become `.req_params.Country` in templates and `Country` in AMQP/Lambda keys.

## Official links

- [Non-REST Connectivity overview](https://www.krakend.io/docs/enterprise/non-rest-connectivity/)
- [KrakenD schema](https://www.krakend.io/schema/)
