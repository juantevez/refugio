package domain

import (
	"time"

	"github.com/google/uuid"
)

type DonationSource string

const (
	SourceMercadoPago DonationSource = "MERCADO_PAGO"
	SourcePayPal      DonationSource = "PAYPAL"
	SourceTransfer    DonationSource = "BANK_TRANSFER"
)

type Donation struct {
	ID              uuid.UUID
	Amount          float64
	Currency        string // "ARS", "USD"
	Source          DonationSource
	ReferenceNumber string     // ID de transacción externa
	DonorName       string     // Opcional (Anónimo)
	AnimalID        *uuid.UUID // Opcional (Si es apadrinamiento)
	CreatedAt       time.Time
}
