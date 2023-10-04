package dao

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	uniqueMembersCountCheckDelay = time.Hour
)

type MembersCountWorker struct {
	service *Service
}

func NewMembersCountWorker(s *Service) *MembersCountWorker {
	return &MembersCountWorker{
		service: s,
	}
}

func (w *MembersCountWorker) ProcessNew(ctx context.Context) error {
	for {
		err := w.service.processNewVoters(ctx)
		if err != nil {
			log.Error().Err(err).Msg("process new voters")
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(uniqueMembersCountCheckDelay):
		}
	}
}
