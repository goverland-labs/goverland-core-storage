package dao

import (
	"fmt"

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
	return r.db.Save(&dao).Error
}

func (r *Repo) GetByID(id string) (*Dao, error) {
	dao := Dao{ID: id}
	request := r.db.Take(&dao)
	if err := request.Error; err != nil {
		return nil, fmt.Errorf("get dao by id #%s: %w", id, err)
	}

	return &dao, nil
}
