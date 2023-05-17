package vote

import (
	"context"
)

type DataProvider interface {
	BatchCreate(data []Vote) error
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
