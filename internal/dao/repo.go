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

func (r *Repo) GetCountByFilters(filters []Filter) (int64, error) {
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
		return 0, fmt.Errorf("db.Count: %w", err)
	}

	return cnt, nil
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
	where dao_id = ? and spam is not true and state != 'canceled'
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
	where dao_id = ? and state = 'active' and spam is not true
) cnt
where daos.id = ?
`, id, id).Error
}

func (r *Repo) UpdateActiveVotesAll() error {
	return r.db.Exec(`
update daos
set active_votes = cnt.active_votes
from (
	select dao_id, count(id) filter (where state = 'active' and spam is not true) as active_votes
	from proposals
	group by dao_id
) cnt
where daos.id = cnt.dao_id
`).Error
}

// GetRecommended returns the list of available dao strategies in our system
func (r *Repo) GetRecommended() ([]Recommendation, error) {
	query := `
select original_id,
       uuid,
       name                 strategy_name,
       params ->> 'symbol'  symbol,
       network              network_id,
       params ->> 'address' address
from (select daos.original_id,
             daos.id        uuid,
             daos.network,
             st."Name"   as name,
             st."Params" as params
      from daos,
           jsonb_to_recordset(daos.strategies) AS st("Name" text, "Params" jsonb)
      where st."Name" <> 'multichain'
        and verified is true) data
where name in
      ('erc20-votes', 'erc20-balance-of',
       'eth-balance', 'erc721', 'eth-with-balance',
       'contract-call', 'erc1155-balance-of', 'ens-domains-owned')
  and params ->> 'symbol' is not null
  and params ->> 'address' is not null

union all

select data.original_id,
       data.uuid,
       st.name as              strategy_name,
       st.params ->> 'symbol'  symbol,
       data.network            network_id,
       st.params ->> 'address' address
from (select daos.original_id,
             daos.id        uuid,
             daos.network,
             st."Name"   as name,
             st."Params" as params
      from daos,
           jsonb_to_recordset(daos.strategies) AS st("Name" text, "Params" jsonb)
      where st."Name" = 'multichain'
        and verified is true) data,
     jsonb_to_recordset(data.params -> 'strategies') AS st("name" text, "params" jsonb)
where st.params ->> 'symbol' is not null
  and st.params ->> 'address' is not null

order by original_id, symbol`

	rows, err := r.db.Raw(query).Rows()
	if err != nil {
		return nil, fmt.Errorf("get active by user: %w", err)
	}

	defer rows.Close()

	defaultSize := 500
	list := make([]Recommendation, 0, defaultSize)
	for rows.Next() {
		rec := Recommendation{}

		err = rows.Scan(
			&rec.OriginalId,
			&rec.InternalId,
			&rec.Name,
			&rec.Symbol,
			&rec.NetworkId,
			&rec.Address,
		)
		if err != nil {
			return nil, fmt.Errorf("convert row: %w", err)
		}

		list = append(list, rec)
	}

	return list, nil
}
