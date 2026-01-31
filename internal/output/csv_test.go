package output

import (
	"bytes"
	"encoding/csv"
	"testing"

	"github.com/ExecutiveOrder6102/postgres-bloat/internal/bloat"
)

func TestWriteCSV(t *testing.T) {
	indexes := []bloat.IndexBloat{{
		Schema:    "public",
		Table:     "orders",
		Index:     "orders_pkey",
		RealSize:  1024,
		BloatSize: 512,
		BloatPct:  50.25,
	}}
	tables := []bloat.TableBloat{{
		Schema:    "public",
		Table:     "orders",
		RealSize:  1024,
		BloatSize: 512,
		BloatPct:  50.25,
	}}

	var buf bytes.Buffer
	if err := WriteCSV(&buf, indexes, tables, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	reader := csv.NewReader(bytes.NewReader(buf.Bytes()))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("failed to read csv: %v", err)
	}
	if len(records) != 3 {
		t.Fatalf("expected 3 records, got %d", len(records))
	}
	if records[0][0] != "record_type" {
		t.Fatalf("unexpected header: %v", records[0])
	}
	if records[1][0] != "index" {
		t.Fatalf("expected index record, got %v", records[1])
	}
	if records[2][0] != "table" {
		t.Fatalf("expected table record, got %v", records[2])
	}
	if records[1][7] != "REINDEX INDEX CONCURRENTLY \"public\".\"orders_pkey\"" {
		t.Fatalf("unexpected index reindex command: %v", records[1][7])
	}
	if records[2][7] != "REINDEX TABLE CONCURRENTLY \"public\".\"orders\"" {
		t.Fatalf("unexpected table reindex command: %v", records[2][7])
	}
}
