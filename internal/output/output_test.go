package output

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ExecutiveOrder6102/postgres-bloat/internal/bloat"
	"github.com/ExecutiveOrder6102/postgres-bloat/internal/config"
)

func TestWriteToFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.csv")

	cfg := config.Config{
		Output:     "csv",
		OutputFile: path,
	}
	indexes := []bloat.IndexBloat{{Schema: "public", Table: "t1", Index: "i1", RealSize: 1, BloatSize: 1, BloatPct: 100}}
	tables := []bloat.TableBloat{{Schema: "public", Table: "t1", RealSize: 1, BloatSize: 1, BloatPct: 100}}

	if err := Write(cfg, indexes, tables, true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("expected output file: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("expected output file to be non-empty")
	}
}
