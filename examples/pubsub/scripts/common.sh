#!/usr/bin/env bash
set -euo pipefail

GATEWAY="${GATEWAY:-http://127.0.0.1:8080}"

wait_gateway() {
  local url="${1:-${GATEWAY}/__health}"
  echo "==> Waiting for Velonetics gateway at ${url}"
  for i in $(seq 1 90); do
    if curl -fsS "$url" >/dev/null 2>&1; then
      echo "OK: gateway ready"
      return 0
    fi
    sleep 2
    if [[ "$i" -eq 90 ]]; then
      echo "gateway not ready" >&2
      exit 1
    fi
  done
}

assert_json_field() {
  local resp="$1"
  local field="$2"
  local expected="$3"
  printf '%s' "$resp" | python3 -c "import json,sys; d=json.load(sys.stdin); assert d.get('${field}') == '${expected}', d"
}

pubsub_roundtrip_smoke() {
  local publish_path="${1:-/publish}"
  local subscribe_path="${2:-/subscribe}"
  local payload="${3:-{\"event\":\"smoke-test\"}}"
  local field="${4:-event}"
  local expected="${5:-smoke-test}"

  echo "==> Publishing message"
  curl -fsS -X POST "${GATEWAY}${publish_path}" \
    -H 'Content-Type: application/json' \
    -d "$payload" >/dev/null

  echo "==> Subscribing to message"
  local resp
  for i in $(seq 1 30); do
    if resp="$(curl -fsS "${GATEWAY}${subscribe_path}" 2>/dev/null)"; then
      echo "$resp"
      if assert_json_field "$resp" "$field" "$expected"; then
        echo "All pub/sub smoke checks passed."
        return 0
      fi
    fi
    sleep 1
  done
  echo "unexpected or missing response: ${resp:-<empty>}" >&2
  exit 1
}
