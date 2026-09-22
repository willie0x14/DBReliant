-- Run this in Terminal B with psql.
-- Start after Terminal A has locked account 1 and stopped at its prompt.
\set ON_ERROR_STOP off

SET application_name = 'deadlock_b';

BEGIN;

-- Lock account 2 first, creating the opposite lock order from Terminal A.
SELECT id, balance
FROM accounts
WHERE id = 2
FOR UPDATE;

-- Stop here and return to Terminal A.
-- Continue A so it waits for account 2, then return here.
\prompt 'Press Enter after Terminal A is waiting for account 2: ' continue

-- This completes the wait cycle. PostgreSQL aborts one transaction with 40P01.
SELECT id, balance
FROM accounts
WHERE id = 1
FOR UPDATE;

-- Clean up whether this transaction completed its lock request or was aborted.
ROLLBACK;
