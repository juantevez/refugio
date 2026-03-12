package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/juantevez/refugio-core/internal/inventory/domain"
)

type InventoryService struct {
	repo domain.InventoryRepository
}

func NewInventoryService(r domain.InventoryRepository) *InventoryService {
	return &InventoryService{repo: r}
}

// RegisterMovement gestiona la entrada o salida de insumos
func (s *InventoryService) RegisterMovement(ctx context.Context, productID uuid.UUID, qty float64, mType domain.MovementType, reason string) error {
	if qty <= 0 {
		return domain.ErrInvalidQuantity
	}

	product, err := s.repo.GetProductByID(ctx, productID)
	if err != nil {
		return err
	}

	// Validar stock si es una salida
	if mType == domain.MovementOut {
		if product.CurrentStock < qty {
			return domain.ErrInsufficientStock
		}
		product.CurrentStock -= qty
	} else {
		product.CurrentStock += qty
	}

	// Crear el registro del movimiento
	movement := &domain.StockMovement{
		ID:        uuid.New(),
		ProductID: productID,
		Quantity:  qty,
		Type:      mType,
		Reason:    reason,
		CreatedAt: time.Now(),
	}

	// Persistir ambos (En el adaptador Postgres esto debería ser una transacción)
	if err := s.repo.AddMovement(ctx, movement); err != nil {
		return err
	}

	product.UpdatedAt = time.Now()
	return s.repo.UpdateProduct(ctx, product)
}
