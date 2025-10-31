package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	queues "hotels/clients_hotels"
	controllers "hotels/controllers_hotels"
	repositories "hotels/repositories_hotels"
	services "hotels/services_hotels"
)

func main() {
	// Usamos repo en memoria para probar
	mainRepository := repositories.NewMock()
	_ = mainRepository.SeedFromJSON("db/hotels.json")

	// Cola “dummy” (no falla si no hay Rabbit)
	eventsQueue := queues.NewRabbit(queues.RabbitConfig{
		Host:      "rabbitmq",
		Port:      "5672",
		Username:  "user",
		Password:  "root",
		QueueName: "hotels-news",
	})

	service := services.NewService(mainRepository, eventsQueue)
	controller := controllers.NewController(service)

	router := gin.Default()
	_ = router.SetTrustedProxies(nil) // <<--- agrega esto para sacar el warning

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/hotels/:id", controller.GetHotelByID)
	router.GET("/hotels", controller.GetHotels)
	router.POST("/createHotel", controller.Create)
	router.PUT("/edit/:id", controller.Update)

	if err := router.Run(":8085"); err != nil { // <<--- cambia el puerto a 8085
		log.Fatalf("error running application: %v", err)
	}
}
