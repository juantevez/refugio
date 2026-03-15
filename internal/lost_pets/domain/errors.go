package domain

import "errors"

var (
	// ErrReportNotFound cuando se busca un reporte por ID y no existe.
	ErrReportNotFound = errors.New("reporte no encontrado")

	// ErrInvalidLocation cuando lat/long están fuera de rango o ausentes.
	ErrInvalidLocation = errors.New("ubicación inválida: se requiere latitud y longitud")

	// ErrInvalidReportType cuando el tipo no es LOST ni FOUND.
	ErrInvalidReportType = errors.New("el tipo de reporte debe ser LOST o FOUND")

	// ErrReportAlreadyResolved cuando se intenta modificar un reporte ya resuelto.
	ErrReportAlreadyResolved = errors.New("el reporte ya fue resuelto")
)
