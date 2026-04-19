package handle_uploads

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	minioOnce   sync.Once
	minioClient *minio.Client
	minioErr    error
)

func parseBoolEnv(value string, fallback bool) bool {
	v := strings.TrimSpace(value)
	if v == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return parsed
}

func initMinIO(endpoint string, accessKey string, secretKey string, useSSL bool) (*minio.Client, error) {
	minioOnce.Do(func() {
		minioClient, minioErr = minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
			Secure: useSSL,
		})
	})
	return minioClient, minioErr
}

func SaveImage(file *multipart.FileHeader) (string, error) {
	endpoint := strings.TrimSpace(os.Getenv("MINIO_ENDPOINT"))
	accessKey := strings.TrimSpace(os.Getenv("MINIO_ACCESS_KEY"))
	secretKey := strings.TrimSpace(os.Getenv("MINIO_SECRET_KEY"))
	bucket := strings.TrimSpace(os.Getenv("MINIO_BUCKET"))
	publicBaseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("MINIO_PUBLIC_BASE_URL")), "/")
	useSSL := parseBoolEnv(os.Getenv("MINIO_USE_SSL"), false)

	if endpoint == "" || accessKey == "" || secretKey == "" || bucket == "" || publicBaseURL == "" {
		return "", fmt.Errorf("minio storage is not configured")
	}

	client, err := initMinIO(endpoint, accessKey, secretKey, useSSL)
	if err != nil {
		return "", fmt.Errorf("failed to initialize minio client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fileReader, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open image file: %w", err)
	}
	defer fileReader.Close()

	objectName := fmt.Sprintf("devices/%d_%s", time.Now().Unix(), filepath.Base(file.Filename))

	_, err = client.PutObject(ctx, bucket, objectName, fileReader, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to minio: %w", err)
	}

	return fmt.Sprintf("%s/%s/%s", publicBaseURL, bucket, objectName), nil
}
