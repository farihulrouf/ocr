package configs

import (
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

	RedisAddr string
	RedisPass string
	RedisDB   string

	AWSRegion  string
	S3Bucket   string
	S3Endpoint string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env not found")
	}

	cfg := &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBPort:     os.Getenv("DB_PORT"),

		RedisAddr: os.Getenv("REDIS_ADDR"),
		RedisPass: os.Getenv("REDIS_PASSWORD"),
		RedisDB:   os.Getenv("REDIS_DB"),

		AWSRegion:  os.Getenv("AWS_REGION"),
		S3Bucket:   os.Getenv("S3_BUCKET"),
		S3Endpoint: os.Getenv("S3_ENDPOINT"),
	}

	// ✅ Validasi sederhana (penting banget)
	if cfg.DBHost == "" || cfg.DBUser == "" {
		log.Fatal("Database config missing")
	}

	return cfg
}
