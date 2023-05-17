package vote

import (
	"time"

	pevents "github.com/goverland-labs/platform-events/events/aggregator"
)

type Vote struct {
	ID         string `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Ipfs       string
	ProposalID string
	Voter      string
	Created    int
	Reason     string
}

func convertToInternal(pl pevents.VotesPayload) []Vote {
	res := make([]Vote, len(pl))
	for i, item := range pl {
		res[i] = Vote{
			ID:         item.ID,
			Ipfs:       item.Ipfs,
			ProposalID: item.ProposalID,
			Voter:      item.Voter,
			Created:    item.Created,
			Reason:     item.Reason,
		}
	}

	return res
}
