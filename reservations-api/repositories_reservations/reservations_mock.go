package repositories_reservations

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"reservations/domain_reservations"
)

type Mock struct {
	mu sync.RWMutex
	db map[string]domain_reservations.Reservation
}

func NewMock() *Mock {
	return &Mock{db: make(map[string]domain_reservations.Reservation)}
}

func (m *Mock) Create(ctx context.Context, r domain_reservations.Reservation) (domain_reservations.Reservation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if r.ID == "" {
		r.ID = strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	if r.CreatedAt == "" {
		r.CreatedAt = time.Now().Format(time.RFC3339)
	}

	// Evitar solapamientos (asumimos formato YYYY-MM-DD)
	if m.checkOverlapUnsafe(r.HotelID, r.CheckIn, r.CheckOut, "") {
		return domain_reservations.Reservation{}, fmt.Errorf("date overlap for hotel %s", r.HotelID)
	}

	m.db[r.ID] = r
	return r, nil
}

func (m *Mock) Update(ctx context.Context, id string, r domain_reservations.Reservation) (domain_reservations.Reservation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	existing, ok := m.db[id]
	if !ok {
		return domain_reservations.Reservation{}, fmt.Errorf("not found")
	}

	// aplicar cambios parciales
	if r.HotelID != "" {
		existing.HotelID = r.HotelID
	}
	if r.UserID != "" {
		existing.UserID = r.UserID
	}
	if r.CheckIn != "" {
		existing.CheckIn = r.CheckIn
	}
	if r.CheckOut != "" {
		existing.CheckOut = r.CheckOut
	}
	if r.Guests != 0 {
		existing.Guests = r.Guests
	}
	if r.RoomType != "" {
		existing.RoomType = r.RoomType
	}
	if r.TotalPrice != 0 {
		existing.TotalPrice = r.TotalPrice
	}
	if r.Status != "" {
		existing.Status = r.Status
	}

	// Evitar solapamientos (excluyendo el propio id)
	if m.checkOverlapUnsafe(existing.HotelID, existing.CheckIn, existing.CheckOut, id) {
		return domain_reservations.Reservation{}, fmt.Errorf("date overlap for hotel %s", existing.HotelID)
	}

	m.db[id] = existing
	return existing, nil
}

func (m *Mock) GetByID(ctx context.Context, id string) (domain_reservations.Reservation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	r, ok := m.db[id]
	if !ok {
		return domain_reservations.Reservation{}, fmt.Errorf("not found")
	}
	return r, nil
}

func (m *Mock) List(ctx context.Context, hotelID, userID, status string) ([]domain_reservations.Reservation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]domain_reservations.Reservation, 0, len(m.db))
	for _, v := range m.db {
		if hotelID != "" && v.HotelID != hotelID {
			continue
		}
		if userID != "" && v.UserID != userID {
			continue
		}
		if status != "" && v.Status != status {
			continue
		}
		out = append(out, v)
	}
	return out, nil
}

func (m *Mock) SeedFromJSON(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		// si no existe, no es error crítico en P1
		return nil
	}
	var list []domain_reservations.Reservation
	if err := json.Unmarshal(b, &list); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, r := range list {
		if r.ID == "" {
			r.ID = strconv.FormatInt(time.Now().UnixNano(), 10)
		}
		m.db[r.ID] = r
	}
	return nil
}

// Chequeo simple de intervalos con strings YYYY-MM-DD (orden lexicográfico coincide)
func (m *Mock) checkOverlapUnsafe(hotelID, checkIn, checkOut, excludeID string) bool {
	if hotelID == "" || checkIn == "" || checkOut == "" {
		return false
	}
	for id, existing := range m.db {
		if id == excludeID {
			continue
		}
		if existing.HotelID != hotelID {
			continue
		}
		// Overlap si !(newOut <= exIn || newIn >= exOut)
		if !(checkOut <= existing.CheckIn || checkIn >= existing.CheckOut) {
			return true
		}
	}
	return false
}
