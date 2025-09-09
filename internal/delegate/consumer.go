package delegate

import (
	"context"
	"fmt"
	"time"

	events "github.com/goverland-labs/goverland-platform-events/events/aggregator"
	client "github.com/goverland-labs/goverland-platform-events/pkg/natsclient"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-core-storage/internal/config"
)

const (
	groupName                = "delegates"
	maxPendingAckPerConsumer = 10
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
		if err := c.service.handleDelegate(context.TODO(), convertToInternal(payload)); err != nil {
			log.Error().Err(err).Msg("delegates: process delegates info")

			return fmt.Errorf("delegates: process delegates info: %w", err)
		}

		log.Debug().Msgf("delegates: event was processed: %d %s", payload.BlockNumber, payload.ChainID)

		return nil
	}
}

func (c *Consumer) handleERC20Delegates() events.ERC20DelegateHandler {
	return func(payload events.ERC20DelegatePayload) error {
		if err := c.service.handleDelegate(context.TODO(), convertERC20ToInternal(payload)); err != nil {
			log.Error().Err(err).Msg("erc20 delegates: process delegates info")

			return fmt.Errorf("erc20 delegates: process delegates info: %w", err)
		}

		log.Debug().Msgf("erc20 delegates: event was processed: %d %s", payload.BlockNumber, payload.Network)

		return nil
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
	vfc, err := client.NewConsumer(ctx, c.conn, group, events.SubjectProposalVotesFetched, c.handleVotesFetched(), client.WithAckWait(time.Minute*5), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, events.SubjectProposalVotesFetched, err)
	}
	erc20c, err := client.NewConsumer(ctx, c.conn, group, events.SubjectDelegateERC20, c.handleERC20Delegates(), client.WithAckWait(time.Minute*5), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, events.SubjectDelegateERC20, err)
	}

	c.consumers = append(c.consumers, de, pr, vc, vfc, erc20c)

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
