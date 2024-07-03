package delegate

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goverland-labs/goverland-core-storage/protocol/storagepb"
)

type Server struct {
	storagepb.UnimplementedDelegateServer

	sp *Service
}

func NewServer(sp *Service) *Server {
	return &Server{
		sp: sp,
	}
}

func (s *Server) GetDelegates(ctx context.Context, req *storagepb.GetDelegatesRequest) (*storagepb.GetDelegatesResponse, error) {
	delegatesResponse, err := s.sp.GetDelegates(ctx, GetDelegatesRequest{
		Addresses: req.GetAddresses(),
		Sort:      req.GetSort(),
		Limit:     int(req.GetLimit()),
		Offset:    int(req.GetOffset()),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get delegates: %v", err)
	}

	delegatesResult := make([]*storagepb.DelegateEntry, 0, len(delegatesResponse.Delegates))
	for _, d := range delegatesResponse.Delegates {
		delegatesResult = append(delegatesResult, &storagepb.DelegateEntry{
			Address:                  d.Address,
			DelegatorCount:           d.DelegatorCount,
			PercentOfDelegators:      d.PercentOfDelegators,
			VotingPower:              d.VotingPower,
			PercentOfVotingPower:     d.PercentOfVotingPower,
			About:                    d.About,
			Statement:                d.Statement,
			UserDelegatedVotingPower: d.UserDelegatedVotingPower,
			VotesCount:               d.VotesCount,
			ProposalsCount:           d.ProposalsCount,
			CreateProposalsCount:     d.CreateProposalsCount,
		})
	}

	return &storagepb.GetDelegatesResponse{
		Delegates: delegatesResult,
	}, nil
}
