package delegate

import (
	"context"

	"github.com/goverland-labs/goverland-datasource-snapshot/protocol/delegatepb"
)

type Service struct {
	delegateClient delegatepb.DelegateClient
}

func NewService(dc delegatepb.DelegateClient) *Service {
	return &Service{
		delegateClient: dc,
	}
}

func (s *Service) GetDelegates(ctx context.Context, request GetDelegatesRequest) (*GetDelegatesResponse, error) {
	resp, err := s.delegateClient.GetDelegates(ctx, &delegatepb.GetDelegatesRequest{
		Addresses: request.Addresses,
		Sort:      request.Sort,
		Limit:     int32(request.Limit),
		Offset:    int32(request.Offset),
	})
	if err != nil {
		return nil, err
	}

	delegates := make([]Delegate, 0, len(resp.Delegates))
	for _, d := range resp.GetDelegates() {
		delegates = append(delegates, Delegate{
			Address:                  d.GetAddress(),
			DelegatorCount:           d.GetDelegatorCount(),
			PercentOfDelegators:      d.GetPercentOfDelegators(),
			VotingPower:              d.GetVotingPower(),
			PercentOfVotingPower:     d.GetPercentOfVotingPower(),
			About:                    d.GetAbout(),
			Statement:                d.GetStatement(),
			UserDelegatedVotingPower: d.GetUserDelegatedVotingPower(),
			VotesCount:               0,
			ProposalsCount:           0,
			CreateProposalsCount:     0,
		})
	}

	return &GetDelegatesResponse{
		Delegates: delegates,
	}, nil
}
