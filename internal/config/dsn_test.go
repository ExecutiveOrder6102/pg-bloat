package config

import (
	"net/url"
	"testing"
)

func TestBuildDSNRequiresFields(t *testing.T) {
	_, err := BuildDSN(Config{User: "alice", Database: "bloat"})
	if err == nil {
		t.Fatal("expected error for missing host")
	}
}

func TestBuildDSNWithPassword(t *testing.T) {
	dsn, err := BuildDSN(Config{Host: "localhost", Port: 0, User: "alice", Password: "secret", Database: "bloat", SSLMode: "require"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("failed to parse dsn: %v", err)
	}
	if parsed.Scheme != "postgres" {
		t.Fatalf("expected postgres scheme, got %q", parsed.Scheme)
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

func TestBuildDSNWithoutPassword(t *testing.T) {
	dsn, err := BuildDSN(Config{Host: "db", Port: 5433, User: "alice", Database: "bloat"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	parsed, err := url.Parse(dsn)
	if err != nil {
		t.Fatalf("failed to parse dsn: %v", err)
	}
	if parsed.Host != "db:5433" {
		t.Fatalf("expected host db:5433, got %q", parsed.Host)
	}
	if parsed.User.Username() != "alice" {
		t.Fatalf("expected user alice, got %q", parsed.User.Username())
	}
	if _, ok := parsed.User.Password(); ok {
		t.Fatal("expected no password in dsn")
	}
}
