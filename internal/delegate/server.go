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

const (
	defaultLimit = 10
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

	limit := int(req.GetLimit())
	if limit == 0 {
		limit = defaultLimit
	}

	delegatesResponse, err := s.sp.GetDelegates(ctx, GetDelegatesRequest{
		DaoID:          daoID,
		QueryAccounts:  req.GetQueryAccounts(),
		Sort:           req.Sort,
		Limit:          limit,
		Offset:         int(req.GetOffset()),
		ChainID:        req.ChainId,
		DelegationType: convertDelegationType(req.GetDelegationType()),
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
			DelegationType:        req.GetDelegationType(),
			ChainId:               req.ChainId,
		})
	}

	return &storagepb.GetDelegatesResponse{
		Delegates: delegatesResult,
		Total:     delegatesResponse.Total,
	}, nil
}

func convertDelegationType(dt storagepb.DelegationType) DelegationType {
	switch dt {
	case storagepb.DelegationType_DELEGATION_TYPE_SPLIT_DELEGATION:
		return DelegationTypeSplitDelegation
	case storagepb.DelegationType_DELEGATION_TYPE_DELEGATION:
		return DelegationTypeDelegation
	case storagepb.DelegationType_DELEGATION_TYPE_ERC20_VOTES:
		return DelegationTypeERC20Votes
	default:
		return DelegationTypeUnrecognized
	}
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

func (s *Server) GetTopDelegates(ctx context.Context, req *storagepb.GetTopDelegatesRequest) (*storagepb.GetTopDelegatesResponse, error) {
	if req.GetAddress() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}

	// delegations [dao_id: [summary, ...]]
	delegations, err := s.sp.getTopDelegates(ctx, req.GetAddress())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get delegations")
	}

	if len(delegations) == 0 {
		return &storagepb.GetTopDelegatesResponse{}, nil
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

	response := &storagepb.GetTopDelegatesResponse{
		List: make([]*storagepb.DelegatesSummary, 0, len(delegations)),
	}

	delegationsCnt := 0
	for _, di := range daoList.Daos {
		list, ok := delegations[di.ID.String()]
		if !ok {
			log.Warn().Msgf("dao info not found: %s", di.ID.String())
			continue
		}

		delegationsCnt += len(list)
		var delegatesInDao int32 = 0
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
		if len(list) > 0 {
			delegatesInDao += int32(list[0].MaxCnt)
		}

		response.List = append(response.List, &storagepb.DelegatesSummary{
			Dao:        dao.ConvertDaoToAPI(&di),
			List:       dl,
			TotalCount: delegatesInDao,
		})
	}

	response.TotalDelegatesCount = int32(delegationsCnt)

	return response, nil
}

func (s *Server) GetTopDelegators(ctx context.Context, req *storagepb.GetTopDelegatorsRequest) (*storagepb.GetTopDelegatorsResponse, error) {
	if req.GetAddress() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}

	// delegators [dao_id: [summary, ...]]
	delegators, err := s.sp.getTopDelegators(ctx, req.GetAddress())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get delegators")
	}

	if len(delegators) == 0 {
		return &storagepb.GetTopDelegatorsResponse{}, nil
	}

	daoIDs := slices.Collect(maps.Keys(delegators))
	daoList, err := s.ds.GetByFilters([]dao.Filter{
		dao.DaoIDsFilter{DaoIDs: daoIDs},
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get dao info")
	}

	addresses := make([]string, 0, len(delegators))
	for _, d := range delegators {
		for _, info := range d {
			addresses = append(addresses, info.AddressFrom)
		}
	}
	ensNames, err := s.sp.resolveAddressesName(addresses)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to resolve ens names")
	}

	response := &storagepb.GetTopDelegatorsResponse{
		List: make([]*storagepb.DelegatorSummary, 0, len(delegators)),
	}

	delegatorsCnt := 0
	for _, di := range daoList.Daos {
		list, ok := delegators[di.ID.String()]
		if !ok {
			log.Warn().Msgf("dao info not found: %s", di.ID.String())
			continue
		}

		delegatorsCnt += len(list)
		var delegatorsInDao int32 = 0
		dl := make([]*storagepb.DelegationDetails, 0, len(list))
		for _, d := range list {
			var expires *timestamppb.Timestamp
			if d.ExpiresAt != 0 {
				expires = timestamppb.New(time.Unix(d.ExpiresAt, 0))
			}

			dl = append(dl, &storagepb.DelegationDetails{
				Address:             d.AddressFrom,
				EnsName:             ensNames[d.AddressFrom],
				PercentOfDelegators: int32(d.Weight),
				Expiration:          expires,
			})
		}
		if len(list) > 0 {
			delegatorsInDao = int32(list[0].MaxCnt)
		}

		response.List = append(response.List, &storagepb.DelegatorSummary{
			Dao:        dao.ConvertDaoToAPI(&di),
			List:       dl,
			TotalCount: delegatorsInDao,
		})
	}

	response.TotalDelegatorsCount = int32(delegatorsCnt)

	return response, nil
}

