package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/thienkb1123/go-clean-arch/config"
	//"github.com/thienkb1123/go-clean-arch/config"
)

type RedisClient struct {
	rdbClient *redis.Client
}

// Returns new redis client
func NewRedisClient(cfg *config.RedisClient) (*RedisClient, error) {
	redisHost := cfg.RedisAddr

	if redisHost == "" {
		redisHost = ":6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:         redisHost,
		MinIdleConns: cfg.MinIdleConns,
		PoolSize:     cfg.PoolSize,
		PoolTimeout:  time.Duration(cfg.PoolTimeout) * time.Second,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	return &RedisClient{
		rdbClient: client,
	}, nil
}

func (r *RedisClient) Get(ctx context.Context, key string) ([]byte, error) {
	res, err := r.rdbClient.Get(ctx, key).Bytes()
	if err == redis.Nil {
		fmt.Println("nil")
		availableKey, err := r.rdbClient.Keys(ctx, key + "*").Result()
		if err != nil {
			return res, err
		}
		if len(availableKey) == 0{
			return res, redis.Nil
		}
		return r.rdbClient.Get(ctx, availableKey[0]).Bytes()
	}
	return res, err
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.rdbClient.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Incr(ctx context.Context, key string) (int64, error) {
	return r.rdbClient.Incr(ctx, key).Result()
}

func (r *RedisClient) Del(ctx context.Context, keys ...string) error {
	return r.rdbClient.Del(ctx, keys...).Err()
}

func (r *RedisClient) Close() error {
	return r.rdbClient.Close()
}

func (r *RedisClient) Ping(ctx context.Context) error {
	return r.rdbClient.Ping(ctx).Err()
}
