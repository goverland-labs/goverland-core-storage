package delegate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	protoany "github.com/golang/protobuf/ptypes/any"
	"github.com/google/uuid"
	"github.com/goverland-labs/goverland-datasource-snapshot/protocol/delegatepb"
	events "github.com/goverland-labs/goverland-platform-events/events/core"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/goverland-labs/goverland-core-storage/internal/dao"
	"github.com/goverland-labs/goverland-core-storage/internal/ensresolver"
)

var errNoResolved = errors.New("no addresses resolved")

type DaoProvider interface {
	GetByID(id uuid.UUID) (*dao.Dao, error)
	GetIDByOriginalID(string) (uuid.UUID, error)
}

type Publisher interface {
	PublishJSON(ctx context.Context, subject string, obj any) error
}

type EventRegistered interface {
	EventExist(_ context.Context, id, t, event string) (bool, error)
	RegisterEvent(_ context.Context, id, t, event string) error
}

type EnsResolver interface {
	GetByNames(names []string) ([]ensresolver.EnsName, error)
	GetByAddresses(addresses []string) ([]ensresolver.EnsName, error)
	AddRequests(list []string)
}

type Service struct {
	delegateClient delegatepb.DelegateClient
	daoProvider    DaoProvider
	ensResolver    EnsResolver
	publisher      Publisher
	er             EventRegistered
	repo           *Repo
}

func NewService(repo *Repo, dc delegatepb.DelegateClient, daoProvider DaoProvider, ensResolver EnsResolver, ep Publisher, er EventRegistered) *Service {
	return &Service{
		delegateClient: dc,
		daoProvider:    daoProvider,
		ensResolver:    ensResolver,
		publisher:      ep,
		er:             er,
		repo:           repo,
	}
}

