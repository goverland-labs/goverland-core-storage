package proposal

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	pevents "github.com/goverland-labs/platform-events/events/aggregator"
	coreevents "github.com/goverland-labs/platform-events/events/core"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/goverland-labs/core-storage/internal/metrics"
)

const (
	startVotingWindow = -time.Hour
	endVotingWindow   = -time.Hour
)

//go:generate mockgen -destination=mocks_test.go -package=proposal . DataProvider,Publisher,EventRegistered,DaoProvider

type Publisher interface {
	PublishJSON(ctx context.Context, subject string, obj any) error
}

type DataProvider interface {
	Create(Proposal) error
	Update(proposal Proposal) error
	GetByID(string) (*Proposal, error)
	GetAvailableForVoting(time.Duration) ([]*Proposal, error)
	GetByFilters(filters []Filter) (ProposalList, error)
	GetTop(filters []Filter) (ProposalList, error)
}

type DaoProvider interface {
	GetIDByOriginalID(string) (uuid.UUID, error)
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
	dp        DaoProvider
}

func NewService(r DataProvider, p Publisher, er EventRegistered, dp DaoProvider) (*Service, error) {
	return &Service{
		repo:      r,
		publisher: p,
		er:        er,
		dp:        dp,
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

func (s *Service) HandleProposalTimeline(_ context.Context, id string, tl Timeline) error {
	pr, err := s.repo.GetByID(id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("handle: %w", err)
	}

	pr.Timeline = tl.ActualizeTimeline()

	if err := s.repo.Update(*pr); err != nil {
		return fmt.Errorf("update timeline: %w", err)
	}

	return nil
}

func (s *Service) processNew(ctx context.Context, p Proposal) error {
	daoID, err := s.dp.GetIDByOriginalID(p.DaoOriginalID)
	if err != nil {
		return fmt.Errorf("get dao by name: %s: %w", p.DaoOriginalID, err)
	}

	p.DaoID = daoID
	err = s.repo.Create(p)
	if err != nil {
		return fmt.Errorf("create proposal: %w", err)
	}

	s.registerEvent(ctx, p, groupName, coreevents.SubjectProposalCreated)

	if err = s.publisher.PublishJSON(ctx, coreevents.SubjectCheckActivitySince, pevents.DaoPayload{ID: p.DaoOriginalID}); err != nil {
		log.Error().Err(err).Msgf("publish dao event #%s", daoID.String())
	}

	return nil
}

func (s *Service) processExisted(ctx context.Context, new, existed Proposal) error {
	equal := compare(new, existed)
	if equal {
		return nil
	}

	new.DaoID = existed.DaoID
	new.CreatedAt = existed.CreatedAt
	err := s.repo.Update(new)
	if err != nil {
		return fmt.Errorf("update proposal #%s: %w", new.ID, err)
	}

	s.registerEvent(ctx, new, groupName, coreevents.SubjectProposalUpdated)
	s.checkSpecificUpdate(ctx, new, existed)

	return nil
}

func (s *Service) checkSpecificUpdate(ctx context.Context, new, existed Proposal) {
	if new.Quorum > 0 && float64(new.ScoresTotal) >= new.Quorum {
		go s.registerEventOnce(ctx, new, groupName, coreevents.SubjectProposalVotingQuorumReached)
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
	p1.DaoID = p2.DaoID
	p1.DaoOriginalID = p2.DaoOriginalID

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
		endsAt := time.Unix(int64(pr.End), 0)

		// voting has started
		if time.Now().After(startsAt) && time.Now().Before(endsAt) {
			go s.registerEventOnce(ctx, *pr, groupName, coreevents.SubjectProposalVotingStarted)
		}

		// voting has ended
		if time.Now().After(endsAt) {
			go s.registerEventOnce(ctx, *pr, groupName, coreevents.SubjectProposalVotingEnded)
		}

		// voting will start soon
		if time.Since(startsAt) > startVotingWindow && startsAt.After(time.Now()) {
			go s.registerEventOnce(ctx, *pr, groupName, coreevents.SubjectProposalVotingStartsSoon)
		}

		// voting will end soon
		if time.Since(endsAt) > endVotingWindow &&
			startsAt.Before(time.Now()) &&
			endsAt.After(time.Now()) {
			go s.registerEventOnce(ctx, *pr, groupName, coreevents.SubjectProposalVotingEndsSoon)
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

func (s *Service) GetTop(filters []Filter) (ProposalList, error) {
	list, err := s.repo.GetTop(filters)
	if err != nil {
		return ProposalList{}, fmt.Errorf("get top: %w", err)
	}

	return list, nil
}
