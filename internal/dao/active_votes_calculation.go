package dao

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	activeVotesCheckDelay = time.Hour
)

type ActiveVotesWorker struct {
	service *Service
}

func NewActiveVotesWorker(s *Service) *ActiveVotesWorker {
	return &ActiveVotesWorker{
		service: s,
	}
}

func (w *ActiveVotesWorker) Process(ctx context.Context) error {
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
