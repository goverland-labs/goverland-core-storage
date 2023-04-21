package dao

import (
	"context"
	"fmt"

	pevents "github.com/goverland-labs/platform-events/events/aggregator"
	client "github.com/goverland-labs/platform-events/pkg/natsclient"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

const (
	groupName = "dao"
)

type Consumer struct {
	ctx       context.Context
	conn      *nats.Conn
	service   *Service
	consumers []*client.Consumer
}

func NewConsumer(ctx context.Context, nc *nats.Conn, s *Service) (*Consumer, error) {
	c := &Consumer{
		ctx:       ctx,
		conn:      nc,
		service:   s,
		consumers: make([]*client.Consumer, 0),
	}

	return c, nil
}

// fixme: measure execution time
// fixme: measure processed events: [subject?, cnt]
// fixme: add counting err
func (c *Consumer) handleCreate() pevents.DaoHandler {
	return func(payload pevents.DaoPayload) error {
		err := c.service.HandleDao(c.ctx, convertToDao(payload))
		if err != nil {
			log.Error().Err(err).Msg("process dao")
		}

		return err
	}
}

func (c *Consumer) Start() error {
	cs, err := client.NewConsumer(c.ctx, c.conn, groupName, pevents.SubjectDaoCreated, c.handleCreate())
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", groupName, pevents.SubjectDaoCreated, err)
	}

	c.consumers = append(c.consumers, cs)

	log.Info().Msg("dao consumers is started")

	select {
	case <-c.ctx.Done():
		return nil
	}
}

func (c *Consumer) Stop() error {
	for _, cs := range c.consumers {
		if err := cs.Close(); err != nil {
			log.Error().Err(err).Msg("cant close dao consumer")
		}
	}

	return nil
}
