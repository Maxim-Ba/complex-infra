// Package main реализует сервис для обработки сообщений из Kafka.
//
// Сервис предоставляет:
// - HTTP API для мониторинга состояния сервиса
// - Тестовые эндпоинты для проверки работы Producer/Consumer Kafka
// - Механизм graceful shutdown при получении сигналов завершения
//
// Основные компоненты:
// - Kafka Producer - отправка тестовых сообщений
// - Kafka Consumer - чтение и обработка сообщений
// - HTTP Server - REST API для взаимодействия
package main

import (
	"context"
	"errors"
	"fmt"
	"go-messages/cmd/wire"
	"go-messages/internal/router"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("start go-messages")
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))

	deps, err := wire.Initialize()
	if err != nil {
		panic(fmt.Sprintf("Error on wire.Initialize() %v", err))
	}
	r := router.New(deps.KafkaHendler, deps.MessageHandler)

	// Запуск консьюмера в горутине
	go deps.Consumer.StartRead()

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	fmt.Println("Server started at :8080")
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			panic(fmt.Sprintf("Error on httpServer.ListenAndServe() %v", err))
		}
	}()
	<-exit

	if err := httpServer.Close(); err != nil {
		slog.Error(err.Error())
	}
	deps.Consumer.Close()
	deps.Producer.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	deps.MongoRepository.Close(ctx)

}
