# Changelog

## v0.1.0 — 2026-03-23

Initial release. Built from scratch as a clean-room implementation of an S3-backed log aggregation system.

### Core

- **Log ingestion** via `POST /logs` — accepts JSON batches with timestamp, service, level, message, and arbitrary labels
- **In-memory buffer** with configurable size threshold (default 5MB) and time window (default 30s) flushing
- **Gzip and Snappy compression** — configurable per deployment, auto-detected on read
- **S3 storage** with hierarchical key layout: `<service>/YYYY/MM/DD/<unix>-<uuid>.<ext>`
- **Write-Ahead Log (WAL)** — optional on-disk durability; entries are fsync'd before buffering, truncated after confirmed S3 upload, replayed on crash recovery
- **Retry with exponential backoff** — all S3 and OpenSearch operations retry up to 3 times with jitter
- **Flush failure retention** — buffer is NOT cleared on S3 write failure; entries stay in memory + WAL for retry

### Label Index

- **Minimal label index** written to S3 alongside each chunk (`index/YYYY/MM/DD/<service>-<id>.json`)
- Records unique label sets, time bounds, service, entry count, and size per chunk
- Query engine reads index files first, then fetches only matching chunks — avoids full S3 scan

### Query

- **Query API** (`GET /query`) — search logs by time range, Loki-style label selectors (`{app="nginx", env="prod"}`), level, with pagination
- Returns matching entries plus performance stats (index files scanned, chunks matched/fetched, entries scanned/matched, duration)
- **Label selector parser** — supports `{key="value", key2="value2"}` syntax

### Web UI

- **Embedded SPA** served at `/ui` on the exporter (no separate deployment)
- Dark theme, monospace font, GitHub-dark color palette
- Time range pickers, label query input, level filter, result limit
- Color-coded log levels, inline label tags, query stats bar
- **Live tail** toggle — streams entries via WebSocket with auto-scroll

### Live Tail

- **WebSocket endpoint** on the ingester (`WS /tail?query={...}&level=...`)
- Pub/sub broker fans out ingested entries to subscribers with per-subscriber label + level filtering
- 256-entry buffered channel per subscriber; non-blocking drops for slow consumers

### OpenSearch Export

- **Export jobs** via `POST /export` — async, bounded worker pool (configurable concurrency, default 4)
- **Job management** — `GET /export/{id}` (poll status), `GET /export` (list all), `DELETE /export/{id}` (cancel)
- Progress tracking: chunks total/processed, logs exported, error count, lifecycle status

### Security

- **API key authentication** — `X-API-Key` header, configurable via `auth.api_keys`
- **Rate limiting** — token-bucket rate limiter on the ingester
- **Request body limits** — `MaxBytesReader` (default 10MB) + max entries per request (default 10,000)
- **TLS support** — optional `ListenAndServeTLS` on the ingester

### Observability

- **Structured logging** — Go `log/slog` with JSON output to stderr
- **Prometheus metrics** — 17 metrics covering ingestion, buffer, WAL, S3, exports, and HTTP
- **Access logging** — every HTTP request logged with method, path, status, duration, bytes

### Operations

- **Config validation** — required fields, sane defaults, TLS/auth consistency checks
- **Graceful shutdown** — SIGINT/SIGTERM flushes buffer, drains export jobs
- **Multitenancy** — `tenant_id` scopes all S3 paths (data + index)
- **Log level sampling** — `min_level` drops entries below a threshold on ingest

### Local Development

- **Docker Compose** — MinIO + OpenSearch for local dev
- **`quickstart.sh`** — one-command local setup with sample data
- **`config.local.json`** — pre-configured for MinIO + local OpenSearch
- **Makefile** — build, test, test-integration, dev-up, dev-down, run-ingester, run-exporter

### Testing

- **78 unit tests** across 7 packages — config validation, buffer, handler, WAL, broker, exporter, middleware (auth, rate limit, body limit), index (label matching, time overlap), query (label parser), retry
- **Integration tests** — S3 round-trip against MinIO (build-tagged `integration`)
- **`test.sh`** — one-command test runner with build, vet, unit tests, integration tests, and per-package coverage summary
