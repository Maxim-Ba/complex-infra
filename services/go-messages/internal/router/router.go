package router

import (
	"go-messages/internal/handlers"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func New(kh *handlers.KafkaHendler, mh *handlers.MessageHandler) *httprouter.Router {

	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Write([]byte("Service is running"))
	})

	// Эндпоинт для проверки producer
	router.GET("/produce", kh.Produce)

	// Эндпоинт для проверки consumer
	router.GET("/consume", kh.Consume)

	router.GET("/messages/:groupID", mh.Get)
	return router
}
