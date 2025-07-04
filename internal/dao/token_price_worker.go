package dao

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	coreevents "github.com/goverland-labs/goverland-platform-events/events/core"
	"golang.org/x/exp/maps"

	"github.com/goverland-labs/goverland-core-storage/pkg/sdk/zerion"

	"github.com/rs/zerolog/log"
)

const (
	tokenPriceCheckDelay = 2 * time.Hour
	MAX_IDS_BY_REQUEST   = 25
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
		filters := []Filter{FungibleIdFilter{}}

		list, err := w.service.GetByFilters(filters)
		if err != nil {
			log.Error().Err(err).Msg("getTokenPrice")
		}
		fm := make(map[string]uuid.UUID)
		for _, d := range list.Daos {
			fm[d.FungibleId] = d.ID
		}
		if len(fm) > 0 {
			fids := maps.Keys(fm)
			idsCount := len(fids)
			for idsCount > 0 {
				ids := fids
				if idsCount > MAX_IDS_BY_REQUEST {
					ids = fids[:MAX_IDS_BY_REQUEST]
					fids = fids[MAX_IDS_BY_REQUEST:]
				}
				idsCount = idsCount - MAX_IDS_BY_REQUEST
				l, err := w.zerionClient.GetFungibleList(strings.Join(ids, ","), "")
				if err != nil {
					log.Error().Err(err).Msg("zerion client error")
					break
				}
				if err := w.service.events.PublishJSON(ctx, coreevents.DaoTokenPriceUpdated, convertToCorePaylod(l.List, fm)); err != nil {
					log.Error().Err(err).Msgf("publish token prices event")
				}
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
				Time:  time.Now(),
				Price: list[i].Attributes.MarketData.Price,
			})
		}
	}

	return res
}
