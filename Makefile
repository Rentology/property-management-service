DATABASE_URL := postgres://postgres:123@localhost:5432/db?sslmode=disable
MIGRATIONS_PATH := ./migrations



migrate:
	go run ./cmd/migrator -database-url "$(DATABASE_URL)" -migrations-path "$(MIGRATIONS_PATH)"