func (s *Service) GetDelegates(ctx context.Context, request GetDelegatesRequest) (*GetDelegatesResponse, error) {
	daoEntity, err := s.daoProvider.GetByID(request.DaoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dao: %w", err)
	}

	delegationStrategyJson, err := s.getDelegationStrategy(daoEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to get delegation strategy: %w", err)
	}

	queryAddresses, err := s.resolveQueryAccounts(request.QueryAccounts)
	if errors.Is(err, errNoResolved) {
		return &GetDelegatesResponse{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to resolve query accounts: %w", err)
	}

	resp, err := s.delegateClient.GetDelegates(ctx, &delegatepb.GetDelegatesRequest{
		DaoOriginalId: daoEntity.OriginalID,
		Strategy: &protoany.Any{
			Value: delegationStrategyJson,
		},
		Addresses: queryAddresses,
		Sort:      request.Sort,
		Limit:     int32(request.Limit),
		Offset:    int32(request.Offset),
	})
	if err != nil {
		return nil, err
	}

	respAddresses := make([]string, 0, len(resp.GetDelegates()))
	for _, d := range resp.GetDelegates() {
		respAddresses = append(respAddresses, d.GetAddress())
	}
	ensNames, err := s.resolveAddressesName(respAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve addresses names: %w", err)
	}

	delegates := make([]Delegate, 0, len(resp.Delegates))
	for _, d := range resp.GetDelegates() {
		delegates = append(delegates, Delegate{
			Address:               d.GetAddress(),
			ENSName:               ensNames[d.GetAddress()],
			DelegatorCount:        d.GetDelegatorCount(),
			PercentOfDelegators:   d.GetPercentOfDelegators(),
			VotingPower:           d.GetVotingPower(),
			PercentOfVotingPower:  d.GetPercentOfVotingPower(),
			About:                 d.GetAbout(),
			Statement:             d.GetStatement(),
			VotesCount:            0,
			CreatedProposalsCount: 0,
		})
	}

	return &GetDelegatesResponse{
		Delegates: delegates,
		Total:     resp.GetTotal(),
	}, nil
}

func (s *Service) GetDelegateProfile(ctx context.Context, request GetDelegateProfileRequest) (GetDelegateProfileResponse, error) {
	daoEntity, err := s.daoProvider.GetByID(request.DaoID)
	if err != nil {
		return GetDelegateProfileResponse{}, fmt.Errorf("failed to get dao: %w", err)
	}

	delegationStrategyJson, err := s.getDelegationStrategy(daoEntity)
	if err != nil {
		return GetDelegateProfileResponse{}, fmt.Errorf("failed to get delegation strategy: %w", err)
	}

	resp, err := s.delegateClient.GetDelegateProfile(ctx, &delegatepb.GetDelegateProfileRequest{
		DaoOriginalId: daoEntity.OriginalID,
		Strategy: &protoany.Any{
			Value: delegationStrategyJson,
		},
		Address: request.Address,
	})
	if err != nil {
		return GetDelegateProfileResponse{}, fmt.Errorf("failed to get delegate profile: %w", err)
	}

	delegatesAddresses := make([]string, 0, len(resp.GetDelegates()))
	for _, d := range resp.GetDelegates() {
		delegatesAddresses = append(delegatesAddresses, d.GetAddress())
	}
	ensNames, err := s.resolveAddressesName(delegatesAddresses)
	if err != nil {
		return GetDelegateProfileResponse{}, fmt.Errorf("failed to resolve addresses names: %w", err)
	}

	delegates := make([]ProfileDelegateItem, 0, len(resp.GetDelegates()))
	for _, d := range resp.GetDelegates() {
		delegates = append(delegates, ProfileDelegateItem{
			Address:        d.GetAddress(),
			ENSName:        ensNames[d.GetAddress()],
			Weight:         d.GetWeight(),
			DelegatedPower: d.GetDelegatedPower(),
		})
	}

	var expiration *time.Time
	if resp.GetExpiration() != nil {
		expirationTime := resp.GetExpiration().AsTime()
		expiration = &expirationTime
	}

	return GetDelegateProfileResponse{
		Address:              resp.GetAddress(),
		VotingPower:          resp.GetVotingPower(),
		IncomingPower:        resp.GetIncomingPower(),
		OutgoingPower:        resp.GetOutgoingPower(),
		PercentOfVotingPower: resp.GetPercentOfVotingPower(),
		PercentOfDelegators:  resp.GetPercentOfDelegators(),
		Delegates:            delegates,
		Expiration:           expiration,
	}, nil
}

func (s *Service) getDelegationStrategy(daoEntity *dao.Dao) ([]byte, error) {
	var delegationStrategy *dao.Strategy
	for _, strategy := range daoEntity.Strategies {
		if strategy.Name == "split-delegation" {
			delegationStrategy = &strategy
			break
		}
	}

	if delegationStrategy == nil {
		return nil, fmt.Errorf("delegation strategy not found for dao")
	}

	// TODO: avoid wrong naming, fix it in another way
	type marshalStrategy struct {
		Name    string                 `json:"name"`
		Network string                 `json:"network"`
		Params  map[string]interface{} `json:"params"`
	}

	ms := marshalStrategy{
		Name:    delegationStrategy.Name,
		Network: delegationStrategy.Network,
		Params:  delegationStrategy.Params,
	}

	delegationStrategyJson, err := json.Marshal(ms)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal delegation strategy: %w", err)
	}

	return delegationStrategyJson, nil
}

func (s *Service) resolveQueryAccounts(accs []string) ([]string, error) {
	var addresses []string
	var names []string
	for _, a := range accs {
		if common.IsHexAddress(a) {
			addresses = append(addresses, a)
			continue
		}

		names = append(names, a)
	}

	ensNames, err := s.ensResolver.GetByNames(names)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve ens names: %w", err)
	}

	for _, n := range ensNames {
		addresses = append(addresses, n.Address)
	}

	if len(accs) > 0 && len(addresses) == 0 {
		return nil, errNoResolved
	}

	return addresses, nil
}

