\set ON_ERROR_STOP on

BEGIN;

-- Reset the workload tables so the same payment_count produces the same rows.
TRUNCATE TABLE transfers, payments, accounts, merchants RESTART IDENTITY;

INSERT INTO merchants (name, created_at)
SELECT
    format('Merchant %s', lpad(merchant_number::text, 3, '0')),
    TIMESTAMPTZ '2025-01-01 00:00:00+00' + merchant_number * INTERVAL '1 second'
FROM generate_series(1, 100) AS merchant_number;

INSERT INTO accounts (merchant_id, balance, currency, created_at)
SELECT
    ((account_number - 1) % 100) + 1,
    0,
    'TWD',
    TIMESTAMPTZ '2025-01-02 00:00:00+00' + account_number * INTERVAL '1 second'
FROM generate_series(1, 1000) AS account_number;

INSERT INTO payments (account_id, amount, status, created_at)
SELECT
    ((payment_number - 1) % 1000) + 1,
    ((payment_number * 7919) % 100000) + 1,
    CASE
        WHEN ((payment_number - 1) / 1000) % 10 < 7 THEN 'completed'
        WHEN ((payment_number - 1) / 1000) % 10 < 9 THEN 'processing'
        ELSE 'failed'
    END,
    TIMESTAMPTZ '2025-01-03 00:00:00+00' + payment_number * INTERVAL '1 second'
FROM generate_series(1, :'payment_count'::bigint) AS payment_number;

COMMIT;
