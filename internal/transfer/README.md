# Transfer Package

The package implements the transactional account transfer flow used by
`POST /transfers`.

- Validates the amount and rejects same-account transfers before opening a
  transaction.
- Locks both accounts with `SELECT ... FOR UPDATE` in ascending account ID order.
- Checks the source balance while holding both locks.
- Updates both balances and inserts the completed transfer in one transaction.
- Rolls back the transaction on query, update, insert, or commit errors.

The service returns domain errors for invalid amounts, same-account transfers,
missing accounts, and insufficient funds. HTTP status mapping remains in
`internal/http`.
