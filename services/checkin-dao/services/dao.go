package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/google/uuid"

	"github.com/sex-worker-safety/checkin-dao/models"
)

type DAOService struct {
	pool *pgxpool.Pool
}

func NewDAOService(pool *pgxpool.Pool) *DAOService {
	return &DAOService{pool: pool}
}

func (s *DAOService) CreateProposal(ctx context.Context, proposerID string, input models.NewProposal) (*models.DAOProposal, error) {
	if input.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if input.VoteDeadline.Before(time.Now()) {
		return nil, fmt.Errorf("vote_deadline must be in the future")
	}
	if input.VoteQuorum < 1 {
		input.VoteQuorum = 5
	}

	proposal := &models.DAOProposal{
		ID:           uuid.New().String(),
		ProposerID:   proposerID,
		Title:        input.Title,
		Description:  input.Description,
		Status:       "active",
		VoteQuorum:   input.VoteQuorum,
		VoteDeadline: input.VoteDeadline,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err := s.pool.Exec(ctx, `
		INSERT INTO dao_proposals (id, proposer_id, title, description, status, vote_quorum, vote_deadline, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		proposal.ID, proposal.ProposerID, proposal.Title, proposal.Description,
		proposal.Status, proposal.VoteQuorum, proposal.VoteDeadline,
		proposal.CreatedAt, proposal.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("creating proposal: %w", err)
	}

	// Register the proposer's automatic yes vote
	_, err = s.pool.Exec(ctx, `
		INSERT INTO dao_votes (id, proposal_id, user_id, vote, weight, voted_at)
		VALUES ($1, $2, $3, 'yes', 1, $4) ON CONFLICT (proposal_id, user_id) DO NOTHING`,
		uuid.New().String(), proposal.ID, proposerID, time.Now(),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to register proposer vote")
	}

	log.Info().Str("proposal_id", proposal.ID).Str("proposer", proposerID).Msg("proposal created")
	return proposal, nil
}

func (s *DAOService) GetProposal(ctx context.Context, proposalID string) (*models.DAOProposal, error) {
	var p models.DAOProposal
	err := s.pool.QueryRow(ctx, `
		SELECT id, proposer_id, title, description, status, vote_quorum, vote_deadline, created_at, updated_at
		FROM dao_proposals WHERE id = $1`, proposalID,
	).Scan(&p.ID, &p.ProposerID, &p.Title, &p.Description, &p.Status,
		&p.VoteQuorum, &p.VoteDeadline, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("proposal not found: %w", err)
	}
	return &p, nil
}

func (s *DAOService) ListProposals(ctx context.Context, status string) ([]models.DAOProposal, error) {
	query := `SELECT id, proposer_id, title, description, status, vote_quorum, vote_deadline, created_at, updated_at
		FROM dao_proposals`
	args := []interface{}{}

	if status != "" {
		query += " WHERE status = $1"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC"

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing proposals: %w", err)
	}
	defer rows.Close()

	var proposals []models.DAOProposal
	for rows.Next() {
		var p models.DAOProposal
		if err := rows.Scan(&p.ID, &p.ProposerID, &p.Title, &p.Description, &p.Status,
			&p.VoteQuorum, &p.VoteDeadline, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning proposal: %w", err)
		}
		proposals = append(proposals, p)
	}

	return proposals, rows.Err()
}

func (s *DAOService) CastVote(ctx context.Context, proposalID, userID, vote string) (*models.DAOVote, error) {
	validVotes := map[string]bool{"yes": true, "no": true, "abstain": true}
	if !validVotes[vote] {
		return nil, fmt.Errorf("invalid vote: must be yes, no, or abstain")
	}

	proposal, err := s.GetProposal(ctx, proposalID)
	if err != nil {
		return nil, err
	}

	if proposal.Status != "active" {
		return nil, fmt.Errorf("proposal is not active")
	}

	if time.Now().After(proposal.VoteDeadline) {
		return nil, fmt.Errorf("voting deadline has passed")
	}

	v := &models.DAOVote{
		ID:         uuid.New().String(),
		ProposalID: proposalID,
		UserID:     userID,
		Vote:       vote,
		Weight:     1,
		VotedAt:    time.Now(),
	}

	_, err = s.pool.Exec(ctx, `
		INSERT INTO dao_votes (id, proposal_id, user_id, vote, weight, voted_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (proposal_id, user_id) DO UPDATE SET vote = $4, weight = $5, voted_at = $6`,
		v.ID, v.ProposalID, v.UserID, v.Vote, v.Weight, v.VotedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("casting vote: %w", err)
	}

	log.Info().Str("proposal_id", proposalID).Str("user_id", userID).Str("vote", vote).Msg("vote cast")
	return v, nil
}

type VoteResults struct {
	Proposal   models.DAOProposal `json:"proposal"`
	YesVotes   int                `json:"yes_votes"`
	NoVotes    int                `json:"no_votes"`
	Abstentions int               `json:"abstentions"`
	TotalVotes int                `json:"total_votes"`
	QuorumMet  bool               `json:"quorum_met"`
	Passed     *bool              `json:"passed,omitempty"`
}

func (s *DAOService) GetResults(ctx context.Context, proposalID string) (*VoteResults, error) {
	proposal, err := s.GetProposal(ctx, proposalID)
	if err != nil {
		return nil, err
	}

	var yesCount, noCount, abstainCount int
	err = s.pool.QueryRow(ctx, `
		SELECT
			COALESCE(SUM(CASE WHEN vote = 'yes' THEN weight ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN vote = 'no' THEN weight ELSE 0 END), 0),
			COALESCE(SUM(CASE WHEN vote = 'abstain' THEN weight ELSE 0 END), 0)
		FROM dao_votes WHERE proposal_id = $1`, proposalID,
	).Scan(&yesCount, &noCount, &abstainCount)
	if err != nil {
		return nil, fmt.Errorf("tallying votes: %w", err)
	}

	totalVotes := yesCount + noCount + abstainCount
	quorumMet := totalVotes >= proposal.VoteQuorum

	var passed *bool
	if proposal.Status == "passed" || proposal.Status == "rejected" {
		result := proposal.Status == "passed"
		passed = &result
	}

	return &VoteResults{
		Proposal:    *proposal,
		YesVotes:    yesCount,
		NoVotes:     noCount,
		Abstentions: abstainCount,
		TotalVotes:  totalVotes,
		QuorumMet:   quorumMet,
		Passed:      passed,
	}, nil
}
