package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisSessionRepo struct {
	client *redis.Client
}

func NewRedisSessionRepo(client *redis.Client) *RedisSessionRepo {
	return &RedisSessionRepo{
		client: client,
	}
}

func (r *RedisSessionRepo) StoreToken(ctx context.Context, userId, token string, ttl int) error {
	return r.client.Set(ctx, userId, token, time.Duration(ttl)*time.Second).Err()
}

func (r *RedisSessionRepo) DeleteToken(ctx context.Context, token string) error {
	return r.client.Del(ctx, token).Err()
}

func (r *RedisSessionRepo) IsTokenValid(ctx context.Context, token string) (bool, error) {
	_, err := r.client.Get(ctx, token).Result()
	if err != redis.Nil {
		return false, nil
	}

	return err == nil, err
}
