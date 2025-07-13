package main

import (
	"errors"
	"fmt"
	"go-auth/internal/app"
	"go-auth/internal/config"
	"go-auth/internal/router"
	"go-auth/internal/services"
	"go-auth/internal/storage"
	"go-auth/pkg/metrics"
	"go-auth/pkg/redis"
	"go-auth/pkg/tracer"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/dig"
)

func main() {
	fmt.Println("start go-auth")
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	initAppContainer()

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
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			slog.Error(err.Error())
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			panic(err)
		}
	}()

	<-exit

	if err := httpServer.Close(); err != nil {
		slog.Error(err.Error())
	}

	if err := app.AppContainer.Invoke(func(rds *redis.Redis) {
		rds.Close()
	}); err != nil {
		slog.Error(err.Error())
	}

}

func initAppContainer() {
	app.AppContainer = dig.New()

	if err := app.AppContainer.Provide(services.AuthNew, dig.As(new(app.AppAuthService))); err != nil {
		panic(fmt.Sprintf("auth service can not be provided: %s", err.Error()))
	}
	if err := app.AppContainer.Provide(config.New, dig.As(new(app.AppConfig))); err != nil {
		panic(fmt.Sprintf("config can not be provided:%s", err.Error()))
	}
	if err := app.AppContainer.Provide(redis.New, dig.As(new(app.AppRedis))); err != nil {
		panic(fmt.Sprintf("redis can not be provided: %s", err.Error()))
	}
	if err := app.AppContainer.Provide(storage.NewTokenStorage, dig.As(new(app.AppTokenStorage))); err != nil {
		panic(fmt.Sprintf("token storage can not be provided: %s", err.Error()))
	}
	if err := app.AppContainer.Provide(storage.NewUserStorage, dig.As(new(app.AppUserStorage))); err != nil {
		panic(fmt.Sprintf("user storage not be provided: %s", err.Error()))
	}
}
