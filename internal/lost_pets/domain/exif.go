package domain

// ExifExtractor es el puerto de salida para extraer metadatos GPS de una imagen.
// Devuelve nil si la imagen no tiene datos EXIF de ubicación.
type ExifExtractor interface {
	ExtractLocation(data []byte) (*Point, error)
}
