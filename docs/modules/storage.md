# Storage Module

`internal/db` opens and verifies the pgx pool, executes the SQL migration, and idempotently creates demo credentials and historical samples.

`migrations/001_init.sql` defines users, probes, endpoints, agents, and samples. Endpoints store an optional `logo_icon` path for the host grid. The seed creates one probe (`probe1` / Irancell) and 17 monitored sites (named hosts HTTP-only; `1.1.1.1` and `8.8.8.8` ICMP-only). Samples are unique per agent, endpoint, protocol, and minute timestamp. Successful samples require non-negative latency; failed samples require null latency.
