package events

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"gorm.io/gorm"
)

const (
	defaultCacheSize = 1000
)

// todo: add unit tests
type DataProvider interface {
	Create(RegisteredEvent) error
	Update(RegisteredEvent) error
	GetByTypeAndEvent(string, string, string) (*RegisteredEvent, error)
	GetLast(limit int) ([]*RegisteredEvent, error)
}

type Service struct {
	repo  DataProvider
	mu    sync.RWMutex
	cache map[string]struct{}
}

func NewService(r DataProvider) (*Service, error) {
	s := &Service{
		repo:  r,
		cache: make(map[string]struct{}, defaultCacheSize),
	}

	s.fillCache(defaultCacheSize)

	return s, nil
}

func (s *Service) fillCache(limit int) {
	items, err := s.repo.GetLast(limit)
	if err != nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, p := range items {
		key := prepareKey(p.TypeID, p.Type, p.Event)
		s.cache[key] = struct{}{}
	}
}

func (s *Service) EventExist(_ context.Context, id, t, event string) (bool, error) {
	key := prepareKey(id, t, event)

	s.mu.RLock()
	_, ok := s.cache[key]
	s.mu.RUnlock()

	if ok {
		return true, nil
	}

	_, err := s.repo.GetByTypeAndEvent(id, t, event)
	if err == nil {
		s.mu.Lock()
		s.cache[key] = struct{}{}
		s.mu.Unlock()

		return true, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return false, err
}

func (s *Service) RegisterEvent(_ context.Context, id, t, event string) error {
	err := s.repo.Create(RegisteredEvent{
		Type:   t,
		TypeID: id,
		Event:  event,
	})
	if err != nil {
		return err
	}

	key := prepareKey(id, t, event)
	s.mu.Lock()
	s.cache[key] = struct{}{}
	s.mu.Unlock()

	return nil
}

func prepareKey(id, t, event string) string {
	return fmt.Sprintf("%s_%s_%s", id, t, event)
}
