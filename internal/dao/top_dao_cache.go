package dao

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	topDAOCacheReloadDelay = 5 * time.Minute
	topDaoCategoryLimit    = 20
)

type TopDAOCache struct {
	repo DataProvider

	cacheLock sync.RWMutex
	cache     map[string]topList
}

func NewTopDAOCache(repo DataProvider) *TopDAOCache {
	return &TopDAOCache{
		repo:  repo,
		cache: make(map[string]topList),
	}
}

func (w *TopDAOCache) Start(ctx context.Context) error {
	for {
		err := w.reload()
		if err != nil {
			log.Error().Err(err).Msg("failed to reload top DAO cache")
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(topDAOCacheReloadDelay):
		}
	}
}

func (w *TopDAOCache) GetTopList(limit uint) map[string]topList {
	w.cacheLock.RLock()
	defer w.cacheLock.RUnlock()

	return makeCopy(w.cache, limit)
}

func (w *TopDAOCache) reload() error {
	categories, err := w.repo.GetCategories()
	if err != nil {
		return fmt.Errorf("get categories: %w", err)
	}

	list := make(map[string]topList)
	for _, category := range categories {
		filters := []Filter{
			CategoryFilter{Category: category},
			PageFilter{Limit: topDaoCategoryLimit, Offset: 0},
			OrderByPopularityIndexFilter{},
		}

		data, err := w.repo.GetByFilters(filters, true)
		if err != nil {
			return fmt.Errorf("get by category %s: %w", category, err)
		}

		list[category] = topList{
			List:  data.Daos,
			Total: data.TotalCount,
		}
	}

	w.cacheLock.Lock()
	defer w.cacheLock.Unlock()

	w.cache = list
	return nil
}

func makeCopy(src map[string]topList, limit uint) map[string]topList {
	copied := map[string]topList{}
	for k, v := range src {
		newLen := min(limit, uint(len(v.List)))
		copied[k] = topList{
			List:  make([]Dao, newLen),
			Total: v.Total,
		}

		for i := range limit { // nolint:gosimple
			copied[k].List[i] = v.List[i]
		}
	}

	return copied
}
