package dao

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/goverland-labs/goverland-core-storage/pkg/sdk/zerion"

	"github.com/goverland-labs/goverland-core-storage/protocol/storagepb"
)

const (
	defaultDaoLimit           = 50
	defaultOffset             = 0
	defaultTopCategoriesLimit = 10
	maxPerTop                 = 20
)

type Server struct {
	storagepb.UnimplementedDaoServer

	sp *Service
}

func NewServer(sp *Service) *Server {
	return &Server{
		sp: sp,
	}
}

func (s *Server) GetByID(_ context.Context, req *storagepb.DaoByIDRequest) (*storagepb.DaoByIDResponse, error) {
	if req.GetDaoId() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid dao ID")
	}

	var (
		dao *Dao
		err error
	)
	if id, errV := uuid.Parse(req.GetDaoId()); errV == nil {
		dao, err = s.sp.GetByID(id)
	} else {
		dao, err = s.sp.GetDaoByOriginalID(req.GetDaoId())
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.InvalidArgument, "invalid dao ID")
	}

	if err != nil {
		log.Error().Err(err).Msgf("get dao by id: %s", req.GetDaoId())
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &storagepb.DaoByIDResponse{
		Dao: ConvertDaoToAPI(dao),
	}, nil
}

func (s *Server) GetByFilter(_ context.Context, req *storagepb.DaoByFilterRequest) (*storagepb.DaoByFilterResponse, error) {
	limit, offset := defaultDaoLimit, defaultOffset
	if req.GetLimit() > 0 {
		limit = int(req.GetLimit())
	}
	if req.GetOffset() > 0 {
		offset = int(req.GetOffset())
	}
	filters := []Filter{
		PageFilter{Limit: limit, Offset: offset},
		OrderByPopularityIndexFilter{},
	}

	if req.GetQuery() != "" {
		filters = append(filters, NameFilter{Name: req.GetQuery()})
	}

	if req.GetCategory() != "" {
		filters = append(filters, CategoryFilter{Category: req.GetCategory()})
	}

	if len(req.GetDaoIds()) != 0 {
		filters = append(filters, DaoIDsFilter{
			DaoIDs: req.GetDaoIds(),
		})
	}

	list, err := s.sp.GetByFilters(filters)
	if err != nil {
		log.Error().Err(err).Msgf("get daos by filter: %+v", req)
		return nil, status.Error(codes.Internal, "internal error")
	}

	res := &storagepb.DaoByFilterResponse{
		Daos:       make([]*storagepb.DaoInfo, len(list.Daos)),
		TotalCount: uint64(list.TotalCount),
	}

	for i, info := range list.Daos {
		res.Daos[i] = ConvertDaoToAPI(&info)
	}

	return res, nil
}

func (s *Server) GetTopByCategories(ctx context.Context, req *storagepb.TopByCategoriesRequest) (*storagepb.TopByCategoriesResponse, error) {
	limit := 10
	if req.GetLimit() != 0 {
		limit = int(req.GetLimit())
	}
	if limit > maxPerTop {
		limit = maxPerTop
	}

	list, err := s.sp.GetTopByCategories(ctx, limit)
	if err != nil {
		log.Error().Err(err).Msgf("get top by categories: %+v", req)
		return nil, status.Error(codes.Internal, "internal error")
	}

	res := &storagepb.TopByCategoriesResponse{
		Categories: make([]*storagepb.TopCategory, len(list)),
	}

	idx := 0
	for cat, details := range list {
		info := &storagepb.TopCategory{
			Category:   cat,
			TotalCount: uint64(details.Total),
			Daos:       make([]*storagepb.DaoInfo, len(details.List)),
		}
		for i, dao := range details.List {
			info.Daos[i] = ConvertDaoToAPI(&dao)
		}

		res.Categories[idx] = info

		idx++
	}

	return res, nil
}

func ConvertDaoToAPI(dao *Dao) *storagepb.DaoInfo {
	return &storagepb.DaoInfo{
		Id:                 dao.ID.String(),
		Alias:              dao.OriginalID,
		CreatedAt:          timestamppb.New(dao.CreatedAt),
		UpdatedAt:          timestamppb.New(dao.UpdatedAt),
		Name:               dao.Name,
		Private:            dao.Private,
		About:              dao.About,
		Avatar:             dao.Avatar,
		Terms:              dao.Terms,
		Location:           dao.Location,
		Website:            dao.Website,
		Twitter:            dao.Twitter,
		Github:             dao.Github,
		Coingeko:           dao.Coingecko,
		Email:              dao.Email,
		Network:            dao.Network,
		Symbol:             dao.Symbol,
		Skin:               dao.Skin,
		Domain:             dao.Domain,
		Strategies:         convertStrategiesToAPI(dao.Strategies),
		Voting:             convertVotingToAPI(dao.Voting),
		Categories:         dao.Categories,
		Treasuries:         convertTreasuriesToAPI(dao.Treasures),
		FollowersCount:     uint64(dao.VotersCount),
		ProposalsCount:     uint64(dao.ProposalsCount),
		Guidelines:         dao.Guidelines,
		Template:           dao.Template,
		ActivitySince:      uint64(dao.ActivitySince),
		VotersCount:        uint64(dao.VotersCount),
		ActiveVotes:        uint64(dao.ActiveVotes),
		Verified:           dao.Verified,
		PopularityIndex:    dao.PopularityIndex,
		ActiveProposalsIds: dao.ActiveProposalsIDs,
		TokenExist:         dao.FungibleId != "",
		TokenSymbol:        dao.TokenSymbol,
		// TODO: parentID
	}
}

