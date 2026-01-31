package bloat

import "testing"

func TestReindexIndexCommand(t *testing.T) {
	cmd := ReindexIndexCommand("public", "idx", true)
	if cmd != "REINDEX INDEX CONCURRENTLY \"public\".\"idx\"" {
		t.Fatalf("unexpected command: %s", cmd)
	}
}

func TestReindexTableCommand(t *testing.T) {
	cmd := ReindexTableCommand("public", "tbl", false)
	if cmd != "REINDEX TABLE \"public\".\"tbl\"" {
		t.Fatalf("unexpected command: %s", cmd)
	}
}

func TestQuoteIdent(t *testing.T) {
	quoted := quoteIdent("foo\"bar")
	if quoted != "\"foo\"\"bar\"" {
		t.Fatalf("unexpected quote: %s", quoted)
	}
}
