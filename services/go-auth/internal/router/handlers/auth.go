package handlers

import (
	"encoding/json"
	"go-auth/internal/app"
	"go-auth/internal/models"
	"go-auth/internal/services"
	"log/slog"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var u models.UserCreate
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var authService *services.AuthService
	var redisDB *redis.Client
	err := app.AppContainer.Invoke(func(as *services.AuthService, rds *redis.Client) {
		redisDB = rds
		authService = as
	})
	if err != nil {
		http.Error(w, "Failed to resolve AuthService", http.StatusInternalServerError)
		return
	}
	jwt, err := authService.Create(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO refresh and access tokens
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    jwt,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		// Secure:   true, // Только для HTTPS
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    jwt,
		Path:     "/",
		MaxAge:   30 * 24 * 3600, // 30 дней
		HttpOnly: true,
		// Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	pipe := redisDB.Pipeline()
	pipe.Set(r.Context(), "access:"+jwt, jwt, 3600)
	pipe.Set(r.Context(), "refresh:"+jwt, jwt, 30*24*3600)
	_, err = pipe.Exec(r.Context())
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Redis error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logged in successfully"))

}

func Logout(w http.ResponseWriter, r *http.Request) {
	// обнулить cookie
	// убрать из redis
}

func Login(w http.ResponseWriter, r *http.Request) {
	// получить значения из body
	// валидировать значения login password
	// установить куки
	// добавить в redis

}
