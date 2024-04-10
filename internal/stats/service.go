package stats

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/goverland-labs/goverland-core-storage/internal/dao"
	"github.com/goverland-labs/goverland-core-storage/internal/proposal"
)

type DaoProvider interface {
	GetCountByFilters(filters []dao.Filter) (int64, error)
}

type ProposalProvider interface {
	GetCountByFilters(filters []proposal.Filter) (int64, error)
}

type Service struct {
	cache Totals
	lock  sync.RWMutex

	dp DaoProvider
	pp ProposalProvider
}

func NewService(dp DaoProvider, pp ProposalProvider) *Service {
	return &Service{
		dp:    dp,
		pp:    pp,
		lock:  sync.RWMutex{},
		cache: Totals{},
	}
}

func (s *Service) GetTotals() Totals {
	s.lock.RLock()
	totals := s.cache
	s.lock.RUnlock()

	return totals
}

func (s *Service) refreshTotals() error {
	totals, err := s.calcTotals()
	if err != nil {
		return fmt.Errorf("calc totals: %w", err)
	}

	s.lock.Lock()
	s.cache = totals
	s.lock.Unlock()

	return nil
}

func (s *Service) calcTotals() (Totals, error) {
	var (
		daoTotal, daoVerified, prTotal int64
	)

	group, _ := errgroup.WithContext(context.TODO())
	group.Go(func() error {
		cnt, err := s.dp.GetCountByFilters(nil)
		if err != nil {
			return fmt.Errorf("s.dp.GetCountByFilters: %w", err)
		}

		daoTotal = cnt

		return nil
	})
	group.Go(func() error {
		cnt, err := s.dp.GetCountByFilters([]dao.Filter{
			dao.VerifiedFilter{},
		})
		if err != nil {
			return fmt.Errorf("s.dp.GetCountByFilters: %w", err)
		}

		daoVerified = cnt

		return nil
	})
	group.Go(func() error {
		cnt, err := s.pp.GetCountByFilters([]proposal.Filter{
			proposal.SkipCanceled{},
			proposal.SkipSpamFilter{},
		})
		if err != nil {
			return fmt.Errorf("s.dp.GetCountByFilters: %w", err)
		}

		prTotal = cnt

		return nil
	})

	if err := group.Wait(); err != nil {
		return Totals{}, fmt.Errorf("get counters: %w", err)
	}

	return Totals{
		Dao: Dao{
			Total:         daoTotal,
			TotalVerified: daoVerified,
		},
		Proposals: Proposals{
			Total: prTotal,
		},
	}, nil
}
