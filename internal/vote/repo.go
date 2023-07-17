package vote

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	defaultBatchSize = 100
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
		UpdateAll: true,
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
