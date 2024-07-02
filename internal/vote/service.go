package vote

import (
	"context"
	"fmt"
	"slices"
	"time"

	protoany "github.com/golang/protobuf/ptypes/any"
	"github.com/google/uuid"
	"github.com/goverland-labs/goverland-datasource-snapshot/protocol/votingpb"
	coreevents "github.com/goverland-labs/goverland-platform-events/events/core"
	"github.com/rs/zerolog/log"
)

type Publisher interface {
	PublishJSON(ctx context.Context, subject string, obj any) error
}

type DataProvider interface {
	BatchCreate(data []Vote) error
	GetByFilters(filters []Filter, limit int, offset int, firstVoter string) (List, error)
	UpdateVotes(list []ResolvedAddress) error
	GetUnique(string, int64) ([]string, error)
	GetByVoter(string) ([]string, error)
}

type DaoProvider interface {
	GetIDByOriginalID(string) (uuid.UUID, error)
}

type EnsResolver interface {
	AddRequests(list []string)
}

type Service struct {
	repo        DataProvider
	dao         DaoProvider
	events      Publisher
	ensResolver EnsResolver
	dsClient    votingpb.VotingClient
}

func NewService(
	r DataProvider,
	dp DaoProvider,
	p Publisher,
	er EnsResolver,
	dsClient votingpb.VotingClient,
) (*Service, error) {
	return &Service{
		repo:        r,
		dao:         dp,
		events:      p,
		ensResolver: er,
		dsClient:    dsClient,
	}, nil
}

func (s *Service) HandleVotes(ctx context.Context, votes []Vote) error {
	list := make(map[string]uuid.UUID)
	now := time.Now()
	authors := make([]string, 0, len(votes))
	for i := range votes {
		if !slices.Contains(authors, votes[i].Voter) {
			authors = append(authors, votes[i].Voter)
		}

		if daoID, ok := list[votes[i].OriginalDaoID]; ok {
			votes[i].DaoID = daoID
			continue
		}

		daoID, err := s.dao.GetIDByOriginalID(votes[i].OriginalDaoID)
		if err != nil {
			log.Error().Err(err).Msgf("get dao by original id: %s", votes[i].OriginalDaoID)

			return err
		}

		list[votes[i].OriginalDaoID] = daoID
		votes[i].DaoID = daoID
	}
	log.Info().Msgf("Gy80sHESRX: prepare votes: %f", time.Since(now).Seconds())

	now = time.Now()
	if err := s.repo.BatchCreate(votes); err != nil {
		return fmt.Errorf("can't create votes: %w", err)
	}
	log.Info().Msgf("Gy80sHESRX: create votes in DB: %f", time.Since(now).Seconds())

	now = time.Now()
	if err := s.events.PublishJSON(ctx, coreevents.SubjectVoteCreated, convertToCoreEvent(votes)); err != nil {
		log.Error().Err(err).Msgf("publish votes event")
	}
	log.Info().Msgf("Gy80sHESRX: publishing took: %f", time.Since(now).Seconds())

	s.ensResolver.AddRequests(authors)

	return nil
}

func (s *Service) GetByFilters(filters []Filter, limit int, offset int, firstVoter string) (List, error) {
	list, err := s.repo.GetByFilters(filters, limit, offset, firstVoter)
	if err != nil {
		return List{}, fmt.Errorf("get by filters: %w", err)
	}

	return list, nil
}

func (s *Service) Validate(ctx context.Context, req ValidateRequest) (ValidateResponse, error) {
	resp, err := s.dsClient.Validate(ctx, &votingpb.ValidateRequest{
		Voter:    req.Voter,
		Proposal: req.Proposal,
	})
	if err != nil {
		return ValidateResponse{}, fmt.Errorf("validate: %w", err)
	}

	var validationError *ValidationError
	if resp.ValidationError != nil {
		validationError = &ValidationError{
			Message: resp.GetValidationError().GetMessage(),
			Code:    resp.GetValidationError().GetCode(),
		}
	}

	return ValidateResponse{
		OK:              resp.GetOk(),
		VotingPower:     resp.GetVotingPower(),
		ValidationError: validationError,
	}, nil
}

func (s *Service) Prepare(ctx context.Context, req PrepareRequest) (PrepareResponse, error) {
	resp, err := s.dsClient.Prepare(ctx, &votingpb.PrepareRequest{
		Voter:    req.Voter,
		Proposal: req.Proposal,
		Choice: &protoany.Any{
			Value: req.Choice,
		},
		Reason: req.Reason,
	})
	if err != nil {
		return PrepareResponse{}, fmt.Errorf("prepare: %w", err)
	}

	return PrepareResponse{
		ID:        resp.GetId(),
		TypedData: resp.GetTypedData(),
	}, nil
}

func (s *Service) Vote(ctx context.Context, req VoteRequest) (VoteResponse, error) {
	resp, err := s.dsClient.Vote(ctx, &votingpb.VoteRequest{
		Id:  req.ID,
		Sig: req.Sig,
	})
	if err != nil {
		return VoteResponse{}, fmt.Errorf("vote: %w", err)
	}

	return VoteResponse{
		ID:   resp.GetId(),
		IPFS: resp.GetIpfs(),
		Relayer: Relayer{
			Address: resp.GetRelayer().GetAddress(),
			Receipt: resp.GetRelayer().GetReceipt(),
		},
	}, nil
}

func (s *Service) FetchAndStoreVote(ctx context.Context, id string) *Vote {
	vote, err := s.dsClient.GetVote(ctx, &votingpb.GetVoteRequest{
		Id: id,
	})
	if err != nil {
		log.Error().Err(err).Msgf("fetch vote: %s", id)
		return nil
	}

	internal := convertFromProtoToInternal(vote)
	err = s.HandleVotes(ctx, []Vote{internal})
	if err != nil {
		log.Error().Err(err).Msgf("store vote: %s", id)
	}

	return &internal
}

func (s *Service) HandleResolvedAddresses(list []ResolvedAddress) error {
	if len(list) == 0 {
		return nil
	}

	if err := s.repo.UpdateVotes(list); err != nil {
		return fmt.Errorf("s.repo.UpdateVotes: %w", err)
	}

	return nil
}

func (s *Service) GetDaosVotedIn(voter string) ([]string, error) {
	return s.repo.GetByVoter(voter)
}

func convertFromProtoToInternal(pl *votingpb.GetVoteResponse) Vote {
	return Vote{
		ID:            pl.GetId(),
		Ipfs:          pl.GetIpfs(),
		OriginalDaoID: pl.GetOriginalDaoId(),
		ProposalID:    pl.GetProposalId(),
		Voter:         pl.GetVoter(),
		Created:       int(pl.GetCreated()),
		Reason:        pl.GetReason(),
		Choice:        pl.GetChoice().GetValue(),
		App:           pl.GetApp(),
		Vp:            pl.GetVp(),
		VpByStrategy:  pl.GetVpByStrategy(),
		VpState:       pl.GetVpState(),
	}
}
