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

type ProposalIDsFilter struct {
	ProposalIDs []string
}

func (f ProposalIDsFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("proposal_id IN ?", f.ProposalIDs)
}

type VoterFilter struct {
	Voter string
}

func (f VoterFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("voter = ?", f.Voter)
}

type OrderByCreatedFilter struct {
}

func (f OrderByCreatedFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Order("created desc")
}
