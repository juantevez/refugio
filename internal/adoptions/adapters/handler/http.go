package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/juantevez/refugio-core/internal/adoptions/domain"
	"github.com/juantevez/refugio-core/internal/adoptions/service"
)

type AdoptionHandler struct {
	service *service.AdoptionService
}

func NewAdoptionHandler(s *service.AdoptionService) *AdoptionHandler {
	return &AdoptionHandler{service: s}
}

// Request structs para validar la entrada
type startAdoptionRequest struct {
	AnimalID  uuid.UUID `json:"animal_id" binding:"required"`
	AdopterID uuid.UUID `json:"adopter_id" binding:"required"`
}

type followUpRequest struct {
	Notes string   `json:"notes" binding:"required"`
	Media []string `json:"media"`
}

// StartAdoption maneja POST /api/v1/adoptions
func (h *AdoptionHandler) StartAdoption(c *gin.Context) {
	var req startAdoptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ERROR REAL: %v", err) // Esto te va a decir la verdad en la terminal
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adoption, err := h.service.StartAdoptionProcess(c.Request.Context(), req.AnimalID, req.AdopterID)
	if err != nil {
		log.Printf("ERROR REAL: %v", err) // Esto te va a decir la verdad en la terminal
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo iniciar el proceso"})
		return
	}

	c.JSON(http.StatusCreated, adoption)
}

// SubmitFollowUp maneja POST /api/v1/adoptions/follow-up/:token
func (h *AdoptionHandler) SubmitFollowUp(c *gin.Context) {
	token := c.Param("token")
	var req followUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ERROR REAL: %v", err) // Esto te va a decir la verdad en la terminal
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.SubmitFollowUp(c.Request.Context(), token, req.Notes, req.Media)
	if err != nil {
		if err == domain.ErrInvalidToken {
			log.Printf("ERROR REAL: %v", err) // Esto te va a decir la verdad en la terminal
			c.JSON(http.StatusNotFound, gin.H{"error": "Token de seguimiento inválido"})
			return
		}
		log.Printf("ERROR REAL: %v", err) // Esto te va a decir la verdad en la terminal
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar seguimiento"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Seguimiento registrado con éxito"})
}
