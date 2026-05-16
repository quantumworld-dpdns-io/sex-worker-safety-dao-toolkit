package services

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/google/uuid"

	"github.com/sex-worker-safety/checkin-dao/models"
)

type RegistryService struct {
	pool         *pgxpool.Pool
	qdrantClient *QdrantClient
}

func NewRegistryService(pool *pgxpool.Pool, qdrant *QdrantClient) *RegistryService {
	return &RegistryService{
		pool:         pool,
		qdrantClient: qdrant,
	}
}

func (s *RegistryService) SubmitReport(ctx context.Context, reporterID string, input models.NewReport) (*models.BadClientReport, error) {
	if input.ClientIdentifier == "" {
		return nil, fmt.Errorf("client_identifier is required")
	}

	report := &models.BadClientReport{
		ID:               uuid.New().String(),
		ReporterID:       reporterID,
		ClientIdentifier: input.ClientIdentifier,
		Category:         input.Category,
		Description:      input.Description,
		Status:           "pending",
		CreatedAt:        now(),
		UpdatedAt:        now(),
	}

	_, err := s.pool.Exec(ctx, `
		INSERT INTO bad_client_reports (id, reporter_id, client_identifier, category, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		report.ID, report.ReporterID, report.ClientIdentifier, report.Category,
		report.Description, report.Status, report.CreatedAt, report.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("inserting report: %w", err)
	}

	log.Info().
		Str("report_id", report.ID).
		Str("reporter_id", reporterID).
		Str("client", input.ClientIdentifier).
		Msg("bad client report submitted")

	go func() {
		if s.qdrantClient != nil {
			c := context.Background()
			vector := hashToVector(input.ClientIdentifier + input.Description)
			payload := map[string]interface{}{
				"report_id":         report.ID,
				"client_identifier": input.ClientIdentifier,
				"category":          input.Category,
				"status":            report.Status,
				"created_at":        report.CreatedAt.String(),
			}
			if err := s.qdrantClient.StoreEmbedding("bad_client_reports", report.ID, vector, payload); err != nil {
				log.Error().Err(err).Msg("failed to store embedding for report")
			}
		}
	}()

	return report, nil
}

func (s *RegistryService) GetReport(ctx context.Context, reportID string) (*models.BadClientReport, error) {
	var r models.BadClientReport
	err := s.pool.QueryRow(ctx, `
		SELECT id, reporter_id, client_identifier, category, description, status, moderator_id, moderated_at, created_at, updated_at
		FROM bad_client_reports WHERE id = $1`, reportID,
	).Scan(&r.ID, &r.ReporterID, &r.ClientIdentifier, &r.Category, &r.Description,
		&r.Status, &r.ModeratorID, &r.ModeratedAt, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("report not found: %w", err)
	}
	return &r, nil
}

func (s *RegistryService) SearchReports(ctx context.Context, query, category, status string) ([]models.BadClientReport, error) {
	sql := `SELECT id, reporter_id, client_identifier, category, description, status, moderator_id, moderated_at, created_at, updated_at
		FROM bad_client_reports WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if query != "" {
		sql += fmt.Sprintf(" AND (client_identifier ILIKE $%d OR description ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+query+"%")
		argIdx++
	}
	if category != "" {
		sql += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, category)
		argIdx++
	}
	if status != "" {
		sql += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	sql += " ORDER BY created_at DESC"

	rows, err := s.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("searching reports: %w", err)
	}
	defer rows.Close()

	var reports []models.BadClientReport
	for rows.Next() {
		var r models.BadClientReport
		if err := rows.Scan(&r.ID, &r.ReporterID, &r.ClientIdentifier, &r.Category, &r.Description,
			&r.Status, &r.ModeratorID, &r.ModeratedAt, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning report: %w", err)
		}
		reports = append(reports, r)
	}

	return reports, rows.Err()
}

func (s *RegistryService) UpdateReportStatus(ctx context.Context, reportID, moderatorID, status string) (*models.BadClientReport, error) {
	validStatuses := map[string]bool{"verified": true, "dismissed": true}
	if !validStatuses[status] {
		return nil, fmt.Errorf("invalid status: %s", status)
	}

	now := now()
	_, err := s.pool.Exec(ctx, `
		UPDATE bad_client_reports SET status = $1, moderator_id = $2, moderated_at = $3, updated_at = $4
		WHERE id = $5`,
		status, moderatorID, now, now, reportID,
	)
	if err != nil {
		return nil, fmt.Errorf("updating report status: %w", err)
	}

	return s.GetReport(ctx, reportID)
}

func (s *RegistryService) FindDuplicates(ctx context.Context, clientIdentifier, description string) ([]MockPoint, error) {
	if s.qdrantClient == nil {
		return nil, nil
	}

	vector := hashToVector(clientIdentifier + description)
	return s.qdrantClient.SearchSimilar("bad_client_reports", vector, 5)
}

func hashToVector(input string) []float32 {
	// Simple hash-based vector generation for demo purposes
	v := make([]float32, 128)
	for i, b := range []byte(input) {
		v[i%128] += float32(b) / 255.0
	}
	// Normalize
	var norm float64
	for _, val := range v {
		norm += float64(val * val)
	}
	norm = sqrt(norm)
	if norm > 0 {
		for i := range v {
			v[i] = float32(float64(v[i]) / norm)
		}
	}
	return v
}

func now() time.Time {
	return time.Now()
}
