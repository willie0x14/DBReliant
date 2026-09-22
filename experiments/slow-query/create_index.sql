CREATE INDEX idx_payments_account_status_created_at
ON payments (account_id, status, created_at DESC);