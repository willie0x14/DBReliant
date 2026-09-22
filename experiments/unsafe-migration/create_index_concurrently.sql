-- Run this in the migration_ddl terminal. Do not wrap this script in BEGIN/COMMIT.
-- Keep error-stop disabled so the final cleanup can remove an INVALID index.
\set ON_ERROR_STOP off

SET application_name = 'migration_ddl';

-- Clean up before starting the writer transaction in another terminal.
DROP INDEX IF EXISTS idx_payments_amount_lab;

-- In the writer terminal, BEGIN and update payment 1 without committing.
\prompt 'Press Enter after the writer transaction is open: ' continue

CREATE INDEX CONCURRENTLY idx_payments_amount_lab
ON payments (amount);

DROP INDEX IF EXISTS idx_payments_amount_lab;
