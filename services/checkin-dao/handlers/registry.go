package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/sex-worker-safety/checkin-dao/models"
	"github.com/sex-worker-safety/checkin-dao/services"
)

type RegistryHandler struct {
	svc *services.RegistryService
}

func NewRegistryHandler(svc *services.RegistryService) *RegistryHandler {
	return &RegistryHandler{svc: svc}
}

func (h *RegistryHandler) SubmitReport(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var input models.NewReport
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	report, err := h.svc.SubmitReport(r.Context(), userID, input)
	if err != nil {
		log.Error().Err(err).Msg("failed to submit report")
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, report)
}

func (h *RegistryHandler) SearchReports(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	status := r.URL.Query().Get("status")

	reports, err := h.svc.SearchReports(r.Context(), query, category, status)
	if err != nil {
		log.Error().Err(err).Msg("failed to search reports")
		respondError(w, http.StatusInternalServerError, "failed to search reports")
		return
	}

	if reports == nil {
		reports = []models.BadClientReport{}
	}

	respondJSON(w, http.StatusOK, reports)
}

func (h *RegistryHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	reportID := chi.URLParam(r, "id")

	report, err := h.svc.GetReport(r.Context(), reportID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get report")
		respondError(w, http.StatusNotFound, "report not found")
		return
	}

	respondJSON(w, http.StatusOK, report)
}

func (h *RegistryHandler) UpdateReportStatus(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	reportID := chi.URLParam(r, "id")

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	report, err := h.svc.UpdateReportStatus(r.Context(), reportID, userID, req.Status)
	if err != nil {
		log.Error().Err(err).Msg("failed to update report status")
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, report)
}
