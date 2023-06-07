package dao

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	coreevents "github.com/goverland-labs/platform-events/events/core"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=mocks_test.go -package=dao . DataProvider,Publisher

type Publisher interface {
	PublishJSON(ctx context.Context, subject string, obj any) error
}

type DataProvider interface {
	Create(dao Dao) error
	Update(dao Dao) error
	GetByID(id string) (*Dao, error)
	GetByFilters(filters []Filter, count bool) (DaoList, error)
	GetCategories() ([]string, error)
}

type Service struct {
	repo   DataProvider
	events Publisher
}

func NewService(r DataProvider, p Publisher) (*Service, error) {
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

	new.CreatedAt = existed.CreatedAt
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

func (s *Service) GetByID(id string) (*Dao, error) {
	dao, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}

	return dao, nil
}

func (s *Service) GetByFilters(filters []Filter) (DaoList, error) {
	list, err := s.repo.GetByFilters(filters, true)
	if err != nil {
		return DaoList{}, fmt.Errorf("get by filters: %w", err)
	}

	return list, nil
}

// todo: use caching here!
func (s *Service) GetTopByCategories(_ context.Context, limit int) (map[string][]Dao, error) {
	categories, err := s.repo.GetCategories()
	if err != nil {
		return nil, fmt.Errorf("get categories: %w", err)
	}

	list := make(map[string][]Dao)
	for _, category := range categories {
		filters := []Filter{
			CategoryFilter{Category: category},
			PageFilter{Limit: limit, Offset: 0},
			OrderByFollowersFilter{},
		}

		data, err := s.repo.GetByFilters(filters, false)
		if err != nil {
			return nil, fmt.Errorf("get by category %s: %w", category, err)
		}

		list[category] = data.Daos
	}

	return list, nil
}
