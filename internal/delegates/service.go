package delegates

import (
	"context"
	"errors"
	"fmt"
	"strings"

	events "github.com/goverland-labs/goverland-platform-events/events/core"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type DaoProvider interface {
	GetIDByOriginalID(string) (uuid.UUID, error)
}

type Publisher interface {
	PublishJSON(ctx context.Context, subject string, obj any) error
}

type Service struct {
	repo      *Repo
	dp        DaoProvider
	publisher Publisher
}

func NewService(repo *Repo, dp DaoProvider, p Publisher) *Service {
	return &Service{
		repo:      repo,
		dp:        dp,
		publisher: p,
	}
}

func (s *Service) handleDelegate(_ context.Context, hr History) error {
	if err := s.repo.CallInTx(func(tx *gorm.DB) error {
		if hr.OriginalSpaceID == "" {
			log.Warn().Msgf("skip processing block %d from %s cause dao id is empty", hr.BlockNumber, hr.ChainID)

			return nil
		}

		// store to history
		if err := s.repo.CreateHistory(tx, hr); err != nil {
			return fmt.Errorf("repo.CreateHistory: %w", err)
		}

		// get space id by provided original_space_id
		daoID, err := s.dp.GetIDByOriginalID(hr.OriginalSpaceID)
		if err != nil {
			return fmt.Errorf("dp.GetIDByOriginalID: %w", err)
		}

		bts, err := s.repo.GetSummaryBlockTimestamp(tx, strings.ToLower(hr.AddressFrom), daoID.String())
		if err != nil {
			return fmt.Errorf("s.repo.GetSummaryBlockTimestamp: %w", err)
		}

		// skip this block due to already processed
		if bts != 0 && bts >= hr.BlockTimestamp {
			log.Warn().Msgf("skip processing block %d from %s due to invalid timestamp", hr.BlockNumber, hr.ChainID)

			return nil
		}

		if hr.Action == actionExpire {
			if err := s.repo.UpdateSummaryExpiration(tx, strings.ToLower(hr.AddressFrom), daoID.String(), hr.Delegations.Expiration, hr.BlockTimestamp); err != nil {
				return fmt.Errorf("UpdateSummaryExpiration: %w", err)
			}

			return nil
		}

		if err := s.repo.RemoveSummary(tx, strings.ToLower(hr.AddressFrom), daoID.String()); err != nil {
			return fmt.Errorf("UpdateSummaryExpiration: %w", err)
		}

		if hr.Action == actionClear {
			return nil
		}

		for _, info := range hr.Delegations.Details {
			if err = s.repo.CreateSummary(Summary{
				AddressFrom:        strings.ToLower(hr.AddressFrom),
				AddressTo:          strings.ToLower(info.Address),
				DaoID:              daoID.String(),
				Weight:             info.Weight,
				LastBlockTimestamp: hr.BlockTimestamp,
				ExpiresAt:          int64(hr.Delegations.Expiration),
			}); err != nil {
				return fmt.Errorf("createSummary [%s/%s/%s]: %w", hr.AddressFrom, info.Address, daoID.String(), err)
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("repo.CallInTx: %w", err)
	}

	return nil
}

func (s *Service) handleProposalCreated(ctx context.Context, pr Proposal) error {
	// get space id by provided original_space_id
	daoID, err := s.dp.GetIDByOriginalID(pr.OriginalDaoID)
	if err != nil {
		return fmt.Errorf("dp.GetIDByOriginalID: %w", err)
	}

	// find delegator by author in specific space id
	summary, err := s.repo.FindDelegator(daoID.String(), strings.ToLower(pr.Author))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("repo.FindDelegator: %w", err)
	}

	// author doesn't have any delegation relations
	if summary == nil {
		return nil
	}

	if summary.SelfDelegation() {
		return nil
	}

	// delegation is expired
	if summary.Expired() {
		return nil
	}

	// make an event
	event := events.DelegatePayload{
		Initiator:  strings.ToLower(pr.Author),
		Delegator:  summary.AddressFrom,
		DaoID:      daoID,
		ProposalID: pr.ID,
	}
	if err = s.publisher.PublishJSON(ctx, events.SubjectDelegateCreateProposal, event); err != nil {
		return fmt.Errorf("s.publisher.PublishJSON: %w", err)
	}

	return nil
}

func (s *Service) handleVotesCreated(ctx context.Context, batch []Vote) error {
	summary, err := s.repo.FindDelegatorsByVotes(batch)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("repo.FindDelegatorsByVotes: %w", err)
	}

	for _, info := range summary {
		if info.SelfDelegation() {
			continue
		}

		// delegation is expired
		if info.Expired() {
			continue
		}

		// make an event
		event := events.DelegatePayload{
			Initiator:  strings.ToLower(info.AddressTo),
			Delegator:  info.AddressFrom,
			DaoID:      uuid.MustParse(info.DaoID),
			ProposalID: info.ProposalID,
		}

		if err = s.publisher.PublishJSON(ctx, events.SubjectDelegateVotingVoted, event); err != nil {
			log.Err(err).Msgf("publish delegate voted: %s %s", info.AddressTo, info.ProposalID)
		}
	}

	return nil
}
