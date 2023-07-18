package vote

import (
	"time"

	"github.com/google/uuid"
	pevents "github.com/goverland-labs/platform-events/events/aggregator"
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
	Created       int
	Reason        string
	Choice        int
	App           string
	Vp            float64
	VpByStrategy  []float64 `gorm:"serializer:json"`
	VpState       string
}

func convertToInternal(pl pevents.VotesPayload) []Vote {
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
