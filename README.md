# DBReliant

DBReliant is a production-style Go and PostgreSQL reliability lab for reproducing, diagnosing, and mitigating common database failure modes under a simplified payment workload.

## Quick Start

```bash
cp .env.example .env
make setup
make migrate
make seed
make build
make run
```

`make seed` creates 100 merchants, 1,000 TWD accounts, and 100,000
deterministic payments by default. Override the payment count when needed:

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

## Planned Experiments

TODO - experiments not implemented yet.

- slow query / index optimization
- locking / deadlock
- connection exhaustion
- unsafe schema migration

## PostgreSQL Notes

PostgreSQL is configured with `shared_preload_libraries=pg_stat_statements` so the extension can be created manually later:

```sql
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
```

No performance indexes, diagnostic SQL, or reliability experiments have been added yet.
