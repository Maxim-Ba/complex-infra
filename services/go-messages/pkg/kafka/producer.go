package kafka

import (
	"go-messages/internal/app"
	"log"
	"log/slog"

	"github.com/IBM/sarama"
)

type Producer struct {
}

var producer sarama.SyncProducer

// TODO заменить логер
func NewProducer(cfg app.AppConfig) *Producer {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	brokers := []string{cfg.GetConfig().KafkaAddr}

	p, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}

	producer = p
	return &Producer{}
}

func (p Producer) Close() {
	if err := producer.Close(); err != nil {
		slog.Error(err.Error())
	}

}

func (p Producer) Produce(topic string , value string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(value),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Message sent to partition %d at offset %d", partition, offset)
	return nil
}
