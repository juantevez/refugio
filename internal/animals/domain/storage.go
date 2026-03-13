package domain

import "context"

type StorageRepository interface {
	UploadPhoto(ctx context.Context, animalID string, filename string, data []byte, contentType string) (url string, key string, err error)
	GetPresignedURL(ctx context.Context, s3Key string) (string, error)
}
