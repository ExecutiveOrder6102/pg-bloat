package config

import (
	"strings"
	"testing"
)

func TestParseRequiresConnectionInfo(t *testing.T) {
	_, err := Parse([]string{})
	if err == nil {
		t.Fatal("expected error for missing connection args")
	}
	if !strings.Contains(err.Error(), "dsn or host, user, dbname required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseRejectsDSNWithCloudSQL(t *testing.T) {
	_, err := Parse([]string{"--dsn", "postgres://example", "--cloudsql-instance", "proj:region:inst"})
	if err == nil {
		t.Fatal("expected error for dsn with cloudsql")
	}
	if !strings.Contains(err.Error(), "use either --dsn or --cloudsql-instance") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseRejectsInvalidOutput(t *testing.T) {
	_, err := Parse([]string{"--host", "localhost", "--user", "alice", "--dbname", "bloat", "--output", "json"})
	if err == nil {
		t.Fatal("expected error for invalid output")
	}
	if !strings.Contains(err.Error(), "unsupported output") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseRejectsInvalidLimit(t *testing.T) {
	_, err := Parse([]string{"--host", "localhost", "--user", "alice", "--dbname", "bloat", "--limit", "0"})
	if err == nil {
		t.Fatal("expected error for limit < 1")
	}
	if !strings.Contains(err.Error(), "limit must be >= 1") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseRejectsInvalidCloudSQLIPType(t *testing.T) {
	_, err := Parse([]string{"--host", "localhost", "--user", "alice", "--dbname", "bloat", "--cloudsql-ip-type", "loopback"})
	if err == nil {
		t.Fatal("expected error for invalid cloudsql ip type")
	}
	if !strings.Contains(err.Error(), "cloudsql-ip-type must be public or private") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseAcceptsValidArgs(t *testing.T) {
	cfg, err := Parse([]string{"--host", "localhost", "--user", "alice", "--dbname", "bloat"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Port != 5432 {
		t.Fatalf("expected default port 5432, got %d", cfg.Port)
	}
	if cfg.Output != "console" {
		t.Fatalf("expected default output console, got %q", cfg.Output)
	}
	if cfg.SSLMode != "disable" {
		t.Fatalf("expected default sslmode disable, got %q", cfg.SSLMode)
	}
}
