package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type AdoptionStatus string

const (
	AdoptionPending   AdoptionStatus = "PENDING"
	AdoptionApproved  AdoptionStatus = "APPROVED"
	AdoptionCompleted AdoptionStatus = "COMPLETED"
	AdoptionRejected  AdoptionStatus = "REJECTED"
)

// Adoptante (La persona)
type Adopter struct {
	ID        uuid.UUID
	FullName  string
	Email     string
	Phone     string
	Address   string
	CreatedAt time.Time
}

// Adopción (La relación)
type Adoption struct {
	ID            uuid.UUID
	AnimalID      uuid.UUID // Referencia lógica al módulo de Animales
	AdopterID     uuid.UUID
	Status        AdoptionStatus
	TrackingToken string // El "Magic Link" para el seguimiento
	AdoptedAt     time.Time
}

// Seguimiento (Las actualizaciones post-adopción)
type FollowUp struct {
	ID         uuid.UUID
	AdoptionID uuid.UUID
	Notes      string
	MediaURLs  []string // Links a fotos/videos en el storage
	Status     string   // "Healthy", "Requires Visit", etc.
	CreatedAt  time.Time
}

var (
	ErrInvalidToken = errors.New("invalid or expired tracking token")
)
