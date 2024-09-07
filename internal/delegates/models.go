package delegates

import (
	"encoding/json"
	"strings"
	"time"

	events "github.com/goverland-labs/goverland-platform-events/events/core"
)

var (
	actionClear  string = "clear"
	actionExpire string = "expire"
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
