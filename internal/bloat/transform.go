package bloat

import (
	"sort"
	"strings"
)

type TableBloat struct {
	Schema    string
	Table     string
	RealSize  int64
	BloatSize int64
	BloatPct  float64
}

func FilterIndexes(rows []IndexBloat, minPct float64, minBytes int64, includeSystem bool) []IndexBloat {
	filtered := make([]IndexBloat, 0, len(rows))
	for _, row := range rows {
		if !includeSystem && isSystemSchema(row.Schema) {
			continue
		}
		if row.BloatPct < minPct {
			continue
		}
		if row.BloatSize < minBytes {
			continue
		}
		filtered = append(filtered, row)
	}
	return filtered
}

func RollupTables(rows []IndexBloat) []TableBloat {
	agg := map[string]*TableBloat{}
	for _, row := range rows {
		key := row.Schema + "." + row.Table
		entry, ok := agg[key]
		if !ok {
			entry = &TableBloat{Schema: row.Schema, Table: row.Table}
			agg[key] = entry
		}
		entry.RealSize += row.RealSize
		entry.BloatSize += row.BloatSize
	}

	results := make([]TableBloat, 0, len(agg))
	for _, entry := range agg {
		if entry.RealSize > 0 {
			entry.BloatPct = float64(entry.BloatSize) * 100 / float64(entry.RealSize)
		}
		results = append(results, *entry)
	}
	return results
}

func FilterTables(rows []TableBloat, minPct float64, minBytes int64, includeSystem bool) []TableBloat {
	filtered := make([]TableBloat, 0, len(rows))
	for _, row := range rows {
		if !includeSystem && isSystemSchema(row.Schema) {
			continue
		}
		if row.BloatPct < minPct {
			continue
		}
		if row.BloatSize < minBytes {
			continue
		}
		filtered = append(filtered, row)
	}
	return filtered
}

func SortIndexes(rows []IndexBloat) {
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].BloatSize == rows[j].BloatSize {
			return rows[i].BloatPct > rows[j].BloatPct
		}
		return rows[i].BloatSize > rows[j].BloatSize
	})
}

func SortTables(rows []TableBloat) {
	sort.SliceStable(rows, func(i, j int) bool {
		if rows[i].BloatSize == rows[j].BloatSize {
			return rows[i].BloatPct > rows[j].BloatPct
		}
		return rows[i].BloatSize > rows[j].BloatSize
	})
}

func ApplyLimitIndexes(rows []IndexBloat, limit int) []IndexBloat {
	if limit <= 0 || len(rows) <= limit {
		return rows
	}
	return rows[:limit]
}

func ApplyLimitTables(rows []TableBloat, limit int) []TableBloat {
	if limit <= 0 || len(rows) <= limit {
		return rows
	}
	return rows[:limit]
}

func isSystemSchema(schema string) bool {
	if schema == "pg_catalog" || schema == "information_schema" {
		return true
	}
	return strings.HasPrefix(schema, "pg_")
}
