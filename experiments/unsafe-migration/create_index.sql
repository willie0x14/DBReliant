-- Run this in the migration_ddl terminal.
-- Keep error-stop disabled so the final cleanup can run after a DDL error.
\set ON_ERROR_STOP off

SET application_name = 'migration_ddl';

-- Clean up before starting the writer transaction in another terminal.
DROP INDEX IF EXISTS idx_payments_amount_lab;

-- In the writer terminal, BEGIN and update payment 1 without committing.
\prompt 'Press Enter after the writer transaction is open: ' continue

CREATE INDEX idx_payments_amount_lab
ON payments (amount);

DROP INDEX IF EXISTS idx_payments_amount_lab;
