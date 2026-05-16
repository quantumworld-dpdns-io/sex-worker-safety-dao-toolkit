-- Sex Worker Safety DAO Toolkit - PostgreSQL Schema

BEGIN;

-- ── Users ──
CREATE TABLE IF NOT EXISTS users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_address  TEXT UNIQUE NOT NULL,
    role            TEXT NOT NULL DEFAULT 'member' CHECK (role IN ('admin', 'moderator', 'member')),
    display_name    TEXT,
    encrypted_contact TEXT,
    is_active       BOOLEAN NOT NULL DEFAULT true,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ── Attestations ──
CREATE TABLE IF NOT EXISTS attestations (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id_hash    TEXT NOT NULL,
    encrypted_payload TEXT NOT NULL,
    proof_type      TEXT NOT NULL DEFAULT 'noir' CHECK (proof_type IN ('noir', 'risc0')),
    proof_data      TEXT NOT NULL,
    proof_hash      TEXT UNIQUE NOT NULL,
    circuit_type    TEXT NOT NULL,
    metadata_json   JSONB DEFAULT '{}',
    is_verified     BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_attestations_user_hash ON attestations(user_id_hash);

-- ── Check-ins ──
CREATE TABLE IF NOT EXISTS check_ins (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    scheduled_at    TIMESTAMPTZ NOT NULL,
    window_minutes  INT NOT NULL DEFAULT 60,
    status          TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'missed', 'emergency')),
    completed_at    TIMESTAMPTZ,
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_check_ins_user ON check_ins(user_id);
CREATE INDEX idx_check_ins_status ON check_ins(status);
CREATE INDEX idx_check_ins_scheduled ON check_ins(scheduled_at);

-- ── Emergency Alerts ──
CREATE TABLE IF NOT EXISTS emergency_alerts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id),
    trigger_type    TEXT NOT NULL CHECK (trigger_type IN ('missed_checkin', 'manual', 'sos_gesture')),
    location_data   JSONB DEFAULT '{}',
    contact_method  TEXT NOT NULL DEFAULT 'sms',
    resolved_at     TIMESTAMPTZ,
    resolution_note TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_emergency_user ON emergency_alerts(user_id);

-- ── Bad Client Reports ──
CREATE TABLE IF NOT EXISTS bad_client_reports (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attestation_id    UUID REFERENCES attestations(id),
    encrypted_details TEXT NOT NULL,
    similarity_hash   TEXT,
    location_region   TEXT,
    report_category   TEXT NOT NULL DEFAULT 'unsafe' CHECK (report_category IN ('unsafe', 'theft', 'assault', 'harassment', 'other')),
    confidence_score  REAL DEFAULT 0.0,
    status            TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'verified', 'dismissed')),
    moderator_id      UUID REFERENCES users(id),
    moderator_notes   TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_bcr_status ON bad_client_reports(status);
CREATE INDEX idx_bcr_category ON bad_client_reports(report_category);
CREATE INDEX idx_bcr_similarity ON bad_client_reports(similarity_hash);

-- ── DAO Proposals ──
CREATE TABLE IF NOT EXISTS dao_proposals (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title           TEXT NOT NULL,
    description     TEXT NOT NULL,
    proposal_type   TEXT NOT NULL DEFAULT 'general' CHECK (proposal_type IN ('general', 'membership', 'budget', 'parameter')),
    status          TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('draft', 'active', 'passed', 'rejected', 'executed')),
    voting_type     TEXT NOT NULL DEFAULT 'simple_majority' CHECK (voting_type IN ('simple_majority', 'super_majority', 'quadratic')),
    starts_at       TIMESTAMPTZ NOT NULL,
    ends_at         TIMESTAMPTZ NOT NULL,
    created_by      UUID NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_dao_proposals_status ON dao_proposals(status);

-- ── DAO Votes ──
CREATE TABLE IF NOT EXISTS dao_votes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    proposal_id     UUID NOT NULL REFERENCES dao_proposals(id),
    voter_id        UUID NOT NULL REFERENCES users(id),
    vote            TEXT NOT NULL CHECK (vote IN ('yes', 'no', 'abstain')),
    voting_weight   REAL NOT NULL DEFAULT 1.0,
    tx_hash         TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(proposal_id, voter_id)
);
CREATE INDEX idx_dao_votes_proposal ON dao_votes(proposal_id);

-- ── Audit Log ──
CREATE TABLE IF NOT EXISTS audit_log (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id        UUID REFERENCES users(id),
    action          TEXT NOT NULL,
    resource_type   TEXT NOT NULL,
    resource_id     TEXT,
    details_json    JSONB DEFAULT '{}',
    ip_address      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_audit_actor ON audit_log(actor_id);
CREATE INDEX idx_audit_action ON audit_log(action);
CREATE INDEX idx_audit_created ON audit_log(created_at);

COMMIT;
