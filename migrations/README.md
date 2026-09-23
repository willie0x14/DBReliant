# Migrations

`001_init.sql` creates the merchants, accounts, payments, and transfers tables,
including the composite payments index used by `GET /payments`. Apply it with:

```bash
make migrate
```

Seed data is intentionally kept outside migrations in `scripts/seed.sql`.
