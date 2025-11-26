package delegate

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	proto "github.com/goverland-labs/goverland-core-storage/protocol/storagepb"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetDelegatesV2(ctx context.Context, req *proto.GetDelegatesV2Request) (*proto.GetDelegatesV2Response, error) {
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

	delegatesResponse, err := s.sp.getDelegatesMixed(ctx, GetDelegatesMixedRequest{
		DaoID:          daoID,
		QueryAccounts:  req.GetQueryAccounts(),
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

	list := make([]*proto.DelegatesWrapper, 0, len(delegatesResponse.List))
	var totalCnt int32
	for _, info := range delegatesResponse.List {
		delegates := make([]*proto.DelegateEntryV2, 0, len(info.Delegates))
		for _, d := range info.Delegates {
			delegates = append(delegates, convertDelegateToProto(d, info.DelegationType))
		}

		list = append(list, &proto.DelegatesWrapper{
			DaoId:          info.DaoID.String(),
			DelegationType: convertDelegationTypeToProto(info.DelegationType),
			ChainId:        info.ChainID,
			Delegates:      delegates,
			TotalCnt:       info.Total,
		})

		totalCnt += info.Total
	}

	return &proto.GetDelegatesV2Response{
		List:     list,
		TotalCnt: totalCnt,
	}, nil
}

func convertDelegateToProto(d Delegate, dt DelegationType) *proto.DelegateEntryV2 {
	entry := &proto.DelegateEntryV2{
		Address: d.Address,
		EnsName: d.ENSName,

		// todo: think about convertoring
		DelegatorCount:        &d.DelegatorCount,
		PercentOfDelegators:   &d.PercentOfDelegators,
		PercentOfVotingPower:  &d.PercentOfVotingPower,
		About:                 &d.About,
		Statement:             &d.Statement,
		VotesCount:            &d.VotesCount,
		CreatedProposalsCount: &d.CreatedProposalsCount,
	}

	switch dt {
	case DelegationTypeSplitDelegation:
		entry.VotingPower = &d.VotingPower
		// think about expiration, we don't have it in the response
		entry.Expiration = nil
	case DelegationTypeERC20Votes:
		// get dao, get strategy - do we need it?
		entry.TokenValue = &proto.TokenValue{
			Value:    fmt.Sprintf("%f", d.VotingPower),
			Symbol:   "",
			Decimals: 0,
		}
	}

	return entry
}

func convertDelegationTypeToProto(dt DelegationType) proto.DelegationType {
	switch dt {
	case DelegationTypeSplitDelegation:
		return proto.DelegationType_DELEGATION_TYPE_SPLIT_DELEGATION
	case DelegationTypeDelegation:
		return proto.DelegationType_DELEGATION_TYPE_DELEGATION
	case DelegationTypeERC20Votes:
		return proto.DelegationType_DELEGATION_TYPE_ERC20_VOTES
	default:
		return proto.DelegationType_DELEGATION_TYPE_UNRECOGNIZED
	}
}

func (s *Server) GetDelegatorsV2(ctx context.Context, req *proto.GetDelegatorsV2Request) (*proto.GetDelegatorsV2Response, error) {
	if req.GetAddress() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}

	daoID, err := uuid.Parse(req.GetDaoId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid dao ID format")
	}

	var queryAccs []string
	if req.GetAddress() != "" {
		queryAccs = append(queryAccs, req.GetAddress())
	}

	var chainID *string
	if req.GetChainId() != "" {
		chainID = &req.ChainId
	}

	delegatesResponse, err := s.sp.getDelegatorsMixed(ctx, GetDelegatesMixedRequest{
		DaoID:          daoID,
		QueryAccounts:  queryAccs,
		Limit:          int(req.GetLimit()),
		Offset:         int(req.GetOffset()),
		ChainID:        chainID,
		DelegationType: convertDelegationType(req.GetDelegationType()),
	})
	if err != nil {
		log.Error().
			Err(err).
			Str("dao_id", daoID.String()).
			Msg("failed to get getDelegatorsMixed")

		return nil, status.Errorf(codes.Internal, "failed to get delegators: %v", err)
	}

	list := make([]*proto.DelegatesWrapper, 0, len(delegatesResponse.List))
	var totalCnt int32
	for _, info := range delegatesResponse.List {
		delegates := make([]*proto.DelegateEntryV2, 0, len(info.Delegates))
		for _, d := range info.Delegates {
			delegates = append(delegates, convertDelegateToProto(d, info.DelegationType))
		}

		list = append(list, &proto.DelegatesWrapper{
			DaoId:          info.DaoID.String(),
			DelegationType: convertDelegationTypeToProto(info.DelegationType),
			ChainId:        info.ChainID,
			Delegates:      delegates,
			TotalCnt:       info.Total,
		})

		totalCnt += info.Total
	}

	return &proto.GetDelegatorsV2Response{
		List:     list,
		TotalCnt: totalCnt,
	}, nil
}

func (s *Server) GetTopDelegatesV2(ctx context.Context, req *proto.GetTopDelegatesV2Request) (*proto.GetTopDelegatesV2Response, error) {
	if req.GetAddress() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}

	// delegations [dao_id: [summary, ...]]
	//delegations, err := s.sp.getTopDelegates(ctx, req.GetAddress())
	//if err != nil {
	//	return nil, status.Error(codes.Internal, "failed to get delegations")
	//}
	//
	//if len(delegations) == 0 {
	//	return &storagepb.GetTopDelegatesResponse{}, nil
	//}
	//
	//daoIDs := slices.Collect(maps.Keys(delegations))
	//daoList, err := s.ds.GetByFilters([]dao.Filter{
	//	dao.DaoIDsFilter{DaoIDs: daoIDs},
	//})
	//if err != nil {
	//	return nil, status.Error(codes.Internal, "failed to get dao info")
	//}
	//
	//addresses := make([]string, 0, len(delegations))
	//for _, d := range delegations {
	//	for _, info := range d {
	//		addresses = append(addresses, info.AddressTo)
	//	}
	//}
	//ensNames, err := s.sp.resolveAddressesName(addresses)
	//if err != nil {
	//	return nil, status.Error(codes.Internal, "failed to resolve ens names")
	//}
	//
	//response := &storagepb.GetTopDelegatesResponse{
	//	List: make([]*storagepb.DelegatesSummary, 0, len(delegations)),
	//}
	//
	//delegationsCnt := 0
	//for _, di := range daoList.Daos {
	//	list, ok := delegations[di.ID.String()]
	//	if !ok {
	//		log.Warn().Msgf("dao info not found: %s", di.ID.String())
	//		continue
	//	}
	//
	//	delegationsCnt += len(list)
	//	var delegatesInDao int32 = 0
	//	dl := make([]*storagepb.DelegationDetails, 0, len(list))
	//	for _, d := range list {
	//		var expires *timestamppb.Timestamp
	//		if d.ExpiresAt != 0 {
	//			expires = timestamppb.New(time.Unix(d.ExpiresAt, 0))
	//		}
	//
	//		dl = append(dl, &storagepb.DelegationDetails{
	//			Address:             d.AddressTo,
	//			EnsName:             ensNames[d.AddressTo],
	//			PercentOfDelegators: int32(d.Weight),
	//			Expiration:          expires,
	//		})
	//	}
	//	if len(list) > 0 {
	//		delegatesInDao += int32(list[0].MaxCnt)
	//	}
	//
	//	response.List = append(response.List, &storagepb.DelegatesSummary{
	//		Dao:        dao.ConvertDaoToAPI(&di),
	//		List:       dl,
	//		TotalCount: delegatesInDao,
	//	})
	//}
	//
	//response.TotalDelegatesCount = int32(delegationsCnt)

	return nil, errors.New("implement me")
}

func (s *Server) GetTopDelegatorsV2(ctx context.Context, req *proto.GetTopDelegatorsV2Request) (*proto.GetTopDelegatorsV2Response, error) {
	if req.GetAddress() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid address")
	}

	resp, err := s.sp.getTopDelegatorsMixed(ctx, req.GetAddress(), req.GetDaoId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get delegators")
	}

	if len(resp.List) == 0 {
		return &proto.GetTopDelegatorsV2Response{}, nil
	}

	list := make([]*proto.DelegatesWrapper, 0, len(resp.List))
	var totalCnt int32
	for _, info := range resp.List {
		delegates := make([]*proto.DelegateEntryV2, 0, len(info.Delegates))
		for _, d := range info.Delegates {
			delegates = append(delegates, convertDelegateToProto(d, info.DelegationType))
		}

		list = append(list, &proto.DelegatesWrapper{
			DaoId:          info.DaoID.String(),
			DelegationType: convertDelegationTypeToProto(info.DelegationType),
			ChainId:        info.ChainID,
			Delegates:      delegates,
			TotalCnt:       info.Total,
		})

		totalCnt += info.Total
	}

	return &proto.GetTopDelegatorsV2Response{
		List:     list,
		TotalCnt: totalCnt,
	}, nil
}
