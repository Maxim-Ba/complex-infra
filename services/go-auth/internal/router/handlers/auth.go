package handlers

import (
	"encoding/json"
	"errors"
	"go-auth/internal/app"
	"go-auth/internal/constants"
	"go-auth/internal/models"
	"go-auth/internal/services"
	"log/slog"
	"net/http"
)

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя и возвращает JWT токены в cookies
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param input body models.UserCreateReq true "Данные для регистрации"
// @Success 204 "Успешная регистрация, токены установлены в cookies"
// @Failure 400 {string} string "Неверный запрос или неверные логин/пароль"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /register [post]
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
		if errors.Is(err, services.ErrLoginAndPasswordAreRequired){
			http.Error(w, "Wrong login or password", http.StatusBadRequest)
			return 
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     constants.AccessTokenCookie,
		Value:    jwt.Access,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		// Secure:   true, // Только для HTTPS
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     constants.RefreshTokenCookie,
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

// Logout godoc
// @Summary Выход из системы
// @Description Удаляет JWT токены из cookies и хранилища
// @Tags Аутентификация
// @Produce json
// @Success 204 "Успешный выход, токены удалены"
// @Failure 400 {string} string "Неверный запрос"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /logout [get]
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

// Login godoc
// @Summary Аутентификация пользователя
// @Description Проверяет учетные данные и возвращает JWT токены в cookies
// @Tags Аутентификация
// @Accept json
// @Produce json
// @Param input body models.UserCreateReq true "Учетные данные"
// @Success 204 "Успешная аутентификация, токены установлены в cookies"
// @Failure 400 {string} string "Неверный запрос или неверные логин/пароль"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /login [post]
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
