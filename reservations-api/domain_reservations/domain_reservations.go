package domain_reservations

import "time"

type Reservation struct {
	ID         string    `json:"id"`
	HotelID    string    `json:"hotel_id"`
	UserID     string    `json:"user_id"`
	CheckIn    time.Time `json:"check_in"`
	CheckOut   time.Time `json:"check_out"`
	Guests     int       `json:"guests"`
	RoomType   string    `json:"room_type"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"` // "confirmed", "cancelled", "completed"
	CreatedAt  time.Time `json:"created_at"`
}

type Repository interface {
	Create(r Reservation) (Reservation, error)
	GetByID(id string) (Reservation, error)
	GetByUserID(userID string) ([]Reservation, error)
	List() ([]Reservation, error)
	Update(id string, r Reservation) (Reservation, error)
	Delete(id string) error
	Cancel(id string) (Reservation, error)
	CheckOverlap(hotelID string, checkIn, checkOut time.Time, excludeID string) (bool, error)
	SeedFromJSON(path string) error
}

type Service interface {
	Create(r Reservation) (Reservation, error)
	GetByID(id string) (Reservation, error)
	GetByUserID(userID string) ([]Reservation, error)
	List() ([]Reservation, error)
	Update(id string, r Reservation) (Reservation, error)
	Delete(id string) error
	Cancel(id string) (Reservation, error)
}
