package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/sex-worker-safety/checkin-dao/config"
	"github.com/sex-worker-safety/checkin-dao/db"
	"github.com/sex-worker-safety/checkin-dao/handlers"
	"github.com/sex-worker-safety/checkin-dao/routes"
	"github.com/sex-worker-safety/checkin-dao/services"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	cfg := config.Load()
	log.Info().Str("port", cfg.Port).Msg("starting checkin-dao service")

	if cfg.DatabaseURL == "" {
		log.Fatal().Msg("DATABASE_URL is required")
	}
	if cfg.JWTSecret == "" {
		log.Fatal().Msg("JWT_SECRET is required")
	}

	pool, err := db.NewPool(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer pool.Close()

	if err := db.RunMigrations(pool); err != nil {
		log.Fatal().Err(err).Msg("failed to run migrations")
	}

	// Services
	qdrantClient := services.NewQdrantClient(cfg.QdrantHost, cfg.QdrantPort)

	emergencySvc := services.NewEmergencyService(pool, cfg)
	checkInSvc := services.NewCheckInService(pool, emergencySvc)
	registrySvc := services.NewRegistryService(pool, qdrantClient)
	daoSvc := services.NewDAOService(pool)

	// Handlers
	h := &routes.Handlers{
		Auth:      handlers.NewAuthHandler(pool, cfg),
		CheckIn:   handlers.NewCheckinHandler(checkInSvc),
		Emergency: handlers.NewEmergencyHandler(emergencySvc),
		Registry:  handlers.NewRegistryHandler(registrySvc),
		DAO:       handlers.NewDAOHandler(daoSvc),
	}

	// Router
	r := chi.NewRouter()
	routes.SetupRoutes(r, h, cfg)

	// Start background check-in monitor
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go checkInSvc.CheckMissedCheckIns(ctx)

	// HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		log.Info().Msg("shutting down...")
		cancel()

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Msg("server forced to shutdown")
		}
	}()

	log.Info().Str("addr", srv.Addr).Msg("server listening")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("server error")
	}

	log.Info().Msg("server stopped")
}
