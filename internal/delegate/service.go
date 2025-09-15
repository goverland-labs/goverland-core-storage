package delegate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
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
	"github.com/goverland-labs/goverland-core-storage/internal/proposal"
)

var (
	errNoResolved = errors.New("no addresses resolved")
)

const (
	unrecognizedStrategyName = "unrecognized"
)

type DaoProvider interface {
	GetByID(id uuid.UUID) (*dao.Dao, error)
	GetIDByOriginalID(string) (uuid.UUID, error)
	GetDaoByOriginalID(id string) (*dao.Dao, error)
}

type ProposalProvider interface {
	GetByID(id string) (*proposal.Proposal, error)
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
	delegateClient   delegatepb.DelegateClient
	daoProvider      DaoProvider
	proposalProvider ProposalProvider
	ensResolver      EnsResolver
	publisher        Publisher
	er               EventRegistered
	repo             *Repo

	mu            sync.RWMutex
	allowedDaoIDs []uuid.UUID
}

func NewService(repo *Repo, dc delegatepb.DelegateClient, daoProvider DaoProvider, prProvider ProposalProvider, ensResolver EnsResolver, ep Publisher, er EventRegistered) *Service {
	return &Service{
		delegateClient:   dc,
		daoProvider:      daoProvider,
		proposalProvider: prProvider,
		ensResolver:      ensResolver,
		publisher:        ep,
		er:               er,
		repo:             repo,
		allowedDaoIDs:    make([]uuid.UUID, 0),
	}
}

func (s *Service) UpdateAllowedDaos(ctx context.Context) error {
	for {
		if err := s.updateAllowedDaos(); err != nil {
			log.Error().Err(err).Msg("delegates lifetime check failed")
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(updateAllowedDaoTTL):
		}
	}
}

