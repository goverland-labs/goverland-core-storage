package proposal

import (
	"context"
	"fmt"
	"time"

	pevents "github.com/goverland-labs/platform-events/events/aggregator"
	client "github.com/goverland-labs/platform-events/pkg/natsclient"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/core-storage/internal/config"
	"github.com/goverland-labs/core-storage/internal/metrics"
)

const (
	groupName                = "proposal"
	maxPendingAckPerConsumer = 1000
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

func (c *Consumer) handler() pevents.ProposalHandler {
	return func(payload pevents.ProposalPayload) error {
		var err error
		defer func(start time.Time) {
			metricHandleHistogram.
				WithLabelValues("handle_proposal", metrics.ErrLabelValue(err)).
				Observe(time.Since(start).Seconds())
		}(time.Now())

		err = c.service.HandleProposal(context.TODO(), convertToProposal(payload))
		if err != nil {
			log.Error().Err(err).Msg("process proposal")
		}

		log.Debug().Msgf("proposal was processed: %s", payload.ID)

		return err
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	group := config.GenerateGroupName(groupName)
	cc, err := client.NewConsumer(ctx, c.conn, group, pevents.SubjectProposalCreated, c.handler(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, pevents.SubjectProposalCreated, err)
	}
	cu, err := client.NewConsumer(ctx, c.conn, group, pevents.SubjectProposalUpdated, c.handler(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, pevents.SubjectProposalUpdated, err)
	}

	c.consumers = append(c.consumers, cc, cu)

	log.Info().Msg("proposal consumers is started")

	// todo: handle correct stopping the consumer by context
	<-ctx.Done()
	return c.stop()
}

func (c *Consumer) stop() error {
	for _, cs := range c.consumers {
		if err := cs.Close(); err != nil {
			log.Error().Err(err).Msg("cant close proposal consumer")
		}
	}

	return nil
}
