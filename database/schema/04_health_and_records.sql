-- 04_health_and_records.sql
-- Vaccinations, medical history, medications, reminders, and timeline events

CREATE TABLE IF NOT EXISTS vaccination_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    species pet_species NOT NULL,
    name VARCHAR(150) NOT NULL,
    code VARCHAR(50) NOT NULL, -- e.g., 'RABIES', 'DHPP', 'FVRCP'
    is_core BOOLEAN NOT NULL DEFAULT TRUE,
    recommended_interval_months INT NOT NULL DEFAULT 12,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_species_vac_code UNIQUE (species, code)
);

CREATE TABLE IF NOT EXISTS vaccination_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vaccination_type_id UUID NOT NULL REFERENCES vaccination_types(id) ON DELETE CASCADE,
    dose_number INT NOT NULL DEFAULT 1,
    min_age_weeks INT NOT NULL,
    booster_interval_months INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pet_vaccinations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    vaccination_type_id UUID REFERENCES vaccination_types(id) ON DELETE SET NULL,
    custom_vaccine_name VARCHAR(150),
    administered_at DATE NOT NULL,
    next_due_date DATE,
    vet_name VARCHAR(150),
    clinic_name VARCHAR(200),
    batch_lot_number VARCHAR(100),
    manufacturer VARCHAR(150),
    document_id UUID REFERENCES pet_documents(id) ON DELETE SET NULL,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pet_medical_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    event_type medical_event_type NOT NULL DEFAULT 'CONSULTATION',
    event_date DATE NOT NULL,
    title VARCHAR(200) NOT NULL,
    diagnosis TEXT,
    treatment_summary TEXT,
    vet_name VARCHAR(150),
    clinic_name VARCHAR(200),
    cost_amount NUMERIC(10, 2),
    cost_currency VARCHAR(3) DEFAULT 'USD',
    document_id UUID REFERENCES pet_documents(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pet_medications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    medication_name VARCHAR(150) NOT NULL,
    dosage VARCHAR(100) NOT NULL, -- e.g., '10mg'
    frequency VARCHAR(100) NOT NULL, -- e.g., 'Twice daily with meals'
    start_date DATE NOT NULL,
    end_date DATE,
    prescribed_by VARCHAR(150),
    reason TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pet_reminders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reminder_type reminder_type NOT NULL DEFAULT 'CUSTOM',
    title VARCHAR(200) NOT NULL,
    description TEXT,
    due_date TIMESTAMPTZ NOT NULL,
    recurrence recurrence_rule NOT NULL DEFAULT 'ONCE',
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    completed_at TIMESTAMPTZ,
    last_notified_at TIMESTAMPTZ,
    related_entity_type VARCHAR(50), -- 'VACCINATION', 'MEDICATION', etc.
    related_entity_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pet_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    event_category VARCHAR(50) NOT NULL, -- 'VACCINE', 'VET_VISIT', 'WEIGHT', 'LOST', 'FOUND', 'STATUS_CHANGE', 'NOTE'
    event_timestamp TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    icon_name VARCHAR(50),
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_pet_vaccinations_pet ON pet_vaccinations(pet_id, administered_at DESC);
CREATE INDEX idx_pet_vaccinations_next_due ON pet_vaccinations(next_due_date) WHERE next_due_date IS NOT NULL;
CREATE INDEX idx_pet_medical_pet ON pet_medical_events(pet_id, event_date DESC);
CREATE INDEX idx_pet_medications_pet ON pet_medications(pet_id) WHERE is_active = TRUE;
CREATE INDEX idx_pet_reminders_due ON pet_reminders(due_date) WHERE is_completed = FALSE;
CREATE INDEX idx_pet_events_timeline ON pet_events(pet_id, event_timestamp DESC);
