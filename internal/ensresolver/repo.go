package ensresolver

import (
	"time"

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

type EnsName struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Address   string
	Name      string
}

func (r *Repo) BatchCreate(data []EnsName) error {
	return r.db.Model(&EnsName{}).Clauses(clause.OnConflict{
		DoNothing: true,
	}).CreateInBatches(data, defaultBatchSize).Error
}

func (r *Repo) GetByAddresses(addresses []string) ([]EnsName, error) {
	db := r.db.Model(&EnsName{}).Where("address IN ?", addresses)
	var list []EnsName
	err := db.Find(&list).Error

	return list, err
}

func (r *Repo) GetByNames(names []string) ([]EnsName, error) {
	db := r.db.Model(&EnsName{}).Where("name IN ?", names)
	var list []EnsName
	err := db.Find(&list).Error

	return list, err
}
