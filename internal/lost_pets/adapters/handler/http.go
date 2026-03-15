package handler

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/juantevez/refugio-core/internal/lost_pets/domain"
	"github.com/juantevez/refugio-core/internal/lost_pets/service"
)

type LostPetsHandler struct {
	service *service.LostPetsService
}

func NewLostPetsHandler(s *service.LostPetsService) *LostPetsHandler {
	return &LostPetsHandler{service: s}
}

// POST /api/v1/lost-pets/reports
// Recibe multipart/form-data con campos de texto + foto opcional.
func (h *LostPetsHandler) CreateReport(c *gin.Context) {
	reportType := domain.ReportType(strings.ToUpper(c.PostForm("type")))
	species := domain.PetSpecies(strings.ToUpper(c.PostForm("species")))
	description := c.PostForm("description")
	contactName := c.PostForm("contact_name")
	contactEmail := c.PostForm("contact_email")
	contactPhone := c.PostForm("contact_phone")

	if contactName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "contact_name es requerido"})
		return
	}

	// Ubicación del cliente (GPS del celular/navegador)
	var clientLocation *domain.Point
	latStr := c.PostForm("lat")
	longStr := c.PostForm("long")
	if latStr != "" && longStr != "" {
		lat, errLat := strconv.ParseFloat(latStr, 64)
		long, errLong := strconv.ParseFloat(longStr, 64)
		if errLat != nil || errLong != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "lat y long deben ser números válidos"})
			return
		}
		clientLocation = &domain.Point{Lat: lat, Long: long}
	}

	// Foto opcional
	var imageData []byte
	var imageFilename, imageContentType string
	file, header, err := c.Request.FormFile("photo")
	if err == nil {
		defer file.Close()

		imageContentType = header.Header.Get("Content-Type")
		if imageContentType != "image/jpeg" && imageContentType != "image/png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "solo se aceptan imágenes JPEG o PNG"})
			return
		}
		if header.Size > 5<<20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "la foto no puede superar 5MB"})
			return
		}

		imageData, err = io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error leyendo la foto"})
			return
		}
		imageFilename = header.Filename
	}

	report, err := h.service.CreateReport(
		c.Request.Context(),
		reportType,
		species,
		description,
		contactName,
		contactEmail,
		contactPhone,
		clientLocation,
		imageData,
		imageFilename,
		imageContentType,
	)
	if err != nil {
		switch err {
		case domain.ErrInvalidReportType:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case domain.ErrInvalidLocation:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			log.Printf("ERROR al crear reporte: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo crear el reporte"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":              report.ID,
		"type":            report.Type,
		"species":         report.Species,
		"status":          report.Status,
		"location_source": report.LocationSource,
		"contact_name":    report.ContactName,
		"reported_at":     report.ReportedAt,
	})
}

// GET /api/v1/lost-pets/reports/search?lat=X&long=Y&radius=Z
func (h *LostPetsHandler) SearchNearby(c *gin.Context) {
	lat, errLat := strconv.ParseFloat(c.Query("lat"), 64)
	long, errLong := strconv.ParseFloat(c.Query("long"), 64)
	if errLat != nil || errLong != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat y long son requeridos y deben ser números válidos"})
		return
	}

	radius := 1000 // radio por defecto: 1km
	if r := c.Query("radius"); r != "" {
		parsed, err := strconv.Atoi(r)
		if err != nil || parsed <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "radius debe ser un entero positivo en metros"})
			return
		}
		radius = parsed
	}

	reports, err := h.service.SearchNearby(c.Request.Context(), domain.SearchArea{
		Center:       domain.Point{Lat: lat, Long: long},
		RadiusMeters: radius,
	})
	if err != nil {
		log.Printf("ERROR en búsqueda: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error en la búsqueda"})
		return
	}

	c.JSON(http.StatusOK, reports)
}

// GET /api/v1/lost-pets/reports/:id
func (h *LostPetsHandler) GetReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de reporte inválido"})
		return
	}

	report, err := h.service.GetReport(c.Request.Context(), id)
	if err != nil {
		if err == domain.ErrReportNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		log.Printf("ERROR al obtener reporte: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error obteniendo el reporte"})
		return
	}

	c.JSON(http.StatusOK, report)
}

// PATCH /api/v1/lost-pets/reports/:id/resolve
func (h *LostPetsHandler) ResolveReport(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de reporte inválido"})
		return
	}

	if err := h.service.ResolveReport(c.Request.Context(), id); err != nil {
		switch err {
		case domain.ErrReportNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case domain.ErrReportAlreadyResolved:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			log.Printf("ERROR al resolver reporte: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo resolver el reporte"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reporte marcado como resuelto"})
}
