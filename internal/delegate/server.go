package delegate

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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
	if req.GetDaoId() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid dao ID")
	}

	daoID, err := uuid.Parse(req.GetDaoId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid dao ID format")
	}

	delegatesResponse, err := s.sp.GetDelegates(ctx, GetDelegatesRequest{
		DaoID:         daoID,
		QueryAccounts: req.GetQueryAccounts(),
		Sort:          req.Sort,
		Limit:         int(req.GetLimit()),
		Offset:        int(req.GetOffset()),
	})
	if err != nil {
		log.Error().
			Err(err).
			Str("dao_id", daoID.String()).
			Msg("failed to get delegates")

		return nil, status.Errorf(codes.Internal, "failed to get delegates: %v", err)
	}

	delegatesResult := make([]*storagepb.DelegateEntry, 0, len(delegatesResponse.Delegates))
	for _, d := range delegatesResponse.Delegates {
		delegatesResult = append(delegatesResult, &storagepb.DelegateEntry{
			Address:                  d.Address,
			EnsName:                  d.ENSName,
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
