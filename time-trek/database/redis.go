package database

import (
	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

// InitRedis initializes the Redis connection
func InitRedis() *redis.Client {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Update with your Redis server details
		Password: "",               // No password by default
		DB:       0,                // Default DB
	})

	return redisClient
}
