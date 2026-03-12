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

// Animal representa la entidad principal del refugio
type Animal struct {
	ID             uuid.UUID
	Name           string
	Species        Species
	Breed          string
	Status         AnimalStatus
	RescueDate     time.Time
	MedicalHistory []MedicalRecord
	CreatedAt      time.Time
}

// MedicalRecord representa una entrada en la historia clínica
type MedicalRecord struct {
	ID          uuid.UUID
	AnimalID    uuid.UUID
	Description string
	Treatment   string
	VetName     string
	Date        time.Time
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
