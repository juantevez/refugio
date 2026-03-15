package exif

import (
	"bytes"
	"fmt"

	goexif "github.com/rwcarlsen/goexif/exif"
	"github.com/juantevez/refugio-core/internal/lost_pets/domain"
)

type Extractor struct{}

func NewExtractor() *Extractor {
	return &Extractor{}
}

// ExtractLocation lee los metadatos GPS de los bytes de una imagen.
// Devuelve nil, nil si la imagen no tiene datos EXIF de ubicación — no es un error.
func (e *Extractor) ExtractLocation(data []byte) (*domain.Point, error) {
	x, err := goexif.Decode(bytes.NewReader(data))
	if err != nil {
		// Sin EXIF en la imagen — caso normal, no es un error de negocio
		return nil, nil
	}

	lat, long, err := x.LatLong()
	if err != nil {
		// Tiene EXIF pero sin coordenadas GPS — también caso normal
		return nil, nil
	}

	if lat == 0 && long == 0 {
		return nil, nil
	}

	if lat < -90 || lat > 90 || long < -180 || long > 180 {
		return nil, fmt.Errorf("coordenadas EXIF fuera de rango: lat=%.6f long=%.6f", lat, long)
	}

	return &domain.Point{Lat: lat, Long: long}, nil
}
