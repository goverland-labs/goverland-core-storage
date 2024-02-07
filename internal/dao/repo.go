package dao

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

// Create creates one dao object
// todo: check creating error/unique and others
func (r *Repo) Create(dao Dao) error {
	return r.db.Create(&dao).Error
}

// Update single dao object in database
// todo: think about updating fields to default value(boolean, string etc)
func (r *Repo) Update(dao Dao) error {
	return r.db.Omit("name", "original_id").Save(&dao).Error
}

func (r *Repo) GetByID(id uuid.UUID) (*Dao, error) {
	dao := Dao{ID: id}
	request := r.db.Take(&dao)
	if err := request.Error; err != nil {
		return nil, fmt.Errorf("get dao by id #%s: %w", id, err)
	}

	return &dao, nil
}

func (r *Repo) GetByName(name string) (*Dao, error) {
	var dao Dao
	request := r.db.Where(&Dao{Name: name}).First(&dao)
	if err := request.Error; err != nil {
		return nil, fmt.Errorf("get dao by name #%s: %w", name, err)
	}

	return &dao, nil
}

func (r *Repo) GetByOriginalID(id string) (*Dao, error) {
	var dao Dao
	request := r.db.Where(&Dao{OriginalID: id}).First(&dao)
	if err := request.Error; err != nil {
		return nil, fmt.Errorf("get dao by original id #%s: %w", id, err)
	}

	return &dao, nil
}

type DaoList struct {
	Daos       []Dao
	TotalCount int64
}

// todo: add order by
func (r *Repo) GetByFilters(filters []Filter, count bool) (DaoList, error) {
	db := r.db.Model(&Dao{})
	for _, f := range filters {
		if _, ok := f.(PageFilter); ok {
			continue
		}
		db = f.Apply(db)
	}

	var cnt int64
	err := db.Count(&cnt).Error
	if err != nil {
		return DaoList{}, err
	}

	for _, f := range filters {
		if _, ok := f.(PageFilter); ok {
			db = f.Apply(db)
		}
	}

	var daos []Dao
	err = db.Find(&daos).Error
	if err != nil {
		return DaoList{}, err
	}

	return DaoList{
		Daos:       daos,
		TotalCount: cnt,
	}, nil
}

func (r *Repo) GetCategories() ([]string, error) {
	var res []string
	err := r.db.Raw(`SELECT distinct JSONB_ARRAY_ELEMENTS_TEXT(categories) AS category FROM daos ORDER BY category`).Scan(&res).Error

	return res, err
}

func (r *Repo) UpdateProposalCnt(id uuid.UUID) error {
	return r.db.Exec(`
update daos
set proposals_count = cnt.proposals_count
from (
	select count(*) as proposals_count
	from proposals
	where dao_id = ?
) cnt
where daos.id = ?
`, id, id).Error
}

func (r *Repo) UpdateActiveVotes(id uuid.UUID) error {
	return r.db.Exec(`
update daos
set active_votes = cnt.active_votes
from (
	select count(*) as active_votes
	from proposals
	where dao_id = ? and state = 'active'
) cnt
where daos.id = ?
`, id, id).Error
}

func (r *Repo) UpdateActiveVotesAll() error {
	return r.db.Exec(`
update daos
set active_votes = cnt.active_votes
from (
	select dao_id, count(id) filter (where state = 'active') as active_votes
	from proposals
	group by dao_id
) cnt
where daos.id = cnt.dao_id
`).Error
}
