# GO-Weather

Full-stack weather dashboard with a Go HTTP API, React + TypeScript frontend, Open-Meteo weather data, in-memory cache, and optional PostgreSQL persistence.

## Features

- Current conditions, hourly forecast, and 7-day forecast from Open-Meteo.
- City search using Open-Meteo Geocoding.
- Temperature, wind, humidity, and pressure charts.
- Celsius/Fahrenheit display, refresh, and responsive charcoal UI.
- Thread-safe in-memory weather cache with stale-data fallback.
- PostgreSQL persistence for saved cities and user settings.
- No Open-Meteo API key required for the configured endpoints.

## Run with Docker Compose

Requires Docker Engine with the Compose plugin.

```sh
cp .env.example .env
# Optionally change POSTGRES_PASSWORD in .env
docker compose up --build
```

Open <http://localhost:8080>. Compose starts PostgreSQL, waits for its health check, and the Go server applies SQL migrations on startup. PostgreSQL data is stored in the `postgres-data` volume.

## Local development

Requirements: Go 1.23+, Node.js 22+, and PostgreSQL if you want persistence.

Start the backend:

```sh
export DATABASE_URL='postgres://weather:weather-local@localhost:5432/weather?sslmode=disable'
go run ./cmd/server
```

Start the frontend in another terminal:

```sh
cd web
npm ci
npm run dev
```

Vite proxies `/api` requests to `http://localhost:8080`. Without `DATABASE_URL`, weather and city search work, while saved-city and settings endpoints are not registered.

## Configuration

| Variable | Default | Purpose |
| --- | --- | --- |
| `PORT` | `8080` | Backend listen port |
| `WEB_DIST` | `web/dist` | Built frontend directory |
| `ALLOWED_ORIGIN` | `http://localhost:5173` | Allowed browser origin |
| `DATABASE_URL` | unset | PostgreSQL connection string; enables persistence endpoints |
| `POSTGRES_PASSWORD` | `weather-local` in Compose | PostgreSQL password; set a private value in `.env` |

## API

- `GET /api/health`
- `GET /api/weather/current?latitude=50.08&longitude=14.43`
- `GET /api/weather/hourly?latitude=50.08&longitude=14.43`
- `GET /api/weather/daily?latitude=50.08&longitude=14.43`
- `GET /api/weather/overview?latitude=50.08&longitude=14.43`
- `GET /api/cities/search?name=Prague`
- With PostgreSQL: `GET|POST /api/cities/favorites`, `DELETE /api/cities/favorites?latitude=...&longitude=...`, `GET|PUT /api/settings`

Weather data is provided by [Open-Meteo](https://open-meteo.com/).

## Checks

```sh
go test ./...
go build ./...
cd web && npm ci && npm run build
```
