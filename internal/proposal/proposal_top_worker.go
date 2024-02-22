package proposal

import (
	"context"
	"time"
)

const (
	topCheckDelay = 5 * time.Minute
)

type TopWorker struct {
	service *Service
}

func NewTopWorker(s *Service) *TopWorker {
	return &TopWorker{
		service: s,
	}
}

func (w *TopWorker) Start(ctx context.Context) error {
	for {
		w.service.prepareTop()

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(topCheckDelay):
		}
	}
}
