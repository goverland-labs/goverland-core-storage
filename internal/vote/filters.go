package vote

import (
	"gorm.io/gorm"
)

type Filter interface {
	Apply(*gorm.DB) *gorm.DB
}

type PageFilter struct {
	Offset int
	Limit  int
}

func (f PageFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Offset(f.Offset).Limit(f.Limit)
}

type ProposalFilter struct {
	ProposalID string
}

func (f ProposalFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("proposal_id = ?", f.ProposalID)
}

type OrderByCreatedFilter struct {
}

func (f OrderByCreatedFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Order("created desc")
}
