//go:build wireinject
// +build wireinject

package wire

import (
	"go-game/internal/app"
	"go-game/internal/config"
	"go-game/internal/services"
	"go-game/pkg/kafka"
	"go-game/pkg/webrtc"

	"github.com/google/wire"
)

type Dependenсies struct {
	Producer       app.KProducer
	Consumer       app.KConsumer
	MessageService *app.MessageService
}

func Initialize() (*Dependenсies, error) {
	wire.Build(
		config.New,
		wire.Bind(new(app.AppConfig), new(*config.Config)),
		kafka.NewProducer,
		wire.Bind(new(app.KProducer), new(*kafka.Producer)),
		kafka.NewConsumer,
		wire.Bind(new(app.KConsumer), new(*kafka.Consumer)),

		webrtc.NewRTCManager,
		services.NewMessageService,
		wire.Bind(new(app.MessageService), new(*services.MessageService)),
		wire.Struct(new(Dependenсies), "*"),
	)
	return &Dependenсies{}, nil
}
