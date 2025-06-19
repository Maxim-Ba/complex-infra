package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	fmt.Println("start go-auth")
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))
	router := chi.NewRouter()
	router.Get("/", handler)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	slog.Info("server start")
	if err := httpServer.ListenAndServe(); err != nil {
		slog.Error(err.Error())
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	slog.Info("handler start")
	slog.Info("handler end")
	w.WriteHeader(http.StatusOK)
}
