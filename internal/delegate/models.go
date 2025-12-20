package delegate

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	events "github.com/goverland-labs/goverland-platform-events/events/aggregator"
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
		// [1.8M holders | 108M transactions],
		strings.ToLower("0x912ce59144191c1204e64559fe8253a0e49e6548"): {
			OriginalID: "arbitrumfoundation.eth",
			ChainID:    "42161",
		},
		// [25K holders | 1M transactions] ,
		strings.ToLower("0x7189fb5B6504bbfF6a852B13B7B82a3c118fDc27"): {
			OriginalID: "etherfi-dao.eth",
			ChainID:    "42161",
		},
		// [0,1K holders | 1K transactions],
		strings.ToLower("0x54b6e28a869a56f4e34d1187ae0a35b7dd3be111"): {
			OriginalID: "maiadao.eth",
			ChainID:    "42161",
		},
		// [28K holders | 272K transactions],
		strings.ToLower("0xCa14007Eff0dB1f8135f4C25B34De49AB0d42766"): {
			OriginalID: "integration-test.eth", // "starknet.eth" // todo: rollback after testing
			ChainID:    "1",
		},
		// [ 65K holders | 1.3M transactions] ,
		strings.ToLower("0xC18360217D8F7Ab5e7c516566761Ea12Ce7F9D72"): {
			OriginalID: "ens.eth",
			ChainID:    "1",
		},
		// [25K holders | 370K transactions],
		strings.ToLower("0x3c3a81e81dc49A522A592e7622A7E711c06bf354"): {
			OriginalID: "bitdao.eth",
			ChainID:    "1",
		},
		// [4K holders | 159K transactions],
		strings.ToLower("0xd9Fcd98c322942075A5C3860693e9f4f03AAE07b"): {
			OriginalID: "eulerdao.eth",
			ChainID:    "1",
		},
		// [12K holders | 122K transactions],
		strings.ToLower("0xc5102fE9359FD9a28f877a67E36B0F050d81a3CC"): {
			OriginalID: "hop.eth",
			ChainID:    "1",
		},
		// [110K holders | 1M transactions],
		strings.ToLower("0xfe0c30065b384f05761f15d0cc899d4f9f9cc0eb"): {
			OriginalID: "etherfi-dao.eth",
			ChainID:    "1",
		},
		// [16K holders | 50K transfers],
		strings.ToLower("0xF5E3D1290FDBFC50ec436f021ad516D0Bcac5d28"): {
			OriginalID: "integration-test.eth",
			ChainID:    "56",
		},
		// wormhole connected to parason due to comparing with previous version
		strings.ToLower("0xB0fFa8000886e57F86dd5264b9582b2Ad87b2b91"): {
			OriginalID: "parason.eth",
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
		Source:  sourceSplitDelegation,
		Payload: pl,
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
	LogIndex           int
	ExpiresAt          int64
	CreatedAt          time.Time
	ChainID            *string
	Type               string

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

// Erc20EventHistory storing erc20 votes history
type Erc20EventHistory struct {
	ID            string
	OriginalDaoID string `gorm:"original_dao_id"`
	ChainID       string
	BlockNumber   int
	LogIndex      int
	Type          string
	Payload       json.RawMessage
	CreatedAt     time.Time
}

func (Erc20EventHistory) TableName() string {
	return "erc20_event_history"
}

type ERC20Delegate struct {
	ID             uint64
	Address        string
	DaoID          uuid.UUID
	ChainID        string
	VP             string
	BlockNumber    int
	LogIndex       int
	RepresentedCnt int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (ERC20Delegate) TableName() string {
	return "erc20_delegates"
}

type ERC20Balance struct {
	ID        uint
	Address   string
	DaoID     uuid.UUID
	ChainID   string
	Value     string `gorm:"type:numeric(78,0);not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ERC20Balance) TableName() string {
	return "erc20_balances"
}

type ERC20Totals struct {
	ID              uint
	DaoID           uuid.UUID
	ChainID         string
	VotingPower     string `gorm:"type:numeric(78,0);not null;default:0"`
	TotalDelegators int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (ERC20Totals) TableName() string {
	return "erc20_totals"
}
