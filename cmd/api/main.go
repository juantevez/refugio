package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Driver de Postgres

	// Importamos nuestros módulos internos
	animalAdapter "github.com/juantevez/refugio/internal/animals/adapters/repository"
	animalService "github.com/juantevez/refugio/internal/animals/service"
	donationAdapter "github.com/juantevez/refugio/internal/donations/adapters/repository"
	donationService "github.com/juantevez/refugio/internal/donations/service"
)

func main() {
	// 1. Conexión a Base de Datos
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error conectando a la DB: %v", err)
	}
	defer db.Close()

	// 2. Instanciar Adaptadores (Infraestructura)
	animalRepo := animalAdapter.NewPostgresRepository(db)
	donationRepo := donationAdapter.NewPostgresRepository(db)

	// 3. Instanciar Servicios (Dominio/Aplicación)
	animalSvc := animalService.NewAnimalService(animalRepo)
	donSvc := donationService.NewDonationService(donationRepo)
	adoptSvc := donationService.NewAdoptionService(donationRepo)

	// 4. Configurar Servidor HTTP (Framework Adapter)
	r := gin.Default()

	// Agrupamos rutas para mantener orden
	api := r.Group("/api/v1")
	{
		// Podrías pasar el servicio a un Handler de cada módulo
		api.GET("/animals", func(c *gin.Context) {
			animals, _ := animalSvc.List(c.Request.Context())
			c.JSON(http.StatusOK, animals)
		})

		api.POST("/donations", func(c *gin.Context) {
			// Lógica para registrar donación usando donSvc...
		})
	}

	log.Println("Servidor corriendo en :8080")
	r.Run(":8080")
}
