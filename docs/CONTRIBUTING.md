# Contributing

## Getting Started

1. Fork the repository and create a feature branch from `main`.
2. Copy `.env.example` to `.env` and fill in required values.
3. Start infrastructure: `docker compose up db qdrant`
4. Run migrations: `psql $DATABASE_URL -f migrations/001_initial_schema.sql`

## Development Workflow

```bash
# Rust services
cd services/gateway && cargo check     # API Gateway
cd services/attestation && cargo check # Attestation Service

# Go service
cd services/checkin-dao && go build ./...

# Julia service
cd services/analytics && julia run.jl

# Frontend
cd frontend && npm run dev
```

## Code Style

- **Rust**: Run `cargo fmt` and `cargo clippy` before committing
- **Go**: Run `go fmt ./...` and `go vet ./...`
- **Julia**: Follow [Julia style guide](https://docs.julialang.org/en/v1/manual/style-guide/)
- **TypeScript**: ESLint + Prettier via `npm run lint`

## Pull Requests

1. Keep changes focused — one feature/fix per PR.
2. Update tests when adding or modifying functionality.
3. Update the README if changing public API or architecture.
4. Ensure all CI checks pass (lint, test, build).
5. Request review from a maintainer.