func (s *Service) resolveAddressesName(addresses []string) (map[string]string, error) {
	ensNames, err := s.ensResolver.GetByAddresses(addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve ens names: %w", err)
	}

	res := make(map[string]string, len(ensNames))
	for _, n := range ensNames {
		// todo: use lowercase in all places
		for _, addr := range addresses {
			if strings.EqualFold(addr, n.Address) {
				res[addr] = n.Name
			}
		}
	}

	return res, nil
}

func (s *Service) handleDelegate(ctx context.Context, hr History) error {
	if err := s.repo.CallInTx(func(tx *gorm.DB) error {
		if hr.OriginalSpaceID == "" {
			log.Warn().Msgf("skip processing block %d from %s cause dao id is empty", hr.BlockNumber, hr.ChainID)

			return nil
		}

		// store to history
		if err := s.repo.CreateHistory(tx, hr); err != nil {
			return fmt.Errorf("repo.CreateHistory: %w", err)
		}

		// get space id by provided original_space_id
		daoID, err := s.daoProvider.GetIDByOriginalID(hr.OriginalSpaceID)
		if err != nil {
			return fmt.Errorf("dp.GetIDByOriginalID: %w", err)
		}

		bts, err := s.repo.GetSummaryBlockTimestamp(tx, strings.ToLower(hr.AddressFrom), daoID.String())
		if err != nil {
			return fmt.Errorf("s.repo.GetSummaryBlockTimestamp: %w", err)
		}

		// skip this block due to already processed
		if bts != 0 && bts >= hr.BlockTimestamp {
			log.Warn().Msgf("delegates: skip processing block %d from %s due to invalid timestamp", hr.BlockNumber, hr.ChainID)

			return nil
		}

		if hr.Action == actionExpire {
			if err := s.repo.UpdateSummaryExpiration(tx, strings.ToLower(hr.AddressFrom), daoID.String(), hr.Delegations.Expiration, hr.BlockTimestamp); err != nil {
				return fmt.Errorf("UpdateSummaryExpiration: %w", err)
			}

			return nil
		}

		if err := s.repo.RemoveSummary(tx, strings.ToLower(hr.AddressFrom), daoID.String()); err != nil {
			return fmt.Errorf("UpdateSummaryExpiration: %w", err)
		}

		if hr.Action == actionClear {
			return nil
		}

		addresses := make([]string, 0, len(hr.Delegations.Details)+1)
		addresses = append(addresses, hr.AddressFrom)

		for _, info := range hr.Delegations.Details {
			addresses = append(addresses, info.Address)
			if err = s.repo.CreateSummary(tx, Summary{
				AddressFrom:        strings.ToLower(hr.AddressFrom),
				AddressTo:          strings.ToLower(info.Address),
				DaoID:              daoID.String(),
				Weight:             info.Weight,
				LastBlockTimestamp: hr.BlockTimestamp,
				ExpiresAt:          int64(hr.Delegations.Expiration),
			}); err != nil {
				return fmt.Errorf("createSummary [%s/%s/%s]: %w", hr.AddressFrom, info.Address, daoID.String(), err)
			}
		}

		for _, info := range hr.Delegations.Details {
			event := events.DelegatePayload{
				Initiator: strings.ToLower(hr.AddressFrom),
				Delegator: strings.ToLower(info.Address),
				DaoID:     daoID,
			}

			if err = s.publisher.PublishJSON(ctx, events.SubjectDelegateCreated, event); err != nil {
				log.Warn().Err(err).Msgf("failed to publish delegate payload")
			}
		}

		go func(list []string) {
			s.ensResolver.AddRequests(list)
		}(addresses)

		return nil
	}); err != nil {
		return fmt.Errorf("repo.CallInTx: %w", err)
	}

	return nil
}

