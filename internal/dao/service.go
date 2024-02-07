package dao

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	coreevents "github.com/goverland-labs/platform-events/events/core"

	"github.com/goverland-labs/core-storage/internal/proposal"
)

//go:generate mockgen -destination=mocks_test.go -package=dao . DataProvider,Publisher,DaoIDProvider

const (
	newDaoCategoryName     = "new_daos"
	popularDaoCategoryName = "popular_daos"
)

var (
	systemCategories = []string{
		newDaoCategoryName,
		popularDaoCategoryName,
	}
)

type Publisher interface {
	PublishJSON(ctx context.Context, subject string, obj any) error
}

type DataProvider interface {
	Create(dao Dao) error
	Update(dao Dao) error
	GetByID(id uuid.UUID) (*Dao, error)
	UpdateProposalCnt(id uuid.UUID) error
	UpdateActiveVotes(id uuid.UUID) error
	UpdateActiveVotesAll() error
	GetByFilters(filters []Filter, count bool) (DaoList, error)
	GetCategories() ([]string, error)
}

type DaoIDProvider interface {
	GetOrCreate(originID string) (uuid.UUID, error)
	GetAll() ([]DaoID, error)
}

type UniqueVoterProvider interface {
	BatchCreate([]UniqueVoter) error
	UpdateVotersCount() error
}

type ProposalProvider interface {
	GetEarliestByDaoID(uuid.UUID) (*proposal.Proposal, error)
}

type Service struct {
	daoIds map[string]uuid.UUID
	daoMu  sync.RWMutex

	repo       DataProvider
	uniqueRepo UniqueVoterProvider
	events     Publisher
	idProvider DaoIDProvider
	proposals  ProposalProvider
}

func NewService(r DataProvider, ur UniqueVoterProvider, ip DaoIDProvider, p Publisher, pp ProposalProvider) (*Service, error) {
	return &Service{
		repo:       r,
		uniqueRepo: ur,
		events:     p,
		idProvider: ip,
		proposals:  pp,
		daoIds:     make(map[string]uuid.UUID),
		daoMu:      sync.RWMutex{},
	}, nil
}

func (s *Service) PrefillDaoIDs() error {
	list, err := s.idProvider.GetAll()
	if err != nil {
		return err
	}

	s.daoMu.Lock()
	s.daoIds = make(map[string]uuid.UUID)
	for i := range list {
		item := list[i]
		s.daoIds[item.OriginalID] = item.InternalID
	}
	s.daoMu.Unlock()

	return nil
}

func (s *Service) HandleDao(ctx context.Context, dao Dao) error {
	id, err := s.GetIDByOriginalID(dao.OriginalID)
	if err != nil {
		return fmt.Errorf("getting/generating dao id: %w", err)
	}
	dao.ID = id

	existed, err := s.repo.GetByID(id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("handle: %w", err)
	}

	if existed == nil {
		return s.processNew(ctx, dao)
	}

	return s.processExisted(ctx, dao, *existed)
}

func (s *Service) processNew(ctx context.Context, dao Dao) error {
	if err := s.repo.Create(dao); err != nil {
		return fmt.Errorf("can't create dao: %w", err)
	}

	defer func(id uuid.UUID) {
		if err := s.repo.UpdateProposalCnt(id); err != nil {
			log.Warn().Err(err).Msgf("repo.UpdateProposalCnt: %s", id.String())
		}
		if err := s.repo.UpdateActiveVotes(id); err != nil {
			log.Warn().Err(err).Msgf("repo.UpdateActiveVotes: %s", id.String())
		}
	}(dao.ID)

	if err := s.events.PublishJSON(ctx, coreevents.SubjectDaoCreated, convertToCoreEvent(dao)); err != nil {
		log.Error().Err(err).Msgf("publish dao event #%s", dao.ID)
	}

	if err := s.events.PublishJSON(ctx, coreevents.SubjectCheckActivitySince, convertToCoreEvent(dao)); err != nil {
		log.Error().Err(err).Msgf("publish dao event #%s", dao.ID)
	}

	return nil
}

