package api

import (
	"net/http/httptest"
	"testing"
)

func TestCoords(t *testing.T) {
	tests := []struct {
		lat, lon string
		ok       bool
	}{{"50.08", "14.43", true}, {"91", "0", false}, {"0", "181", false}, {"x", "0", false}}
	for _, tt := range tests {
		r := httptest.NewRequest("GET", "/?latitude="+tt.lat+"&longitude="+tt.lon, nil)
		_, _, ok := coords(r)
		if ok != tt.ok {
			t.Errorf("coords(%q,%q)=%v want %v", tt.lat, tt.lon, ok, tt.ok)
		}
	}
}
