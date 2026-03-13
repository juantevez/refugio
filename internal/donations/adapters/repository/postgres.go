package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/juantevez/refugio-core/internal/donations/domain"
)

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) SaveDonation(ctx context.Context, d *domain.Donation) error {
	query := `
        INSERT INTO adoptions_donations.donations (
            id, amount, currency, source, reference_number, 
            animal_id, donor_name, donor_email, tax_id, 
            province, account_number, account_alias, created_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    `

	// Manejo de punteros para campos opcionales
	var taxID, province, account, alias interface{}
	if d.TransferDetails != nil {
		taxID = d.TransferDetails.TaxID
		province = d.TransferDetails.Province
		account = d.TransferDetails.Account
		alias = d.TransferDetails.Alias
	}

	_, err := r.db.ExecContext(ctx, query,
		d.ID,
		d.Amount,
		d.Currency,
		d.Source,
		d.ReferenceNumber,
		d.AnimalID, // sqlx maneja el *uuid.UUID como NULL si es nil
		d.DonorName,
		d.DonorEmail,
		taxID,
		province,
		account,
		alias,
		time.Now(),
	)
	return err
}

// GetTotalDonations adaptado a la nueva tabla
func (r *PostgresRepository) GetTotalDonations(ctx context.Context) (float64, error) {
	var total float64
	query := `SELECT COALESCE(SUM(amount), 0) FROM animal_management.donations`

	err := r.db.GetContext(ctx, &total, query)
	return total, err
}

// Mantengo CreateAdoption por ahora, pero recordá que deberá apuntar al esquema correcto
func (r *PostgresRepository) CreateAdoption(ctx context.Context, a *domain.Adoption) error {
	query := `
        INSERT INTO animal_management.adoptions (id, animal_id, adopter_id, tracking_token, status, adopted_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := r.db.ExecContext(ctx, query,
		a.ID, a.AnimalID, a.AdopterID, a.TrackingToken, a.Status, a.AdoptedAt,
	)
	return err
}