func (s *Service) processExisted(ctx context.Context, new, existed Dao) error {
	equal := compare(new, existed)
	if equal {
		return nil
	}

	new.CreatedAt = existed.CreatedAt
	new.ActivitySince = existed.ActivitySince
	new.PopularityIndex = existed.PopularityIndex
	new.Categories = enrichWithSystemCategories(new.Categories, existed.Categories)
	err := s.repo.Update(new)
	if err != nil {
		return fmt.Errorf("update dao #%s: %w", new.ID, err)
	}

	defer func(id uuid.UUID) {
		if err = s.repo.UpdateProposalCnt(id); err != nil {
			log.Warn().Err(err).Msgf("repo.UpdateProposalCnt: %s", id.String())
		}
		if err = s.repo.UpdateActiveVotes(id); err != nil {
			log.Warn().Err(err).Msgf("repo.UpdateActiveVotes: %s", id.String())
		}
	}(new.ID)

	go func(dao Dao) {
		if err := s.events.PublishJSON(ctx, coreevents.SubjectDaoUpdated, convertToCoreEvent(dao)); err != nil {
			log.Error().Err(err).Msgf("publish dao event #%s", dao.ID)
		}

		if err := s.events.PublishJSON(ctx, coreevents.SubjectCheckActivitySince, convertToCoreEvent(dao)); err != nil {
			log.Error().Err(err).Msgf("publish dao event #%s", dao.ID)
		}
	}(new)

	return nil
}

func enrichWithSystemCategories(list, existed []string) []string {
	for _, category := range systemCategories {
		if !slices.Contains(list, category) &&
			slices.Contains(existed, category) {
			list = append(list, category)
		}
	}

	return list
}

func compare(d1, d2 Dao) bool {
	d1.CreatedAt = d2.CreatedAt
	d1.UpdatedAt = d2.UpdatedAt
	d1.ActivitySince = d2.ActivitySince
	d1.PopularityIndex = d2.PopularityIndex

	return reflect.DeepEqual(d1, d2)
}

func (s *Service) GetByID(id uuid.UUID) (*Dao, error) {
	dao, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}

	return dao, nil
}

func (s *Service) GetIDByOriginalID(id string) (uuid.UUID, error) {
	s.daoMu.RLock()
	val, ok := s.daoIds[id]
	s.daoMu.RUnlock()
	if ok {
		return val, nil
	}

	val, err := s.idProvider.GetOrCreate(id)
	if err != nil {
		return val, fmt.Errorf("get or create: %w", err)
	}

	s.daoMu.Lock()
	s.daoIds[id] = val
	s.daoMu.Unlock()

	return val, nil
}

func (s *Service) GetByFilters(filters []Filter) (DaoList, error) {
	list, err := s.repo.GetByFilters(filters, true)
	if err != nil {
		return DaoList{}, fmt.Errorf("get by filters: %w", err)
	}

	return list, nil
}

type topList struct {
	List  []Dao
	Total int64
}

// todo: use caching here!
func (s *Service) GetTopByCategories(_ context.Context, limit int) (map[string]topList, error) {
	categories, err := s.repo.GetCategories()
	if err != nil {
		return nil, fmt.Errorf("get categories: %w", err)
	}

	list := make(map[string]topList)
	for _, category := range categories {
		filters := []Filter{
			CategoryFilter{Category: category},
			PageFilter{Limit: limit, Offset: 0},
			OrderByPopularityIndexFilter{},
		}

		data, err := s.repo.GetByFilters(filters, true)
		if err != nil {
			return nil, fmt.Errorf("get by category %s: %w", category, err)
		}

		list[category] = topList{
			List:  data.Daos,
			Total: data.TotalCount,
		}
	}

	return list, nil
}

func (s *Service) HandleActivitySince(_ context.Context, id uuid.UUID) (*Dao, error) {
	dao, err := s.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Warn().Str("dao_id", id.String()).Msg("dao is not ready yet")
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting dao by id: %w", err)
	}

	pr, err := s.proposals.GetEarliestByDaoID(dao.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("get earliest proposal: %w", err)
	}

	if pr == nil {
		return nil, nil
	}

	if dao.ActivitySince != 0 && pr.Created > dao.ActivitySince {
		return nil, nil
	}

	if dao.ActivitySince == 0 {
		dao.ActivitySince = pr.Created
	}

	if dao.ActivitySince > pr.Created {
		dao.ActivitySince = pr.Created
	}

	return dao, s.repo.Update(*dao)
}

