package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	adoptionsDomain "github.com/juantevez/refugio-core/internal/adoptions/domain"
)

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateAdoption(ctx context.Context, a *adoptionsDomain.Adoption) error {
	query := `INSERT INTO adoptions_donations.adoptions (id, animal_id, adopter_id, status, tracking_token, adopted_at) 
              VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, a.ID, a.AnimalID, a.AdopterID, a.Status, a.TrackingToken, a.AdoptedAt)
	return err
}

func (r *PostgresRepository) GetAdoptionByToken(ctx context.Context, token string) (*adoptionsDomain.Adoption, error) {
	var a adoptionsDomain.Adoption
	query := `SELECT id, animal_id, adopter_id, status, tracking_token, adopted_at 
              FROM adoptions_donations.adoptions WHERE tracking_token = $1`
	err := r.db.GetContext(ctx, &a, query, token)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *PostgresRepository) AddFollowUp(ctx context.Context, f *adoptionsDomain.FollowUp) error {
	query := `INSERT INTO adoptions_donations.follow_ups (id, adoption_id, notes, media_urls, created_at) 
              VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, f.ID, f.AdoptionID, f.Notes, f.MediaURLs, f.CreatedAt)
	return err
}
