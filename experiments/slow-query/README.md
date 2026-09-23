# Slow Query / Index Optimization Experiment

## Goal

Test a query that returns the latest 100 completed payments for one account,
then compare the plan before and after adding a composite index.

Dataset:

- 100 merchants
- 1,000 accounts
- 1,000,000 payments
- About 70% completed, 20% processing, and 10% failed per account
- Deterministic seed data

`ANALYZE payments` was run before collecting plans.

## Query

```sql
SELECT
    id,
    account_id,
    amount,
    status,
    created_at
FROM payments
WHERE account_id = 123
  AND status = 'completed'
ORDER BY created_at DESC
LIMIT 100;
```

## Before

The only index was `payments_pkey` on `id`. Nothing supported the filter and
sort used by this query.

```diff
Parallel Seq Scan
-> Sort
-> Gather Merge
-> Limit
```

After a warm-up run, five measurements were recorded:

```text
72.701 ms
51.845 ms
31.057 ms
93.013 ms
33.121 ms
```

- Median: **51.845 ms**
- Shared buffer hits: **9,321**
- Rows filtered: approximately **999,300**

The main issue was the plan shape, not one slow timing sample. PostgreSQL scanned
almost the entire table to find about 700 matching rows, sorted them, and
returned 100.

## Index

```sql
CREATE INDEX idx_payments_account_status_created_at
ON payments (account_id, status, created_at DESC);
```

- `account_id` is the main lookup key.
- `status` narrows the matching range.
- `created_at DESC` matches the requested order.
- Matching the sort order removes the explicit sort and lets `LIMIT 100` stop
  earlier.

The order follows this query's access pattern. It is not based on a rule that
the most selective column must always come first.

## After

```text
Index Scan
-> Limit
```

```text
Index Cond:
account_id = 123
AND status = 'completed'
```

The sequential scan, sort, gather merge, and large filter discard were removed.
After a warm-up run, five measurements were recorded:

```text
0.269 ms
3.033 ms
1.190 ms
2.377 ms
0.632 ms
```

- Median: **1.190 ms**
- Shared buffer hits: **104**

## Result

| Metric | Before | After | Change |
| --- | ---: | ---: | ---: |
| Median execution time | 51.845 ms | 1.190 ms | About 43.6x faster |
| Median latency | 51.845 ms | 1.190 ms | About 97.7% lower |
| Shared buffer hits | 9,321 | 104 | About 98.9% lower |

These are local lab results, not production performance claims.

## What I Learned

- Query plan shape matters more than one timing sample.
- Latency varies between runs, so buffer usage is useful additional evidence.
- Composite index order should match real filters and ordering requirements.
- The read improvement has a cost: the index uses storage and adds maintenance
  work to inserts, updates, and deletes.

## Files / Reproduction

- [`../../scripts/seed.sql`](../../scripts/seed.sql) - deterministic seed data
- [`create_index.sql`](create_index.sql) - index used in the experiment
- [`before.txt`](before.txt) - baseline plan
- [`after.txt`](after.txt) - optimized plan

Seed the 1,000,000-row dataset:

```bash
make seed PAYMENT_COUNT=1000000
```

Normal application setup installs the optimized index through migrations. For
the baseline plan, remove it in the local lab database first:

```sql
DROP INDEX IF EXISTS idx_payments_account_status_created_at;
```

Refresh statistics, then run the query above with plan and buffer details:

```sql
ANALYZE payments;
```

Prefix the query with `EXPLAIN (ANALYZE, BUFFERS)`. Record the baseline after a
warm-up run, apply `create_index.sql`, then repeat the same measurement process.
Reapplying the index also restores the normal application setup.
