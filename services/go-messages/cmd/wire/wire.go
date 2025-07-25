//go:build wireinject
// +build wireinject

package wire

import (
	"go-messages/internal/app"
	"go-messages/internal/config"
	"go-messages/internal/handlers"
	"go-messages/pkg/kafka"

	"github.com/google/wire"
)

type Dependenсies struct {
	Producer  *kafka.Producer
	Consumer *kafka.Consumer
	KafkaHendler  *handlers.KafkaHendler
}

func Initialize() *Dependenсies {
	wire.Build(
		config.New,
		wire.Bind(new(app.AppConfig), new(*config.Config)),
		kafka.NewProducer,
		wire.Bind(new(app.KProducer), new(*kafka.Producer)),
		kafka.NewConsumer,
		wire.Bind(new(app.KConsumer), new(*kafka.Consumer)),
		handlers.InitKafkaHandlers,
		wire.Struct(new(Dependenсies), "*"),
	)
	return &Dependenсies{} // Эта строка не выполнится, Wire заменит её.
}
