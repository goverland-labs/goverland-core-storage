package vote

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	aggevents "github.com/goverland-labs/goverland-platform-events/events/aggregator"
	events "github.com/goverland-labs/goverland-platform-events/events/core"
)

type Vote struct {
	ID            string `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Ipfs          string
	OriginalDaoID string `gorm:"-"`
	DaoID         uuid.UUID
	ProposalID    string
	Voter         string
	EnsName       string
	Created       int
	Reason        string
	Choice        json.RawMessage
	App           string
	Vp            float64
	VpByStrategy  []float64 `gorm:"serializer:json"`
	VpState       string
}

func convertToInternal(pl aggevents.VotesPayload) []Vote {
	res := make([]Vote, len(pl))
	for i, item := range pl {
		res[i] = Vote{
			ID:            item.ID,
			Ipfs:          item.Ipfs,
			OriginalDaoID: item.OriginalDaoID,
			ProposalID:    item.ProposalID,
			Voter:         item.Voter,
			Created:       item.Created,
			Reason:        item.Reason,
			Choice:        item.Choice,
			App:           item.App,
			Vp:            item.Vp,
			VpByStrategy:  item.VpByStrategy,
			VpState:       item.VpState,
		}
	}

	return res
}

func convertToCoreEvent(votes []Vote) events.VotesPayload {
	res := make([]events.VotePayload, len(votes))
	for i, item := range votes {
		res[i] = events.VotePayload{
			ID:           item.ID,
			Ipfs:         item.Ipfs,
			DaoID:        item.DaoID,
			ProposalID:   item.ProposalID,
			Voter:        item.Voter,
			Created:      item.Created,
			Reason:       item.Reason,
			Choice:       item.Choice,
			App:          item.App,
			Vp:           item.Vp,
			VpByStrategy: item.VpByStrategy,
			VpState:      item.VpState,
		}
	}

	return res
}

type ResolvedAddress struct {
	Address string
	Name    string
}
