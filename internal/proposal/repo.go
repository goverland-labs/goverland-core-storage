package proposal

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

// Create creates one proposal object
// todo: check creating error/unique and others
func (r *Repo) Create(p Proposal) error {
	return r.db.Create(&p).Error
}

// Update single proposal object in database
// todo: think about updating fields to default value(boolean, string etc)
func (r *Repo) Update(p Proposal) error {
	return r.db.Save(&p).Error
}

func (r *Repo) GetByID(id string) (*Proposal, error) {
	var p Proposal
	request := r.db.Where(&Proposal{ID: id}).First(&p)
	if err := request.Error; err != nil {
		return nil, fmt.Errorf("get proposal by id #%s: %w", id, err)
	}

	return &p, nil
}

// todo: think about limits, add pagination or cursor
func (r *Repo) GetAvailableForVoting(window time.Duration) ([]*Proposal, error) {
	var items []*Proposal
	err := r.db.Raw("select * from proposals p where p.end > ?", time.Now().Add(-window).Unix()).Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("find active: %w", err)
	}

	return items, nil
}

func (r *Repo) GetEarliestByDaoID(daoID uuid.UUID) (*Proposal, error) {
	var pr *Proposal
	err := r.db.Raw("select * from proposals p where p.dao_id = ? order by created asc limit 1", daoID).First(&pr).Error
	if err != nil {
		return nil, fmt.Errorf("find active: %w", err)
	}

	return pr, nil
}

func (r *Repo) GetLatestByDaoID(daoID uuid.UUID) (*Proposal, error) {
	var pr *Proposal
	err := r.db.Raw("select * from proposals p where p.dao_id = ? order by created desc limit 1", daoID).First(&pr).Error
	if err != nil {
		return nil, fmt.Errorf("find active: %w", err)
	}

	return pr, nil
}

type ProposalList struct {
	Proposals  []Proposal
	TotalCount int64
}

// todo: add order by
func (r *Repo) GetByFilters(filters []Filter) (ProposalList, error) {
	db := r.db.Model(&Proposal{}).InnerJoins("inner join daos on daos.id = proposals.dao_id")
	for _, f := range filters {
		if _, ok := f.(PageFilter); ok {
			continue
		}
		db = f.Apply(db)
	}

	var cnt int64
	err := db.Count(&cnt).Error
	if err != nil {
		return ProposalList{}, err
	}

	return getProposalList(db, filters, cnt)
}

func (r *Repo) GetCountByFilters(filters []Filter) (int64, error) {
	db := r.db.
		Model(&Proposal{}).
		InnerJoins("inner join daos on daos.id = proposals.dao_id")

	for _, f := range filters {
		if _, ok := f.(PageFilter); ok {
			continue
		}

		db = f.Apply(db)
	}

	var cnt int64
	err := db.Count(&cnt).Error
	if err != nil {
		return 0, fmt.Errorf("db.Count: %w", err)
	}

	return cnt, nil
}

func (r *Repo) GetTop(filters []Filter) (ProposalList, error) {
	db := r.db.
		Model(&Proposal{}).
		InnerJoins("inner join (select id, row_number() over (partition by dao_id order by votes/(EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)-start) desc) as r " +
			"from proposals where state = 'active' and spam is not true and votes >= 30 and ((to_timestamp(proposals.start)>=now() -interval '14 days') or (now() >=to_timestamp(proposals.end) -interval '14 days'))) pr on pr.id = proposals.id and pr.r <= 3").
		InnerJoins("inner join daos on daos.id = proposals.dao_id").
		Order("daos.verified desc").
		Order("pr.r").
		Order("votes/(EXTRACT(EPOCH FROM CURRENT_TIMESTAMP)-start) desc")

	var (
		cnt   int64
		total int64
	)
	for _, f := range filters {
		if pf, ok := f.(PageFilter); ok {
			cnt = int64(pf.Limit)
			continue
		}
		db = f.Apply(db)
	}

	err := db.Count(&total).Error
	if err != nil {
		return ProposalList{}, err
	}

	if cnt > total {
		cnt = total
	}

	return getProposalList(db, filters, cnt)
}

func (r *Repo) UpdateVotes(list []ResolvedAddress) error {
	if len(list) == 0 {
		return nil
	}

	placeholders := make([]string, 0, len(list))
	args := make([]any, 0, 2*len(list))
	for i := range list {
		placeholders = append(placeholders, "(?,?)")
		args = append(args, list[i].Address, list[i].Name)
	}

	query := fmt.Sprintf(`
			update proposals
			set ens_name = rs.name
			from (
				values %s
			) as rs (address, name)
			where rs.address = proposals.author
`, strings.Join(placeholders, ","))

	return r.db.Exec(query, args...).Error
}

func getProposalList(db *gorm.DB, filters []Filter, cnt int64) (ProposalList, error) {
	for _, f := range filters {
		if _, ok := f.(PageFilter); ok {
			db = f.Apply(db)
		}
	}

	var list []Proposal
	if err := db.Find(&list).Error; err != nil {
		return ProposalList{}, err
	}

	return ProposalList{
		Proposals:  list,
		TotalCount: cnt,
	}, nil
}

func (r *Repo) GetSucceededChoices(daoId uuid.UUID) []string {
	var sc DaoSucceededChoices
	request := r.db.Where(&DaoSucceededChoices{DaoID: daoId}).First(&sc)
	if err := request.Error; err != nil {
		return nil
	}
	return sc.Choices
}
