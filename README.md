# Sex Worker Safety DAO Toolkit

> Cooperative shared bad-client reporting and emergency check-ins with encrypted incident attestations via zero-knowledge proofs.

A production-ready, multi-language microservices platform for anonymous safety reporting, emergency check-ins, and DAO governance — designed for sex worker communities.

## Architecture

```
                    ┌──────────────────────┐
                    │   Next.js App         │  ← Vercel
                    │   (Frontend)          │
                    └──────────┬───────────┘
                               │ HTTPS
                    ┌──────────▼───────────┐
                    │  Rust API Gateway     │  Image 1: Rust (pingora)
                    │  Reverse Proxy + Auth │  JWT, rate-limit, TLS, encrypt/decrypt
                    └──────────┬───────────┘
                               │
          ┌────────────────────┼────────────────────┐
          ▼                    ▼                     ▼
┌──────────────────┐  ┌───────────────┐  ┌──────────────────────┐
│  Rust Attestation│  │  Go Check-in  │  │  Julia Analytics     │
│  Service         │  │  & Emergency  │  │  + Qdrant Similarity │
│  (Noir + RISC-Z) │  │  + DAO/Reg    │  │  Engine              │
│  Image 2         │  │  Image 3      │  │  Image 4             │
└────────┬─────────┘  └──────┬────────┘  └──────────┬───────────┘
         │                   │                       │
         └───────────────────┼───────────────────────┘
                             ▼
                    ┌────────────────┐
                    │  PostgreSQL    │  ← AlwaysData
                    └────────────────┘

                    ┌────────────────┐
                    │    Qdrant      │  Vector DB (similarity matching)
                    └────────────────┘
```

| Component | Language | Role | Docker Image | Deploy Target |
|---|---|---|---|---|
| **API Gateway** | Rust (pingora) | Reverse proxy, JWT auth, rate limiting, TLS, request encryption | Image 1 | Choreo |
| **Attestation Service** | Rust (axum) | ZK proof generation/verification (Noir + RISC Zero), encrypted attestations | Image 2 | Choreo |
| **Check-in & Emergency + DAO** | Go (chi) | Check-in scheduling, emergency alerts, bad-client registry, DAO governance | Image 3 | Choreo |
| **Analytics Engine** | Julia (Genie) | Qdrant vector similarity, report dedup, trend/DAO analytics | Image 4 | Choreo |
| **Frontend** | Next.js 14 (TypeScript) | User interface, attestation forms, dashboard, emergency button | Vercel | Vercel |
| **Database** | PostgreSQL | Primary data store | N/A (managed) | AlwaysData |
| **Vector DB** | Qdrant | Similarity search for bad-client reports | Sidecar | Choreo / managed |

## Quick Start

### Prerequisites

- Rust 1.78+
- Go 1.22+
- Julia 1.10+
- Node.js 20+
- Docker & Docker Compose
- PostgreSQL 16

### Local Development

```bash
# 1. Start infrastructure (PostgreSQL + Qdrant)
docker compose up db qdrant

# 2. Set up database
psql -h localhost -U postgres -d safety-dao -f migrations/001_initial_schema.sql

# 3. Run services (in separate terminals)

# API Gateway (Rust)
cd services/gateway && cargo run

# Attestation Service (Rust)
cd services/attestation && cargo run

# Check-in & DAO (Go)
cd services/checkin-dao && go run .

# Analytics (Julia)
cd services/analytics && julia run.jl

# Frontend (Next.js)
cd frontend && npm install && npm run dev
```

### Docker Compose (all services)

```bash
# Copy environment template
cp .env.example .env
# Edit .env with your configuration

# Build and start everything
docker compose up --build
```

### Environment Variables

See [`.env.example`](.env.example) for all required and optional configuration.

## Services

### 1. API Gateway (`services/gateway/`)

Full-featured reverse proxy built on Cloudflare's [pingora](https://github.com/cloudflare/pingora) framework.

| Feature | Implementation |
|---|---|
| Reverse proxying | Layer 7 path-based routing to upstream services |
| JWT auth | HS256 token validation, injects user claims as headers |
| Rate limiting | Per-IP token bucket (governor) |
| TLS termination | Automatic ACME/Let's Encrypt via rustls |
| Request encryption | AES-256-GCM encrypt/decrypt for sensitive payloads |
| Structured logging | JSON tracing with request IDs |
| Health checking | Passive + active probes to backends |

**Endpoints:**

| Path | Auth | Proxy To |
|---|---|---|
| `GET /health` | No | Gateway itself |
| `POST /api/v1/auth/*` | No | Check-in service |
| `POST /api/v1/attestations/*` | Yes | Attestation service |
| `POST /api/v1/checkins/*` | Yes | Check-in service |
| `POST /api/v1/emergency/*` | Yes | Check-in service |
| `GET/POST /api/v1/dao/*` | Yes | Check-in service |
| `GET/POST /api/v1/registry/*` | Yes | Check-in service |
| `GET/POST /api/v1/analytics/*` | Yes | Analytics service |

### 2. Attestation Service (`services/attestation/`)

Zero-knowledge proof attestation engine. Encrypts incident reports and generates verifiable proofs without revealing reporter identity.

| Feature | Status |
|---|---|
| AES-256-GCM payload encryption | ✅ Implemented |
| ZK proof generation (mock) | ✅ Trait interface ready |
| Noir circuit integration | 🔜 Interface defined, wiring pending |
| RISC Zero zkVM integration | 🔜 Interface defined, wiring pending |
| Proof verification | ✅ Mock verification complete |
| PostgreSQL persistence | ✅ sqlx async |

**Proof Engine Architecture:**

