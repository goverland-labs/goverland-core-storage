package vote

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/goverland-labs/goverland-core-storage/internal/proposal"

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
		Columns: []clause.Column{
			{Name: "proposal_id"},
			{Name: "voter"},
		},
		UpdateAll: true,
	}).CreateInBatches(data, defaultBatchSize).Error
}

type List struct {
	Votes      []Vote
	TotalCount int64
	TotalVp    float32
}

func (r *Repo) GetByFilters(filters []Filter, limit int, offset int, firstVoter string) (List, error) {
	db := r.db.Model(&Vote{}).InnerJoins("inner join daos on daos.id = votes.dao_id")
	for _, f := range filters {
		if _, ok := f.(proposal.OrderFilter); ok {
			continue
		}
		db = f.Apply(db)
	}
	var totals Totals
	err := db.Select([]string{"count(*) as Votes", "sum(vp) as Vp"}).Scan(&totals).Error

	if err != nil {
		return List{}, err
	}

	var list []Vote
	if firstVoter == "" {
		db = r.db.Model(&Vote{}).InnerJoins("inner join daos on daos.id = votes.dao_id")
		filters = append(filters, PageFilter{Limit: limit, Offset: offset})
		for _, f := range filters {
			db = f.Apply(db)
		}

		err = db.Find(&list).Error
		if err != nil {
			return List{}, err
		}
	} else {
		db = r.db.Model(&Vote{}).InnerJoins("inner join daos on daos.id = votes.dao_id")
		for _, f := range filters {
			db = f.Apply(db)
		}
		db = VoterFilter{Voter: firstVoter}.Apply(db)
		var v Vote
		request := db.First(&v)
		if err := request.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return List{}, err
		}

		db = r.db.Model(&Vote{}).InnerJoins("inner join daos on daos.id = votes.dao_id")
		prepend := false
		if request.RowsAffected > 0 {
			filters = append(filters, ExcludeVoterFilter{Voter: firstVoter})
			if offset == 0 {
				limit = limit - 1
				prepend = true
			} else {
				offset = offset - 1
			}
		}
		filters = append(filters, PageFilter{Limit: limit, Offset: offset})

		for _, f := range filters {
			db = f.Apply(db)
		}

		err = db.Find(&list).Error
		if err != nil {
			return List{}, err
		}
		if prepend {
			list = append([]Vote{v}, list...)
		}
	}

	return List{
		Votes:      list,
		TotalCount: totals.Votes,
		TotalVp:    totals.Vp,
	}, nil
}

func (r *Repo) GetLastItems(lastUpdatedAt time.Time, limit int) ([]Vote, error) {
	var list []Vote

	err := r.db.
		Where("updated_at > ?", lastUpdatedAt).
		Order("updated_at asc").
		Limit(limit).
		Find(&list).
		Error

	return list, err
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

func (r *Repo) GetByVoter(voter string) ([]string, error) {
	var res []string
	request := r.db.Model(&Vote{}).InnerJoins("inner join daos on daos.id = votes.dao_id").Distinct("dao_id").Where(&Vote{Voter: voter}).Find(&res)
	if err := request.Error; err != nil {
		return nil, fmt.Errorf("get dao list by voter #%s: %w", voter, err)
	}

	return res, nil
}
