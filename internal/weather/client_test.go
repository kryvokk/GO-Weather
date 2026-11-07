package weather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchMapsProviderPayload(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"timezone":"Europe/Prague","current":{"temperature_2m":12,"relative_humidity_2m":70,"apparent_temperature":11,"is_day":1,"weather_code":2},"hourly":{"time":["2026-11-07T12:00"],"temperature_2m":[12],"relative_humidity_2m":[70]},"daily":{"time":["2026-11-07"],"weather_code":[2],"temperature_2m_max":[15],"temperature_2m_min":[8]}}`))
	}))
	defer srv.Close()
	c := NewClient()
	c.base = srv.URL
	o, e := c.Fetch(context.Background(), 50, 14)
	if e != nil {
		t.Fatal(e)
	}
	if o.Current.Temperature != 12 || len(o.Hourly) != 1 || o.Hourly[0].Humidity != 70 || len(o.Daily) != 1 || o.Daily[0].TemperatureMax != 15 {
		t.Fatalf("unexpected mapping: %+v", o)
	}
}
