-- 01_extensions_and_enums.sql
-- Extensions and custom enumerated types for Pet Digital Identity Platform

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
-- Enable PostGIS if available (falls back gracefully if extension not loaded)
DO $$
BEGIN
    CREATE EXTENSION IF NOT EXISTS "postgis";
EXCEPTION
    WHEN others THEN
        RAISE NOTICE 'PostGIS not installed; coordinates will use FLOAT columns';
END $$;

-- User Roles
CREATE TYPE user_role AS ENUM (
    'PET_OWNER',
    'VET',
    'SHELTER_ADMIN',
    'SYSTEM_ADMIN'
);

-- Pet Status
CREATE TYPE pet_status AS ENUM (
    'ACTIVE',
    'LOST',
    'FOUND',
    'REUNITED',
    'DECEASED',
    'ARCHIVED'
);

-- Pet Species
CREATE TYPE pet_species AS ENUM (
    'DOG',
    'CAT',
    'BIRD',
    'RABBIT',
    'HORSE',
    'REPTILE',
    'OTHER'
);

-- Pet Gender
CREATE TYPE pet_gender AS ENUM (
    'MALE',
    'FEMALE',
    'UNKNOWN'
);

-- Pet Co-Ownership Roles
CREATE TYPE pet_user_role AS ENUM (
    'PRIMARY_OWNER',
    'CO_OWNER',
    'CAREGIVER',
    'VET_ACCESS'
);

-- QR Code Status
CREATE TYPE qr_status AS ENUM (
    'ACTIVE',
    'INACTIVE',
    'REVOKED',
    'REPLACED'
);

-- Medical Event Type
CREATE TYPE medical_event_type AS ENUM (
    'CONSULTATION',
    'DIAGNOSIS',
    'TREATMENT',
    'MEDICATION',
    'SURGERY',
    'LAB_TEST',
    'DENTAL',
    'WEIGHT_CHECK',
    'OTHER'
);

-- Reminder Type
CREATE TYPE reminder_type AS ENUM (
    'VACCINATION',
    'MEDICATION',
    'DEWORMING',
    'FLEA_TICK',
    'VET_APPOINTMENT',
    'GROOMING',
    'CUSTOM'
);

-- Reminder Frequency
CREATE TYPE recurrence_rule AS ENUM (
    'ONCE',
    'DAILY',
    'WEEKLY',
    'MONTHLY',
    'QUARTERLY',
    'BIANNUAL',
    'YEARLY'
);

-- Notification Channels
CREATE TYPE notification_channel AS ENUM (
    'PUSH',
    'EMAIL',
    'SMS',
    'IN_APP'
);

-- Notification Types
CREATE TYPE notification_type AS ENUM (
    'VACCINATION_DUE',
    'MEDICATION_DUE',
    'REMINDER',
    'QR_SCANNED',
    'PET_FOUND',
    'PET_LOST',
    'OWNERSHIP_TRANSFER',
    'SECURITY',
    'SYSTEM'
);

-- Job Status
CREATE TYPE job_status AS ENUM (
    'QUEUED',
    'PROCESSING',
    'COMPLETED',
    'FAILED',
    'DEAD_LETTER'
);

-- Device Platform
CREATE TYPE device_platform AS ENUM (
    'ANDROID',
    'IOS',
    'WEB'
);
-- 02_users_and_auth.sql
-- User management, authentication, sessions, devices, and OTP verification

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone_number VARCHAR(30),
    phone_verified BOOLEAN NOT NULL DEFAULT FALSE,
    email_verified BOOLEAN NOT NULL DEFAULT FALSE,
    role user_role NOT NULL DEFAULT 'PET_OWNER',
    avatar_url TEXT,
    timezone VARCHAR(50) NOT NULL DEFAULT 'UTC',
    preferred_language VARCHAR(10) NOT NULL DEFAULT 'en',
    allow_direct_calls BOOLEAN NOT NULL DEFAULT FALSE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token_hash VARCHAR(255) NOT NULL UNIQUE,
    user_agent TEXT,
    client_ip VARCHAR(45),
    expires_at TIMESTAMPTZ NOT NULL,
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_token TEXT NOT NULL,
    platform device_platform NOT NULL,
    device_name VARCHAR(150),
    app_version VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_user_device_token UNIQUE (user_id, device_token)
);

CREATE TABLE IF NOT EXISTS user_email_verifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_otp_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    target VARCHAR(255) NOT NULL, -- email or phone
    otp_hash VARCHAR(255) NOT NULL,
    purpose VARCHAR(50) NOT NULL, -- 'LOGIN', 'VERIFY_PHONE', 'RESET_PASSWORD'
    attempts_count INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 5,
    expires_at TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_phone ON users(phone_number) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_sessions_user ON user_sessions(user_id) WHERE NOT is_revoked;
