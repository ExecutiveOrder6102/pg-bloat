package bloat

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestServerVersion(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	mock.ExpectQuery("SHOW server_version_num").WillReturnRows(
		sqlmock.NewRows([]string{"server_version_num"}).AddRow("120004"),
	)

	version, err := ServerVersion(context.Background(), db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != 120004 {
		t.Fatalf("expected 120004, got %d", version)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestFetchIndexBloat(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	query := regexp.QuoteMeta(indexBloatQuery)
	rows := sqlmock.NewRows([]string{"schemaname", "tblname", "idxname", "real_size", "bloat_size", "bloat_pct"}).
		AddRow("public", "t1", "i1", int64(1000), int64(200), 20.5)

	mock.ExpectQuery(query).WillReturnRows(rows)

	result, err := FetchIndexBloat(context.Background(), db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 row, got %d", len(result))
	}
	if result[0].Index != "i1" || result[0].BloatSize != 200 {
		t.Fatalf("unexpected row: %+v", result[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestFetchStaleStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	query := regexp.QuoteMeta(staleStatsQuery)
	rows := sqlmock.NewRows([]string{"schemaname", "relname", "last_analyze"}).
		AddRow("public", "t1", time.Now())

	mock.ExpectQuery(query).WithArgs("7 days").WillReturnRows(rows)

	result, err := FetchStaleStats(context.Background(), db, 7*24*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 row, got %d", len(result))
	}
	if !result[0].LastAnalyze.Valid {
		t.Fatalf("expected last analyze to be valid: %+v", result[0])
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
