include ./database/.env

REDIS_URL = redis://$(REDIS_USER):$(REDIS_PASSWORD)@$(REDIS_HOST):$(REDIS_PORT)
DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

MIGRATE_PATH=database/migrations
SEEDER_PATH=database/seeders

migrate-create:
	migrate create -ext sql -dir $(MIGRATE_PATH) -seq create_$(NAME)_table

migrate-up:
	migrate -database $(DB_URL) -path $(MIGRATE_PATH) up

migrate-down:
	migrate -database $(DB_URL) -path $(MIGRATE_PATH) down $(N)

db-seed:
	cat $(SEEDER_PATH)/*.sql | psql $(DB_URL)

up:
	docker compose up -d --build
	@echo "Waiting for PostgreSQL to be ready..."
	@until docker exec db-erp pg_isready -U $(POSTGRES_USER) -d $(POSTGRES_DB) > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo "Running migrations..."
	@make migrate-up
	@echo "Running seeders..."
	@make db-seed
	@echo "Database setup completed successfully!"

down:
	docker compose down --rmi local

ps:
	docker compose ps -a

swag:
	swag init -g backend/cmd/main.go -o docs