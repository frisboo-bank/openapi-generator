# OpenAPI Generator

A tool for generating OpenAPI specifications from frisboo-core-banking APIs.

## Structure

- `api/` - Go backend (OpenAPI generator service)
- `frontend/` - React/TypeScript frontend application (Rsbuild)

## Build

The root `Makefile` is the single build entry point for the whole repo
(Go backend in `api/`, React/TypeScript frontend in `frontend/`).

```bash
make help        # list all available targets
make install     # install all dependencies (frontend + go modules)
make build       # build backend + frontend
make test        # run all tests
make lint        # lint backend + frontend
```

Per-sub-project commands:

```bash
make run            # run the API server (backend)
make dev-frontend   # start the frontend dev server
```

## Getting Started

### Frontend

```bash
cd frontend
pnpm install
pnpm run dev
```

### Backend (Go)

```bash
cd api
go run ./cmd/main.go
```

## Observability (Local Dev Stack)

The local OTel stack lives in `api/docker-compose.yaml` and `api/resources/configs/`.
It covers all three signals: traces (Tempo), metrics (Prometheus), logs (Loki).

### Default stack

```bash
cd api
docker compose up
```

Services: postgres, redis, nats, otel-collector, tempo, prometheus, loki.

### ⚠️ Local vs. prod differences

- **Traces are 100% sampled locally.** Prod uses tail sampling (~10% retention).
  Do not build volume, cost, or cardinality assumptions on local trace data.
- **PII redaction is NOT enabled locally.** Prod enables `transform/redact-pii`
  (see the commented-out stanza in `otel-collector-config.yaml`). Uncomment it
  locally to test redaction behavior.
- **Local Tempo/Loki run in monolithic mode with filesystem storage.** Prod runs
  distributed with S3/GCS. Trace routing, compaction, and multi-tenant isolation
  are not replicated locally.
- **Tempo uses multi-tenant mode** with `X-Scope-OrgID` headers. The collector
  sets `dev-tenant-1` on all trace and log exports. Prod uses the same header
  propagation pattern.

### Optional parity profiles

```bash
# S3-backed Loki testing (bucket naming, auth, path-style access, retries)
docker compose --profile storage-parity up

# Mimir remote-write testing (cardinality limits, remote-write rejection)
docker compose --profile metrics-parity up
```

`telemetrygen` is included under the `load-test` profile — it is a test tool,
not a service, so it is not started by default.

### Verifying the persistent queue

```bash
cd api
./test-persistent-queue.sh 100000
```

Blasts the collector with a bounded burst of logs while restarting Tempo, then
asserts the delivered count matches the sent count (delta = 0). Without this
script the persistent queue is just YAML — this is the only way to prove
backpressure, disk buffering, and flushing actually work.

## License

Internal use only.
