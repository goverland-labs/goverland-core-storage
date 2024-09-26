package delegate

import (
	"gorm.io/gorm"
)

type Filter interface {
	Apply(*gorm.DB) *gorm.DB
}

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
