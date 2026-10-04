# Database Partitioning Strategy

## High-Volume Target Tables
In high-volume enterprise deployments (Stage 4), three tables grow continuously:
1. `qr_scans` — Every physical collar tag scan creates a log entry with location telemetry.
2. `outbox_events` — Every database mutation generates an outbox row.
3. `audit_logs` — Every security-sensitive action writes an audit entry.

## Native PostgreSQL Range Partitioning Design

### 1. `qr_scans` Partitioning Implementation
```sql
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

-- Monthly partitions
CREATE TABLE qr_scans_2026_10 PARTITION OF qr_scans_partitioned
    FOR VALUES FROM ('2026-10-01 00:00:00+00') TO ('2026-11-01 00:00:00+00');

CREATE TABLE qr_scans_2026_11 PARTITION OF qr_scans_partitioned
    FOR VALUES FROM ('2026-11-01 00:00:00+00') TO ('2026-12-01 00:00:00+00');
```

### 2. Automated Partition Maintenance
Use `pg_partman` or Kubernetes cron jobs to pre-create partitions 3 months in advance and detach/archive partitions older than 2 years into cold object storage (Parquet / S3 Glacier).
