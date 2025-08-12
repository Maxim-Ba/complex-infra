package app

import (
	"go-transport/internal/config"
)



type AppConfig interface {
	GetConfig() *config.Config
}



