package delegate

import (
	"context"
	"maps"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/goverland-labs/goverland-core-storage/protocol/storagepb"

	"github.com/goverland-labs/goverland-core-storage/internal/dao"
)

type DaoSearcher interface {
	GetByFilters(filters []dao.Filter) (dao.DaoList, error)
}

type Server struct {
	storagepb.UnimplementedDelegateServer

	sp *Service
	ds DaoSearcher
}

func NewServer(sp *Service, ds DaoSearcher) *Server {
	return &Server{
		sp: sp,
		ds: ds,
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
		Total:     delegatesResponse.Total,
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

	var expiration *timestamppb.Timestamp
	if profile.Expiration != nil {
		expiration = timestamppb.New(*profile.Expiration)
	}

	return &storagepb.GetDelegateProfileResponse{
		Address:              profile.Address,
		VotingPower:          profile.VotingPower,
		IncomingPower:        profile.IncomingPower,
		OutgoingPower:        profile.OutgoingPower,
		PercentOfVotingPower: profile.PercentOfVotingPower,
		PercentOfDelegators:  profile.PercentOfDelegators,
		Delegates:            delegates,
		Expiration:           expiration,
	}, nil
}

func (s *Server) GetAllDelegations(ctx context.Context, req *storagepb.GetAllDelegationsRequest) (*storagepb.GetAllDelegationsResponse, error) {
	if req.GetAddress() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}

	// delegations [dao_id: [summary, ...]]
	delegations, err := s.sp.getAllDelegations(ctx, req.GetAddress())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get delegations")
	}

	if len(delegations) == 0 {
		return &storagepb.GetAllDelegationsResponse{}, nil
	}

	daoIDs := slices.Collect(maps.Keys(delegations))
	daoList, err := s.ds.GetByFilters([]dao.Filter{
		dao.DaoIDsFilter{DaoIDs: daoIDs},
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get dao info")
	}

	addresses := make([]string, 0, len(delegations))
	for _, d := range delegations {
		for _, info := range d {
			addresses = append(addresses, info.AddressTo)
		}
	}
	ensNames, err := s.sp.resolveAddressesName(addresses)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to resolve ens names")
	}

	response := &storagepb.GetAllDelegationsResponse{
		Delegations: make([]*storagepb.DelegationSummary, 0, len(delegations)),
	}

	delegationsCnt := 0
	for _, di := range daoList.Daos {
		list, ok := delegations[di.ID.String()]
		if !ok {
			log.Warn().Msgf("dao info not found: %s", di.ID.String())
			continue
		}

		delegationsCnt += len(list)
		dl := make([]*storagepb.DelegationDetails, 0, len(list))
		for _, d := range list {
			var expires *timestamppb.Timestamp
			if d.ExpiresAt != 0 {
				expires = timestamppb.New(time.Unix(d.ExpiresAt, 0))
			}

			dl = append(dl, &storagepb.DelegationDetails{
				Address:             d.AddressTo,
				EnsName:             ensNames[d.AddressTo],
				PercentOfDelegators: int32(d.Weight),
				Expiration:          expires,
			})
		}

		response.Delegations = append(response.Delegations, &storagepb.DelegationSummary{
			Dao:         dao.ConvertDaoToAPI(&di),
			Delegations: dl,
		})
	}

	response.TotalDelegationsCount = int32(delegationsCnt)

	return response, nil
}
