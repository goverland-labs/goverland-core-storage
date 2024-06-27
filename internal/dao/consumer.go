package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	pevents "github.com/goverland-labs/goverland-platform-events/events/aggregator"
	"github.com/goverland-labs/goverland-platform-events/events/core"
	coreevents "github.com/goverland-labs/goverland-platform-events/events/core"
	client "github.com/goverland-labs/goverland-platform-events/pkg/natsclient"

	"github.com/goverland-labs/goverland-core-storage/internal/config"
	"github.com/goverland-labs/goverland-core-storage/internal/metrics"
)

const (
	groupName                = "dao"
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

func (c *Consumer) uniqueVoters() coreevents.VotesHandler {
	return func(payload coreevents.VotesPayload) error {
		var err error
		defer func(start time.Time) {
			metricHandleHistogram.
				WithLabelValues("handle_unique_voters", metrics.ErrLabelValue(err)).
				Observe(time.Since(start).Seconds())
		}(time.Now())

		err = c.service.ProcessUniqueVoters(context.TODO(), convertVoteToInternal(payload))
		if err != nil {
			log.Error().Err(err).Msg("process dao unique voters")

			return err
		}

		log.Debug().Msg("dao unique voters was processed")

		return err
	}
}

func (c *Consumer) handleProposalCreated() pevents.ProposalHandler {
	return func(payload pevents.ProposalPayload) error {
		var err error
		defer func(start time.Time) {
			metricHandleHistogram.
				WithLabelValues("handle_proposal_created", metrics.ErrLabelValue(err)).
				Observe(time.Since(start).Seconds())
		}(time.Now())

		err = c.service.ProcessNewProposal(context.TODO(), payload.DaoID)
		if err != nil {
			log.Error().Err(err).Msg("process dao proposal created")

			return err
		}

		log.Debug().Msg("dao proposal created was processed")

		return err
	}
}

func (c *Consumer) handleProposalDeleted() pevents.ProposalHandler {
	return func(payload pevents.ProposalPayload) error {
		pr, err := c.service.proposals.GetByID(payload.ID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		if err != nil {
			log.Error().Err(err).Msg("process dao proposal deleted")

			return fmt.Errorf("proposals.GetByID: %s: %w", payload.ID, err)
		}

		if err := c.service.ProcessDeletedProposal(context.TODO(), pr.DaoID); err != nil {
			log.Error().Err(err).Msg("process dao proposal created")

			return err
		}

		log.Debug().Msg("dao consumer: proposal deleted event was processed")

		return nil
	}
}

func (c *Consumer) handleProposalUpdated() pevents.ProposalHandler {
	return func(payload pevents.ProposalPayload) error {
		var err error
		defer func(start time.Time) {
			metricHandleHistogram.
				WithLabelValues("handle_proposal_updated", metrics.ErrLabelValue(err)).
				Observe(time.Since(start).Seconds())
		}(time.Now())

		err = c.service.ProcessExistedProposal(context.TODO(), payload.DaoID)
		if err != nil {
			log.Error().Err(err).Msg("process dao proposal updated")

			return err
		}

		log.Debug().Msg("dao proposal updated was processed")

		return err
	}
}

func (c *Consumer) popularityIndexHandler() core.DaoHandler {
	return func(payload core.DaoPayload) error {
		var err error
		defer func(start time.Time) {
			metricHandleHistogram.
				WithLabelValues("handle_popularity_index_update", metrics.ErrLabelValue(err)).
				Observe(time.Since(start).Seconds())
		}(time.Now())

		err = c.service.ProcessPopularityIndexUpdate(context.TODO(), payload.ID, *payload.PopularityIndex)
		if err != nil {
			log.Error().Err(err).Msg("process popularity index updated")

			return err
		}

		log.Debug().Msg("dao popularity index updated was processed")

		return err
	}
}

func convertVoteToInternal(pl coreevents.VotesPayload) []UniqueVoter {
	res := make([]UniqueVoter, len(pl))
	for i := range pl {
		res[i] = UniqueVoter{
			DaoID: pl[i].DaoID,
			Voter: pl[i].Voter,
		}
	}

	return res
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
	vc, err := client.NewConsumer(ctx, c.conn, group, coreevents.SubjectVoteCreated, c.uniqueVoters(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, coreevents.SubjectVoteCreated, err)
	}
	pc, err := client.NewConsumer(ctx, c.conn, group, pevents.SubjectProposalCreated, c.handleProposalCreated(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, pevents.SubjectProposalCreated, err)
	}
	pu, err := client.NewConsumer(ctx, c.conn, group, pevents.SubjectProposalUpdated, c.handleProposalUpdated(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, pevents.SubjectProposalUpdated, err)
	}
	piu, err := client.NewConsumer(ctx, c.conn, group, core.SubjectPopularityIndexUpdated, c.popularityIndexHandler(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, core.SubjectPopularityIndexUpdated, err)
	}
	pd, err := client.NewConsumer(ctx, c.conn, group, pevents.SubjectProposalDeleted, c.handleProposalDeleted(), client.WithMaxAckPending(maxPendingAckPerConsumer))
	if err != nil {
		return fmt.Errorf("consume for %s/%s: %w", group, pevents.SubjectProposalDeleted, err)
	}

	c.consumers = append(
		c.consumers,
		cc,
		cu,
		cac,
		vc,
		pc,
		pu,
		piu,
		pd,
	)

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
