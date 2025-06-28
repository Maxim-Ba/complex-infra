package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	slog.Info("handler start")
	// _, span:= 	tracer.Tracer.Start(r.Context(), "handler")
	t := otel.Tracer("handler")
	_, span := t.Start(r.Context(), "handler")
	defer span.End()
	slog.Info("handler end")
	w.WriteHeader(http.StatusOK)
}

func EmitError(w http.ResponseWriter, r *http.Request) {
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
