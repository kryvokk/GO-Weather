package main

import (
	"github.com/kryvna/go-weather/internal/api"
	"github.com/kryvna/go-weather/internal/city"
	"github.com/kryvna/go-weather/internal/config"
	"github.com/kryvna/go-weather/internal/weather"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
	api.New(weather.NewClient(), city.NewClient()).Register(mux)
	mux.Handle("GET /", http.FileServer(http.Dir(cfg.WebDist)))
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", cfg.AllowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(204)
			return
		}
		mux.ServeHTTP(w, r)
	})
	addr := ":" + cfg.Port
	s := &http.Server{Addr: addr, Handler: handler, ReadHeaderTimeout: 5 * time.Second}
	logger.Info("server starting", "addr", addr)
	if e := s.ListenAndServe(); e != nil && e != http.ErrServerClosed {
		logger.Error("server failed", "error", e)
		os.Exit(1)
	}
}
