APP_NAME=devices-api
API_DIR=./cmd/api
REPO_TEST_DIR=./internal/repository

DATABASE_URL?=postgres://user:pass@localhost:5432/devices?sslmode=disable
TEST_DATABASE_URL?=postgres://user:pass@localhost:5432/devices?sslmode=disable

.PHONY: help build run dev test test-unit test-repo test-cover fmt vet tidy lint clean \
	docker-up docker-down docker-build logs db-up api-up migrate-up migrate-down \
	wait-db migrate-test

help:
	@echo "Available targets:"
	@echo "  build        Build the API binary"
	@echo "  run          Run the API locally"
	@echo "  dev          Run with go run"
	@echo "  test         Run all tests"
	@echo "  test-unit    Run unit tests only"
	@echo "  test-repo    Run repository tests (with DB + migrations)"
	@echo "  test-cover   Run tests with coverage"
	@echo "  fmt          Format code"
	@echo "  vet          Run go vet"
	@echo "  tidy         Run go mod tidy"
	@echo "  lint         fmt + vet"
	@echo "  clean        Remove binary"
	@echo "  docker-up    Start all containers"
	@echo "  docker-down  Stop all containers"
	@echo "  docker-build Rebuild containers"
	@echo "  logs         Show logs"
	@echo "  db-up        Start DB only"
	@echo "  migrate-up   Run migrations"
	@echo "  test-repo    Run repository tests with setup"

build:
	go build -o $(APP_NAME) $(API_DIR)

run: build
	DATABASE_URL=$(DATABASE_URL) ./$(APP_NAME)

dev:
	DATABASE_URL=$(DATABASE_URL) go run $(API_DIR)

test:
	go test ./...

test-unit:
	go test $$(go list ./... | grep -v $(REPO_TEST_DIR))

# 🔥 IMPORTANT TARGET
test-repo: db-up wait-db migrate-test
	TEST_DATABASE_URL=$(TEST_DATABASE_URL) go test $(REPO_TEST_DIR) -v

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy

lint: fmt vet

clean:
	rm -f $(APP_NAME)
	rm -f coverage.out

docker-up:
	docker compose up

docker-down:
	docker compose down

docker-build:
	docker compose up --build

logs:
	docker compose logs -f

db-up:
	docker compose up -d db

# ⏳ Wait for Postgres to be ready
wait-db:
	@echo "Waiting for database..."
	@until docker exec $$(docker ps -qf "name=db") pg_isready -U user > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo "Database ready"

# 🚀 Run migrations against TEST DB
migrate-test:
	migrate -path migrations -database "$(TEST_DATABASE_URL)" up

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down
