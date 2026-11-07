# GO-Weather

A full-stack weather dashboard project. The first milestone provides a Go health endpoint and a React + TypeScript + Vite frontend styled with Tailwind CSS.

## Requirements

- Go 1.23+
- Node.js 22+
- Docker with Docker Compose (optional)

## Development

Run the backend in one terminal:

```sh
go run ./cmd/server
```

Run the frontend in another:

```sh
cd web
npm install
npm run dev
```

Vite proxies `/api` requests to `http://localhost:8080`.

## Docker

```sh
docker compose up --build
```

The app listens on port 8080; health check: `GET /api/health`.
