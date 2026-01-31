package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/ExecutiveOrder6102/postgres-bloat/internal/bloat"
)

func WriteConsole(w io.Writer, indexes []bloat.IndexBloat, tables []bloat.TableBloat, concurrentSupported bool) error {
	fmt.Fprintln(w, "Index Bloat")
	indexRows := make([][]string, 0, len(indexes))
	for _, row := range indexes {
		indexRows = append(indexRows, []string{
			row.Schema,
			row.Table,
			row.Index,
			formatBytes(row.RealSize),
			formatBytes(row.BloatSize),
			formatPct(row.BloatPct),
			bloat.ReindexIndexCommand(row.Schema, row.Index, concurrentSupported),
		})
	}
	renderTable(w, []string{"schema", "table", "index", "real_size", "bloat_size", "bloat_pct", "reindex_cmd"}, indexRows)

	fmt.Fprintln(w)
	fmt.Fprintln(w, "Table Rollup")
	tableRows := make([][]string, 0, len(tables))
	for _, row := range tables {
		tableRows = append(tableRows, []string{
			row.Schema,
			row.Table,
			formatBytes(row.RealSize),
			formatBytes(row.BloatSize),
			formatPct(row.BloatPct),
			bloat.ReindexTableCommand(row.Schema, row.Table, concurrentSupported),
		})
	}
	renderTable(w, []string{"schema", "table", "real_size", "bloat_size", "bloat_pct", "reindex_cmd"}, tableRows)

	if !concurrentSupported {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Note: concurrent reindex is not supported on this Postgres version.")
	}

	return nil
}

func renderTable(w io.Writer, headers []string, rows [][]string) {
	widths := make([]int, len(headers))
	for i, header := range headers {
		widths[i] = len(header)
	}
	for _, row := range rows {
		for i, col := range row {
			if len(col) > widths[i] {
				widths[i] = len(col)
			}
		}
	}

	writeRow(w, headers, widths)
	separators := make([]string, len(headers))
	for i, width := range widths {
		separators[i] = strings.Repeat("-", width)
	}
	writeRow(w, separators, widths)
	for _, row := range rows {
		writeRow(w, row, widths)
	}
}

func writeRow(w io.Writer, row []string, widths []int) {
	for i, col := range row {
		if i > 0 {
			fmt.Fprint(w, "  ")
		}
		fmt.Fprint(w, padRight(col, widths[i]))
	}
	fmt.Fprintln(w)
}

func padRight(value string, width int) string {
	if len(value) >= width {
		return value
	}
	return value + strings.Repeat(" ", width-len(value))
}

func formatBytes(value int64) string {
	units := []string{"B", "KiB", "MiB", "GiB", "TiB"}
	size := float64(value)
	unit := 0
	for size >= 1024 && unit < len(units)-1 {
		size /= 1024
		unit++
	}
	if unit == 0 {
		return fmt.Sprintf("%d %s", value, units[unit])
	}
	return fmt.Sprintf("%.1f %s", size, units[unit])
}

func formatPct(value float64) string {
	return fmt.Sprintf("%.1f%%", value)
}
