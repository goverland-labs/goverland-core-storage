package delegates

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type DaoProvider interface {
	GetIDByOriginalID(string) (uuid.UUID, error)
}

type Service struct {
	repo *Repo
	dp   DaoProvider
}

func NewService(repo *Repo, dp DaoProvider) *Service {
	return &Service{
		repo: repo,
		dp:   dp,
	}
}

func (s *Service) HandleDelegate(_ context.Context, hr History) error {
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

		bts, err := s.repo.GetSummaryBlockTimestamp(tx, hr.AddressFrom, daoID.String())
		if err != nil {
			return fmt.Errorf("s.repo.GetSummaryBlockTimestamp: %w", err)
		}

		// skip this block due to already processed
		if bts != 0 && bts >= hr.BlockTimestamp {
			log.Warn().Msgf("skip processing block %d from %s due to invalid timestamp", hr.BlockNumber, hr.ChainID)

			return nil
		}

		if hr.Action == actionExpire {
			if err := s.repo.UpdateSummaryExpiration(tx, hr.AddressFrom, daoID.String(), hr.Delegations.Expiration, hr.BlockTimestamp); err != nil {
				return fmt.Errorf("UpdateSummaryExpiration: %w", err)
			}

			return nil
		}

		if err := s.repo.RemoveSummary(tx, hr.AddressFrom, daoID.String()); err != nil {
			return fmt.Errorf("UpdateSummaryExpiration: %w", err)
		}

		if hr.Action == actionClear {
			return nil
		}

		for _, info := range hr.Delegations.Details {
			if err = s.repo.CreateSummary(Summary{
				AddressFrom:        hr.AddressFrom,
				AddressTo:          info.Address,
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
