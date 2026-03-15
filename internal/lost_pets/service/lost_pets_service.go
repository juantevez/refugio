package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/juantevez/refugio-core/internal/lost_pets/domain"
)

type LostPetsService struct {
	repo    domain.PetReportRepository
	storage domain.ImageStorage
	exif    domain.ExifExtractor
}

func NewLostPetsService(r domain.PetReportRepository, s domain.ImageStorage, e domain.ExifExtractor) *LostPetsService {
	return &LostPetsService{repo: r, storage: s, exif: e}
}

// CreateReport crea un nuevo reporte de mascota perdida o encontrada.
// Si la imagen tiene GPS en EXIF, usa esa ubicación (mayor confianza).
// Si no, usa las coordenadas enviadas por el cliente (GPS del celular).
func (s *LostPetsService) CreateReport(
	ctx context.Context,
	reportType domain.ReportType,
	species domain.PetSpecies,
	description string,
	contactName string,
	contactEmail string,
	contactPhone string,
	clientLocation *domain.Point,
	imageData []byte,
	imageFilename string,
	imageContentType string,
) (*domain.PetReport, error) {

	if reportType != domain.ReportTypeLost && reportType != domain.ReportTypeFound {
		return nil, domain.ErrInvalidReportType
	}

	// Determinar ubicación y su fuente de confianza
	location, locationSource, err := s.resolveLocation(imageData, clientLocation)
	if err != nil {
		return nil, err
	}

	// Subir imagen a S3 (opcional: el reporte puede no tener foto)
	var photoKey string
	if len(imageData) > 0 {
		photoKey, err = s.storage.Upload(ctx, "lost_pets", imageFilename, imageData, imageContentType)
		if err != nil {
			return nil, err
		}
	}

	report := &domain.PetReport{
		ID:             uuid.New(),
		Type:           reportType,
		Species:        species,
		Description:    description,
		PhotoS3Key:     photoKey,
		Location:       *location,
		LocationSource: locationSource,
		RadiusMeters:   500, // radio de búsqueda por defecto
		Status:         domain.ReportStatusActive,
		ContactName:    contactName,
		ContactEmail:   contactEmail,
		ContactPhone:   contactPhone,
		ReportedAt:     time.Now(),
		CreatedAt:      time.Now(),
	}

	if err := s.repo.Save(ctx, report); err != nil {
		return nil, err
	}

	return report, nil
}

// SearchNearby devuelve reportes activos dentro del área de búsqueda.
func (s *LostPetsService) SearchNearby(ctx context.Context, area domain.SearchArea) ([]domain.PetReport, error) {
	if area.Center.Lat == 0 && area.Center.Long == 0 {
		return nil, domain.ErrInvalidLocation
	}
	return s.repo.SearchNearby(ctx, area)
}

// GetReport devuelve un reporte por ID.
func (s *LostPetsService) GetReport(ctx context.Context, id uuid.UUID) (*domain.PetReport, error) {
	return s.repo.GetByID(ctx, id)
}

// ResolveReport marca un reporte como resuelto (mascota encontrada/reunida con dueño).
func (s *LostPetsService) ResolveReport(ctx context.Context, id uuid.UUID) error {
	report, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.ErrReportNotFound
	}
	if report.Status == domain.ReportStatusResolved {
		return domain.ErrReportAlreadyResolved
	}
	return s.repo.UpdateStatus(ctx, id, domain.ReportStatusResolved)
}

// resolveLocation decide qué ubicación usar y su nivel de confianza.
// Jerarquía: EXIF (mayor confianza) > GPS del cliente > error si ninguna.
func (s *LostPetsService) resolveLocation(imageData []byte, clientLocation *domain.Point) (*domain.Point, domain.LocationSource, error) {
	// Intentar extraer GPS del EXIF de la imagen
	if len(imageData) > 0 && s.exif != nil {
		if exifPoint, err := s.exif.ExtractLocation(imageData); err == nil && exifPoint != nil {
			return exifPoint, domain.LocationSourceEXIF, nil
		}
	}

	// Fallback: usar ubicación del cliente (GPS del celular/navegador)
	if clientLocation != nil && (clientLocation.Lat != 0 || clientLocation.Long != 0) {
		return clientLocation, domain.LocationSourceGPS, nil
	}

	return nil, "", domain.ErrInvalidLocation
}
