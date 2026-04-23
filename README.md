# Devices API

A small REST API for managing devices. It supports creating, listing, updating and deleting devices, with simple business rules around device state.

The project is written in Go, uses Gin for HTTP, pgx for PostgreSQL access, and golang-migrate for schema migrations.

## Features

* Create devices with name, brand and state
* List devices with optional filters (brand and state)
* Update devices with validation rules
* Delete devices with safety checks
* Structured logging with zerolog
* Docker setup for local development
* Unit and integration tests

## Requirements

* Go 1.22 or newer
* Docker and Docker Compose
* migrate CLI (for running migrations from Makefile)

Install migrate if you do not have it:

```bash
brew install golang-migrate
```

or

```bash
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

## Running with Docker

Start everything:

```bash
make docker-up
```

This will build the API and start both the API and PostgreSQL.

The API will be available at:

```
http://localhost:8080
```

To stop everything:

```bash
make docker-down
```

## Running locally (without Docker)

Start only the database:

```bash
make db-up
```

Run migrations:

```bash
make migrate-up
```

Start the API:

```bash
make dev
```

Or build and run:

```bash
make build
make run
```

## Environment variables

The API uses:

```
DATABASE_URL=postgres://user:pass@localhost:5432/devices?sslmode=disable
```

Repository tests use a separate database in the same Postgres instance:

```
TEST_DATABASE_URL=postgres://user:pass@localhost:5432/devices_test?sslmode=disable
```

The `devices_test` database is created automatically via the init script mounted in the Postgres container.

## API endpoints

Create device

```
POST /devices
```

Get device

```
GET /devices/:id
```

List devices

```
GET /devices?brand=Apple&state=available
```

Update device

```
PATCH /devices/:id
```

Delete device

```
DELETE /devices/:id
```

## Example requests

Create a device:

```bash
curl -X POST http://localhost:8080/devices \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15",
    "brand": "Apple",
    "state": "available"
  }'
```

List all devices:

```bash
curl http://localhost:8080/devices
```

Filter by brand:

```bash
curl "http://localhost:8080/devices?brand=Apple"
```

Filter by state:

```bash
curl "http://localhost:8080/devices?state=available"
```

Update a device:

```bash
curl -X PATCH http://localhost:8080/devices/<device-id> \
  -H "Content-Type: application/json" \
  -d '{
    "state": "inactive"
  }'
```

Delete a device:

```bash
curl -X DELETE http://localhost:8080/devices/<device-id>
```

## Business rules

* A device in use cannot have its name or brand updated
* A device in use cannot be deleted
* State must be one of: available, in-use, inactive

## Testing

Run all tests:

```bash
make test
```

Run only unit tests:

```bash
make test-unit
```

Run repository tests against a real database:

```bash
make test-repo
```

This will:

* start the database
* wait until it is ready
* run migrations on the test database
* execute repository tests

## Project structure

```
cmd/api             entrypoint
internal/handler    HTTP layer
internal/service    business logic
internal/repository database access
internal/model      domain models
internal/dto        request/response types
internal/db         database setup and migrations
internal/middleware logging middleware
```

## Notes

* Repository tests use a separate `devices_test` database
* Migrations are located in the `migrations` directory
* Logging is configured via zerolog and can be adjusted with environment variables

