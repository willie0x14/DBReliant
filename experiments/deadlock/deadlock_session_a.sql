-- Run this in Terminal A with psql.
-- Keep ON_ERROR_STOP disabled so ROLLBACK runs if PostgreSQL aborts this transaction.
\set ON_ERROR_STOP off

SET application_name = 'deadlock_a';

BEGIN;

-- Lock account 1 first.
SELECT id, balance
FROM accounts
WHERE id = 1
FOR UPDATE;

-- Stop here and switch to Terminal B.
-- Start deadlock_session_b.sql and wait at its prompt before continuing.
\prompt 'Press Enter after Terminal B has locked account 2: ' continue

-- This waits because Terminal B holds account 2.
SELECT id, balance
FROM accounts
WHERE id = 2
FOR UPDATE;

-- PostgreSQL aborts one transaction after Terminal B requests account 1.
-- ROLLBACK also clears the failed transaction state if this session is aborted.
ROLLBACK;
