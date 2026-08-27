![CodeRabbit Pull Request Reviews](https://img.shields.io/coderabbit/prs/github/devrapture/stack-trace?utm_source=oss&utm_medium=github&utm_campaign=devrapture%2Fstack-trace&labelColor=171717&color=FF570A&link=https%3A%2F%2Fcoderabbit.ai&label=CodeRabbit+Reviews)
# Stack Trace

Stack Trace is a Go service built with Gin, Zap, and PostgreSQL. The repository
includes separate development and production container targets in one Docker
Compose file.

## Requirements

- Docker with Docker Compose 2.22 or newer
- Go 1.26.6 for running the backend outside Docker
- Air for hot reload when running the backend outside Docker
- `make` for the convenience commands

## Quick start

### Development

Start PostgreSQL and the API with Compose Watch and Air hot reload:

```sh
make dev
```

The API is available at `http://localhost:8080`. Source changes under
`apps/backend` are synchronized into the container, where Air rebuilds and
restarts the API. Changes to `go.mod` or `go.sum` rebuild the development image.

Press `Ctrl+C` to stop the attached development stack.

### Production

Build and start the production containers in the background:

```sh
make prod
```

The production image contains only the compiled API and its runtime
dependencies. The API runs as a non-root user, and PostgreSQL is not published
to the host.

Stop either stack with:

```sh
make down
```

## Compose profiles

The development and production configurations live in `docker-compose.yml`:

```sh
# Development: Compose Watch + Air
docker compose --profile dev up --build --watch --remove-orphans

# Production: compiled binary
docker compose --profile prod up -d --build
```

Development and production use separate PostgreSQL volumes. The development
profile keeps compatibility with the existing `stack-trace_wordsearch-pgdata`
volume.

## Configuration

The following environment variables can be set in your shell before running
Compose:

| Variable | Default | Purpose |
| --- | --- | --- |
| `APP_PORT` | `8080` | Host port mapped to the API |
| `POSTGRES_PORT` | `5434` | Development PostgreSQL host port (avoids colliding with local Postgres on 5432 and other Compose stacks on 5433) |
| `POSTGRES_USER` | `wordsearch` | PostgreSQL user |
| `POSTGRES_PASSWORD` | `wordsearch` | PostgreSQL password |
| `POSTGRES_DB` | `wordsearch` | PostgreSQL database |
| `DATABASE_URL` | Generated from the PostgreSQL settings | Production database connection string |

For production, provide secure PostgreSQL credentials and, when using an
external database, set `DATABASE_URL` explicitly:

```sh
POSTGRES_USER=stacktrace \
POSTGRES_PASSWORD='replace-me' \
POSTGRES_DB=stacktrace \
DATABASE_URL='postgresql://stacktrace:replace-me@database:5432/stacktrace?sslmode=require' \
make prod
```

The backend also provides `apps/backend/.env.example` for running directly on
the host. Copy it to `apps/backend/.env` and adjust the values as needed.

## API

Check that the API is running:

```sh
curl http://localhost:8080/api/v1/healthz
```

Expected response:

```json
{"message":"OK"}
```

Readiness check (returns 200 when the database is reachable):

```sh
curl http://localhost:8080/api/v1/readyz
```

## Local backend development

To run the API without Docker, start PostgreSQL and then launch Air from the
backend directory:

```sh
make dev-up
cd apps/backend
air
```

`make dev-up` starts Postgres and creates the Atlas development database
(`wordsearch_dev`) if it is missing. Make sure `DATABASE_URL` in
`apps/backend/.env` points to the published development database port
(`5434` by default). Host-side tools (Atlas, `air` outside Compose) must
use that port; `localhost:5432` is often a separate Homebrew or
Postgres.app cluster that does not have the `wordsearch` role.

### Database migrations

Install the [Atlas CLI](https://atlasgo.io/getting-started#installation), run
`make dev-up`, then:

```sh
make migrate-diff NAME=<description>  # generate a migration from GORM models
make migrate-up                       # apply pending migrations
make migrate-down                     # rollback the last migration
make migrate-status                   # show migration status
```

Production equivalents (`migrate-prod-up`, `migrate-prod-down`,
`migrate-prod-status`) read `DATABASE_URL` and `DEV_DATABASE_URL` from the
environment.

## Database helpers

```sh
make db-up   # Start Postgres
make db-down # Stop Postgres
make psql    # Open a psql shell against local Postgres
```

## Quality checks

```sh
make tidy           # Sync go.mod/go.sum
make fmt            # Format Go source
make vet            # Run go vet
make lint           # Run the pinned golangci-lint version
make test           # Run tests
make test-verbose   # Run tests with verbose output
make test-coverage  # Run tests with coverage
```

The continuous integration workflow runs backend tests, `go vet`, and
golangci-lint.

## Project structure

```text
.
├── apps/
│   ├── backend/
│   │   ├── cmd/api/           # API entry point
│   │   ├── cmd/worker/        # Worker entry point (not yet implemented)
│   │   └── internal/platform/ # Configuration, HTTP, logging, and utilities
│   └── frontend/              # Frontend application (scaffolded)
├── docker-compose.yml         # Development and production profiles
└── Makefile                   # Development and quality commands
```
