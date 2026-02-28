DB_CONTAINER_NAME=medibrain-db
DB_HOST=localhost
DB_PORT=5432
DB_NAME=medibrain
DB_USER=root
DB_PASSWORD=1234
SSL_MODE=disable
SQL_FILE?=./internal/database/migrations/schemas.sql
MIGRATION_NAME?=db_update
STEPS?=1
VERSION?=1

PSQL_CMD=docker exec -i $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME)
PSQL_URI=postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE)

.PHONY: help db-reset db-drop db-create db-schema db-wait db-seed run-database run-api

help:
	@echo "Available commands:"
	@echo "  make db-create-migrate MIGRATION_NAME=migration_name - Create a new migration files (up and down)"
	@echo "  make db-migrate-up STEPS=<number>  - run migrations by step (default is 1, you can provide empty after argument to up all )"
	@echo "  make db-migrate-down STEPS=<number> - down migrations by step (default is 1, you can provide empty space after argument to down all)"
	@echo "  make db-migrate-force - force migration to last clean migrate"
	@echo "  make qdrant-seed - Drops any existing collections and reconstruct collections from go"
	@echo "  make db-wait     - Wait for database to be ready"
	@echo "  make db-seed     - Seed the database with data"
	@echo "  make run-database - Run your database seeding command"
	@echo "  make db-tables - List all database tables"

db-create-migrate: db-wait
	@echo "Creating migration..."
	@if [ -d "internal/database/migrate" ]; then \
		cd ./internal/database/migrate && migrate create -ext sql -dir . -seq $(MIGRATION_NAME); \
	else \
		echo "migrate folder not found at internal/database/migrate"; \
	fi

db-migrate-up: db-wait
	@echo "Running migrations..."
	@if [ -d "internal/database/migrate" ]; then \
		cd internal/database/migrate && migrate -database $(PSQL_URI) -path ./ up $(STEPS); \
		echo "success!"; \
	else \
		echo "migrate folder not found at internal/database/migrate"; \
	fi

db-migrate-down: db-wait
	@echo "reverting last migration..."
	@if [ -d "internal/database/migrate" ]; then \
		cd ./internal/database/migrate && migrate -database $(PSQL_URI) -path . down $(STEPS); \
	else \
		echo "migrate folder not found at internal/database/migrate"; \
	fi

db-migrate-force: db-wait
	@echo "rolling back to last good migration..."
	@if [ -d "internal/database/migrate" ]; then \
		cd ./internal/database/migrate && migrate -database $(PSQL_URI) -path . force $(VERSION); \
	else \
		echo "migrate folder not found at internal/database/migrate"; \
	fi

qdrant-seed:
	@echo "Dropping and recreating qdrant collections with Go program..."
	@cd cmd/qdrant && go run main.go

db-seed: db-wait
	@echo "Seeding database with Go program..."
	@cd cmd/database && go run main.go

db-tables: db-wait
	@echo "Listing all tables..."
	@docker exec -i $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME) -c "\dt"

db-info: db-wait
	@echo "Database information:"
	@docker exec -i $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME) -c "\l medibrain"

# db-backup:
# 	@echo "Backing up database..."
# 	@docker exec $(DB_CONTAINER_NAME) pg_dump -U $(DB_USER) $(DB_NAME) > backup_$(shell date +%Y%m%d_%H%M%S).sql
# 	@echo "Backup created!"

# fresh: db-reset db-seed
# 	@echo "Fresh database ready!"

db-wait:
	@echo "Waiting for database to be ready..."
	@until docker exec $(DB_CONTAINER_NAME) pg_isready -U $(DB_USER) -d $(DB_NAME); do \
		echo "Waiting for database..."; \
		sleep 2; \
	done
	@echo "Database is ready!"
