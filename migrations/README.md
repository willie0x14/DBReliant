# Migrations

`001_init.sql` creates the initial merchants, accounts, payments, and transfers
tables. Apply it with:

```bash
make migrate
```

Seed data is intentionally kept outside migrations in `scripts/seed.sql`.
