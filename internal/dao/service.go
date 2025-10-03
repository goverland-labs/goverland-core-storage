package dao

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"sort"
	"sync"
	"time"

	"golang.org/x/text/language"
	message2 "golang.org/x/text/message"

	"github.com/goverland-labs/goverland-core-storage/internal/discord"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	coreevents "github.com/goverland-labs/goverland-platform-events/events/core"

	"github.com/goverland-labs/goverland-core-storage/internal/proposal"
	"github.com/goverland-labs/goverland-core-storage/pkg/sdk/zerion"
)

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
	GetByOriginalID(id string) (*Dao, error)
	UpdateProposalCnt(id uuid.UUID) error
	UpdateActiveVotes(id uuid.UUID) error
	UpdateActiveVotesAll() error
	UpdateProposalCntAll() error
	GetByFilters(filters []Filter, count bool) (DaoList, error)
	GetCategories() ([]string, error)
	GetRecommended() ([]Recommendation, error)
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
	GetLatestByDaoID(uuid.UUID) (*proposal.Proposal, error)
	GetByID(string) (*proposal.Proposal, error)
}

type cachedDAO struct {
	value     *Dao
	expiresAt time.Time
}

type Service struct {
	daoIds map[string]uuid.UUID
	daoMu  sync.RWMutex

	spacesByOriginalID map[string]cachedDAO
	spacesLock         sync.RWMutex

	recommendations   []Recommendation
	recommendationsMu sync.RWMutex

	repo              DataProvider
	fungibleChainRepo *FungibleChainRepo
	uniqueRepo        UniqueVoterProvider
	events            Publisher
	idProvider        DaoIDProvider
	proposals         ProposalProvider

	topDAOCache   *TopDAOCache
	zerionClient  *zerion.Client
	discordSender *discord.Sender
}

