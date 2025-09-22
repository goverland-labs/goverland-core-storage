package delegate

import (
	"time"

	"github.com/google/uuid"
)

const (
	DelegationTypeSplitDelegation DelegationType = "split-delegation"
	DelegationTypeDelegation      DelegationType = "delegation"
	DelegationTypeERC20Votes      DelegationType = "erc20-votes"
	DelegationTypeUnrecognized    DelegationType = "unrecognised"
)

type DelegationType string

type GetDelegatesRequest struct {
	DaoID          uuid.UUID
	QueryAccounts  []string
	Sort           *string
	Limit          int
	Offset         int
	DelegationType DelegationType
	ChainID        *string
}

type GetDelegatesResponse struct {
	Delegates []Delegate
	Total     int32
}

type Delegate struct {
	Address               string
	ENSName               string
	DelegatorCount        int32
	PercentOfDelegators   float64
	VotingPower           float64
	PercentOfVotingPower  float64
	About                 string
	Statement             string
	VotesCount            int32
	CreatedProposalsCount int32
}

type GetDelegateProfileRequest struct {
	DaoID   uuid.UUID
	Address string
}

type GetDelegateProfileResponse struct {
	Address              string
	VotingPower          float64
	IncomingPower        float64
	OutgoingPower        float64
	PercentOfVotingPower float64
	PercentOfDelegators  float64
	Delegates            []ProfileDelegateItem
	Expiration           *time.Time
}

type ProfileDelegateItem struct {
	Address        string
	ENSName        string
	Weight         float64
	DelegatedPower float64
}
