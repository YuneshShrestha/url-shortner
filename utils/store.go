package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

// Create Redis connection
func NewRedisClient() *redis.Client {
	fmt.Println("Connecting to redis server on:", os.Getenv("REDIS_HOST"))

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return rdb
}

func NextID(ctx context.Context, rdb *redis.Client) (uint64, error) {
	id, err := rdb.Incr(ctx, "url-shortener:next-id").Result()
	if err != nil {
		return 0, fmt.Errorf("failed to generate URL ID: %w", err)
	}
	return uint64(id), nil
}

func SetLongURL(ctx context.Context, rdb *redis.Client, id uint64, value string) error {
	if err := rdb.Set(ctx, fmt.Sprintf("url:%d", id), value, 0).Err(); err != nil {
		return fmt.Errorf("failed to save URL: %w", err)
	}
	return nil
}

// Get original URL using short URL
func GetLongURL(
	ctx context.Context,
	rdb *redis.Client,
	shortURL string,
) (string, error) {
	id, err := DecodeID(shortURL)
	if err != nil {
		return "", fmt.Errorf("invalid short code: %w", err)
	}

	longURL, err := rdb.Get(ctx, fmt.Sprintf("url:%d", id)).Result()

	if err == redis.Nil {
		return "", fmt.Errorf("short URL not found")
	}

	if err != nil {
		return "", fmt.Errorf("failed to retrieve from Redis: %v", err)
	}

	return longURL, nil
}
