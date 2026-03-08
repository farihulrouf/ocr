package storage

import (
	"context"
	"io"
	"log"
	"net/url"
	"ocr-saas-backend/configs"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	Client *minio.Client
	Bucket string
}

// NewMinioStorage buat client MinIO
func NewMinioStorage(cfg *configs.MinioConfigStruct) *MinioStorage {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		log.Fatalln("MinIO init error:", err)
	}

	// pastikan bucket ada
	exists, err := client.BucketExists(context.Background(), cfg.Bucket)
	if err != nil {
		log.Fatalln("MinIO bucket check error:", err)
	}
	if !exists {
		err = client.MakeBucket(context.Background(), cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalln("MinIO create bucket error:", err)
		}
	}

	return &MinioStorage{
		Client: client,
		Bucket: cfg.Bucket,
	}
}

// Upload file ke MinIO
func (m *MinioStorage) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := m.Client.PutObject(ctx, m.Bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// GetFileURL menghasilkan URL sementara yang bisa diakses publik
func (m *MinioStorage) GetFileURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)
	// Generate URL yang berlaku selama 'expiry'
	presignedURL, err := m.Client.PresignedGetObject(ctx, m.Bucket, objectName, expiry, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
