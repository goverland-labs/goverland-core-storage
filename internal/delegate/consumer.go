package delegate

import (
	"context"
	"fmt"
	"strings"

	events "github.com/goverland-labs/goverland-platform-events/events/aggregator"
	client "github.com/goverland-labs/goverland-platform-events/pkg/natsclient"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/goverland-labs/goverland-core-storage/internal/config"
)

const (
	groupName                = "delegates"
	maxPendingAckPerConsumer = 10
)

const (
	nullAddress = "0x0000000000000000000000000000000000000000"
)

type closable interface {
	Close() error
}

type Consumer struct {
	conn      *nats.Conn
	service   *Service
	consumers []closable
}

func NewConsumer(nc *nats.Conn, s *Service) (*Consumer, error) {
	c := &Consumer{
		conn:      nc,
		service:   s,
		consumers: make([]closable, 0),
	}

	return c, nil
}

func (c *Consumer) handleDelegates() events.DelegateHandler {
	return func(payload events.DelegatePayload) error {
		if err := c.service.handleSplitDelegation(context.TODO(), convertToInternal(payload)); err != nil {
			log.Error().Err(err).Msg("delegates: process delegates info")

			return fmt.Errorf("delegates: process delegates info: %w", err)
		}

		log.Debug().Msgf("delegates: event was processed: %d %s", payload.BlockNumber, payload.ChainID)

		return nil
	}
}

func (c *Consumer) handleERC20Delegates() events.ERC20DelegationHandler {
	return func(payload events.ERC20DelegationPayload) error {
		erc20Info, ok := tokenErc20Set[strings.ToLower(payload.Token)]
		if !ok {
			log.Warn().Msgf("erc20 mapping not found for token %s", payload.Token)

			return nil
		}

		event := ERC20Delegation{
			Token:            strings.ToLower(payload.Token),
			DelegatorAddress: strings.ToLower(payload.Delegator),
			AddressFrom:      strings.ToLower(payload.AddressFrom),
			AddressTo:        strings.ToLower(payload.AddressTo),
			ChainID:          erc20Info.ChainID,
			BlockNumber:      int(payload.BlockNumber),
			BlockTimestamp:   int(payload.BlockTimestamp),
			LogIndex:         int(payload.LogIndex),
		}

		processor := func(ctx context.Context, tx *gorm.DB) error {
			err := c.service.handleERC20Delegation(ctx, event, tx)
			if err != nil {
				return fmt.Errorf("c.service.handleERC20Delegation: %w", err)
			}

			if event.AddressTo == nullAddress {
				if err = c.service.UpdateERC20Totals(tx, ERC20TotalChanges{
					Token:           event.Token,
					ChainID:         event.ChainID,
					VPDelta:         "0",
					DelegatorsDelta: -1,
				}); err != nil {
					return fmt.Errorf("c.service.UpdateERC20Totals: %w", err)
				}
			} else {
				increaseCnt := 1
				if err = c.service.UpdateERC20Delegate(tx, ERC20DelegateUpdate{
					Address:  event.AddressTo,
					Token:    event.Token,
					ChainID:  event.ChainID,
					CntDelta: &increaseCnt,
				}); err != nil {
					return fmt.Errorf("c.service.UpdateERC20Delegate: increase: %w", err)
				}
			}

			if event.AddressFrom == nullAddress {
				if err = c.service.UpdateERC20Totals(tx, ERC20TotalChanges{
					Token:           event.Token,
					ChainID:         event.ChainID,
					VPDelta:         "0",
					DelegatorsDelta: 1,
				}); err != nil {
					return fmt.Errorf("c.service.UpdateERC20Totals: %w", err)
				}

				return nil
			}

			decreaseCnt := -1
			if err = c.service.UpdateERC20Delegate(tx, ERC20DelegateUpdate{
				Address:  event.AddressFrom,
				Token:    event.Token,
				ChainID:  event.ChainID,
				CntDelta: &decreaseCnt,
			}); err != nil {
				return fmt.Errorf("c.service.UpdateERC20Delegate: decrease: %w", err)
			}

			return nil
		}

		return c.service.processErc20Event(context.TODO(), event, processor)
	}
}

func (c *Consumer) handleERC20VPChanges() events.ERC20VPChangesHandler {
	return func(payload events.ERC20VPChangesPayload) error {
		erc20Info, ok := tokenErc20Set[strings.ToLower(payload.Token)]
		if !ok {
			log.Warn().Msgf("erc20 mapping not found for token %s", payload.Token)

			return nil
		}

		event := ERC20VPChanges{
			Token:          strings.ToLower(payload.Token),
			ChainID:        erc20Info.ChainID,
			Address:        strings.ToLower(payload.Address),
			BlockNumber:    int(payload.BlockNumber),
			BlockTimestamp: int(payload.BlockTimestamp),
			LogIndex:       int(payload.LogIndex),
			VP:             payload.VotingPower,
			Delta:          payload.DeltaPower,
		}

		processor := func(ctx context.Context, tx *gorm.DB) error {
			if err := c.service.UpdateERC20Delegate(tx, ERC20DelegateUpdate{
				Token:   event.Token,
				ChainID: event.ChainID,
				Address: event.Address,
				VPUpdate: &VPUpdate{
					Value:       event.VP,
					BlockNumber: event.BlockNumber,
					LogIndex:    event.LogIndex,
				},
			}); err != nil {
				return fmt.Errorf("c.service.UpdateERC20Delegate: vp_changes: %w", err)
			}

			if err := c.service.UpdateERC20Totals(tx, ERC20TotalChanges{
				Token:           event.Token,
				ChainID:         event.ChainID,
				VPDelta:         event.Delta,
				DelegatorsDelta: 0,
			}); err != nil {
				return fmt.Errorf("c.service.UpdateERC20Totals: %w", err)
			}

			return nil
		}

		return c.service.processErc20Event(context.TODO(), event, processor)
	}
}

