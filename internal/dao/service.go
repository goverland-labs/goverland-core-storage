package dao

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	coreevents "github.com/goverland-labs/platform-events/events/core"
	"github.com/jinzhu/gorm"
	"github.com/rs/zerolog/log"
)

type Publisher interface {
	PublishJSON(ctx context.Context, subject string, obj any) error
}

// todo: convert types to interfaces for unit testing
type Service struct {
	repo   *Repo
	events Publisher
}

func NewService(r *Repo, p Publisher) (*Service, error) {
	return &Service{
		repo:   r,
		events: p,
	}, nil
}

func (s *Service) HandleDao(ctx context.Context, dao Dao) error {
	existed, err := s.repo.GetByID(dao.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("handle: %w", err)
	}

	if existed == nil {
		return s.processNew(ctx, dao)
	}

	return s.processExisted(ctx, dao, *existed)
}

func (s *Service) processNew(ctx context.Context, dao Dao) error {
	err := s.repo.Create(dao)
	if err != nil {
		return fmt.Errorf("can't create dao: %w", err)
	}

	go func(dao Dao) {
		if err := s.events.PublishJSON(ctx, coreevents.SubjectDaoCreated, convertToCoreEvent(dao)); err != nil {
			log.Error().Err(err).Msgf("publish dao event #%s", dao.ID)
		}
	}(dao)

	return nil
}

func (s *Service) processExisted(ctx context.Context, new, existed Dao) error {
	equal := compare(new, existed)
	if equal {
		return nil
	}

	err := s.repo.Update(new)
	if err != nil {
		return fmt.Errorf("update dao #%s: %w", new.ID, err)
	}

	go func(dao Dao) {
		if err := s.events.PublishJSON(ctx, coreevents.SubjectDaoUpdated, convertToCoreEvent(dao)); err != nil {
			log.Error().Err(err).Msgf("publish dao event #%s", dao.ID)
		}
	}(new)

	return nil
}

func compare(d1, d2 Dao) bool {
	d1.CreatedAt = d2.CreatedAt
	d1.UpdatedAt = d2.UpdatedAt

	return reflect.DeepEqual(d1, d2)
}
