#!/usr/bin/env bash
# Verify the OTel Collector's persistent queue actually buffers to disk
# and flushes on backend recovery.
#
# Usage: ./test-persistent-queue.sh [num_logs]
#   num_logs defaults to 100000
#
# The test:
#   1. Blasts the collector with a bounded burst of logs.
#   2. Restarts Loki AND Tempo mid-burst.
#   3. Polls both backends' readiness endpoints until recovered.
#   4. Asserts the delivered count matches the sent count (delta = 0).
#
# Without this script the persistent queue is just YAML — this is the only
# way to prove backpressure, disk buffering, and flushing actually work.
set -euo pipefail

NUM_LOGS="${1:-100000}"
LOKI_BASE="http://loki:3100"
LOKI_QUERY="${LOKI_BASE}/loki/api/v1/query"
LOKI_READY="${LOKI_BASE}/ready"
TEMPO_BASE="http://tempo:3200"
TEMPO_READY="${TEMPO_BASE}/ready"
COLLECTOR_HTTP="http://otel-collector:4318"
WAIT_TIMEOUT="${WAIT_TIMEOUT:-120}"

echo "=== Persistent queue test: ${NUM_LOGS} records ==="

# 1. Ensure the stack is up (default profile).
echo "[1/5] Ensuring stack is up..."
docker compose up -d --wait >/dev/null 2>&1 || true
docker compose up -d otel-collector loki tempo >/dev/null 2>&1

# Wait for collector to be ready.
for i in $(seq 1 30); do
  if curl -sf "http://otel-collector:13133/" >/dev/null 2>&1; then
    echo "      Collector ready."
    break
  fi
  sleep 2
done

# 2. Record pre-restart baseline log count.
#    Uses the instant query endpoint (not query_range) — query_range returns
#    values as [[timestamp, "value"], ...] pairs, which broke float() parsing.
#    The instant query returns a single [timestamp, "value"] pair.
echo "[2/5] Recording pre-restart baseline..."
BASELINE=$(curl -s "${LOKI_QUERY}" \
  --data-urlencode 'query=count_over_time({service_name="openapi-generator"}[1m])' \
  2>/dev/null | python3 -c "
import sys, json
d = json.load(sys.stdin)
try:
    print(int(float(d['data']['result'][0]['value'][1])))
except (KeyError, IndexError, ValueError):
    print(0)
" 2>/dev/null || echo "0")
echo "      Baseline log count: ${BASELINE}"

# 3. Blast the collector with a bounded burst while restarting Loki AND Tempo.
echo "[3/5] Blasting ${NUM_LOGS} records while restarting Loki + Tempo..."
docker compose run --rm --profile load-test telemetrygen \
  logs \
    --logs "${NUM_LOGS}" \
    --workers 10 \
    --rate 1000 \
    --otlp-http \
    --otlp-endpoint "${COLLECTOR_HTTP}" \
    --otlp-insecure \
    --service openapi-generator \
    --duration 60s &
TELEMETRYGEN_PID=$!

# Give telemetrygen a moment to start, then restart both backends.
sleep 5
echo "      Restarting Loki..."
docker compose restart loki >/dev/null 2>&1
echo "      Restarting Tempo..."
docker compose restart tempo >/dev/null 2>&1

# Wait for telemetrygen to finish.
wait "${TELEMETRYGEN_PID}" || true
echo "      Blast complete."

# 4. Poll Loki AND Tempo readiness endpoints until recovered.
echo "[4/5] Waiting for backends to recover (timeout ${WAIT_TIMEOUT}s)..."
DEADLINE=$(( $(date +%s) + WAIT_TIMEOUT ))
while true; do
  NOW=$(date +%s)
  if [ "${NOW}" -ge "${DEADLINE}" ]; then
    echo "      ERROR: Backends did not recover within ${WAIT_TIMEOUT}s."
    exit 1
  fi
  LOKI_OK=$(curl -sf "${LOKI_READY}" >/dev/null 2>&1 && echo "1" || echo "0")
  TEMPO_OK=$(curl -sf "${TEMPO_READY}" >/dev/null 2>&1 && echo "1" || echo "0")
  if [ "${LOKI_OK}" = "1" ] && [ "${TEMPO_OK}" = "1" ]; then
    echo "      Loki and Tempo ready."
    break
  fi
  sleep 2
done

# Give the collector a moment to flush the queue.
sleep 10

# 5. Assert delta = 0.
echo "[5/5] Asserting delivery..."
AFTER=$(curl -s "${LOKI_QUERY}" \
  --data-urlencode 'query=count_over_time({service_name="openapi-generator"}[1h])' \
  2>/dev/null | python3 -c "
import sys, json
d = json.load(sys.stdin)
try:
    print(int(float(d['data']['result'][0]['value'][1])))
except (KeyError, IndexError, ValueError):
    print(0)
" 2>/dev/null || echo "0")
echo "      Delivered: ${AFTER} (baseline ${BASELINE}, sent ${NUM_LOGS})"
DELTA=$(python3 -c "print(int(${AFTER}) - int(${BASELINE}) - ${NUM_LOGS})" 2>/dev/null || echo "unknown")
echo "      Delta (delivered - baseline - sent): ${DELTA}"

if [ "${DELTA}" = "0" ]; then
  echo "PASS: all ${NUM_LOGS} records delivered after backend restart."
  exit 0
else
  echo "FAIL: delta ${DELTA} != 0. Queue did not flush correctly."
  exit 1
fi