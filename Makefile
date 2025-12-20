.PHONY: help setup run test migrate

help:
	@echo "Available commands:"
	@echo "  make setup    - Install dependencies"
	@echo "  make migrate  - Run database migrations"
	@echo "  make run      - Run the server"
	@echo "  make test     - Run tests"

setup:
	go mod download

migrate:
	@echo "Running migrations..."
	@psql -d assignment -f migrations/001_initial_schema.sql || echo "Please ensure PostgreSQL is running and database 'assignment' exists"

run:
	go run cmd/server/main.go

test:
	go test ./...

