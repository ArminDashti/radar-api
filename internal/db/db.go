package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ArminDashti/radar-api/internal/auth"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const (
	seedUsername   = "armin"
	seedPassword   = "dopadopa123"
	seedAgentToken = "radar-demo-agent-token-probe1"
)

func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return pool, nil
}

func RunMigration(ctx context.Context, pool *pgxpool.Pool, path string) error {
	sql, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read migration: %w", err)
	}
	if _, err := pool.Exec(ctx, string(sql)); err != nil {
		return fmt.Errorf("execute migration: %w", err)
	}
	return nil
}

func Seed(ctx context.Context, pool *pgxpool.Pool) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(seedPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash seed password: %w", err)
	}
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err = tx.Exec(ctx, `
		INSERT INTO users (username, password_hash) VALUES ($1, $2)
		ON CONFLICT (username) DO NOTHING`, seedUsername, string(passwordHash)); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `
		INSERT INTO agents (name, probe_id, token_hash)
		SELECT 'Demo agent probe1', id, $1 FROM probes WHERE code = 'probe1'
		ON CONFLICT (name) DO UPDATE SET token_hash = EXCLUDED.token_hash`,
		auth.HashAgentToken(seedAgentToken)); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `
		INSERT INTO samples (agent_id, endpoint_id, protocol, observed_at, latency_ms, ok)
		SELECT a.id, e.id, protocol.name, bucket,
		       CASE WHEN generated.ok THEN 20 + generated.latency_seed * 180 ELSE NULL END,
		       generated.ok
		FROM agents a
		JOIN probes p ON p.id = a.probe_id AND p.code = 'probe1'
		CROSS JOIN endpoints e
		CROSS JOIN (VALUES ('http'), ('icmp')) AS protocol(name)
		CROSS JOIN generate_series(
			date_trunc('minute', now()) - interval '119 minutes',
			date_trunc('minute', now()),
			interval '1 minute'
		) AS bucket
		CROSS JOIN LATERAL (
			SELECT random() >= 0.08 AS ok, random() AS latency_seed
			WHERE bucket IS NOT NULL
		) AS generated
		WHERE (protocol.name = 'http' AND e.http_enabled)
		   OR (protocol.name = 'icmp' AND e.icmp_enabled)
		ON CONFLICT (agent_id, endpoint_id, protocol, observed_at) DO NOTHING`); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, `UPDATE agents SET last_seen_at = $1 WHERE name = 'Demo agent probe1'`, time.Now()); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
