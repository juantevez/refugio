package domain

import "errors"

var (
	// ErrInvalidAmount se dispara cuando el monto es cero o negativo
	// ErrInvalidAmount = errors.New("el monto de la donación debe ser mayor a cero")

	// ErrInvalidCurrency para cuando se intenta donar en una moneda no soportada
	ErrInvalidCurrency = errors.New("la moneda especificada no es válida")

	// ErrDonationNotFound para consultas de historial
	ErrDonationNotFound = errors.New("no se encontró la donación especificada")

	// ErrInvalidToken para el seguimiento de adopciones
	ErrInvalidToken = errors.New("el token de seguimiento es inválido o ha expirado")
)
