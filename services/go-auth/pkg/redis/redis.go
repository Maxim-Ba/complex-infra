package redis

import (
	"context"
	"fmt"
	"go-auth/internal/app"
	"go-auth/internal/config"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

func New() *redis.Client {
	var redisAddr string
	err := app.AppContainer.Invoke(func(cfg *config.Config) {
		redisAddr = cfg.RedisAddr
	})
	if err != nil {
		panic(fmt.Errorf("error with invoke config, %w", err))
	}
	slog.Info("Connecting to Redis", "address", redisAddr)

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr, // адрес Redis сервера
		Password: "",        // пароль, если есть
		DB:       0,         // номер базы данных
		PoolSize: 10,        // ???

	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(pong)
	slog.Info("Successfully connected to Redis")
	// TODO defer close 
	return rdb
}
