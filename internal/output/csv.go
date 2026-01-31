package output

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/ExecutiveOrder6102/postgres-bloat/internal/bloat"
)

func WriteCSV(w io.Writer, indexes []bloat.IndexBloat, tables []bloat.TableBloat, concurrentSupported bool) error {
	writer := csv.NewWriter(w)
	if err := writer.Write([]string{
		"record_type",
		"schema",
		"table",
		"index",
		"real_size_bytes",
		"bloat_size_bytes",
		"bloat_pct",
		"reindex_command",
	}); err != nil {
		return err
	}

	for _, row := range indexes {
		if err := writer.Write([]string{
			"index",
			row.Schema,
			row.Table,
			row.Index,
			fmt.Sprintf("%d", row.RealSize),
			fmt.Sprintf("%d", row.BloatSize),
			fmt.Sprintf("%.2f", row.BloatPct),
			bloat.ReindexIndexCommand(row.Schema, row.Index, concurrentSupported),
		}); err != nil {
			return err
		}
	}

	for _, row := range tables {
		if err := writer.Write([]string{
			"table",
			row.Schema,
			row.Table,
			"",
			fmt.Sprintf("%d", row.RealSize),
			fmt.Sprintf("%d", row.BloatSize),
			fmt.Sprintf("%.2f", row.BloatPct),
			bloat.ReindexTableCommand(row.Schema, row.Table, concurrentSupported),
		}); err != nil {
			return err
		}
	}

	writer.Flush()
	return writer.Error()
}
