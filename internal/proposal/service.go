package proposal

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	coreevents "github.com/goverland-labs/platform-events/events/core"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/goverland-labs/core-storage/internal/metrics"
)

const (
	startVotingWindow = time.Hour
)

//go:generate mockgen -destination=mocks_test.go -package=proposal . DataProvider,Publisher,EventRegistered

type Publisher interface {
	PublishJSON(ctx context.Context, subject string, obj any) error
}

type DataProvider interface {
	Create(Proposal) error
	Update(proposal Proposal) error
	GetByID(string) (*Proposal, error)
	GetAvailableForVoting(time.Duration) ([]*Proposal, error)
	GetByFilters(filters []Filter) (ProposalList, error)
}

type EventRegistered interface {
	EventExist(_ context.Context, id, t, event string) (bool, error)
	RegisterEvent(_ context.Context, id, t, event string) error
}

// todo: convert types to interfaces for unit testing
type Service struct {
	repo      DataProvider
	publisher Publisher
	er        EventRegistered
}

func NewService(r DataProvider, p Publisher, er EventRegistered) (*Service, error) {
	return &Service{
		repo:      r,
		publisher: p,
		er:        er,
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

	go s.registerEvent(ctx, p, groupName, coreevents.SubjectProposalCreated)

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

	go s.registerEvent(ctx, new, groupName, coreevents.SubjectProposalUpdated)
	go s.checkSpecificUpdate(ctx, new, existed)

	return nil
}

func (s *Service) checkSpecificUpdate(ctx context.Context, new, existed Proposal) {
	if float64(new.ScoresTotal) >= new.Quorum {
		go s.registerEventOnce(ctx, new, groupName, coreevents.SubjectProposalVotingReached)
	}

	if new.State != existed.State {
		go s.registerEvent(ctx, new, groupName, coreevents.SubjectProposalUpdatedState)
	}
}

func (s *Service) registerEventOnce(ctx context.Context, p Proposal, group, subject string) {
	var err error
	if ok, err := s.er.EventExist(ctx, p.ID, group, subject); ok || err != nil {
		return
	}

	s.registerEvent(ctx, p, group, subject)

	if err = s.er.RegisterEvent(ctx, p.ID, group, subject); err != nil {
		log.Error().Err(err).Msgf("register event #%s", p.ID)
		return
	}
}

func (s *Service) registerEvent(ctx context.Context, p Proposal, group, subject string) {
	var err error
	defer func(group, subject string) {
		metricSendEventGauge.
			WithLabelValues(group, subject, metrics.ErrLabelValue(err)).
			Inc()
	}(group, subject)

	if err = s.publisher.PublishJSON(ctx, subject, convertToCoreEvent(p)); err != nil {
		log.Error().Err(err).Msgf("publish event #%s", p.ID)
	}
}

func compare(p1, p2 Proposal) bool {
	p1.CreatedAt = p2.CreatedAt
	p1.UpdatedAt = p2.UpdatedAt

	return reflect.DeepEqual(p1, p2)
}

// todo: parallel it
// todo: think how to move rules to the separate logic and apply all of them to the proposals
func (s *Service) processAvailableForVoting(ctx context.Context) error {
	active, err := s.repo.GetAvailableForVoting(time.Hour)
	if err != nil {
		return fmt.Errorf("get available for voting: %w", err)
	}

	for _, pr := range active {
		startsAt := time.Unix(int64(pr.Start), 0)
		endedAt := time.Unix(int64(pr.End), 0)

		// voting has started
		// do not spam with voting started if proposal was created in our system after start voting
		if pr.CreatedAt.Before(startsAt) && time.Now().After(startsAt) && time.Now().Before(endedAt) {
			go s.registerEventOnce(ctx, *pr, groupName, coreevents.SubjectProposalVotingStarted)
		}

		// voting has ended
		if pr.CreatedAt.Before(endedAt) && time.Now().After(endedAt) {
			go s.registerEventOnce(ctx, *pr, groupName, coreevents.SubjectProposalVotingEnded)
		}

		// voting is coming
		if time.Now().Sub(startsAt) < startVotingWindow {
			go s.registerEventOnce(ctx, *pr, groupName, coreevents.SubjectProposalVotingComing)
		}
	}

	return nil
}

func (s *Service) GetByID(id string) (*Proposal, error) {
	pro, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}

	return pro, nil
}

func (s *Service) GetByFilters(filters []Filter) (ProposalList, error) {
	list, err := s.repo.GetByFilters(filters)
	if err != nil {
		return ProposalList{}, fmt.Errorf("get by filters: %w", err)
	}

	return list, nil
}
