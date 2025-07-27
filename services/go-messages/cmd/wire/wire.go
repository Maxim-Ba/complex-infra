//go:build wireinject
// +build wireinject

package wire

import (
	"go-messages/internal/app"
	"go-messages/internal/config"
	"go-messages/internal/handlers"
	"go-messages/internal/services"
	"go-messages/pkg/kafka"

	"github.com/google/wire"
)

type Dependenсies struct {
	Producer  *kafka.Producer
	Consumer *kafka.Consumer
	KafkaHendler  *handlers.KafkaHendler
}

func Initialize() (*Dependenсies, error) {
	wire.Build(
		config.New,
		wire.Bind(new(app.AppConfig), new(*config.Config)),
		kafka.NewProducer,
		wire.Bind(new(app.KProducer), new(*kafka.Producer)),
		kafka.NewConsumer,
		wire.Bind(new(app.KConsumer), new(*kafka.Consumer)),
		services.New,
				wire.Bind(new(app.MessageService), new(*services.MessageService)),

		handlers.InitKafkaHandlers,
		wire.Struct(new(Dependenсies), "*"),
	)
	return &Dependenсies{}, nil // Эта строка не выполнится, Wire заменит её.
}
