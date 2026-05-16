package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/google/uuid"

	"github.com/sex-worker-safety/checkin-dao/config"
	"github.com/sex-worker-safety/checkin-dao/models"
)

type EmergencyService struct {
	pool  *pgxpool.Pool
	cfg   *config.Config
	httpC *http.Client
}

func NewEmergencyService(pool *pgxpool.Pool, cfg *config.Config) *EmergencyService {
	return &EmergencyService{
		pool:  pool,
		cfg:   cfg,
		httpC: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *EmergencyService) TriggerAlert(ctx context.Context, userID, alertType string, lat, lng *float64) (*models.EmergencyAlert, error) {
	alert := &models.EmergencyAlert{
		ID:        uuid.New().String(),
		UserID:    userID,
		AlertType: alertType,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if lat != nil && lng != nil {
		alert.LocationLat = lat
		alert.LocationLng = lng
	}

	_, err := s.pool.Exec(ctx, `
		INSERT INTO emergency_alerts (id, user_id, alert_type, location_lat, location_lng, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		alert.ID, alert.UserID, alert.AlertType, alert.LocationLat, alert.LocationLng,
		alert.Status, alert.CreatedAt, alert.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("inserting emergency alert: %w", err)
	}

	log.Warn().
		Str("alert_id", alert.ID).
		Str("user_id", userID).
		Str("alert_type", alertType).
		Msg("emergency alert triggered")

	go s.dispatchNotifications(alert)

	return alert, nil
}

func (s *EmergencyService) ResolveAlert(ctx context.Context, alertID, resolverID string) (*models.EmergencyAlert, error) {
	now := time.Now()
	_, err := s.pool.Exec(ctx, `
		UPDATE emergency_alerts SET status = 'resolved', resolved_by = $1, resolved_at = $2, updated_at = $3
		WHERE id = $4 AND status = 'active'`,
		resolverID, now, now, alertID,
	)
	if err != nil {
		return nil, fmt.Errorf("resolving emergency alert: %w", err)
	}

	var alert models.EmergencyAlert
	err = s.pool.QueryRow(ctx, `
		SELECT id, user_id, alert_type, location_lat, location_lng, status, resolved_by, resolved_at, created_at, updated_at
		FROM emergency_alerts WHERE id = $1`, alertID,
	).Scan(&alert.ID, &alert.UserID, &alert.AlertType, &alert.LocationLat, &alert.LocationLng,
		&alert.Status, &alert.ResolvedBy, &alert.ResolvedAt, &alert.CreatedAt, &alert.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("fetching resolved alert: %w", err)
	}

	log.Info().Str("alert_id", alertID).Str("resolved_by", resolverID).Msg("emergency alert resolved")
	return &alert, nil
}

func (s *EmergencyService) ListAlerts(ctx context.Context, userID, status string) ([]models.EmergencyAlert, error) {
	query := `SELECT id, user_id, alert_type, location_lat, location_lng, status, resolved_by, resolved_at, created_at, updated_at
		FROM emergency_alerts WHERE user_id = $1`
	args := []interface{}{userID}

	if status != "" {
		query += " AND status = $2"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC"

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("listing alerts: %w", err)
	}
	defer rows.Close()

	var alerts []models.EmergencyAlert
	for rows.Next() {
		var a models.EmergencyAlert
		if err := rows.Scan(&a.ID, &a.UserID, &a.AlertType, &a.LocationLat, &a.LocationLng,
			&a.Status, &a.ResolvedBy, &a.ResolvedAt, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning alert: %w", err)
		}
		alerts = append(alerts, a)
	}

	return alerts, rows.Err()
}

func (s *EmergencyService) dispatchNotifications(alert *models.EmergencyAlert) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if s.cfg.TwilioAccountSID != "" && s.cfg.TwilioAuthToken != "" && s.cfg.TwilioFromNumber != "" {
		phone, err := s.getUserPhone(ctx, alert.UserID)
		if err != nil {
			log.Error().Err(err).Str("user_id", alert.UserID).Msg("failed to get user phone for SMS")
		} else if phone != "" {
			s.sendTwilioSMS(ctx, phone, alert)
		}
	}

	if s.cfg.EmergencyWebhook != "" {
		s.sendWebhook(ctx, alert)
	}
}

func (s *EmergencyService) getUserPhone(ctx context.Context, userID string) (string, error) {
	var phone string
	err := s.pool.QueryRow(ctx, `SELECT phone FROM users WHERE id = $1`, userID).Scan(&phone)
	return phone, err
}

func (s *EmergencyService) sendTwilioSMS(ctx context.Context, to string, alert *models.EmergencyAlert) {
	msg := fmt.Sprintf("🚨 EMERGENCY ALERT (%s) - User %s requires immediate assistance.", alert.AlertType, alert.UserID)

	reqBody := fmt.Sprintf("To=%s&From=%s&Body=%s", to, s.cfg.TwilioFromNumber, msg)
	req, err := http.NewRequestWithContext(ctx, "POST",
		fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", s.cfg.TwilioAccountSID),
		bytes.NewBufferString(reqBody))
	if err != nil {
		log.Error().Err(err).Msg("failed to create Twilio request")
		return
	}
	req.SetBasicAuth(s.cfg.TwilioAccountSID, s.cfg.TwilioAuthToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpC.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Twilio SMS request failed")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		log.Error().Int("status", resp.StatusCode).Str("body", string(body)).Msg("Twilio API error")
		return
	}

	log.Info().Str("alert_id", alert.ID).Str("to", to).Msg("Twilio SMS sent")
}

func (s *EmergencyService) sendWebhook(ctx context.Context, alert *models.EmergencyAlert) {
	payload, err := json.Marshal(map[string]interface{}{
		"event":      "emergency_alert",
		"alert_id":   alert.ID,
		"user_id":    alert.UserID,
		"alert_type": alert.AlertType,
		"status":     alert.Status,
		"timestamp":  alert.CreatedAt,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal webhook payload")
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", s.cfg.EmergencyWebhook, bytes.NewBuffer(payload))
	if err != nil {
		log.Error().Err(err).Msg("failed to create webhook request")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpC.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("webhook request failed")
		return
	}
	defer resp.Body.Close()

	log.Info().Str("alert_id", alert.ID).Int("status", resp.StatusCode).Msg("webhook dispatched")
}
