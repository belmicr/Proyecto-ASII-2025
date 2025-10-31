package repositories_hotels

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"hotels/domain_hotels"
)

type Mock struct {
	mu sync.RWMutex
	db map[string]domain_hotels.Hotel
}

func NewMock() *Mock { return &Mock{db: map[string]domain_hotels.Hotel{}} }

func (m *Mock) Create(ctx context.Context, h domain_hotels.Hotel) (domain_hotels.Hotel, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if h.ID == "" {
		h.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	m.db[h.ID] = h
	return h, nil
}

func (m *Mock) Update(ctx context.Context, id string, h domain_hotels.Hotel) (domain_hotels.Hotel, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.db[id]; !ok {
		return domain_hotels.Hotel{}, fmt.Errorf("not found")
	}
	h.ID = id
	m.db[id] = h
	return h, nil
}

func (m *Mock) GetByID(ctx context.Context, id string) (domain_hotels.Hotel, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	h, ok := m.db[id]
	if !ok {
		return domain_hotels.Hotel{}, fmt.Errorf("not found")
	}
	return h, nil
}

func (m *Mock) List(ctx context.Context, q string) ([]domain_hotels.Hotel, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Si no hay query, devolvemos todo
	if strings.TrimSpace(q) == "" {
		out := make([]domain_hotels.Hotel, 0, len(m.db))
		for _, v := range m.db {
			out = append(out, v)
		}
		return out, nil
	}

	// Filtro simple por nombre o ciudad (case-insensitive)
	q = strings.ToLower(q)
	out := make([]domain_hotels.Hotel, 0, len(m.db))
	for _, v := range m.db {
		if strings.Contains(strings.ToLower(v.Name), q) || strings.Contains(strings.ToLower(v.City), q) {
			out = append(out, v)
		}
	}
	return out, nil
}

// SeedFromJSON carga hoteles desde un archivo JSON (ej: "db/hotels.json")
func (m *Mock) SeedFromJSON(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var items []domain_hotels.Hotel
	if err := json.NewDecoder(f).Decode(&items); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	for _, h := range items {
		if h.ID == "" {
			continue
		}
		m.db[h.ID] = h
	}
	return nil
}
