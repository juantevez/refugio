package domain

import (
	"context"
)

type AdoptionRepository interface {
	// Adopciones
	SaveAdopter(ctx context.Context, adopter *Adopter) error
	CreateAdoption(ctx context.Context, adoption *Adoption) error
	GetAdoptionByToken(ctx context.Context, token string) (*Adoption, error)

	// Seguimiento
	AddFollowUp(ctx context.Context, followUp *FollowUp) error

	// Donaciones
	SaveDonation(ctx context.Context, donation *Donation) error
}
