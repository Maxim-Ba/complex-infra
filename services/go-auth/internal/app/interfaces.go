package app

import (
	"context"
	"go-auth/internal/config"
	"go-auth/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type AppAuthService interface {
	Create(user models.UserCreateReq) (*models.TokenDto, error)
	Login(user models.UserCreateReq) (*models.TokenDto, error)
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
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Pipeline() redis.Pipeliner
}
type AppUserStorage interface {
	Save(user models.UserCreateReq) models.UserCreateRes
	Get(user models.UserCreateDto) (*models.UserCreateRes, error)
}
