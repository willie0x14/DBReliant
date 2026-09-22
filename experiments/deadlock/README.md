# Deadlock and Row Locking Lab

## Goal

Reproduce PostgreSQL row-lock blocking and a real deadlock, inspect the blocked
sessions, and verify that consistent lock ordering prevents the same cycle.

The lab also checks how `SELECT ... FOR UPDATE` protects a concurrent balance
validation.

## Deadlock Reproduction

Two transactions locked the same accounts in opposite order:

```text
Transaction A: account 1 -> account 2
Transaction B: account 2 -> account 1
```

After A locked account 1 and B locked account 2, A requested account 2. A was
blocked, but this was not a deadlock yet because B could still finish.

When B then requested account 1, the wait became a cycle:

```text
A waits for B
B waits for A
```

PostgreSQL detected the deadlock and aborted one transaction with SQLSTATE
`40P01`. The observed error included:

```text
ERROR: deadlock detected

Process B waits for a lock held by process A.
Process A waits for a lock held by process B.
```

Process and transaction IDs change between runs. The aborted transaction stays
in a failed state until it is cleaned up with `ROLLBACK`.

## Diagnosis

[`../../queries/diagnostics/blocking.sql`](../../queries/diagnostics/blocking.sql)
uses `pg_stat_activity` and `pg_blocking_pids()` to show both sessions and their
blockers.

While A was waiting for account 2, the observed state for `deadlock_a` was:

```text
state: active
wait_event_type: Lock
wait_event: transactionid
blocked_by: deadlock_b
```

## Fix: Consistent Lock Ordering

Both transactions were repeated with the same order:

```text
account 1 -> account 2
```

The second transaction could still block while the first held account 1, but
there was no cycle and no deadlock.

```text
Blocking: B waits for A; A can finish.
Deadlock: A waits for B; B waits for A.
```

The practical mitigation is to lock shared resources in a deterministic order.
For a future transfer API, lock accounts by ascending ID regardless of transfer
direction:

```text
Transfer 20 -> 10: lock 10 -> 20
Transfer 10 -> 20: lock 10 -> 20
```

## FOR UPDATE and Concurrent Balance Checks

A normal `SELECT` could still read an account while another transaction held a
`FOR UPDATE` lock. PostgreSQL MVCC allows that read to use a visible row version.
A conflicting `UPDATE` or row-lock request waited until the first transaction
committed or rolled back.

The lock also fixed a concurrent read/modify/write race:

```text
Initial balance: 1000
Transaction A reads: 1000
Transaction B reads: 1000
Both approve an 800-unit debit
Both write: 200
Final balance: 200
```

Both requests made their decision from the same stale balance. With
`SELECT ... FOR UPDATE`, transaction B waited. After A committed the change from
1000 to 200, B continued, saw 200, and could reject another 800-unit debit.
Balance validation must happen while holding the row lock, which remains held
until `COMMIT` or `ROLLBACK`.

## What I Learned

- Blocking and deadlock are different.
- `SELECT ... FOR UPDATE` holds a row lock until the transaction ends.
- Inconsistent lock ordering can create a deadlock cycle.
- PostgreSQL reports deadlocks as SQLSTATE `40P01`.
- Consistent lock ordering avoids this deadlock pattern.
- Balance validation must happen while holding the appropriate lock.

## Reproduction

Open two `psql` terminals. Run the paired scripts with `\i` and follow their
prompts to switch terminals:

```psql
\i experiments/deadlock/deadlock_session_a.sql
\i experiments/deadlock/deadlock_session_b.sql
```

While A is waiting for account 2, inspect the blocking relationship from a
third `psql` terminal:

```psql
\i queries/diagnostics/blocking.sql
```

Use these scripts for the consistent-order version:

```psql
\i experiments/deadlock/ordered_session_a.sql
\i experiments/deadlock/ordered_session_b.sql
```

The scripts only acquire row locks and end with `ROLLBACK` or `COMMIT`; they do
not change account balances.
