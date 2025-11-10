package controllers_reservations

import (
	"context"
	"net/http"
	"strings"

	"reservations/domain_reservations"

	"github.com/gin-gonic/gin"
)

type Service interface {
	GetByID(ctx context.Context, id string) (domain_reservations.Reservation, error)
	List(ctx context.Context, hotelID, userID, status string) ([]domain_reservations.Reservation, error)
	Create(ctx context.Context, r domain_reservations.Reservation) (domain_reservations.Reservation, error)
	Update(ctx context.Context, id string, r domain_reservations.Reservation) (domain_reservations.Reservation, error)
}

type Controller struct{ s Service }

func NewController(s Service) *Controller { return &Controller{s: s} }

// GET /reservations/:id
func (c *Controller) GetByID(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Param("id"))
	out, err := c.s.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.String(http.StatusNotFound, "not found")
		return
	}
	ctx.JSON(http.StatusOK, out)
}

// GET /reservations?hotel_id=&user_id=&status=
func (c *Controller) List(ctx *gin.Context) {
	hotelID := ctx.Query("hotel_id")
	userID := ctx.Query("user_id")
	status := ctx.Query("status")
	out, err := c.s.List(ctx.Request.Context(), hotelID, userID, status)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, out)
}

// POST /createReservation
func (c *Controller) Create(ctx *gin.Context) {
	var in domain_reservations.Reservation
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.String(http.StatusBadRequest, "bad request")
		return
	}
	// Validaciones mínimas
	if strings.TrimSpace(in.HotelID) == "" ||
		strings.TrimSpace(in.UserID) == "" ||
		strings.TrimSpace(in.CheckIn) == "" ||
		strings.TrimSpace(in.CheckOut) == "" ||
		in.Guests <= 0 {
		ctx.String(http.StatusBadRequest, "invalid reservation payload")
		return
	}
	// Normalizo y valido status
	if in.Status == "" {
		in.Status = domain_reservations.StatusPending
	} else if !domain_reservations.IsValidStatus(strings.ToLower(in.Status)) {
		ctx.String(http.StatusBadRequest, "invalid status")
		return
	}

	out, err := c.s.Create(ctx.Request.Context(), in)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, out)
}

// PUT /edit/:id
func (c *Controller) Update(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Param("id"))

	var in domain_reservations.Reservation
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.String(http.StatusBadRequest, "bad request")
		return
	}
	// Validaciones rápidas
	if in.Guests < 0 {
		ctx.String(http.StatusBadRequest, "invalid reservation payload")
		return
	}
	// Valido status si vino
	if in.Status != "" && !domain_reservations.IsValidStatus(strings.ToLower(in.Status)) {
		ctx.String(http.StatusBadRequest, "invalid status")
		return
	}

	out, err := c.s.Update(ctx.Request.Context(), id, in)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, out)
}
