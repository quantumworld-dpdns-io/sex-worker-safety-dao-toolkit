package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/sex-worker-safety/checkin-dao/config"
	"github.com/sex-worker-safety/checkin-dao/models"
)

type AuthHandler struct {
	pool *pgxpool.Pool
	cfg  *config.Config
}

func NewAuthHandler(pool *pgxpool.Pool, cfg *config.Config) *AuthHandler {
	return &AuthHandler{pool: pool, cfg: cfg}
}

type registerRequest struct {
	WalletAddress string `json:"wallet_address"`
	DisplayName   string `json:"display_name"`
	Phone         string `json:"phone"`
}

type authResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.WalletAddress == "" {
		respondError(w, http.StatusBadRequest, "wallet_address is required")
		return
	}

	user := &models.User{
		ID:            uuid.New().String(),
		WalletAddress: req.WalletAddress,
		Role:          "user",
		DisplayName:   req.DisplayName,
		Phone:         req.Phone,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err := h.pool.Exec(r.Context(), `
		INSERT INTO users (id, wallet_address, role, display_name, phone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (wallet_address) DO UPDATE SET display_name = $4, phone = $5, updated_at = $7`,
		user.ID, user.WalletAddress, user.Role, user.DisplayName, user.Phone,
		user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user")
		respondError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	// Re-fetch in case of upsert
	err = h.pool.QueryRow(r.Context(), `SELECT id, wallet_address, role, display_name, phone, created_at, updated_at FROM users WHERE wallet_address = $1`,
		req.WalletAddress,
	).Scan(&user.ID, &user.WalletAddress, &user.Role, &user.DisplayName, &user.Phone, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Error().Err(err).Msg("failed to fetch user after upsert")
		respondError(w, http.StatusInternalServerError, "registration failed")
		return
	}

	token, err := h.generateToken(user.ID, user.Role)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate token")
		respondError(w, http.StatusInternalServerError, "authentication failed")
		return
	}

	log.Info().Str("user_id", user.ID).Str("wallet", user.WalletAddress).Msg("user registered")
	respondJSON(w, http.StatusCreated, authResponse{Token: token, User: user})
}

type loginRequest struct {
	WalletAddress string `json:"wallet_address"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.WalletAddress == "" {
		respondError(w, http.StatusBadRequest, "wallet_address is required")
		return
	}

	var user models.User
	err := h.pool.QueryRow(r.Context(), `
		SELECT id, wallet_address, role, display_name, phone, created_at, updated_at
		FROM users WHERE wallet_address = $1`, req.WalletAddress,
	).Scan(&user.ID, &user.WalletAddress, &user.Role, &user.DisplayName, &user.Phone, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "user not found")
		return
	}

	token, err := h.generateToken(user.ID, user.Role)
	if err != nil {
		log.Error().Err(err).Msg("failed to generate token")
		respondError(w, http.StatusInternalServerError, "authentication failed")
		return
	}

	log.Info().Str("user_id", user.ID).Msg("user logged in")
	respondJSON(w, http.StatusOK, authResponse{Token: token, User: &user})
}

func (h *AuthHandler) generateToken(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return tokenString, nil
}
