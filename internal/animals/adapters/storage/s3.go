package storage

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Repository struct {
	client *s3.Client
	bucket string
}

func NewS3Repository() (*S3Repository, error) {
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

	return &S3Repository{
		client: s3.NewFromConfig(cfg),
		bucket: os.Getenv("AWS_S3_BUCKET"),
	}, nil
}

func (r *S3Repository) UploadPhoto(ctx context.Context, animalID string, filename string, data []byte, contentType string) (string, string, error) {
	key := fmt.Sprintf("animals/%s/%d_%s", animalID, time.Now().Unix(), filename)

	// LOG: verificar qué parámetros está usando
	log.Printf("S3 DEBUG — bucket: %s | region: %s | key: %s | contentType: %s | size: %d bytes",
		r.bucket,
		os.Getenv("AWS_REGION"),
		key,
		contentType,
		len(data),
	)

	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Printf("S3 DEBUG — PutObject error: %v", err)
		return "", "", fmt.Errorf("error subiendo foto a S3: %w", err)
	}

	return "", key, nil
}

func (r *S3Repository) GetPresignedURL(ctx context.Context, s3Key string) (string, error) {
	presigner := s3.NewPresignClient(r.client)
	presigned, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(s3Key),
	}, s3.WithPresignExpires(7*24*time.Hour))
	if err != nil {
		return "", fmt.Errorf("error generando pre-signed URL: %w", err)
	}
	return presigned.URL, nil
}
