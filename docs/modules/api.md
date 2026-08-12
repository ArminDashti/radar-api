# HTTP API Module

`internal/handlers` owns HTTP validation, authentication middleware, database calls, and JSON responses.

- `Server` supplies the PostgreSQL pool and JWT secret.
- `WebAuth` validates signed JWT bearer tokens.
- `AgentAuth` hashes and validates agent bearer tokens.
- Grid handlers return oldest-to-newest cells aligned with bucket timestamps.

The module depends on `auth`, `models`, `rollup`, Gin, and pgx.
