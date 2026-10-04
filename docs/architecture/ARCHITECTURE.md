# Pet Digital Identity Platform — Architecture Documentation

## 1. Executive Architecture Overview

PetIdentity is an API-first, event-driven pet digital identity and lifetime-records platform designed for enterprise scale. Every pet receives a permanent, tamper-resistant digital identity linked to physical QR / NFC collar medals.

```text
                       [ Cloudflare CDN / WAF ]
                                  │
                       [ Ingress / Load Balancer ]
                                  │
         ┌────────────────────────┼────────────────────────┐
         ▼                        ▼                        ▼
    Go REST API              Go REST API              Go REST API
  (Stateless Pod)          (Stateless Pod)          (Stateless Pod)
         │                        │                        │
         └────────────────────────┼────────────────────────┘
                                  │
        ┌─────────────────────────┼─────────────────────────┐
        ▼                         ▼                         ▼
   [ PostgreSQL ]              [ Redis ]            [ Object Storage ]
 (Managed RDS + PostGIS)     (Cluster Cache)          (S3 / R2 Bucket)
        │
   Outbox Events
        │
        ▼
   [ Outbox Publisher ] ──► [ Event Bus ]
                                  │
             ┌────────────────────┼────────────────────┐
             ▼                    ▼                    ▼
     [ Push Worker ]       [ Email Worker ]     [ Analytics / DLQ ]
```

---

## 2. API-First Invariant

The Go REST API is the single source of truth and business logic:
* **Web Client**: Ultra-fast HTML5 + CSS3 + Vanilla JS consuming `/api/v1/*`.
* **Mobile Client**: Production-grade Flutter application compiling natively to Android and iOS.
* **API Contract**: Formalized in `backend/openapi.yaml` prior to client implementation.

---

## 3. Transactional Outbox Pattern

All critical database mutations write events into `outbox_events` in the **same database transaction**:

```sql
BEGIN;
  -- Mutation
  UPDATE pets SET status = 'LOST' WHERE id = '...';
  -- Outbox Event
  INSERT INTO outbox_events (id, event_type, aggregate_type, aggregate_id, payload)
  VALUES (gen_random_uuid(), 'PetMarkedLost', 'PET', '...', '{"petId": "..."}');
COMMIT;
```

A dedicated lightweight background publisher scans `outbox_events` and enqueues jobs with exponential backoff and dead-letter queues (DLQ).

---

## 4. Privacy & Data Segregation Model

1. **Public Level**: Safe information shown when tag is scanned (Pet photo, name, species, breed, medical warnings, safe in-app contact form).
2. **Private Level**: Microchip numbers (unless owner toggles ON), vet notes, clinical invoices, private home address.
3. **System Level**: Password hashes, refresh token signatures, OTP rate-limiting attempts, security audit logs.
