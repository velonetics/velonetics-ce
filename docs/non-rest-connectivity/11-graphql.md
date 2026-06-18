# GraphQL Integration

**Edition:** Community + Enterprise  
**Namespace:** `backend/graphql`  
**Official docs:** [GraphQL Backend Integration](https://www.krakend.io/docs/enterprise/backends/graphql/)

## What it does

Integrates KrakenD with GraphQL servers in two modes:

1. **REST → GraphQL adapter** — clients call REST; KrakenD builds and sends GraphQL queries/mutations
2. **GraphQL proxy** — KrakenD validates/rate-limits; forwards original GraphQL to backend

## Key capabilities

- Simple GraphQL federation (aggregate multiple subgraphs in config)
- JWT validation before hitting GraphQL servers
- Rate limiting GraphQL usage
- Hide GraphQL complexity behind REST
- Parallel subgraph queries with `group` aggregation

## Mode 1: REST to GraphQL transformation

```
Client POST /review/1500  (REST JSON)
    │
    ▼
KrakenD loads .graphql file, merges variables
    │
    ▼
POST { query, variables, operationName } → GraphQL server
    │
    ▼
GraphQL response → JSON/XML/RSS to client
```

### Configuration

```json
{
  "endpoint": "/marketing/{user_id}",
  "method": "POST",
  "backend": [{
    "timeout": "4100ms",
    "url_pattern": "/graphql?timeout=4s",
    "extra_config": {
      "backend/graphql": {
        "type": "mutation",
        "query_path": "./graphql/mutations/marketing.graphql",
        "variables": {
          "user": "{user_id}",
          "other_static_variables": { "foo": false, "bar": true }
        },
        "operationName": "addMktPreferencesForUser"
      }
    }
  }]
}
```

### `backend/graphql` fields

| Field | Description |
|-------|-------------|
| `type` * | `query` (read) or `mutation` (write) |
| `query_path` | Path to `.graphql` file loaded at startup |
| `query` | Inline GraphQL query (alternative to `query_path`) |
| `variables` | Static + dynamic variables (`{placeholders}` from URL) |
| `operationName` | Operation name for multi-operation documents |

### Method + type matrix

| Endpoint method | GraphQL type | Parameter source | Sent to GraphQL as |
|-----------------|--------------|------------------|-------------------|
| GET | query | URL `{variables}` | Query string |
| GET | mutation | Request body | Query string |
| POST | query | URL `{variables}` | JSON body |
| POST | mutation | Body + config variables | JSON body (user overrides config on collision) |

## Mode 2: GraphQL proxy

Forward client GraphQL queries unchanged; add gateway middleware:

```json
{
  "endpoint": "/graphql",
  "method": "POST",
  "input_query_strings": ["query", "operationName", "variables"],
  "backend": [{
    "timeout": "4100ms",
    "host": ["http://your-graphql.server:4000"],
    "url_pattern": "/graphql?timeout=4s"
  }]
}
```

For browser tools (Apollo Studio), also add:
- `auto_options: true` for OPTIONS preflight
- CORS configuration

## Mode 3: GraphQL federation

Aggregate multiple subgraphs in parallel using `group`:

```json
{
  "endpoint": "/user-data/{id_user}",
  "backend": [
    {
      "group": "user",
      "host": ["http://user-graph:4000"],
      "url_pattern": "/graphql",
      "extra_config": {
        "backend/graphql": {
          "type": "query",
          "query_path": "./graphql/queries/user.graphql",
          "variables": { "user": "{user_id}" },
          "operationName": "getUserData"
        }
      }
    },
    {
      "group": "user_metadata",
      "host": ["http://metadata:4000"],
      "url_pattern": "/graphql",
      "extra_config": {
        "backend/graphql": {
          "type": "query",
          "query_path": "./graphql/queries/user_metadata.graphql",
          "variables": { "user": "{user_id}" },
          "operationName": "getUserMetadata"
        }
      }
    }
  ]
}
```

Each backend can have its own timeout, validation, and rate limiting.

## Notes

- Since KrakenD CE v2.3.3, Content-Type `application/json` is set automatically for GraphQL backends.
- `query_path` files are loaded at startup only — changes require restart.
