-- 11_indexes_and_partitioning.sql
-- Performance indexes and time-range partitioning templates for hyper-scale tables:
--   - qr_scans
--   - outbox_events
--   - audit_logs

-- Composite indexes for high-frequency queries
CREATE INDEX IF NOT EXISTS idx_pet_vaccinations_lookup 
    ON pet_vaccinations (pet_id, next_due_date ASC);

CREATE INDEX IF NOT EXISTS idx_pet_reminders_active_due 
    ON pet_reminders (user_id, due_date ASC) 
    WHERE is_completed = FALSE;

CREATE INDEX IF NOT EXISTS idx_qr_scans_analytics 
    ON qr_scans (pet_id, scanned_at DESC);

-- Partitioning notes for scale:
-- In Stage 4 high-scale deployments, the qr_scans, audit_logs, and outbox_events tables
-- should be created as PARTITION BY RANGE (created_at / scanned_at).
-- Example implementation for PostgreSQL native partitioning:
/*
CREATE TABLE qr_scans_partitioned (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    qr_code_id UUID NOT NULL,
    pet_id UUID NOT NULL,
    scanned_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    approx_latitude DOUBLE PRECISION,
    approx_longitude DOUBLE PRECISION,
    accuracy_meters DOUBLE PRECISION,
    approx_address TEXT,
    city VARCHAR(100),
    country_code VARCHAR(10),
    ip_hash VARCHAR(64),
    user_agent TEXT,
    referrer TEXT,
    is_suspicious BOOLEAN NOT NULL DEFAULT FALSE,
    metadata JSONB,
    PRIMARY KEY (id, scanned_at)
) PARTITION BY RANGE (scanned_at);

-- Monthly partition creation automation via pg_partman or cron
CREATE TABLE qr_scans_y2026m10 PARTITION OF qr_scans_partitioned
    FOR VALUES FROM ('2026-10-01 00:00:00+00') TO ('2026-11-01 00:00:00+00');
*/
