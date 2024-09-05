package delegates

import (
	"context"
	"fmt"

	events "github.com/goverland-labs/goverland-platform-events/events/core"
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

func (c *Consumer) handler() events.DelegateHandler {
	return func(payload events.DelegatePayload) error {
		if err := c.service.HandleDelegate(context.TODO(), convertToInternal(payload)); err != nil {
			log.Error().Err(err).Msg("process delegates info")

			return fmt.Errorf("process delegates info: %w", err)
		}

		log.Debug().Msgf("event was processed: %d %s", payload.BlockNumber, payload.ChainID)

		return nil
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	group := config.GenerateGroupName(groupName)
	de, err := client.NewConsumer(ctx, c.conn, group, events.SubjectDelegateUpsert, c.handler(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, events.SubjectDelegateUpsert, err)
	}
	c.consumers = append(c.consumers, de)

	log.Info().Msg("delegates consumers is started")

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
