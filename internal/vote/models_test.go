package vote

import (
	"github.com/google/uuid"
	"testing"

	aggevents "github.com/goverland-labs/platform-events/events/aggregator"
	events "github.com/goverland-labs/platform-events/events/core"
	"github.com/stretchr/testify/assert"
)

var (
	id1 = uuid.New()
	id2 = uuid.New()
)

func TestUnitConvertToInternal(t *testing.T) {
	for name, tc := range map[string]struct {
		in  aggevents.VotesPayload
		out []Vote
	}{
		"empty input": {
			in:  nil,
			out: []Vote{},
		},
		"correct converting": {
			in: aggevents.VotesPayload{
				{
					ID:         "id-1",
					Ipfs:       "id-1",
					ProposalID: "proposal-1",
					Voter:      "voter-1",
					Created:    123,
					Reason:     "reason-1",
				},
				{
					ID:         "id-2",
					Ipfs:       "id-2",
					ProposalID: "proposal-2",
					Voter:      "voter-2",
					Created:    1234,
					Reason:     "reason-2",
				},
			},
			out: []Vote{
				{
					ID:         "id-1",
					Ipfs:       "id-1",
					ProposalID: "proposal-1",
					Voter:      "voter-1",
					Created:    123,
					Reason:     "reason-1",
				},
				{
					ID:         "id-2",
					Ipfs:       "id-2",
					ProposalID: "proposal-2",
					Voter:      "voter-2",
					Created:    1234,
					Reason:     "reason-2",
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			expected := convertToInternal(tc.in)
			assert.Equal(t, expected, tc.out)
		})
	}
}

func TestUnitConvertToCoreEvent(t *testing.T) {
	for name, tc := range map[string]struct {
		in  []Vote
		out events.VotesPayload
	}{
		"empty input": {
			in:  nil,
			out: []events.VotePayload{},
		},
		"correct converting": {
			in: []Vote{
				{
					ID:            "id-1",
					Ipfs:          "id-1",
					ProposalID:    "proposal-1",
					OriginalDaoID: "original",
					DaoID:         id1,
					Voter:         "voter-1",
					Created:       123,
					Reason:        "reason-1",
				},
				{
					ID:            "id-2",
					Ipfs:          "id-2",
					ProposalID:    "proposal-2",
					OriginalDaoID: "original2",
					DaoID:         id2,
					Voter:         "voter-2",
					Created:       1234,
					Reason:        "reason-2",
				},
			},
			out: events.VotesPayload{
				{
					ID:         "id-1",
					Ipfs:       "id-1",
					DaoID:      id1,
					ProposalID: "proposal-1",
					Voter:      "voter-1",
					Created:    123,
					Reason:     "reason-1",
				},
				{
					ID:         "id-2",
					Ipfs:       "id-2",
					DaoID:      id2,
					ProposalID: "proposal-2",
					Voter:      "voter-2",
					Created:    1234,
					Reason:     "reason-2",
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			expected := convertToCoreEvent(tc.in)
			assert.Equal(t, expected, tc.out)
		})
	}
}
