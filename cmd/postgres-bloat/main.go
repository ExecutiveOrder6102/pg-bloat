package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/term"

	"github.com/ExecutiveOrder6102/postgres-bloat/internal/bloat"
	"github.com/ExecutiveOrder6102/postgres-bloat/internal/config"
	"github.com/ExecutiveOrder6102/postgres-bloat/internal/connect"
	"github.com/ExecutiveOrder6102/postgres-bloat/internal/output"
)

func main() {
	cfg, err := config.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if cfg.DebugSQL {
		fmt.Fprintln(os.Stderr, "Bloat SQL:")
		fmt.Fprintln(os.Stderr, bloat.IndexBloatQuery())
	}

	var db *sql.DB
	var cleanup func()
	if cfg.CloudSQLInstance != "" {
		if err := maybePromptPassword(&cfg); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		var err error
		db, cleanup, err = connect.OpenCloudSQL(ctx, connect.CloudSQLOptions{
			Instance: cfg.CloudSQLInstance,
			User:     cfg.User,
			Password: cfg.Password,
			Database: cfg.Database,
			SSLMode:  cfg.SSLMode,
			IAMAuthN: cfg.CloudSQLIAMAuthN,
			IPType:   cfg.CloudSQLIPType,
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "open cloud sql:", err)
			os.Exit(1)
		}
	} else {
		if cfg.DSN == "" {
			if err := maybePromptPassword(&cfg); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}
		}
		dsn, err := resolveDSN(cfg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		db, err = sql.Open("pgx", dsn)
		if err != nil {
			fmt.Fprintln(os.Stderr, "open database:", err)
			os.Exit(1)
		}
	}
	defer db.Close()
	if cleanup != nil {
		defer cleanup()
	}

	if err := db.PingContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, "ping database:", err)
		os.Exit(1)
	}

	serverVersion, err := bloat.ServerVersion(ctx, db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "read server version:", err)
		os.Exit(1)
	}

	allIndexes, err := bloat.FetchIndexBloat(ctx, db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "fetch bloat:", err)
		os.Exit(1)
	}

	filteredIndexes := bloat.FilterIndexes(allIndexes, cfg.MinBloatPct, cfg.MinBloatBytes, cfg.IncludeSystemSchemas)
	bloat.SortIndexes(filteredIndexes)
	filteredIndexes = bloat.ApplyLimitIndexes(filteredIndexes, cfg.Limit)

	tables := bloat.RollupTables(allIndexes)
	tables = bloat.FilterTables(tables, cfg.MinBloatPct, cfg.MinBloatBytes, cfg.IncludeSystemSchemas)
	bloat.SortTables(tables)
	tables = bloat.ApplyLimitTables(tables, cfg.Limit)

	concurrentSupported := serverVersion >= 120000

	if cfg.StaleStatsDays > 0 {
		staleStats, err := bloat.FetchStaleStats(ctx, db, time.Duration(cfg.StaleStatsDays)*24*time.Hour)
		if err != nil {
			fmt.Fprintln(os.Stderr, "warn: failed to check stats freshness:", err)
		} else if len(staleStats) > 0 {
			fmt.Fprintf(os.Stderr, "Warning: %d tables have stale or missing stats (older than %d days). Consider ANALYZE for more accurate bloat estimates.\n", len(staleStats), cfg.StaleStatsDays)
		}
	}

	if err := output.Write(cfg, filteredIndexes, tables, concurrentSupported); err != nil {
		fmt.Fprintln(os.Stderr, "write output:", err)
		os.Exit(1)
	}
}

func resolveDSN(cfg config.Config) (string, error) {
	if cfg.DSN != "" {
		return cfg.DSN, nil
	}
	return config.BuildDSN(cfg)
}

func promptPassword() (string, error) {
	stdin := int(os.Stdin.Fd())
	if !term.IsTerminal(stdin) {
		return "", fmt.Errorf("password required: provide --password or use --no-password")
	}
	fmt.Fprint(os.Stderr, "Password: ")
	password, err := term.ReadPassword(stdin)
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}
	return string(password), nil
}

func maybePromptPassword(cfg *config.Config) error {
	if cfg.Password != "" || cfg.NoPassword {
		return nil
	}
	password, err := promptPassword()
	if err != nil {
		return err
	}
	if password == "" {
		cfg.NoPassword = true
		return nil
	}
	cfg.Password = password
	return nil
}
