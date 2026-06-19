# GraphQL example

Local Docker Compose stack demonstrating all three KrakenD-equivalent GraphQL modes:

1. **REST to GraphQL adapter** — `POST /review/{id}` and `GET /hero/{episode}`
2. **GraphQL proxy** — `POST /graphql` forwards client queries unchanged
3. **Simple federation** — `GET /user-data/{id}` aggregates two subgraphs in parallel

## Run

```bash
cd examples/graphql
docker compose up --build -d
chmod +x scripts/smoke.sh
./scripts/smoke.sh
docker compose down -v
```

Or from the CE repo root:

```bash
make graphql-compose-test
```

## Layout

- `mock-backend/` — minimal GraphQL server on port 4000
- `graphql/` — `.graphql` query files loaded via `query_path`
- `pucora.json` — gateway configuration
- `scripts/smoke.sh` — curl-based smoke tests
