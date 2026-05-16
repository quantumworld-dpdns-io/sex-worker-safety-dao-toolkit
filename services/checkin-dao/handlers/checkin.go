package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/sex-worker-safety/checkin-dao/models"
	"github.com/sex-worker-safety/checkin-dao/services"
)

type CheckinHandler struct {
	svc *services.CheckInService
}

func NewCheckinHandler(svc *services.CheckInService) *CheckinHandler {
	return &CheckinHandler{svc: svc}
}

func (h *CheckinHandler) CreateCheckIn(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var input models.NewCheckIn
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	chk, err := h.svc.ScheduleCheckIn(r.Context(), userID, input)
	if err != nil {
		log.Error().Err(err).Msg("failed to create check-in")
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, chk)
}

func (h *CheckinHandler) CompleteCheckIn(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	checkInID := chi.URLParam(r, "id")

	chk, err := h.svc.CompleteCheckIn(r.Context(), checkInID, userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to complete check-in")
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, chk)
}

func (h *CheckinHandler) ListCheckIns(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	status := r.URL.Query().Get("status")

	checkins, err := h.svc.ListCheckIns(r.Context(), userID, status)
	if err != nil {
		log.Error().Err(err).Msg("failed to list check-ins")
		respondError(w, http.StatusInternalServerError, "failed to list check-ins")
		return
	}

	if checkins == nil {
		checkins = []models.CheckIn{}
	}

	respondJSON(w, http.StatusOK, checkins)
}

func (h *CheckinHandler) GetCheckIn(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	checkInID := chi.URLParam(r, "id")

	chk, err := h.svc.GetCheckIn(r.Context(), checkInID, userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get check-in")
		respondError(w, http.StatusNotFound, "check-in not found")
		return
	}

	respondJSON(w, http.StatusOK, chk)
}
