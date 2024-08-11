package delegate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	protoany "github.com/golang/protobuf/ptypes/any"
	"github.com/google/uuid"
	"github.com/goverland-labs/goverland-datasource-snapshot/protocol/delegatepb"

	"github.com/goverland-labs/goverland-core-storage/internal/dao"
	"github.com/goverland-labs/goverland-core-storage/internal/ensresolver"
)

var errNoResolved = errors.New("no addresses resolved")

type DaoProvider interface {
	GetByID(id uuid.UUID) (*dao.Dao, error)
}

type EnsResolver interface {
	GetByNames(names []string) ([]ensresolver.EnsName, error)
	GetByAddresses(addresses []string) ([]ensresolver.EnsName, error)
}

type Service struct {
	delegateClient delegatepb.DelegateClient
	daoProvider    DaoProvider
	ensResolver    EnsResolver
}

func NewService(dc delegatepb.DelegateClient, daoProvider DaoProvider, ensResolver EnsResolver) *Service {
	return &Service{
		delegateClient: dc,
		daoProvider:    daoProvider,
		ensResolver:    ensResolver,
	}
}

func (s *Service) GetDelegates(ctx context.Context, request GetDelegatesRequest) (*GetDelegatesResponse, error) {
	daoEntity, err := s.daoProvider.GetByID(request.DaoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dao: %w", err)
	}

	delegationStrategyJson, err := s.getDelegationStrategy(daoEntity, err)
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
	}, nil
}

func (s *Service) GetDelegateProfile(ctx context.Context, request GetDelegateProfileRequest) (GetDelegateProfileResponse, error) {
	daoEntity, err := s.daoProvider.GetByID(request.DaoID)
	if err != nil {
		return GetDelegateProfileResponse{}, fmt.Errorf("failed to get dao: %w", err)
	}

	delegationStrategyJson, err := s.getDelegationStrategy(daoEntity, err)
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

	return GetDelegateProfileResponse{
		Address:              resp.GetAddress(),
		VotingPower:          resp.GetVotingPower(),
		IncomingPower:        resp.GetIncomingPower(),
		OutgoingPower:        resp.GetOutgoingPower(),
		PercentOfVotingPower: resp.GetPercentOfVotingPower(),
		PercentOfDelegators:  resp.GetPercentOfDelegators(),
		Delegates:            delegates,
	}, nil
}

func (s *Service) getDelegationStrategy(daoEntity *dao.Dao, err error) ([]byte, error) {
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
		res[n.Address] = n.Name
	}

	return res, nil
}
