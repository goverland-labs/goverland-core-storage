package proposal

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	coreevents "github.com/goverland-labs/platform-events/events/core"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=mocks_test.go -package=proposal . DataProvider,Publisher

type Publisher interface {
	PublishJSON(ctx context.Context, subject string, obj any) error
}

type DataProvider interface {
	Create(Proposal) error
	Update(proposal Proposal) error
	GetByID(string) (*Proposal, error)
}

// todo: convert types to interfaces for unit testing
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

func (s *Service) HandleProposal(ctx context.Context, pro Proposal) error {
	existed, err := s.repo.GetByID(pro.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("handle: %w", err)
	}

	if existed == nil {
		return s.processNew(ctx, pro)
	}

	return s.processExisted(ctx, pro, *existed)
}

func (s *Service) processNew(ctx context.Context, p Proposal) error {
	err := s.repo.Create(p)
	if err != nil {
		return fmt.Errorf("can't create proposal: %w", err)
	}

	go func(p Proposal) {
		if err := s.events.PublishJSON(ctx, coreevents.SubjectProposalCreated, convertToCoreEvent(p)); err != nil {
			log.Error().Err(err).Msgf("publish proposal event #%s", p.ID)
		}
	}(p)

	return nil
}

func (s *Service) processExisted(ctx context.Context, new, existed Proposal) error {
	equal := compare(new, existed)
	if equal {
		return nil
	}

	new.CreatedAt = existed.CreatedAt
	err := s.repo.Update(new)
	if err != nil {
		return fmt.Errorf("update proposal #%s: %w", new.ID, err)
	}

	go func(p Proposal) {
		if err := s.events.PublishJSON(ctx, coreevents.SubjectProposalUpdated, convertToCoreEvent(p)); err != nil {
			log.Error().Err(err).Msgf("publish proposal event #%s", p.ID)
		}
	}(new)

	return nil
}

func compare(p1, p2 Proposal) bool {
	p1.CreatedAt = p2.CreatedAt
	p1.UpdatedAt = p2.UpdatedAt

	return reflect.DeepEqual(p1, p2)
}
