package cache

import (
	"testing"
	"time"
)

func TestMemoryMarksExpiredEntryStale(t *testing.T) {
	m := NewMemory[string]()
	m.Set("city", "forecast", time.Millisecond)
	a, ok := m.Get("city")
	if !ok || a.Stale {
		t.Fatalf("expected fresh entry, got %+v", a)
	}
	time.Sleep(2 * time.Millisecond)
	b, ok := m.Get("city")
	if !ok || !b.Stale || b.Value != "forecast" {
		t.Fatalf("expected retained stale value, got %+v", b)
	}
}
