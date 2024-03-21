package stats

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	calcCheckDelay = time.Minute * 10
)

type CalcTotalsWorker struct {
	service *Service
}

func NewCalcTotalsWorker(s *Service) *CalcTotalsWorker {
	return &CalcTotalsWorker{
		service: s,
	}
}

func (w *CalcTotalsWorker) Start(ctx context.Context) error {
	for {
		if err := w.service.refreshTotals(); err != nil {
			log.Error().Err(err).Msg("calc totals worker")
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(calcCheckDelay):
		}
	}
}
