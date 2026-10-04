-- 06_lost_and_found.sql
-- Status tracking, lost pet alerts, and community found reports

CREATE TABLE IF NOT EXISTS pet_status_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    previous_status pet_status NOT NULL,
    new_status pet_status NOT NULL,
    changed_by UUID REFERENCES users(id) ON DELETE SET NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS lost_pet_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    reported_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lost_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_seen_address TEXT NOT NULL,
    last_seen_city VARCHAR(100),
    last_seen_state VARCHAR(100),
    last_seen_postal_code VARCHAR(20),
    last_seen_latitude DOUBLE PRECISION,
    last_seen_longitude DOUBLE PRECISION,
    circumstances TEXT,
    collar_description TEXT,
    reward_amount VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS found_pet_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID REFERENCES pets(id) ON DELETE SET NULL,
    qr_code_id UUID REFERENCES pet_qr_codes(id) ON DELETE SET NULL,
    finder_name VARCHAR(150),
    finder_phone VARCHAR(30),
    finder_email VARCHAR(255),
    message TEXT NOT NULL,
    current_pet_status VARCHAR(100), -- 'SAFE_WITH_ME', 'AT_VET', 'SIGHTED_ONLY'
    found_latitude DOUBLE PRECISION,
    found_longitude DOUBLE PRECISION,
    found_address TEXT,
    photo_url TEXT,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    notified_owner_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_lost_pet_active ON lost_pet_reports(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_lost_pet_pet ON lost_pet_reports(pet_id);
CREATE INDEX idx_found_pet_pet ON found_pet_reports(pet_id, created_at DESC);
