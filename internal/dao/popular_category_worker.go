package dao

import (
	"context"
	"github.com/rs/zerolog/log"
	"time"
)

const (
	popularCategoryCheckDelay = 12 * time.Hour
)

type PopularCategoryWorker struct {
	service *Service
}

func NewPopularCategoryWorker(s *Service) *PopularCategoryWorker {
	return &PopularCategoryWorker{
		service: s,
	}
}

func (w *PopularCategoryWorker) Process(ctx context.Context) error {
	for {
		err := w.service.processPopularCategory(ctx)
		if err != nil {
			log.Error().Err(err).Msg("process popular category")
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(popularCategoryCheckDelay):
		}
	}
}
