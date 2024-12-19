package vote

import (
	"context"
	"errors"

	"gorm.io/gorm"

	protoany "github.com/golang/protobuf/ptypes/any"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goverland-labs/goverland-core-storage/protocol/storagepb"

	"github.com/goverland-labs/goverland-core-storage/internal/proposal"
)

const (
	defaultLimit  = 50
	defaultOffset = 0
)

type Server struct {
	storagepb.UnimplementedVoteServer

	sp *Service
}

func NewServer(sp *Service) *Server {
	return &Server{
		sp: sp,
	}
}

func (s *Server) GetVotes(_ context.Context, req *storagepb.VotesFilterRequest) (*storagepb.VotesFilterResponse, error) {
	limit, offset := defaultLimit, defaultOffset
	if req.GetLimit() > 0 {
		limit = int(req.GetLimit())
	}
	if req.GetOffset() > 0 {
		offset = int(req.GetOffset())
	}
	var filters []Filter

	if req.GetOrderByVoter() != "" {
		filters = append(filters, proposal.OrderFilter{
			Orders: []proposal.Order{OrderByVp, OrderByCreated},
		})
	} else {
		filters = append(filters, proposal.OrderFilter{
			Orders: []proposal.Order{OrderByCreated},
		})
	}

	if req.GetProposalIds() != nil {
		filters = append(filters, ProposalIDsFilter{ProposalIDs: req.GetProposalIds()})
	}

	if req.GetVoter() != "" {
		filters = append(filters, VoterFilter{Voter: req.GetVoter()})
	}

	if req.GetQuery() != "" {
		filters = append(filters, QueryFilter{Query: req.GetQuery()})
	}

	if req.GetDaoId() != "" {
		filters = append(filters, DaoIDFilter{DaoID: req.GetDaoId()})
	}

	list, err := s.sp.GetByFilters(filters, limit, offset, req.GetOrderByVoter())
	if err != nil {
		log.Error().Err(err).Msgf("get votes by filter: %+v", req)
		return nil, status.Error(codes.Internal, "internal error")
	}

	res := &storagepb.VotesFilterResponse{
		Votes:      make([]*storagepb.VoteInfo, len(list.Votes)),
		TotalCount: uint64(list.TotalCount),
		TotalVp:    list.TotalVp,
	}

	for i, info := range list.Votes {
		res.Votes[i] = convertVoteToAPI(&info)
	}

	return res, nil
}

func (s *Server) Validate(ctx context.Context, req *storagepb.ValidateRequest) (*storagepb.ValidateResponse, error) {
	validateResp, err := s.sp.Validate(ctx, ValidateRequest{
		Proposal: req.GetProposal(),
		Voter:    req.GetVoter(),
	})
	if err != nil {
		log.Error().Err(err).Msgf("validate vote: %+v", req)

		return nil, status.Error(codes.Internal, "failed to validate vote")
	}

	var validationError *storagepb.ValidationError
	if validateResp.ValidationError != nil {
		validationError = &storagepb.ValidationError{
			Message: validateResp.ValidationError.Message,
			Code:    validateResp.ValidationError.Code,
		}
	}

	return &storagepb.ValidateResponse{
		Ok:              validateResp.OK,
		VotingPower:     validateResp.VotingPower,
		ValidationError: validationError,
		VoteStatus: &storagepb.VoteStatus{
			Voted: validateResp.VoteStatus.Voted,
			Choice: &protoany.Any{
				Value: validateResp.VoteStatus.Choice,
			},
		},
	}, nil
}

func (s *Server) Prepare(ctx context.Context, req *storagepb.PrepareRequest) (*storagepb.PrepareResponse, error) {
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

	return &storagepb.PrepareResponse{
		Id:        prepareResp.ID,
		TypedData: prepareResp.TypedData,
	}, nil
}

func (s *Server) Vote(ctx context.Context, req *storagepb.VoteRequest) (*storagepb.VoteResponse, error) {
	voteResp, err := s.sp.Vote(ctx, VoteRequest{
		ID:  req.GetId(),
		Sig: req.GetSig(),
	})
	if err != nil {
		log.Error().Err(err).Msgf("vote: %+v", req)

		return nil, status.Error(codes.Internal, "failed to vote")
	}

	vote := s.sp.FetchAndStoreVote(ctx, voteResp.ID)
	var proposalID string
	if vote != nil {
		proposalID = vote.ProposalID
	}

	return &storagepb.VoteResponse{
		Id:         voteResp.ID,
		Ipfs:       voteResp.IPFS,
		ProposalId: proposalID,
		Relayer: &storagepb.Relayer{
			Address: voteResp.Relayer.Address,
			Receipt: voteResp.Relayer.Receipt,
		},
	}, nil
}

func (s *Server) GetDaosVotedIn(_ context.Context, req *storagepb.DaosVotedInRequest) (*storagepb.DaosVotedInResponse, error) {
	if req.GetVoter() == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid voter")
	}

	daos, err := s.sp.GetDaosVotedIn(req.GetVoter())
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &storagepb.DaosVotedInResponse{
			DaoIds:     []string{},
			TotalCount: 0,
		}, nil
	}

	if err != nil {
		log.Error().Err(err).Msgf("get daos by voter: %s", req.GetVoter())
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &storagepb.DaosVotedInResponse{
		DaoIds:     daos,
		TotalCount: uint64(len(daos)),
	}, nil
}

func convertVoteToAPI(info *Vote) *storagepb.VoteInfo {
	vpByStrategies := make([]float32, len(info.VpByStrategy))
	for i := range info.VpByStrategy {
		vpByStrategies[i] = float32(info.VpByStrategy[i])
	}

	return &storagepb.VoteInfo{
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
