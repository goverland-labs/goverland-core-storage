package delegate

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	events "github.com/goverland-labs/goverland-platform-events/events/aggregator"
	"github.com/rs/zerolog/log"
)

const (
	actionClear  = "clear"
	actionExpire = "expire"
)

const (
	sourceSplitDelegation = "split-delegation"
	sourceErc20Votes      = "erc20-votes"
)

var (
	daoErc20Set = map[string]Erc20Mapping{
		strings.ToLower("0x912ce59144191c1204e64559fe8253a0e49e6548"): {
			OriginalID: "arbitrumfoundation.eth",
			ChainID:    "42161",
		},
		strings.ToLower("0xCa14007Eff0dB1f8135f4C25B34De49AB0d42766"): {
			OriginalID: "starknet.eth",
			ChainID:    "1",
		},
		strings.ToLower("0xC18360217D8F7Ab5e7c516566761Ea12Ce7F9D72"): {
			OriginalID: "ens.eth",
			ChainID:    "1",
		},
	}
)

type DelegationDetails struct {
	Address string
	Weight  int
}

type Erc20Mapping struct {
	OriginalID string
	ChainID    string
}

type Delegations struct {
	Details    []DelegationDetails
	Expiration int
}

// History storing delegate actions history
type History struct {
	Action          string
	AddressFrom     string
	OriginalSpaceID string
	ChainID         string
	BlockNumber     int
	BlockTimestamp  int
	Source          string
	Payload         json.RawMessage
	Delegations     Delegations `gorm:"-"`
	VotingPower     string      `gorm:"-"`
}

func (History) TableName() string {
	return "delegates_history"
}

func convertToInternal(payload events.DelegatePayload) History {
	pl, _ := json.Marshal(payload.Delegations)

	delegations := make([]DelegationDetails, 0, len(payload.Delegations.Details))

	totalWeight := 0
	for _, d := range payload.Delegations.Details {
		totalWeight += d.Weight
	}

	normalizeMultiplier := 1
	if totalWeight <= 100 {
		normalizeMultiplier = 100
	}

	for _, d := range payload.Delegations.Details {
		delegations = append(delegations, DelegationDetails{
			Address: d.Address,
			Weight:  d.Weight * normalizeMultiplier,
		})
	}

	return History{
		Action:          payload.Action,
		AddressFrom:     payload.AddressFrom,
		ChainID:         payload.ChainID,
		OriginalSpaceID: prepareForUTF(payload.OriginalSpaceID),
		BlockNumber:     payload.BlockNumber,
		BlockTimestamp:  payload.BlockTimestamp,
		Delegations: Delegations{
			Expiration: payload.Delegations.Expiration,
			Details:    delegations,
		},
		Source:      sourceSplitDelegation,
		Payload:     pl,
		VotingPower: "0",
	}
}

func convertERC20ToInternal(payload events.ERC20DelegatePayload) History {
	pl, _ := json.Marshal(payload)
	daoInfo, ok := daoErc20Set[strings.ToLower(payload.Token)]
	if !ok {
		log.Warn().Msgf("dao erc20 mapping not found for token %s", payload.Token)
	}

	return History{
		Action:          events.DelegateActionSet,
		AddressFrom:     payload.AddressFrom,
		ChainID:         daoInfo.ChainID,
		OriginalSpaceID: daoInfo.OriginalID,
		BlockNumber:     int(payload.BlockNumber),
		BlockTimestamp:  int(payload.BlockTimestamp),
		Delegations: Delegations{
			Details: []DelegationDetails{
				{
					Address: payload.AddressTo,
					Weight:  10000,
				},
			},
		},
		Source:      sourceErc20Votes,
		VotingPower: payload.VotingPower,
		Payload:     pl,
	}
}

func prepareForUTF(s string) string {
	return strings.ReplaceAll(s, "\u0000", "")
}

type Summary struct {
	AddressFrom        string
	AddressTo          string
	DaoID              string
	Weight             int
	LastBlockTimestamp int
	ExpiresAt          int64
	CreatedAt          time.Time
	ChainID            *string
	Type               string
	VotingPower        string

	// virtual property
	MaxCnt     int    `gorm:"-"`
	ProposalID string `gorm:"-"`
}

func (Summary) TableName() string {
	return "delegates_summary"
}

func (s *Summary) Expired() bool {
	if s.ExpiresAt == 0 {
		return false
	}

	return time.Now().Unix() > s.ExpiresAt
}

func (s *Summary) SelfDelegation() bool {
	return s.AddressTo == s.AddressFrom
}

type Proposal struct {
	ID            string
	OriginalDaoID string
	Author        string
}

func convertEventToProposal(event events.ProposalPayload) Proposal {
	return Proposal{
		ID:            event.ID,
		OriginalDaoID: event.DaoID,
		Author:        event.Author,
	}
}

type Vote struct {
	Voter         string
	OriginalDaoID string
	ProposalID    string
}

func convertEventToVoteDetails(event events.VotesPayload) []Vote {
	votes := make([]Vote, 0, len(event))
	for _, info := range event {
		votes = append(votes, Vote{
			Voter:         info.Voter,
			OriginalDaoID: info.OriginalDaoID,
			ProposalID:    info.ProposalID,
		})
	}

	return votes
}

type summaryByVote struct {
	Summary

	ProposalID string
}

type AllowedDao struct {
	DaoName    string
	CreatedAt  time.Time
	InternalID uuid.UUID
}

func (a *AllowedDao) TableName() string {
	return "delegate_allowed_daos"
}
