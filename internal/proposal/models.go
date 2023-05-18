package proposal

import (
	"time"

	aggevents "github.com/goverland-labs/platform-events/events/aggregator"
	events "github.com/goverland-labs/platform-events/events/core"
)

type Choices []string

type Scores []float32

type Strategy struct {
	Name    string
	Network string
}

type Strategies []Strategy

func convertToStrategies(list Strategies) []events.StrategyPayload {
	result := make([]events.StrategyPayload, len(list))
	for i, strategy := range list {
		result[i] = events.StrategyPayload{
			Name:    strategy.Name,
			Network: strategy.Network,
		}
	}

	return result
}

// Proposal model
// todo: check queries to the DB and add indexes
type Proposal struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Ipfs      string
	Author    string
	Created   int
	// todo: think about relation to the DAO model.
	// Some proposal events could be processed early then dao created
	DaoID         string
	Network       string
	Symbol        string
	Type          string
	Strategies    Strategies `gorm:"serializer:json"`
	Title         string
	Body          string
	Discussion    string
	Choices       Choices `gorm:"serializer:json"`
	Start         int
	End           int
	Quorum        float64
	Privacy       string
	Snapshot      string
	State         string
	Link          string
	App           string
	Scores        Scores `gorm:"serializer:json"`
	ScoresState   string
	ScoresTotal   float32
	ScoresUpdated int
	Votes         int
}

func convertToCoreEvent(p Proposal) events.ProposalPayload {
	return events.ProposalPayload{
		ID:            p.ID,
		Ipfs:          p.Ipfs,
		Author:        p.Author,
		Created:       p.Created,
		DaoID:         p.DaoID,
		Network:       p.Network,
		Symbol:        p.Symbol,
		Type:          p.Type,
		Strategies:    convertToStrategies(p.Strategies),
		Title:         p.Title,
		Body:          p.Body,
		Discussion:    p.Discussion,
		Choices:       p.Choices,
		Start:         p.Start,
		End:           p.End,
		Quorum:        p.Quorum,
		Privacy:       p.Privacy,
		Snapshot:      p.Snapshot,
		State:         p.State,
		Link:          p.Link,
		App:           p.App,
		Scores:        p.Scores,
		ScoresState:   p.ScoresState,
		ScoresTotal:   p.ScoresTotal,
		ScoresUpdated: p.ScoresUpdated,
		Votes:         p.Votes,
	}
}

func convertToProposal(p aggevents.ProposalPayload) Proposal {
	return Proposal{
		ID:            p.ID,
		Ipfs:          p.Ipfs,
		Author:        p.Author,
		Created:       p.Created,
		DaoID:         p.DaoID,
		Network:       p.Network,
		Symbol:        p.Symbol,
		Type:          p.Type,
		Strategies:    convertToInternalStrategies(p.Strategies),
		Title:         p.Title,
		Body:          p.Body,
		Discussion:    p.Discussion,
		Choices:       p.Choices,
		Start:         p.Start,
		End:           p.End,
		Quorum:        p.Quorum,
		Privacy:       p.Privacy,
		Snapshot:      p.Snapshot,
		State:         p.State,
		Link:          p.Link,
		App:           p.App,
		Scores:        p.Scores,
		ScoresState:   p.ScoresState,
		ScoresTotal:   p.ScoresTotal,
		ScoresUpdated: p.ScoresUpdated,
		Votes:         p.Votes,
	}
}

func convertToInternalStrategies(s []aggevents.StrategyPayload) Strategies {
	res := make(Strategies, len(s))
	for i, item := range s {
		res[i] = Strategy{
			Name:    item.Name,
			Network: item.Network,
		}
	}

	return res
}
