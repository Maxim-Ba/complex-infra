package main

import (
	"errors"
	"fmt"
	"go-auth/internal/app"
	"go-auth/internal/config"
	"go-auth/internal/router"
	"go-auth/internal/services"
	"go-auth/pkg/metrics"
	"go-auth/pkg/redis"
	"go-auth/pkg/tracer"
	"log/slog"
	"net/http"
	"os"

	"go.uber.org/dig"
)

func main() {
	fmt.Println("start go-auth")

	app.AppContainer = dig.New()

	if err := app.AppContainer.Provide(services.AuthNew); err != nil {
		panic(fmt.Sprintf("auth service can not be provided %s", err.Error()))
	}
	if err := app.AppContainer.Provide(config.New); err != nil {
		panic(fmt.Sprintf("config can not be provided %s", err.Error()))
	}
	if err := app.AppContainer.Provide(redis.New); err != nil {
		panic(fmt.Sprintf("redis can not be provided %s", err.Error()))
	}

	_, err := tracer.InitTracer("jaeger:4318", "go-auth service")
	if err != nil {
		slog.Error(err.Error())
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))

	metrics.Start(":8081")

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router.New(),
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
