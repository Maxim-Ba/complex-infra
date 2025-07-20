package router

import (
	"go-auth/internal/app"
	"go-auth/internal/router/handlers"
	"go-auth/internal/router/middlewares"
	"go-auth/pkg/metrics"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func New() *chi.Mux {
	var authService app.AppAuthService

	if err := app.AppContainer.Invoke(func(as app.AppAuthService) {
		authService = as
	}); err != nil {
		panic("router New, get authService: " + err.Error())
	}
	router := chi.NewRouter()
	router.Use(metrics.MetricsMiddleware)
	router.Get("/handler", handlers.Handler)
	router.Get("/error", handlers.EmitError)
	router.Group(func(r chi.Router) {
		r.Use(middlewares.WithNoAuthOnly)
		r.Post("/register", handlers.Register)
	})

	router.Group(func(r chi.Router) {
		r.Use(func(h http.Handler) http.Handler {
			return middlewares.WithAuth(h, authService)
		})
		r.Post("/login", handlers.Login)
		r.Get("/logout", handlers.Logout)
	})

	return router
}
