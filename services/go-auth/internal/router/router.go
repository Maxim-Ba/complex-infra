package router

import (
	"go-auth/internal/router/handlers"
	"go-auth/internal/router/middlewares"
	"go-auth/pkg/metrics"

	"github.com/go-chi/chi/v5"
)

func New() *chi.Mux {
	router := chi.NewRouter()
	router.Use(metrics.MetricsMiddleware)
	router.Get("/handler", handlers.Handler)
	router.Get("/error", handlers.EmitError)
		router.Group(func(r chi.Router) {
			r.Use(middlewares.WithNoAuthOnly)
			r.Post("/register", handlers.Register)
		})

	router.Group(func(r chi.Router) {
		r.Use(middlewares.WithAuth)
		r.Post("/login", handlers.Login)
		r.Get("/logout", handlers.Logout)
	})
	
	return router
}
