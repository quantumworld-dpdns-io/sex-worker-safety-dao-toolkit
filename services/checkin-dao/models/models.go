package models

import (
	"time"
)

type User struct {
	ID            string    `json:"id" db:"id"`
	WalletAddress string    `json:"wallet_address" db:"wallet_address"`
	Role          string    `json:"role" db:"role"`
	DisplayName   string    `json:"display_name" db:"display_name"`
	Phone         string    `json:"phone" db:"phone"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type CheckIn struct {
	ID          string     `json:"id" db:"id"`
	UserID      string     `json:"user_id" db:"user_id"`
	ScheduledAt time.Time  `json:"scheduled_at" db:"scheduled_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	Status      string     `json:"status" db:"status"`
	Note        string     `json:"note" db:"note"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type EmergencyAlert struct {
	ID          string     `json:"id" db:"id"`
	UserID      string     `json:"user_id" db:"user_id"`
	AlertType   string     `json:"alert_type" db:"alert_type"`
	LocationLat *float64   `json:"location_lat,omitempty" db:"location_lat"`
	LocationLng *float64   `json:"location_lng,omitempty" db:"location_lng"`
	Status      string     `json:"status" db:"status"`
	ResolvedBy  *string    `json:"resolved_by,omitempty" db:"resolved_by"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty" db:"resolved_at"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type BadClientReport struct {
	ID               string     `json:"id" db:"id"`
	ReporterID       string     `json:"reporter_id" db:"reporter_id"`
	ClientIdentifier string     `json:"client_identifier" db:"client_identifier"`
	Category         string     `json:"category" db:"category"`
	Description      string     `json:"description" db:"description"`
	Status           string     `json:"status" db:"status"`
	ModeratorID      *string    `json:"moderator_id,omitempty" db:"moderator_id"`
	ModeratedAt      *time.Time `json:"moderated_at,omitempty" db:"moderated_at"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

type DAOProposal struct {
	ID           string    `json:"id" db:"id"`
	ProposerID   string    `json:"proposer_id" db:"proposer_id"`
	Title        string    `json:"title" db:"title"`
	Description  string    `json:"description" db:"description"`
	Status       string    `json:"status" db:"status"`
	VoteQuorum   int       `json:"vote_quorum" db:"vote_quorum"`
	VoteDeadline time.Time `json:"vote_deadline" db:"vote_deadline"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type DAOVote struct {
	ID         string    `json:"id" db:"id"`
	ProposalID string    `json:"proposal_id" db:"proposal_id"`
	UserID     string    `json:"user_id" db:"user_id"`
	Vote       string    `json:"vote" db:"vote"`
	Weight     int       `json:"weight" db:"weight"`
	VotedAt    time.Time `json:"voted_at" db:"voted_at"`
}

type AuditLog struct {
	ID           string    `json:"id" db:"id"`
	UserID       *string   `json:"user_id,omitempty" db:"user_id"`
	Action       string    `json:"action" db:"action"`
	ResourceType string    `json:"resource_type" db:"resource_type"`
	ResourceID   *string   `json:"resource_id,omitempty" db:"resource_id"`
	Details      string    `json:"details" db:"details"`
	IPAddress    string    `json:"ip_address" db:"ip_address"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type NewCheckIn struct {
	ScheduledAt time.Time `json:"scheduled_at"`
	Note        string    `json:"note"`
}

type NewEmergencyAlert struct {
	AlertType   string   `json:"alert_type"`
	LocationLat *float64 `json:"location_lat,omitempty"`
	LocationLng *float64 `json:"location_lng,omitempty"`
}

type NewReport struct {
	ClientIdentifier string `json:"client_identifier"`
	Category         string `json:"category"`
	Description      string `json:"description"`
}

type NewProposal struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	VoteQuorum   int       `json:"vote_quorum"`
	VoteDeadline time.Time `json:"vote_deadline"`
}

type NewVote struct {
	Vote string `json:"vote"`
}
