-- Run this in Terminal B after Terminal A stops at its prompt.

SET application_name = 'ordered_b';

BEGIN;

-- This may wait for Terminal A, but both sessions use the same lock order.
SELECT id, balance
FROM accounts
WHERE id = 1
FOR UPDATE;

SELECT id, balance
FROM accounts
WHERE id = 2
FOR UPDATE;

COMMIT;
