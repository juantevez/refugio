package domain

import (
	"context"

	"github.com/google/uuid"
)

// AnimalRepository es el puerto de salida para la persistencia
type AnimalRepository interface {
	Save(ctx context.Context, animal *Animal) error
	GetByID(ctx context.Context, id uuid.UUID) (*Animal, error)
	List(ctx context.Context, filter map[string]interface{}) ([]Animal, error)
	Update(ctx context.Context, animal *Animal) error

	// Para los registros médicos
	AddMedicalRecord(ctx context.Context, record *MedicalRecord) error
}
