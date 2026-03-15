package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Tipos para estados y especies
type AnimalStatus string
type Species string

const (
	StatusRescued   AnimalStatus = "RESCUED"
	StatusSheltered AnimalStatus = "SHELTERED"
	StatusAdopted   AnimalStatus = "ADOPTED"

	SpeciesDog Species = "DOG"
	SpeciesCat Species = "CAT"
)

type Animal struct {
	ID         uuid.UUID    `json:"id" db:"id"`
	Name       string       `json:"name" db:"name"`
	Species    string       `db:"species"`
	Breed      string       `json:"breed" db:"breed"`
	Status     AnimalStatus `json:"status" db:"status"`
	RescueDate time.Time    `json:"rescue_date" db:"rescue_date"`
	CreatedAt  time.Time    `json:"created_at" db:"created_at"`
}

// MedicalRecord representa una entrada en la historia clínica
type MedicalRecord struct {
	ID          uuid.UUID `db:"id"`
	AnimalID    uuid.UUID `db:"animal_id"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"` 
}

// Errores de dominio
var (
	ErrAnimalNotFound = errors.New("animal not found")
	ErrInvalidStatus  = errors.New("invalid animal status transition")
)

// Regla de negocio: Un animal adoptado no puede volver a estar en rescate sin un proceso nuevo
func (a *Animal) MarkAsAdopted() error {
	if a.Status == StatusAdopted {
		return ErrInvalidStatus
	}
	a.Status = StatusAdopted
	return nil
}

type AnimalPhoto struct {
	ID         uuid.UUID `db:"id"`
	AnimalID   uuid.UUID `db:"animal_id"`
	S3URL      string    `db:"s3_url"`
	S3Key      string    `db:"s3_key"`
	PhotoOrder int16     `db:"photo_order"`
	CreatedAt  time.Time `db:"created_at"`
}

var ErrMaxPhotosReached = errors.New("el animal ya tiene 4 fotos")
