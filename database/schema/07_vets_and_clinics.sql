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
