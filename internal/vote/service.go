package vote

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type DataProvider interface {
	BatchCreate(data []Vote) error
	GetByFilters(filters []Filter) (List, error)
}

type DaoProvider interface {
	GetIDByOriginalID(string) (uuid.UUID, error)
}

type Service struct {
	repo DataProvider
	dao  DaoProvider
}

func NewService(r DataProvider, dp DaoProvider) (*Service, error) {
	return &Service{
		repo: r,
		dao:  dp,
	}, nil
}

func (s *Service) HandleVotes(_ context.Context, votes []Vote) error {
	list := make(map[string]uuid.UUID)
	for i := range votes {
		if daoID, ok := list[votes[i].OriginalDaoID]; ok {
			votes[i].DaoID = daoID
			continue
		}

		daoID, err := s.dao.GetIDByOriginalID(votes[i].OriginalDaoID)
		if err != nil {
			log.Error().Err(err).Msgf("get dao by original id: %s", votes[i].OriginalDaoID)

			return err
		}

		list[votes[i].OriginalDaoID] = daoID
		votes[i].DaoID = daoID
	}

	return s.repo.BatchCreate(votes)
}

func (s *Service) GetByFilters(filters []Filter) (List, error) {
	list, err := s.repo.GetByFilters(filters)
	if err != nil {
		return List{}, fmt.Errorf("get by filters: %w", err)
	}

	return list, nil
}
