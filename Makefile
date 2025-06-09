.PHONY: test generate clean docker-build docker-up docker-down dev

# Variables
BINARY_NAME=themenu
DOCKER_COMPOSE=docker-compose.yml

# Comandos principales
test:
	go test -v ./...

generate:
	sqlc generate
	templ generate

clean:
	rm -rf bin/
	go clean

# Comandos de Docker
docker-build:
	docker compose -f $(DOCKER_COMPOSE) build --no-cache

docker-up: docker-build
	docker compose -f $(DOCKER_COMPOSE) up -d

docker-down:
	docker compose -f $(DOCKER_COMPOSE) down

# Comandos de desarrollo
dev: docker-up

# Comandos de base de datos
db-migrate:
	psql -h localhost -U postgres -d themenu -f schema.sql

db-seed:
	psql -h localhost -U postgres -d themenu -f seed.sql
