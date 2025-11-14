package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	ctrl "reservations/controllers_reservations"
	"reservations/domain_reservations"
	repo "reservations/repositories_reservations"
	svc "reservations/services_reservations"
)

func main() {
	// Repo en memoria
	mainRepository := repo.NewMock()
	_ = mainRepository.SeedFromJSON("db/reservations.json")

	// Cola de eventos opcional (nil tipado)
	var q domain_reservations.EventQueue = nil

	// Service: repo + queue + cliente de hoteles (nil -> noop)
	service := svc.NewService(mainRepository, q, nil)

	// Controller
	controller := ctrl.NewController(service)

	// Router = r
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/reservations/:id", controller.GetByID)
	r.GET("/reservations", controller.List)
	r.POST("/createReservation", controller.Create)
	r.PUT("/edit/:id", controller.Update)

	if err := r.Run(":8086"); err != nil {
		log.Fatalf("error running application: %v", err)
	}
}
