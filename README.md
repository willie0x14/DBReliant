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

## Current Experiments

### [Slow Query / Index Optimization](experiments/slow-query/README.md)

- 1,000,000-row payments dataset
- `Parallel Seq Scan` -> `Index Scan`
- Median latency: `51.845 ms` -> `1.190 ms`
- Shared buffer hits: `9,321` -> `104`

These are local lab results, not production benchmarks.

### [Deadlock / Row Locking](experiments/deadlock/README.md)

- Reproduced PostgreSQL `40P01` with conflicting row-lock order
- Diagnosed blocking with `pg_stat_activity` and `pg_blocking_pids()`
- Prevented the reproduced pattern with consistent lock ordering
- Verified `SELECT ... FOR UPDATE` for concurrent balance checks

### [Connection Exhaustion / Connection Pool](experiments/connections/README.md)

- Reproduced PostgreSQL connection exhaustion with an unbounded Go pool
- Added application-side backpressure with `SetMaxOpenConns`
- Compared queueing, idle connection reuse, and context timeouts

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
  endpoint, slow-query experiment, deadlock / row-locking experiment, and
  connection pool experiment
- Placeholder: diagnostic CLI and load generator

## Future Work

- Unsafe schema migration experiment

## PostgreSQL Notes

PostgreSQL preloads `pg_stat_statements`. Create the extension when needed:

```sql
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
```
