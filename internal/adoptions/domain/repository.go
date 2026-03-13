package domain

import (
	"context"
)

type AdoptionRepository interface {
	CreateAdoption(ctx context.Context, adoption *Adoption) error
	GetAdoptionByToken(ctx context.Context, token string) (*Adoption, error)
	AddFollowUp(ctx context.Context, followUp *FollowUp) error
}
