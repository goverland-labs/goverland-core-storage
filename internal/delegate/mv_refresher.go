package delegate

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	refreshTimeout = 5 * time.Minute
)

type RefreshWorker struct {
	service *Service
}

func NewRefreshWorker(s *Service) *RefreshWorker {
	return &RefreshWorker{
		service: s,
	}
}

func (w *RefreshWorker) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(refreshTimeout):
		}

		err := w.service.refreshDelegatesMV(ctx)
		if err != nil {
			log.Error().Err(err).Msg("process available for voting")
		}
	}
}
