CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS attestations (
    id UUID PRIMARY KEY,
    circuit_type VARCHAR(255) NOT NULL,
    proof_data TEXT NOT NULL DEFAULT '',
    proof_hash VARCHAR(64) NOT NULL DEFAULT '',
    user_id_hash VARCHAR(64) NOT NULL,
    payload BYTEA NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_attestations_user_id_hash ON attestations(user_id_hash);
CREATE INDEX IF NOT EXISTS idx_attestations_circuit_type ON attestations(circuit_type);
