package proposal

import (
	"time"

	"github.com/google/uuid"
	aggevents "github.com/goverland-labs/platform-events/events/aggregator"
	coreevents "github.com/goverland-labs/platform-events/events/core"
	events "github.com/goverland-labs/platform-events/events/core"
)

const (
	StatePending   = "pending"
	StateActive    = "active"
	StateCancelled = "canceled"
	StateFailed    = "failed"
	StateSucceeded = "succeeded"
	StateDefeated  = "defeated"
)

type State string

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
	ID            string `gorm:"primary_key"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Ipfs          string
	Author        string
	Created       int
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
	State         State
	OriginalState string
	Link          string
	App           string
	Scores        Scores `gorm:"serializer:json"`
	ScoresState   string
	ScoresTotal   float32
	ScoresUpdated int
	Votes         int
	Timeline      Timeline `gorm:"serializer:json"`
	EnsName       string
	Spam          bool
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
		State:         string(p.State), // todo: replace state in core as enum
		Link:          p.Link,
		App:           p.App,
		Scores:        p.Scores,
		ScoresState:   p.ScoresState,
		ScoresTotal:   p.ScoresTotal,
		ScoresUpdated: p.ScoresUpdated,
		Votes:         p.Votes,
		EnsName:       p.EnsName,
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
		OriginalState: p.State,
		Link:          p.Link,
		App:           p.App,
		Scores:        p.Scores,
		ScoresState:   p.ScoresState,
		ScoresTotal:   p.ScoresTotal,
		ScoresUpdated: p.ScoresUpdated,
		Votes:         p.Votes,
		Spam:          p.Flagged,
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

type ResolvedAddress struct {
	Address string
	Name    string
}

func convertToResolvedAddresses(list []coreevents.EnsNamePayload) []ResolvedAddress {
	res := make([]ResolvedAddress, 0, len(list))
	for i := range list {
		res = append(res, ResolvedAddress{
			Address: list[i].Address,
			Name:    list[i].Name,
		})
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

func (p *Proposal) CalculateState() State {
	if p.Deleted() {
		return StateCancelled
	}

	if p.Pending() {
		return StatePending
	}

	if p.InProgress() {
		return StateActive
	}

	if p.Votes == 0 ||
		(p.QuorumSpecified() && !p.QuorumReached()) {
		return StateFailed
	}

	if p.IsVotedAgainst() {
		return StateDefeated
	}

	return StateSucceeded
}

func (p *Proposal) InProgress() bool {
	startsAt := time.Unix(int64(p.Start), 0)
	endsAt := time.Unix(int64(p.End), 0)

	return time.Now().After(startsAt) && time.Now().Before(endsAt)
}

func (p *Proposal) Pending() bool {
	startsAt := time.Unix(int64(p.Start), 0)

	return startsAt.After(time.Now())
}

func (p *Proposal) Deleted() bool {
	return p.State == StateCancelled
}

func (p *Proposal) QuorumSpecified() bool {
	return p.Quorum > 0
}

func (p *Proposal) QuorumReached() bool {
	if p.Quorum == 0 {
		return false
	}

	return float64(p.ScoresTotal) >= p.Quorum
}

func (p *Proposal) IsBasic() bool {
	return p.Type == "basic"
}

func (p *Proposal) IsVotedAgainst() bool {
	if !p.IsBasic() {
		return false
	}

	// invalid data, should collect votes for
	// 0 - For
	// 1 - Against
	// 2 - Abstain
	if len(p.Scores) != 3 {
		return false
	}

	forVotes := p.Scores[0]
	againstVotes := p.Scores[1]

	return againstVotes > forVotes
}
