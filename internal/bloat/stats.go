package bloat

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type StaleStats struct {
	Schema      string
	Table       string
	LastAnalyze sql.NullTime
}

const staleStatsQuery = `
SELECT
    schemaname,
    relname,
    COALESCE(last_analyze, last_autoanalyze) AS last_analyze
FROM pg_stat_user_tables
WHERE COALESCE(last_analyze, last_autoanalyze) IS NULL
   OR COALESCE(last_analyze, last_autoanalyze) < now() - $1::interval
ORDER BY schemaname, relname
`

func FetchStaleStats(ctx context.Context, db *sql.DB, maxAge time.Duration) ([]StaleStats, error) {
	interval := formatInterval(maxAge)
	rows, err := db.QueryContext(ctx, staleStatsQuery, interval)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []StaleStats
	for rows.Next() {
		var row StaleStats
		if err := rows.Scan(&row.Schema, &row.Table, &row.LastAnalyze); err != nil {
			return nil, err
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func formatInterval(value time.Duration) string {
	if value <= 0 {
		return "0 seconds"
	}
	hours := int(value.Hours())
	if hours >= 24 {
		days := hours / 24
		return fmt.Sprintf("%d days", days)
	}
	if hours >= 1 {
		return fmt.Sprintf("%d hours", hours)
	}
	minutes := int(value.Minutes())
	if minutes >= 1 {
		return fmt.Sprintf("%d minutes", minutes)
	}
	return "1 minute"
}
