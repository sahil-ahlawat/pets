-- 05_qr_and_scanning.sql
-- QR tag management, token generation, and secure scanning logs

CREATE TABLE IF NOT EXISTS pet_qr_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    token VARCHAR(64) NOT NULL UNIQUE, -- Cryptographically secure, unguessable random token
    status qr_status NOT NULL DEFAULT 'ACTIVE',
    tag_label VARCHAR(100) DEFAULT 'Standard Collar Tag',
    physical_tag_type VARCHAR(50) DEFAULT 'QR_COLLAR', -- 'QR_COLLAR', 'NFC_TAG', 'SMART_MEDAL'
    scan_count INT NOT NULL DEFAULT 0,
    last_scanned_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    revocation_reason VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS qr_scans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    qr_code_id UUID NOT NULL REFERENCES pet_qr_codes(id) ON DELETE CASCADE,
    pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    scanned_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    approx_latitude DOUBLE PRECISION,
    approx_longitude DOUBLE PRECISION,
    accuracy_meters DOUBLE PRECISION,
    approx_address TEXT,
    city VARCHAR(100),
    country_code VARCHAR(10),
    ip_hash VARCHAR(64), -- Hashed for GDPR / privacy compliance
    user_agent TEXT,
    referrer TEXT,
    is_suspicious BOOLEAN NOT NULL DEFAULT FALSE,
    metadata JSONB
);

CREATE INDEX idx_pet_qr_token ON pet_qr_codes(token);
CREATE INDEX idx_pet_qr_pet ON pet_qr_codes(pet_id);
CREATE INDEX idx_qr_scans_qr ON qr_scans(qr_code_id, scanned_at DESC);
CREATE INDEX idx_qr_scans_pet ON qr_scans(pet_id, scanned_at DESC);
