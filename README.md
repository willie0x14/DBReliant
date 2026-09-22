# DBReliant

DBReliant is a Go + PostgreSQL reliability lab for reproducing common database
problems and testing how to diagnose and fix them.

The project focuses on backend engineering, PostgreSQL performance, concurrency,
locking, and operational troubleshooting using a small payment domain.

## Components

- Go HTTP service with a health endpoint and graceful shutdown
- PostgreSQL 18 in Docker Compose
- SQL migrations and deterministic seed data
- Experiment notes with raw `EXPLAIN (ANALYZE, BUFFERS)` output
- `pg_stat_statements` preloaded for future diagnostics

The diagnostic CLI and workload generator are currently placeholders.

## Current Experiment

### [Slow Query / Index Optimization](experiments/slow-query/README.md)

- 1,000,000-row payments dataset
- `Parallel Seq Scan` -> `Index Scan`
- Median latency: `51.845 ms` -> `1.190 ms`
- Shared buffer hits: `9,321` -> `104`

These are local lab results, not production benchmarks.

## Quick Start

```bash
cp .env.example .env
make setup
make migrate
make seed
make build
make run
```

`make seed` creates 100 merchants, 1,000 TWD accounts, and 100,000 deterministic
payments by default. Set a different payment count with:

```bash
make seed PAYMENT_COUNT=250000
```

The API exposes:

```text
GET /health
```

Expected response:

```json
{"status":"ok"}
```

Useful commands:

```bash
make test
make psql
make down
```

## Status

- Implemented: local PostgreSQL setup, schema, deterministic seed data, health
  endpoint, and slow-query experiment
- Placeholder: diagnostic CLI and load generator

## Future Work

- Locking and deadlock experiment
- Connection exhaustion experiment
- Unsafe schema migration experiment

## PostgreSQL Notes

PostgreSQL preloads `pg_stat_statements`. Create the extension when needed:

```sql
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
```
