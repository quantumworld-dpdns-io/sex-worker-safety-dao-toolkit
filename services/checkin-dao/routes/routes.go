package routes

import (
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/sex-worker-safety/checkin-dao/config"
	"github.com/sex-worker-safety/checkin-dao/handlers"
	"github.com/sex-worker-safety/checkin-dao/middleware"
)

type Handlers struct {
	Auth     *handlers.AuthHandler
	CheckIn  *handlers.CheckinHandler
	Emergency *handlers.EmergencyHandler
	Registry *handlers.RegistryHandler
	DAO      *handlers.DAOHandler
}

func SetupRoutes(r chi.Router, h *Handlers, cfg *config.Config) {
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(middleware.RequestLogging)
	r.Use(middleware.CORSMiddleware)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.Auth.Register)
			r.Post("/login", h.Auth.Login)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(cfg))

			r.Route("/checkins", func(r chi.Router) {
				r.Post("/", h.CheckIn.CreateCheckIn)
				r.Get("/", h.CheckIn.ListCheckIns)
				r.Get("/{id}", h.CheckIn.GetCheckIn)
				r.Post("/{id}/complete", h.CheckIn.CompleteCheckIn)
			})

			r.Route("/emergency", func(r chi.Router) {
				r.Post("/", h.Emergency.TriggerAlert)
				r.Get("/", h.Emergency.ListAlerts)
				r.Post("/{id}/resolve", h.Emergency.ResolveAlert)
			})

			r.Route("/registry", func(r chi.Router) {
				r.Route("/reports", func(r chi.Router) {
					r.Post("/", h.Registry.SubmitReport)
					r.Get("/", h.Registry.SearchReports)
					r.Get("/{id}", h.Registry.GetReport)

					r.Group(func(r chi.Router) {
						r.Use(middleware.RoleAuth("moderator", "admin"))
						r.Patch("/{id}/status", h.Registry.UpdateReportStatus)
					})
				})
			})

			r.Route("/dao", func(r chi.Router) {
				r.Route("/proposals", func(r chi.Router) {
					r.Post("/", h.DAO.CreateProposal)
					r.Get("/", h.DAO.ListProposals)
					r.Get("/{id}", h.DAO.GetProposal)
					r.Post("/{id}/votes", h.DAO.CastVote)
					r.Get("/{id}/results", h.DAO.GetResults)
				})
			})
		})
	})
}
