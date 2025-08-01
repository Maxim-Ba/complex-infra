package redis

import (
	"context"
	"go-auth/internal/app"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func (r *Redis) Close() {
	r.client.Close()
}
func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(ctx, key, value, expiration)
}

func (r *Redis) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.client.Get(ctx, key)
}

func (r *Redis) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return r.client.Del(ctx, keys...)
}

func (r *Redis) Pipeline() redis.Pipeliner {
	return r.client.Pipeline()
}

func New(cfg app.AppConfig) *Redis {
	var redisAddr = cfg.GetConfig().RedisAddr

	slog.Info("Connecting to Redis", "address", redisAddr)

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
		PoolSize: 10,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	slog.Info("Successfully connected to Redis")

	return &Redis{client: rdb}
}
