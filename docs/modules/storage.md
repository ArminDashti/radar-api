# Storage Module

`internal/db` opens and verifies the pgx pool, executes the SQL migration, and idempotently creates demo credentials and historical samples.

`migrations/001_init.sql` defines users, probes, endpoints, agents, and samples. Samples are unique per agent, endpoint, protocol, and minute timestamp. Successful samples require non-negative latency; failed samples require null latency.
