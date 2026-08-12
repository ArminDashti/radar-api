CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS probes (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    flag_icon TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS endpoints (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    host TEXT NOT NULL,
    http_enabled BOOLEAN NOT NULL DEFAULT true,
    icmp_enabled BOOLEAN NOT NULL DEFAULT true,
    probe_id BIGINT REFERENCES probes(id) ON DELETE SET NULL,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS agents (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    probe_id BIGINT NOT NULL REFERENCES probes(id) ON DELETE RESTRICT,
    token_hash TEXT NOT NULL UNIQUE,
    last_seen_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS samples (
    id BIGSERIAL PRIMARY KEY,
    agent_id BIGINT NOT NULL REFERENCES agents(id) ON DELETE CASCADE,
    endpoint_id BIGINT NOT NULL REFERENCES endpoints(id) ON DELETE CASCADE,
    protocol TEXT NOT NULL CHECK (protocol IN ('http', 'icmp')),
    observed_at TIMESTAMPTZ NOT NULL,
    latency_ms DOUBLE PRECISION,
    ok BOOLEAN NOT NULL,
    CONSTRAINT samples_latency_check CHECK (
        (ok AND latency_ms IS NOT NULL AND latency_ms >= 0)
        OR (NOT ok AND latency_ms IS NULL)
    ),
    UNIQUE (agent_id, endpoint_id, protocol, observed_at)
);

CREATE INDEX IF NOT EXISTS samples_grid_idx
    ON samples (protocol, observed_at, endpoint_id, agent_id);

INSERT INTO probes (code, name, flag_icon) VALUES
    ('probe1', 'Helsinki', '/flags/probe1.svg'),
    ('probe2', 'Frankfurt', '/flags/probe2.svg'),
    ('probe3', 'Tehran', '/flags/probe3.svg')
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name, flag_icon = EXCLUDED.flag_icon;

INSERT INTO endpoints (name, host, http_enabled, icmp_enabled, probe_id, active)
SELECT seed.name, seed.host, true, true, p.id, true
FROM (VALUES
    ('google.com', 'google.com', 'probe1'),
    ('cloudflare.com 1.1.1.1', '1.1.1.1', 'probe2'),
    ('github.com', 'github.com', 'probe3')
) AS seed(name, host, probe_code)
JOIN probes p ON p.code = seed.probe_code
WHERE NOT EXISTS (SELECT 1 FROM endpoints e WHERE e.host = seed.host);
