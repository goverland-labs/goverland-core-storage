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
			Address:               d.Address,
			EnsName:               d.ENSName,
			DelegatorCount:        d.DelegatorCount,
			PercentOfDelegators:   d.PercentOfDelegators,
			VotingPower:           d.VotingPower,
			PercentOfVotingPower:  d.PercentOfVotingPower,
			About:                 d.About,
			Statement:             d.Statement,
			VotesCount:            d.VotesCount,
			CreatedProposalsCount: d.CreatedProposalsCount,
		})
	}

	return &storagepb.GetDelegatesResponse{
		Delegates: delegatesResult,
	}, nil
}

func (s *Server) GetDelegateProfile(ctx context.Context, req *storagepb.GetDelegateProfileRequest) (*storagepb.GetDelegateProfileResponse, error) {
	if req.GetDaoId() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid dao ID")
	}

	daoID, err := uuid.Parse(req.GetDaoId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid dao ID format")
	}

	profile, err := s.sp.GetDelegateProfile(ctx, GetDelegateProfileRequest{
		DaoID:   daoID,
		Address: req.GetAddress(),
	})
	if err != nil {
		log.Error().
			Err(err).
			Str("dao_id", daoID.String()).
			Str("address", req.GetAddress()).
			Msg("failed to get delegate profile")

		return nil, status.Errorf(codes.Internal, "failed to get delegate profile: %v", err)
	}

	delegates := make([]*storagepb.ProfileDelegateItem, 0, len(profile.Delegates))
	for _, d := range profile.Delegates {
		delegates = append(delegates, &storagepb.ProfileDelegateItem{
			Address:        d.Address,
			EnsName:        d.ENSName,
			Weight:         d.Weight,
			DelegatedPower: d.DelegatedPower,
		})
	}

	return &storagepb.GetDelegateProfileResponse{
		Address:              profile.Address,
		VotingPower:          profile.VotingPower,
		IncomingPower:        profile.IncomingPower,
		OutgoingPower:        profile.OutgoingPower,
		PercentOfVotingPower: profile.PercentOfVotingPower,
		PercentOfDelegators:  profile.PercentOfDelegators,
	}, nil
}
