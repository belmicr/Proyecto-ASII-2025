package services_reservations

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	domain "reservations/domain_reservations"
	"time"
)

type EventQueue interface {
	Publish(event string) error
}

type Service struct {
	repo   domain.Repository
	events EventQueue
}

func NewService(repo domain.Repository, events EventQueue) *Service {
	return &Service{
		repo:   repo,
		events: events,
	}
}

func (s *Service) Create(r domain.Reservation) (domain.Reservation, error) {
	// Validaciones
	if err := s.validateReservation(r); err != nil {
		return domain.Reservation{}, err
	}

	// Validar que el usuario existe
	if err := s.validateUserExists(r.UserID); err != nil {
		return domain.Reservation{}, fmt.Errorf("invalid user: %w", err)
	}

	// Validar que el hotel existe
	if err := s.validateHotelExists(r.HotelID); err != nil {
		return domain.Reservation{}, fmt.Errorf("invalid hotel: %w", err)
	}

	// Crear
	created, err := s.repo.Create(r)
	if err != nil {
		return domain.Reservation{}, err
	}

	// Publicar evento
	_ = s.events.Publish(fmt.Sprintf("reservation.created:%s", created.ID))

	return created, nil
}

func (s *Service) validateReservation(r domain.Reservation) error {
	// Validar fechas
	if r.CheckIn.IsZero() {
		return errors.New("check-in date is required")
	}
	if r.CheckOut.IsZero() {
		return errors.New("check-out date is required")
	}

	// Check-in antes que check-out
	if !r.CheckIn.Before(r.CheckOut) {
		return errors.New("check-in must be before check-out")
	}

	// No reservar en el pasado (con 24h de tolerancia)
	if r.CheckIn.Before(time.Now().Add(-24 * time.Hour)) {
		return errors.New("check-in cannot be in the past")
	}

	// Validar huéspedes
	if r.Guests < 1 {
		return errors.New("at least one guest is required")
	}
	if r.Guests > 10 {
		return errors.New("maximum 10 guests allowed")
	}

	// Validar precio
	if r.TotalPrice < 0 {
		return errors.New("total price cannot be negative")
	}

	// Validar IDs
	if r.HotelID == "" {
		return errors.New("hotel_id is required")
	}
	if r.UserID == "" {
		return errors.New("user_id is required")
	}

	return nil
}

func (s *Service) validateUserExists(userID string) error {
	url := fmt.Sprintf("http://users-api:8080/users/%s", userID)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error contacting users API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return errors.New("user not found")
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("users API returned status %d", resp.StatusCode)
	}

	return nil
}

func (s *Service) validateHotelExists(hotelID string) error {
	url := fmt.Sprintf("http://hotels-api:8082/hotels/%s", hotelID)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error contacting hotels API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return errors.New("hotel not found")
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("hotels API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (s *Service) GetByID(id string) (domain.Reservation, error) {
	if id == "" {
		return domain.Reservation{}, errors.New("reservation ID is required")
	}
	return s.repo.GetByID(id)
}

func (s *Service) GetByUserID(userID string) ([]domain.Reservation, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}
	return s.repo.GetByUserID(userID)
}

func (s *Service) List() ([]domain.Reservation, error) {
	return s.repo.List()
}

func (s *Service) Update(id string, r domain.Reservation) (domain.Reservation, error) {
	if id == "" {
		return domain.Reservation{}, errors.New("reservation ID is required")
	}

	// Validar datos
	if err := s.validateReservation(r); err != nil {
		return domain.Reservation{}, err
	}

	// Actualizar
	updated, err := s.repo.Update(id, r)
	if err != nil {
		return domain.Reservation{}, err
	}

	// Publicar evento
	_ = s.events.Publish(fmt.Sprintf("reservation.updated:%s", id))

	return updated, nil
}

func (s *Service) Delete(id string) error {
	if id == "" {
		return errors.New("reservation ID is required")
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// Publicar evento
	_ = s.events.Publish(fmt.Sprintf("reservation.deleted:%s", id))

	return nil
}

func (s *Service) Cancel(id string) (domain.Reservation, error) {
	if id == "" {
		return domain.Reservation{}, errors.New("reservation ID is required")
	}

	cancelled, err := s.repo.Cancel(id)
	if err != nil {
		return domain.Reservation{}, err
	}

	// Publicar evento
	_ = s.events.Publish(fmt.Sprintf("reservation.cancelled:%s", id))

	return cancelled, nil
}
