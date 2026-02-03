# Variables
DB_CONTAINER_NAME=medibrain-db
DB_HOST=localhost
DB_PORT=5432
DB_NAME=medibrain
DB_USER=root
DB_PASSWORD=1234
SSL_MODE=disable
SQL_FILE?=./internal/database/migrations/schemas.sql

# Connection string for psql
PSQL_CMD=docker exec -i $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME)
PSQL_URI=postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(SSL_MODE)

.PHONY: help db-reset db-drop db-create db-schema db-wait db-seed run-database run-api

help:
	@echo "Available commands:"
	@echo "  make db-reset    - Drop all tables and recreate schema"
	@echo "  make db-drop     - Drop all tables in the database"
	@echo "  make db-create   - Create tables from SQL file"
	@echo "  make db-schema   - Run specific SQL file (use SQL_FILE=path/to/file.sql)"
	@echo "  make db-wait     - Wait for database to be ready"
	@echo "  make db-seed     - Seed the database with data"
	@echo "  make run-database - Run your database seeding command"
	@echo "  make run-api     - Run the API server"

# Wait for database to be ready
db-wait:
	@echo "Waiting for database to be ready..."
	@until docker exec $(DB_CONTAINER_NAME) pg_isready -U $(DB_USER) -d $(DB_NAME); do \
		echo "Waiting for database..."; \
		sleep 2; \
	done
	@echo "Database is ready!"

# Drop all tables in the database
db-drop: db-wait
	@echo "Dropping all tables..."
	@docker exec -i $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME) -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO public;"
	@echo "All tables dropped!"

# Create tables from SQL file
db-create: db-wait
	@echo "Creating tables from $(SQL_FILE)..."
	@if [ -f "$(SQL_FILE)" ]; then \
		cat $(SQL_FILE) | docker exec -i $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME); \
		echo "Tables created successfully!"; \
	else \
		echo "Error: SQL file $(SQL_FILE) not found!"; \
		echo "Please specify with: make db-create SQL_FILE=./path/to/schema.sql"; \
		exit 1; \
	fi

# Main reset command - drop and recreate
db-reset: db-drop db-create
	@echo "Database reset complete!"

# Run a specific SQL file
db-schema: db-wait
	@echo "Running SQL file: $(SQL_FILE)..."
	@if [ -f "$(SQL_FILE)" ]; then \
		cat $(SQL_FILE) | $(PSQL_CMD); \
		echo "SQL file executed successfully!"; \
	else \
		echo "Error: SQL file $(SQL_FILE) not found!"; \
		echo "Please specify with: make db-schema SQL_FILE=./path/to/file.sql"; \
		exit 1; \
	fi

# Run your Go database seeding command
db-seed: db-wait
	@echo "Seeding database with Go program..."
	@cd cmd/database && go run main.go

# Just run your database seeding command
run-database: db-wait
	@echo "Running database seeding program..."
	@cd cmd/database && go run main.go

# Run your API
run-api:
	@echo "Starting API server..."
	@cd cmd/api && go run main.go

# Connect to the database with psql
db-connect:
	@echo "Connecting to database..."
	@docker exec -it $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME)

# List all tables
db-tables: db-wait
	@echo "Listing all tables..."
	@docker exec -i $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME) -c "\dt"

# Show database info
db-info: db-wait
	@echo "Database information:"
	@docker exec -i $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME) -c "\l medibrain"

# Backup database to file
db-backup:
	@echo "Backing up database..."
	@docker exec $(DB_CONTAINER_NAME) pg_dump -U $(DB_USER) $(DB_NAME) > backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "Backup created!"

# Restore database from backup
db-restore:
	@echo "Restoring database from backup..."
	@if [ -f "$(BACKUP_FILE)" ]; then \
		cat $(BACKUP_FILE) | docker exec -i $(DB_CONTAINER_NAME) psql -U $(DB_USER) -d $(DB_NAME); \
		echo "Database restored from $(BACKUP_FILE)"; \
	else \
		echo "Error: Backup file $(BACKUP_FILE) not found!"; \
		echo "Usage: make db-restore BACKUP_FILE=./backup.sql"; \
		exit 1; \
	fi

# Development: run both API and database seed (in background)
dev:
	@echo "Starting development environment..."
	@echo "Run 'make run-api' in one terminal"
	@echo "Run 'make run-database' in another terminal"

# Quick reset and seed (useful for development)
fresh: db-reset db-seed
	@echo "Fresh database ready!"

# If you're using sqlx and want to run migrations
db-migrate: db-wait
	@echo "Running migrations..."
	@if [ -f "cmd/database/main.go" ]; then \
		cd cmd/database && go run main.go migrate; \
	else \
		echo "Database program not found at cmd/database/main.go"; \
	fi
