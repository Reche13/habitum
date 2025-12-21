.PHONY: run tidy migrate-up migrate-down migrate-status migrate-create

MIGRATIONS_DIR=internal/database/migrations

run:
	go run ./cmd/habitum

tidy:
	@echo "Formatting .go files..."
	go fmt ./...
	@echo "Tidying module dependencies..."
	go mod tidy
	@echo "Verifying dependencies..."
	go mod verify

migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres $$DATABASE_URL up

migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres $$DATABASE_URL down

migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres $$DATABASE_URL status

migrate-create:
ifeq ($(strip $(name)),)
	$(error name is required)
endif
	goose -dir $(MIGRATIONS_DIR) create $(name) sql