package proposal

import (
	"time"

	"github.com/google/uuid"
	aggevents "github.com/goverland-labs/platform-events/events/aggregator"
	coreevents "github.com/goverland-labs/platform-events/events/core"
	events "github.com/goverland-labs/platform-events/events/core"
)

type Choices []string

type Scores []float32

type Strategy struct {
	Name    string
	Network string
	Params  map[string]interface{}
}

type Strategies []Strategy

func convertToStrategies(list Strategies) []events.StrategyPayload {
	result := make([]events.StrategyPayload, len(list))
	for i, strategy := range list {
		result[i] = events.StrategyPayload{
			Name:    strategy.Name,
			Network: strategy.Network,
			Params:  strategy.Params,
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
	DaoOriginalID string `gorm:"-"`
	DaoID         uuid.UUID
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
	Timeline      Timeline `gorm:"serializer:json"`
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
		DaoOriginalID: p.DaoID,
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

func convertToTimeline(tl []coreevents.TimelineItem) Timeline {
	if len(tl) == 0 {
		return Timeline{}
	}

	res := make(Timeline, len(tl))
	for i := range tl {
		res[i] = TimelineItem{
			CreatedAt: tl[i].CreatedAt,
			Action:    TimelineAction(tl[i].Action),
		}
	}

	return res
}

func convertToInternalStrategies(s []aggevents.StrategyPayload) Strategies {
	res := make(Strategies, len(s))
	for i, item := range s {
		res[i] = Strategy{
			Name:    item.Name,
			Network: item.Network,
			Params:  item.Params,
		}
	}

	return res
}
