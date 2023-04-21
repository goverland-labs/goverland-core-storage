package dao

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
	groupName = "dao"
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

func (c *Consumer) handleCreate() pevents.DaoHandler {
	return func(payload pevents.DaoPayload) error {
		var err error
		defer func(start time.Time) {
			metricHandleHistogram.
				WithLabelValues("dao_create", metrics.ErrLabelValue(err)).
				Observe(time.Since(start).Seconds())
		}(time.Now())

		err = c.service.HandleDao(context.TODO(), convertToDao(payload))
		if err != nil {
			log.Error().Err(err).Msg("process dao")
		}

		return err
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	cs, err := client.NewConsumer(ctx, c.conn, groupName, pevents.SubjectDaoCreated, c.handleCreate())
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", groupName, pevents.SubjectDaoCreated, err)
	}

	c.consumers = append(c.consumers, cs)

	log.Info().Msg("dao consumers is started")

	// todo: handle correct stopping the consumer by context
	<-ctx.Done()
	return c.stop()
}

func (c *Consumer) stop() error {
	for _, cs := range c.consumers {
		if err := cs.Close(); err != nil {
			log.Error().Err(err).Msg("cant close dao consumer")
		}
	}

	return nil
}
