package services_reservations

import (
	"context"
	"reservations/domain_reservations"
)

type Repository interface {
	Create(ctx context.Context, r domain_reservations.Reservation) (domain_reservations.Reservation, error)
	Update(ctx context.Context, id string, r domain_reservations.Reservation) (domain_reservations.Reservation, error)
	GetByID(ctx context.Context, id string) (domain_reservations.Reservation, error)
	List(ctx context.Context, hotelID, userID, status string) ([]domain_reservations.Reservation, error)
}

type Events interface {
	Publish(event any) error
}

type Service struct {
	repo Repository
	ev   Events
}

func NewService(r Repository, e Events) *Service { return &Service{repo: r, ev: e} }

func (s *Service) Create(ctx context.Context, r domain_reservations.Reservation) (domain_reservations.Reservation, error) {
	out, err := s.repo.Create(ctx, r)
	if err == nil {
		_ = s.ev.Publish(map[string]any{"op": "create", "id": out.ID})
	}
	return out, err
}

func (s *Service) Update(ctx context.Context, id string, r domain_reservations.Reservation) (domain_reservations.Reservation, error) {
	out, err := s.repo.Update(ctx, id, r)
	if err == nil {
		_ = s.ev.Publish(map[string]any{"op": "update", "id": id})
	}
	return out, err
}

func (s *Service) GetByID(ctx context.Context, id string) (domain_reservations.Reservation, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, hotelID, userID, status string) ([]domain_reservations.Reservation, error) {
	return s.repo.List(ctx, hotelID, userID, status)
}
