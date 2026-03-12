package domain

import (
	"context"

	"github.com/google/uuid"
)

type DonationRepository interface {
	// ... (métodos de adopción previos)

	// Donaciones
	SaveDonation(ctx context.Context, donation *Donation) error
	GetDonationsByAnimal(ctx context.Context, animalID uuid.UUID) ([]Donation, error)
	GetTotalDonations(ctx context.Context) (float64, error)
}
