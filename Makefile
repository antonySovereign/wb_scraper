.PHONY: up down migrate-up migrate-down migrate-create migrate-force migrate-goto

include .env
export 

CURR_DIR := $(shell pwd)
MIGRATE := docker compose run --rm migrate

db_up:
	docker compose up -d

db_down:
	docker compose down

migrate_up:
	$(MIGRATE) up

migrate_down:
	$(MIGRATE) down 1

migrate_create:
	@docker run --rm \
	-u $$(id -u):$$(id -g) \
	-v $(CURR_DIR)/migrations:/migrations migrate/migrate:v4.19.1 \
	create \
	-ext sql \
	-dir /migrations \
	-seq $(name)

migrate_force:
	$(MIGRATE) force $(v)

migrate_goto:
	$(MIGRATE) goto $(v)

db:
	docker exec -it postgres_db psql -U $(DB_USER) -d $(DB_NAME)

run:
	@go run cmd/scraper/main.go

build:
	@go build -o scraper cmd/scraper/main.go

clean:
	@docker compose down --remove-orphans -v
	@docker container prune -f

db_reset: clean db_up
