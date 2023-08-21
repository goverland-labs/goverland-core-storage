package proposal

import (
	"sort"
	"strings"
	"time"
)

const (
	None                        TimelineAction = ""
	ProposalCreated             TimelineAction = "proposal.created"
	ProposalUpdated             TimelineAction = "proposal.updated"
	ProposalVotingStartsSoon    TimelineAction = "proposal.voting.starts_soon"
	ProposalVotingEndsSoon      TimelineAction = "proposal.voting.ends_soon"
	ProposalVotingStarted       TimelineAction = "proposal.voting.started"
	ProposalVotingQuorumReached TimelineAction = "proposal.voting.quorum_reached"
	ProposalVotingEnded         TimelineAction = "proposal.voting.ended"
)

type Timeline []TimelineItem

func (t *Timeline) AddUniqueAction(createdAt time.Time, action TimelineAction) (isNew bool) {
	if *t == nil {
		*t = make(Timeline, 0, 1)
	}

	for i := range *t {
		if (*t)[i].Action.Equals(action) {
			return false
		}
	}

	*t = append(*t, TimelineItem{
		CreatedAt: createdAt,
		Action:    action,
	})

	return true
}

func (t *Timeline) AddNonUniqueAction(createdAt time.Time, action TimelineAction) {
	if *t == nil {
		*t = make(Timeline, 0, 1)
	}

	*t = append(*t, TimelineItem{
		CreatedAt: createdAt,
		Action:    action,
	})
}

func (t *Timeline) ContainsAction(action TimelineAction) bool {
	if t == nil || len(*t) == 0 {
		return false
	}

	for _, item := range *t {
		if item.Action.Equals(action) {
			return true
		}
	}

	return false
}

func (t *Timeline) Sort() {
	if t == nil || len(*t) == 0 {
		return
	}

	sort.SliceStable(*t, func(i, j int) bool {
		if (*t)[i].CreatedAt.Equal((*t)[j].CreatedAt) {
			return actionWeight((*t)[i].Action) < actionWeight((*t)[j].Action)
		}

		return (*t)[i].CreatedAt.Before((*t)[j].CreatedAt)
	})
}

func actionWeight(a TimelineAction) int {
	switch a {
	case ProposalCreated:
		return 1
	case ProposalVotingQuorumReached:
		return 3
	case ProposalVotingEnded:
		return 4
	default:
		return 2
	}
}

func (t *Timeline) LastAction() TimelineAction {
	if t == nil {
		return None
	}

	if len(*t) == 0 {
		return None
	}

	return (*t)[len(*t)-1].Action
}

type TimelineItem struct {
	CreatedAt time.Time      `json:"created_at"`
	Action    TimelineAction `json:"action"`
}

type TimelineAction string

func (a TimelineAction) Equals(action TimelineAction) bool {
	return strings.EqualFold(string(a), string(action))
}

func (t *Timeline) ActualizeTimeline() Timeline {
	actual := make(Timeline, 0, len(*t))

	var minCreatedAt, maxCreatedAt, createdAt, finishedAt, quorumReachedAt time.Time
	var createdIdx, quorumReachedIdx int
	for idx, info := range *t {
		action := getAction(info.Action)

		if action == ProposalCreated {
			createdAt = info.CreatedAt
			createdIdx = idx
		}

		if action == ProposalVotingQuorumReached {
			quorumReachedAt = info.CreatedAt
			quorumReachedIdx = idx
		}

		if action == ProposalVotingEnded {
			finishedAt = info.CreatedAt
		}

		if minCreatedAt.IsZero() || minCreatedAt.After(info.CreatedAt) {
			minCreatedAt = info.CreatedAt
		}

		if maxCreatedAt.Before(info.CreatedAt) {
			maxCreatedAt = info.CreatedAt
		}

		actual = append(actual, TimelineItem{
			CreatedAt: info.CreatedAt,
			Action:    action,
		})
	}

	if !createdAt.IsZero() && !createdAt.Equal(minCreatedAt) {
		actual[createdIdx].CreatedAt = minCreatedAt
	}

	if !quorumReachedAt.IsZero() && !finishedAt.IsZero() && quorumReachedAt.After(finishedAt) {
		actual[quorumReachedIdx].CreatedAt = finishedAt
	}

	actual.Sort()

	return actual
}

func getAction(action TimelineAction) TimelineAction {
	switch action {

	default:
		return action
	}
}
