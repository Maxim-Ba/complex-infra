package handlers

import (
	"fmt"
	"go-messages/internal/app"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type KafkaHendler struct {
	consumer app.KConsumer
	provider app.KProducer
}

func InitKafkaHandlers(consumer app.KConsumer, provider app.KProducer) *KafkaHendler {
	return &KafkaHendler{
		consumer: consumer,
		provider: provider,
	}
}

func (h *KafkaHendler) Produce(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := h.provider.Produce("test_topic", "Hello, Kafka!")
	if err != nil {
		http.Error(w, fmt.Sprintf("Producer error: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Message produced successfully"))
}
func (h *KafkaHendler) Consume(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO make check
	w.Write([]byte("Consumer is running (check logs for received messages)"))
}
