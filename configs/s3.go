package configs

import (
	"context"
	storage "ocr-saas-backend/internal/storage/s3"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *storage.S3Storage

func InitS3(cfg *Config) *storage.S3Storage {
	awsCfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(cfg.AWSRegion),
	)
	if err != nil {
		panic(err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.S3Endpoint)
		o.UsePathStyle = true
	})

	S3Client = storage.NewS3Storage(client, cfg.S3Bucket)

	return S3Client
}
