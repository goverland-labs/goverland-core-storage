package vote

import (
	"context"

	protoany "github.com/golang/protobuf/ptypes/any"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goverland-labs/goverland-core-storage/protocol/storagebp"
)

const (
	defaultLimit  = 50
	defaultOffset = 0
)

type Server struct {
	storagebp.UnimplementedVoteServer

	sp *Service
}

func NewServer(sp *Service) *Server {
	return &Server{
		sp: sp,
	}
}

func (s *Server) GetVotes(_ context.Context, req *storagebp.VotesFilterRequest) (*storagebp.VotesFilterResponse, error) {
	limit, offset := defaultLimit, defaultOffset
	if req.GetLimit() > 0 {
		limit = int(req.GetLimit())
	}
	if req.GetOffset() > 0 {
		offset = int(req.GetOffset())
	}
	filters := []Filter{
		PageFilter{Limit: limit, Offset: offset},
	}

	if req.GetOrderByVoter() != "" {
		filters = append(filters, OrderByVoterAndCreatedFilter{Voter: req.GetOrderByVoter()})
	} else {
		filters = append(filters, OrderByCreatedFilter{})
	}

	if req.GetProposalIds() != nil {
		filters = append(filters, ProposalIDsFilter{ProposalIDs: req.GetProposalIds()})
	}

	if req.GetVoter() != "" {
		filters = append(filters, VoterFilter{Voter: req.GetVoter()})
	}

	list, err := s.sp.GetByFilters(filters)
	if err != nil {
		log.Error().Err(err).Msgf("get votes by filter: %+v", req)
		return nil, status.Error(codes.Internal, "internal error")
	}

	res := &storagebp.VotesFilterResponse{
		Votes:      make([]*storagebp.VoteInfo, len(list.Votes)),
		TotalCount: uint64(list.TotalCount),
	}

	for i, info := range list.Votes {
		res.Votes[i] = convertVoteToAPI(&info)
	}

	return res, nil
}

func (s *Server) Validate(ctx context.Context, req *storagebp.ValidateRequest) (*storagebp.ValidateResponse, error) {
	validateResp, err := s.sp.Validate(ctx, ValidateRequest{
		Proposal: req.GetProposal(),
		Voter:    req.GetVoter(),
	})
	if err != nil {
		log.Error().Err(err).Msgf("validate vote: %+v", req)

		return nil, status.Error(codes.Internal, "failed to validate vote")
	}

	var validationError *storagebp.ValidationError
	if validateResp.ValidationError != nil {
		validationError = &storagebp.ValidationError{
			Message: validateResp.ValidationError.Message,
			Code:    validateResp.ValidationError.Code,
		}
	}

	return &storagebp.ValidateResponse{
		Ok:              validateResp.OK,
		VotingPower:     validateResp.VotingPower,
		ValidationError: validationError,
	}, nil
}

func (s *Server) Prepare(ctx context.Context, req *storagebp.PrepareRequest) (*storagebp.PrepareResponse, error) {
	prepareResp, err := s.sp.Prepare(ctx, PrepareRequest{
		Voter:    req.GetVoter(),
		Proposal: req.GetProposal(),
		Choice:   req.GetChoice().GetValue(),
		Reason:   req.Reason,
	})
	if err != nil {
		log.Error().Err(err).Msgf("prepare vote: %+v", req)

		return nil, status.Error(codes.Internal, "failed to prepare vote")
	}

	return &storagebp.PrepareResponse{
		Id:        prepareResp.ID,
		TypedData: prepareResp.TypedData,
	}, nil
}

func (s *Server) Vote(ctx context.Context, req *storagebp.VoteRequest) (*storagebp.VoteResponse, error) {
	voteResp, err := s.sp.Vote(ctx, VoteRequest{
		ID:  req.GetId(),
		Sig: req.GetSig(),
	})
	if err != nil {
		log.Error().Err(err).Msgf("vote: %+v", req)

		return nil, status.Error(codes.Internal, "failed to vote")
	}

	return &storagebp.VoteResponse{
		Id:   voteResp.ID,
		Ipfs: voteResp.IPFS,
		Relayer: &storagebp.Relayer{
			Address: voteResp.Relayer.Address,
			Receipt: voteResp.Relayer.Receipt,
		},
	}, nil
}

func convertVoteToAPI(info *Vote) *storagebp.VoteInfo {
	vpByStrategies := make([]float32, len(info.VpByStrategy))
	for i := range info.VpByStrategy {
		vpByStrategies[i] = float32(info.VpByStrategy[i])
	}

	return &storagebp.VoteInfo{
		Id:         info.ID,
		Ipfs:       info.Ipfs,
		Voter:      info.Voter,
		EnsName:    info.EnsName,
		Created:    uint64(info.Created),
		DaoId:      info.DaoID.String(),
		ProposalId: info.ProposalID,
		Choice: &protoany.Any{
			Value: info.Choice,
		},
		Reason:       info.Reason,
		App:          info.App,
		Vp:           float32(info.Vp),
		VpByStrategy: vpByStrategies,
		VpState:      info.VpState,
	}
}
