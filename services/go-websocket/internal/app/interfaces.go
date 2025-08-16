package app

import (
	"context"
	"go-websocket/internal/config"
	"go-websocket/internal/models"
)

type AppConfig interface {
	GetConfig() *config.Config
}

type KProducer interface {
	Produce(topic string, value string) error
	Close()
}

type KConsumer interface {
	StartRead()
	Close()
}

type MessageService interface {
	HandleConfirmationMessage(ctx context.Context, m models.MessageDTO) error
	HandleWebRTCResponse(ctx context.Context, m models.MessageDTO) error
}

type WebSocketService interface {
	SendMessage(m models.MessageDTO) error
}
