package kafka

import (
	"context"
	"encoding/json"
	"go-game/internal/app"
	"go-game/internal/models"

	"log/slog"
	"os"
	"os/signal"
	"sync"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumerGroup sarama.ConsumerGroup
	handler       app.MessageService
	ready         chan bool
	topics []string
}


func NewConsumer(cfg app.AppConfig, handler app.MessageService) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0 
	config.Consumer.Return.Errors = true
	c := cfg.GetConfig()

	// Настройка ручного управления offset'ами
	config.Consumer.Offsets.AutoCommit.Enable = false     // Отключаем авто-коммит
	config.Consumer.Offsets.Initial = sarama.OffsetNewest // Начинаем чтение с новых сообщений

	consumerGroup, err := sarama.NewConsumerGroup(c.KafkaBrokers, c.KafkaGroupId, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumerGroup: consumerGroup,
		handler:       handler,
		ready:         make(chan bool),
		topics:  []string{ c.GetConfig().RTCSignalTopic},
	}, nil
}

func (c Consumer) Close() {
	if err := c.consumerGroup.Close(); err != nil {
		slog.Error(err.Error())
	}

}

func (c Consumer) StartRead() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			// Запускаем потребление сообщений
			if err := c.consumerGroup.Consume(ctx, c.topics, &c); err != nil {
				slog.Error("Error from consumer", "error", err)
			}

			// Проверяем, не завершен ли контекст
			if ctx.Err() != nil {
				return
			}
			c.ready = make(chan bool)
		}
	}()

	// Ожидаем, пока потребитель не будет готов
	<-c.ready
	slog.Info("Sarama consumer up and running...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt)

	select {
	case <-ctx.Done():
		slog.Info("Terminating: context cancelled")
	case <-sigterm:
		slog.Info("Terminating: via signal")
	}
	cancel()
	wg.Wait()

	if err := c.consumerGroup.Close(); err != nil {
		slog.Error("Error closing consumer", "error", err)
	}
}

// Setup вызывается при начале новой сессии.
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

// Cleanup вызывается при завершении сессии.
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim обрабатывает сообщения из партиции.
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		slog.Info("New message received",
			"topic", msg.Topic,
			"partition", msg.Partition,
			"offset", msg.Offset,
			"value", string(msg.Value),
		)
		var m models.MessageDTO
		if err := json.Unmarshal(msg.Value, &m); err != nil {
			slog.Error("Failed to unmarshal message", 
				"error", err,
				"value", string(msg.Value),
			)
			continue // Пропускаем сообщение при ошибке декодирования
		}

		// Обрабатываем сообщение
		if err := c.handler.HandleMessage(session.Context(), m); err != nil {
			slog.Error("Failed to handle message", "error", err)
			continue // Не подтверждаем offset при ошибке
		}

		// Вручную подтверждаем успешную обработку
		session.MarkMessage(msg, "")
		session.Commit() // Явный коммит
	}
	return nil
}
