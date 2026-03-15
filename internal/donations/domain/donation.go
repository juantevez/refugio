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
	ID              uuid.UUID        `json:"id"`
	Amount          float64          `json:"amount"`
	Currency        string           `json:"currency"`
	Source          DonationSource   `json:"source"`
	ReferenceNumber string           `json:"reference_number"`
	AnimalID        *uuid.UUID       `json:"animal_id"`
	DonorName       string           `json:"donor_name"`
	DonorEmail      string           `json:"donor_email"`
	TransferDetails *TransferDetails `json:"transfer_details,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
}
