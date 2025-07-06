package storage

import (
	"context"
	"fmt"
	"go-auth/internal/models"

	"github.com/redis/go-redis/v9"
)


type TokenStorage struct {
		redisDB *redis.Client

}

func NewTokenStorage(redisDB *redis.Client) *TokenStorage {
return &TokenStorage{redisDB: redisDB}}

func (s *TokenStorage) SetTokens(ctx context.Context, jwt *models.TokenDto) error {
	pipe := s.redisDB.Pipeline()
	pipe.Set(ctx, "access:"+jwt.Access, jwt.Access, 3600)
	pipe.Set(ctx, "refresh:"+jwt.Refresh, jwt.Refresh, 30*24*3600)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("TokenStorage SetTokens: %w", err )
	}
	return nil
}

func (s *TokenStorage) RemoveToken(ctx context.Context, refresh, access string) error {
	_,err:=s.redisDB.Del(ctx, "access:"+access, "refresh:"+refresh).Result()
	
	if err != nil {
		return fmt.Errorf("TokenStorage RemoveTokens: %w", err )
	}
	return nil
}
