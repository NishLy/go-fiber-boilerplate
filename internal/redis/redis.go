package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/NishLy/go-fiber-boilerplate/internal/config"
	"github.com/NishLy/go-fiber-boilerplate/pkg/logger"
	"github.com/go-redis/redis/v8"
)

var (
	rdb  *redis.Client
	once sync.Once
)

func GetRedis() *redis.Client {
	cfg := config.Get()
	once.Do(func() {
		rdb = redis.NewClient(&redis.Options{
			Addr:     cfg.REDIS_HOST + ":" + cfg.REDIS_PORT,
			Password: "",
			DB:       0,
		})
	})
	return rdb
}

// Get tries to get the value from Redis cache. If it doesn't exist, it calls the fallback function to get the value, stores it in Redis, and returns it.
func Get[T any](key string, expiration time.Duration, fallback func() (T, error)) (T, error) {
	ctx := context.Background()
	var zero T

	val, err := GetRedis().Get(ctx, key).Result()

	logger.Log.Info(fmt.Sprintf("Cache lookup for key: %s", key))

	if err == nil {
		var data T
		if json.Unmarshal([]byte(val), &data) == nil {
			return data, nil
		}
		return zero, fmt.Errorf("failed to unmarshal cached value for key %s", key)
	}

	if err != redis.Nil {
		return zero, err
	}

	data, err := fallback()
	if err != nil {
		return zero, err
	}

	bytes, _ := json.Marshal(data)
	_ = GetRedis().Set(ctx, key, bytes, expiration).Err()

	return data, nil
}

// Delete deletes the specified key from Redis cache.
func Delete(key string) error {
	ctx := context.Background()
	return GetRedis().Del(ctx, key).Err()
}
