package configs

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func ConnectRedis() {

	addr := os.Getenv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	db := 0
	if dbStr != "" {
		v, err := strconv.Atoi(dbStr)
		if err == nil {
			db = v
		}
	}

	if addr == "" {
		addr = "localhost:6379"
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := RedisClient.Ping(Ctx).Err(); err != nil {
		log.Fatalf("failed to connect redis: %v", err)
	}

	log.Println("✅ Redis connected")
}