```rust
trait ProofEngine {
    fn generate_proof(&self, circuit_type: &str, private: &[u8], public: &[u8])
        -> Result<ProofOutput>;
    fn verify_proof(&self, proof_data: &str, public: &[u8])
        -> Result<bool>;
}
```

Implementations: `NoirProver` (Noir/Barretenberg), `Risc0Prover` (RISC Zero zkVM).

### 3. Check-in & Emergency + DAO (`services/checkin-dao/`)

Go service handling all CRUD-heavy operations.

**Check-in Module:**
- Schedule periodic safety check-ins with configurable windows
- Auto-detect missed check-ins (background goroutine)
- Emergency alert dispatch via SMS (Twilio) and webhooks

**Bad Client Registry:**
- Submit encrypted reports with category classification
- Qdrant-powered similarity deduplication
- Moderator verification workflow
- Confidence scoring based on cross-referencing

**DAO Governance:**
- Proposal lifecycle (draft → active → passed/rejected → executed)
- Voting types: simple majority, super majority, quadratic
- Vote tally with quorum enforcement
- Audit log for all governance actions

### 4. Analytics Engine (`services/analytics/`)

Julia-based analytics and similarity matching service.

| Endpoint | Purpose |
|---|---|
| `POST /api/v1/analytics/dedup` | Check report for Qdrant duplicates |
| `GET /api/v1/analytics/trends` | Aggregated report statistics |
| `GET /api/v1/analytics/dao` | DAO governance metrics |
| `POST /api/v1/analytics/embed` | Text → embedding vector |

**Similarity Flow:**
1. New report text → deterministic 128-dim embedding
2. Search Qdrant for cosine similarity > 0.85
3. Flag potential duplicates for moderator review
4. Store new embedding if no match found

### 5. Frontend (`frontend/`)

Next.js 14 App Router with TypeScript and Tailwind CSS.

**Pages:**

| Route | Feature |
|---|---|
| `/` | Landing page with hero and features |
| `/login` | Wallet-based authentication (simulated) |
| `/dashboard` | Stats cards, activity feed, quick actions |
| `/report` | Anonymous attestation submission form |
| `/registry` | Searchable bad-client report list |
| `/checkin` | Safety check-in scheduler and history |
| `/emergency` | SOS panic button with alert history |
| `/dao` | DAO proposal list and creation |
| `/dao/proposals/[id]` | Proposal detail with voting |

## Deployment

### Choreo (Backend)

4 Docker images are deployed to [console.choreo.dev](https://console.choreo.dev):

```bash
# Build images
docker compose -f docker-compose.choreo.yml build

# Push to Choreo registry
# (follow Choreo CLI instructions)
```

See [`docker-compose.choreo.yml`](docker-compose.choreo.yml) for the Choreo-specific compose configuration.

### Vercel (Frontend)

```bash
cd frontend
npm install
npm run build

# Connect your GitHub repo to Vercel:
# - Framework: Next.js
# - Root directory: frontend/
# - Environment variables from .env.example
```

### AlwaysData (PostgreSQL)

```bash
# Create database on AlwaysData control panel
# Run migration:
psql $DATABASE_URL -f migrations/001_initial_schema.sql
```

## Database Schema

7 core tables: `users`, `attestations`, `check_ins`, `emergency_alerts`, `bad_client_reports`, `dao_proposals`, `dao_votes` + `audit_log`.

Full schema at [`migrations/001_initial_schema.sql`](migrations/001_initial_schema.sql).

## CI/CD

GitHub Actions workflow (`.github/workflows/ci.yml`):
- **lint-rust** — `cargo fmt` + `cargo clippy` for gateway + attestation
- **test-rust** — `cargo test` for gateway + attestation
- **lint-go** — `go vet` for checkin-dao
- **build-go** — `go build` for checkin-dao
- **lint-frontend** — `npm run lint`
- **build-frontend** — `npm run build`
- **docker-build** — Build all Docker images with BuildKit cache

## Security Model

- **Anonymity-by-design**: ZK proofs (Noir/RISC Zero) prevent identity leakage
- **Encryption at rest**: AES-256-GCM for stored attestation payloads
- **Encryption in transit**: TLS termination at gateway (rustls + ACME)
- **Authentication**: JWT tokens with role-based access control
- **Rate limiting**: Per-IP token bucket at gateway
- **Audit trail**: Immutable logs for all registry and governance changes
- **PQC readiness**: CIRCL/liboqs integration planned (Phase 2)

## Project Structure

```
.
├── services/
│   ├── gateway/          # Rust API Gateway (pingora)
│   ├── attestation/      # Rust ZK Attestation Service (axum)
│   ├── checkin-dao/      # Go Check-in & DAO Service (chi)
│   └── analytics/        # Julia Analytics Engine (Genie)
├── frontend/             # Next.js 14 (TypeScript, Tailwind)
├── circuits/
│   ├── noir/             # Noir ZK circuits
│   └── risc-zero/        # RISC Zero guest programs
├── migrations/           # PostgreSQL schema migrations
├── docker/               # Dockerfiles per service
├── .github/workflows/    # CI/CD pipeline
├── docker-compose.yml    # Local development
├── docker-compose.choreo.yml  # Choreo deployment
├── Justfile              # Common task shortcuts
└── .env.example          # Environment variable template
```

## Developing

### Adding a New Noir Circuit

1. Create circuit in `circuits/noir/<name>/`
2. Implement `ProofEngine` for the new circuit type
3. Register in attestation service handler

### Adding a New DAO Proposal Type

1. Add enum variant in Go `models/models.go`
2. Add validation logic in `services/dao.go`
3. Add handler in `handlers/dao.go`
4. Update frontend `types/index.ts` + proposal form

## Contributing

Please read [CONTRIBUTING.md](docs/CONTRIBUTING.md) before opening a pull request.

## License

[MIT](LICENSE)
