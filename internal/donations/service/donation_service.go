package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/juantevez/refugio-core/internal/donations/domain"
)

type DonationService struct {
	repo domain.DonationRepository
}

func NewDonationService(r domain.DonationRepository) *DonationService {
	return &DonationService{repo: r}
}

// RegisterDonation procesa y guarda una nueva entrada de dinero
func (s *DonationService) RegisterDonation(ctx context.Context, amount float64, currency string, source domain.DonationSource, ref string, animalID *uuid.UUID) (*domain.Donation, error) {

	// Regla de negocio básica: no aceptamos montos negativos o cero
	if amount <= 0 {
		return nil, domain.ErrInvalidAmount // Deberías definir este error en domain
	}

	donation := &domain.Donation{
		ID:              uuid.New(),
		Amount:          amount,
		Currency:        currency,
		Source:          source,
		ReferenceNumber: ref,
		AnimalID:        animalID,
		CreatedAt:       time.Now(),
	}

	if err := s.repo.SaveDonation(ctx, donation); err != nil {
		return nil, err
	}

	// Aquí podrías disparar un evento: "DonationReceived"
	// para que el módulo de notificaciones envíe un agradecimiento.

	return donation, nil
}

func (s *DonationService) GetImpactReport(ctx context.Context) (float64, error) {
	return s.repo.GetTotalDonations(ctx)
}
