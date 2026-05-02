package s3

import (
	"context"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type S3Storage struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	bucket        string
}

func NewS3Storage(client *s3.Client, bucket string) *S3Storage {
	return &S3Storage{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		bucket:        bucket,
	}
}

/*
Helper: normalize URL -> key
support:
- receipts/123.jpg
- https://xxx/receipts/123.jpg
*/
func normalizeKey(input string) string {
	if strings.HasPrefix(input, "http") {
		u, err := url.Parse(input)
		if err == nil {
			return strings.TrimPrefix(u.Path, "/")
		}
	}
	return input
}

/*
Upload (low level)
*/
func (s *S3Storage) Upload(
	ctx context.Context,
	key string,
	body io.Reader,
	size int64,
	contentType string,
) error {

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

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
Upload file langsung dari path
*/
func (s *S3Storage) UploadFile(
	ctx context.Context,
	key string,
	filePath string,
) error {

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	return s.Upload(
		ctx,
		key,
		file,
		stat.Size(),
		"application/octet-stream",
	)
}

/*
Download file ke local
*/
func (s *S3Storage) Download(
	ctx context.Context,
	key string,
	dst string,
) error {

	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	key = normalizeKey(key)

	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	defer out.Body.Close()

	if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, out.Body)
	return err
}

/*
🔥 VERSI BARU (RECOMMENDED)
Download ke temp + cleanup
*/
func (s *S3Storage) DownloadToTempWithCleanup(
	ctx context.Context,
	key string,
) (string, func(), error) {

	key = normalizeKey(key)

	tmpFile := filepath.Join(
		os.TempDir(),
		uuid.New().String()+filepath.Ext(key),
	)

	err := s.Download(ctx, key, tmpFile)
	if err != nil {
		return "", nil, err
	}

	cleanup := func() {
		_ = os.Remove(tmpFile)
	}

	return tmpFile, cleanup, nil
}

/*
🔥 VERSI LAMA (BIAR KODE KAMU NGGAK ERROR)
*/
func (s *S3Storage) DownloadToTemp(
	ctx context.Context,
	key string,
) (string, error) {

	path, _, err := s.DownloadToTempWithCleanup(ctx, key)
	return path, err
}

/*
Generate presigned URL
*/
func (s *S3Storage) GetFileURL(
	ctx context.Context,
	key string,
	expiry time.Duration,
) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	key = normalizeKey(key)

	req, err := s.presignClient.PresignGetObject(
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

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	key = normalizeKey(key)

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	return err
}
