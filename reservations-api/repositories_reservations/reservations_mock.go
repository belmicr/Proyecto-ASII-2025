package repositories_reservations

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"reservations/domain_reservations"
)

type Mock struct {
	mu sync.RWMutex
	db map[string]domain_reservations.Reservation
}

func NewMock() *Mock { return &Mock{db: map[string]domain_reservations.Reservation{}} }

func (m *Mock) Create(ctx context.Context, r domain_reservations.Reservation) (domain_reservations.Reservation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if r.ID == "" {
		r.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	m.db[r.ID] = r
	return r, nil
}

func (m *Mock) Update(ctx context.Context, id string, r domain_reservations.Reservation) (domain_reservations.Reservation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.db[id]; !ok {
		return domain_reservations.Reservation{}, fmt.Errorf("not found")
	}
	r.ID = id
	m.db[id] = r
	return r, nil
}

func (m *Mock) GetByID(ctx context.Context, id string) (domain_reservations.Reservation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.db[id]
	if !ok {
		return domain_reservations.Reservation{}, fmt.Errorf("not found")
	}
	return v, nil
}

// Filtros básicos por hotel_id, user_id, status
func (m *Mock) List(ctx context.Context, hotelID, userID, status string) ([]domain_reservations.Reservation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	trim := func(s string) string { return strings.TrimSpace(strings.ToLower(s)) }
	hotelID = trim(hotelID)
	userID = trim(userID)
	status = trim(status)

	out := make([]domain_reservations.Reservation, 0, len(m.db))
	for _, v := range m.db {
		if hotelID != "" && strings.ToLower(v.HotelID) != hotelID {
			continue
		}
		if userID != "" && strings.ToLower(v.UserID) != userID {
			continue
		}
		if status != "" && strings.ToLower(v.Status) != status {
			continue
		}
		out = append(out, v)
	}
	return out, nil
}

// Semilla opcional desde db/reservations.json
func (m *Mock) SeedFromJSON(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var items []domain_reservations.Reservation
	if err := json.NewDecoder(f).Decode(&items); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	for _, r := range items {
		if r.ID == "" {
			continue
		}
		m.db[r.ID] = r
	}
	return nil
}
