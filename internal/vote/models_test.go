package vote

import (
	"testing"

	pevents "github.com/goverland-labs/platform-events/events/aggregator"
	"github.com/stretchr/testify/assert"
)

func TestUnitConvertToInternal(t *testing.T) {
	for name, tc := range map[string]struct {
		in  pevents.VotesPayload
		out []Vote
	}{
		"empty input": {
			in:  nil,
			out: []Vote{},
		},
		"correct converting": {
			in: pevents.VotesPayload{
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
