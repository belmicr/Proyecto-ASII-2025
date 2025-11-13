package controllers_reservations

import (
	"net/http"
	domain "reservations/domain_reservations"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	svc domain.Service
}

func NewController(svc domain.Service) *Controller {
	return &Controller{svc: svc}
}

func (c *Controller) Create(ctx *gin.Context) {
	var req domain.Reservation
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	created, err := c.svc.Create(req)
	if err != nil {
		status := http.StatusInternalServerError

		// Determinar código de estado apropiado
		switch err.Error() {
		case "user not found", "hotel not found":
			status = http.StatusNotFound
		case "check-in must be before check-out",
			"check-in cannot be in the past",
			"las fechas se solapan con una reserva existente":
			status = http.StatusBadRequest
		}

		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, created)
}

func (c *Controller) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")

	reservation, err := c.svc.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Reservation not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, reservation)
}

func (c *Controller) List(ctx *gin.Context) {
	// Si hay user_id query param, filtrar por usuario
	userID := ctx.Query("user_id")

	if userID != "" {
		reservations, err := c.svc.GetByUserID(userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, reservations)
		return
	}

	// Sino, listar todas
	reservations, err := c.svc.List()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, reservations)
}

func (c *Controller) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req domain.Reservation
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	updated, err := c.svc.Update(id, req)
	if err != nil {
		status := http.StatusInternalServerError

		if err.Error() == "reservation not found" {
			status = http.StatusNotFound
		} else if err.Error() == "cannot modify cancelled reservation" ||
			err.Error() == "las fechas se solapan con otra reserva" {
			status = http.StatusBadRequest
		}

		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updated)
}

func (c *Controller) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.svc.Delete(id); err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "reservation not found" {
			status = http.StatusNotFound
		}

		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Reservation deleted successfully",
	})
}

// NUEVO: Endpoint de cancelación
func (c *Controller) Cancel(ctx *gin.Context) {
	id := ctx.Param("id")

	cancelled, err := c.svc.Cancel(id)
	if err != nil {
		status := http.StatusInternalServerError

		if err.Error() == "reservation not found" {
			status = http.StatusNotFound
		} else if err.Error() == "reservation already cancelled" ||
			err.Error() == "cannot cancel reservation that has already started" {
			status = http.StatusBadRequest
		}

		ctx.JSON(status, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cancelled)
}
