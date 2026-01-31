# postgres-bloat

Find bloated Postgres indexes, rank the highest-value tables to reindex, and generate `REINDEX` commands. Output can be console or CSV.

This tool uses the AWS Prescriptive Guidance btree bloat estimate query (adapted for general Postgres). It is an estimate and depends on up-to-date stats (`ANALYZE`).

## Build

```bash
go build -o postgres-bloat ./cmd/postgres-bloat
```

## Usage

```bash
./postgres-bloat --dsn "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
```

Or use psql-style flags:

```bash
./postgres-bloat --host localhost --port 5432 --user postgres --dbname bloatdb
```

Console output is the default. For CSV:

```bash
./postgres-bloat --dsn "postgres://user:pass@localhost:5432/dbname?sslmode=disable" --output csv --output-file bloat.csv
```

## Connection examples

Regular Postgres (DSN):

```bash
./postgres-bloat --dsn "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
```

Regular Postgres (psql-style flags):

```bash
./postgres-bloat --host localhost --port 5432 --user postgres --dbname bloatdb
```

Cloud SQL (Connector + ADC):

```bash
./postgres-bloat \
  --cloudsql-instance "project:region:instance" \
  --user user \
  --dbname dbname \
  --cloudsql-ip-type public \
  --cloudsql-iam-authn true
```

### Flags

- `--dsn`: Postgres connection string.
- `--host`: Postgres host (when not using `--dsn`).
- `--port`: Postgres port (default: 5432).
- `--user`: Postgres user.
- `--password`: Postgres password (if omitted, you will be prompted unless `--no-password` is set; empty input skips password).
- `--dbname`: Postgres database name.
- `--sslmode`: SSL mode (default: `disable`).
- `--no-password`: Do not prompt for a password (useful for IAM auth).
- `--output`: `console` or `csv` (default: `console`).
- `--output-file`: Write output to a file instead of stdout.
- `--min-bloat-pct`: Minimum bloat percent (default: 20).
- `--min-bloat-bytes`: Minimum bloat size in bytes (default: 0).
- `--limit`: Max rows per section (default: 50).
- `--include-system-schemas`: Include `pg_catalog` and other system schemas.
- `--debug-sql`: Print SQL used for bloat detection.
- `--stale-stats-days`: Warn if stats are older than this many days (default: 7).
- `--cloudsql-instance`: Cloud SQL instance connection name (project:region:instance).
- `--cloudsql-iam-authn`: Use IAM database authentication (default: true).
- `--cloudsql-ip-type`: Cloud SQL IP type (`public` or `private`, default: `public`).

### Output

Console output includes two sections:

- Index bloat detail with reindex commands.
- Table rollup with the highest-value tables to reindex.

CSV output includes both record types with a `record_type` column (`index` or `table`).

## Limitations

- Estimates bloat for btree indexes only.
- Relies on `pg_stats`; run `ANALYZE` or `VACUUM (ANALYZE)` for accurate results.
- `REINDEX ... CONCURRENTLY` is only available for Postgres 12+.

## Cloud SQL Connector

If you pass `--cloudsql-instance`, the tool uses the Cloud SQL Go Connector to connect directly. This relies on Application Default Credentials (ADC).

ADC setup options:

- `gcloud auth application-default login`
- `GOOGLE_APPLICATION_CREDENTIALS` pointing at a service account key

See the connection examples above for a full Cloud SQL command.

## Docker test environment

Start Postgres 18 and generate sample bloat:

```bash
docker compose up -d postgres
docker compose run --rm bloat-gen
```

To scale the seed data and bloat passes:

```bash
BLOAT_MULTIPLIER=10 docker compose run --rm bloat-gen
```

Then run the tool against the local database:

```bash
./postgres-bloat --dsn "postgres://postgres:postgres@localhost:5432/bloatdb?sslmode=disable"
```
