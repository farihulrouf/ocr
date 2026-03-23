package configs

import (
	"context"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func ConnectRedis(cfg *Config) *redis.Client {
	addr := cfg.RedisAddr
	password := cfg.RedisPass
	dbStr := cfg.RedisDB

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
		log.Fatalf("failed connect redis: %v", err)
	}

	log.Println("✅ Redis Connected")

	return RedisClient
}
