package proposal

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	votingCheckDelay = time.Minute
)

type VotingWorker struct {
	service *Service
}

func NewVotingWorker(s *Service) *VotingWorker {
	return &VotingWorker{
		service: s,
	}
}

func (w *VotingWorker) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(votingCheckDelay):
		}

		err := w.service.processAvailableForVoting(ctx)
		if err != nil {
			log.Error().Err(err).Msg("process available for voting")
		}
	}
}
