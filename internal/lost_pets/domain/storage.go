package domain

import "context"

// ImageStorage es el puerto de salida para almacenar imágenes en S3.
// Solo sube bytes y devuelve la key — no sabe nada de EXIF ni de negocio.
type ImageStorage interface {
	Upload(ctx context.Context, folder string, filename string, data []byte, contentType string) (key string, err error)
	GetPresignedURL(ctx context.Context, s3Key string) (string, error)
}
