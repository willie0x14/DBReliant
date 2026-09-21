COMPOSE ?= docker compose

.PHONY: setup up down reset build run test psql

setup: up

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

# Destructive: removes the local PostgreSQL volume and recreates PostgreSQL.
reset:
	$(COMPOSE) down -v
	$(COMPOSE) up -d

build:
	go build ./...

run:
	go run ./cmd/api

test:
	go test ./...

psql:
	$(COMPOSE) exec postgres psql -U $${POSTGRES_USER:-dbreliant} -d $${POSTGRES_DB:-dbreliant}
