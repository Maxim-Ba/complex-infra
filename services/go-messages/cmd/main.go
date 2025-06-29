package main

import (
	"fmt"
	"go-messages/pkg/kafka"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
fmt.Println("start go-messages")

	// Инициализация продюсера и консьюмера
	p := kafka.NewProducer()

	defer p.Close()

	c := kafka.NewConsumer()
	
	defer c.Close()

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

	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", router)
}
