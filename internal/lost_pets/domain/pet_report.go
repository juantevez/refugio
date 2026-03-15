package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReportType string
type ReportStatus string
type LocationSource string
type PetSpecies string

const (
	ReportTypeLost  ReportType = "LOST"
	ReportTypeFound ReportType = "FOUND"

	ReportStatusActive   ReportStatus = "ACTIVE"
	ReportStatusResolved ReportStatus = "RESOLVED"
	ReportStatusExpired  ReportStatus = "EXPIRED"

	LocationSourceGPS    LocationSource = "GPS"
	LocationSourceEXIF   LocationSource = "EXIF"
	LocationSourceManual LocationSource = "MANUAL"

	SpeciesDog   PetSpecies = "DOG"
	SpeciesCat   PetSpecies = "CAT"
	SpeciesOther PetSpecies = "OTHER"
)

type PetReport struct {
	ID             uuid.UUID      `db:"id"`
	Type           ReportType     `db:"type"`
	Species        PetSpecies     `db:"species"`
	Description    string         `db:"description"`
	PhotoS3Key     string         `db:"photo_s3_key"`
	Location       Point          `db:"location"`
	LocationSource LocationSource `db:"location_source"`
	RadiusMeters   int            `db:"radius_meters"`
	Status         ReportStatus   `db:"status"`
	ContactName    string         `db:"contact_name"`
	ContactEmail   string         `db:"contact_email"`
	ContactPhone   string         `db:"contact_phone"`
	ReportedAt     time.Time      `db:"reported_at"`
	CreatedAt      time.Time      `db:"created_at"`
}
