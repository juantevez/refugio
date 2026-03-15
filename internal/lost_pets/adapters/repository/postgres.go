package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/juantevez/refugio-core/internal/lost_pets/domain"
)

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// petReportRow es el struct intermedio para leer desde PostGIS.
// ST_X devuelve longitud, ST_Y devuelve latitud — no se puede mapear
// GEOGRAPHY directamente con sqlx.
type petReportRow struct {
	ID             uuid.UUID      `db:"id"`
	Type           string         `db:"type"`
	Species        string         `db:"species"`
	Description    string         `db:"description"`
	PhotoS3Key     string         `db:"photo_s3_key"`
	Lat            float64        `db:"lat"`
	Long           float64        `db:"long"`
	LocationSource string         `db:"location_source"`
	RadiusMeters   int            `db:"radius_meters"`
	Status         string         `db:"status"`
	ContactName    string         `db:"contact_name"`
	ContactEmail   string         `db:"contact_email"`
	ContactPhone   string         `db:"contact_phone"`
	ReportedAt     time.Time      `db:"reported_at"`
	CreatedAt      time.Time      `db:"created_at"`
}

func (row *petReportRow) toDomain() *domain.PetReport {
	return &domain.PetReport{
		ID:             row.ID,
		Type:           domain.ReportType(row.Type),
		Species:        domain.PetSpecies(row.Species),
		Description:    row.Description,
		PhotoS3Key:     row.PhotoS3Key,
		Location:       domain.Point{Lat: row.Lat, Long: row.Long},
		LocationSource: domain.LocationSource(row.LocationSource),
		RadiusMeters:   row.RadiusMeters,
		Status:         domain.ReportStatus(row.Status),
		ContactName:    row.ContactName,
		ContactEmail:   row.ContactEmail,
		ContactPhone:   row.ContactPhone,
		ReportedAt:     row.ReportedAt,
		CreatedAt:      row.CreatedAt,
	}
}

// selectColumns es el SELECT reutilizable que extrae lat/long desde GEOGRAPHY.
const selectColumns = `
	id, type, species, description, photo_s3_key,
	ST_Y(location::geometry) AS lat,
	ST_X(location::geometry) AS long,
	location_source, radius_meters, status,
	contact_name, contact_email, contact_phone,
	reported_at, created_at`

// Save persiste un nuevo reporte. Usa ST_MakePoint(long, lat) — orden PostGIS: X=long, Y=lat.
func (r *PostgresRepository) Save(ctx context.Context, report *domain.PetReport) error {
	query := `
		INSERT INTO lost_pets.pet_reports (
			id, type, species, description, photo_s3_key,
			location, location_source, radius_meters, status,
			contact_name, contact_email, contact_phone,
			reported_at, created_at
		) VALUES (
			$1, $2, $3, $4, $5,
			ST_MakePoint($6, $7)::geography, $8, $9, $10,
			$11, $12, $13,
			$14, $15
		)`

	_, err := r.db.ExecContext(ctx, query,
		report.ID,
		report.Type,
		report.Species,
		report.Description,
		report.PhotoS3Key,
		report.Location.Long, // $6 — ST_MakePoint(X, Y) = (long, lat)
		report.Location.Lat,  // $7
		report.LocationSource,
		report.RadiusMeters,
		report.Status,
		report.ContactName,
		report.ContactEmail,
		report.ContactPhone,
		report.ReportedAt,
		report.CreatedAt,
	)
	return err
}

// GetByID busca un reporte por UUID.
func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.PetReport, error) {
	query := `SELECT` + selectColumns + ` FROM lost_pets.pet_reports WHERE id = $1`

	var row petReportRow
	if err := r.db.GetContext(ctx, &row, query, id); err != nil {
		return nil, err
	}
	return row.toDomain(), nil
}

// SearchNearby devuelve reportes ACTIVE dentro del radio especificado en metros.
// ST_DWithin sobre GEOGRAPHY opera en metros directamente.
func (r *PostgresRepository) SearchNearby(ctx context.Context, area domain.SearchArea) ([]domain.PetReport, error) {
	query := `
		SELECT` + selectColumns + `
		FROM lost_pets.pet_reports
		WHERE status = 'ACTIVE'
		  AND ST_DWithin(
		        location,
		        ST_MakePoint($1, $2)::geography,
		        $3
		      )
		ORDER BY reported_at DESC`

	var rows []petReportRow
	if err := r.db.SelectContext(ctx, &rows, query,
		area.Center.Long, // $1 — ST_MakePoint(X, Y) = (long, lat)
		area.Center.Lat,  // $2
		area.RadiusMeters, // $3 — metros
	); err != nil {
		return nil, err
	}

	reports := make([]domain.PetReport, 0, len(rows))
	for _, row := range rows {
		reports = append(reports, *row.toDomain())
	}
	return reports, nil
}

// UpdateStatus actualiza el estado de un reporte (ACTIVE → RESOLVED / EXPIRED).
func (r *PostgresRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ReportStatus) error {
	query := `UPDATE lost_pets.pet_reports SET status = $2 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, status)
	return err
}
