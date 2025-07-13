package handlers

import (
	"encoding/json"
	"errors"
	"go-auth/internal/app"
	"go-auth/internal/models"
	"go-auth/internal/services"
	"log/slog"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var u models.UserCreateReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var authService app.AppAuthService
	var tokenStore app.AppTokenStorage

	if err := app.AppContainer.Invoke(func(as app.AppAuthService, s app.AppTokenStorage) {
		tokenStore = s
		authService = as
	}); err != nil {
		http.Error(w, "Failed to resolve AuthService", http.StatusInternalServerError)
		return
	}
	jwt, err := authService.Create(u)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    jwt.Access,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		// Secure:   true, // Только для HTTPS
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    jwt.Refresh,
		Path:     "/",
		MaxAge:   30 * 24 * 3600, // 30 дней
		HttpOnly: true,
		// Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	if err := tokenStore.SetTokens(r.Context(), jwt); err != nil {
		slog.Error(err.Error())
		http.Error(w, "token store error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))

}

func Logout(w http.ResponseWriter, r *http.Request) {
	access, err := r.Cookie("access_token")
	if err != nil {

		slog.Error(err.Error())
		http.Error(w, "Failed to read cookie access", http.StatusInternalServerError)
		return
	}
	refresh, err := r.Cookie("refresh_token")
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Failed to read cookie refresh", http.StatusInternalServerError)
		return
	}

	var tokenStore app.AppTokenStorage
	err = app.AppContainer.Invoke(func(s app.AppTokenStorage) {
		tokenStore = s
	})
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Failed to resolve AppTokenStorage", http.StatusInternalServerError)
		return
	}

	err = tokenStore.RemoveToken(r.Context(), refresh.Value, access.Value)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Failed to remove token", http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		// Secure:   true, // Только для HTTPS
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		// Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
}

func Login(w http.ResponseWriter, r *http.Request) {

	var u models.UserCreateReq
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var authService app.AppAuthService
	var tokenStore app.AppTokenStorage
	if err := app.AppContainer.Invoke(func(as app.AppAuthService, s app.AppTokenStorage) {
		tokenStore = s
		authService = as
	}); err != nil {
		http.Error(w, "Failed to resolve AuthService & TokenStorage", http.StatusInternalServerError)
		return
	}
	jwt, err := authService.Login(u)
	if err != nil {
		if errors.Is(err, services.ErrWrongLoginOrPassword) || errors.Is(err, services.ErrLoginAndPasswordAreRequired){
			http.Error(w, "Wrong login or password", http.StatusBadRequest)
			return 
		}
		slog.Error(err.Error())
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    jwt.Access,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		// Secure:   true, // Только для HTTPS
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    jwt.Refresh,
		Path:     "/",
		MaxAge:   30 * 24 * 3600, // 30 дней
		HttpOnly: true,
		// Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	if err := tokenStore.SetTokens(r.Context(), jwt); err != nil {
		slog.Error(err.Error())
		http.Error(w, "token store error", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(""))
}
