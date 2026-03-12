package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/juantevez/refugio-core/internal/donations/domain"
)

type AdoptionService struct {
	repo domain.AdoptionRepository
}

func NewAdoptionService(r domain.AdoptionRepository) *AdoptionService {
	return &AdoptionService{repo: r}
}

// StartAdoptionProcess inicia el vínculo y genera el token de seguimiento
func (s *AdoptionService) StartAdoptionProcess(ctx context.Context, animalID uuid.UUID, adopterID uuid.UUID) (*domain.Adoption, error) {

	// Generar un token seguro para el Magic Link
	tokenBytes := make([]byte, 16)
	rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)

	adoption := &domain.Adoption{
		ID:            uuid.New(),
		AnimalID:      animalID,
		AdopterID:     adopterID,
		Status:        domain.AdoptionApproved,
		TrackingToken: token,
		AdoptedAt:     time.Now(),
	}

	if err := s.repo.CreateAdoption(ctx, adoption); err != nil {
		return nil, err
	}

	return adoption, nil
}

// SubmitFollowUp permite al adoptante subir noticias sin login, solo con el token
func (s *AdoptionService) SubmitFollowUp(ctx context.Context, token string, notes string, urls []string) error {
	adoption, err := s.repo.GetAdoptionByToken(ctx, token)
	if err != nil {
		return domain.ErrInvalidToken
	}

	followUp := &domain.FollowUp{
		ID:         uuid.New(),
		AdoptionID: adoption.ID,
		Notes:      notes,
		MediaURLs:  urls,
		CreatedAt:  time.Now(),
	}

	return s.repo.AddFollowUp(ctx, followUp)
}
