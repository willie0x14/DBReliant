# Unsafe Migration / Online Indexing Lab

## Goal

Compare write availability during `CREATE INDEX` and `CREATE INDEX CONCURRENTLY`.
These are local lab observations on the payments table.

## Normal CREATE INDEX

Session `migration_writer` opened a transaction and left it active:

```sql
BEGIN;
UPDATE payments SET amount = amount + 1 WHERE id = 1;
```

Session `migration_ddl` then ran:

```sql
CREATE INDEX idx_payments_amount_lab ON payments (amount);
```

The index build waited for the writer. `pg_stat_activity` showed it as active,
with `wait_event_type: Lock`, `wait_event: relation`, and
`blocked_by: migration_writer`. After the writer ended, the build completed.
Normal `CREATE INDEX` needs a lock that conflicts with writes, so an existing
writer can make the DDL wait.

## Lock Queue

The test was repeated with `migration_writer_1` holding the first write open.
While `migration_ddl` waited, `migration_writer_2` attempted:

```sql
BEGIN;
UPDATE payments SET amount = amount + 1 WHERE id = 2;
```

The new writer also waited:

```text
migration_ddl      blocked_by migration_writer_1
migration_writer_2 blocked_by migration_ddl
```

```text
long transaction -> DDL waits -> new writes queue behind DDL -> latency grows
```

Both writer sessions used `ROLLBACK`, so payment amounts were not changed.

## CREATE INDEX CONCURRENTLY

The concurrent version was tested while a writer remained active:

```sql
CREATE INDEX CONCURRENTLY idx_payments_amount_lab ON payments (amount);
```

The DDL session waited on `virtualxid`, blocked by `migration_writer`.
`pg_stat_progress_create_index` reported:

```text
command: CREATE INDEX CONCURRENTLY
phase: waiting for writers before build
```

The writer could still run another update on payment 2. `CONCURRENTLY` is not
lock-free and can wait for transactions during specific phases. Its benefit is
that normal inserts, updates, and deletes can continue during most of the build.

## Transaction Block Limitation

This failed as expected:

```sql
BEGIN;
CREATE INDEX CONCURRENTLY idx_payments_amount_lab ON payments (amount);
-- ERROR: CREATE INDEX CONCURRENTLY cannot run inside a transaction block
ROLLBACK;
```

Concurrent index creation has multiple internal phases. Migration tooling must
support running this operation outside an explicit transaction.

## Comparison

| `CREATE INDEX` | `CREATE INDEX CONCURRENTLY` |
| --- | --- |
| Simpler and usually faster | Keeps normal DML available during most of the build |
| Can block normal writes | Can still wait for existing transactions |
| Runs in a transaction | Cannot run inside an explicit transaction block |
| | Usually takes longer and does more work |
| | A failed build can leave an `INVALID` index to inspect and remove |

## What I Learned

- DDL can become a source of production blocking.
- Waiting DDL can cause later application writes to queue.
- `CREATE INDEX CONCURRENTLY` improves write availability but is not lock-free.
- Long-running transactions matter during online schema changes.
- Migration tooling must support operations that cannot run in a transaction.

## Reproduction

Open separate `psql` terminals and set the session names shown above. Use
[`create_index.sql`](create_index.sql) for the normal build or
[`create_index_concurrently.sql`](create_index_concurrently.sql) for the online
build. Both scripts include cleanup for `idx_payments_amount_lab`.

Set the observer's `application_name` to `migration_observer`, then run
[`../../queries/diagnostics/blocking.sql`](../../queries/diagnostics/blocking.sql).
For a concurrent build, also run:

```sql
SELECT pid, command, phase
FROM pg_stat_progress_create_index;
```

Finish open writer transactions with `ROLLBACK`.
