package vote

import (
	"context"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	proto "github.com/goverland-labs/core-storage/protobuf/internalapi"
)

const (
	defaultLimit  = 50
	defaultOffset = 0
)

type Server struct {
	proto.UnimplementedVoteServer

	sp *Service
}

func NewServer(sp *Service) *Server {
	return &Server{
		sp: sp,
	}
}

func (s *Server) GetVotes(_ context.Context, req *proto.VotesFilterRequest) (*proto.VotesFilterResponse, error) {
	limit, offset := defaultLimit, defaultOffset
	if req.GetLimit() > 0 {
		limit = int(req.GetLimit())
	}
	if req.GetOffset() > 0 {
		offset = int(req.GetOffset())
	}
	filters := []Filter{
		PageFilter{Limit: limit, Offset: offset},
		OrderByCreatedFilter{},
	}

	if req.GetProposalId() != "" {
		filters = append(filters, ProposalFilter{ProposalID: req.GetProposalId()})
	}

	list, err := s.sp.GetByFilters(filters)
	if err != nil {
		log.Error().Err(err).Msgf("get votes by filter: %+v", req)
		return nil, status.Error(codes.Internal, "internal error")
	}

	res := &proto.VotesFilterResponse{
		Votes:      make([]*proto.VoteInfo, len(list.Votes)),
		TotalCount: uint64(list.TotalCount),
	}

	for i, info := range list.Votes {
		res.Votes[i] = convertVoteToAPI(&info)
	}

	return res, nil
}

func convertVoteToAPI(info *Vote) *proto.VoteInfo {
	return &proto.VoteInfo{
		ProposalId: info.ProposalID,
		Ipfs:       info.Ipfs,
		Voter:      info.Voter,
		Created:    uint64(info.Created),
		Reason:     info.Reason,
	}
}
