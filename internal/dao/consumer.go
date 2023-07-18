package dao

import (
	"context"
	"fmt"
	"time"

	pevents "github.com/goverland-labs/platform-events/events/aggregator"
	"github.com/goverland-labs/platform-events/events/core"
	client "github.com/goverland-labs/platform-events/pkg/natsclient"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/core-storage/internal/config"
	"github.com/goverland-labs/core-storage/internal/metrics"
)

const (
	groupName                = "dao"
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

func (c *Consumer) handler() pevents.DaoHandler {
	return func(payload pevents.DaoPayload) error {
		var err error
		defer func(start time.Time) {
			metricHandleHistogram.
				WithLabelValues("handle_dao", metrics.ErrLabelValue(err)).
				Observe(time.Since(start).Seconds())
		}(time.Now())

		err = c.service.HandleDao(context.TODO(), convertToDao(payload))
		if err != nil {
			log.Error().Err(err).Msg("process dao")
		}

		log.Debug().Msgf("dao was processed: %s", payload.ID)

		return err
	}
}

func (c *Consumer) activitySinceHandler() pevents.DaoHandler {
	return func(payload pevents.DaoPayload) error {
		var err error
		defer func(start time.Time) {
			metricHandleHistogram.
				WithLabelValues("handle_activity_since", metrics.ErrLabelValue(err)).
				Observe(time.Since(start).Seconds())
		}(time.Now())

		daoID, err := c.service.GetIDByOriginalID(payload.ID)
		if err != nil {
			log.Error().Err(err).Str("original_id", payload.ID).Msg("unable to get internal dao id by original")

			return err
		}

		updated, err := c.service.HandleActivitySince(context.TODO(), daoID)
		if err != nil {
			log.Error().Err(err).Msg("process dao activity since")

			return err
		}

		if updated == nil {
			return nil
		}

		if err = c.service.events.PublishJSON(context.TODO(), core.SubjectDaoUpdated, convertToCoreEvent(*updated)); err != nil {
			log.Error().Err(err).Msgf("publish dao event #%s", updated.ID)
		}

		log.Debug().Msgf("dao activity since was processed: %s", payload.ID)

		return err
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	group := config.GenerateGroupName(groupName)
	cc, err := client.NewConsumer(ctx, c.conn, group, pevents.SubjectDaoCreated, c.handler(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, pevents.SubjectDaoCreated, err)
	}
	cu, err := client.NewConsumer(ctx, c.conn, group, pevents.SubjectDaoUpdated, c.handler(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, pevents.SubjectDaoUpdated, err)
	}
	cac, err := client.NewConsumer(ctx, c.conn, group, core.SubjectCheckActivitySince, c.activitySinceHandler(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, core.SubjectCheckActivitySince, err)
	}

	c.consumers = append(c.consumers, cc, cu, cac)

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