// todo: add transaction here to avoid concurrent update
func (s *Service) processNewCategory(_ context.Context) error {
	list, err := s.repo.GetByFilters([]Filter{
		NotCategoryFilter{Category: newDaoCategoryName},
		ActivitySinceRangeFilter{From: time.Now().Add(-30 * 24 * time.Hour)},
		PageFilter{Limit: 300},
	}, false)
	if err != nil {
		return fmt.Errorf("get by filters: %w", err)
	}

	for i := range list.Daos {
		dao := list.Daos[i]
		dao.Categories = append(dao.Categories, newDaoCategoryName)

		if err = s.repo.Update(dao); err != nil {
			return fmt.Errorf("update dao: %s: %w", dao.ID.String(), err)
		}
	}

	return nil
}

// todo: add transaction here to avoid concurrent update
func (s *Service) processOutdatedNewCategory(_ context.Context) error {
	list, err := s.repo.GetByFilters([]Filter{
		CategoryFilter{Category: newDaoCategoryName},
		ActivitySinceRangeFilter{To: time.Now().Add(-31 * 24 * time.Hour)},
	}, false)
	if err != nil {
		return fmt.Errorf("get by filters: %w", err)
	}

	for i := range list.Daos {
		dao := list.Daos[i]

		dao.Categories = remove(dao.Categories, newDaoCategoryName)

		if err = s.repo.Update(dao); err != nil {
			return fmt.Errorf("update dao: %s: %w", dao.ID.String(), err)
		}
	}

	return nil
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func (s *Service) ProcessUniqueVoters(_ context.Context, voters []UniqueVoter) error {
	err := s.uniqueRepo.BatchCreate(voters)
	if err != nil {
		return fmt.Errorf("batchCreate: %w", err)
	}

	return nil
}

func (s *Service) processNewVoters(_ context.Context) error {
	err := s.uniqueRepo.UpdateVotersCount()
	if err != nil {
		return fmt.Errorf("UpdateVotersCount: %w", err)
	}

	return nil
}

func (s *Service) ProcessNewProposal(_ context.Context, originalDaoID string) error {
	daoID, err := s.idProvider.GetOrCreate(originalDaoID)
	if err != nil {
		return fmt.Errorf("idProvider.GetOrCreate: %w", err)
	}

	if err = s.repo.UpdateProposalCnt(daoID); err != nil {
		return fmt.Errorf("UpdateProposalCnt: %w", err)
	}

	if err = s.repo.UpdateActiveVotes(daoID); err != nil {
		return fmt.Errorf("UpdateActiveVotes: %w", err)
	}

	return nil
}

func (s *Service) ProcessExistedProposal(_ context.Context, originalDaoID string) error {
	daoID, err := s.idProvider.GetOrCreate(originalDaoID)
	if err != nil {
		return fmt.Errorf("idProvider.GetOrCreate: %w", err)
	}

	if err = s.repo.UpdateActiveVotes(daoID); err != nil {
		return fmt.Errorf("UpdateActiveVotes: %w", err)
	}

	return nil
}

func (s *Service) processPopularCategory(_ context.Context) error {
	listCurrent, err := s.repo.GetByFilters([]Filter{
		CategoryFilter{Category: popularDaoCategoryName},
	}, false)

	if err != nil {
		return fmt.Errorf("get by filters: %w", err)
	}

	for i := range listCurrent.Daos {
		dao := listCurrent.Daos[i]

		dao.Categories = remove(dao.Categories, popularDaoCategoryName)

		if err = s.repo.Update(dao); err != nil {
			return fmt.Errorf("update dao: %s: %w", dao.ID.String(), err)
		}
	}

	listNew, err := s.repo.GetByFilters([]Filter{
		OrderByPopularityIndexFilter{},
		PageFilter{Limit: 100},
	}, false)
	if err != nil {
		return fmt.Errorf("get by filters: %w", err)
	}

	for i := range listNew.Daos {
		dao := listNew.Daos[i]
		dao.Categories = append(dao.Categories, popularDaoCategoryName)

		if err = s.repo.Update(dao); err != nil {
			return fmt.Errorf("update dao: %s: %w", dao.ID.String(), err)
		}
	}

	return nil
}

func (s *Service) ProcessPopularityIndexUpdate(_ context.Context, id uuid.UUID, index float64) error {
	existed, err := s.repo.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("handle: %w", err)
	}

	if existed != nil {
		existed.PopularityIndex = index
		err := s.repo.Update(*existed)
		if err != nil {
			return fmt.Errorf("update dao #%s: %w", id, err)
		}
	}

	return nil
}

func (s *Service) processActiveVotes(_ context.Context) error {
	err := s.repo.UpdateActiveVotesAll()
	if err != nil {
		return fmt.Errorf("UpdateActiveVotesAll: %w", err)
	}

	return nil
}
