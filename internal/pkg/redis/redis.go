package redis

import (
	"context"
	"encoding/json"
	"time"
	"worker-service/config"
	"worker-service/infrastructure"

	"github.com/redis/go-redis/v9"
)

type RedisClient[T any] struct {
	client *redis.Client
	key    string
	TTL    time.Duration
}

func NewRedisClient[T any](cfg config.AppConfig, key string, TTL time.Duration) *RedisClient[T] {
	return &RedisClient[T]{
		client: infrastructure.InitializeRedisConnection(cfg),
		key:    key,
		TTL:    TTL,
	}
}

func (r *RedisClient[T]) Enqueue(ctx context.Context, task T) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return r.client.RPush(ctx, r.key, data).Err()
}

func (r *RedisClient[T]) Dequeue(ctx context.Context) (T, error) {
	var task T
	res, err := r.client.BRPop(ctx, 0, r.key).Result()
	if err != nil {
		return task, err
	}
	if err := json.Unmarshal([]byte(res[1]), &task); err != nil {
		return task, err
	}
	return task, nil
}

func (r *RedisClient[T]) Set(ctx context.Context, suffixKey string, task T) error {
	var key string = r.key
	if suffixKey != "" {
		key = r.key + ":" + suffixKey
	}
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, r.TTL).Err()
}

func (r *RedisClient[T]) Get(ctx context.Context, suffixKey string) (T, error) {
	var task T
	var key string = r.key
	if suffixKey != "" {
		key = r.key + ":" + suffixKey
	}
	res, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return task, err
	}
	if err := json.Unmarshal([]byte(res), &task); err != nil {
		return task, err
	}
	return task, nil
}
