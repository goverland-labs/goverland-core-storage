package events

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

func (r *Repo) Create(e RegisteredEvent) error {
	return r.db.Create(&e).Error
}

func (r *Repo) Update(e RegisteredEvent) error {
	return r.db.Save(&e).Error
}

func (r *Repo) GetByTypeAndEvent(id, t, event string) (*RegisteredEvent, error) {
	var re RegisteredEvent
	request := r.db.Where(RegisteredEvent{
		Type:   t,
		TypeID: id,
		Event:  event,
	}).First(&re)
	if err := request.Error; err != nil {
		return nil, fmt.Errorf("get registered event #%s_%s_%s: %w", t, id, event, err)
	}

	return &re, nil
}

func (r *Repo) GetLast(limit int) ([]*RegisteredEvent, error) {
	var res []*RegisteredEvent
	request := r.db.Order("id desc").Limit(limit).Find(&res)
	if err := request.Error; err != nil {
		return nil, fmt.Errorf("get last #%d: %w", limit, err)
	}

	return res, nil
}
