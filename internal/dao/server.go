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

	id, err := uuid.Parse(req.GetDaoId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid dao ID format")
	}

	dao, err := s.sp.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.InvalidArgument, "invalid dao ID")
	}

	if err != nil {
		log.Error().Err(err).Msgf("get dao by id: %s", req.GetDaoId())
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &storagepb.DaoByIDResponse{
		Dao: convertDaoToAPI(dao),
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
		res.Daos[i] = convertDaoToAPI(&info)
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
			info.Daos[i] = convertDaoToAPI(&dao)
		}

		res.Categories[idx] = info

		idx++
	}

	return res, nil
}

func convertDaoToAPI(dao *Dao) *storagepb.DaoInfo {
	return &storagepb.DaoInfo{
		Id:             dao.ID.String(),
		Alias:          dao.OriginalID,
		CreatedAt:      timestamppb.New(dao.CreatedAt),
		UpdatedAt:      timestamppb.New(dao.UpdatedAt),
		Name:           dao.Name,
		Private:        dao.Private,
		About:          dao.About,
		Avatar:         dao.Avatar,
		Terms:          dao.Terms,
		Location:       dao.Location,
		Website:        dao.Website,
		Twitter:        dao.Twitter,
		Github:         dao.Github,
		Coingeko:       dao.Coingecko,
		Email:          dao.Email,
		Network:        dao.Network,
		Symbol:         dao.Symbol,
		Skin:           dao.Skin,
		Domain:         dao.Domain,
		Strategies:     convertStrategiesToAPI(dao.Strategies),
		Voting:         convertVotingToAPI(dao.Voting),
		Categories:     dao.Categories,
		Treasuries:     convertTreasuriesToAPI(dao.Treasures),
		FollowersCount: uint64(dao.VotersCount),
		ProposalsCount: uint64(dao.ProposalsCount),
		Guidelines:     dao.Guidelines,
		Template:       dao.Template,
		ActivitySince:  uint64(dao.ActivitySince),
		VotersCount:    uint64(dao.VotersCount),
		ActiveVotes:    uint64(dao.ActiveVotes),
		Verified:       dao.Verified,
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
