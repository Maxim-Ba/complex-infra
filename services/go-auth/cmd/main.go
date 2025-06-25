package main

import (
	"errors"
	"fmt"
	"go-auth/pkg/tracer"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	fmt.Println("start go-auth")

	_, err := tracer.InitTracer("http://localhost:14268/api/traces", "Note Service")
if err != nil {
	slog.Error("init tracer", err)
}
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
_, span:= 	tracer.Tracer.Start(r.Context(), "handler")
defer span.End()
	slog.Info("handler end")
	w.WriteHeader(http.StatusOK)
}


