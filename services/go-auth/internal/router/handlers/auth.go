package handlers

import (
	"encoding/json"
	"go-auth/internal/app"
	"go-auth/internal/models"
	"go-auth/internal/services"
	"go-auth/internal/storage"
	"log/slog"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var u models.UserCreateReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var authService *services.AuthService
	var tokenStore *storage.TokenStorage
	err := app.AppContainer.Invoke(func(as *services.AuthService, s *storage.TokenStorage) {
		tokenStore = s
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

	if err := tokenStore.SetTokens(r.Context(), jwt, jwt); err != nil {
		slog.Error(err.Error())
		http.Error(w, "token store error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))

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