func (s *Service) updateAllowedDaos() error {
	allowedDaos, err := s.repo.AllowedDaos()
	if err != nil {
		return fmt.Errorf("s.daoProvider.GetByID: %w", err)
	}

	allowed := make([]uuid.UUID, 0, len(allowedDaos))
	for _, info := range allowedDaos {
		allowed = append(allowed, info.InternalID)
	}

	s.mu.Lock()
	s.allowedDaoIDs = allowed
	s.mu.Unlock()

	return nil
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
		respAddresses = append(respAddresses, strings.ToLower(d.GetAddress()))
	}
	ensNames, err := s.resolveAddressesName(respAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve addresses names: %w", err)
	}

	votesCnt, err := s.repo.GetVotesCnt(daoEntity.ID, respAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to get votes count: %w", err)
	}

	prCnt, err := s.repo.GetProposalsCnt(daoEntity.ID, respAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to get votes count: %w", err)
	}

	delegates := make([]Delegate, 0, len(resp.Delegates))
	for _, d := range resp.GetDelegates() {
		address := strings.ToLower(d.GetAddress())
		votes := votesCnt[address]
		proposals := prCnt[address]

		delegates = append(delegates, Delegate{
			Address:               address,
			ENSName:               ensNames[address],
			DelegatorCount:        d.GetDelegatorCount(),
			PercentOfDelegators:   d.GetPercentOfDelegators(),
			VotingPower:           d.GetVotingPower(),
			PercentOfVotingPower:  d.GetPercentOfVotingPower(),
			About:                 d.GetAbout(),
			Statement:             d.GetStatement(),
			VotesCount:            int32(votes),
			CreatedProposalsCount: int32(proposals),
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
	strategy := daoEntity.GetStrategyByName(dao.StrategyNameSplitDelegation)
	if strategy == nil {
		return nil, fmt.Errorf("delegation strategy not found for dao")
	}

	// TODO: avoid wrong naming, fix it in another way
	type marshalStrategy struct {
		Name    string                 `json:"name"`
		Network string                 `json:"network"`
		Params  map[string]interface{} `json:"params"`
	}

	ms := marshalStrategy{
		Name:    strategy.Name,
		Network: strategy.Network,
		Params:  strategy.Params,
	}

	converted, err := json.Marshal(ms)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal delegation strategy: %w", err)
	}

	return converted, nil
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
	logger := log.With().
		Str("source", "handle_delegates").
		Str("block_number", fmt.Sprintf("%d", hr.BlockNumber)).
		Str("address_from", hr.AddressFrom).
		Logger()

	if err := s.repo.CallInTx(func(tx *gorm.DB) error {
		if hr.OriginalSpaceID == "" {
			logger.Warn().Msg("skip processing block from empty original_space_id")

			return nil
		}

		// store to history
		if err := s.repo.CreateHistory(tx, hr); err != nil {
			return fmt.Errorf("repo.CreateHistory: %w", err)
		}

		// get space by provided original_space_id
		delegatedDao, err := s.daoProvider.GetDaoByOriginalID(hr.OriginalSpaceID)
		if err != nil {
			return fmt.Errorf("dp.GetIDByOriginalID: %w", err)
		}

		strategy := getDaoPrimaryStrategy(delegatedDao, hr.Source)
		if strategy == nil {
			logger.Warn().Msgf("no strategy found for delegated dao %s", delegatedDao.OriginalID)

			return nil
		}

		var chainID *string
		if strategy.Name != dao.StrategyNameSplitDelegation {
			chainID = &hr.ChainID
		}

		bts, err := s.repo.GetSummaryBlockTimestamp(tx, strings.ToLower(hr.AddressFrom), delegatedDao.ID.String(), chainID)
		if err != nil {
			return fmt.Errorf("s.repo.GetSummaryBlockTimestamp: %w", err)
		}

		// skip this block due to already processed
		if bts != 0 && bts >= hr.BlockTimestamp {
			logger.Warn().Msgf("delegates: skip processing block %d from %s due to invalid timestamp", hr.BlockNumber, hr.ChainID)

			return nil
		}

		if hr.Action == actionExpire {
			if err := s.repo.UpdateSummaryExpiration(tx, strings.ToLower(hr.AddressFrom), delegatedDao.ID.String(), hr.Delegations.Expiration, hr.BlockTimestamp); err != nil {
				return fmt.Errorf("UpdateSummaryExpiration: %w", err)
			}

			return nil
		}

		if err = s.repo.RemoveSummary(tx, strings.ToLower(hr.AddressFrom), delegatedDao.ID.String(), chainID); err != nil {
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
				DaoID:              delegatedDao.ID.String(),
				Weight:             info.Weight,
				LastBlockTimestamp: hr.BlockTimestamp,
				ExpiresAt:          int64(hr.Delegations.Expiration),
				Type:               strategy.Name,
				ChainID:            chainID,
				VotingPower:        hr.VotingPower,
			}); err != nil {
				return fmt.Errorf("createSummary [%s/%s/%s]: %w", hr.AddressFrom, info.Address, delegatedDao.ID.String(), err)
			}
		}

		for _, info := range hr.Delegations.Details {
			event := events.DelegatePayload{
				Initiator: strings.ToLower(hr.AddressFrom),
				Delegator: strings.ToLower(info.Address),
				DaoID:     delegatedDao.ID,
			}

			if err = s.publisher.PublishJSON(ctx, events.SubjectDelegateCreated, event); err != nil {
				logger.Warn().Err(err).Msgf("failed to publish delegate payload")
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

func getDaoPrimaryStrategy(space *dao.Dao, source string) *dao.Strategy {
	var strategy *dao.Strategy
	if strategy = space.GetStrategyByName(dao.StrategyNameSplitDelegation); strategy != nil {
		return strategy
	}

	switch source {
	case sourceSplitDelegation:
		strategy = space.GetStrategyByName(dao.StrategyNameDelegation)
	case sourceErc20Votes:
		strategy = space.GetStrategyByName(dao.StrategyNameErc20Votes)
	}

	if strategy != nil {
		return strategy
	}

	return &dao.Strategy{
		Name: unrecognizedStrategyName,
	}
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

func (s *Service) handleVotesFetched(ctx context.Context, prId string) error {
	pr, err := s.proposalProvider.GetByID(prId)
	if err != nil {
		return fmt.Errorf("proposalProvider.GetByID: %w", err)
	}

	s.mu.RLock()
	allowedDaos := make([]uuid.UUID, len(s.allowedDaoIDs))
	copy(allowedDaos, s.allowedDaoIDs)
	s.mu.RUnlock()

	if !slices.Contains(allowedDaos, pr.DaoID) {
		return nil
	}

	offset, limit := 0, 500
	for {
		delegates, err := s.repo.FindDelegates(pr.DaoID.String(), offset, limit)
		if err != nil {
			return fmt.Errorf("repo.FindDelegates: %w", err)
		}

		if len(delegates) == 0 {
			return nil
		}

		delegators := make(map[string]Summary, len(delegates))
		addresses := make([]string, 0, len(delegates))

		for _, info := range delegates {
			if info.SelfDelegation() || info.Expired() {
				continue
			}

			delegators[strings.ToLower(info.AddressTo)] = info
			addresses = append(addresses, strings.ToLower(info.AddressTo))
		}

		voters, err := s.repo.GetVotersByAddresses(pr.DaoID.String(), pr.ID, addresses)
		if err != nil {
			return fmt.Errorf("repo.GetVotersByAddresses: %w", err)
		}

		if len(voters) == len(addresses) {
			return nil
		}

		for _, address := range addresses {
			if slices.Contains(voters, address) {
				continue
			}

			info := delegators[address]
			info.ProposalID = prId
			s.registerEventOnce(ctx, info, events.SubjectDelegateVotingSkipVote)
		}

		if len(delegates) < limit {
			break
		}

		offset += limit
	}

	return nil
}

// getTopDelegates returns list of delegations grouped by dao
func (s *Service) getTopDelegates(_ context.Context, address string) (map[string][]Summary, error) {
	limitPerDao := 100
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
