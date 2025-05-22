package dao

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-core-storage/pkg/sdk/zerion"
)

const (
	fungibleUpdateDelay = 2 * time.Hour
	updateExpired       = 7 * 24 * time.Hour
)

type FungibleChainWorker struct {
	service      *Service
	zerionClient *zerion.Client
	fungibleRepo *FungibleChainRepo
}

func NewFungibleChainWorker(zerionClient *zerion.Client, service *Service, fungibleRepo *FungibleChainRepo) *FungibleChainWorker {
	return &FungibleChainWorker{
		zerionClient: zerionClient,
		service:      service,
		fungibleRepo: fungibleRepo,
	}
}

func (c *FungibleChainWorker) Start(ctx context.Context) error {
	for {
		if err := c.process(); err != nil {
			log.Error().Err(err).Msg("chains cache update error")
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(fungibleUpdateDelay):
			continue
		}
	}
}

func (c *FungibleChainWorker) process() error {
	chains, err := c.zerionClient.GetChains()
	if err != nil {
		return fmt.Errorf("get chains: %w", err)
	}

	chainsMap := make(map[string]zerion.ChainData)
	for _, chain := range chains {
		chainsMap[chain.ID] = chain
	}

	filters := []Filter{FungibleIdFilter{}}
	daoList, err := c.service.GetByFilters(filters)
	if err != nil {
		return fmt.Errorf("get daos: %w", err)
	}

	for _, dao := range daoList.Daos {
		if dao.FungibleId == "" {
			log.Warn().Msg("dao has no fungible id")
			continue
		}

		needsUpdate, err := c.fungibleRepo.NeedsUpdate(dao.FungibleId, time.Now().Add(-updateExpired))
		if err != nil {
			log.Error().Err(err).Msg("check if needs update")
			continue
		}

		if !needsUpdate {
			continue
		}

		fData, err := c.zerionClient.GetFungibleData(dao.FungibleId)
		if err != nil {
			log.Error().Err(err).Msg("get fungible data")
			continue
		}

		chainItems := fData.Attributes.Implementations
		chainItems = c.filterByDaoStrategies(chainItems, dao.Strategies)

		for _, chainItem := range chainItems {
			chainExInfo, ok := chainsMap[chainItem.ChainID]
			if !ok {
				log.Error().Msgf("chain %s not found in chains map", chainItem.ChainID)
				continue
			}

			err := c.fungibleRepo.Save(FungibleChain{
				FungibleID: dao.FungibleId,
				ChainID:    chainItem.ChainID,
				ExternalID: chainExInfo.Attributes.ExternalID,
				ChainName:  chainExInfo.Attributes.Name,
				IconURL:    chainExInfo.Attributes.Icon.URL,
				Address:    chainItem.Address,
				Decimals:   chainItem.Decimals,
			})
			if err != nil {
				log.Error().Err(err).Msg("save fungible chain")
			}
		}
	}

	return nil
}

func (c *FungibleChainWorker) filterByDaoStrategies(items []zerion.Implementations, strategies Strategies) []zerion.Implementations {
	allChains := map[string]struct{}{}
	for _, strategy := range strategies {
		allChains[strings.ToLower(strategy.Network)] = struct{}{}
	}

	result := make([]zerion.Implementations, 0, len(items))
	for _, item := range items {
		_, ok := allChains[item.ChainID]
		if !ok {
			continue
		}

		result = append(result, item)
	}

	return result
}
