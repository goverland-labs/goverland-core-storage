package delegate

import (
	"github.com/google/uuid"
)

type GetDelegatesRequest struct {
	DaoID     uuid.UUID
	Addresses []string
	Sort      string
	Limit     int
	Offset    int
}

type GetDelegatesResponse struct {
	Delegates []Delegate
}

type Delegate struct {
	Address                  string
	DelegatorCount           int32
	PercentOfDelegators      int32
	VotingPower              float64
	PercentOfVotingPower     int32
	About                    string
	Statement                string
	UserDelegatedVotingPower float64
	VotesCount               int32
	ProposalsCount           int32
	CreateProposalsCount     int32
}