CREATE INDEX idx_user_devices_user ON user_devices(user_id) WHERE is_active = TRUE;
CREATE INDEX idx_user_otp_target ON user_otp_attempts(target, purpose) WHERE consumed_at IS NULL;
-- 03_pets_and_profiles.sql
-- Pet entities, profiles, emergency settings, identifications, photos, and documents

CREATE TABLE IF NOT EXISTS pets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    species pet_species NOT NULL DEFAULT 'DOG',
    breed VARCHAR(150),
    gender pet_gender NOT NULL DEFAULT 'UNKNOWN',
    is_neutered BOOLEAN NOT NULL DEFAULT FALSE,
    date_of_birth DATE,
    approximate_age_years INT,
    color_markings VARCHAR(255),
    weight_kg NUMERIC(5, 2),
    avatar_url TEXT,
    status pet_status NOT NULL DEFAULT 'ACTIVE',
    is_archived BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS pet_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role pet_user_role NOT NULL DEFAULT 'PRIMARY_OWNER',
    can_edit BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_pet_user UNIQUE (pet_id, user_id)
);

CREATE TABLE IF NOT EXISTS pet_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL UNIQUE REFERENCES pets(id) ON DELETE CASCADE,
    bio TEXT,
    dietary_preferences TEXT,
    favorite_toys_activities TEXT,
    behavioral_notes TEXT,
    allergies TEXT,
    chronic_conditions TEXT,
    -- Privacy Settings (Control what public sees)
    show_owner_name BOOLEAN NOT NULL DEFAULT TRUE,
    show_owner_phone BOOLEAN NOT NULL DEFAULT FALSE,
    show_microchip BOOLEAN NOT NULL DEFAULT TRUE,
    show_allergies_emergency BOOLEAN NOT NULL DEFAULT TRUE,
    show_vet_info BOOLEAN NOT NULL DEFAULT FALSE,
    custom_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pet_emergency_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL UNIQUE REFERENCES pets(id) ON DELETE CASCADE,
    emergency_contact_name VARCHAR(150),
    emergency_contact_phone VARCHAR(30),
    emergency_contact_relation VARCHAR(50),
    critical_medical_alert TEXT,
    special_care_instructions TEXT,
    insurance_provider VARCHAR(150),
    insurance_policy_number VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pet_identifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    id_type VARCHAR(50) NOT NULL, -- 'MICROCHIP', 'TATTOO', 'MUNICIPAL_LICENSE', 'PEDIGREE_REG'
    id_value VARCHAR(100) NOT NULL,
    issuer_registry VARCHAR(150),
    implant_date DATE,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_pet_id_type UNIQUE (pet_id, id_type, id_value)
);

CREATE TABLE IF NOT EXISTS pet_photos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    photo_url TEXT NOT NULL,
    thumbnail_url TEXT,
    caption VARCHAR(255),
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pet_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    uploaded_by UUID NOT NULL REFERENCES users(id),
    title VARCHAR(200) NOT NULL,
    category VARCHAR(50) NOT NULL, -- 'VACCINE_CERT', 'ADOPTION_PAPERS', 'INSURANCE', 'LAB_RESULT', 'PRESCRIPTION'
    storage_key TEXT NOT NULL,
    file_url TEXT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    file_size_bytes BIGINT NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_pets_status ON pets(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_pet_users_user ON pet_users(user_id);
CREATE INDEX idx_pet_users_pet ON pet_users(pet_id);
CREATE INDEX idx_pet_photos_pet ON pet_photos(pet_id);
CREATE INDEX idx_pet_documents_pet ON pet_documents(pet_id);
CREATE INDEX idx_pet_identifications_val ON pet_identifications(id_value);
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
-- 07_vets_and_clinics.sql
-- Veterinary practices, clinics directory, and pet-vet associations

CREATE TABLE IF NOT EXISTS vets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_name VARCHAR(200) NOT NULL,
    doctor_name VARCHAR(150),
    license_number VARCHAR(100),
    phone_number VARCHAR(30) NOT NULL,
    emergency_phone VARCHAR(30),
    email VARCHAR(255),
    website_url TEXT,
    address_line1 VARCHAR(255) NOT NULL,
    address_line2 VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    postal_code VARCHAR(20),
    country_code VARCHAR(10) NOT NULL DEFAULT 'US',
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    is_24h_emergency BOOLEAN NOT NULL DEFAULT FALSE,
    specializations TEXT[],
    rating NUMERIC(3, 2) DEFAULT 5.00,
    review_count INT DEFAULT 0,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pet_vets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    vet_id UUID NOT NULL REFERENCES vets(id) ON DELETE CASCADE,
    is_primary BOOLEAN NOT NULL DEFAULT TRUE,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_pet_vet UNIQUE (pet_id, vet_id)
);

