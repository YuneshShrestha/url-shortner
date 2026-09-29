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

// Save short URL → original URL
func SetKey(
	ctx *context.Context,
	rdb *redis.Client,
	key string,
	value string,
	ttl int,
) {
	fmt.Println("Setting key", key, "to", value, "in Redis")

	rdb.Set(*ctx, key, value, 0)

	fmt.Println("The key", key, "has been set to", value, "successfully")
}

// Get original URL using short URL
func GetLongURL(
	ctx *context.Context,
	rdb *redis.Client,
	shortURL string,
) (string, error) {

	longURL, err := rdb.Get(*ctx, shortURL).Result()

	if err == redis.Nil {
		return "", fmt.Errorf("short URL not found")
	}

	if err != nil {
		return "", fmt.Errorf("failed to retrieve from Redis: %v", err)
	}

	return longURL, nil
}
