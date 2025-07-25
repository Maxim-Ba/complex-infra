package kafka

import (
	"fmt"
	"go-messages/internal/app"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
)

type Consumer struct {
}

var partitionConsumer sarama.PartitionConsumer
var consumer sarama.Consumer

func NewConsumer(cfg app.AppConfig) *Consumer {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	brokers := []string{cfg.GetConfig().KafkaAddr}

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	topic := "test_topic"
	c, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	partitionConsumer = c

	return &Consumer{}
}

func (c Consumer) Close() {
	if err := consumer.Close(); err != nil {
		slog.Error(err.Error())
	}

}

func (c Consumer) StartRead() {
	// Обработка сигналов для graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	// Чтение сообщений
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			slog.Info(fmt.Sprintf("Received message: %s (Partition: %d, Offset: %d)",
				string(msg.Value), msg.Partition, msg.Offset))
		case err := <-partitionConsumer.Errors():
			slog.Error(err.Error())
		case <-signals:
			return
		}
	}
}
