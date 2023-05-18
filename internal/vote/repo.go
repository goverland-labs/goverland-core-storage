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
	return r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(data, defaultBatchSize).Error
}
