package main

import (
	"errors"
	"fmt"
	"go-auth/pkg/metrics"
	"go-auth/pkg/tracer"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	fmt.Println("start go-auth")

	_, err := tracer.InitTracer("jaeger:4318", "go-auth service")
	if err != nil {
		slog.Error(err.Error())
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))

	metrics.Start(":8081")
	
	router := chi.NewRouter()
	router.Use(metrics.MetricsMiddleware)
	router.Get("/", handler)
	router.Get("/error", emitError)

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
	// _, span:= 	tracer.Tracer.Start(r.Context(), "handler")
	t := otel.Tracer("handler")
	_, span := t.Start(r.Context(), "handler")
	defer span.End()
	slog.Info("handler end")
	w.WriteHeader(http.StatusOK)
}

func emitError(w http.ResponseWriter, r *http.Request) {
	slog.Info("emitError start")

	t := otel.Tracer("emitError")
	ctx, span := t.Start(r.Context(), "emitError")
	defer span.End()
	someError := errors.New("some error")
	if sp := trace.SpanFromContext(ctx); sp.IsRecording() {
		sp.RecordError(someError)
		sp.SetStatus(codes.Error, someError.Error())
	}

	slog.Info("emitError end")
	http.Error(w, "Error handler", http.StatusBadRequest)

}
