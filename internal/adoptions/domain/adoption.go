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
	ID            uuid.UUID      `db:"id"`
	AnimalID      uuid.UUID      `db:"animal_id"`
	AdopterID     uuid.UUID      `db:"adopter_id"`
	Status        AdoptionStatus `db:"status"`
	TrackingToken string         `db:"tracking_token"`
	AdoptedAt     time.Time      `db:"adopted_at"`
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
