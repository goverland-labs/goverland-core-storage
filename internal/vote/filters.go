package vote

import (
	"github.com/goverland-labs/goverland-core-storage/internal/proposal"
	"gorm.io/gorm"
)

var (
	OrderByVp = proposal.Order{
		Field:     "vp",
		Direction: proposal.DirectionDesc,
	}
	OrderByCreated = proposal.Order{
		Field:     "created",
		Direction: proposal.DirectionDesc,
	}
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

type ExcludeVoterFilter struct {
	Voter string
}

func (f ExcludeVoterFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("voter != ?", f.Voter)
}

type QueryFilter struct {
	Query string
}

func (f QueryFilter) Apply(db *gorm.DB) *gorm.DB {
	return db.Where("voter = ? or ens_name = ?", f.Query, f.Query)
}
