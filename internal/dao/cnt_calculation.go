package dao

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	activeVotesCheckDelay = time.Hour
	cntProposalsCheckDelay = 3 * time.Hour
)

type CntCalculationWorker struct {
	service *Service
}

func NewCntCalculationWorker(s *Service) *CntCalculationWorker {
	return &CntCalculationWorker{
		service: s,
	}
}

func (w *CntCalculationWorker) ProcessVotes(ctx context.Context) error {
	for {
		err := w.service.processActiveVotes(ctx)
		if err != nil {
			log.Error().Err(err).Msg("process active votes")
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(activeVotesCheckDelay):
		}
	}
}

func (w *CntCalculationWorker) ProcessProposalCounters(ctx context.Context) error {
	for {
		err := w.service.processProposalsCnt(ctx)
		if err != nil {
			log.Error().Err(err).Msg("process active votes")
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(cntProposalsCheckDelay):
		}
	}
}
