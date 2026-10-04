# PostgreSQL Database Schema Documentation

## Core Tables Overview

### 1. Identity & Auth
* `users` — Primary user identity, authentication hashes, and global privacy settings.
* `user_sessions` — Refresh token hashes and active session tracking.
* `user_devices` — Registered FCM/APNs push notification device tokens.
* `user_email_verifications` — Secure email confirmation tokens.
* `user_otp_attempts` — Rate-limited one-time password attempts.

### 2. Companion Profiles
* `pets` — Core pet entity (name, species, breed, gender, neutered, DOB, weight, avatar).
* `pet_users` — Multi-owner relationship (Primary owner, Co-owner, Vet, Caregiver).
* `pet_profiles` — Personality notes, dietary preferences, and Privacy Shield toggles.
* `pet_emergency_profiles` — Secondary emergency contacts, critical medical warnings.
* `pet_identifications` — ISO microchip numbers, municipal licenses, tattoo IDs.
* `pet_photos` — Gallery images and display order.
* `pet_documents` — File metadata pointing to private S3/R2 object storage.

### 3. Health & Preventive Records
* `vaccination_types` — Global WSAVA-aligned vaccine catalog.
* `vaccination_schedules` — Standard booster timing protocols.
* `pet_vaccinations` — Immunization records with lot numbers and due dates.
* `pet_medical_events` — Vet visits, diagnoses, surgeries, and dental exams.
* `pet_medications` — Active prescriptions, dosages, and administration frequencies.
* `pet_reminders` — Preventive care reminders (deworming, flea/tick, grooming).
* `pet_events` — Consolidated chronological timeline history.

### 4. Smart QR Tags & Scanning
* `pet_qr_codes` — Random 128-bit unguessable tokens, tag status, and scan counts.
* `qr_scans` — Scan telemetry (timestamps, IP hashes, approximate coordinates, user agents).

### 5. Lost & Found Network
* `pet_status_history` — State machine audits (ACTIVE, LOST, FOUND, REUNITED).
* `lost_pet_reports` — Last seen locations, rewards, and circumstances.
* `found_pet_reports` — Community sighting submissions with coordinates and photos.

### 6. Veterinary Practices
* `vets` — Directory with addresses, 24h emergency tags, and coordinates.
* `pet_vets` — Association between companion and their primary vet practice.

### 7. Scalability & System Integrity
* `outbox_events` — Transactional outbox ledger for zero-loss asynchronous events.
* `outbox_consumer_tracking` — Consumer idempotency table.
* `jobs` — Asynchronous priority job queues with exponential backoff and DLQ.
* `job_attempts` — Execution diagnostics.
* `audit_logs` — System security audit trails.
