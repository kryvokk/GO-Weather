FROM golang:1.23-alpine AS backend-build
WORKDIR /app
COPY go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal
COPY migrations ./migrations
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

FROM node:22-alpine AS frontend-build
WORKDIR /app
COPY web/package*.json ./
RUN npm ci
COPY web ./
RUN npm run build

FROM alpine:3.21
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=backend-build /server /app/server
COPY --from=frontend-build /app/dist /app/web/dist
EXPOSE 8080
CMD ["/app/server"]
