package domain

import (
	"time"

	"github.com/google/uuid"
)

type AdoptionStatus string

const (
	AdoptionPending   AdoptionStatus = "PENDING"
	AdoptionApproved  AdoptionStatus = "APPROVED"
	AdoptionCompleted AdoptionStatus = "COMPLETED"
)

type Adopter struct {
	ID        uuid.UUID `db:"id"`
	FullName  string    `db:"full_name"`
	Email     string    `db:"email"`
	Phone     string    `db:"phone"`
	Address   string    `db:"address"`
	CreatedAt time.Time `db:"created_at"`
}

type Adoption struct {
	ID            uuid.UUID      `db:"id"`
	AnimalID      uuid.UUID      `db:"animal_id"` // Referencia lógica al módulo de Animals
	AdopterID     uuid.UUID      `db:"adopter_id"`
	Status        AdoptionStatus `db:"status"`
	TrackingToken string         `db:"tracking_token"`
	AdoptedAt     time.Time      `db:"adopted_at"`
}

type FollowUp struct {
	ID         uuid.UUID `db:"id"`
	AdoptionID uuid.UUID `db:"adoption_id"`
	Notes      string    `db:"notes"`
	MediaURLs  []string  `db:"media_urls"` // Postgres lo maneja como ARRAY o JSONB
	Status     string    `db:"status"`     // Ej: "Saludable", "Necesita visita"
	CreatedAt  time.Time `db:"created_at"`
}
