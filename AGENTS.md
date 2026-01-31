# AGENTS

This repository is a Go CLI for identifying Postgres index bloat and ranking tables for reindexing.

## Project layout

- `cmd/postgres-bloat`: CLI entrypoint.
- `internal/bloat`: Bloat query, filtering, and rollups.
- `internal/output`: Console/CSV output.
- `docker/`: Local test environment SQL.

## Running tests

There are no automated tests yet. When they exist, use:

```bash
go test ./...
```

## Local test environment

Start Postgres 18 and generate sample bloat:

```bash
docker compose up -d postgres
docker compose run --rm bloat-gen
```

Then run the CLI:

```bash
go build -o postgres-bloat ./cmd/postgres-bloat
./postgres-bloat --dsn "postgres://postgres:postgres@localhost:5432/bloatdb?sslmode=disable"
```
