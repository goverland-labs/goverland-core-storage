package delegate

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	coreevents "github.com/goverland-labs/goverland-platform-events/events/core"
	"github.com/rs/zerolog/log"
)

const (
	ltCheckDelay        = 30 * time.Minute
	endDelegationWindow = -6 * time.Hour
	updateAllowedDaoTTL = 15 * time.Minute
)

type LifeTimeWorker struct {
	service *Service
}

func NewLifeTimeWorker(s *Service) *LifeTimeWorker {
	return &LifeTimeWorker{
		service: s,
	}
}

func (w *LifeTimeWorker) Start(ctx context.Context) error {
	for {
		if err := w.service.checkLifeTime(ctx); err != nil {
			log.Error().Err(err).Msg("delegates lifetime check failed")
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(ltCheckDelay):
		}
	}
}

func (s *Service) checkLifeTime(ctx context.Context) error {
	offset, limit := 0, 100
	for {
		list, err := s.repo.GetDelegatesWithExpirations(offset, limit)
		if err != nil {
			return fmt.Errorf("s.repo.GetDelegatesWithExpirations: %w", err)
		}

		offset += limit

		for _, info := range list {
			endsAt := time.Unix(info.ExpiresAt, 0)

			// delegation expired
			if time.Now().After(endsAt) {
				go s.registerEventOnce(ctx, info, coreevents.SubjectDelegateDelegationExpired)
				continue
			}

			// delegation will end soon
			if time.Since(endsAt) > endDelegationWindow &&
				endsAt.After(time.Now()) {
				go s.registerEventOnce(ctx, info, coreevents.SubjectDelegateDelegationExpiringSoon)
			}
		}

		if len(list) < limit {
			break
		}
	}

	return nil
}

func (s *Service) registerEventOnce(ctx context.Context, delegate MixedDelegation, subject string) {
	group := fmt.Sprintf("delegation_%s_%s_%d_%s", subject, delegate.DaoID, delegate.LastBlockTimestamp, delegate.ProposalID)

	var err error
	if ok, err := s.er.EventExist(ctx, delegate.AddressTo, group, subject); ok || err != nil {
		return
	}

	if err = s.publisher.PublishJSON(ctx, subject, convertToCoreEvent(delegate)); err != nil {
		log.Error().Err(err).Msg("register delegate event")
	}

	if err = s.er.RegisterEvent(ctx, delegate.AddressTo, group, subject); err != nil {
		log.Error().Err(err).Msg("register delegate event")
		return
	}
}

func convertToCoreEvent(info MixedDelegation) coreevents.DelegatePayload {
	var dueDate *time.Time
	if info.ExpiresAt != 0 {
		expAt := time.Unix(info.ExpiresAt, 0)
		dueDate = &expAt
	}

	return coreevents.DelegatePayload{
		Initiator:  info.AddressTo,
		Delegator:  info.AddressFrom,
		DaoID:      uuid.MustParse(info.DaoID),
		DueDate:    dueDate,
		ProposalID: info.ProposalID,
	}
}
