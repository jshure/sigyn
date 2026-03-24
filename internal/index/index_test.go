package index

import (
	"testing"
	"time"
)

func TestMatchLabels_EmptyQuery(t *testing.T) {
	sets := []map[string]string{{"app": "nginx"}}
	if !MatchLabels(sets, nil) {
		t.Error("empty query should match everything")
	}
}

func TestMatchLabels_EmptySets(t *testing.T) {
	query := map[string]string{"app": "nginx"}
	if MatchLabels(nil, query) {
		t.Error("empty label sets should not match a non-empty query")
	}
}

func TestMatchLabels_ExactMatch(t *testing.T) {
	sets := []map[string]string{{"app": "nginx", "env": "prod"}}
	query := map[string]string{"app": "nginx", "env": "prod"}
	if !MatchLabels(sets, query) {
		t.Error("expected match for exact label set")
	}
}

func TestMatchLabels_SubsetMatch(t *testing.T) {
	sets := []map[string]string{{"app": "nginx", "env": "prod", "region": "us-west-2"}}
	query := map[string]string{"app": "nginx"}
	if !MatchLabels(sets, query) {
		t.Error("query is a subset of label set, should match")
	}
}

func TestMatchLabels_Mismatch(t *testing.T) {
	sets := []map[string]string{{"app": "nginx", "env": "staging"}}
	query := map[string]string{"env": "prod"}
	if MatchLabels(sets, query) {
		t.Error("expected no match for wrong value")
	}
}

func TestMatchLabels_MultipleSetOneMatches(t *testing.T) {
	sets := []map[string]string{
		{"app": "nginx", "env": "staging"},
		{"app": "nginx", "env": "prod"},
	}
	query := map[string]string{"env": "prod"}
	if !MatchLabels(sets, query) {
		t.Error("expected match when one of multiple sets matches")
	}
}

func TestMatchLabels_MissingKey(t *testing.T) {
	sets := []map[string]string{{"app": "nginx"}}
	query := map[string]string{"env": "prod"}
	if MatchLabels(sets, query) {
		t.Error("expected no match when query key is missing from set")
	}
}

func TestTimeOverlaps_FullyInside(t *testing.T) {
	idx := timeRange(10, 20)
	q := timeRange(5, 25)
	if !timeOverlaps(idx.min, idx.max, q.min, q.max) {
		t.Error("index fully inside query should overlap")
	}
}

func TestTimeOverlaps_PartialLeft(t *testing.T) {
	idx := timeRange(5, 15)
	q := timeRange(10, 20)
	if !timeOverlaps(idx.min, idx.max, q.min, q.max) {
		t.Error("partial left overlap should match")
	}
}

func TestTimeOverlaps_PartialRight(t *testing.T) {
	idx := timeRange(15, 25)
	q := timeRange(10, 20)
	if !timeOverlaps(idx.min, idx.max, q.min, q.max) {
		t.Error("partial right overlap should match")
	}
}

func TestTimeOverlaps_NoOverlap(t *testing.T) {
	idx := timeRange(1, 5)
	q := timeRange(10, 20)
	if timeOverlaps(idx.min, idx.max, q.min, q.max) {
		t.Error("non-overlapping ranges should not match")
	}
}

func TestTimeOverlaps_EdgeExact(t *testing.T) {
	idx := timeRange(10, 10)
	q := timeRange(10, 20)
	if !timeOverlaps(idx.min, idx.max, q.min, q.max) {
		t.Error("touching at edge should overlap")
	}
}

func TestIndexKey_Format(t *testing.T) {
	s := &Store{prefix: "index/"}
	key := s.indexKey("myapp/2026/03/23/1679558400-abc123.gz")
	want := "index/2026/03/23/myapp-1679558400-abc123.json"
	if key != want {
		t.Errorf("indexKey = %q, want %q", key, want)
	}
}

func TestDayPrefixes(t *testing.T) {
	s := &Store{prefix: "index/"}
	start := time.Date(2026, 3, 22, 12, 0, 0, 0, time.UTC)
	end := time.Date(2026, 3, 24, 6, 0, 0, 0, time.UTC)

	pfx := s.dayPrefixes(start, end)
	if len(pfx) != 3 {
		t.Fatalf("expected 3 prefixes, got %d: %v", len(pfx), pfx)
	}
	expected := []string{"index/2026/03/22/", "index/2026/03/23/", "index/2026/03/24/"}
	for i, want := range expected {
		if pfx[i] != want {
			t.Errorf("prefix[%d] = %q, want %q", i, pfx[i], want)
		}
	}
}

type tr struct{ min, max time.Time }

func timeRange(minHour, maxHour int) tr {
	base := time.Date(2026, 3, 23, 0, 0, 0, 0, time.UTC)
	return tr{base.Add(time.Duration(minHour) * time.Hour), base.Add(time.Duration(maxHour) * time.Hour)}
}
