package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSLMode  string
	DBTimeZone string

	RedisAddr string
	RedisPass string
	RedisDB   string

	AWSRegion  string
	S3Bucket   string
	S3Endpoint string
}

func LoadConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("❌ .env NOT FOUND:", err)
	} else {
		log.Println("✅ .env loaded")
	}

	cfg := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBPort:     os.Getenv("DB_PORT"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
		DBTimeZone: os.Getenv("DB_TIMEZONE"),

		RedisAddr: os.Getenv("REDIS_ADDR"),
		RedisPass: os.Getenv("REDIS_PASSWORD"),
		RedisDB:   os.Getenv("REDIS_DB"),

		AWSRegion:  os.Getenv("AWS_REGION"),
		S3Bucket:   os.Getenv("S3_BUCKET"),
		S3Endpoint: os.Getenv("S3_ENDPOINT"),
	}

	// 🔥 DEBUG (AMAN – tanpa password)
	fmt.Println("========== DEBUG CONFIG ==========")
	fmt.Println("DB_HOST     :", cfg.DBHost)
	fmt.Println("DB_PORT     :", cfg.DBPort)
	fmt.Println("DB_USER     :", cfg.DBUser)
	fmt.Println("DB_NAME     :", cfg.DBName)
	fmt.Println("DB_SSLMODE  :", cfg.DBSSLMode)
	fmt.Println("REDIS_ADDR  :", cfg.RedisAddr)
	fmt.Println("S3_ENDPOINT :", cfg.S3Endpoint)
	fmt.Println("==================================")

	if cfg.DBHost == "" || cfg.DBUser == "" {
		log.Fatal("❌ Database config missing")
	}

	return cfg
}
