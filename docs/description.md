# Radar API

Radar API is the Gin/PostgreSQL hub for collecting agent latency samples and serving WebUI grids. It uses Go 1.22, pgx, JWT WebUI authentication, SHA-256 agent tokens, PostgreSQL 16, and Docker Compose.

Run PostgreSQL with `docker compose up -d`, copy `.env.example` to `.env`, then start `go run ./cmd/server`.
