# OpenAPI Generator Service

Internal use only.

## Backend (Go)

```bash
cd api
go run ./cmd/main.go
```

## Observability (Local Dev Stack)

The local OTel stack lives in `api/docker-compose.yaml` and `api/resources/configs/`.
It covers all three signals: traces (Tempo), metrics (Prometheus + Mimir), logs (Loki).

### Default stack

```bash
cd api
docker compose up
```

Services: postgres, redis, nats, otel-collector, tempo, prometheus, loki, mimir.

Signal flow:

- **Traces**: app -> otel-collector (OTLP) -> Tempo (multi-tenant, `X-Scope-OrgID: dev-tenant-1`).
- **Metrics**: app -> otel-collector (OTLP) -> Prometheus (pull, `otel-collector:8889`)
  and Mimir (remote-write push, `http://mimir:8080/api/v1/push`).
- **Logs**: app -> otel-collector (OTLP) -> Loki, with PII redaction (`transform/redact-pii`).

### Local vs. prod differences

- **Traces are 100% sampled locally.** Prod uses tail sampling (~10% retention).
- **PII redaction IS enabled locally** using the same `transform/redact-pii` processor
  as prod (redacting `user.email` / `user.phone`). Prod adds more fields.
- **Tempo/Loki/Mimir run in monolithic mode with filesystem storage.** Prod runs
  distributed with S3/GCS. Trace routing, compaction, and multi-tenant isolation are
  not fully replicated locally.
- **Tempo runs multi-tenant** (`multitenancy_enabled: true`); the collector sets
  `X-Scope-OrgID: dev-tenant-1` on trace and log exports, matching prod header
  propagation.

### Image versions

Pinned to avoid silent breaking changes:

- `otel/opentelemetry-collector-contrib:0.161.0`
- `grafana/tempo:3.0.3`
- `grafana/loki:3.7.8`
- `grafana/mimir:3.2.1`
- `prom/prometheus:v3.14.0`
- `ghcr.io/open-telemetry/opentelemetry-collector-contrib/telemetrygen:0.161.0`

Bump deliberately.

### Running the app against the stack

The app runs on the host (`go run ./cmd/main.go`) and sends OTLP to
`127.0.0.1:4318` (the collector's host-mapped HTTP port). When the app itself is
containerized, set the endpoint to the collector service via environment:

```bash
APP_METRICS_MAIN_ENDPOINT=otel-collector:4318 \
APP_TRACER_MAIN_ENDPOINT=otel-collector:4318 \
  go run ./cmd/main.go
```

### No container healthcheck on otel-collector

The distroless image has no `curl`/`wget`. The `health_check` extension at port
`13133` is the collector's own readiness signal — rely on it (or `docker compose ps`)
instead of a Docker healthcheck.

`telemetrygen` is included under the `load-test` profile — it is a test tool,
not a service, so it is not started by default.

## License

Internal use only.
