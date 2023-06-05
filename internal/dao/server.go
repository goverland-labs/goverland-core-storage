package dao

import (
	"context"
	"errors"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	proto "github.com/goverland-labs/core-storage/protobuf/internalapi"
)

type Servicer interface {
	GetByID(id string) (*Dao, error)
}

type Server struct {
	proto.UnimplementedDaoServer

	sp Servicer
}

func NewServer(sp Servicer) *Server {
	return &Server{
		sp: sp,
	}
}

func (s *Server) GetByID(_ context.Context, req *proto.DaoByIDRequest) (*proto.DaoResponse, error) {
	if req.GetDaoId() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid dao ID")
	}

	dao, err := s.sp.GetByID(req.GetDaoId())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.InvalidArgument, "invalid dao ID")
	}

	if err != nil {
		log.Error().Err(err).Msgf("get dao by id: %s", req.GetDaoId())
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &proto.DaoResponse{
		Dao: convertDaoToAPI(dao),
	}, nil
}

func convertDaoToAPI(dao *Dao) *proto.DaoInfo {
	return &proto.DaoInfo{
		Id:             dao.ID,
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
		FollowersCount: uint64(dao.FollowersCount),
		ProposalsCount: uint64(dao.ProposalsCount),
		Guidelines:     dao.Guidelines,
		Template:       dao.Template,
		ParentId:       dao.ParentID,
	}
}

func convertStrategiesToAPI(data Strategies) []*proto.Strategy {
	res := make([]*proto.Strategy, len(data))
	for i, info := range data {
		res[i] = &proto.Strategy{
			Name:    info.Name,
			Network: info.Network,
		}
	}

	return res
}

func convertTreasuriesToAPI(data Treasuries) []*proto.Treasury {
	res := make([]*proto.Treasury, len(data))
	for i, info := range data {
		res[i] = &proto.Treasury{
			Name:    info.Name,
			Address: info.Address,
			Network: info.Network,
		}
	}

	return res
}

func convertVotingToAPI(voting Voting) *proto.Voting {
	return &proto.Voting{
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
