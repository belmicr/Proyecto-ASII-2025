package services_hotels

import (
	"context"

	"hotels/domain_hotels"
)

// El repo que usa el service (lo satisface Mock y el repo Mongo)
type Repository interface {
	Create(ctx context.Context, h domain_hotels.Hotel) (domain_hotels.Hotel, error)
	Update(ctx context.Context, id string, h domain_hotels.Hotel) (domain_hotels.Hotel, error)
	GetByID(ctx context.Context, id string) (domain_hotels.Hotel, error)
	List(ctx context.Context, q string) ([]domain_hotels.Hotel, error)
}

// La cola de eventos (para ahora puede ser un stub que no haga nada)
type Events interface {
	Publish(event any) error
}

type Service struct {
	repo Repository
	ev   Events
}

func NewService(r Repository, e Events) *Service { return &Service{repo: r, ev: e} }

func (s *Service) Create(ctx context.Context, h domain_hotels.Hotel) (domain_hotels.Hotel, error) {
	out, err := s.repo.Create(ctx, h)
	if err == nil {
		_ = s.ev.Publish(map[string]any{"op": "create", "id": out.ID})
	}
	return out, err
}

func (s *Service) Update(ctx context.Context, id string, h domain_hotels.Hotel) (domain_hotels.Hotel, error) {
	out, err := s.repo.Update(ctx, id, h)
	if err == nil {
		_ = s.ev.Publish(map[string]any{"op": "update", "id": id})
	}
	return out, err
}

func (s *Service) GetByID(ctx context.Context, id string) (domain_hotels.Hotel, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context, q string) ([]domain_hotels.Hotel, error) {
	return s.repo.List(ctx, q)
}
