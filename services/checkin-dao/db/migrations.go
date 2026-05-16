package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

const schemaSQL = `
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_address TEXT NOT NULL UNIQUE,
    role TEXT NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'moderator', 'admin')),
    display_name TEXT NOT NULL DEFAULT '',
    phone TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS check_ins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    scheduled_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'missed', 'cancelled')),
    note TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS emergency_alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    alert_type TEXT NOT NULL DEFAULT 'general' CHECK (alert_type IN ('general', 'safety', 'medical', 'legal')),
    location_lat DOUBLE PRECISION,
    location_lng DOUBLE PRECISION,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'resolved', 'false_alarm')),
    resolved_by UUID REFERENCES users(id),
    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS bad_client_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reporter_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    client_identifier TEXT NOT NULL,
    category TEXT NOT NULL DEFAULT 'other' CHECK (category IN ('violent', 'harassment', 'theft', 'stalking', 'non_payment', 'other')),
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'verified', 'dismissed')),
    moderator_id UUID REFERENCES users(id),
    moderated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS dao_proposals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    proposer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'passed', 'rejected', 'cancelled')),
    vote_quorum INT NOT NULL DEFAULT 5,
    vote_deadline TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS dao_votes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    proposal_id UUID NOT NULL REFERENCES dao_proposals(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    vote TEXT NOT NULL CHECK (vote IN ('yes', 'no', 'abstain')),
    weight INT NOT NULL DEFAULT 1,
    voted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(proposal_id, user_id)
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    action TEXT NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id TEXT,
    details JSONB NOT NULL DEFAULT '{}',
    ip_address TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_check_ins_user_id ON check_ins(user_id);
CREATE INDEX IF NOT EXISTS idx_check_ins_status ON check_ins(status);
CREATE INDEX IF NOT EXISTS idx_check_ins_scheduled_at ON check_ins(scheduled_at);
CREATE INDEX IF NOT EXISTS idx_emergency_alerts_user_id ON emergency_alerts(user_id);
CREATE INDEX IF NOT EXISTS idx_emergency_alerts_status ON emergency_alerts(status);
CREATE INDEX IF NOT EXISTS idx_bad_client_reports_reporter_id ON bad_client_reports(reporter_id);
CREATE INDEX IF NOT EXISTS idx_bad_client_reports_category ON bad_client_reports(category);
CREATE INDEX IF NOT EXISTS idx_bad_client_reports_status ON bad_client_reports(status);
CREATE INDEX IF NOT EXISTS idx_dao_proposals_status ON dao_proposals(status);
CREATE INDEX IF NOT EXISTS idx_dao_votes_proposal_id ON dao_votes(proposal_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
`

func RunMigrations(pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err := pool.Exec(ctx, schemaSQL); err != nil {
		return fmt.Errorf("running migrations: %w", err)
	}

	log.Info().Msg("database migrations completed")
	return nil
}
