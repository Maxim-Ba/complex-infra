package main

import (
	"errors"
	"fmt"
	"go-messages/pkg/kafka"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/julienschmidt/httprouter"
)

func main() {
	fmt.Println("start go-messages")
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})))

	// Инициализация продюсера и консьюмера
	p := kafka.NewProducer()

	c := kafka.NewConsumer()

	// Запуск консьюмера в горутине
	go c.StartRead()

	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Write([]byte("Service is running"))
	})

	// Эндпоинт для проверки producer
	router.GET("/produce", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		err := p.ProduceTest()
		if err != nil {
			http.Error(w, fmt.Sprintf("Producer error: %v", err), http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Message produced successfully"))
	})

	// Эндпоинт для проверки consumer
	router.GET("/consume", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Write([]byte("Consumer is running (check logs for received messages)"))
	})

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
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
	c.Close()
	p.Close()

}
