package vote

import (
	"context"
	"fmt"
)

type DataProvider interface {
	BatchCreate(data []Vote) error
	GetByFilters(filters []Filter) (List, error)
}

type Service struct {
	repo DataProvider
}

func NewService(r DataProvider) (*Service, error) {
	return &Service{
		repo: r,
	}, nil
}

func (s *Service) HandleVotes(_ context.Context, votes []Vote) error {
	return s.repo.BatchCreate(votes)
}

func (s *Service) GetByFilters(filters []Filter) (List, error) {
	list, err := s.repo.GetByFilters(filters)
	if err != nil {
		return List{}, fmt.Errorf("get by filters: %w", err)
	}

	return list, nil
}
