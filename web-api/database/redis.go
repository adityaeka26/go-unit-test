package database

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis interface {
	GetClient() RedisClient
}
type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}
type RedisImpl struct {
	redisClient *redis.Client
}

func NewRedis(addr string, password string) Redis {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})
	return &RedisImpl{
		redisClient: redisClient,
	}
}

func (redis RedisImpl) GetClient() RedisClient {
	return redis.redisClient
}
