# Radar API

Gin and PostgreSQL backend for the Radar latency hub. It accepts HTTP/ICMP latency samples from authenticated agents and serves hierarchical minute, hour, day, and month grids to the WebUI.

## Run locally

Requirements: Go 1.22+, Docker, and Docker Compose.

```powershell
Copy-Item .env.example .env
docker compose up -d
go run ./cmd/server
```

Postgres is published on host port **5436** (`DATABASE_URL` in `.env.example`). The server runs at `http://localhost:8088`. Startup applies `migrations/001_init.sql`, creates the demo credentials, and inserts two hours of sample history. Seeding is idempotent.

WebUI login:

- Username: `armin`
- Password: `dopadopa123`

Demo agent bearer token:

```text
radar-demo-agent-token-probe1
```

Only the SHA-256 token hash is stored in PostgreSQL. Replace all demo credentials and `JWT_SECRET` outside local development.

## API

Public:

- `POST /api/auth/login`

JWT bearer authentication:

- `GET /api/probes`
- `GET /api/endpoints`
- `POST /api/endpoints`
- `GET /api/grid/endpoints?interval=minutes&protocol=http&probe=all`
- `GET /api/grid/probes?interval=hours&protocol=icmp`

Agent bearer authentication:

- `GET /api/agent/targets`
- `POST /api/agent/samples`

Intervals support `minutes`, `hours`, `days`, and `months`. Protocols support `http` and `icmp`. Sample timestamps are normalized to UTC minute buckets, and duplicate submissions update the existing minute sample.

Example login:

```powershell
Invoke-RestMethod -Method Post http://localhost:8088/api/auth/login `
  -ContentType application/json `
  -Body '{"username":"armin","password":"dopadopa123"}'
```

Example agent sample:

```json
{
  "samples": [
    {
      "endpoint_id": 1,
      "protocol": "http",
      "observed_at": "2026-08-13T01:00:00Z",
      "latency_ms": 42.5,
      "ok": true
    }
  ]
}
```

For timeouts, send `"ok": false` and omit or set `"latency_ms": null`.
