package domain_reservations

const (
	StatusPending   = "pending"
	StatusConfirmed = "confirmed"
	StatusCancelled = "cancelled"
)

type Reservation struct {
	ID       string `json:"id"`
	HotelID  string `json:"hotel_id"`
	UserID   string `json:"user_id"`
	CheckIn  string `json:"check_in"`  // YYYY-MM-DD
	CheckOut string `json:"check_out"` // YYYY-MM-DD
	Guests   int    `json:"guests"`
	Status   string `json:"status"` // "pending", "confirmed", "cancelled"
}

//Validar el estado

func IsValidStatus(s string) bool {
	switch s {
	case StatusPending, StatusConfirmed, StatusCancelled:
		return true
	default:
		return false
	}
}
