package proposal

import (
	"context"
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	proto "github.com/goverland-labs/core-storage/protobuf/internalapi"
)

const (
	defaultDaoLimit = 50
	defaultOffset   = 0
)

type Server struct {
	proto.UnimplementedProposalServer

	sp *Service
}

func NewServer(sp *Service) *Server {
	return &Server{
		sp: sp,
	}
}

func (s *Server) GetByID(_ context.Context, req *proto.ProposalByIDRequest) (*proto.ProposalByIDResponse, error) {
	if req.GetProposalId() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid proposal ID")
	}

	dao, err := s.sp.GetByID(req.GetProposalId())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.InvalidArgument, "invalid proposal ID")
	}

	if err != nil {
		log.Error().Err(err).Msgf("get dao by id: %s", req.GetProposalId())
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &proto.ProposalByIDResponse{
		Proposal: convertProposalToAPI(dao),
	}, nil
}

func (s *Server) GetByFilter(_ context.Context, req *proto.ProposalByFilterRequest) (*proto.ProposalByFilterResponse, error) {
	limit, offset := defaultDaoLimit, defaultOffset
	if req.GetLimit() > 0 {
		limit = int(req.GetLimit())
	}
	if req.GetOffset() > 0 {
		offset = int(req.GetOffset())
	}
	filters := []Filter{
		PageFilter{Limit: limit, Offset: offset},
	}

	if req.GetCategory() != "" {
		filters = append(filters, CategoriesFilter{Category: req.GetCategory()})
	}

	if req.GetDao() != "" {
		daos := strings.Split(req.GetDao(), ",")
		filters = append(filters, DaoIDsFilter{DaoIDs: daos})
	}

	list, err := s.sp.GetByFilters(filters)
	if err != nil {
		log.Error().Err(err).Msgf("get proposals by filter: %+v", req)
		return nil, status.Error(codes.Internal, "internal error")
	}

	res := &proto.ProposalByFilterResponse{
		Proposals:  make([]*proto.ProposalInfo, len(list.Proposals)),
		TotalCount: uint64(list.TotalCount),
	}

	for i, info := range list.Proposals {
		res.Proposals[i] = convertProposalToAPI(&info)
	}

	return res, nil
}

func convertProposalToAPI(info *Proposal) *proto.ProposalInfo {
	return &proto.ProposalInfo{
		Id:            info.ID,
		CreatedAt:     timestamppb.New(info.CreatedAt),
		UpdatedAt:     timestamppb.New(info.UpdatedAt),
		Ipfs:          info.Ipfs,
		Author:        info.Author,
		DaoId:         info.DaoID,
		Created:       uint64(info.Created),
		Network:       info.Network,
		Symbol:        info.Symbol,
		Type:          info.Type,
		Strategies:    convertStrategiesToAPI(info.Strategies),
		Title:         info.Title,
		Body:          info.Body,
		Discussion:    info.Discussion,
		Choices:       info.Choices,
		Start:         uint64(info.Start),
		End:           uint64(info.End),
		Quorum:        float32(info.Quorum),
		Privacy:       info.Privacy,
		Snapshot:      info.Snapshot,
		State:         info.State,
		Link:          info.Link,
		App:           info.App,
		Scores:        info.Scores,
		ScoresState:   info.ScoresState,
		ScoresTotal:   info.ScoresTotal,
		ScoresUpdated: uint64(info.ScoresUpdated),
		Votes:         uint64(info.Votes),
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
