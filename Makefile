DATABASE_URL := postgres://postgres:123@localhost:5432/property?sslmode=disable
MIGRATIONS_PATH := ./migrations

ifeq ($(ENV),production)
    DATABASE_URL := postgres://postgres:123@property_db:5432/property?sslmode=disable
endif


migrate:
	go run ./cmd/migrator -database-url "$(DATABASE_URL)" -migrations-path "$(MIGRATIONS_PATH)"