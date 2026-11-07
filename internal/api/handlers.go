package api

import (
	"encoding/json"
	"github.com/kryvna/go-weather/internal/cache"
	"github.com/kryvna/go-weather/internal/city"
	"github.com/kryvna/go-weather/internal/weather"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	weather       *weather.Client
	cities        *city.Client
	forecastCache cache.Cache[weather.Overview]
}

func New(w *weather.Client, c *city.Client) *Handler {
	return &Handler{weather: w, cities: c, forecastCache: cache.NewMemory[weather.Overview]()}
}
func (h *Handler) Register(m *http.ServeMux) {
	m.HandleFunc("GET /api/weather/overview", h.overview)
	m.HandleFunc("GET /api/weather/current", h.current)
	m.HandleFunc("GET /api/weather/hourly", h.hourly)
	m.HandleFunc("GET /api/weather/daily", h.daily)
	m.HandleFunc("GET /api/cities/search", h.search)
}
func coords(r *http.Request) (float64, float64, bool) {
	a, e := strconv.ParseFloat(r.URL.Query().Get("latitude"), 64)
	if e != nil {
		return 0, 0, false
	}
	b, e := strconv.ParseFloat(r.URL.Query().Get("longitude"), 64)
	return a, b, e == nil && !math.IsNaN(a) && !math.IsNaN(b) && a >= -90 && a <= 90 && b >= -180 && b <= 180
}
func (h *Handler) overview(w http.ResponseWriter, r *http.Request) {
	a, b, ok := coords(r)
	if !ok {
		fail(w, 400, "valid latitude and longitude are required")
		return
	}
	v, e := h.getForecast(r, a, b)
	if e != nil {
		fail(w, 502, "weather provider is unavailable")
		return
	}
	jsonOut(w, 200, v)
}
func (h *Handler) current(w http.ResponseWriter, r *http.Request) { h.part(w, r, 0) }
func (h *Handler) hourly(w http.ResponseWriter, r *http.Request)  { h.part(w, r, 1) }
func (h *Handler) daily(w http.ResponseWriter, r *http.Request)   { h.part(w, r, 2) }
func (h *Handler) part(w http.ResponseWriter, r *http.Request, n int) {
	a, b, ok := coords(r)
	if !ok {
		fail(w, 400, "valid latitude and longitude are required")
		return
	}
	v, e := h.getForecast(r, a, b)
	if e != nil {
		fail(w, 502, "weather provider is unavailable")
		return
	}
	switch n {
	case 0:
		jsonOut(w, 200, v.Current)
	case 1:
		jsonOut(w, 200, v.Hourly)
	default:
		jsonOut(w, 200, v.Daily)
	}
}
func (h *Handler) getForecast(r *http.Request, lat, lon float64) (weather.Overview, error) {
	key := strconv.FormatFloat(lat, 'f', 4, 64) + ":" + strconv.FormatFloat(lon, 'f', 4, 64)
	if entry, ok := h.forecastCache.Get(key); ok && !entry.Stale {
		return entry.Value, nil
	}
	v, err := h.weather.Fetch(r.Context(), lat, lon)
	if err == nil {
		h.forecastCache.Set(key, v, 5*time.Minute)
		return v, nil
	}
	if entry, ok := h.forecastCache.Get(key); ok {
		entry.Value.Updated = entry.FetchedAt.Format(time.RFC3339) + " (stale cache)"
		return entry.Value, nil
	}
	return weather.Overview{}, err
}
func (h *Handler) search(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(r.URL.Query().Get("name"))
	if len([]rune(name)) < 2 || len([]rune(name)) > 100 {
		fail(w, 400, "city name must contain 2 to 100 characters")
		return
	}
	v, e := h.cities.Search(r.Context(), name)
	if e != nil {
		fail(w, 502, "city search is unavailable")
		return
	}
	jsonOut(w, 200, v)
}
func jsonOut(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func fail(w http.ResponseWriter, status int, msg string) {
	jsonOut(w, status, weather.APIError{Error: msg})
}

var _ = time.Second
