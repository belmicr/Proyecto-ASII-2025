package domain_reservations

import "context"

// NOTA: usamos string para CheckIn/CheckOut porque en los ejemplos
// mandamos "YYYY-MM-DD". Si querés time.Time después, lo cambiamos
// junto con el binding y el parseo en el service.
type Reservation struct {
	ID         string  `json:"id"`
	HotelID    string  `json:"hotel_id"`
	UserID     string  `json:"user_id"`
	CheckIn    string  `json:"check_in"`
	CheckOut   string  `json:"check_out"`
	Guests     int     `json:"guests"`
	RoomType   string  `json:"room_type,omitempty"`
	TotalPrice float64 `json:"total_price,omitempty"`
	Status     string  `json:"status,omitempty"`     // p.ej: pending, confirmed, cancelled
	CreatedAt  string  `json:"created_at,omitempty"` // ISO string opcional
}

// Interfaz que implementan repositories (mock/mysql)
type Repository interface {
	Create(ctx context.Context, r Reservation) (Reservation, error)
	Update(ctx context.Context, id string, r Reservation) (Reservation, error)
	GetByID(ctx context.Context, id string) (Reservation, error)
	List(ctx context.Context, hotelID, userID, status string) ([]Reservation, error)
	SeedFromJSON(path string) error
}

// Opcional: si publicamos eventos (Rabbit/Noop)
type EventQueue interface {
	PublishReservationCreated(ctx context.Context, r Reservation) error
}

// Interfaz que expone el servicio a los controllers
type Service interface {
	GetByID(ctx context.Context, id string) (Reservation, error)
	List(ctx context.Context, hotelID, userID, status string) ([]Reservation, error)
	Create(ctx context.Context, r Reservation) (Reservation, error)
	Update(ctx context.Context, id string, r Reservation) (Reservation, error)
}
