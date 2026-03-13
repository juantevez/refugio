package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	adoptionsDomain "github.com/juantevez/refugio-core/internal/adoptions/domain"
)

// Antes usaba donationsDomain.DonationRepository, ahora:
type AdoptionService struct {
	repo adoptionsDomain.AdoptionRepository
}

func NewAdoptionService(repo adoptionsDomain.AdoptionRepository) *AdoptionService {
	return &AdoptionService{repo: repo}
}

// StartAdoptionProcess inicia el trámite y genera el token de seguimiento
func (s *AdoptionService) StartAdoptionProcess(ctx context.Context, animalID, adopterID uuid.UUID) (*adoptionsDomain.Adoption, error) {
	token, _ := generateRandomToken(16)

	adoption := &adoptionsDomain.Adoption{
		ID:            uuid.New(),
		AnimalID:      animalID,
		AdopterID:     adopterID,
		Status:        adoptionsDomain.AdoptionPending,
		TrackingToken: token,
		AdoptedAt:     time.Now(),
	}

	if err := s.repo.CreateAdoption(ctx, adoption); err != nil {
		return nil, err
	}

	return adoption, nil
}

// SubmitFollowUp permite al adoptante subir actualizaciones
func (s *AdoptionService) SubmitFollowUp(ctx context.Context, token, notes string, media []string) error {
	// 1. Buscar la adopción por token
	adoption, err := s.repo.GetAdoptionByToken(ctx, token)
	if err != nil {
		return adoptionsDomain.ErrInvalidToken
	}

	// 2. Crear el registro de seguimiento
	followUp := &adoptionsDomain.FollowUp{
		ID:         uuid.New(),
		AdoptionID: adoption.ID,
		Notes:      notes,
		MediaURLs:  media,
		CreatedAt:  time.Now(),
	}

	return s.repo.AddFollowUp(ctx, followUp)
}

func generateRandomToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
