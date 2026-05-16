package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/sex-worker-safety/checkin-dao/models"
	"github.com/sex-worker-safety/checkin-dao/services"
)

type EmergencyHandler struct {
	svc *services.EmergencyService
}

func NewEmergencyHandler(svc *services.EmergencyService) *EmergencyHandler {
	return &EmergencyHandler{svc: svc}
}

func (h *EmergencyHandler) TriggerAlert(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var input models.NewEmergencyAlert
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.AlertType == "" {
		input.AlertType = "general"
	}

	alert, err := h.svc.TriggerAlert(r.Context(), userID, input.AlertType, input.LocationLat, input.LocationLng)
	if err != nil {
		log.Error().Err(err).Msg("failed to trigger alert")
		respondError(w, http.StatusInternalServerError, "failed to trigger alert")
		return
	}

	respondJSON(w, http.StatusCreated, alert)
}

func (h *EmergencyHandler) ResolveAlert(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	alertID := chi.URLParam(r, "id")

	alert, err := h.svc.ResolveAlert(r.Context(), alertID, userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to resolve alert")
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, alert)
}

func (h *EmergencyHandler) ListAlerts(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	status := r.URL.Query().Get("status")

	alerts, err := h.svc.ListAlerts(r.Context(), userID, status)
	if err != nil {
		log.Error().Err(err).Msg("failed to list alerts")
		respondError(w, http.StatusInternalServerError, "failed to list alerts")
		return
	}

	if alerts == nil {
		alerts = []models.EmergencyAlert{}
	}

	respondJSON(w, http.StatusOK, alerts)
}
