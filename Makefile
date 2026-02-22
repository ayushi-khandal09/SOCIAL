include .env
MIGRATIONS_PATH = ./cmd/migrate/migrations

.PHONY: migrate-create
migration:
	@migrate create -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY: test
test:
	@go test -v ./...

.PHONY: seed
seed: 
	@DB_ADDR=$(DB_ADDR) go run cmd/migrate/seed/main.go

.PHONY: gen-docs
gen-docs:
	rm -rf docs
	@swag init -g cmd/api/main.go -d . -o docs -parseInternal

.PHONY: migrate-version migrate-force

migrate-version:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) -verbose version

migrate-force:
	@V=$(V); if [ -z "$$V" ]; then echo "usage: make migrate-force V=<version>"; exit 2; fi; \
	migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) -verbose force $$V

.PHONY: migrate-down-all migrate-reset-all migrate-drop

migrate-down-all:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) -verbose down -all

migrate-reset-all:
	@echo "→ Resetting DB to 0…"
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) -verbose down -all || \
	migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) -verbose drop -f
	@echo "→ Applying all migrations…"
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) -verbose up
	@$(MAKE) migrate-version

migrate-drop:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DB_ADDR) -verbose drop -f