func convertStrategiesToAPI(data Strategies) []*storagepb.Strategy {
	res := make([]*storagepb.Strategy, len(data))
	for i, info := range data {
		params, _ := json.Marshal(info.Params)

		res[i] = &storagepb.Strategy{
			Name:    info.Name,
			Network: info.Network,
			Params:  params,
		}
	}

	return res
}

func convertTreasuriesToAPI(data Treasuries) []*storagepb.Treasury {
	res := make([]*storagepb.Treasury, len(data))
	for i, info := range data {
		res[i] = &storagepb.Treasury{
			Name:    info.Name,
			Address: info.Address,
			Network: info.Network,
		}
	}

	return res
}

func convertVotingToAPI(voting Voting) *storagepb.Voting {
	return &storagepb.Voting{
		Delay:       uint64(voting.Delay),
		Period:      uint64(voting.Period),
		Type:        voting.Type,
		Quorum:      voting.Quorum,
		Blind:       voting.Blind,
		HideAbstain: voting.HideAbstain,
		Privacy:     voting.Privacy,
		Aliased:     voting.Aliased,
	}
}

func (s *Server) GetRecommendationsList(
	_ context.Context,
	_ *storagepb.GetRecommendationsListRequest,
) (*storagepb.GetRecommendationsListResponse, error) {
	list := s.sp.getRecommendations()
	resp := &storagepb.GetRecommendationsListResponse{
		List: make([]*storagepb.DaoRecommendationDetails, 0, len(list)),
	}

	for _, details := range list {
		resp.List = append(resp.List, &storagepb.DaoRecommendationDetails{
			OriginalId: details.OriginalId,
			InternalId: details.InternalId,
			Name:       details.Name,
			Symbol:     details.Symbol,
			NetworkId:  details.NetworkId,
			Address:    details.Address,
		})
	}

	return resp, nil
}

func (s *Server) GetTokenInfo(_ context.Context, req *storagepb.TokenInfoRequest) (*storagepb.TokenInfoResponse, error) {
	id, err := uuid.Parse(req.GetDaoId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid dao ID")
	}

	data, err := s.sp.GetTokenInfo(id)
	if err != nil {
		log.Error().Err(err).Msgf("get token info: %+v", req)
		return nil, status.Error(codes.Internal, "internal error")
	}

	tokenChains := make([]*storagepb.TokenChainInfo, 0, len(data.Chains))
	for _, info := range data.Chains {
		tokenChains = append(tokenChains, &storagepb.TokenChainInfo{
			ChainId:  info.ChainID,
			Name:     info.Name,
			Decimals: uint32(info.Decimals),
			IconUrl:  info.IconURL,
			Address:  info.Address,
		})
	}

	return &storagepb.TokenInfoResponse{
		Name:                  data.Name,
		Symbol:                data.Symbol,
		TotalSupply:           data.TotalSupply,
		CirculatingSupply:     data.CirculatingSupply,
		MarketCap:             data.MarketCap,
		FullyDilutedValuation: data.FullyDilutedValuation,
		Price:                 data.Price,
		FungibleId:            data.FungibleID,
		Chains:                tokenChains,
	}, nil
}

func (s *Server) GetTokenChart(_ context.Context, req *storagepb.TokenChartRequest) (*storagepb.TokenChartResponse, error) {
	id, err := uuid.Parse(req.GetDaoId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid dao ID")
	}
	data, err := s.sp.GetTokenChart(id, req.GetPeriod())
	if err != nil {
		log.Error().Err(err).Msgf("get token chart: %+v", req)
		return nil, status.Error(codes.Internal, "internal error")
	}
	return convertChartToAPI(data), nil
}

func (s *Server) PopulateTokenPrices(ctx context.Context, req *storagepb.TokenPricesRequest) (*storagepb.TokenPricesResponse, error) {
	id, err := uuid.Parse(req.GetDaoId())
	if err != nil {
		return &storagepb.TokenPricesResponse{Status: false}, status.Error(codes.InvalidArgument, "invalid dao ID")
	}
	ok, err := s.sp.PopulateTokenPrices(ctx, id)
	if err != nil {
		log.Error().Err(err).Msgf("populate token prices: %+v", req)
		return &storagepb.TokenPricesResponse{Status: false}, status.Error(codes.Internal, "internal error")
	}
	return &storagepb.TokenPricesResponse{Status: ok}, nil
}

func convertChartToAPI(data *zerion.ChartData) *storagepb.TokenChartResponse {
	var pc float64
	if data.ChartAttributes.Stats.First != 0 {
		pc = (data.ChartAttributes.Stats.Last - data.ChartAttributes.Stats.First) / data.ChartAttributes.Stats.First
	} else {
		pc = 0
	}

	points := make([]*storagepb.Point, len(data.ChartAttributes.Points))
	for i, info := range data.ChartAttributes.Points {
		points[i] = &storagepb.Point{
			Time:  timestamppb.New(info.Time),
			Price: info.Price,
		}
	}

	return &storagepb.TokenChartResponse{
		Price:        data.ChartAttributes.Stats.Last,
		PriceChanges: pc,
		Points:       points,
	}
}
