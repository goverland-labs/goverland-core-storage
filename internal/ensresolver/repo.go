package ensresolver

import (
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	return r.db.
		Model(&EnsName{}).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "address"},
			},
			UpdateAll: true,
		}).
		CreateInBatches(data, 500).
		Error
}

func (r *Repo) GetByAddresses(addresses []string) ([]EnsName, error) {
	args := make([]string, 0, len(addresses))
	for _, addr := range addresses {
		args = append(args, strings.ToLower(addr))
	}

	db := r.db.Model(&EnsName{}).Where("lower(address) IN ?", args)
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
