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

ALTER TABLE endpoints ADD COLUMN IF NOT EXISTS logo_icon TEXT NOT NULL DEFAULT '';

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
    ('probe1', 'Irancell', '/flags/probe1.svg')
ON CONFLICT (code) DO UPDATE
SET name = EXCLUDED.name, flag_icon = EXCLUDED.flag_icon;

UPDATE endpoints SET probe_id = (SELECT id FROM probes WHERE code = 'probe1');
UPDATE agents SET probe_id = (SELECT id FROM probes WHERE code = 'probe1');
DELETE FROM probes WHERE code IN ('probe2', 'probe3');

INSERT INTO endpoints (name, host, http_enabled, icmp_enabled, probe_id, active, logo_icon)
SELECT seed.name, seed.host, seed.http_enabled, seed.icmp_enabled, p.id, true, seed.logo_icon
FROM (VALUES
    ('Google', 'google.com', true, false, '/logos/google.svg'),
    ('GitHub', 'github.com', true, false, '/logos/github.svg'),
    ('1.1.1.1', '1.1.1.1', false, true, '/logos/cloudflare.svg'),
    ('8.8.8.8', '8.8.8.8', false, true, '/logos/google-dns.svg'),
    ('YouTube', 'youtube.com', true, false, '/logos/youtube.svg'),
    ('ChatGPT', 'chatgpt.com', true, false, '/logos/chatgpt.svg'),
    ('Claude', 'claude.ai', true, false, '/logos/claude.svg'),
    ('DeepSeek', 'deepseek.com', true, false, '/logos/deepseek.svg'),
    ('Microsoft', 'microsoft.com', true, false, '/logos/microsoft.svg'),
    ('Apple', 'apple.com', true, false, '/logos/apple.svg'),
    ('Play Store', 'play.google.com', true, false, '/logos/play-store.svg'),
    ('App Store', 'apps.apple.com', true, false, '/logos/app-store.svg'),
    ('Docker', 'docker.com', true, false, '/logos/docker.svg'),
    ('Spotify', 'spotify.com', true, false, '/logos/spotify.svg'),
    ('Grok', 'grok.com', true, false, '/logos/grok.svg'),
    ('Arena.ai', 'arena.ai', true, false, '/logos/arena.svg'),
    ('Ondpline.com', 'ondpline.com', true, false, '/logos/ondpline.svg'),
    ('GitLab', 'gitlab.com', true, false, '/logos/gitlab.svg'),
    ('digikala.ir', 'digikala.ir', true, false, '/logos/digikala.svg'),
    ('snapp.ir', 'snapp.ir', true, false, '/logos/snapp.svg'),
    ('tapsi.ir', 'tapsi.ir', true, false, '/logos/tapsi.svg'),
    ('divar.ir', 'divar.ir', true, false, '/logos/divar.svg'),
    ('bale.ai', 'bale.ai', true, false, '/logos/bale.svg'),
    ('Devin AI', 'devin.ai', true, false, '/logos/devin.svg'),
    ('Neshan.org', 'neshan.org', true, false, '/logos/neshan.svg'),
    ('Wikipedia', 'wikipedia.org', true, false, '/logos/wikipedia.svg'),
    ('WhatsApp', 'whatsapp.com', true, false, '/logos/whatsapp.svg'),
    ('zoomit.ir', 'zoomit.ir', true, false, '/logos/zoomit.svg'),
    ('digiato.ir', 'digiato.ir', true, false, '/logos/digiato.svg'),
    ('Let''s Encrypt', 'letsencrypt.org', true, false, '/logos/letsencrypt.svg'),
    ('Steam', 'steampowered.com', true, false, '/logos/steam.svg'),
    ('torob.ir', 'torob.ir', true, false, '/logos/torob.svg'),
    ('PlayStation', 'playstation.com', true, false, '/logos/playstation.svg'),
    ('Xbox', 'xbox.com', true, false, '/logos/xbox.svg'),
    ('Gmail', 'gmail.com', true, false, '/logos/gmail.svg'),
    ('Figma', 'figma.com', true, false, '/logos/figma.svg'),
    ('Windows Update', 'windowsupdate.microsoft.com', true, false, '/logos/windows-update.svg'),
    ('Bing', 'bing.com', true, false, '/logos/bing.svg'),
    ('aparat.ir', 'aparat.ir', true, false, '/logos/aparat.svg'),
    ('PyPI', 'pypi.org', true, false, '/logos/pypi.svg'),
    ('npm', 'registry.npmjs.org', true, false, '/logos/npm.svg'),
    ('WordPress', 'wordpress.org', true, false, '/logos/wordpress.svg'),
    ('Go modules', 'proxy.golang.org', true, false, '/logos/golang.svg'),
    ('crates.io', 'crates.io', true, false, '/logos/crates.svg'),
    ('VS Code', 'update.code.visualstudio.com', true, false, '/logos/vscode.svg'),
    ('NuGet', 'api.nuget.org', true, false, '/logos/nuget.svg'),
    ('Instagram', 'instagram.com', true, false, '/logos/instagram.svg'),
    ('Telegram', 'telegram.org', true, false, '/logos/telegram.svg')
) AS seed(name, host, http_enabled, icmp_enabled, logo_icon)
JOIN probes p ON p.code = 'probe1'
WHERE NOT EXISTS (SELECT 1 FROM endpoints e WHERE e.host = seed.host);

UPDATE endpoints e SET
    name = seed.name,
    http_enabled = seed.http_enabled,
    icmp_enabled = seed.icmp_enabled,
    probe_id = p.id,
    active = true,
    logo_icon = seed.logo_icon
