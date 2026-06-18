# Pub/Sub examples

Docker Compose stacks demonstrating KrakenD-compatible `backend/pubsub/publisher` and `backend/pubsub/subscriber` backends.

| Driver | Directory | Broker | Env variable |
|--------|-----------|--------|--------------|
| NATS | [nats/](nats/) | `nats:2-alpine` | `NATS_SERVER_URL` |
| Kafka | [kafka/](kafka/) | Redpanda | `KAFKA_BROKERS` |
| RabbitMQ | [rabbit/](rabbit/) | `rabbitmq:3-management` | `RABBIT_SERVER_URL` |
| GCP Pub/Sub | [gcp/](gcp/) | Pub/Sub emulator | `PUBSUB_EMULATOR_HOST` |
| AWS SNS/SQS | [aws/](aws/) | LocalStack | AWS credentials |
| Azure Service Bus | [azure/](azure/) | Service Bus emulator | `SERVICEBUS_CONNECTION_STRING` |
| Kafka Advanced | [kafka-advanced/](kafka-advanced/) | Redpanda + SASL | explicit `cluster` config |
| Async Kafka | [async-kafka/](async-kafka/) | Redpanda | `async/kafka` agent |

## Quick start (NATS)

```bash
cd examples/pubsub/nats
docker compose up --build -d
./scripts/smoke.sh
docker compose down -v
```

Or from the CE root:

```bash
make pubsub-nats-compose-test
```

Each subfolder has its own `docker-compose.yml`, `velonetics.json`, and `scripts/smoke.sh`.
