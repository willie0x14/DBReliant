# DBReliant

DBReliant is a production-style Go and PostgreSQL reliability lab for reproducing, diagnosing, and mitigating common database failure modes under a simplified payment workload.

## Quick Start

```bash
cp .env.example .env
make setup
make build
make run
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

No application schema, tables, constraints, indexes, diagnostic SQL, or workload implementation has been added yet.
