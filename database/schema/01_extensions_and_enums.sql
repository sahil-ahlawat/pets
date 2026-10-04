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
