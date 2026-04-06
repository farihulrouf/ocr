package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// App Config
	AppPort    string
	JWTSecret  string
	MistralKey string

	// Database Config
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSLMode  string
	DBTimeZone string // <--- Penting agar ConnectDB tidak error

	// Redis Config
	RedisAddr string
	RedisPass string
	RedisDB   string

	// Storage Config (INI YANG BIKIN ERROR TADI)
	AWSRegion  string // <--- Harus ada ini!
	S3Bucket   string
	S3Endpoint string
}

func LoadConfig() *Config {
	// Cek file .env
	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load(".env")
		log.Println("✅ .env file loaded")
	} else {
		log.Println("ℹ️  Using System Environment (Docker Mode)")
	}

	cfg := &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		JWTSecret:  getEnv("JWT_SECRET", "super-secret-key"),
		MistralKey: os.Getenv("MISTRAL_API_KEY"),

		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		DBTimeZone: getEnv("DB_TIMEZONE", "Asia/Jakarta"),

		RedisAddr: os.Getenv("REDIS_ADDR"),
		RedisPass: os.Getenv("REDIS_PASSWORD"),
		RedisDB:   getEnv("REDIS_DB", "0"),

		// Mapping Field untuk S3
		AWSRegion:  getEnv("AWS_REGION", "us-east-1"),
		S3Bucket:   os.Getenv("S3_BUCKET"),
		S3Endpoint: os.Getenv("S3_ENDPOINT"),
	}

	fmt.Println("\n========== 🛠️  CONFIG LOADED ==========")
	fmt.Printf("🌐 PORT      : %s\n", cfg.AppPort)
	fmt.Printf("🗄️  DB HOST   : %s\n", cfg.DBHost)
	fmt.Printf("☁️  S3 REGION : %s\n", cfg.AWSRegion)
	fmt.Println("======================================\n")

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
