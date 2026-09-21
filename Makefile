COMPOSE ?= docker compose
PAYMENT_COUNT ?= 100000

.PHONY: setup up down reset migrate seed build run test psql

setup: up

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

# Destructive: removes the local PostgreSQL volume and recreates PostgreSQL.
reset:
	$(COMPOSE) down -v
	$(COMPOSE) up -d

migrate:
	$(COMPOSE) exec -T postgres psql -v ON_ERROR_STOP=1 -U $${POSTGRES_USER:-dbreliant} -d $${POSTGRES_DB:-dbreliant} < migrations/001_init.sql

seed:
	@case "$(PAYMENT_COUNT)" in ''|*[!0-9]*) echo "PAYMENT_COUNT must be a positive integer"; exit 1;; esac
	@test "$(PAYMENT_COUNT)" -gt 0 || (echo "PAYMENT_COUNT must be a positive integer"; exit 1)
	$(COMPOSE) exec -T postgres psql -v ON_ERROR_STOP=1 -v payment_count="$(PAYMENT_COUNT)" -U $${POSTGRES_USER:-dbreliant} -d $${POSTGRES_DB:-dbreliant} < scripts/seed.sql

build:
	go build ./...

run:
	go run ./cmd/api

test:
	go test ./...

psql:
	$(COMPOSE) exec postgres psql -U $${POSTGRES_USER:-dbreliant} -d $${POSTGRES_DB:-dbreliant}
