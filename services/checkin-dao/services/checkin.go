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

type CheckInService struct {
	pool            *pgxpool.Pool
	emergencySvc    *EmergencyService
}

func NewCheckInService(pool *pgxpool.Pool, emergencySvc *EmergencyService) *CheckInService {
	return &CheckInService{
		pool:         pool,
		emergencySvc: emergencySvc,
	}
}

func (s *CheckInService) ScheduleCheckIn(ctx context.Context, userID string, input models.NewCheckIn) (*models.CheckIn, error) {
	if input.ScheduledAt.Before(time.Now()) {
		return nil, fmt.Errorf("scheduled_at must be in the future")
	}

	chk := &models.CheckIn{
		ID:          uuid.New().String(),
		UserID:      userID,
		ScheduledAt: input.ScheduledAt,
		Status:      "pending",
		Note:        input.Note,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := s.pool.Exec(ctx, `
		INSERT INTO check_ins (id, user_id, scheduled_at, status, note, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		chk.ID, chk.UserID, chk.ScheduledAt, chk.Status, chk.Note, chk.CreatedAt, chk.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("inserting check-in: %w", err)
	}

	log.Info().Str("checkin_id", chk.ID).Str("user_id", userID).Msg("check-in scheduled")
	return chk, nil
}

func (s *CheckInService) CompleteCheckIn(ctx context.Context, checkInID, userID string) (*models.CheckIn, error) {
	var chk models.CheckIn
	err := s.pool.QueryRow(ctx, `
		SELECT id, user_id, scheduled_at, completed_at, status, note, created_at, updated_at
		FROM check_ins WHERE id = $1 FOR UPDATE`, checkInID).Scan(
		&chk.ID, &chk.UserID, &chk.ScheduledAt, &chk.CompletedAt, &chk.Status,
		&chk.Note, &chk.CreatedAt, &chk.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("check-in not found: %w", err)
	}

	if chk.UserID != userID {
		return nil, fmt.Errorf("check-in does not belong to user")
	}

	if chk.Status != "pending" {
		return nil, fmt.Errorf("check-in is not in pending status")
	}

	now := time.Now()
	_, err = s.pool.Exec(ctx, `
		UPDATE check_ins SET status = 'completed', completed_at = $1, updated_at = $2 WHERE id = $3`,
		now, now, checkInID,
	)
	if err != nil {
		return nil, fmt.Errorf("completing check-in: %w", err)
	}

	chk.Status = "completed"
	chk.CompletedAt = &now
	chk.UpdatedAt = now

	log.Info().Str("checkin_id", checkInID).Msg("check-in completed")
	return &chk, nil
}

func (s *CheckInService) GetCheckIn(ctx context.Context, checkInID, userID string) (*models.CheckIn, error) {
	var chk models.CheckIn
	err := s.pool.QueryRow(ctx, `
		SELECT id, user_id, scheduled_at, completed_at, status, note, created_at, updated_at
		FROM check_ins WHERE id = $1 AND user_id = $2`,
		checkInID, userID,
	).Scan(&chk.ID, &chk.UserID, &chk.ScheduledAt, &chk.CompletedAt, &chk.Status,
		&chk.Note, &chk.CreatedAt, &chk.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("check-in not found: %w", err)
	}
	return &chk, nil
}

func (s *CheckInService) ListCheckIns(ctx context.Context, userID string, status string) ([]models.CheckIn, error) {
	query := `SELECT id, user_id, scheduled_at, completed_at, status, note, created_at, updated_at
		FROM check_ins WHERE user_id = $1`
	args := []interface{}{userID}

	if status != "" {
		query += " AND status = $2"
		args = append(args, status)
	}

	query += " ORDER BY scheduled_at DESC"

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing check-ins: %w", err)
	}
	defer rows.Close()

	var checkins []models.CheckIn
	for rows.Next() {
		var chk models.CheckIn
		if err := rows.Scan(&chk.ID, &chk.UserID, &chk.ScheduledAt, &chk.CompletedAt,
			&chk.Status, &chk.Note, &chk.CreatedAt, &chk.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning check-in: %w", err)
		}
		checkins = append(checkins, chk)
	}

	return checkins, rows.Err()
}

func (s *CheckInService) CheckMissedCheckIns(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("check-in monitor shutting down")
			return
		case <-ticker.C:
			s.processMissedCheckIns(ctx)
		}
	}
}

func (s *CheckInService) processMissedCheckIns(ctx context.Context) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, user_id, scheduled_at FROM check_ins
		WHERE status = 'pending' AND scheduled_at < NOW()
		FOR UPDATE SKIP LOCKED`)
	if err != nil {
		log.Error().Err(err).Msg("failed to query missed check-ins")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, userID string
		var scheduledAt time.Time
		if err := rows.Scan(&id, &userID, &scheduledAt); err != nil {
			log.Error().Err(err).Msg("failed to scan missed check-in")
			continue
		}

		now := time.Now()
		_, err := s.pool.Exec(ctx,
			`UPDATE check_ins SET status = 'missed', updated_at = $1 WHERE id = $2`,
			now, id)
		if err != nil {
			log.Error().Err(err).Str("checkin_id", id).Msg("failed to mark check-in as missed")
			continue
		}

		log.Warn().Str("checkin_id", id).Str("user_id", userID).Msg("check-in missed, triggering alert")

		if s.emergencySvc != nil {
			if _, err := s.emergencySvc.TriggerAlert(ctx, userID, "safety", nil, nil); err != nil {
				log.Error().Err(err).Msg("failed to trigger emergency alert for missed check-in")
			}
		}
	}
}
