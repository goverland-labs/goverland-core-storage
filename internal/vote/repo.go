package vote

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	defaultBatchSize = 1000
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

// BatchCreate creates votes in batch
func (r *Repo) BatchCreate(data []Vote) error {
	return r.db.Model(&Vote{}).Clauses(clause.OnConflict{
		DoNothing: true,
	}).CreateInBatches(data, defaultBatchSize).Error
}

type List struct {
	Votes      []Vote
	TotalCount int64
}

func (r *Repo) GetByFilters(filters []Filter) (List, error) {
	db := r.db.Model(&Vote{})
	for _, f := range filters {
		if _, ok := f.(PageFilter); ok {
			continue
		}
		db = f.Apply(db)
	}

	var cnt int64
	err := db.Count(&cnt).Error
	if err != nil {
		return List{}, err
	}

	for _, f := range filters {
		if _, ok := f.(PageFilter); ok {
			db = f.Apply(db)
		}
	}

	var list []Vote
	err = db.Find(&list).Error
	if err != nil {
		return List{}, err
	}

	return List{
		Votes:      list,
		TotalCount: cnt,
	}, nil
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
			update votes
			set ens_name = rs.name
			from (
				values %s
			) as rs (address, name)
			where rs.address = votes.voter
`, strings.Join(placeholders, ","))

	return r.db.Exec(query, args...).Error
}

func (r *Repo) GetUnique(cursor string, limit int64) ([]string, error) {
	if limit == 0 {
		return nil, nil
	}

	query := `
select distinct dao_voter.voter author
from dao_voter
where dao_voter.voter > ?
order by author
limit ?
`

	var list []string
	err := r.db.Debug().Raw(query, cursor, limit).Scan(&list).Error

	return list, err
}
