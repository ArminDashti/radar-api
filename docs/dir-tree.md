# Directory Tree

```text
cmd/server/main.go                 # Configures and starts Gin server
internal/auth/auth.go              # JWT and agent-token helpers
internal/config/config.go          # Loads environment configuration
internal/config/config_test.go     # Tests configuration behavior
internal/db/db.go                  # Connects, migrates, and seeds database
internal/handlers/agent.go         # Handles agent target and sample routes
internal/handlers/auth.go          # Handles WebUI login
internal/handlers/endpoints.go     # Handles endpoint routes
internal/handlers/grid.go          # Builds hierarchical latency grids
internal/handlers/middleware.go    # Auth middleware and server dependencies
internal/handlers/probes.go        # Handles probe listing
internal/models/models.go          # Shared request and response models
internal/rollup/rollup.go          # Defines intervals and time buckets
internal/rollup/rollup_test.go     # Tests interval and bucket behavior
migrations/001_init.sql            # Creates and seeds PostgreSQL schema
docs/                              # Project reference documentation
.env.example                       # Local environment template
.gitignore                         # Excludes local and generated files
docker-compose.yml                 # Runs PostgreSQL 16 locally
go.mod                             # Go module and direct dependencies
go.sum                             # Dependency checksums
README.md                          # Setup, credentials, and API guide
```
