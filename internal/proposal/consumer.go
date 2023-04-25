package proposal

import (
	"context"
	"fmt"
	"time"

	pevents "github.com/goverland-labs/platform-events/events/aggregator"
	client "github.com/goverland-labs/platform-events/pkg/natsclient"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/core-storage/internal/metrics"
)

const (
	groupName = "proposal"
)

type Consumer struct {
	conn      *nats.Conn
	service   *Service
	consumers []*client.Consumer
}

func NewConsumer(nc *nats.Conn, s *Service) (*Consumer, error) {
	c := &Consumer{
		conn:      nc,
		service:   s,
		consumers: make([]*client.Consumer, 0),
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
	cc, err := client.NewConsumer(ctx, c.conn, groupName, pevents.SubjectProposalCreated, c.handler())
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", groupName, pevents.SubjectProposalCreated, err)
	}
	cu, err := client.NewConsumer(ctx, c.conn, groupName, pevents.SubjectProposalUpdated, c.handler())
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", groupName, pevents.SubjectProposalUpdated, err)
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
