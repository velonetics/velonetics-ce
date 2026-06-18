# GraphQL Backend Integration

Velonetics integrates with GraphQL servers in three modes, matching [KrakenD GraphQL](https://www.krakend.io/docs/backends/graphql/) behavior.

## Namespace

```json
"extra_config": {
  "backend/graphql": { ... }
}
```

## Mode 1: REST to GraphQL adapter

Clients call REST endpoints; Velonetics builds and sends GraphQL queries or mutations to the upstream server.

```json
{
  "endpoint": "/review/{id_show}",
  "method": "POST",
  "backend": [{
    "host": ["http://graphql:4000"],
    "url_pattern": "/graphql",
    "extra_config": {
      "backend/graphql": {
        "type": "mutation",
        "query_path": "./graphql/mutations/review.graphql",
        "variables": {
          "ep": "JEDI",
          "review": { "stars": 3, "commentary": "meh" }
        },
        "operationName": "CreateReviewForEpisode"
      }
    }
  }]
}
```

### Fields

| Field | Required | Description |
|-------|----------|-------------|
| `type` | yes | `query` (read) or `mutation` (write) |
| `query` | one of | Inline GraphQL document |
| `query_path` | one of | Path to `.graphql` file (loaded at startup) |
| `variables` | no | Static variables; use `{param}` for URL placeholders |
| `operationName` | no | Operation name for multi-operation documents |

### Method + type matrix

| Endpoint method | GraphQL type | Parameter source | Sent to GraphQL as |
|-----------------|--------------|------------------|-------------------|
| GET | query | URL `{variables}` | Query string |
| GET | mutation | Request body | Query string |
| POST | query | URL `{variables}` | JSON body (client body ignored) |
| POST | mutation | Body + config variables | JSON body (user overrides config on collision) |

HTTP method for the upstream call follows: `extra_config.method` → `backend.method` → endpoint `method` → `POST`.

`Content-Type: application/json` is set automatically on GraphQL backend requests.

## Mode 2: GraphQL proxy

Forward client GraphQL unchanged; apply gateway middleware (JWT, rate limiting, CORS).

```json
{
  "endpoint": "/graphql",
  "method": "POST",
  "input_query_strings": ["query", "operationName", "variables"],
  "backend": [{
    "host": ["http://graphql:4000"],
    "url_pattern": "/graphql"
  }]
}
```

For browser tools (Apollo Studio), add CORS and enable OPTIONS on the router.

## Mode 3: Simple federation

Aggregate multiple subgraphs in parallel using backend `group`:

```json
{
  "endpoint": "/user-data/{id_user}",
  "backend": [
    {
      "group": "user",
      "method": "GET",
      "host": ["http://user-graph:4000"],
      "url_pattern": "/graphql",
      "extra_config": {
        "backend/graphql": {
          "type": "query",
          "query_path": "./graphql/queries/user.graphql",
          "variables": { "user": "{id_user}" },
          "operationName": "getUserData"
        }
      }
    },
    {
      "group": "user_metadata",
      "method": "GET",
      "host": ["http://metadata:4000"],
      "url_pattern": "/graphql",
      "extra_config": {
        "backend/graphql": {
          "type": "query",
          "query_path": "./graphql/queries/user_metadata.graphql",
          "variables": { "user": "{id_user}" },
          "operationName": "getUserMetadata"
        }
      }
    }
  ]
}
```

Each backend can have its own timeout, validation, and rate limiting.

## Notes

- `query_path` files are loaded at startup only; changes require a gateway restart.
- URL path parameters are always sent as strings in GraphQL variables (same as KrakenD).
- GraphQL subscriptions are not supported in v1.

## Local example

```bash
make graphql-compose-test
```

See [examples/graphql/README.md](../examples/graphql/README.md).

## Implementation

Runtime lives in `github.com/velonetics/lura/v2/transport/http/client/graphql` and `proxy.NewGraphQLMiddleware`. No separate Velonetics module is required.
