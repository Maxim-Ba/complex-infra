package app

import (
	"context"
	"database/sql"
	"go-messages/internal/config"
	"go-messages/internal/models"

	"github.com/redis/go-redis/v9"
)

type AppConfig interface {
	GetConfig() *config.Config
}

type AppRedis interface {
	Close()
	Get(ctx context.Context, key string) *redis.StringCmd
}

type DB interface {
	Close()
	GetConnection() *sql.DB
}

type KProducer interface {
	ProduceTest() error
	Close()
}

type KConsumer interface {
	StartRead(topics []string)
	Close()
}

type MessageService interface {
	HandleMessage(ctx context.Context, m models.MessageDTO) error
}
type MongoRepository interface {
	SaveMessage(ctx context.Context, m models.MessageDTO) error
	GetMessages(ctx context.Context) ([]models.MessageDTO, error)
	Close(ctx context.Context) error
}
