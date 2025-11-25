package delegate

import (
	"encoding/json"
	"fmt"
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

type GetDelegatesMixedRequest struct {
	DaoID          uuid.UUID
	QueryAccounts  []string
	Sort           *string
	Limit          int
	Offset         int
	DelegationType DelegationType
	ChainID        *string
}

type DelegatesWrapper struct {
	DaoID          uuid.UUID
	Delegates      []Delegate
	DelegationType DelegationType
	ChainID        *string
	Total          int32
}

type GetDelegatesMixedResponse struct {
	List  []DelegatesWrapper
	Total int32
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
	DaoID          uuid.UUID
	Address        string
	DelegationType DelegationType
	ChainID        string
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

type ERC20Event interface {
	GetKey() string
	ConvertToHistory() *Erc20EventHistory
}

type ERC20Delegation struct {
	DelegatorAddress string
	AddressFrom      string
	AddressTo        string
	OriginalSpaceID  string
	ChainID          string
	BlockNumber      int
	BlockTimestamp   int
	LogIndex         int
}

func (e ERC20Delegation) GetKey() string {
	return generateUniqueKey(e.ChainID, e.BlockNumber, e.LogIndex)
}

func (e ERC20Delegation) ConvertToHistory() *Erc20EventHistory {
	payload, _ := json.Marshal(map[string]any{
		"delegator_address": e.DelegatorAddress,
		"address_from":      e.AddressFrom,
		"address_to":        e.AddressTo,
	})

	return &Erc20EventHistory{
		ID:            e.GetKey(),
		OriginalDaoID: e.OriginalSpaceID,
		ChainID:       e.ChainID,
		BlockNumber:   e.BlockNumber,
		LogIndex:      e.LogIndex,
		Type:          "delegation",
		Payload:       payload,
		CreatedAt:     time.Now(),
	}
}

func generateUniqueKey(chainID string, blockNumber, logIndex int) string {
	return fmt.Sprintf("%s_%d_%d", chainID, blockNumber, logIndex)
}

type ERC20VPChanges struct {
	Address         string
	OriginalSpaceID string
	ChainID         string
	BlockNumber     int
	BlockTimestamp  int
	LogIndex        int
	VP              string
	Delta           string
}

func (e ERC20VPChanges) GetKey() string {
	return generateUniqueKey(e.ChainID, e.BlockNumber, e.LogIndex)
}

func (e ERC20VPChanges) ConvertToHistory() *Erc20EventHistory {
	payload, _ := json.Marshal(map[string]any{
		"address": e.Address,
		"vp":      e.VP,
		"delta":   e.Delta,
	})

	return &Erc20EventHistory{
		ID:            e.GetKey(),
		OriginalDaoID: e.OriginalSpaceID,
		ChainID:       e.ChainID,
		BlockNumber:   e.BlockNumber,
		LogIndex:      e.LogIndex,
		Type:          "vp_changes",
		Payload:       payload,
		CreatedAt:     time.Now(),
	}
}

type ERC20Transfer struct {
	AddressFrom     string
	AddressTo       string
	OriginalSpaceID string
	ChainID         string
	BlockNumber     int
	BlockTimestamp  int
	LogIndex        int
	Amount          string
}

func (e ERC20Transfer) GetKey() string {
	return generateUniqueKey(e.ChainID, e.BlockNumber, e.LogIndex)
}

func (e ERC20Transfer) ConvertToHistory() *Erc20EventHistory {
	payload, _ := json.Marshal(map[string]any{
		"address_from": e.AddressFrom,
		"address_to":   e.AddressTo,
		"amount":       e.Amount,
	})

	return &Erc20EventHistory{
		ID:            e.GetKey(),
		OriginalDaoID: e.OriginalSpaceID,
		ChainID:       e.ChainID,
		BlockNumber:   e.BlockNumber,
		LogIndex:      e.LogIndex,
		Type:          "transfer",
		Payload:       payload,
		CreatedAt:     time.Now(),
	}
}

type VPUpdate struct {
	Value       string
	BlockNumber int
	LogIndex    int
}

type ERC20DelegateUpdate struct {
	Address    string
	OriginalID string
	ChainID    string
	VPUpdate   *VPUpdate
	CntDelta   *int
}

type ERC20TotalChanges struct {
	OriginalID      string
	ChainID         string
	VPDelta         string
	DelegatorsDelta int64
}

type ERC20DelegatorsRequest struct {
	Address string
	ChainID string
	DaoID   uuid.UUID
	Limit   int
	Offset  int
}

type AddressValue struct {
	Address    string
	TokenValue string
}
