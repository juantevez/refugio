package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/juantevez/refugio-core/internal/animals/domain"
)

type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository crea una nueva instancia del adaptador de persistencia
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Save persiste un nuevo animal en la base de datos
func (r *PostgresRepository) Save(ctx context.Context, animal *domain.Animal) error {
	// 1. Asignamos el tiempo al objeto para que el puntero quede actualizado
	animal.CreatedAt = time.Now()

	// Si rescue_date también viene en cero, podrías inicializarlo aquí si fuera necesario,
	// pero lo ideal es que venga desde el Service/DTO.

	query := `
        INSERT INTO animal_management.animals (id, name, type, breed, status, rescue_date, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.ExecContext(ctx, query,
		animal.ID,
		animal.Name,
		animal.Species,
		animal.Breed,
		animal.Status,
		animal.RescueDate,
		animal.CreatedAt, // <--- Usamos el valor del objeto actualizado
	)
	return err
}

// GetByID busca un animal por su UUID
func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Animal, error) {
	var animal domain.Animal

	// Asegurate de incluir created_at y rescue_date en el SELECT
	query := `
        SELECT id, name, species, breed, status, rescue_date, created_at 
        FROM animal_management.animals 
        WHERE id = $1`

	// Si usás sqlx, el Get mapea automáticamente usando los tags db:"..."
	err := r.db.GetContext(ctx, &animal, query, id)
	if err != nil {
		return nil, err
	}

	return &animal, nil
}

// Cambiá la firma del método List para que acepte el map de filtros
func (r *PostgresRepository) List(ctx context.Context, filters map[string]interface{}) ([]domain.Animal, error) {
	var animals []domain.Animal
	query := `SELECT id, name, species, breed, status, rescue_date, created_at  FROM animal_management.animals`

	// Por ahora ignoramos los filters para que compile, luego podés implementarlos
	err := r.db.SelectContext(ctx, &animals, query)
	return animals, err
}

/*
// En tu método List o GetAll
query := `SELECT id, name, species, breed, status, rescue_date, created_at
          FROM animal_management.animals`
*/
// Update actualiza el estado o datos de un animal
func (r *PostgresRepository) Update(ctx context.Context, animal *domain.Animal) error {
	query := `
		UPDATE animal_management.animals 
		SET name = $2, status = $3, breed = $4 
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, animal.ID, animal.Name, animal.Status, animal.Breed)
	return err
}

// AddMedicalRecord guarda una nueva entrada en el historial clínico
func (r *PostgresRepository) AddMedicalRecord(ctx context.Context, record *domain.MedicalRecord) error {
	query := `
		INSERT INTO animal_management.medical_records (id, animal_id, description, created_at)
		VALUES ($1, $2, $3, $4)`

	_, err := r.db.ExecContext(ctx, query,
		record.ID,
		record.AnimalID,
		record.Description,
		record.CreatedAt,
	)
	return err
}
