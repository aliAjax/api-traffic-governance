#!/bin/sh
set -eu
base=${TRAFFIC_BASE_URL:-http://127.0.0.1:8085}
curl -fsS "$base/healthz" >/dev/null
curl -fsS -X POST "$base/api/v1/services" -H 'content-type: application/json' -d '{"tenant_id":"demo","name":"catalog","region":"local"}' | grep -q 'catalog'
echo "traffic governance smoke ok"
