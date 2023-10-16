package vote

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
	groupName                = "vote"
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

func (c *Consumer) handler() pevents.VotesHandler {
	return func(payload pevents.VotesPayload) error {
		var err error
		defer func(start time.Time) {
			metricHandleHistogram.
				WithLabelValues("handle_votes", metrics.ErrLabelValue(err)).
				Observe(time.Since(start).Seconds())
		}(time.Now())

		err = c.service.HandleVotes(context.TODO(), convertToInternal(payload))
		if err != nil {
			log.Error().Err(err).Msg("process votes")

			return err
		}

		log.Debug().Msgf("vote was processed: %d", len(payload))

		return nil
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	group := config.GenerateGroupName(groupName)
	cc, err := client.NewConsumer(ctx, c.conn, group, pevents.SubjectVoteCreated, c.handler(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, pevents.SubjectVoteCreated, err)
	}

	c.consumers = append(c.consumers, cc)

	log.Info().Msg("vote consumers is started")

	// todo: handle correct stopping the consumer by context
	<-ctx.Done()
	return c.stop()
}

func (c *Consumer) stop() error {
	for _, cs := range c.consumers {
		if err := cs.Close(); err != nil {
			log.Error().Err(err).Msg("cant close vote consumer")
		}
	}

	return nil
}
