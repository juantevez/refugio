package domain

import (
	"context"

	"github.com/google/uuid"
)

type InventoryRepository interface {
	// Productos
	SaveProduct(ctx context.Context, product *Product) error
	UpdateProduct(ctx context.Context, product *Product) error
	GetProductByID(ctx context.Context, id uuid.UUID) (*Product, error)
	ListProducts(ctx context.Context) ([]Product, error)

	// Movimientos
	AddMovement(ctx context.Context, movement *StockMovement) error
	GetMovementsByProduct(ctx context.Context, productID uuid.UUID) ([]StockMovement, error)
}
