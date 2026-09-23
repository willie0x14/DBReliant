# Transfer API / Load Experiment

## Goal

Verify the transaction-safe transfer flow, then observe HTTP latency and
database pool contention under a controlled local workload.

## Transfer Flow

`POST /transfers` runs the following work in one database transaction:

1. Lock both account rows with `SELECT ... FOR UPDATE`.
2. Always lock the lower account ID first.
3. Validate the source balance while holding the locks.
4. Debit the source, credit the destination, and insert a completed transfer.
5. Commit all changes together.

Deterministic lock ordering mitigates the two-account deadlock pattern reproduced
in the row-locking lab.

## Functional Checks

- A successful transfer updated both balances and inserted a completed row.
- Insufficient funds returned HTTP `409` and left balances unchanged.
- Simultaneous 1→2 and 2→1 transfers both completed without reproducing a
  deadlock.

## Workload

```yaml
Endpoint: POST /transfers
Requests: 5,000
Concurrent workers: 50
Transfer amount: 1
DB_MAX_OPEN_CONNS: 10
Accounts: 1 and 2, alternating direction
Result: all 5,000 requests returned HTTP 201
```

This was a local controlled experiment, not a production benchmark.

## Results

| Metric | Observed value |
| --- | ---: |
| `http_requests_total{method="POST",status="201"}` | 5,000 |
| `db_wait_count_total` | 4,990 |
| `db_wait_duration_seconds_total` | ~243.25 s |
| Estimated p50 handler latency | ~46.44 ms |
| Estimated p95 handler latency | ~210.51 ms |
| Estimated p99 handler latency | ~249.42 ms |

Average connection-pool wait per wait event was approximately:

```text
243.25 s / 4,990 = 48.7 ms
```

## Interpretation

- The bounded pool provides backpressure: 50 workers compete for at most 10
  database connections.
- `db_wait_count_total` counts `database/sql` pool wait events. It does not count
  PostgreSQL row-lock waits and does not necessarily represent unique requests.
- `db_wait_duration_seconds_total` is cumulative across waiters, not workload
  wall-clock duration.
- Every transaction locks accounts 1 and 2, intentionally creating hot-row
  contention.
- Pool waiting, transaction work, and row-lock contention can all contribute to
  HTTP latency. The metrics do not attribute all latency to one cause.
- `histogram_quantile` values are estimates based on histogram buckets.

## Reproduction

Run the API with a 10-connection pool, then start the configurable load generator:

The deterministic seed starts account balances at zero. Fund accounts 1 and 2
in the local lab database before running successful transfer workloads:

```sql
UPDATE accounts SET balance = 100000 WHERE id IN (1, 2);
```

```bash
DB_MAX_OPEN_CONNS=10 make run
go run ./cmd/loadgen -requests 5000 -workers 50 -amount 1
```

The generator alternates transfers between accounts 1 and 2 and treats only
HTTP `201` as success.