func NewService(r DataProvider, ur UniqueVoterProvider, ip DaoIDProvider, p Publisher, pp ProposalProvider, topDAOCache *TopDAOCache, fungibleChainRepo *FungibleChainRepo, zerionClient *zerion.Client, discordSender *discord.Sender) (*Service, error) {
	return &Service{
		repo:               r,
		uniqueRepo:         ur,
		events:             p,
		idProvider:         ip,
		proposals:          pp,
		daoIds:             make(map[string]uuid.UUID),
		daoMu:              sync.RWMutex{},
		spacesByOriginalID: make(map[string]cachedDAO),
		spacesLock:         sync.RWMutex{},
		topDAOCache:        topDAOCache,
		fungibleChainRepo:  fungibleChainRepo,
		zerionClient:       zerionClient,
		discordSender:      discordSender,
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
	new.FungibleId = existed.FungibleId
	new.TokenSymbol = existed.TokenSymbol
	new.VerificationStatus = existed.VerificationStatus
	new.VerificationComment = existed.VerificationComment
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

func (s *Service) getFungibleId(strategies Strategies) (string, string) {
	for _, strategy := range strategies {
		if strategy.Name != "erc20-balance-of" && strategy.Name != "erc20-votes" && strategy.Name != "erc20-balance-of-delegation" && strategy.Name != "comp-like-votes" {
			continue
		}
		adr := strategy.Params["address"].(string)
		if adr == "" {
			continue
		}
		l, err := s.zerionClient.GetFungibleList("", adr)
		if err == nil && l != nil && len(l.List) == 1 {
			data := l.List[0]
			if data.Attributes.MarketData.Price != 0 {
				return data.ID, data.Attributes.Symbol
			}
		}
	}

	return "", ""
}

func compare(d1, d2 Dao) bool {
	d1.CreatedAt = d2.CreatedAt
	d1.UpdatedAt = d2.UpdatedAt
	d1.ActivitySince = d2.ActivitySince
	d1.PopularityIndex = d2.PopularityIndex
	d1.FungibleId = d2.FungibleId
	d1.TokenSymbol = d2.TokenSymbol
	d1.VerificationStatus = d2.VerificationStatus
	d1.VerificationComment = d2.VerificationComment

	return reflect.DeepEqual(d1, d2)
}

func (s *Service) GetByID(id uuid.UUID) (*Dao, error) {
	dao, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}

	return dao, nil
}

func (s *Service) GetDaoByOriginalID(id string) (*Dao, error) {
	s.spacesLock.RLock()
	data, ok := s.spacesByOriginalID[id]
	s.spacesLock.RUnlock()
	if ok && data.expiresAt.After(time.Now()) {
		return data.value, nil
	}

	dao, err := s.repo.GetByOriginalID(id)
	if err != nil {
		return nil, fmt.Errorf("get by original id: %w", err)
	}

	// simple caching for 5 minutes
	s.spacesLock.Lock()
	s.spacesByOriginalID[id] = cachedDAO{
		value:     dao,
		expiresAt: time.Now().Add(time.Minute * 5),
	}
	s.spacesLock.Unlock()

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

func (s *Service) GetTopByCategories(_ context.Context, limit int) (map[string]topList, error) {
	return s.topDAOCache.GetTopList(uint(limit)), nil
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

func (s *Service) UpdateFungibleId(_ context.Context, id uuid.UUID) {
	dao, err := s.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Warn().Str("dao_id", id.String()).Msg("dao is not ready yet")
		return
	}
	if err != nil {
		return
	}

	if dao.FungibleId != "" || dao.VerificationStatus == "declined" {
		return
	}

	fi, ts := s.getFungibleId(dao.Strategies)
	if fi == "" {
		return

	}

	dao.VerificationStatus = "pending"
	dao.FungibleId = fi
	dao.TokenSymbol = ts
	_ = s.repo.Update(*dao)
	s.sendDaoToDiscord(*dao)
}

// todo: add transaction here to avoid concurrent update
func (s *Service) processNewCategory(_ context.Context) error {
	list, err := s.repo.GetByFilters([]Filter{
		NotCategoryFilter{Category: newDaoCategoryName},
		ActivitySinceRangeFilter{From: time.Now().Add(-90 * 24 * time.Hour)},
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
		ActivitySinceRangeFilter{To: time.Now().Add(-91 * 24 * time.Hour)},
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

func (s *Service) ProcessDeletedProposal(_ context.Context, daoID uuid.UUID) error {
	if err := s.repo.UpdateProposalCnt(daoID); err != nil {
		return fmt.Errorf("UpdateProposalCnt: %w", err)
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

func (s *Service) processProposalsCnt(_ context.Context) error {
	err := s.repo.UpdateProposalCntAll()
	if err != nil {
		return fmt.Errorf("UpdateProposalCntAll: %w", err)
	}

	return nil
}

func (s *Service) getRecommendations() []Recommendation {
	s.recommendationsMu.RLock()
	data := make([]Recommendation, len(s.recommendations))
	copy(data, s.recommendations)
	s.recommendationsMu.RUnlock()

	return data
}

func (s *Service) syncRecommendations(_ context.Context) error {
	list, err := s.repo.GetRecommended()
	if err != nil {
		return fmt.Errorf("getRecommended: %w", err)
	}

	s.recommendationsMu.Lock()
	s.recommendations = list
	s.recommendationsMu.Unlock()

	return nil
}

func (s *Service) GetTokenInfo(id uuid.UUID) (*TokenInfo, error) {
	dao, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get dao: %w", err)
	}

	if dao.FungibleId == "" {
		return nil, status.Error(codes.Internal, "can't receive the information from Zerion")
	}

	data, err := s.zerionClient.GetFungibleData(dao.FungibleId)
	if err != nil {
		return nil, fmt.Errorf("failed to get token info: %w", err)
	}

	var chains []TokenChainInfo
	if dao.FungibleId != "" {
		fChains, err := s.fungibleChainRepo.GetByFungibleID(dao.FungibleId)
		if err != nil {
			return nil, fmt.Errorf("failed to get chains: %w", err)
		}

		for _, chain := range fChains {
			chains = append(chains, TokenChainInfo{
				ChainID:  chain.ChainID,
				Name:     chain.ChainName,
				Decimals: chain.Decimals,
				IconURL:  chain.IconURL,
				Address:  chain.Address,
			})
		}
	}

	return &TokenInfo{
		Name:                  data.Attributes.Name,
		Symbol:                data.Attributes.Symbol,
		TotalSupply:           data.Attributes.MarketData.TotalSupply,
		CirculatingSupply:     data.Attributes.MarketData.CirculatingSupply,
		MarketCap:             data.Attributes.MarketData.MarketCap,
		FullyDilutedValuation: data.Attributes.MarketData.FullyDilutedValuation,
		Price:                 data.Attributes.MarketData.Price,
		FungibleID:            dao.FungibleId,
		Chains:                chains,
	}, nil
}

func (s *Service) GetTokenPrice(id uuid.UUID, created int) float64 {
	data, err := s.GetTokenChart(id, "hour")
	if err != nil || data == nil {
		return 0
	}

	points := data.ChartAttributes.Points
	if len(points) == 0 {
		return 0
	}

	createdTime := time.Unix(int64(created), 0)

	idx := sort.Search(len(points), func(i int) bool {
		return points[i].Time.After(createdTime)
	})

	if idx == 0 {
		// All points are after created time
		return 0
	}

	return points[idx-1].Price
}

func (s *Service) GetTokenChart(id uuid.UUID, period string) (*zerion.ChartData, error) {
	dao, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get dao: %w", err)
	}

	if dao.FungibleId == "" {
		return nil, status.Error(codes.Internal, "can't receive the information from Zerion")
	}

	data, err := s.zerionClient.GetFungibleChart(dao.FungibleId, period)
	if err != nil {
		return nil, fmt.Errorf("failed to get token chart: %w", err)
	}

	return data, nil
}

func (s *Service) PopulateTokenPrices(ctx context.Context, id uuid.UUID) (bool, error) {
	data, err := s.GetTokenChart(id, "month")
	if err != nil || data == nil {
		return false, fmt.Errorf("failed to get token prices: %w", err)
	}
	if err := s.events.PublishJSON(ctx, coreevents.DaoTokenPriceUpdated, convertToTokenPricesPayload(data.ChartAttributes.Points, id)); err != nil {
		return false, fmt.Errorf("publish token prices event: %w", err)
	}
	return true, nil
}

func convertToTokenPricesPayload(list []zerion.Point, daoId uuid.UUID) coreevents.TokenPricesPayload {
	res := make(coreevents.TokenPricesPayload, 0, len(list))
	for _, point := range list {
		res = append(res, coreevents.TokenPricePayload{
			DaoID: daoId,
			Time:  point.Time,
			Price: point.Price,
		})
	}

	return res
}

func (s *Service) UpdateFungibleIds(_ context.Context, category string) (bool, error) {
	filters := []Filter{
		FungibleIdEmptyFilter{},
	}
	if category == "new" {
		filters = append(filters, CategoryFilter{Category: newDaoCategoryName})
	} else {
		filters = append(filters, VerifiedFilter{})
	}
	daos, err := s.repo.GetByFilters(filters, false)
	if err != nil {
		return false, fmt.Errorf("get daos: %w", err)
	}
	for _, dao := range daos.Daos {
		if dao.VerificationStatus != "declined" {
			fi, ts := s.getFungibleId(dao.Strategies)
			if fi != "" {
				dao.VerificationStatus = "pending"
				dao.FungibleId = fi
				dao.TokenSymbol = ts
				_ = s.repo.Update(dao)
				s.sendDaoToDiscord(dao)
			}
		}
	}

	return true, nil
}

func (s *Service) sendDaoToDiscord(dao Dao) {
	pr, err := s.proposals.GetLatestByDaoID(dao.ID)
	vp := float32(0)
	if err == nil && pr != nil {
		vp = pr.ScoresTotal
	}
	price := float64(0)
	fdp := float64(0)
	data, err := s.zerionClient.GetFungibleData(dao.FungibleId)
	if err == nil && data != nil {
		if data.Attributes.MarketData.Price != 0 {
			price = data.Attributes.MarketData.Price
			fdp = data.Attributes.MarketData.FullyDilutedValuation
		}
	}
	pEN := message2.NewPrinter(language.English)
	message := pEN.Sprintf("DAO: https://gl.app/dao/%s\nCreated: %s\n FungibleId: %s\nToken %s: Price - %v, FDV - %v\nLast Proposal VP: %v",
		dao.OriginalID, dao.CreatedAt.Format(time.DateOnly), dao.FungibleId, dao.TokenSymbol, price, fdp, vp)
	if err := s.discordSender.SendMessage(message); err != nil {
		log.Warn().Err(err).Msgf("can't send discord message for dao: %s", dao.ID.String())
	}
}
