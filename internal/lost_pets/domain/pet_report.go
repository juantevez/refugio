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
	ID             uuid.UUID      `db:"id"              json:"id"`
	Type           ReportType     `db:"type"            json:"type"`
	Species        PetSpecies     `db:"species"         json:"species"`
	Description    string         `db:"description"     json:"description"`
	PhotoS3Keys    []string       `db:"photo_s3_key"    json:"photo_s3_key"`
	Location       Point          `db:"location"        json:"location"`
	LocationSource LocationSource `db:"location_source" json:"location_source"`
	RadiusMeters   int            `db:"radius_meters"   json:"radius_meters"`
	Status         ReportStatus   `db:"status"          json:"status"`
	ContactName    string         `db:"contact_name"    json:"contact_name"`
	ContactEmail   string         `db:"contact_email"   json:"contact_email"`
	ContactPhone   string         `db:"contact_phone"   json:"contact_phone"`
	ReportedAt     time.Time      `db:"reported_at"     json:"reported_at"`
	CreatedAt      time.Time      `db:"created_at"      json:"created_at"`
}

// Nuevo tipo en domain:
type ImageInput struct {
	Data        []byte
	Filename    string
	ContentType string
}
