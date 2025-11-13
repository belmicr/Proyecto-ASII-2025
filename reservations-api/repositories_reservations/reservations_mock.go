package repositories_reservations

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	domain "reservations/domain_reservations"
	"sync"
	"time"
)

type Mock struct {
	data   map[string]domain.Reservation
	nextID int
	mu     sync.RWMutex
}

func NewMock() *Mock {
	return &Mock{
		data:   make(map[string]domain.Reservation),
		nextID: 1,
	}
}

// SeedFromJSON carga datos iniciales desde un archivo JSON
func (m *Mock) SeedFromJSON(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	file, err := os.ReadFile(path)
	if err != nil {
		// Si el archivo no existe, no es error crítico
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("error reading seed file: %w", err)
	}

	var seedData struct {
		Reservations []domain.Reservation `json:"reservations"`
	}

	if err := json.Unmarshal(file, &seedData); err != nil {
		return fmt.Errorf("error parsing seed JSON: %w", err)
	}

	// Cargar reservas
	for _, res := range seedData.Reservations {
		m.data[res.ID] = res

		// Actualizar nextID
		var numID int
		fmt.Sscanf(res.ID, "%d", &numID)
		if numID >= m.nextID {
			m.nextID = numID + 1
		}
	}

	return nil
}

func (m *Mock) Create(r domain.Reservation) (domain.Reservation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Validar solapamiento
	hasOverlap, err := m.checkOverlapUnsafe(r.HotelID, r.CheckIn, r.CheckOut, "")
	if err != nil {
		return domain.Reservation{}, err
	}
	if hasOverlap {
		return domain.Reservation{}, errors.New("las fechas se solapan con una reserva existente")
	}

	// Generar ID
	r.ID = fmt.Sprintf("%d", m.nextID)
	m.nextID++

	// Establecer valores por defecto
	if r.Status == "" {
		r.Status = "confirmed"
	}
	r.CreatedAt = time.Now()

	m.data[r.ID] = r
	return r, nil
}

func (m *Mock) GetByID(id string) (domain.Reservation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res, exists := m.data[id]
	if !exists {
		return domain.Reservation{}, errors.New("reservation not found")
	}
	return res, nil
}

func (m *Mock) GetByUserID(userID string) ([]domain.Reservation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []domain.Reservation
	for _, res := range m.data {
		if res.UserID == userID {
			result = append(result, res)
		}
	}
	return result, nil
}

func (m *Mock) List() ([]domain.Reservation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]domain.Reservation, 0, len(m.data))
	for _, res := range m.data {
		result = append(result, res)
	}
	return result, nil
}

func (m *Mock) Update(id string, r domain.Reservation) (domain.Reservation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	existing, exists := m.data[id]
	if !exists {
		return domain.Reservation{}, errors.New("reservation not found")
	}

	// No permitir modificar reservas canceladas
	if existing.Status == "cancelled" {
		return domain.Reservation{}, errors.New("cannot modify cancelled reservation")
	}

	// Validar solapamiento (excluyendo la reserva actual)
	hasOverlap, err := m.checkOverlapUnsafe(r.HotelID, r.CheckIn, r.CheckOut, id)
	if err != nil {
		return domain.Reservation{}, err
	}
	if hasOverlap {
		return domain.Reservation{}, errors.New("las fechas se solapan con otra reserva")
	}

	// Mantener ID y CreatedAt
	r.ID = id
	r.CreatedAt = existing.CreatedAt

	m.data[id] = r
	return r, nil
}

func (m *Mock) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[id]; !exists {
		return errors.New("reservation not found")
	}

	delete(m.data, id)
	return nil
}

func (m *Mock) Cancel(id string) (domain.Reservation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	res, exists := m.data[id]
	if !exists {
		return domain.Reservation{}, errors.New("reservation not found")
	}

	if res.Status == "cancelled" {
		return domain.Reservation{}, errors.New("reservation already cancelled")
	}

	// Validar que no haya empezado
	if res.CheckIn.Before(time.Now()) {
		return domain.Reservation{}, errors.New("cannot cancel reservation that has already started")
	}

	res.Status = "cancelled"
	m.data[id] = res
	return res, nil
}

func (m *Mock) CheckOverlap(hotelID string, checkIn, checkOut time.Time, excludeID string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.checkOverlapUnsafe(hotelID, checkIn, checkOut, excludeID)
}

// checkOverlapUnsafe debe llamarse con el mutex ya tomado
func (m *Mock) checkOverlapUnsafe(hotelID string, checkIn, checkOut time.Time, excludeID string) (bool, error) {
	for id, existing := range m.data {
		// Excluir la reserva actual (para updates) y reservas canceladas
		if id == excludeID || existing.Status == "cancelled" {
			continue
		}

		// Solo verificar mismo hotel
		if existing.HotelID != hotelID {
			continue
		}

		// Verificar solapamiento
		// Caso 1: Nueva empieza durante existente
		if (checkIn.After(existing.CheckIn) || checkIn.Equal(existing.CheckIn)) && checkIn.Before(existing.CheckOut) {
			return true, nil
		}

		// Caso 2: Nueva termina durante existente
		if checkOut.After(existing.CheckIn) && (checkOut.Before(existing.CheckOut) || checkOut.Equal(existing.CheckOut)) {
			return true, nil
		}

		// Caso 3: Nueva contiene completamente a existente
		if (checkIn.Before(existing.CheckIn) || checkIn.Equal(existing.CheckIn)) &&
			(checkOut.After(existing.CheckOut) || checkOut.Equal(existing.CheckOut)) {
			return true, nil
		}
	}

	return false, nil
}
