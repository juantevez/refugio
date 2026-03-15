package storage

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Adapter struct {
	client *s3.Client
	bucket string
}

func NewS3Adapter() (*S3Adapter, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("error configurando AWS: %w", err)
	}

	return &S3Adapter{
		client: s3.NewFromConfig(cfg),
		bucket: os.Getenv("AWS_S3_BUCKET"),
	}, nil
}

// Upload sube una imagen a S3 y devuelve la key.
// La key tiene el formato: {folder}/{timestamp}_{filename}
// Ej: lost_pets/1718000000_perro.jpg
func (a *S3Adapter) Upload(ctx context.Context, folder string, filename string, data []byte, contentType string) (string, error) {
	key := fmt.Sprintf("%s/%d_%s", folder, time.Now().Unix(), filename)

	_, err := a.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(a.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("error subiendo imagen a S3: %w", err)
	}

	return key, nil
}

// GetPresignedURL genera una URL temporal para acceder a la imagen (TTL 7 días).
func (a *S3Adapter) GetPresignedURL(ctx context.Context, s3Key string) (string, error) {
	presigner := s3.NewPresignClient(a.client)
	presigned, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(a.bucket),
		Key:    aws.String(s3Key),
	}, s3.WithPresignExpires(7*24*time.Hour))
	if err != nil {
		return "", fmt.Errorf("error generando pre-signed URL: %w", err)
	}
	return presigned.URL, nil
}
