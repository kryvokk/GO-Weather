package config

import "os"

type Config struct{ Port, WebDist, AllowedOrigin, DatabaseURL string }

func Load() Config {
	return Config{Port: get("PORT", "8080"), WebDist: get("WEB_DIST", "web/dist"), AllowedOrigin: get("ALLOWED_ORIGIN", "http://localhost:5173"), DatabaseURL: os.Getenv("DATABASE_URL")}
}
func get(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
