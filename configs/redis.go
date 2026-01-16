package configs

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func ConnectRedis() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	if err := RedisClient.Ping(Ctx).Err(); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	log.Println("✅ Redis connected")
}
