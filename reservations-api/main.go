package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	queues "reservations/clients_reservations"
	controllers "reservations/controllers_reservations"
	repositories "reservations/repositories_reservations"
	services "reservations/services_reservations"
)

func main() {
	// Inicializar repositorio
	repo := repositories.NewMock()

	// Cargar datos semilla
	if err := repo.SeedFromJSON("db/reservations.json"); err != nil {
		log.Printf("Warning: Could not load seed data: %v", err)
	} else {
		log.Println("Seed data loaded successfully")
	}

	// Inicializar RabbitMQ
	events := queues.NewRabbit(queues.RabbitConfig{
		Host:      "rabbitmq",
		Port:      "5672",
		Username:  "user",
		Password:  "root",
		QueueName: "reservations-events",
	})

	// Inicializar servicio y controlador
	svc := services.NewService(repo, events)
	ctrl := controllers.NewController(svc)

	// Configurar Gin
	r := gin.Default()
	_ = r.SetTrustedProxies(nil)

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"service": "reservations-api",
		})
	})

	// Rutas NUEVAS (RESTful)
	r.GET("/reservations/:id", ctrl.GetByID)
	r.GET("/reservations", ctrl.List) // Soporta ?user_id=X
	r.POST("/reservations", ctrl.Create)
	r.PUT("/reservations/:id", ctrl.Update)
	r.DELETE("/reservations/:id", ctrl.Delete)      // ← NUEVA
	r.POST("/reservations/:id/cancel", ctrl.Cancel) // ← NUEVA

	// Mantener compatibilidad con rutas antiguas
	r.POST("/createReservation", ctrl.Create)
	r.PUT("/edit/:id", ctrl.Update)

	log.Println("Reservations API running on :8086")
	if err := r.Run(":8086"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
