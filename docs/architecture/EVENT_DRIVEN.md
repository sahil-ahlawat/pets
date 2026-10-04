# Event-Driven Architecture & Outbox System

## 1. Event Catalog

| Event Name | Aggregate | Trigger Condition | Primary Downstream Consumers |
| :--- | :--- | :--- | :--- |
| `UserRegistered` | `USER` | Successful account creation | Welcome email, device sync |
| `PetCreated` | `PET` | Companion profile saved | Initial QR generation, welcome push |
| `PetUpdated` | `PET` | Profile or privacy toggled | Cache invalidation, tag profile update |
| `QRCreated` | `QR` | New QR code minted | Vector image pre-render |
| `QRScanned` | `QR` | Tag scanned by smartphone | Owner push notification, geo audit |
| `PetMarkedLost` | `PET` | Owner activates Lost Mode | Community push broadcast, shelter alerts |
| `PetFound` | `PET` | Finder submits found report | Urgent owner push + SMS dispatch |
| `PetReunited` | `PET` | Companion marked reunited | Alert resolution, success timeline entry |
| `VaccinationAdded` | `VACCINATION` | Shot logged | Booster countdown scheduling |
| `DocumentUploaded` | `DOCUMENT` | File uploaded to S3 | Malware scan, thumbnail generator |

---

## 2. Consumer Idempotency

Duplicate message deliveries are prevented through the `outbox_consumer_tracking` table:

```sql
CREATE TABLE outbox_consumer_tracking (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES outbox_events(id) ON DELETE CASCADE,
    consumer_name VARCHAR(100) NOT NULL,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_event_consumer UNIQUE (event_id, consumer_name)
);
```

Workers check `uq_event_consumer` before executing external side-effects (such as sending transactional emails or firing APNs push notifications).
