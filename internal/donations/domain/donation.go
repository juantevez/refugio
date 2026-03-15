package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type DonationSource string

const (
	SourceMercadoPago DonationSource = "MERCADO_PAGO"
	SourceTransfer    DonationSource = "TRANSFERENCIA"
)

var ErrInvalidAmount = errors.New("el monto de la donación debe ser mayor a cero")

type TransferDetails struct {
	TaxID    string // CUIT/CUIL
	Province string
	Account  string // CBU/CVU
	Alias    string
}

type Donation struct {
	ID              uuid.UUID
	Amount          float64
	Currency        string
	Source          DonationSource
	ReferenceNumber string
	AnimalID        *uuid.UUID
	DonorName       string
	DonorEmail      string
	TransferDetails *TransferDetails
	CreatedAt       time.Time
}
