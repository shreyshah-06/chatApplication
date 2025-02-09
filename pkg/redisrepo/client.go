package redisrepo

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

// InitialiseRedis initializes the Redis connection
func InitialiseRedis() *redis.Client {
	conn := redis.NewClient(&redis.Options{
		Addr:        os.Getenv("REDIS_CONNECTION_STRING"),
		// Password:    os.Getenv("REDIS_PASSWORD"),
		DB:          0,
		DialTimeout: 5 * time.Second,
	})

	// Check if Redis is connected
	if pong, err := conn.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Redis Connection Failed: %v", err)
	} else {
		log.Println("Redis Successfully Connected. Ping:", pong)
	}

	redisClient = conn
	return redisClient
}