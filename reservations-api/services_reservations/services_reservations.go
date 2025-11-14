package services_reservations

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"reservations/domain_reservations"
)

// Cliente mínimo para consultar hotels-api
type hotelsClient interface {
	HotelExists(id string) bool
}

// Fallback por si no inyectan cliente (no bloquea P1)
type noopHotels struct{}

func (n *noopHotels) HotelExists(_ string) bool { return true }

// Service implementa la lógica de reservas
type Service struct {
	repo   domain_reservations.Repository
	q      domain_reservations.EventQueue
	hotels hotelsClient
}

// NewService ahora recibe también el cliente de hoteles
func NewService(r domain_reservations.Repository, q domain_reservations.EventQueue, h hotelsClient) *Service {
	if h == nil {
		h = &noopHotels{}
	}
	return &Service{repo: r, q: q, hotels: h}
}

func (s *Service) GetByID(ctx context.Context, id string) (domain_reservations.Reservation, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain_reservations.Reservation{}, fmt.Errorf("empty id")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, hotelID, userID, status string) ([]domain_reservations.Reservation, error) {
	return s.repo.List(ctx, strings.TrimSpace(hotelID), strings.TrimSpace(userID), strings.TrimSpace(status))
}

func (s *Service) Create(ctx context.Context, r domain_reservations.Reservation) (domain_reservations.Reservation, error) {
	// Validaciones mínimas
	if strings.TrimSpace(r.HotelID) == "" ||
		strings.TrimSpace(r.UserID) == "" ||
		strings.TrimSpace(r.CheckIn) == "" ||
		strings.TrimSpace(r.CheckOut) == "" ||
		r.Guests <= 0 {
		return domain_reservations.Reservation{}, fmt.Errorf("invalid reservation payload")
	}

	// Validar que el hotel exista en hotels-api (si el cliente está disponible)
	if s.hotels != nil && !s.hotels.HotelExists(r.HotelID) {
		return domain_reservations.Reservation{}, fmt.Errorf("hotel not found")
	}

	// Orden fechas (formato YYYY-MM-DD o ISO8601 como string)
	if r.CheckOut <= r.CheckIn {
		return domain_reservations.Reservation{}, fmt.Errorf("invalid date range")
	}
	if r.Status == "" {
		r.Status = "pending"
	}
	if r.ID == "" {
		r.ID = strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	if r.CreatedAt == "" {
		r.CreatedAt = time.Now().Format(time.RFC3339)
	}

	out, err := s.repo.Create(ctx, r)
	if err != nil {
		return domain_reservations.Reservation{}, err
	}

	// Publicación opcional (no falla P1 si no hay cola)
	if s.q != nil {
		_ = s.q.PublishReservationCreated(ctx, out)
	}
	return out, nil
}

func (s *Service) Update(ctx context.Context, id string, in domain_reservations.Reservation) (domain_reservations.Reservation, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain_reservations.Reservation{}, fmt.Errorf("empty id")
	}

	// Validaciones básicas cuando vienen campos
	if in.CheckIn != "" && in.CheckOut != "" && in.CheckOut <= in.CheckIn {
		return domain_reservations.Reservation{}, fmt.Errorf("invalid date range")
	}
	if in.Guests < 0 {
		return domain_reservations.Reservation{}, fmt.Errorf("invalid guests")
	}
	if in.Status != "" {
		stat := strings.ToLower(in.Status)
		switch stat {
		case "pending", "confirmed", "cancelled":
		default:
			return domain_reservations.Reservation{}, fmt.Errorf("invalid status")
		}
	}

	return s.repo.Update(ctx, id, in)
}
