package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// UnitType define cómo medimos el stock (Kg, Unidades, Litros)
type UnitType string

const (
	UnitKg    UnitType = "KG"
	UnitUnit  UnitType = "UNIT"
	UnitLitre UnitType = "LITRE"
)

// Category define el tipo de insumo
type Category string

const (
	CatFood     Category = "FOOD"
	CatMedicine Category = "MEDICINE"
	CatCleaning Category = "CLEANING"
)

// Product representa un artículo del inventario
type Product struct {
	ID           uuid.UUID
	Name         string
	Category     Category
	Unit         UnitType
	MinThreshold float64 // Cantidad mínima para disparar alertas
	CurrentStock float64
	UpdatedAt    time.Time
}

// MovementType indica si entra o sale stock
type MovementType string

const (
	MovementIn  MovementType = "IN"  // Compra o donación
	MovementOut MovementType = "OUT" // Consumo o pérdida
)

// StockMovement registra la trazabilidad de los cambios
type StockMovement struct {
	ID        uuid.UUID
	ProductID uuid.UUID
	Quantity  float64
	Type      MovementType
	Reason    string // Ej: "Donación recibida", "Cena perros adultos"
	CreatedAt time.Time
}

// Errores de dominio
var (
	ErrInsufficientStock = errors.New("insufficient stock for this operation")
	ErrInvalidQuantity   = errors.New("quantity must be greater than zero")
)
