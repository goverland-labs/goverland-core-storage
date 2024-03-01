package ensresolver

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/goverland-labs/goverland-helpers-ens-resolver/protocol/enspb"
	coreevents "github.com/goverland-labs/goverland-platform-events/events/core"
	"github.com/rs/zerolog/log"
)

const (
	maxBatchCount = 30
	deadline      = time.Minute
)

type Publisher interface {
	PublishJSON(ctx context.Context, subject string, obj any) error
}

type Service struct {
	client     enspb.EnsClient
	repo       *Repo
	mu         sync.Mutex
	inProgress atomic.Int64

	publisher Publisher

	// store addresses
	queue []string
}

func NewService(repo *Repo, cl enspb.EnsClient, pl Publisher) (*Service, error) {
	return &Service{
		repo:      repo,
		client:    cl,
		publisher: pl,
		queue:     make([]string, 0, maxBatchCount),
	}, nil
}

func (s *Service) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			s.resolve()

			for s.inProgress.Load() != 0 {
				<-time.After(time.Second)
			}

			return ctx.Err()
		case <-time.After(deadline):
			s.resolve()
		}
	}
}

func (s *Service) resolve() {
	s.mu.Lock()
	requests := make([]string, len(s.queue))
	copy(requests, s.queue)
	s.queue = make([]string, 0, maxBatchCount)
	s.mu.Unlock()

	if len(requests) == 0 {
		return
	}

	s.inProgress.Add(1)
	defer s.inProgress.Add(-1)

	for _, chunk := range chunkSlice(requests, 50) {
		s.processAddresses(chunk)
	}
}

func (s *Service) processAddresses(list []string) {
	if len(list) == 0 {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	res, err := s.client.ResolveDomains(ctx, &enspb.ResolveDomainsRequest{
		Addresses: list,
	})
	if err != nil {
		// todo: add to queue?
		log.Error().Err(err).Msg("resolve domains")

		return
	}

	result := make([]EnsName, 0, len(res.Addresses))
	for _, address := range res.Addresses {
		result = append(result, EnsName{
			Address: address.Address,
			Name:    address.EnsName,
		})
	}

	if err := s.repo.BatchCreate(result); err != nil {
		log.Error().Err(err).Msg("ens batch create in db")

		return
	}

	if err := s.publisher.PublishJSON(ctx, coreevents.SubjectEnsResolverResolved, convertToCoreEvent(result)); err != nil {
		log.Error().Err(err).Msgf("publish ens names event")
	}

	log.Info().Msgf("processed %d items", len(list))
}

func convertToCoreEvent(list []EnsName) coreevents.EnsNamesPayload {
	res := make(coreevents.EnsNamesPayload, 0, len(list))
	for i := range list {
		res = append(res, coreevents.EnsNamePayload{
			Address: list[i].Address,
			Name:    list[i].Name,
		})
	}

	return res
}

// AddRequests add requests to resolve ens names. The results will be putted to the queue
func (s *Service) AddRequests(list []string) {
	s.mu.Lock()
	s.queue = append(s.queue, list...)
	s.mu.Unlock()

	if len(list) < maxBatchCount {
		return
	}

	go s.resolve()
}

func chunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for {
		if len(slice) == 0 {
			break
		}

		if len(slice) < chunkSize {
			chunkSize = len(slice)
		}

		chunks = append(chunks, slice[0:chunkSize])
		slice = slice[chunkSize:]
	}

	return chunks
}
