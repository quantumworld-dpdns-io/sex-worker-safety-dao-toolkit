package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"github.com/sex-worker-safety/checkin-dao/models"
	"github.com/sex-worker-safety/checkin-dao/services"
)

type DAOHandler struct {
	svc *services.DAOService
}

func NewDAOHandler(svc *services.DAOService) *DAOHandler {
	return &DAOHandler{svc: svc}
}

func (h *DAOHandler) CreateProposal(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var input models.NewProposal
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	proposal, err := h.svc.CreateProposal(r.Context(), userID, input)
	if err != nil {
		log.Error().Err(err).Msg("failed to create proposal")
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, proposal)
}

func (h *DAOHandler) ListProposals(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")

	proposals, err := h.svc.ListProposals(r.Context(), status)
	if err != nil {
		log.Error().Err(err).Msg("failed to list proposals")
		respondError(w, http.StatusInternalServerError, "failed to list proposals")
		return
	}

	if proposals == nil {
		proposals = []models.DAOProposal{}
	}

	respondJSON(w, http.StatusOK, proposals)
}

func (h *DAOHandler) GetProposal(w http.ResponseWriter, r *http.Request) {
	proposalID := chi.URLParam(r, "id")

	proposal, err := h.svc.GetProposal(r.Context(), proposalID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get proposal")
		respondError(w, http.StatusNotFound, "proposal not found")
		return
	}

	respondJSON(w, http.StatusOK, proposal)
}

type castVoteRequest struct {
	Vote string `json:"vote"`
}

func (h *DAOHandler) CastVote(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	proposalID := chi.URLParam(r, "id")

	var req castVoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	vote, err := h.svc.CastVote(r.Context(), proposalID, userID, req.Vote)
	if err != nil {
		log.Error().Err(err).Msg("failed to cast vote")
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, vote)
}

func (h *DAOHandler) GetResults(w http.ResponseWriter, r *http.Request) {
	proposalID := chi.URLParam(r, "id")

	results, err := h.svc.GetResults(r.Context(), proposalID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get results")
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, results)
}
