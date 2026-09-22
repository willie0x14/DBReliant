# Connection Exhaustion and Pool Lab

## Goal

Check how Go's `database/sql` pool behaves under concurrent load and how pool
limits affect PostgreSQL connection usage, queueing, and timeouts.

These are local lab results, not production performance measurements.

## PostgreSQL Limit

The server was configured with:

```ini
max_connections = 100
reserved_connections = 0
superuser_reserved_connections = 3
```

The experiment used the non-superuser role `dbreliant_app`. Under connection
pressure, the observed client connections were:

```yaml
dbreliant     active: 1
dbreliant_app active: 96
total client connections: 97
```

Additional application connections failed with:

```text
FATAL: remaining connection slots are reserved for roles with the SUPERUSER attribute
SQLSTATE 53300
```

The PostgreSQL connection limit is a shared system budget. Every application
instance and other client consumes part of the same limit.

## Unbounded Pool

The first run started 100 goroutines without calling `SetMaxOpenConns`. Each ran:

```sql
SELECT pg_sleep(2);
```

During load:

```yaml
OpenConnections: 100
InUse: 100
Idle: 0
WaitCount: 0
```

Several workers failed with PostgreSQL `53300`.

```text
application concurrency
-> connection attempts
-> pressure reaches PostgreSQL directly
```

## Bounded Pool

The same 100-worker workload was repeated with:

```go
db.SetMaxOpenConns(10)
```

During load:

```yaml
OpenConnections: 10
InUse: 10
Idle: 0
WaitCount: 90
```

After load:

```yaml
OpenConnections: 2
InUse: 0
Idle: 2
WaitCount: 90
WaitDuration: approximately 15 minutes cumulative
```

The workload took approximately **20.1 seconds**. `WaitDuration` is summed
across waiting callers, so it can be much larger than wall-clock runtime.

The pool protected PostgreSQL by moving excess concurrency into an
application-side wait queue.

## Idle Connection Reuse

The next run added:

```go
db.SetMaxIdleConns(10)
```

After load:

```yaml
OpenConnections: 10
InUse: 0
Idle: 10
MaxIdleClosed: 0
```

Idle connections can be reused instead of reopened, but they still consume
PostgreSQL connection slots. `MaxIdleConns` must fit within the same global
database connection budget.

## Context Timeout

Each worker used a 3-second context timeout with `ExecContext`:

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

_, err := db.ExecContext(ctx, "SELECT pg_sleep(2)")
```

Configuration:

```text
workers: 100
MaxOpenConns: 10
query duration: 2 seconds
```

During load, all 10 connections were in use, `WaitCount` initially reached 90,
and many workers returned `context deadline exceeded`.

After a short cleanup wait:

```yaml
OpenConnections: 10
InUse: 0
Idle: 10
WaitDuration: approximately 4m20s cumulative
```

The workload completed around the 3-second deadline. The program prints that
workload duration separately from the later cleanup observation and total
duration.

A context timeout covers both waiting for a pool connection and executing SQL.
Workers may time out at either stage; this run did not establish that every
timeout happened at the same stage.

## Comparison

| Setup | Database effect | Application effect |
| --- | --- | --- |
| Unbounded pool | Connection exhaustion risk | Little pool waiting until PostgreSQL rejects connections |
| Bounded pool | Connection count capped | Requests queue and latency increases |
| Bounded pool + timeout | Connection count capped | Requests fail instead of waiting indefinitely |

## What I Learned

- Database connection limits are shared across application instances.
- `MaxOpenConns` provides application-side backpressure.
- A pool that is too large can exhaust PostgreSQL.
- A pool that is too small increases queueing latency.
- Idle connections are reusable but still consume database capacity.
- Request and database timeouts prevent unbounded waiting.
- Pool sizing must consider database capacity, application instance count, and
  latency requirements.

## Reproduction

The current experiment program contains the bounded pool, idle reuse, and
context timeout configuration:

```bash
go run ./experiments/connections/pool_lab.go
```

It connects as `dbreliant_app`, starts 100 workers, prints a pool snapshot during
load, waits for workers to finish, then prints the cleanup snapshot. The role and
database must already exist before running it.
