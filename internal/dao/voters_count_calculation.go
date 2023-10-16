package dao

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	uniqueVotersCountCheckDelay = time.Hour
)

type VotersCountWorker struct {
	service *Service
}

func NewVotersCountWorker(s *Service) *VotersCountWorker {
	return &VotersCountWorker{
		service: s,
	}
}

func (w *VotersCountWorker) ProcessNew(ctx context.Context) error {
	for {
		err := w.service.processNewVoters(ctx)
		if err != nil {
			log.Error().Err(err).Msg("process new voters")
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(uniqueVotersCountCheckDelay):
		}
	}
}
