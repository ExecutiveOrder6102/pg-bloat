package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ExecutiveOrder6102/postgres-bloat/internal/bloat"
)

func TestWriteConsole(t *testing.T) {
	indexes := []bloat.IndexBloat{{
		Schema:    "public",
		Table:     "orders",
		Index:     "orders_pkey",
		RealSize:  1024,
		BloatSize: 512,
		BloatPct:  50,
	}}
	tables := []bloat.TableBloat{{
		Schema:    "public",
		Table:     "orders",
		RealSize:  1024,
		BloatSize: 512,
		BloatPct:  50,
	}}

	var buf bytes.Buffer
	if err := WriteConsole(&buf, indexes, tables, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()

	for _, needle := range []string{
		"Index Bloat",
		"Table Rollup",
		"schema",
		"reindex_cmd",
		"REINDEX INDEX \"public\".\"orders_pkey\"",
		"REINDEX TABLE \"public\".\"orders\"",
		"Note: concurrent reindex is not supported",
	} {
		if !strings.Contains(output, needle) {
			t.Fatalf("expected output to contain %q", needle)
		}
	}
}
