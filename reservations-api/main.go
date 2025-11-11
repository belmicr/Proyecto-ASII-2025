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
	repo := repositories.NewMock()
	_ = repo.SeedFromJSON("db/reservations.json")

	events := queues.NewRabbit(queues.RabbitConfig{
		Host: "rabbitmq", Port: "5672", Username: "user", Password: "root", QueueName: "reservations-events",
	})
	svc := services.NewService(repo, events)
	ctrl := controllers.NewController(svc)

	r := gin.Default()
	_ = r.SetTrustedProxies(nil)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/reservations/:id", ctrl.GetByID)
	r.GET("/reservations", ctrl.List)
	r.POST("/createReservation", ctrl.Create)
	r.PUT("/edit/:id", ctrl.Update)

	if err := r.Run(":8086"); err != nil {
		log.Fatalf("error: %v", err)
	}
}
