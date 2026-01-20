package services

import (
	"context"
	"worker-service/internal/pkg/redis"
)

type TokenBucket struct {
	Token int `json:"token"`
}

type RateLimitter struct {
	cache *redis.RedisClient[TokenBucket]
}

func NewRateLimiter(cache *redis.RedisClient[TokenBucket]) *RateLimitter {
	return &RateLimitter{
		cache: cache,
	}
}

func (r *RateLimitter) RefillToken(ctx context.Context, key string, token int) error {
	// Get from redis with specific key
	if _, err := r.cache.Get(ctx, key); err != nil {
		// If the data is empty that means we should set a new token
		errSetCache := r.cache.Set(ctx, key, TokenBucket{
			Token: token,
		})
		if errSetCache != nil {
			return errSetCache
		}
		return nil
	}

	return nil
}

func (r *RateLimitter) Allow(ctx context.Context, key string) bool {
	// Get from redis with specific key
	data, err := r.cache.Get(ctx, key)
	if err != nil {
		return false
	}

	// Update token
	if data.Token > 0 {
		data.Token--
		err := r.cache.Set(ctx, key, data)
		if err != nil {
			return false
		}
		return true
	}

	return false
}
