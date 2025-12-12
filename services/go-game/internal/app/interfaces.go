package app

import (
	"context"
	"go-game/internal/config"
	"go-game/internal/models"
)

type AppConfig interface {
	GetConfig() *config.Config
}
type MessageService interface {
	HandleMessage(context.Context, models.MessageDTO) error
}
type KProducer interface {
	Produce(topic string, value string) error
	Close()
}

type KConsumer interface {
	StartRead()
	Close()
}

