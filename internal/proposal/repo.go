package proposal

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

// Create creates one proposal object
// todo: check creating error/unique and others
func (r *Repo) Create(p Proposal) error {
	return r.db.Create(&p).Error
}

// Update single proposal object in database
// todo: think about updating fields to default value(boolean, string etc)
func (r *Repo) Update(p Proposal) error {
	return r.db.Save(&p).Error
}

func (r *Repo) GetByID(id string) (*Proposal, error) {
	p := Proposal{ID: id}
	request := r.db.Take(&p)
	if err := request.Error; err != nil {
		return nil, fmt.Errorf("get proposal by id #%s: %w", id, err)
	}

	return &p, nil
}

// todo: think about limits, add pagination or cursor
func (r *Repo) GetAvailableForVoting(window time.Duration) ([]*Proposal, error) {
	var items []*Proposal
	err := r.db.Raw("select * from proposals p where p.end > ?", time.Now().Add(-window).Unix()).Find(&items).Error
	if err != nil {
		return nil, fmt.Errorf("find active: %w", err)
	}

	return items, nil
}