FROM (VALUES
    ('Google', 'google.com', true, false, '/logos/google.svg'),
    ('GitHub', 'github.com', true, false, '/logos/github.svg'),
    ('1.1.1.1', '1.1.1.1', false, true, '/logos/cloudflare.svg'),
    ('8.8.8.8', '8.8.8.8', false, true, '/logos/google-dns.svg'),
    ('YouTube', 'youtube.com', true, false, '/logos/youtube.svg'),
    ('ChatGPT', 'chatgpt.com', true, false, '/logos/chatgpt.svg'),
    ('Claude', 'claude.ai', true, false, '/logos/claude.svg'),
    ('DeepSeek', 'deepseek.com', true, false, '/logos/deepseek.svg'),
    ('Microsoft', 'microsoft.com', true, false, '/logos/microsoft.svg'),
    ('Apple', 'apple.com', true, false, '/logos/apple.svg'),
    ('Play Store', 'play.google.com', true, false, '/logos/play-store.svg'),
    ('App Store', 'apps.apple.com', true, false, '/logos/app-store.svg'),
    ('Docker', 'docker.com', true, false, '/logos/docker.svg'),
    ('Spotify', 'spotify.com', true, false, '/logos/spotify.svg'),
    ('Grok', 'grok.com', true, false, '/logos/grok.svg'),
    ('Arena.ai', 'arena.ai', true, false, '/logos/arena.svg'),
    ('Ondpline.com', 'ondpline.com', true, false, '/logos/ondpline.svg'),
    ('GitLab', 'gitlab.com', true, false, '/logos/gitlab.svg'),
    ('digikala.ir', 'digikala.ir', true, false, '/logos/digikala.svg'),
    ('snapp.ir', 'snapp.ir', true, false, '/logos/snapp.svg'),
    ('tapsi.ir', 'tapsi.ir', true, false, '/logos/tapsi.svg'),
    ('divar.ir', 'divar.ir', true, false, '/logos/divar.svg'),
    ('bale.ai', 'bale.ai', true, false, '/logos/bale.svg'),
    ('Devin AI', 'devin.ai', true, false, '/logos/devin.svg'),
    ('Neshan.org', 'neshan.org', true, false, '/logos/neshan.svg'),
    ('Wikipedia', 'wikipedia.org', true, false, '/logos/wikipedia.svg'),
    ('WhatsApp', 'whatsapp.com', true, false, '/logos/whatsapp.svg'),
    ('zoomit.ir', 'zoomit.ir', true, false, '/logos/zoomit.svg'),
    ('digiato.ir', 'digiato.ir', true, false, '/logos/digiato.svg'),
    ('Let''s Encrypt', 'letsencrypt.org', true, false, '/logos/letsencrypt.svg'),
    ('Steam', 'steampowered.com', true, false, '/logos/steam.svg'),
    ('torob.ir', 'torob.ir', true, false, '/logos/torob.svg'),
    ('PlayStation', 'playstation.com', true, false, '/logos/playstation.svg'),
    ('Xbox', 'xbox.com', true, false, '/logos/xbox.svg'),
    ('Gmail', 'gmail.com', true, false, '/logos/gmail.svg'),
    ('Figma', 'figma.com', true, false, '/logos/figma.svg'),
    ('Windows Update', 'windowsupdate.microsoft.com', true, false, '/logos/windows-update.svg'),
    ('Bing', 'bing.com', true, false, '/logos/bing.svg'),
    ('aparat.ir', 'aparat.ir', true, false, '/logos/aparat.svg'),
    ('PyPI', 'pypi.org', true, false, '/logos/pypi.svg'),
    ('npm', 'registry.npmjs.org', true, false, '/logos/npm.svg'),
    ('WordPress', 'wordpress.org', true, false, '/logos/wordpress.svg'),
    ('Go modules', 'proxy.golang.org', true, false, '/logos/golang.svg'),
    ('crates.io', 'crates.io', true, false, '/logos/crates.svg'),
    ('VS Code', 'update.code.visualstudio.com', true, false, '/logos/vscode.svg'),
    ('NuGet', 'api.nuget.org', true, false, '/logos/nuget.svg'),
    ('Instagram', 'instagram.com', true, false, '/logos/instagram.svg'),
    ('Telegram', 'telegram.org', true, false, '/logos/telegram.svg')
) AS seed(name, host, http_enabled, icmp_enabled, logo_icon)
JOIN probes p ON p.code = 'probe1'
WHERE e.host = seed.host;

UPDATE endpoints SET active = false
WHERE host NOT IN (
    'google.com', 'github.com', '1.1.1.1', '8.8.8.8', 'youtube.com',
    'chatgpt.com', 'claude.ai', 'deepseek.com', 'microsoft.com', 'apple.com',
    'play.google.com', 'apps.apple.com', 'docker.com', 'spotify.com',
    'grok.com', 'arena.ai', 'ondpline.com',
    'gitlab.com', 'digikala.ir', 'snapp.ir', 'tapsi.ir', 'divar.ir',
    'bale.ai', 'devin.ai', 'neshan.org', 'wikipedia.org',
    'whatsapp.com', 'zoomit.ir', 'digiato.ir', 'letsencrypt.org', 'steampowered.com',
    'torob.ir', 'playstation.com', 'xbox.com', 'gmail.com', 'figma.com',
    'windowsupdate.microsoft.com', 'bing.com', 'aparat.ir', 'pypi.org',
    'registry.npmjs.org', 'wordpress.org', 'proxy.golang.org', 'crates.io',
    'update.code.visualstudio.com', 'api.nuget.org', 'instagram.com', 'telegram.org'
);
