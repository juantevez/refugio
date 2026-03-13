package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/juantevez/refugio-core/internal/animals/domain"
	animalsDomain "github.com/juantevez/refugio-core/internal/animals/domain"
)

type AnimalService struct {
	repo animalsDomain.AnimalRepository
}

func NewAnimalService(r animalsDomain.AnimalRepository) *AnimalService {
	return &AnimalService{repo: r}
}

// Agregamos rescueDate time.Time como parámetro
func (s *AnimalService) RegisterRescue(ctx context.Context, name string, species string, breed string, rescueDate time.Time) (*animalsDomain.Animal, error) {
	animal := &animalsDomain.Animal{
		ID:         uuid.New(),
		Name:       name,
		Species:    species,
		Breed:      breed,
		Status:     animalsDomain.StatusRescued,
		RescueDate: rescueDate, // <--- AHORA SÍ SE ASIGNA EL DATO DEL REQUEST
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
