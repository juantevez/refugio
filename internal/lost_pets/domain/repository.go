package domain

import (
	"context"

	"github.com/google/uuid"
)

// PetReportRepository es el puerto de salida para la persistencia de reportes.
type PetReportRepository interface {
	Save(ctx context.Context, report *PetReport) error
	GetByID(ctx context.Context, id uuid.UUID) (*PetReport, error)
	SearchNearby(ctx context.Context, area SearchArea) ([]PetReport, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status ReportStatus) error
}
