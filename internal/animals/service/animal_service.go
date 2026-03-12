package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/juantevez/refugio-core/internal/animals/domain"
)

type AnimalService struct {
	repo domain.AnimalRepository
}

func NewAnimalService(r domain.AnimalRepository) *AnimalService {
	return &AnimalService{repo: r}
}

func (s *AnimalService) RegisterRescue(ctx context.Context, name string, species domain.Species, breed string) (*domain.Animal, error) {
	animal := &domain.Animal{
		ID:         uuid.New(),
		Name:       name,
		Species:    species,
		Breed:      breed,
		Status:     domain.StatusRescued,
		RescueDate: time.Now(),
		CreatedAt:  time.Now(),
	}

	if err := s.repo.Save(ctx, animal); err != nil {
		return nil, err
	}

	return animal, nil
}

func (s *AnimalService) GetAnimalDetails(ctx context.Context, id uuid.UUID) (*domain.Animal, error) {
	return s.repo.GetByID(ctx, id)
}