func (s *Service) handleProposalCreated(ctx context.Context, pr Proposal) error {
	// get space id by provided original_space_id
	daoID, err := s.daoProvider.GetIDByOriginalID(pr.OriginalDaoID)
	if err != nil {
		return fmt.Errorf("dp.GetIDByOriginalID: %w", err)
	}

	// find delegator by author in specific space id
	summary, err := s.repo.FindDelegator(daoID.String(), strings.ToLower(pr.Author))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("repo.FindDelegator: %w", err)
	}

	// author doesn't have any delegation relations
	if summary == nil {
		return nil
	}

	if summary.SelfDelegation() {
		return nil
	}

	// delegation is expired
	if summary.Expired() {
		return nil
	}

	// make an event
	event := events.DelegatePayload{
		Initiator:  strings.ToLower(pr.Author),
		Delegator:  summary.AddressFrom,
		DaoID:      daoID,
		ProposalID: pr.ID,
	}
	if err = s.publisher.PublishJSON(ctx, events.SubjectDelegateCreateProposal, event); err != nil {
		return fmt.Errorf("s.publisher.PublishJSON: %w", err)
	}

	return nil
}

func (s *Service) handleVotesCreated(ctx context.Context, batch []Vote) error {
	summary, err := s.repo.FindDelegatorsByVotes(batch)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("repo.FindDelegatorsByVotes: %w", err)
	}

	for _, info := range summary {
		if info.SelfDelegation() {
			continue
		}

		// delegation is expired
		if info.Expired() {
			continue
		}

		// make an event
		event := events.DelegatePayload{
			Initiator:  strings.ToLower(info.AddressTo),
			Delegator:  info.AddressFrom,
			DaoID:      uuid.MustParse(info.DaoID),
			ProposalID: info.ProposalID,
		}

		if err = s.publisher.PublishJSON(ctx, events.SubjectDelegateVotingVoted, event); err != nil {
			log.Err(err).Msgf("publish delegate voted: %s %s", info.AddressTo, info.ProposalID)
		}
	}

	return nil
}

// getTopDelegates returns list of delegations grouped by dao
func (s *Service) getTopDelegates(_ context.Context, address string) (map[string][]Summary, error) {
	limitPerDao := 5
	list, err := s.repo.GetTopDelegatesByAddress(address, limitPerDao)
	if err != nil {
		return nil, fmt.Errorf("repo.GetTopDelegatesByAddress: %w", err)
	}

	result := make(map[string][]Summary, len(list))
	for _, info := range list {
		if _, ok := result[info.DaoID]; !ok {
			result[info.DaoID] = make([]Summary, 0, limitPerDao)
		}

		result[info.DaoID] = append(result[info.DaoID], info)
	}

	return result, nil
}

// getTopDelegators returns list of first 5 delegators grouped by dao
func (s *Service) getTopDelegators(_ context.Context, address string) (map[string][]Summary, error) {
	limitPerDao := 5
	list, err := s.repo.GetTopDelegatorsByAddress(address, limitPerDao)
	if err != nil {
		return nil, fmt.Errorf("repo.GetTopDelegatorsByAddress: %w", err)
	}

	result := make(map[string][]Summary, len(list))
	for _, info := range list {
		if _, ok := result[info.DaoID]; !ok {
			result[info.DaoID] = make([]Summary, 0, limitPerDao)
		}

		result[info.DaoID] = append(result[info.DaoID], info)
	}

	return result, nil
}

// getDelegatorsCnt returns count of delegators based on address
func (s *Service) getDelegatorsCnt(_ context.Context, address string) (int32, error) {
	cnt, err := s.repo.GetCnt(DelegateFilter{Address: address})
	if err != nil {
		return 0, fmt.Errorf("repo.GetByFilters: %w", err)
	}

	return int32(cnt), nil
}

// getDelegatesCnt returns count of delegations based on address
func (s *Service) getDelegatesCnt(_ context.Context, address string) (int32, error) {
	cnt, err := s.repo.GetCnt(DelegatorFilter{Address: address})
	if err != nil {
		return 0, fmt.Errorf("repo.GetByFilters: %w", err)
	}

	return int32(cnt), nil
}

func (s *Service) GetByFilters(filters ...Filter) ([]Summary, error) {
	return s.repo.GetByFilters(filters...)
}

func (s *Service) GetCntByFilters(filters ...Filter) (int64, error) {
	return s.repo.GetCnt(filters...)
}
