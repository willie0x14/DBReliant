-- Run this in Terminal A with psql.

SET application_name = 'ordered_a';

BEGIN;

-- Both sessions lock accounts in ascending ID order.
SELECT id, balance
FROM accounts
WHERE id = 1
FOR UPDATE;

-- Stop here and start ordered_session_b.sql in Terminal B.
-- Terminal B will wait for this transaction to finish.
\prompt 'Press Enter after Terminal B is waiting for account 1: ' continue

SELECT id, balance
FROM accounts
WHERE id = 2
FOR UPDATE;

COMMIT;
