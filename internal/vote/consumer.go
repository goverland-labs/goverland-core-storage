package vote

import (
	"context"
	"fmt"
	"time"

	pevents "github.com/goverland-labs/goverland-platform-events/events/aggregator"
	coreevents "github.com/goverland-labs/goverland-platform-events/events/core"
	client "github.com/goverland-labs/goverland-platform-events/pkg/natsclient"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"

	"github.com/goverland-labs/goverland-core-storage/internal/config"
	"github.com/goverland-labs/goverland-core-storage/internal/metrics"
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

func (c *Consumer) handlerAddressResolved() coreevents.EnsNamesHandler {
	return func(payload coreevents.EnsNamesPayload) error {
		var err error
		defer func(start time.Time) {
			metricHandleHistogram.
				WithLabelValues("handle_address_resolved", metrics.ErrLabelValue(err)).
				Observe(time.Since(start).Seconds())
		}(time.Now())

		err = c.service.HandleResolvedAddresses(convertToResolvedAddresses(payload))
		if err != nil {
			log.Error().Err(err).Msg("process address resolved")
		}

		log.Debug().Msgf("proposal resolved addresses were processed")

		return err
	}
}

func convertToResolvedAddresses(list []coreevents.EnsNamePayload) []ResolvedAddress {
	res := make([]ResolvedAddress, 0, len(list))
	for i := range list {
		res = append(res, ResolvedAddress{
			Address: list[i].Address,
			Name:    list[i].Name,
		})
	}

	return res
}

func (c *Consumer) Start(ctx context.Context) error {
	group := config.GenerateGroupName(groupName)
	cc, err := client.NewConsumer(ctx, c.conn, group, pevents.SubjectVoteCreated, c.handler(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, pevents.SubjectVoteCreated, err)
	}
	cer, err := client.NewConsumer(ctx, c.conn, group, coreevents.SubjectEnsResolverResolved, c.handlerAddressResolved(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, coreevents.SubjectEnsResolverResolved, err)
	}

	c.consumers = append(c.consumers, cc, cer)

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
