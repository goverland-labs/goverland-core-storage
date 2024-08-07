package delegate

import (
	"github.com/google/uuid"
)

type GetDelegatesRequest struct {
	DaoID         uuid.UUID
	QueryAccounts []string
	Sort          *string
	Limit         int
	Offset        int
}

type GetDelegatesResponse struct {
	Delegates []Delegate
}

type Delegate struct {
	Address                  string
	ENSName                  string
	DelegatorCount           int32
	PercentOfDelegators      float64
	VotingPower              float64
	PercentOfVotingPower     float64
	About                    string
	Statement                string
	UserDelegatedVotingPower float64
	VotesCount               int32
	ProposalsCount           int32
	CreateProposalsCount     int32
}