CREATE INDEX idx_vets_city ON vets(city);
CREATE INDEX idx_vets_emergency ON vets(is_24h_emergency) WHERE is_24h_emergency = TRUE;
CREATE INDEX idx_pet_vets_pet ON pet_vets(pet_id);
-- 08_notifications_and_devices.sql
-- Notification management and delivery audit logs

CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pet_id UUID REFERENCES pets(id) ON DELETE SET NULL,
    notification_type notification_type NOT NULL,
    title VARCHAR(200) NOT NULL,
    body TEXT NOT NULL,
    data JSONB,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS notification_deliveries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notification_id UUID NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
    channel notification_channel NOT NULL,
    recipient_target TEXT NOT NULL, -- e.g., device token, email, phone
    status VARCHAR(50) NOT NULL, -- 'PENDING', 'SENT', 'FAILED'
    external_message_id VARCHAR(255),
    error_message TEXT,
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user ON notifications(user_id, is_read, created_at DESC);
CREATE INDEX idx_notifications_pet ON notifications(pet_id);
CREATE INDEX idx_notification_deliveries_status ON notification_deliveries(status);
-- 09_outbox_jobs_and_audit.sql
-- Transactional Outbox pattern, Asynchronous Jobs, and Security Audit Logs

CREATE TABLE IF NOT EXISTS outbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(100) NOT NULL,
    aggregate_type VARCHAR(50) NOT NULL, -- 'PET', 'USER', 'QR', 'VACCINATION'
    aggregate_id UUID NOT NULL,
    payload JSONB NOT NULL,
    headers JSONB,
    published_at TIMESTAMPTZ,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    retry_count INT NOT NULL DEFAULT 0,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS outbox_consumer_tracking (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES outbox_events(id) ON DELETE CASCADE,
    consumer_name VARCHAR(100) NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_event_consumer UNIQUE (event_id, consumer_name)
);

CREATE TABLE IF NOT EXISTS jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    queue VARCHAR(50) NOT NULL DEFAULT 'default',
    job_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    status job_status NOT NULL DEFAULT 'QUEUED',
    run_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    max_retries INT NOT NULL DEFAULT 5,
    retry_count INT NOT NULL DEFAULT 0,
    last_error TEXT,
    locked_by VARCHAR(100),
    locked_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS job_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    attempt_number INT NOT NULL,
    status job_status NOT NULL,
    error_message TEXT,
    duration_ms INT,
    started_at TIMESTAMPTZ NOT NULL,
    finished_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID,
    ip_address VARCHAR(45),
    user_agent TEXT,
    old_state JSONB,
    new_state JSONB,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_outbox_unpublished ON outbox_events(created_at) WHERE is_published = FALSE;
CREATE INDEX idx_jobs_pending ON jobs(status, run_at) WHERE status = 'QUEUED';
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id, created_at DESC);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
-- 10_content_and_blog.sql
-- Educational blog content, pet health guides, and category management

CREATE TABLE IF NOT EXISTS blog_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    icon_name VARCHAR(50),
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS blog_posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID REFERENCES blog_categories(id) ON DELETE SET NULL,
    author_id UUID REFERENCES users(id) ON DELETE SET NULL,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    excerpt TEXT NOT NULL,
    content_markdown TEXT NOT NULL,
    featured_image_url TEXT,
    reading_time_minutes INT DEFAULT 5,
    is_published BOOLEAN NOT NULL DEFAULT FALSE,
    published_at TIMESTAMPTZ,
    views_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_blog_published ON blog_posts(is_published, published_at DESC) WHERE is_published = TRUE;
CREATE INDEX idx_blog_slug ON blog_posts(slug);
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
-- seeds.sql
-- Initial seed data for Pet Digital Identity Platform

-- Standard Core Canine Vaccination Types
INSERT INTO vaccination_types (species, name, code, is_core, recommended_interval_months, description)
VALUES 
    ('DOG', 'Rabies (Canine 1-Year/3-Year)', 'RABIES_CANINE', TRUE, 36, 'Mandatory rabies vaccination protecting against the rabies virus.'),
    ('DOG', 'DHPP (Distemper, Hepatitis, Parvovirus, Parainfluenza)', 'DHPP', TRUE, 12, 'Core 5-in-1 combo vaccine guarding puppies and dogs against fatal viral infections.'),
    ('DOG', 'Bordetella (Kennel Cough)', 'BORDETELLA', FALSE, 6, 'Recommended for dogs frequenting boarding, daycare, or parks.'),
    ('DOG', 'Leptospirosis', 'LEPTO', FALSE, 12, 'Bacterial infection protection, especially for outdoor and water-loving dogs.'),
    ('DOG', 'Lyme Disease', 'LYME', FALSE, 12, 'Protects against tick-borne Borrelia burgdorferi.')