func (c *Consumer) handleERC20Transfers() events.ERC20TransfersHandler {
	return func(payload events.ERC20TransferPayload) error {
		erc20Info, ok := tokenErc20Set[strings.ToLower(payload.Token)]
		if !ok {
			log.Warn().Msgf("erc20 mapping not found for token %s", payload.Token)

			return nil
		}

		event := ERC20Transfer{
			Token:          strings.ToLower(payload.Token),
			AddressFrom:    strings.ToLower(payload.AddressFrom),
			AddressTo:      strings.ToLower(payload.AddressTo),
			ChainID:        erc20Info.ChainID,
			BlockNumber:    int(payload.BlockNumber),
			BlockTimestamp: int(payload.BlockTimestamp),
			LogIndex:       int(payload.LogIndex),
			Amount:         payload.Amount,
		}

		processor := func(ctx context.Context, tx *gorm.DB) error {
			amount := event.Amount
			if err := c.service.UpsertERC20Balance(tx, event.Token, event.ChainID, event.AddressTo, amount); err != nil {
				return fmt.Errorf("c.service.UpsertERC20Balance: increase: %w", err)
			}

			negAmount := "-" + amount
			if err := c.service.UpsertERC20Balance(tx, event.Token, event.ChainID, event.AddressFrom, negAmount); err != nil {
				return fmt.Errorf("c.service.UpsertERC20Balance: decrease: %w", err)
			}

			return nil
		}

		return c.service.processErc20Event(context.TODO(), event, processor)
	}
}

func (c *Consumer) handleProposalCreated() events.ProposalHandler {
	return func(payload events.ProposalPayload) error {
		if err := c.service.handleProposalCreated(context.TODO(), convertEventToProposal(payload)); err != nil {
			log.Error().Err(err).Msg("delegates: process proposal create")

			return fmt.Errorf("delegates: process proposal create: %w", err)
		}

		log.Debug().Msgf("delegates: proposal create event was processed: %s %s %s", payload.ID, payload.DaoID, payload.Author)

		return nil
	}
}

func (c *Consumer) handleVotesCreated() events.VotesHandler {
	return func(payload events.VotesPayload) error {
		if err := c.service.handleVotesCreated(context.TODO(), convertEventToVoteDetails(payload)); err != nil {
			log.Error().Err(err).Msg("delegates: process votes created")

			return fmt.Errorf("delegates: process votes created: %w", err)
		}

		log.Debug().Msgf("delegates: votes created event pack was processed: %d", len(payload))

		return nil
	}
}

func (c *Consumer) handleVotesFetched() events.ProposalHandler {
	return func(payload events.ProposalPayload) error {
		if err := c.service.handleVotesFetched(context.TODO(), payload.ID); err != nil {
			log.Error().Err(err).Msg("delegates: process votes fetched")

			return fmt.Errorf("delegates: process votes fetched: %w", err)
		}

		log.Debug().Msgf("delegates: votes fetched event was processed: %s", payload.ID)

		return nil
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	group := config.GenerateGroupName(groupName)
	de, err := client.NewConsumer(ctx, c.conn, group, events.SubjectDelegateUpsert, c.handleDelegates(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, events.SubjectDelegateUpsert, err)
	}
	pr, err := client.NewConsumer(ctx, c.conn, group, events.SubjectProposalCreated, c.handleProposalCreated(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, events.SubjectProposalCreated, err)
	}
	vc, err := client.NewConsumer(ctx, c.conn, group, events.SubjectVoteCreated, c.handleVotesCreated(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, events.SubjectVoteCreated, err)
	}
	vfc, err := client.NewConsumer(ctx, c.conn, group, events.SubjectProposalVotesFetched, c.handleVotesFetched(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, events.SubjectProposalVotesFetched, err)
	}
	erc20del, err := client.NewConsumer(ctx, c.conn, group, events.SubjectERC20Delegations, c.handleERC20Delegates(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, events.SubjectERC20Delegations, err)
	}
	erc20vpc, err := client.NewConsumer(ctx, c.conn, group, events.SubjectERC20VPChanges, c.handleERC20VPChanges(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, events.SubjectERC20VPChanges, err)
	}
	erc20transfers, err := client.NewConsumer(ctx, c.conn, group, events.SubjectERC20Transfer, c.handleERC20Transfers(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, events.SubjectERC20Transfer, err)
	}

	c.consumers = append(c.consumers, de, pr, vc, vfc, erc20del, erc20vpc, erc20transfers)

	log.Info().Msg("delegates consumers are started")

	// todo: handle correct stopping the consumer by context
	<-ctx.Done()
	return c.stop()
}

func (c *Consumer) stop() error {
	for _, cs := range c.consumers {
		if err := cs.Close(); err != nil {
			log.Error().Err(err).Msg("cant close delegates consumer")
		}
	}

	return nil
}
