package s3

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client *s3.Client
	bucket string
}

func NewS3Storage(client *s3.Client, bucket string) *S3Storage {
	return &S3Storage{
		client: client,
		bucket: bucket,
	}
}

/*
Upload file ke S3
*/
func (s *S3Storage) Upload(
	ctx context.Context,
	key string,
	body io.Reader,
	size int64,
	contentType string,
) error {

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          body,
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(size),
	})

	return err
}

/*
Download file dari S3 ke local path (untuk OCR worker)
*/
func (s *S3Storage) Download(
	ctx context.Context,
	key string,
	dst string,
) error {

	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	defer out.Body.Close()

	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, out.Body)
	return err
}

/*
Generate presigned URL untuk akses file
Biasanya dipakai untuk frontend melihat image receipt
*/
func (s *S3Storage) GetFileURL(
	ctx context.Context,
	key string,
	expiry time.Duration,
) (string, error) {

	presignClient := s3.NewPresignClient(s.client)

	req, err := presignClient.PresignGetObject(
		ctx,
		&s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
		},
		s3.WithPresignExpires(expiry),
	)

	if err != nil {
		return "", err
	}

	return req.URL, nil
}
