package infrastructure

import (
	"context"
	"fmt"
	"time"
	"worker-service/config"

	"github.com/redis/go-redis/v9"
)

func InitializeRedisConnection(cfg config.AppConfig) *redis.Client {
	opt, _ := redis.ParseURL(cfg.Redis.ServerAddress)
	opt = &redis.Options{
		ReadTimeout:     3 * time.Second,
		WriteTimeout:    3 * time.Second,
		PoolTimeout:     4 * time.Second,
		MaxRetries:      10,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,
	}
	db := redis.NewClient(opt)

	status := db.Ping(context.TODO())
	if status.Err() != nil {
		panic(fmt.Sprintf("error connecting to redis: %s", status.Err()))
	}

	return db
}
