package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DSN                  string
	Host                 string
	Port                 int
	User                 string
	Password             string
	Database             string
	SSLMode              string
	NoPassword           bool
	DebugSQL             bool
	StaleStatsDays       int
	CloudSQLInstance     string
	CloudSQLIAMAuthN     bool
	CloudSQLIPType       string
	Output               string
	OutputFile           string
	MinBloatPct          float64
	MinBloatBytes        int64
	Limit                int
	IncludeSystemSchemas bool
}

func Parse(args []string) (Config, error) {
	fs := flag.NewFlagSet("postgres-bloat", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var cfg Config
	fs.StringVar(&cfg.DSN, "dsn", "", "Postgres connection string")
	fs.StringVar(&cfg.Host, "host", "", "Postgres host")
	fs.IntVar(&cfg.Port, "port", 5432, "Postgres port")
	fs.StringVar(&cfg.User, "user", "", "Postgres user")
	fs.StringVar(&cfg.Password, "password", "", "Postgres password")
	fs.StringVar(&cfg.Database, "dbname", "", "Postgres database name")
	fs.StringVar(&cfg.SSLMode, "sslmode", "disable", "SSL mode")
	fs.BoolVar(&cfg.NoPassword, "no-password", false, "Do not prompt for a password")
	fs.BoolVar(&cfg.DebugSQL, "debug-sql", false, "Print SQL used for bloat detection")
	fs.IntVar(&cfg.StaleStatsDays, "stale-stats-days", 7, "Warn if stats are older than this many days")
	fs.StringVar(&cfg.CloudSQLInstance, "cloudsql-instance", "", "Cloud SQL instance connection name (project:region:instance)")
	fs.BoolVar(&cfg.CloudSQLIAMAuthN, "cloudsql-iam-authn", true, "Use IAM database authentication")
	fs.StringVar(&cfg.CloudSQLIPType, "cloudsql-ip-type", "public", "Cloud SQL IP type: public or private")
	fs.StringVar(&cfg.Output, "output", "console", "Output format: console or csv")
	fs.StringVar(&cfg.OutputFile, "output-file", "", "Write output to file instead of stdout")
	fs.Float64Var(&cfg.MinBloatPct, "min-bloat-pct", 20, "Minimum bloat percent")
	fs.Int64Var(&cfg.MinBloatBytes, "min-bloat-bytes", 0, "Minimum bloat size in bytes")
	fs.IntVar(&cfg.Limit, "limit", 50, "Maximum rows per section")
	fs.BoolVar(&cfg.IncludeSystemSchemas, "include-system-schemas", false, "Include pg_catalog and other system schemas")

	if err := fs.Parse(args); err != nil {
		return Config{}, err
	}

	if cfg.CloudSQLInstance != "" && cfg.DSN != "" {
		return Config{}, errors.New("use either --dsn or --cloudsql-instance, not both")
	}

	if cfg.DSN == "" {
		missing := []string{}
		if cfg.Host == "" && cfg.CloudSQLInstance == "" {
			missing = append(missing, "host")
		}
		if cfg.User == "" {
			missing = append(missing, "user")
		}
		if cfg.Database == "" {
			missing = append(missing, "dbname")
		}
		if len(missing) > 0 {
			return Config{}, fmt.Errorf("dsn or %s required", strings.Join(missing, ", "))
		}
	}

	if cfg.Output != "console" && cfg.Output != "csv" {
		return Config{}, fmt.Errorf("unsupported output: %s", cfg.Output)
	}

	if cfg.Limit < 1 {
		return Config{}, errors.New("limit must be >= 1")
	}

	if cfg.MinBloatPct < 0 {
		return Config{}, errors.New("min-bloat-pct must be >= 0")
	}

	if cfg.MinBloatBytes < 0 {
		return Config{}, errors.New("min-bloat-bytes must be >= 0")
	}

	if cfg.StaleStatsDays < 0 {
		return Config{}, errors.New("stale-stats-days must be >= 0")
	}

	if cfg.CloudSQLIPType != "public" && cfg.CloudSQLIPType != "private" {
		return Config{}, errors.New("cloudsql-ip-type must be public or private")
	}

	return cfg, nil
}
