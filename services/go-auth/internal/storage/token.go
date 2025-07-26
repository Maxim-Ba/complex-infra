package storage

import (
	"context"
	"fmt"
	"go-auth/internal/app"
	"go-auth/internal/models"
	"log/slog"
	"time"
)

type TokenStorage struct {
	redisDB app.AppRedis
}

func NewTokenStorage(redisDB app.AppRedis) *TokenStorage {
	return &TokenStorage{redisDB: redisDB}
}

func (s *TokenStorage) SetTokens(ctx context.Context, jwt *models.TokenDto) error {
	slog.Info("TokenStorage SetTokens")
	if jwt.Access == "" || jwt.Refresh == "" {
		slog.Error("Empty tokens provided",
			"access_empty", jwt.Access == "",
			"refresh_empty", jwt.Refresh == "",
		)
		// return fmt.Errorf("empty tokens")
	}
	pipe := s.redisDB.Pipeline()
	pipe.Set(ctx, "access:"+jwt.Access, jwt.Access, 30* time.Second)
	pipe.Set(ctx, "refresh:"+jwt.Refresh, jwt.Refresh, 30*24*3600 * time.Second)
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		slog.Info("TokenStorage SetTokens Redis pipeline failed " + err.Error())
		return fmt.Errorf("TokenStorage SetTokens: %w", err)
	}
	for _, cmd := range cmds {
		slog.Info("Redis command executed",
			"cmd", cmd.String(),
			"args", cmd.Args(),
		)
	}
	return nil
}

func (s *TokenStorage) RemoveToken(ctx context.Context, refresh, access string) error {
	_, err := s.redisDB.Del(ctx, "access:"+access, "refresh:"+refresh).Result()

	if err != nil {
		return fmt.Errorf("TokenStorage RemoveTokens: %w", err)
	}
	return nil
}
