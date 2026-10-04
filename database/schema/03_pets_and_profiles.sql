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
