package delegate

import (
	"gorm.io/gorm"
)

type Filter interface {
	Apply(*gorm.DB) *gorm.DB
}

// DelegatorFilter Who delegate voting power
type DelegatorFilter struct {
	Address string
}

func (f DelegatorFilter) Apply(db *gorm.DB) *gorm.DB {
	var (
		dummy = Summary{}
		_     = dummy.AddressFrom
	)

	return db.Where("lower(address_from) = lower(?)", f.Address)
}

// DelegationFilter Whom delegate voting power
type DelegationFilter struct {
	Address string
}

func (f DelegationFilter) Apply(db *gorm.DB) *gorm.DB {
	var (
		dummy = Summary{}
		_     = dummy.AddressTo
	)

	return db.Where("lower(address_to) = lower(?)", f.Address)
}
