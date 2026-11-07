package main

import (
	"context"
	"github.com/kryvna/go-weather/internal/api"
	"github.com/kryvna/go-weather/internal/city"
	"github.com/kryvna/go-weather/internal/config"
	"github.com/kryvna/go-weather/internal/database"
	"github.com/kryvna/go-weather/internal/settings"
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
	apiHandler := api.New(weather.NewClient(), city.NewClient())
	if cfg.DatabaseURL != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		db, err := database.Connect(ctx, cfg.DatabaseURL)
		cancel()
		if err != nil {
			logger.Error("database connection failed", "error", err)
			os.Exit(1)
		}
		defer db.Close()
		apiHandler.WithPersistence(city.NewRepository(db.Pool), settings.NewRepository(db.Pool))
	}
	apiHandler.Register(mux)
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
