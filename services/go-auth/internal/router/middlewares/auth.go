package middlewares

import (
	"go-auth/internal/app"
	"go-auth/internal/constants"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"time"
)

func WithAuth(next http.Handler, authService app.AppAuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем наличие access token
		accessTokenCookie := getCookie(r, constants.AccessTokenCookie)
		if accessTokenCookie == nil {
			slog.Info("WithAuth accessTokenCookie is nil")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Проверяем срок жизни access token
		if !isTokenValid(accessTokenCookie) {
			slog.Info("WithAuth accessTokenCookie is not valid")
			// Если access token истек, проверяем refresh token
			refreshTokenCookie := getCookie(r, constants.RefreshTokenCookie)
			if refreshTokenCookie == nil || !isTokenValid(refreshTokenCookie) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			jwt, err := authService.RefreshToken(refreshTokenCookie.Value)
			if err != nil {
				slog.Info("WithAuth set refressh token")
				w.WriteHeader(http.StatusUnauthorized)
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
		}

		next.ServeHTTP(w, r)
	})
}
func WithNoAuthOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем наличие хотя бы одного из токенов
		cookies := r.Cookies()
		loggableCookies := make([]map[string]string, len(cookies))
		for i, cookie := range cookies {
			loggableCookies[i] = map[string]string{
				"Name":     cookie.Name,
				"Value":    cookie.Value,
				"Path":     cookie.Path,
				"SameSite": strconv.Itoa(int(cookie.SameSite)),
				"Domain":   cookie.Domain,
				// Можно добавить другие поля, если нужно (Path, Domain, Expires и т.д.)
			}
		}
		slog.Info("Cookies", "cookies", loggableCookies)
		if getCookie(r, constants.AccessTokenCookie) != nil ||
			getCookie(r, constants.RefreshTokenCookie) != nil {
			slog.Info("WithNoAuthOnly must not are tokens")

			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// TODO Rate limiting - защита от брутфорса
// TODO logging not valid tries
func getCookie(r *http.Request, name string) *http.Cookie {
	cookies := r.Cookies()
	idx := slices.IndexFunc(cookies, func(c *http.Cookie) bool {
		return c.Name == name
	})
	if idx == -1 {
		return nil
	}
	return cookies[idx]
}

func isTokenValid(cookie *http.Cookie) bool {
	// Если у куки нет Expires, считаем ее валидной
	if cookie.Expires.IsZero() {
		return true
	}
	// Проверяем не истек ли срок
	return cookie.Expires.After(time.Now())
}
