package middlewares

import "net/http"

func WithAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO проверить есть ли кука
		next.ServeHTTP(w, r)
	})
}
func WithNoAuthOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO проверить есть ли кука, если есть , то выдать код ошибки
		next.ServeHTTP(w, r)
	})
}

// TODO Rate limiting - защита от брутфорса
// TODO logging not valid tries
