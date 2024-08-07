package delegate

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	protoany "github.com/golang/protobuf/ptypes/any"
	"github.com/google/uuid"
	"github.com/goverland-labs/goverland-datasource-snapshot/protocol/delegatepb"

	"github.com/goverland-labs/goverland-core-storage/internal/dao"
	"github.com/goverland-labs/goverland-core-storage/internal/ensresolver"
)

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

	queryAddresses, err := s.resolveQueryAccounts(request.QueryAccounts)
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
			Address:                  d.GetAddress(),
			ENSName:                  ensNames[d.GetAddress()],
			DelegatorCount:           d.GetDelegatorCount(),
			PercentOfDelegators:      d.GetPercentOfDelegators(),
			VotingPower:              d.GetVotingPower(),
			PercentOfVotingPower:     d.GetPercentOfVotingPower(),
			About:                    d.GetAbout(),
			Statement:                d.GetStatement(),
			UserDelegatedVotingPower: d.GetUserDelegatedVotingPower(),
			VotesCount:               0,
			ProposalsCount:           0,
			CreateProposalsCount:     0,
		})
	}

	return &GetDelegatesResponse{
		Delegates: delegates,
	}, nil
}

func (s *Service) resolveQueryAccounts(accs []string) ([]string, error) {
	var addresses []string
	var names []string
	for _, a := range accs {
		if common.IsHexAddress(a) {
			addresses = append(addresses)
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
