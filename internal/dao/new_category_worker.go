package dao

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	newCategoryCheckDelay      = 1 * time.Hour
	outdatedCategoryCheckDelay = 12 * time.Hour
)

type NewCategoryWorker struct {
	service *Service
}

func NewNewCategoryWorker(s *Service) *NewCategoryWorker {
	return &NewCategoryWorker{
		service: s,
	}
}

func (w *NewCategoryWorker) ProcessNew(ctx context.Context) error {
	for {
		err := w.service.processNewCategory(ctx)
		if err != nil {
			log.Error().Err(err).Msg("process new category")
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(newCategoryCheckDelay):
		}
	}
}

func (w *NewCategoryWorker) RemoveOutdated(ctx context.Context) error {
	for {
		err := w.service.processOutdatedNewCategory(ctx)
		if err != nil {
			log.Error().Err(err).Msg("process outdated")
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(outdatedCategoryCheckDelay):
		}
	}
}
