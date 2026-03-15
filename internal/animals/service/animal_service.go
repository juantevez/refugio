package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/juantevez/refugio-core/internal/animals/domain"
	animalsDomain "github.com/juantevez/refugio-core/internal/animals/domain"
)

type AnimalService struct {
	repo    animalsDomain.AnimalRepository
	storage domain.StorageRepository
}

func NewAnimalService(r domain.AnimalRepository, s domain.StorageRepository) *AnimalService {
	return &AnimalService{repo: r, storage: s}
}

// Agregamos rescueDate time.Time como parámetro
func (s *AnimalService) RegisterRescue(ctx context.Context, name string, species string, breed string, rescueDate time.Time) (*animalsDomain.Animal, error) {
	animal := &animalsDomain.Animal{
		ID:         uuid.New(),
		Name:       name,
		Species:    species,
		Breed:      breed,
		Status:     animalsDomain.StatusRescued,
		RescueDate: rescueDate, 
		CreatedAt:  time.Now(),
	}

	if err := s.repo.Save(ctx, animal); err != nil {
		return nil, err
	}

	return animal, nil
}

func (s *AnimalService) GetAnimalDetails(ctx context.Context, id uuid.UUID) (*animalsDomain.Animal, error) {
	return s.repo.GetByID(ctx, id)
}

// List retorna la lista de animales. El segundo parámetro son los filtros.
func (s *AnimalService) List(ctx context.Context, filters map[string]interface{}) ([]domain.Animal, error) {
	// Acá podrías poner lógica de negocio, por ejemplo:
	// "Solo mostrar animales que no estén adoptados"
	return s.repo.List(ctx, filters)
}

func (s *AnimalService) UploadPhoto(ctx context.Context, animalID uuid.UUID, filename string, data []byte, contentType string) (*domain.AnimalPhoto, error) {
	// Validar que el animal existe
	_, err := s.repo.GetByID(ctx, animalID)
	if err != nil {
		return nil, domain.ErrAnimalNotFound
	}

	// Validar máximo 4 fotos
	existing, err := s.repo.GetPhotosByAnimalID(ctx, animalID)
	if err != nil {
		return nil, err
	}
	if len(existing) >= 4 {
		return nil, domain.ErrMaxPhotosReached
	}

	// Subir a S3
	_, key, err := s.storage.UploadPhoto(ctx, animalID.String(), filename, data, contentType)
	if err != nil {
		return nil, err
	}

	photo := &domain.AnimalPhoto{
		ID:         uuid.New(),
		AnimalID:   animalID,
		S3Key:      key,
		PhotoOrder: int16(len(existing)),
		CreatedAt:  time.Now(),
	}

	if err := s.repo.SavePhoto(ctx, photo); err != nil {
		return nil, err
	}

	return photo, nil
}

func (s *AnimalService) GetPhotos(ctx context.Context, animalID uuid.UUID) ([]string, error) {
	photos, err := s.repo.GetPhotosByAnimalID(ctx, animalID)
	if err != nil {
		return nil, err
	}

	urls := make([]string, 0, len(photos))
	for _, p := range photos {
		url, err := s.storage.GetPresignedURL(ctx, p.S3Key)
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}

	return urls, nil
}
