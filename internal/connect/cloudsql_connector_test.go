package connect

import (
	"net/url"
	"testing"
)

func TestBuildLocalDSNRequiresFields(t *testing.T) {
	_, err := buildLocalDSN("", "", "db", "")
	if err == nil {
		t.Fatal("expected error for missing user")
	}
}

func TestBuildLocalDSNWithPassword(t *testing.T) {
	dsn, err := buildLocalDSN("alice", "secret", "bloat", "require")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("failed to parse dsn: %v", err)
	}
	if parsed.Scheme != "postgres" {
		t.Fatalf("expected scheme postgres, got %q", parsed.Scheme)
	}
	if parsed.Host != "localhost:5432" {
		t.Fatalf("expected host localhost:5432, got %q", parsed.Host)
	}
	if parsed.Path != "/bloat" {
		t.Fatalf("expected path /bloat, got %q", parsed.Path)
	}
	if parsed.User.Username() != "alice" {
		t.Fatalf("expected user alice, got %q", parsed.User.Username())
	}
	pass, ok := parsed.User.Password()
	if !ok || pass != "secret" {
		t.Fatalf("expected password secret, got %q", pass)
	}
	if parsed.Query().Get("sslmode") != "require" {
		t.Fatalf("expected sslmode require, got %q", parsed.Query().Get("sslmode"))
	}
}

func TestBuildLocalDSNWithoutPassword(t *testing.T) {
	dsn, err := buildLocalDSN("alice", "", "bloat", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("failed to parse dsn: %v", err)
	}
	if _, ok := parsed.User.Password(); ok {
		t.Fatal("expected no password")
	}
	if parsed.Query().Get("sslmode") != "" {
		t.Fatalf("expected empty sslmode, got %q", parsed.Query().Get("sslmode"))
	}
}
