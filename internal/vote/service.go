package vote

import (
	"context"
	"fmt"
	coreevents "github.com/goverland-labs/platform-events/events/core"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Publisher interface {
	PublishJSON(ctx context.Context, subject string, obj any) error
}

type DataProvider interface {
	BatchCreate(data []Vote) error
	GetByFilters(filters []Filter) (List, error)
}

type DaoProvider interface {
	GetIDByOriginalID(string) (uuid.UUID, error)
}

type Service struct {
	repo   DataProvider
	dao    DaoProvider
	events Publisher
}

func NewService(r DataProvider, dp DaoProvider, p Publisher) (*Service, error) {
	return &Service{
		repo:   r,
		dao:    dp,
		events: p,
	}, nil
}

func (s *Service) HandleVotes(ctx context.Context, votes []Vote) error {
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

	if err := s.repo.BatchCreate(votes); err != nil {
		return fmt.Errorf("can't create votes: %w", err)
	}

	if err := s.events.PublishJSON(ctx, coreevents.SubjectVoteCreated, convertToCoreEvent(votes)); err != nil {
		log.Error().Err(err).Msgf("publish votes event")
	}

	return nil
}

func (s *Service) GetByFilters(filters []Filter) (List, error) {
	list, err := s.repo.GetByFilters(filters)
	if err != nil {
		return List{}, fmt.Errorf("get by filters: %w", err)
	}

	return list, nil
}
