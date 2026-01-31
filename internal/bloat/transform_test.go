package bloat

import (
	"math"
	"testing"
)

func TestFilterIndexes(t *testing.T) {
	rows := []IndexBloat{
		{Schema: "public", Table: "t1", Index: "i1", RealSize: 1000, BloatSize: 200, BloatPct: 20},
		{Schema: "public", Table: "t2", Index: "i2", RealSize: 1000, BloatSize: 50, BloatPct: 5},
		{Schema: "pg_catalog", Table: "t3", Index: "i3", RealSize: 1000, BloatSize: 500, BloatPct: 50},
	}
	filtered := FilterIndexes(rows, 10, 100, false)
	if len(filtered) != 1 {
		t.Fatalf("expected 1 result, got %d", len(filtered))
	}
	if filtered[0].Index != "i1" {
		t.Fatalf("expected i1, got %s", filtered[0].Index)
	}
}

func TestRollupTables(t *testing.T) {
	rows := []IndexBloat{
		{Schema: "public", Table: "t1", Index: "i1", RealSize: 100, BloatSize: 10},
		{Schema: "public", Table: "t1", Index: "i2", RealSize: 50, BloatSize: 5},
	}
	rolled := RollupTables(rows)
	if len(rolled) != 1 {
		t.Fatalf("expected 1 table, got %d", len(rolled))
	}
	if rolled[0].RealSize != 150 || rolled[0].BloatSize != 15 {
		t.Fatalf("unexpected rollup sizes: %+v", rolled[0])
	}
	if math.Abs(rolled[0].BloatPct-10.0) > 0.001 {
		t.Fatalf("unexpected rollup pct: %.3f", rolled[0].BloatPct)
	}
}

func TestFilterTables(t *testing.T) {
	rows := []TableBloat{
		{Schema: "public", Table: "t1", RealSize: 1000, BloatSize: 200, BloatPct: 20},
		{Schema: "public", Table: "t2", RealSize: 1000, BloatSize: 50, BloatPct: 5},
		{Schema: "pg_catalog", Table: "t3", RealSize: 1000, BloatSize: 500, BloatPct: 50},
	}
	filtered := FilterTables(rows, 10, 100, false)
	if len(filtered) != 1 {
		t.Fatalf("expected 1 result, got %d", len(filtered))
	}
	if filtered[0].Table != "t1" {
		t.Fatalf("expected t1, got %s", filtered[0].Table)
	}
}

func TestSortIndexes(t *testing.T) {
	rows := []IndexBloat{
		{Schema: "public", Table: "t1", Index: "i1", BloatSize: 100, BloatPct: 10},
		{Schema: "public", Table: "t1", Index: "i2", BloatSize: 100, BloatPct: 20},
		{Schema: "public", Table: "t1", Index: "i3", BloatSize: 200, BloatPct: 5},
	}
	SortIndexes(rows)
	if rows[0].Index != "i3" {
		t.Fatalf("expected i3 first, got %s", rows[0].Index)
	}
	if rows[1].Index != "i2" {
		t.Fatalf("expected i2 second, got %s", rows[1].Index)
	}
}

func TestSortTables(t *testing.T) {
	rows := []TableBloat{
		{Schema: "public", Table: "t1", BloatSize: 100, BloatPct: 10},
		{Schema: "public", Table: "t2", BloatSize: 100, BloatPct: 20},
		{Schema: "public", Table: "t3", BloatSize: 200, BloatPct: 5},
	}
	SortTables(rows)
	if rows[0].Table != "t3" {
		t.Fatalf("expected t3 first, got %s", rows[0].Table)
	}
	if rows[1].Table != "t2" {
		t.Fatalf("expected t2 second, got %s", rows[1].Table)
	}
}

func TestApplyLimitIndexes(t *testing.T) {
	rows := []IndexBloat{{Index: "i1"}, {Index: "i2"}}
	limited := ApplyLimitIndexes(rows, 1)
	if len(limited) != 1 || limited[0].Index != "i1" {
		t.Fatalf("unexpected limit results: %+v", limited)
	}
}

func TestApplyLimitTables(t *testing.T) {
	rows := []TableBloat{{Table: "t1"}, {Table: "t2"}}
	limited := ApplyLimitTables(rows, 1)
	if len(limited) != 1 || limited[0].Table != "t1" {
		t.Fatalf("unexpected limit results: %+v", limited)
	}
}
