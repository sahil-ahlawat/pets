# Pet Digital Identity Platform

> A premium, scalable pet identity and lifetime-record platform where every companion gets a permanent digital profile accessible through a smart **QR/NFC collar tag**.

---

## 1. System Overview

```text
Create account
    ↓
Create companion profile
    ↓
Add photo + identity + health information
    ↓
Generate smart QR tag
    ↓
Print / attach to collar
    ↓
Someone scans QR
    ↓
Safe public companion profile (No app required)
    ↓
Contact owner safely / report found companion
```

---

## 2. Repository Layout

```text
pets/
├── backend/              # Go REST API, background workers & scheduler
│   ├── cmd/              # api, worker, scheduler main entrypoints
│   ├── internal/         # config, domain, handler, repository, service, outbox, worker
│   ├── migrations/       # PostgreSQL up and down migrations
│   ├── tests/            # Automated unit and integration tests
│   ├── Dockerfile        # Multi-stage production container build
│   ├── go.mod            # Go dependencies (Chi, JWT, PQ, go-qrcode)
│   └── openapi.yaml      # OpenAPI 3.0 API Contract
│
├── frontend/             # Premium Vanilla JS, HTML5 & CSS3 Web Client
│   ├── css/              # Luxury design system and responsive layout
│   ├── js/               # API client, auth session manager, UI utilities
│   ├── pages/            # Dashboard, Pet profile, QR manager, Public scan view
│   ├── nginx.conf        # Production reverse proxy and deep-link router
│   └── Dockerfile        # Lightweight Nginx container
│
├── app/                  # Production-Quality Native Flutter Mobile Application
│   ├── lib/              # Providers, models, screens, theme, networking
│   ├── android/          # Deep link intent filters, camera and location permissions
│   ├── ios/              # Universal link associated domains and privacy keys
│   └── pubspec.yaml      # Flutter dependencies (http, provider, qr_flutter)
│
├── database/
│   └── schema/           # Normalized PostgreSQL schema, seed data, and partitioning
│
├── infra/
│   ├── docker/           # Production container manifests for all services
│   └── helm/             # Kubernetes-ready Helm chart with multi-environment values
│
├── docs/                 # Architecture, event-driven outbox, and API specifications
├── docker-compose.yml    # One-click local development environment
└── README.md
```

---

## 3. Technology Stack

| Domain | Technology |
| :--- | :--- |
| **Backend** | Go 1.24 (Chi Router, lib/pq, go-qrcode, JWT, bcrypt) |
| **API** | REST + OpenAPI 3.0 Specification |
| **Web** | HTML5 + CSS3 + Vanilla JavaScript (Modular ES6) |
| **Mobile** | Flutter (iOS & Android with native deep-linking) |
| **Database** | PostgreSQL 16 + PostGIS extension |
| **Events** | Transactional Outbox Pattern |
| **Workers** | Asynchronous Go Workers with retries and DLQ |
| **Cache** | Redis 7 |
| **Object Storage** | S3 / Cloudflare R2 presigned uploads |
| **Containers** | Docker & Docker Compose |
| **Orchestration** | Kubernetes & Helm (`infra/helm/pet-platform/`) |

---

## 4. Quick Start (Local Development)

### Prerequisites
* [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
* [Go 1.24+](https://golang.org/) (optional for local binary builds)

### 1. Start Complete Stack
```bash
docker compose up --build
```

This starts:
* **PostgreSQL + PostGIS**: `localhost:5432` (auto-applies schema and seeds)
* **Redis**: `localhost:6379`
* **Go REST API**: `http://localhost:8080`
* **Asynchronous Workers**: Processes scan alerts, lost broadcasts, and reminders
* **Scheduler**: Checks due booster vaccines and preventive reminders
* **Web Frontend (Nginx)**: `http://localhost:3000`

### 2. Verify Health
```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
```

---

## 5. Running the Backend Tests

```bash
cd backend
go test -v ./tests/...
```

---

## 6. Public QR Scan Experience

When a collar tag is scanned by any smartphone camera:
1. Resolves token via `/p/:token` or `/pages/public-pet.html?token=:token`.
2. Dedicated **Public DTO** ensures sensitive home addresses and medical records remain private.
3. If pet is marked **LOST**, a prominent emergency banner with last seen address is rendered.
4. Finder can submit a **Found Report** with optional GPS coordinates or message the owner securely.

---

## 7. Kubernetes Deployment with Helm

Production deployment is fully Kubernetes-ready via the included Helm chart:

```bash
# Dry-run template validation
helm template pet-platform infra/helm/pet-platform/ -f infra/helm/pet-platform/values-production.yaml

# Deploy to cluster
helm upgrade --install pet-platform infra/helm/pet-platform/ \
  --namespace pet-platform --create-namespace \
  -f infra/helm/pet-platform/values-production.yaml
```

---

## 8. License

MIT License. Designed and engineered for companion animal safety.