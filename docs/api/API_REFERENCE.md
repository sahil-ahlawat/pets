# API Reference Documentation

All endpoints follow RESTful standards and are rooted under `/api/v1/`.

## Authentication & Identity
* `POST /auth/register` — Register new owner account.
* `POST /auth/login` — Sign in with email and password. Returns JWT access token & rotating refresh token.
* `POST /auth/refresh` — Issue a fresh JWT pair via refresh token.
* `POST /auth/otp/request` — Request 6-digit one-time password.
* `POST /auth/otp/verify` — Validate OTP.
* `GET /auth/me` — Authenticated profile info.

## Companion Management
* `GET /pets` — List authenticated user's pets with cursor pagination.
* `POST /pets` — Create new companion digital identity.
* `GET /pets/:id` — Full private companion details.
* `PATCH /pets/:id` — Update identity and privacy shield toggles.
* `DELETE /pets/:id` — Archive/soft-delete companion.
* `POST /pets/:id/lost` — Activate Lost Mode emergency status.
* `POST /pets/:id/reunited` — Confirm safe reunion and deactivate emergency banner.

## Smart Collar Tags & Scanning
* `GET /pets/:id/qr` — List tags associated with companion.
* `POST /pets/:id/qr/generate` — Generate new cryptographic 128-bit QR tag.
* `POST /qr/:id/revoke` — Revoke lost or damaged tag.
* `GET /public/pets/:token` — **Public endpoint**. Returns privacy-filtered safe DTO for smartphone scans.
* `POST /public/pets/:token/contact` — Securely message owner without phone exposure.
* `POST /public/pets/:token/found` — Report found companion with optional GPS coordinates and photo.
* `GET /public/qr/:token.png` — Render high-res PNG image.

## Health, Clinical & Reminders
* `GET /vaccination-types` — Standardized core/optional vaccine catalog.
* `GET /pets/:id/vaccinations` — List companion immunization records.
* `POST /pets/:id/vaccinations` — Record vaccination dose.
* `GET /pets/:id/medical-events` — Consultations, surgeries, lab tests.
* `POST /pets/:id/medical-events` — Log medical event.
* `GET /pets/:id/medications` — Active prescriptions.
* `POST /pets/:id/medications` — Add medication.
* `GET /pets/:id/timeline` — Consolidated chronological life events.
* `GET /reminders` — List active companion preventive tasks.
* `POST /reminders` — Schedule reminder.
* `POST /reminders/:id/complete` — Mark completed.

## Veterinary Directory
* `GET /vets` — Search clinics and emergency practices.
* `GET /vets/nearby` — PostGIS geographic radius lookup.

## System & Diagnostics
* `GET /healthz` — Service liveness.
* `GET /readyz` — Database and dependency readiness.
* `GET /metrics` — Prometheus metrics scrape target.
