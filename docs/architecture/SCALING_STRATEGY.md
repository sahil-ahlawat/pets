# Scaling Strategy — From MVP to Millions of Tag Scans

To avoid premature complexity while ensuring a clear path to millions of concurrent scans, the platform follows four distinct scale stages:

### Stage 1: MVP Foundation
* Go REST API (1-2 container instances)
* PostgreSQL (Managed primary instance)
* Object Storage (S3 / Cloudflare R2)
* Embedded Go Outbox publisher & background worker
* Static frontend behind Cloudflare CDN

### Stage 2: Growth & Resilience
* Multiple horizontal Go API pods with HPA
* Redis cluster for session cache and query memoization
* Read replica for PostgreSQL reporting queries
* Decoupled worker pods consuming asynchronous job tables

### Stage 3: High Throughput
* Dedicated Kafka / Redpanda event broker
* Specialized workers: Push notification worker, Geo-analytics worker, Document scanner
* PostGIS spatial indexing on scan locations and veterinary practices

### Stage 4: Hyper-Scale Global Architecture
* Range partitioning on high-volume tables (`qr_scans`, `audit_logs`, `outbox_events`)
* Multi-region Cloudflare edge caching for public pet profile DTOs (`/p/:token`)
* Read replicas distributed across geographic availability zones
