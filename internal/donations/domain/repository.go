package domain

import (
	"context"

	"github.com/google/uuid"
)

type DonationRepository interface {
	SaveDonation(ctx context.Context, donation *Donation) error
	GetDonationsByAnimal(ctx context.Context, animalID uuid.UUID) ([]Donation, error)
	GetTotalDonations(ctx context.Context) (float64, error)
	// AddFollowUp NO debe estar aquí, debe estar en adoptions/domain/repository.go
}
