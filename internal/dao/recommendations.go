package dao

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	recommendationCheckDelay = time.Hour
)

type Recommendation struct {
	OriginalId string
	InternalId string
	Name       string
	Symbol     string
	NetworkId  string
	Address    string
}

type RecommendationWorker struct {
	service *Service
}

func NewRecommendationWorker(s *Service) *RecommendationWorker {
	return &RecommendationWorker{
		service: s,
	}
}

func (w *RecommendationWorker) Process(ctx context.Context) error {
	for {
		if err := w.service.syncRecommendations(ctx); err != nil {
			log.Error().Err(err).Msg("sync recommendations")
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(recommendationCheckDelay):
		}
	}
}
