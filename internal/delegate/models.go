package delegate

import (
	"encoding/json"
	"strings"
	"time"

	events "github.com/goverland-labs/goverland-platform-events/events/aggregator"
)

var (
	actionClear  = "clear"
	actionExpire = "expire"
)

type DelegationDetails struct {
	Address string
	Weight  int
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
	Delegations     Delegations `gorm:"-"`
	Payload         json.RawMessage
}

func (History) TableName() string {
	return "delegates_history"
}

func convertToInternal(payload events.DelegatePayload) History {
	pl, _ := json.Marshal(payload.Delegations)

	delegations := make([]DelegationDetails, 0, len(payload.Delegations.Details))
	for _, d := range payload.Delegations.Details {
		delegations = append(delegations, DelegationDetails{
			Address: d.Address,
			Weight:  d.Weight,
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
	ExpiresAt          int64
	CreatedAt          time.Time
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
