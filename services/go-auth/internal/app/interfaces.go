package app

import (
	"context"
	"go-auth/internal/config"
	"go-auth/internal/models"
)

type AppAuthService interface {
	Create(user models.UserCreateReq) (*models.TokenDto, error)
}

type AppConfig interface {
	GetConfig() *config.Config
}

type AppTokenStorage interface {
	SetTokens(ctx context.Context, jwt *models.TokenDto) error
	RemoveToken(ctx context.Context, refresh, access string) error
}

type AppRedis interface {
	Close()
}
type AppUserStorage interface {
	Save(user models.UserCreateReq) models.UserCreateRes
	Get(user models.UserCreateDto) (*models.UserCreateRes, error)
}
