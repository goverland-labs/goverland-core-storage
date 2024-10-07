package ensresolver

import (
	"time"

	"github.com/rs/zerolog/log"
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

func (r *Repo) BatchCreate(data []EnsName) []EnsName {
	var created []EnsName
	for i := range data {
		target := data[i]
		target.UpdatedAt = time.Now()

		err := r.db.Model(&EnsName{}).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "name"}},
			DoUpdates: clause.AssignmentColumns([]string{"address"}),
		}).Create(&target).Error
		if err != nil {
			log.Error().
				Err(err).
				Str("name", target.Name).
				Str("address", target.Address).
				Msg("create ens name in db")

			continue
		}

		created = append(created, target)
	}

	return created
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
