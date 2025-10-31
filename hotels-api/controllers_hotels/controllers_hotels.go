package controllers_hotels

import (
	"context"
	"net/http"
	"strings"

	"hotels/domain_hotels"

	"github.com/gin-gonic/gin"
)

// Esta interfaz la satisface services_hotels.Service
type Service interface {
	GetByID(ctx context.Context, id string) (domain_hotels.Hotel, error)
	List(ctx context.Context, q string) ([]domain_hotels.Hotel, error)
	Create(ctx context.Context, h domain_hotels.Hotel) (domain_hotels.Hotel, error)
	Update(ctx context.Context, id string, h domain_hotels.Hotel) (domain_hotels.Hotel, error)
}

type Controller struct {
	service Service
}

func NewController(s Service) *Controller { return &Controller{service: s} }

// GET /hotels/:id
func (c *Controller) GetHotelByID(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Param("id"))
	h, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.String(http.StatusNotFound, "not found")
		return
	}
	ctx.JSON(http.StatusOK, h)
}

// GET /hotels?q=...
func (c *Controller) GetHotels(ctx *gin.Context) {
	q := ctx.Query("q")
	list, err := c.service.List(ctx.Request.Context(), q)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, list)
}

// POST /createHotel
func (c *Controller) Create(ctx *gin.Context) {
	var in domain_hotels.Hotel
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.String(http.StatusBadRequest, "bad request")
		return
	}

	// Validaciones mínimas
	if strings.TrimSpace(in.Name) == "" ||
		strings.TrimSpace(in.City) == "" ||
		in.PricePerNight <= 0 ||
		in.Stars < 1 || in.Stars > 5 {
		ctx.String(http.StatusBadRequest, "invalid hotel payload")
		return
	}

	out, err := c.service.Create(ctx.Request.Context(), in)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, out)
}

// PUT /edit/:id
func (c *Controller) Update(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Param("id"))

	var in domain_hotels.Hotel
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.String(http.StatusBadRequest, "bad request")
		return
	}

	// Validaciones mínimas (permitimos updates parciales, solo validamos si vienen)
	if in.Name != "" && strings.TrimSpace(in.Name) == "" {
		ctx.String(http.StatusBadRequest, "invalid hotel payload")
		return
	}
	if in.City != "" && strings.TrimSpace(in.City) == "" {
		ctx.String(http.StatusBadRequest, "invalid hotel payload")
		return
	}
	if in.PricePerNight < 0 {
		ctx.String(http.StatusBadRequest, "invalid hotel payload")
		return
	}
	if in.Stars != 0 && (in.Stars < 1 || in.Stars > 5) {
		ctx.String(http.StatusBadRequest, "invalid hotel payload")
		return
	}

	out, err := c.service.Update(ctx.Request.Context(), id, in)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, out)
}
