package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-core-storage/pkg/sdk/zerion"
)

const (
	fungibleUpdateDelay = time.Hour
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

		fData, err := c.zerionClient.GetFungibleData(dao.FungibleId)
		if err != nil {
			log.Error().Err(err).Msg("get fungible data")
			continue
		}

		allDaoChains := map[string]struct{}{}
		for _, strategy := range dao.Strategies {
			allDaoChains[strategy.Network] = struct{}{}
		}

		for _, chainItem := range fData.Attributes.Implementations {
			chainExInfo, ok := chainsMap[chainItem.ChainID]
			if !ok {
				log.Error().Msgf("chain %s not found in chains map", chainItem.ChainID)
				continue
			}

			cutHex, _ := strings.CutPrefix(strings.ToLower(chainExInfo.Attributes.ExternalID), "0x")
			decVal, err := strconv.ParseInt(cutHex, 16, 64)
			if err != nil {
				log.Error().Err(err).Msgf("parse chain external id %s to decimal", chainExInfo.Attributes.ExternalID)
				continue
			}
			decValStr := strconv.FormatInt(decVal, 10)

			if _, ok = allDaoChains[decValStr]; !ok {
				log.Warn().Msgf("dao %s has no strategy for chain %s", dao.Name, chainItem.ChainID)
				continue
			}

			err = c.fungibleRepo.Save(FungibleChain{
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

		time.Sleep(time.Second) // TODO: add rate limiter
	}

	return nil
}
