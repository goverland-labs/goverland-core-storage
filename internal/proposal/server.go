package proposal

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"

	"github.com/goverland-labs/goverland-core-storage/protocol/storagebp"
)

const (
	defaultDaoLimit = 50
	defaultOffset   = 0
)

type Server struct {
	storagebp.UnimplementedProposalServer

	sp *Service
}

func NewServer(sp *Service) *Server {
	return &Server{
		sp: sp,
	}
}

func (s *Server) GetByID(_ context.Context, req *storagebp.ProposalByIDRequest) (*storagebp.ProposalByIDResponse, error) {
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

	return &storagebp.ProposalByIDResponse{
		Proposal: convertProposalToAPI(dao),
	}, nil
}

func (s *Server) GetByFilter(_ context.Context, req *storagebp.ProposalByFilterRequest) (*storagebp.ProposalByFilterResponse, error) {
	limit, offset := defaultDaoLimit, defaultOffset
	if req.GetLimit() > 0 {
		limit = int(req.GetLimit())
	}
	if req.GetOffset() > 0 {
		offset = int(req.GetOffset())
	}
	filters := []Filter{
		SkipCanceled{},
		SkipSpamFilter{},
		PageFilter{Limit: limit, Offset: offset},
	}

	var list ProposalList
	var err error

	if req.GetTop() {
		list, err = s.sp.GetTop(limit, offset)
	} else {
		if req.GetCategory() != "" {
			filters = append(filters, CategoriesFilter{Category: req.GetCategory()})
		}

		if req.GetDao() != "" {
			daos := strings.Split(req.GetDao(), ",")
			filters = append(filters, DaoIDsFilter{DaoIDs: daos})
		}

		if req.GetTitle() != "" {
			filters = append(filters,
				TitleFilter{Title: req.GetTitle()},
				OrderFilter{
					Orders: []Order{
						OrderByStates,
						OrderByVotes,
					},
				})
		} else {
			filters = append(filters, OrderFilter{
				Orders: []Order{OrderByVotes},
			})
		}

		if len(req.GetProposalIds()) != 0 {
			filters = append(filters, ProposalIDsFilter{
				ProposalIDs: req.GetProposalIds(),
			})
		}

		list, err = s.sp.GetByFilters(filters)
	}
	if err != nil {
		log.Error().Err(err).Msgf("get proposals by filter: %+v", req)
		return nil, status.Error(codes.Internal, "internal error")
	}

	res := &storagebp.ProposalByFilterResponse{
		Proposals:  make([]*storagebp.ProposalInfo, len(list.Proposals)),
		TotalCount: uint64(list.TotalCount),
	}

	for i, info := range list.Proposals {
		res.Proposals[i] = convertProposalToAPI(&info)
	}

	return res, nil
}

func convertProposalToAPI(info *Proposal) *storagebp.ProposalInfo {
	return &storagebp.ProposalInfo{
		Id:            info.ID,
		CreatedAt:     timestamppb.New(info.CreatedAt),
		UpdatedAt:     timestamppb.New(info.UpdatedAt),
		Ipfs:          info.Ipfs,
		Author:        info.Author,
		EnsName:       info.EnsName,
		DaoId:         info.DaoID.String(),
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
		State:         string(info.State),
		Link:          info.Link,
		App:           info.App,
		Scores:        info.Scores,
		ScoresState:   info.ScoresState,
		ScoresTotal:   info.ScoresTotal,
		ScoresUpdated: uint64(info.ScoresUpdated),
		Votes:         uint64(info.Votes),
		Timeline:      convertTimelineToAPI(info.Timeline),
	}
}

func convertStrategiesToAPI(data Strategies) []*storagebp.Strategy {
	res := make([]*storagebp.Strategy, len(data))
	for i, info := range data {
		params, _ := json.Marshal(info.Params)

		res[i] = &storagebp.Strategy{
			Name:    info.Name,
			Network: info.Network,
			Params:  params,
		}
	}

	return res
}

func convertTimelineToAPI(tl Timeline) []*storagebp.ProposalTimelineItem {
	if len(tl) == 0 {
		return nil
	}

	res := make([]*storagebp.ProposalTimelineItem, len(tl))
	for i := range tl {
		res[i] = &storagebp.ProposalTimelineItem{
			CreatedAt: timestamppb.New(tl[i].CreatedAt),
			Action:    convertTimelineActionToAPI(tl[i].Action),
		}
	}

	return res
}

func convertTimelineActionToAPI(action TimelineAction) storagebp.ProposalTimelineItem_TimelineAction {
	switch action {
	case ProposalCreated:
		return storagebp.ProposalTimelineItem_ProposalCreated
	case ProposalUpdated:
		return storagebp.ProposalTimelineItem_ProposalUpdated
	case ProposalVotingStartsSoon:
		return storagebp.ProposalTimelineItem_ProposalVotingStartsSoon
	case ProposalVotingStarted:
		return storagebp.ProposalTimelineItem_ProposalVotingStarted
	case ProposalVotingQuorumReached:
		return storagebp.ProposalTimelineItem_ProposalVotingQuorumReached
	case ProposalVotingEnded:
		return storagebp.ProposalTimelineItem_ProposalVotingEnded
	case ProposalVotingEndsSoon:
		return storagebp.ProposalTimelineItem_ProposalVotingEndsSoon
	default:
		return storagebp.ProposalTimelineItem_Unspecified
	}
}
