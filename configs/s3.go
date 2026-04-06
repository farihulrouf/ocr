package configs

import (
	"context"
	"log" // Tambahkan import log
	storage "ocr-saas-backend/internal/storage/s3"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *storage.S3Storage

func InitS3(cfg *Config) *storage.S3Storage {
	region := cfg.AWSRegion
	if region == "" {
		region = "us-east-1" // Default untuk Localstack
	}

	// Tangkap err di sini
	awsCfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
	)

	// ✅ GUNAKAN err-nya (Cek apakah ada error saat load config AWS)
	if err != nil {
		log.Fatalf("❌ Failed to load AWS SDK config: %v", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.S3Endpoint)
		o.UsePathStyle = true
	})

	S3Client = storage.NewS3Storage(client, cfg.S3Bucket)

	return S3Client
}