ON CONFLICT (species, code) DO NOTHING;

-- Standard Core Feline Vaccination Types
INSERT INTO vaccination_types (species, name, code, is_core, recommended_interval_months, description)
VALUES 
    ('CAT', 'Rabies (Feline Purevax/Standard)', 'RABIES_FELINE', TRUE, 12, 'Crucial rabies prophylaxis for domestic and indoor/outdoor cats.'),
    ('CAT', 'FVRCP (Feline Viral Rhinotracheitis, Calicivirus, Panleukopenia)', 'FVRCP', TRUE, 36, 'Core 3-in-1 feline respiratory and distemper protection.'),
    ('CAT', 'FeLV (Feline Leukemia Virus)', 'FELV', FALSE, 12, 'Highly recommended for cats with outdoor access or multiple cat households.')
ON CONFLICT (species, code) DO NOTHING;

-- Seed Sample Veterinary Clinics
INSERT INTO vets (
    clinic_name, doctor_name, license_number, phone_number, emergency_phone,
    email, website_url, address_line1, city, state, postal_code, country_code,
    latitude, longitude, is_24h_emergency, specializations, rating, review_count, is_verified
) VALUES 
(
    'Metropolitan Animal Emergency & Specialty Hospital',
    'Dr. Elena Rostova, DVM, DACVECC',
    'VET-883921',
    '+1 (555) 234-5678',
    '+1 (555) 234-9999',
    'care@metrovet-emergency.com',
    'https://metrovet-emergency.com',
    '742 Evergreen Promenade, Suite 100',
    'San Francisco',
    'CA',
    '94107',
    'US',
    37.7749,
    -122.4194,
    TRUE,
    ARRAY['Emergency Medicine', 'Critical Care', 'Orthopedic Surgery', 'Ultrasound'],
    4.95,
    148,
    TRUE
),
(
    'Paws & Claws Holistic Veterinary Center',
    'Dr. Marcus Vance, MRCVS',
    'VET-419082',
    '+1 (555) 789-0123',
    '+1 (555) 789-0124',
    'hello@pawsandclaws-vet.com',
    'https://pawsandclaws-vet.com',
    '1248 Oak Boulevard',
    'San Francisco',
    'CA',
    '94110',
    'US',
    37.7599,
    -122.4148,
    FALSE,
    ARRAY['General Practice', 'Dental Surgery', 'Feline Medicine', 'Preventative Care'],
    4.88,
    92,
    TRUE
)
ON CONFLICT DO NOTHING;

-- Seed Blog Categories
INSERT INTO blog_categories (slug, name, description, icon_name, display_order)
VALUES 
    ('pet-health', 'Pet Health & Wellness', 'Expert veterinary guidance on vaccinations, preventive care, and wellness.', 'heart-pulse', 1),
    ('lost-and-found', 'Lost Pet Recovery', 'Action plans, search strategies, and recovery stories to keep pets safe.', 'shield-alert', 2),
    ('smart-tags', 'Smart QR & NFC Technology', 'How modern digital pet identity outperforms obsolete metal collar tags.', 'qr-code', 3),
    ('nutrition-lifestyle', 'Nutrition & Lifestyle', 'Dietary recommendations, exercise, and enriching your pet companionship.', 'bone', 4)
ON CONFLICT (slug) DO NOTHING;

-- Seed Sample Blog Post
INSERT INTO blog_posts (
    title, slug, excerpt, content_markdown, reading_time_minutes, is_published, published_at, views_count
) VALUES (
    'Why Metal Tags Fail: The Digital QR Revolution in Pet Safety',
    'why-metal-tags-fail-digital-qr-revolution',
    'Traditional metal tags scratch, wear off, and reveal only a static phone number. Here is how modern QR digital identity instantly re-unites lost pets.',
    '# Why Metal Tags Fail: The Digital QR Revolution in Pet Safety

Every year, millions of family pets wander away or slip out through an open gate. In past decades, their only line of defense was an engraved metal disc stamped with a phone number that quickly becomes scratched and illegible.

### The Limitations of Legacy Tags
* **Static Information**: When you travel, move houses, or change phone numbers, metal tags are immediately obsolete.
* **Zero Medical Context**: If your pet is diabetic or has severe medication allergies, a regular tag communicates nothing to the rescuer.
* **Privacy Risks**: Displaying your private home address and mobile number on an open collar invites privacy and stalking concerns.

### The Power of Dynamic Digital QR Identity
With dynamic QR tags, scanning the collar opens a **secure, real-time responsive profile**. The owner can update contact details, activate **Lost Mode**, and receive instant scan alerts with approximate location coordinates.',
    4,
    TRUE,
    CURRENT_TIMESTAMP,
    342
)
ON CONFLICT (slug) DO NOTHING;
