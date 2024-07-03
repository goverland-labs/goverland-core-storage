package delegate

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/goverland-labs/goverland-datasource-snapshot/protocol/delegatepb"

	"github.com/goverland-labs/goverland-core-storage/internal/dao"
)

type DaoProvider interface {
	GetByID(id uuid.UUID) (*dao.Dao, error)
}

type Service struct {
	delegateClient delegatepb.DelegateClient
	daoProvider    DaoProvider
}

func NewService(dc delegatepb.DelegateClient, daoProvider DaoProvider) *Service {
	return &Service{
		delegateClient: dc,
		daoProvider:    daoProvider,
	}
}

func (s *Service) GetDelegates(ctx context.Context, request GetDelegatesRequest) (*GetDelegatesResponse, error) {
	daoEntity, err := s.daoProvider.GetByID(request.DaoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dao: %w", err)
	}

	resp, err := s.delegateClient.GetDelegates(ctx, &delegatepb.GetDelegatesRequest{
		DaoName:   daoEntity.Name,
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
