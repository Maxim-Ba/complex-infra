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
	"errors"
	"fmt"
	"go-messages/cmd/wire"
	"go-messages/internal/router"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)



func main() {
	fmt.Println("start go-messages")
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))


	deps:= wire.Initialize()
	r:= router.New(deps.KafkaHendler)

	// Запуск консьюмера в горутине
	go deps.Consumer.StartRead()

	
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	fmt.Println("Server started at :8080")
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
	deps.Consumer.Close()
	deps.Producer.Close()

}
