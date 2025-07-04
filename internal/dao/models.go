package dao

import (
	"time"

	"github.com/google/uuid"
	aggevents "github.com/goverland-labs/goverland-platform-events/events/aggregator"
	events "github.com/goverland-labs/goverland-platform-events/events/core"
)

type Treasury struct {
	Name    string
	Address string
	Network string
}

type Treasuries []Treasury

func convertToTreasures(list Treasuries) []events.TreasuryPayload {
	res := make([]events.TreasuryPayload, len(list))
	for i, treasury := range list {
		res[i] = events.TreasuryPayload{
			Name:    treasury.Name,
			Address: treasury.Address,
			Network: treasury.Network,
		}
	}

	return res
}

type Categories []string

type Strategy struct {
	Name    string
	Network string
	Params  map[string]interface{}
}

type Strategies []Strategy

func convertToStrategies(list Strategies) []events.StrategyPayload {
	result := make([]events.StrategyPayload, len(list))
	for i, strategy := range list {
		result[i] = events.StrategyPayload{
			Name:    strategy.Name,
			Network: strategy.Network,
		}
	}

	return result
}

type Voting struct {
	Delay       int
	Period      int
	Type        string
	Quorum      float32
	Blind       bool
	HideAbstain bool
	Privacy     string
	Aliased     bool
}

func convertToVoting(v Voting) events.VotingPayload {
	return events.VotingPayload{
		Delay:       v.Delay,
		Period:      v.Period,
		Type:        v.Type,
		Quorum:      v.Quorum,
		Blind:       v.Blind,
		HideAbstain: v.HideAbstain,
		Privacy:     v.Privacy,
		Aliased:     v.Aliased,
	}
}

type Dao struct {
	ID              uuid.UUID `gorm:"primary_key"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	OriginalID      string
	Name            string
	Private         bool
	About           string
	Avatar          string
	Terms           string
	Location        string
	Website         string
	Twitter         string
	Github          string
	Coingecko       string
	Email           string
	Network         string
	Symbol          string
	Skin            string
	Domain          string
	Strategies      Strategies `gorm:"serializer:json"`
	Voting          Voting     `gorm:"serializer:json"`
	Categories      Categories `gorm:"serializer:json"`
	Treasures       Treasuries `gorm:"serializer:json"`
	FollowersCount  int
	ProposalsCount  int
	Guidelines      string
	Template        string
	ParentID        *uuid.UUID
	ActivitySince   int
	VotersCount     int
	PopularityIndex float64
	// ActiveVotes the number of active proposals
	ActiveVotes int
	// ActiveProposalsIDs the list of active proposals identifiers
	ActiveProposalsIDs  []string `gorm:"serializer:json"`
	Verified            bool
	FungibleId          string
	TokenSymbol         string
	VerificationStatus  string
	VerificationComment string
}

func convertToCoreEvent(dao Dao) events.DaoPayload {
	return events.DaoPayload{
		ID:                 dao.ID,
		CreatedAt:          dao.CreatedAt,
		Alias:              dao.OriginalID,
		Name:               dao.Name,
		Private:            dao.Private,
		About:              dao.About,
		Avatar:             dao.Avatar,
		Terms:              dao.Terms,
		Location:           dao.Location,
		Website:            dao.Website,
		Twitter:            dao.Twitter,
		Github:             dao.Github,
		Coingecko:          dao.Coingecko,
		Email:              dao.Email,
		Network:            dao.Network,
		Symbol:             dao.Symbol,
		Skin:               dao.Skin,
		Domain:             dao.Domain,
		Strategies:         convertToStrategies(dao.Strategies),
		Voting:             convertToVoting(dao.Voting),
		Categories:         dao.Categories,
		Treasures:          convertToTreasures(dao.Treasures),
		FollowersCount:     dao.FollowersCount,
		ProposalsCount:     dao.ProposalsCount,
		Guidelines:         dao.Guidelines,
		Template:           dao.Template,
		ParentID:           dao.ParentID,
		ActiveSince:        &dao.ActivitySince,
		Verified:           dao.Verified,
		PopularityIndex:    &dao.PopularityIndex,
		VerificationStatus: dao.VerificationStatus,
	}
}

func convertToDao(e aggevents.DaoPayload) Dao {
	return Dao{
		OriginalID:     e.ID,
		Name:           e.Name,
		Private:        e.Private,
		About:          e.About,
		Avatar:         e.Avatar,
		Terms:          e.Terms,
		Location:       e.Location,
		Website:        e.Website,
		Twitter:        e.Twitter,
		Github:         e.Github,
		Coingecko:      e.Coingecko,
		Email:          e.Email,
		Network:        e.Network,
		Symbol:         e.Symbol,
		Skin:           e.Skin,
		Domain:         e.Domain,
		Strategies:     convertToInternalStrategies(e.Strategies),
		Voting:         convertToInternalVoting(e.Voting),
		Categories:     e.Categories,
		Treasures:      convertToInternalTreasures(e.Treasures),
		FollowersCount: e.FollowersCount,
		ProposalsCount: e.ProposalsCount,
		Guidelines:     e.Guidelines,
		Template:       e.Template,
		Verified:       e.Verified,

		// TODO: parentID
	}
}

func convertToInternalStrategies(s []aggevents.StrategyPayload) Strategies {
	res := make(Strategies, len(s))
	for i, item := range s {
		res[i] = Strategy{
			Name:    item.Name,
			Network: item.Network,
			Params:  item.Params,
		}
	}

	return res
}

func convertToInternalVoting(v aggevents.VotingPayload) Voting {
	return Voting{
		Delay:       v.Delay,
		Period:      v.Period,
		Type:        v.Type,
		Quorum:      v.Quorum,
		Blind:       v.Blind,
		HideAbstain: v.HideAbstain,
		Privacy:     v.Privacy,
		Aliased:     v.Aliased,
	}
}

func convertToInternalTreasures(list []aggevents.TreasuryPayload) Treasuries {
	res := make(Treasuries, len(list))
	for i, item := range list {
		res[i] = Treasury{
			Name:    item.Name,
			Address: item.Address,
			Network: item.Network,
		}
	}

	return res
}

type TokenInfo struct {
	Name                  string
	Symbol                string
	TotalSupply           float64
	CirculatingSupply     float64
	MarketCap             float64
	FullyDilutedValuation float64
	Price                 float64
	FungibleID            string
	Chains                []TokenChainInfo
}

type TokenChainInfo struct {
	ChainID  string
	Name     string
	Decimals int
	IconURL  string
	Address  string
}