func (s *Server) GetDelegationSummary(ctx context.Context, req *storagepb.GetDelegationSummaryRequest) (*storagepb.GetDelegationSummaryResponse, error) {
	if req.GetAddress() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}

	delegatorsCnt, err := s.sp.getDelegatorsCnt(ctx, req.GetAddress())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get delegators")
	}

	delegationsCnt, err := s.sp.getDelegatesCnt(ctx, req.GetAddress())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get delegators")
	}

	return &storagepb.GetDelegationSummaryResponse{
		TotalDelegatorsCount: delegatorsCnt,
		TotalDelegatesCount:  delegationsCnt,
	}, nil
}

func (s *Server) GetDelegatesByDao(_ context.Context, req *storagepb.GetDelegatesByDaoRequest) (*storagepb.GetDelegatesByDaoResponse, error) {
	if req.GetAddress() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}

	filters := []Filter{
		DelegatorFilter{Address: req.GetAddress()},
		DaoFilter{ID: req.GetDaoId()},
	}

	cnt, err := s.sp.GetCntByFilters(filters...)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get cnt by filters")
	}

	if cnt == 0 {
		return &storagepb.GetDelegatesByDaoResponse{}, nil
	}

	filters = append(filters,
		PageFilter{
			Limit:  int(req.GetLimit()),
			Offset: int(req.GetOffset()),
		},
		OrderByAddressToFilter{},
	)
	list, err := s.sp.GetByFilters(filters...)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get list by filters")
	}

	addresses := make([]string, 0, len(list))
	for _, d := range list {
		addresses = append(addresses, d.AddressTo)
	}
	ensNames, err := s.sp.resolveAddressesName(addresses)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to resolve ens names")
	}

	converted := make([]*storagepb.DelegationDetails, 0, len(list))
	for _, d := range list {
		var expires *timestamppb.Timestamp
		if d.ExpiresAt != 0 {
			expires = timestamppb.New(time.Unix(d.ExpiresAt, 0))
		}

		converted = append(converted, &storagepb.DelegationDetails{
			Address:             d.AddressTo,
			EnsName:             ensNames[d.AddressTo],
			PercentOfDelegators: int32(d.Weight),
			Expiration:          expires,
		})
	}

	return &storagepb.GetDelegatesByDaoResponse{
		List:       converted,
		TotalCount: int32(cnt),
	}, nil
}

func (s *Server) GetDelegatorsByDao(_ context.Context, req *storagepb.GetDelegatorsByDaoRequest) (*storagepb.GetDelegatorsByDaoResponse, error) {
	if req.GetAddress() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}

	filters := []Filter{
		DelegateFilter{Address: req.GetAddress()},
		DaoFilter{ID: req.GetDaoId()},
	}

	cnt, err := s.sp.GetCntByFilters(filters...)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get cnt by filters")
	}

	if cnt == 0 {
		return &storagepb.GetDelegatorsByDaoResponse{}, nil
	}

	filters = append(filters,
		PageFilter{
			Limit:  int(req.GetLimit()),
			Offset: int(req.GetOffset()),
		},
		OrderByAddressToFilter{},
	)
	list, err := s.sp.GetByFilters(filters...)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get list by filters")
	}

	addresses := make([]string, 0, len(list))
	for _, d := range list {
		addresses = append(addresses, d.AddressFrom)
	}
	ensNames, err := s.sp.resolveAddressesName(addresses)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to resolve ens names")
	}

	converted := make([]*storagepb.DelegationDetails, 0, len(list))
	for _, d := range list {
		var expires *timestamppb.Timestamp
		if d.ExpiresAt != 0 {
			expires = timestamppb.New(time.Unix(d.ExpiresAt, 0))
		}

		converted = append(converted, &storagepb.DelegationDetails{
			Address:             d.AddressFrom,
			EnsName:             ensNames[d.AddressFrom],
			PercentOfDelegators: int32(d.Weight),
			Expiration:          expires,
		})
	}

	return &storagepb.GetDelegatorsByDaoResponse{
		List:       converted,
		TotalCount: int32(cnt),
	}, nil
}

// GetDelegators returns list of erc20 delegators based on internal info ordered by voting power desc
func (s *Server) GetDelegators(
	ctx context.Context,
	req *storagepb.GetDelegatorsRequest,
) (*storagepb.GetDelegatorsResponse, error) {
	if req.GetAddress() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}

	daoID, err := uuid.Parse(req.GetDaoId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid dao ID format")
	}

	delegate, err := s.sp.GetErc20Delegate(ctx, daoID, req.GetChainId(), req.GetAddress())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get erc20 delegate")
	}

	list, err := s.sp.GetDelegators(ctx, ERC20DelegatorsRequest{
		Address: req.GetAddress(),
		ChainID: req.GetChainId(),
		DaoID:   daoID,
		Limit:   int(req.GetLimit()),
		Offset:  int(req.GetOffset()),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get delegators")
	}

	addresses := make([]string, 0, len(list))
	for _, d := range list {
		addresses = append(addresses, d.Address)
	}
	ensNames, err := s.sp.resolveAddressesName(addresses)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to resolve ens names")
	}

	converted := make([]*storagepb.DelegatorEntry, 0, len(list))
	for _, info := range list {
		converted = append(converted, &storagepb.DelegatorEntry{
			Address:    info.Address,
			EnsName:    ensNames[info.Address],
			TokenValue: info.TokenValue,
		})
	}

	return &storagepb.GetDelegatorsResponse{
		List:       converted,
		TotalCount: int32(delegate.RepresentedCnt),
	}, nil
}
