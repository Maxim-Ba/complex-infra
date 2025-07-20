package router

import (
	"go-auth/internal/app"
	"go-auth/internal/constants"
	"go-auth/internal/router/handlers"
	"go-auth/internal/router/middlewares"
	"go-auth/pkg/metrics"
	"net/http"

	_ "go-auth/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Go Auth API
// @version 1.0
// @description Authentication service API

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func New() *chi.Mux {
	var authService app.AppAuthService

	if err := app.AppContainer.Invoke(func(as app.AppAuthService) {
		authService = as
	}); err != nil {
		panic("router New, get authService: " + err.Error())
	}
	// TODO add handler middleware
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", constants.AccessTokenCookie, constants.RefreshTokenCookie},
		AllowCredentials: true,
		// MaxAge:           300,
	}))
	router.Use(metrics.MetricsMiddleware)
	router.Get("/handler", handlers.Handler)
	router.Get("/error", handlers.EmitError)
	
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // URL для json документации
	))
	router.Group(func(r chi.Router) {
		r.Use(middlewares.WithNoAuthOnly)
		r.Post("/register", handlers.Register)
		r.Post("/login", handlers.Login)
		r.Get("/logout", handlers.Logout)
	})

	router.Group(func(r chi.Router) {
		r.Use(func(h http.Handler) http.Handler {
			return middlewares.WithAuth(h, authService)
		})
		r.Get("/logout", handlers.Logout)
	})

	return router
}
