package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	adoptionsService "github.com/juantevez/refugio-core/internal/adoptions/service"

	donationDomain "github.com/juantevez/refugio-core/internal/donations/domain"
	donationsService "github.com/juantevez/refugio-core/internal/donations/service"
)

// DTO para el seguimiento (Follow-up)
type followUpRequest struct {
	Notes     string   `json:"notes" binding:"required"`
	MediaURLs []string `json:"media_urls"` // URLs de las fotos subidas a S3/Cloudinary
}

type DonationHandler struct {
	donSvc   *donationsService.DonationService
	adoptSvc *adoptionsService.AdoptionService
}

func NewDonationHandler(ds *donationsService.DonationService, as *adoptionsService.AdoptionService) *DonationHandler {
	return &DonationHandler{
		donSvc:   ds,
		adoptSvc: as,
	}
}

type donationRequest struct {
	Amount    float64                       `json:"amount" binding:"required,gt=0"`
	Currency  string                        `json:"currency" binding:"required"`
	Source    donationDomain.DonationSource `json:"source" binding:"required"`
	Reference string                        `json:"reference" binding:"required"`
	AnimalID  *uuid.UUID                    `json:"animal_id"` // Opcional

	// Campos extra para cuando Source es TRANSFERENCIA
	DonorName  string  `json:"donor_name"`
	DonorEmail string  `json:"donor_email"`
	TaxID      *string `json:"tax_id"` // CUIT/CUIL
	Province   *string `json:"province"`
	Account    *string `json:"account"` // CBU/CVU
	Alias      *string `json:"alias"`
}

// DTO para iniciar una adopción
type adoptionRequest struct {
	AnimalID  uuid.UUID `json:"animal_id" binding:"required"`
	AdopterID uuid.UUID `json:"adopter_id" binding:"required"`
}

// RegisterDonation maneja POST /api/v1/donations
func (h *DonationHandler) RegisterDonation(c *gin.Context) {
	var req donationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ERROR REAL: %v", err) // Esto te va a decir la verdad en la terminal
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos de donación inválidos: " + err.Error()})
		return
	}

	// Armamos el objeto de detalles solo si vienen datos de transferencia
	var transferDetails *donationDomain.TransferDetails
	if req.TaxID != nil || req.Account != nil {
		transferDetails = &donationDomain.TransferDetails{
			TaxID:    derefString(req.TaxID),
			Province: derefString(req.Province),
			Account:  derefString(req.Account),
			Alias:    derefString(req.Alias),
		}
	}

	// Llamamos al Service con los nuevos campos
	donation, err := h.donSvc.RegisterDonation(
		c.Request.Context(),
		req.Amount,
		req.Currency,
		req.Source,
		req.Reference,
		req.AnimalID,
		req.DonorName,   // <--- Nuevo
		req.DonorEmail,  // <--- Nuevo
		transferDetails, // <--- Nuevo (puede ser nil)
	)

	if err != nil {
		log.Printf("ERROR REAL: %v", err) // Esto te va a decir la verdad en la terminal
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, donation)
}

// Función auxiliar (podes ponerla al final del archivo del handler)
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// CreateAdoption maneja POST /api/v1/adoptions
func (h *DonationHandler) CreateAdoption(c *gin.Context) {
	var req adoptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ERROR REAL: %v", err) // Esto te va a decir la verdad en la terminal
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos de adopción inválidos"})
		return
	}

	adoption, err := h.adoptSvc.StartAdoptionProcess(c.Request.Context(), req.AnimalID, req.AdopterID)
	if err != nil {
		log.Printf("ERROR REAL: %v", err) // Esto te va a decir la verdad en la terminal
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al procesar adopción"})
		return
	}

	// Devolvemos la adopción, que incluye el TrackingToken
	c.JSON(http.StatusCreated, adoption)
}

// SubmitFollowUp maneja POST /api/v1/adoptions/follow-up/:token
func (h *DonationHandler) SubmitFollowUp(c *gin.Context) {
	// Extraemos el token de la URL
	token := c.Param("token")
	if token == "" {
		log.Printf("ERROR REAL: %v", token) // Esto te va a decir la verdad en la terminal
		c.JSON(http.StatusBadRequest, gin.H{"error": "token de seguimiento requerido"})
		return
	}

	var req followUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ERROR REAL: %v", err) // Esto te va a decir la verdad en la terminal
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos de seguimiento inválidos"})
		return
	}

	// Llamamos al servicio de adopción que definimos anteriormente
	err := h.adoptSvc.SubmitFollowUp(c.Request.Context(), token, req.Notes, req.MediaURLs)
	if err != nil {
		if errors.Is(err, donationDomain.ErrInvalidToken) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token de seguimiento inválido o vencido"})
			return
		}
		log.Printf("ERROR REAL: %v", err) // Esto te va a decir la verdad en la terminal
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo registrar el seguimiento"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "¡Gracias por la actualización! Seguimiento registrado con éxito."})
}
