package handler

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	animalsDomain "github.com/juantevez/refugio-core/internal/animals/domain"
	"github.com/juantevez/refugio-core/internal/animals/service"
)

type AnimalHandler struct {
	// Fijate que acá diga 'service' (o 'svc')
	service *service.AnimalService
}

// El constructor también debe asignarlo
func NewAnimalHandler(s *service.AnimalService) *AnimalHandler {
	return &AnimalHandler{
		service: s,
	}
}

// Request DTO (Data Transfer Object) para validación
type registerRequest struct {
	Name       string    `json:"name" binding:"required"`
	Species    string    `json:"species" binding:"required"`
	Breed      string    `json:"breed" binding:"required"` // <-- Faltaba esta línea
	Status     string    `json:"status" binding:"required"`
	RescueDate time.Time `json:"rescue_date" binding:"required"`
}

// RegisterRescue maneja el POST /api/v1/animals
func (h *AnimalHandler) RegisterRescue(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("ERROR DE BINDING: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// LLAMADA CORREGIDA: Pasamos los parámetros uno por uno como pide el servicio
	// Hacemos el cast de req.Species a animalsDomain.Species
	animal, err := h.service.RegisterRescue(
		c.Request.Context(),
		req.Name,
		strings.ToUpper(req.Species), // "Dog" → "DOG"
		req.Breed,
		req.RescueDate, // El nuevo parámetro que agregaste
	)

	if err != nil {
		log.Printf("ERROR REAL EN EL SERVICIO: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo registrar el rescate"})
		return
	}

	c.JSON(http.StatusCreated, animal)
}

// GetAnimal maneja el GET /api/v1/animals/:id
func (h *AnimalHandler) GetAnimal(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Printf("ERROR REAL: %v", err) // Esto te va a decir la verdad en la terminal
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de animal inválido"})
		return
	}

	animal, err := h.service.GetAnimalDetails(c.Request.Context(), id)
	if err != nil {
		if err == animalsDomain.ErrAnimalNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "animal no encontrado"})
			return
		}
		log.Printf("ERROR REAL: %v", err) // Esto te va a decir la verdad en la terminal
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, animal)
}

func (h *AnimalHandler) ListAnimals(c *gin.Context) {
	// Por ahora mandamos nil en los filtros para que traiga todo
	animals, err := h.service.List(c.Request.Context(), nil)
	if err != nil {
		log.Printf("ERROR REAL: %v", err) // Esto te va a decir la verdad en la terminal
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, animals)
}
