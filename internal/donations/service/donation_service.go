package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/juantevez/refugio-core/internal/donations/domain"
)

// DonationRepository define la interfaz que debe implementar el dRepo en el main
type DonationRepository interface {
	SaveDonation(ctx context.Context, d *domain.Donation) error
	GetTotalDonations(ctx context.Context) (float64, error)
}

type DonationService struct {
	repo DonationRepository
}

func NewDonationService(r DonationRepository) *DonationService {
	return &DonationService{repo: r}
}

// RegisterDonation orquesta la creación de la donación con los campos de Argentina
func (s *DonationService) RegisterDonation(
	ctx context.Context,
	amount float64,
	currency string,
	source domain.DonationSource,
	ref string,
	animalID *uuid.UUID,
	donorName string,
	donorEmail string,
	details *domain.TransferDetails,
) (*domain.Donation, error) {

	if amount <= 0 {
		return nil, domain.ErrInvalidAmount
	}

	donation := &domain.Donation{
		ID:              uuid.New(),
		Amount:          amount,
		Currency:        currency,
		Source:          source,
		ReferenceNumber: ref,
		AnimalID:        animalID,
		DonorName:       donorName,
		DonorEmail:      donorEmail,
		TransferDetails: details,
		CreatedAt:       time.Now(),
	}

	if err := s.repo.SaveDonation(ctx, donation); err != nil {
		return nil, err
	}

	return donation, nil
}

func (s *DonationService) GetImpactReport(ctx context.Context) (float64, error) {
	return s.repo.GetTotalDonations(ctx)
}
