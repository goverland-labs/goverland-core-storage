package dao

import (
	"context"
	"github.com/google/uuid"
	"github.com/goverland-labs/goverland-core-storage/pkg/sdk/zerion"
	coreevents "github.com/goverland-labs/goverland-platform-events/events/core"
	"golang.org/x/exp/maps"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	tokenPriceCheckDelay = 24 * time.Hour
)

type TokenPriceWorker struct {
	service      *Service
	zerionClient *zerion.Client
}

func NewTokenPriceWorker(s *Service, z *zerion.Client) *TokenPriceWorker {
	return &TokenPriceWorker{
		service:      s,
		zerionClient: z,
	}
}

func (w *TokenPriceWorker) Process(ctx context.Context) error {
	for {
		filters := []Filter{VerifiedFilter{}}

		list, err := w.service.GetByFilters(filters)
		if err != nil {
			log.Error().Err(err).Msg("getTokenPrice")
		}
		fm := make(map[string]uuid.UUID)
		for _, d := range list.Daos {
			if d.FungibleId != "" {
				fm[d.FungibleId] = d.ID
			}
		}
		if len(fm) > 0 {
			l, err := w.zerionClient.GetTokenPrices(strings.Join(maps.Keys(fm), ","))
			if err != nil {
				log.Error().Err(err).Msg("zerion client error")
			}
			if err := w.service.events.PublishJSON(ctx, coreevents.DaoTokenPriceUpdated, convertToCorePaylod(l.List, fm)); err != nil {
				log.Error().Err(err).Msgf("publish token prices event")
			}

		}
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(tokenPriceCheckDelay):
		}
	}
}

func convertToCorePaylod(list []zerion.FungibleData, fungiblesMap map[string]uuid.UUID) coreevents.TokenPricesPayload {
	res := make(coreevents.TokenPricesPayload, 0, len(list))
	for i := range list {
		daoId, exist := fungiblesMap[list[i].ID]
		if exist {
			res = append(res, coreevents.TokenPricePayload{
				DaoID: daoId,
				Price: list[i].Attributes.MarketData.Price,
			})
		}
	}

	return res
}
