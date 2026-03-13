package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	// Capa de Adaptadores (Repositorios)
	adoptRepo "github.com/juantevez/refugio-core/internal/adoptions/adapters/repository"
	animalRepo "github.com/juantevez/refugio-core/internal/animals/adapters/repository"
	donRepo "github.com/juantevez/refugio-core/internal/donations/adapters/repository"

	// Capa de Servicios (Lógica de Negocio)
	adoptService "github.com/juantevez/refugio-core/internal/adoptions/service"
	animalService "github.com/juantevez/refugio-core/internal/animals/service"
	donService "github.com/juantevez/refugio-core/internal/donations/service"

	// Capa de Adaptadores (Handlers HTTP)
	adoptionHandler "github.com/juantevez/refugio-core/internal/adoptions/adapters/handler"
	animalHandler "github.com/juantevez/refugio-core/internal/animals/adapters/handler"
	donationHandler "github.com/juantevez/refugio-core/internal/donations/adapters/handler"
)

func main() {
	// 1. Conexión a Base de Datos
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://shelter_user:shelter_password@localhost:5432/shelter_db?sslmode=disable"
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error conectando a la DB: %v", err)
	}
	defer db.Close()

	// 2. Instanciar Adaptadores de Persistencia
	aRepo := animalRepo.NewPostgresRepository(db)
	dRepo := donRepo.NewPostgresRepository(db)
	adRepo := adoptRepo.NewPostgresRepository(db)

	// 3. Instanciar Servicios
	aSvc := animalService.NewAnimalService(aRepo)
	donSvc := donService.NewDonationService(dRepo)
	adSvc := adoptService.NewAdoptionService(adRepo)

	// 4. Instanciar Handlers
	aHandler := animalHandler.NewAnimalHandler(aSvc)
	dHandler := donationHandler.NewDonationHandler(donSvc, adSvc)
	adHandler := adoptionHandler.NewAdoptionHandler(adSvc)

	// 5. Configurar Servidor HTTP
	r := gin.Default()
	r.SetTrustedProxies(nil) // Esto quita el warning y asegura que no confíe en proxies externos

	// Rutas de la API v1
	api := r.Group("/api/v1")
	{
		// Endpoints de Animales
		animals := api.Group("/animals")
		{
			// ESTA ES LA QUE FALTA:
			animals.GET("", aHandler.ListAnimals)

			animals.POST("", aHandler.RegisterRescue)
			animals.GET("/:id", aHandler.GetAnimal)
		}

		// Endpoints de Donaciones
		api.POST("/donations", dHandler.RegisterDonation)

		// Endpoints de Adopciones (Unificados en un solo grupo)
		adoptions := api.Group("/adoptions")
		{
			adoptions.POST("", adHandler.StartAdoption)
			adoptions.POST("/follow-up/:token", adHandler.SubmitFollowUp)
		}
	}

	log.Println("Servidor corriendo en el puerto :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
